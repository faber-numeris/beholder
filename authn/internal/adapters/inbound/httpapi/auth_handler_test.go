package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/faber-numeris/beholder/authn/internal/core/domain"
	"github.com/faber-numeris/beholder/authn/internal/mocks"
	"github.com/faber-numeris/foundation/beholder/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_ConfirmUserRegistration(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().ConfirmUserRegistration(ctx, "token-123").Return(nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/confirm?token=token-123", nil).WithContext(ctx)
		handler.ConfirmUserRegistration(w, r, api.ConfirmUserRegistrationParams{Token: "token-123"})

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().ConfirmUserRegistration(ctx, "bad-token").Return(errors.New("invalid"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/confirm?token=bad-token", nil).WithContext(ctx)
		handler.ConfirmUserRegistration(w, r, api.ConfirmUserRegistrationParams{Token: "bad-token"})

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_RegisterUser(t *testing.T) {
	ctx := context.Background()
	req := &api.UserCreateRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	t.Run("success", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().RegisterUser(ctx, mock.Anything, []byte("password")).
			Return(&domain.User{ID: "123", Email: "test@example.com"}, nil)

		body := toJSONBody(t, req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/register", body).WithContext(ctx)
		handler.RegisterUser(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().RegisterUser(ctx, mock.Anything, []byte("password")).
			Return(nil, errors.New("db error"))

		body := toJSONBody(t, req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/register", body).WithContext(ctx)
		handler.RegisterUser(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

}

func TestHandler_LoginUser(t *testing.T) {
	ctx := context.Background()
	req := &api.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	t.Run("success", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().VerifyPassword(ctx, "test@example.com", []byte("password")).
			Return(&domain.UserCredentials{Email: "test@example.com"}, nil)
		userService.EXPECT().GetUserByEmail(ctx, "test@example.com").
			Return(&domain.User{ID: "123", Email: "test@example.com"}, nil)

		body := toJSONBody(t, req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/login", body).WithContext(ctx)
		handler.LoginUser(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().VerifyPassword(ctx, "test@example.com", []byte("password")).
			Return(nil, errors.New("invalid"))

		body := toJSONBody(t, req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/login", body).WithContext(ctx)
		handler.LoginUser(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

}

func TestHandler_GetUserByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().GetUserByID(ctx, "123").
			Return(&domain.User{ID: "123", Email: "test@example.com"}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users/123", nil).WithContext(ctx)
		handler.GetUserByID(w, r, "123")

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().GetUserByID(ctx, "404").Return(nil, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users/404", nil).WithContext(ctx)
		handler.GetUserByID(w, r, "404")

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().GetUserByID(ctx, "500").Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/users/500", nil).WithContext(ctx)
		handler.GetUserByID(w, r, "500")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

}

func TestHandler_UpdateUserProfile(t *testing.T) {
	req := &api.UserUpdateRequest{
		FirstName: ptr("New"),
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "123")
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().UpdateUserProfile(ctx, "123", mock.Anything).
			Return(&domain.User{ID: "123", Profile: &domain.UserProfile{FirstName: "New"}}, nil)

		body := toJSONBody(t, req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/profile", body).WithContext(ctx)
		handler.UpdateUserProfile(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		handler := NewHandler(nil, nil, NewHealthChecker())
		body := toJSONBody(t, req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/profile", body).WithContext(context.Background())
		handler.UpdateUserProfile(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "123")
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().UpdateUserProfile(ctx, "123", mock.Anything).
			Return(nil, errors.New("error"))

		body := toJSONBody(t, req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/profile", body).WithContext(ctx)
		handler.UpdateUserProfile(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

}

func TestHandler_GetUserProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "123")
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().GetUserByID(ctx, "123").
			Return(&domain.User{ID: "123", Email: "test@example.com"}, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/profile", nil).WithContext(ctx)
		handler.GetUserProfile(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		handler := NewHandler(nil, nil, NewHealthChecker())
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/profile", nil).WithContext(context.Background())
		handler.GetUserProfile(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("user not found", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "123")
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().GetUserByID(ctx, "123").Return(nil, nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/profile", nil).WithContext(ctx)
		handler.GetUserProfile(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "123")
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().GetUserByID(ctx, "123").Return(nil, errors.New("error"))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/profile", nil).WithContext(ctx)
		handler.GetUserProfile(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

}

func TestHandler_RequestPasswordReset(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().RequestPasswordReset(ctx, "test@example.com").Return("token", nil)

		body := toJSONBody(t, &api.PasswordResetRequest{Email: "test@example.com"})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/password-reset", body).WithContext(ctx)
		handler.RequestPasswordReset(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().RequestPasswordReset(ctx, "test@example.com").Return("", errors.New("not found"))

		body := toJSONBody(t, &api.PasswordResetRequest{Email: "test@example.com"})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/password-reset", body).WithContext(ctx)
		handler.RequestPasswordReset(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandler_ResetPassword(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().ResetPassword(ctx, "token", []byte("new-pass")).Return(nil)

		body := toJSONBody(t, &api.PasswordResetConfirm{Token: "token", NewPassword: "new-pass"})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/password-reset/confirm", body).WithContext(ctx)
		handler.ResetPassword(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, NewHealthChecker())
		userService.EXPECT().ResetPassword(ctx, "token", []byte("new-pass")).Return(errors.New("invalid"))

		body := toJSONBody(t, &api.PasswordResetConfirm{Token: "token", NewPassword: "new-pass"})
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/password-reset/confirm", body).WithContext(ctx)
		handler.ResetPassword(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_LogoutUser(t *testing.T) {
	handler := NewHandler(nil, nil, NewHealthChecker())
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/logout", nil).WithContext(context.Background())
	handler.LogoutUser(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func ptr[T any](v T) *T {
	return &v
}

func toJSONBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, json.NewEncoder(&buf).Encode(v))
	return &buf
}
