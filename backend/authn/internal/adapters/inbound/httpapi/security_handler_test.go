package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequireBearerAuth(t *testing.T) {
	t.Run("no auth scope - passes through", func(t *testing.T) {
		handler := RequireBearerAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		handler := RequireBearerAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandleBearerAuth(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		token := "token-123"
		newCtx, err := HandleBearerAuth(context.Background(), token)
		assert.NoError(t, err)
		assert.Equal(t, token, newCtx.Value(UserIDKey))
	})

	t.Run("missing token", func(t *testing.T) {
		_, err := HandleBearerAuth(context.Background(), "")
		assert.Error(t, err)
		assert.Equal(t, "missing bearer token", err.Error())
	})
}
