package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/faber-numeris/beholder/backend/authn/internal/adapters/inbound/httpapi/gen"
	"github.com/faber-numeris/beholder/backend/authn/internal/core/domain"
	"github.com/faber-numeris/beholder/backend/authn/internal/platform/mapper/generated"
	inboundport "github.com/faber-numeris/beholder/backend/authn/internal/ports/inbound"
	outboundport "github.com/faber-numeris/beholder/backend/authn/internal/ports/outbound"
)

type Handler struct {
	userService    inboundport.UserService
	hashingService outboundport.HashingService
	healthChecker  *HealthChecker
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

func (h *Handler) ConfirmUserRegistration(w http.ResponseWriter, r *http.Request, params api.ConfirmUserRegistrationParams) {
	ctx := r.Context()
	err := h.userService.ConfirmUserRegistration(ctx, params.Token)
	if err != nil {
		slog.Error("Failed to confirm registration", "error", err)
		writeJSON(w, http.StatusBadRequest, api.Error{
			Error:   "INVALID_TOKEN",
			Message: "Invalid or expired confirmation token",
		})
		return
	}

	writeJSON(w, http.StatusOK, api.MessageResponse{
		Message: "User registration confirmed successfully",
	})
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request, id api.ULID) {
	ctx := r.Context()
	user, err := h.userService.GetUserByID(ctx, string(id))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, api.Error{
			Error:   err.Error(),
			Message: "Could not retrieve user. Please try again later.",
		})
		return
	}

	if user == nil {
		writeJSON(w, http.StatusNotFound, api.Error{
			Error:   "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	apiUser, err := converterImpl.UserModelToApiUser(*user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, api.Error{
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
		writeJSON(w, http.StatusUnauthorized, api.Error{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	user, err := h.userService.GetUserByID(ctx, userID.(string))
	if err != nil {
		slog.Error("Failed to get user profile", "error", err)
		writeJSON(w, http.StatusInternalServerError, api.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Could not retrieve user profile",
		})
		return
	}

	if user == nil {
		writeJSON(w, http.StatusUnauthorized, api.Error{
			Error:   "USER_NOT_FOUND",
			Message: "User not found",
		})
		return
	}

	apiUser, err := converterImpl.UserModelToApiUser(*user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, api.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Could not process user response",
		})
		return
	}

	writeJSON(w, http.StatusOK, apiUser)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req api.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, api.Error{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	_, err := h.userService.VerifyPassword(ctx, string(req.Email), []byte(req.Password))
	if err != nil {
		slog.Error("Login failed", "email", req.Email, "error", err)
		writeJSON(w, http.StatusUnauthorized, api.Error{
			Error:   "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
		})
		return
	}

	user, err := h.userService.GetUserByEmail(ctx, string(req.Email))
	if err != nil || user == nil {
		writeJSON(w, http.StatusInternalServerError, api.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Could not retrieve user details",
		})
		return
	}

	apiUser, err := converterImpl.UserModelToApiUser(*user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, api.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Could not process user response",
		})
		return
	}

	writeJSON(w, http.StatusOK, api.LoginResponse{
		AccessToken: user.ID,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		User:        apiUser,
	})
}

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req api.RegisterUserJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, api.Error{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	userModel, err := converterImpl.UserModelFromUserRequest(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, api.Error{
			Error:   err.Error(),
			Message: "Invalid Register User Data",
		})
		return
	}

	userResponse, err := h.userService.RegisterUser(ctx, &userModel, []byte(req.Password))
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			writeJSON(w, http.StatusConflict, api.Error{
				Error:   "USER_ALREADY_EXISTS",
				Message: "User already exists",
			})
			return
		}
		writeJSON(w, http.StatusInternalServerError, api.Error{
			Error:   err.Error(),
			Message: "Could not register user. Please try again later.",
		})
		return
	}

	apiUserResponse, err := converterImpl.UserModelToApiUser(*userResponse)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, api.Error{
			Error:   err.Error(),
			Message: "Could not process user response.",
		})
		return
	}

	writeJSON(w, http.StatusCreated, apiUserResponse)
}

func (h *Handler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req api.RequestPasswordResetJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, api.Error{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	_, err := h.userService.RequestPasswordReset(ctx, string(req.Email))
	if err != nil {
		slog.Error("Password reset request failed", "email", req.Email, "error", err)
	}

	writeJSON(w, http.StatusOK, api.MessageResponse{
		Message: "If the email exists, a password reset link will be sent",
	})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req api.ResetPasswordJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, api.Error{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	err := h.userService.ResetPassword(ctx, req.Token, []byte(req.NewPassword))
	if err != nil {
		slog.Error("Password reset failed", "error", err)
		writeJSON(w, http.StatusBadRequest, api.Error{
			Error:   "INVALID_TOKEN",
			Message: "Invalid or expired password reset token",
		})
		return
	}

	writeJSON(w, http.StatusOK, api.MessageResponse{
		Message: "Password reset successfully",
	})
}

func (h *Handler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(UserIDKey)
	if userID == nil {
		writeJSON(w, http.StatusUnauthorized, api.Error{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	var req api.UpdateUserProfileJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, api.Error{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	userModel := &domain.User{
		Profile: &domain.UserProfile{},
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
		writeJSON(w, http.StatusInternalServerError, api.Error{
			Error:   "INTERNAL_ERROR",
			Message: "Could not update user profile",
		})
		return
	}

	apiUser, err := converterImpl.UserModelToApiUser(*updatedUser)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, api.Error{
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
	writeJSON(w, http.StatusServiceUnavailable, api.Error{
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
