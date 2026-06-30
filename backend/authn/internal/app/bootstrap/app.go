package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpapi2 "github.com/faber-numeris/beholder/backend/authn/internal/adapters/inbound/httpapi"
	"github.com/faber-numeris/beholder/backend/authn/internal/adapters/inbound/httpapi/gen"
	"github.com/faber-numeris/beholder/backend/authn/internal/adapters/outbound/mail"
	"github.com/faber-numeris/beholder/backend/authn/internal/core/services"
	config2 "github.com/faber-numeris/beholder/backend/authn/internal/infrastructure/config"
	"github.com/faber-numeris/beholder/backend/authn/internal/infrastructure/postgres"
	"github.com/faber-numeris/beholder/backend/authn/internal/platform/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	specui "github.com/oaswrap/spec-ui"
)

type App struct {
	server        *http.Server
	db            *sqlx.DB
	healthChecker *httpapi2.HealthChecker
}

func NewApp() *App {
	cfg := util.Must(config2.NewConfig())

	healthChecker := httpapi2.NewHealthChecker()

	hashingSvc := ProvideHashingService()
	healthChecker.RegisterAdapter("hashing", hashingSvc)

	mailer := mail.NewService(cfg)
	healthChecker.RegisterAdapter("mail", mailer)

	repo := ProvideRepository(cfg)
	healthChecker.RegisterAdapter("database", repo)

	userSvc := services.NewUserService(repo, hashingSvc, mailer)

	handler := httpapi2.NewHandler(userSvc, hashingSvc, healthChecker)

	router := buildRouter(handler, cfg)

	address := fmt.Sprintf(":%d", cfg.Port())
	return &App{
		server: &http.Server{
			Addr:    address,
			Handler: router,
		},
		db:            postgres.GetDB(),
		healthChecker: healthChecker,
	}
}

func buildRouter(handler api.ServerInterface, cfg config2.IServiceConfig) http.Handler {
	specuiHandler := specui.NewHandler(
		specui.WithTitle("Beholder API"),
		specui.WithDocsPath("/docs/authn"),
		specui.WithSpecPath("/docs/authn/openapi.yaml"),
		specui.WithSpecFile("authn/internal/adapters/inbound/httpapi/openapi.yaml"),
		specui.WithStoplightElements(),
	)

	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(httpapi2.RequireBearerAuth)
	mux.Get(specuiHandler.DocsPath(), specuiHandler.DocsFunc())
	mux.Get(specuiHandler.SpecPath(), specuiHandler.SpecFunc())

	api.HandlerFromMuxWithBaseURL(handler, mux, "/v1")

	return mux
}

func (a *App) backgroundRetry(ctx context.Context) {
	baseDelay := 1 * time.Second
	maxDelay := 30 * time.Second

	for attempt := 0; ; attempt++ {
		if a.healthChecker.IsReady() {
			attempt = 0
			select {
			case <-ctx.Done():
				return
			case <-time.After(15 * time.Second):
			}
			continue
		}

		delay := time.Duration(math.Min(
			float64(baseDelay)*math.Pow(2, float64(attempt)),
			float64(maxDelay),
		))

		slog.Warn("Dependencies not ready, retrying",
			"attempt", attempt+1,
			"backoff", delay,
		)

		if a.db != nil {
			if err := a.db.Ping(); err != nil {
				slog.Warn("Database ping failed", "error", err)
			} else {
				slog.Info("Database connection established")
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
		}
	}
}

func (a *App) Run() error {
	slog.Info("Starting AuthN service", "address", a.server.Addr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go a.backgroundRetry(ctx)

	srvErrChan := make(chan error, 1)
	go func() {
		if srvErr := a.server.ListenAndServe(); srvErr != nil && !errors.Is(srvErr, http.ErrServerClosed) {
			srvErrChan <- srvErr
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srvErrChan:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
	case <-quit:
		slog.Info("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}

		if a.db != nil {
			if err := a.db.Close(); err != nil {
				slog.Error("failed to close database", "error", err)
			}
		}

		slog.Info("Server exited")
	}

	return nil
}
