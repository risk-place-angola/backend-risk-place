package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	domainErrors "github.com/risk-place-angola/backend-risk-place/internal/domain/errors"
)

type UserHandler struct {
	userUseCase *application.Application
}

func NewUserHandler(userUseCase *application.Application) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// Signup godoc
// @Summary Register a new user
// @Description Register a new user. If X-Device-ID header is provided, anonymous user data will be migrated to the new account.
// @Tags auth
// @Accept json
// @Produce json
// @Param X-Device-ID header string false "Device ID for anonymous data migration"
// @Param user body dto.RegisterUserInput true "User registration data"
// @Success 201 {object} dto.RegisterUserOutput
// @Failure 400 {object} util.ErrorResponse
// @Router /auth/signup [post]
func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	deviceID := r.Header.Get("X-Device-Id")
	if deviceID == "" {
		deviceID = r.Header.Get("Device-Id")
	}

	userOut, err := h.userUseCase.UserUseCase.Signup(r.Context(), req, deviceID)
	if err != nil {
		slog.Error("error during signup", slog.Any("error", err))
		util.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.Response(w, userOut, http.StatusCreated)
}

// Login godoc
// @Summary Login a user
// @Description Login a user. If X-Device-ID header is provided, anonymous user data will be migrated to the authenticated account.
// @Tags auth
// @Accept json
// @Produce json
// @Param X-Device-ID header string false "Device ID for anonymous data migration"
// @Param credentials body dto.LoginInput true "User login credentials"
// @Success 200 {object} dto.UserSignInDTO
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginInput

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	deviceID := r.Header.Get("X-Device-Id")
	if deviceID == "" {
		deviceID = r.Header.Get("Device-Id")
	}

	token, err := h.userUseCase.UserUseCase.Login(r.Context(), req.Email, req.Password, deviceID, req.DeviceFCMToken, req.DeviceLanguage)
	if err != nil {
		switch {
		case errors.Is(err, domainErrors.ErrInvalidCredentials):
			util.Error(w, "invalid credentials", http.StatusUnauthorized)
		case errors.Is(err, domainErrors.ErrEmailNotVerified):
			util.Error(w, "email not verified", http.StatusForbidden)
		default:
			util.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	util.Response(w, token, http.StatusOK)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access and refresh tokens using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body object{refresh_token=string} true "Refresh token"
// @Success 200 {object} dto.UserSignInDTO
// @Failure 400 {object} util.ErrorResponse
// @Failure 401 {object} util.ErrorResponse
// @Router /auth/refresh [post]
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		util.Error(w, "refresh token is required", http.StatusBadRequest)
		return
	}

	token, err := h.userUseCase.UserUseCase.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, domainErrors.ErrExpiredToken):
			util.Error(w, "refresh token expired", http.StatusUnauthorized)
		case errors.Is(err, domainErrors.ErrInvalidToken):
			util.Error(w, "invalid refresh token", http.StatusUnauthorized)
		case errors.Is(err, domainErrors.ErrEmailNotVerified):
			util.Error(w, "email not verified", http.StatusForbidden)
		default:
			util.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	util.Response(w, token, http.StatusOK)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user and invalidate session
// @Tags auth
// @Security BearerAuth
// @Success 200 {string} string "logout successful"
// @Failure 401 {object} util.ErrorResponse
// @Router /auth/logout [post]
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.userUseCase.UserUseCase.Logout(r.Context(), userID); err != nil {
		util.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	util.Response(w, map[string]string{"message": "logout successful"}, http.StatusOK)
}

// ForgotPassword godoc
// @Summary Initiate password reset
// @Description Send a password reset code to the user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param email body object{email=string} true "User email"
// @Success 200 {string} string "password reset code sent"
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /auth/password/forgot [post]
func (h *UserHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.userUseCase.UserUseCase.ForgotPassword(r.Context(), req.Email)
	if err != nil {
		switch {
		case errors.Is(err, domainErrors.ErrUserAccountNotExists):
			util.Error(w, "user not found", http.StatusNotFound)
		default:
			util.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	util.Response(w, "password reset code sent", http.StatusOK)
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset user password using the reset code
// @Tags auth
// @Accept json
// @Produce json
// @Param reset body object{email=string,password=string} true "Password reset data"
// @Success 200 {string} string "password reset successfully"
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /auth/password/reset [post]
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.userUseCase.UserUseCase.ResetPassword(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domainErrors.ErrUserAccountNotExists):
			util.Error(w, "user not found", http.StatusNotFound)
		case errors.Is(err, domainErrors.ErrInvalidCode):
			util.Error(w, "invalid confirmation code", http.StatusBadRequest)
		default:
			util.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	util.Response(w, "password reset successfully", http.StatusOK)
}

// ConfirmSignup godoc
// @Summary Confirm user signup
// @Description Confirm user signup using the verification code
// @Tags auth
// @Accept json
// @Produce json
// @Param confirmation body object{email=string,code=string} true "Signup confirmation data"
// @Success 204 {string} string "signup confirmed successfully"
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /auth/confirm [post]
func (h *UserHandler) ConfirmSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.userUseCase.UserUseCase.VerifyCode(r.Context(), req.Email, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, domainErrors.ErrUserNotFound):
			util.Error(w, "user not found", http.StatusNotFound)
		case errors.Is(err, domainErrors.ErrInvalidCode):
			util.Error(w, "invalid verification code", http.StatusBadRequest)
		case errors.Is(err, domainErrors.ErrExpiredCode):
			util.Error(w, "verification code expired", http.StatusBadRequest)
		default:
			util.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ResendCode godoc
// @Summary Resend verification code
// @Description Resend verification code to user's phone (SMS) or email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object{email=string} true "Email"
// @Success 204
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /auth/resend-code [post]
func (h *UserHandler) ResendCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.UserUseCase.ResendVerificationCode(r.Context(), req.Email); err != nil {
		slog.Error("Failed to resend verification code", "email", req.Email, "error", err)
		util.Error(w, "failed to resend verification code", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Me godoc
// @Summary Get current user info
// @Description Get information about the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserProfileOutput
// @Failure 401 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Error("failed to get user ID from context")
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := dto.ParseUUID(userIDStr)
	if err != nil {
		slog.Error("invalid user ID in context", slog.Any("error", err))
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userOut, err := h.userUseCase.UserUseCase.FindUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrUserNotFound) {
			util.Error(w, "user not found", http.StatusNotFound)
			return
		}
		util.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	util.Response(w, userOut, http.StatusOK)
}

// UpdateProfile godoc
// @Summary Update user profile with saved locations
// @Description Update user profile to save home and work addresses for navigation
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body dto.UpdateProfileRequest true "Profile update with home/work addresses"
// @Success 200 {object} map[string]string "Profile updated successfully"
// @Failure 400 {object} util.ErrorResponse "Invalid request body"
// @Failure 401 {object} util.ErrorResponse "Unauthorized - missing or invalid JWT token"
// @Failure 404 {object} util.ErrorResponse "User not found"
// @Failure 500 {object} util.ErrorResponse "Internal server error"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := util.GetUserIDFromContext(r.Context())
	if !ok {
		slog.Error("failed to get user ID from context")
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := dto.ParseUUID(userIDStr)
	if err != nil {
		slog.Error("invalid user ID in context", slog.Any("error", err))
		util.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.UserUseCase.UpdateUserProfile(r.Context(), userID, &req); err != nil {
		if errors.Is(err, domainErrors.ErrUserNotFound) {
			util.Error(w, "user not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to update user profile", "error", err)
		util.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	util.Response(w, map[string]string{"message": "profile updated successfully"}, http.StatusOK)
}
