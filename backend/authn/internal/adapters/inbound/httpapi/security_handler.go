package httpapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/faber-numeris/beholder/backend/authn/internal/adapters/inbound/httpapi/gen"
)

type userContextKey string

const UserIDKey userContextKey = "user_id"

func RequireBearerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scopes := r.Context().Value(api.BearerAuthScopes)
		if scopes == nil {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error":   "UNAUTHORIZED",
				"message": "Missing authorization header",
			})
			return
		}

		token, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok || token == "" {
			slog.Warn("Empty or invalid token provided")
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error":   "UNAUTHORIZED",
				"message": "Missing bearer token",
			})
			return
		}

		slog.Debug("Token validated", "token", token)
		ctx := context.WithValue(r.Context(), UserIDKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func HandleBearerAuth(ctx context.Context, token string) (context.Context, error) {
	if token == "" {
		slog.Warn("Empty token provided")
		return ctx, errors.New("missing bearer token")
	}
	slog.Debug("Token validated", "token", token)
	ctx = context.WithValue(ctx, UserIDKey, token)
	return ctx, nil
}
