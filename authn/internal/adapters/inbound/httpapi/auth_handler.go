package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	core "github.com/faber-numeris/beholder/authn/internal/core/domain"
	"github.com/faber-numeris/beholder/authn/internal/platform/mapper/generated"
	inboundport "github.com/faber-numeris/beholder/authn/internal/ports/inbound"
	outboundport "github.com/faber-numeris/beholder/authn/internal/ports/outbound"
	foundation "github.com/faber-numeris/foundation/beholder/api"
)

type Handler struct {
	userService    inboundport.UserService
	hashingService outboundport.HashingService
	healthChecker  *HealthChecker
}

func (h *Handler) ListAuditLogs(w http.ResponseWriter, r *http.Request, params foundation.ListAuditLogsParams) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) CheckPermission(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) DeleteRelationship(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) ListRelationships(w http.ResponseWriter, r *http.Request, params foundation.ListRelationshipsParams) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) CreateRelationship(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) DeleteUserProfile(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request, params foundation.ListUsersParams) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) DeactivateUser(w http.ResponseWriter, r *http.Request, id foundation.ULID) {
	//TODO implement me
	panic("implement me")
}

var converterImpl = generated.ConverterImpl{}

func NewHandler(
	userService inboundport.UserService,
	hashingService outboundport.HashingService,
	healthChecker *HealthChecker,
) *Handler {
	return &Handler{
		userService:    userService,
		hashingService: hashingService,
		healthChecker:  healthChecker,
	}
}

func (h *Handler) ConfirmUserRegistration(w http.ResponseWriter, r *http.Request, params foundation.ConfirmUserRegistrationParams) {
	ctx := r.Context()
	err := h.userService.ConfirmUserRegistration(ctx, params.Token)
	if err != nil {
		slog.Error("Failed to confirm registration", "error", err)
		writeJSON(w, http.StatusBadRequest, foundation.Error{
			Error:   "INVALID_TOKEN",
			Message: "Invalid or expired confirmation token",
		})
		return
	}

	writeJSON(w, http.StatusOK, foundation.MessageResponse{
		Message: "User registration confirmed successfully",
	})
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request, id foundation.ULID) {
	ctx := r.Context()
	user, err := h.userService.GetUserByID(ctx, string(id))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, foundation.Error{
			Error:   err.Error(),
			Message: "Could not retrieve user. Please try again later.",
		})
		return
	}

	if user == nil {
		writeJSON(w, http.StatusNotFound, foundation.Error{
			Error:   "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	apiUser, err := converterImpl.UserModelToApiUser(*user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, foundation.Error{
			Error:   err.Error(),
			Message: "Could not process user response.",
		})
		return
	}

	writeJSON(w, http.StatusOK, apiUser)
}

func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(UserIDKey)
	if userID == nil {
		writeJSON(w, http.StatusUnauthorized, foundation.Error{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	user, err := h.userService.GetUserByID(ctx, userID.(string))
	if err != nil {
		slog.Error("Failed to get user profile", "error", err)
		writeJSON(w, http.StatusInternalServerError, foundation.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Could not retrieve user profile",
		})
		return
	}

	if user == nil {
		writeJSON(w, http.StatusUnauthorized, foundation.Error{
			Error:   "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	apiUser, err := converterImpl.UserModelToApiUser(*user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, foundation.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Could not process user response",
		})
		return
	}

	writeJSON(w, http.StatusOK, apiUser)
}

const accessTokenTTLSeconds = 3600

// LoginUser authenticates a browser/SPA client and issues the access token as an
// HttpOnly cookie, never in the response body, so page JavaScript cannot read it
// (mitigates XSS-based token theft). Mobile and server-to-server clients should use
// IssueToken instead.
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req foundation.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, foundation.Error{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	accessToken, err := h.authenticate(ctx, string(req.Email), []byte(req.Password))
	if err != nil {
		writeAuthError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   accessTokenTTLSeconds,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusCreated, foundation.MessageResponse{
		Message: "Authenticated successfully",
	})
}

// IssueToken authenticates a mobile or server-to-server client and returns the
// opaque access token and refresh token in the response body, for clients that
// store tokens in OS-level secure storage rather than a browser document.
func (h *Handler) IssueToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req foundation.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, foundation.Error{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	accessToken, err := h.authenticate(ctx, string(req.Email), []byte(req.Password))
	if err != nil {
		writeAuthError(w, err)
		return
	}

	refreshToken, err := generateOpaqueToken()
	if err != nil {
		slog.Error("Failed to generate refresh token", "error", err)
		writeJSON(w, http.StatusInternalServerError, foundation.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Could not generate refresh token",
		})
		return
	}

	writeJSON(w, http.StatusCreated, foundation.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    accessTokenTTLSeconds,
		RefreshToken: refreshToken,
	})
}

