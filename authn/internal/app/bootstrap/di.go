package bootstrap

import (
	"context"
	"log/slog"

	"github.com/faber-numeris/beholder/authn/internal/adapters/outbound/postgres"
	sqlcgen "github.com/faber-numeris/beholder/authn/internal/adapters/outbound/postgres/gen"
	"github.com/faber-numeris/beholder/authn/internal/core/services"
	"github.com/faber-numeris/beholder/authn/internal/infrastructure/config"
	infrapostgres "github.com/faber-numeris/beholder/authn/internal/infrastructure/postgres"
	"github.com/faber-numeris/beholder/authn/internal/platform/util"
	servicesadapter "github.com/faber-numeris/beholder/authn/internal/ports/outbound"
)

func ProvideHashingService() servicesadapter.HashingService {
	return services.NewHashingService()
}

type repositoryComposition struct {
	servicesadapter.UserRepository
	servicesadapter.UserConfirmationRepository
}

func (r repositoryComposition) Ping() bool {
	return r.UserRepository.Ping()
}

func ProvideRepository(dbConfig config.IDatabaseConfig) servicesadapter.Repository {

	pool := util.Must(infrapostgres.Connect(context.Background(), dbConfig))
	_ = pool

	db := infrapostgres.GetDB()
	if db == nil {
		slog.Warn("Database not connected, repositories will return errors until connection is established")
	}

	// Inject the SQLC dependencies
	queries := sqlcgen.New(pool)
	return repositoryComposition{
		postgresadapter.NewUserRepository(queries),
		postgresadapter.NewUserConfirmationRepository(queries),
	}
}