// authenticate verifies the given credentials and returns the opaque access token
// for the authenticated user, or an error suitable for writeAuthError.
func (h *Handler) authenticate(ctx context.Context, email string, password []byte) (string, error) {
	_, err := h.userService.VerifyPassword(ctx, email, password)
	if err != nil {
		slog.Error("Login failed", "email", email, "error", err)
		return "", errInvalidCredentials
	}

	user, err := h.userService.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errUserLookupFailed
	}

	return user.ID, nil
}

var (
	errInvalidCredentials = errors.New("invalid email or password")
	errUserLookupFailed   = errors.New("could not retrieve user details")
)

func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errInvalidCredentials):
		writeJSON(w, http.StatusUnauthorized, foundation.Error{
			Error:   "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
		})
	default:
		writeJSON(w, http.StatusInternalServerError, foundation.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Could not retrieve user details",
		})
	}
}

// generateOpaqueToken returns a random hex-encoded token, following the same
// pattern used for confirmation tokens in core/services.
func generateOpaqueToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req foundation.RegisterUserJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, foundation.Error{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	userModel, err := converterImpl.UserModelFromUserRequest(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, foundation.Error{
			Error:   err.Error(),
			Message: "Invalid Register User Data",
		})
		return
	}

	userResponse, err := h.userService.RegisterUser(ctx, &userModel, []byte(req.Password))
	if err != nil {
		if errors.Is(err, core.ErrUserAlreadyExists) {
			writeJSON(w, http.StatusConflict, foundation.Error{
				Error:   "USER_ALREADY_EXISTS",
				Message: "User already exists",
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, foundation.Error{
			Error:   err.Error(),
			Message: "Could not register user. Please try again later.",
		})
		return
	}

	apiUserResponse, err := converterImpl.UserModelToApiUser(*userResponse)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, foundation.Error{
			Error:   err.Error(),
			Message: "Could not process user response.",
		})
		return
	}

	writeJSON(w, http.StatusCreated, apiUserResponse)
}

func (h *Handler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req foundation.RequestPasswordResetJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, foundation.Error{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	_, err := h.userService.RequestPasswordReset(ctx, string(req.Email))
	if err != nil {
		slog.Error("Password reset request failed", "email", req.Email, "error", err)
	}

	writeJSON(w, http.StatusOK, foundation.MessageResponse{
		Message: "If the email exists, a password reset link will be sent",
	})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req foundation.ResetPasswordJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, foundation.Error{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	err := h.userService.ResetPassword(ctx, req.Token, []byte(req.NewPassword))
	if err != nil {
		slog.Error("Password reset failed", "error", err)
		writeJSON(w, http.StatusBadRequest, foundation.Error{
			Error:   "INVALID_TOKEN",
			Message: "Invalid or expired password reset token",
		})
		return
	}

	writeJSON(w, http.StatusOK, foundation.MessageResponse{
		Message: "Password reset successfully",
	})
}

func (h *Handler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(UserIDKey)
	if userID == nil {
		writeJSON(w, http.StatusUnauthorized, foundation.Error{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	var req foundation.UpdateUserProfileJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, foundation.Error{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	userModel := &core.User{
		Profile: &core.UserProfile{},
	}

	if req.FirstName != nil && *req.FirstName != "" {
		userModel.Profile.FirstName = *req.FirstName
	}
	if req.LastName != nil && *req.LastName != "" {
		userModel.Profile.LastName = *req.LastName
	}
	if req.Locale != nil && *req.Locale != "" {
		userModel.Profile.Locale = *req.Locale
	}
	if req.Timezone != nil && *req.Timezone != "" {
		userModel.Profile.Timezone = *req.Timezone
	}

	updatedUser, err := h.userService.UpdateUserProfile(ctx, userID.(string), userModel)
	if err != nil {
		slog.Error("Failed to update user profile", "error", err)
		writeJSON(w, http.StatusInternalServerError, foundation.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Could not update user profile",
		})
		return
	}

	apiUser, err := converterImpl.UserModelToApiUser(*updatedUser)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, foundation.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Could not process user response",
		})
		return
	}

	writeJSON(w, http.StatusOK, apiUser)
}

func (h *Handler) CheckLiveness(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "UP"})
}

func (h *Handler) CheckReadiness(w http.ResponseWriter, r *http.Request) {
	if h.healthChecker.IsReady() {
		writeJSON(w, http.StatusOK, map[string]string{"status": "UP"})
		return
	}

	w.Header().Set("Retry-After", "30")
	writeJSON(w, http.StatusServiceUnavailable, foundation.Error{
		Error:   "SERVICE_DEGRADED",
		Message: "Service is temporarily unavailable. Please retry later.",
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
	}
}
