package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	domainErrors "github.com/risk-place-angola/backend-risk-place/internal/domain/errors"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	domainrepository "github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
	domainService "github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

type UserUseCase struct {
	userRepo            domainrepository.UserRepository
	roleRepo            domainrepository.RoleRepository
	token               port.TokenGenerator
	hasher              port.PasswordHasher
	config              *config.Config
	migrationService    domainService.AnonymousMigrationService
	verificationService domainService.VerificationService
}

func NewUserUseCase(
	userRepo domainrepository.UserRepository,
	roleRepo domainrepository.RoleRepository,
	token port.TokenGenerator,
	hasher port.PasswordHasher,
	config *config.Config,
	migrationService domainService.AnonymousMigrationService,
	verificationService domainService.VerificationService,
) *UserUseCase {
	return &UserUseCase{
		userRepo:            userRepo,
		roleRepo:            roleRepo,
		token:               token,
		hasher:              hasher,
		config:              config,
		migrationService:    migrationService,
		verificationService: verificationService,
	}
}

func (uc *UserUseCase) Signup(ctx context.Context, input dto.RegisterUserInput, deviceID string) (*dto.RegisterUserOutput, error) {
	user, err := model.NewUser(input.Name, input.Phone, input.Email, input.Password)
	if err != nil {
		slog.Error("Error creating new user model", "email", input.Email, "error", err)
		return nil, err
	}

	hashed, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashed
	user.DeviceToken = input.DeviceFCMToken
	user.DeviceLanguage = input.DeviceLanguage

	err = uc.userRepo.Save(ctx, user)
	if err != nil {
		slog.Error("Error saving new user", "email", input.Email, "error", err)
		return nil, err
	}

	err = uc.roleRepo.AssignRoleNameToUser(ctx, user.ID, "citizen")
	if err != nil {
		slog.Error("Error assigning owner role to user", "user_id", user.ID, "error", err)
		return nil, err
	}

	if err := uc.verificationService.SendCode(ctx, user.ID, user.Phone, user.Email); err != nil {
		if errors.Is(err, domainErrors.ErrSentViaEmail) {
			return &dto.RegisterUserOutput{ID: user.ID}, err
		}
		slog.Error("Failed to send verification code", "user_id", user.ID, "error", err)
		return nil, fmt.Errorf("failed to send verification code: %w", err)
	}

	if deviceID != "" && uc.migrationService != nil {
		go func(bgCtx context.Context, devID string, usrID uuid.UUID) {
			if err := uc.migrationService.MigrateAnonymousData(
				bgCtx,
				devID,
				usrID,
				"signup",
			); err != nil {
				slog.Error("Anonymous data migration failed but signup succeeded",
					"user_id", usrID,
					"device_id", devID,
					"error", err)
			} else {
				slog.Info("Anonymous data migrated successfully after signup",
					"user_id", usrID,
					"device_id", devID)
			}
		}(context.WithoutCancel(ctx), deviceID, user.ID)
	}

	return &dto.RegisterUserOutput{ID: user.ID}, nil
}

func (uc *UserUseCase) Login(ctx context.Context, identifier, rawPassword, deviceID, fcmToken, deviceLanguage string) (*dto.UserSignInDTO, error) {
	u, err := uc.userRepo.FindByEmailOrPhone(ctx, identifier)
	if err != nil || u == nil {
		slog.Error("User not found", "identifier", identifier, "error", err)
		return nil, domainErrors.ErrInvalidCredentials
	}

	if !uc.hasher.Compare(u.Password, rawPassword) {
		slog.Error("Password mismatch for user", "identifier", identifier)
		return nil, domainErrors.ErrInvalidCredentials
	}

	if !u.AccountVerification.Verified {
		slog.Warn("Login attempt with unverified account, resending code", "identifier", identifier, "user_id", u.ID)
		if err := uc.verificationService.ResendCode(ctx, u.ID, u.Phone, u.Email); err != nil {
			slog.Error("Failed to resend verification code", "user_id", u.ID, "error", err)
		}
		return nil, domainErrors.ErrAccountNotVerified
	}

	if fcmToken != "" || deviceLanguage != "" {
		if err := uc.userRepo.UpdateUserDeviceInfo(ctx, u.ID, fcmToken, deviceLanguage); err != nil {
			slog.Error("Failed to update device info", "user_id", u.ID, "error", err)
		}
	}

	roles, err := uc.roleRepo.GetUserRoles(ctx, u.ID)
	if err != nil || len(roles) == 0 {
		slog.Error("No roles assigned to user", "user_id", u.ID, "error", err)
		return nil, domainErrors.ErrNoRolesAssigned
	}

	highest := model.HighestPriorityRole(roles)

	claims := dto.AccessClaims{
		Sub:        u.ID.String(),
		Email:      u.Email,
		Roles:      dto.ToGetRoleNames(roles),
		ActiveRole: highest.Name,
	}

	accessToken, err := uc.token.SignAccessToken(claims)
	if err != nil {
		slog.Error("Error signing access token", "user_id", u.ID, "error", err)
		return nil, err
	}
	refreshToken, err := uc.token.IssueRefreshToken(u.ID, roles, highest.Name)
	if err != nil {
		slog.Error("Error issuing refresh token", "user_id", u.ID, "error", err)
		return nil, err
	}

	if deviceID != "" && uc.migrationService != nil {
		go func(bgCtx context.Context, devID string, usrID uuid.UUID) {
			if err := uc.migrationService.MigrateAnonymousData(
				bgCtx,
				devID,
				usrID,
				"login",
			); err != nil {
				slog.Error("Anonymous data migration failed but login succeeded",
					"user_id", usrID,
					"device_id", devID,
					"error", err)
			} else {
				slog.Info("Anonymous data migrated successfully after login",
					"user_id", usrID,
					"device_id", devID)
			}
		}(context.WithoutCancel(ctx), deviceID, u.ID)
	}

	return &dto.UserSignInDTO{
		AccessToken:  accessToken,
		ExpiresIn:    time.Now().Add(config.JwtAccessTTL).Unix(),
		RefreshToken: refreshToken,
		TokenType:    config.TokenTypeBearer,
		UserProfileResponse: dto.UserProfileResponse{
			ID:         u.ID.String(),
			ActiveRole: highest.Name,
			Email:      u.Email,
			Name:       u.Name,
			RoleName:   dto.ToGetRoleNames(roles),
		},
	}, nil
}

func (uc *UserUseCase) Refresh(ctx context.Context, refreshToken string) (*dto.UserSignInDTO, error) {
	_, claims, err := uc.token.ParseRefresh(refreshToken)
	if err != nil {
		if errors.Is(err, domainErrors.ErrExpiredToken) {
			return nil, domainErrors.ErrExpiredToken
		}
		return nil, domainErrors.ErrInvalidToken
	}

	uid, err := uuid.Parse(claims.Sub)
	if err != nil {
		return nil, domainErrors.ErrInvalidToken
	}
	user, err := uc.userRepo.FindByID(ctx, uid)
	if err != nil || user == nil || user.ID == uuid.Nil {
		return nil, domainErrors.ErrInvalidToken
	}
	if !user.AccountVerification.Verified {
		return nil, domainErrors.ErrAccountNotVerified
	}

	roles, err := uc.roleRepo.GetUserRoles(ctx, user.ID)
	if err != nil || len(roles) == 0 {
		return nil, domainErrors.ErrNoRolesAssigned
	}
	active := model.HighestPriorityRole(roles)

	accessClaims := dto.AccessClaims{
		Sub:        user.ID.String(),
		Email:      user.Email,
		Roles:      dto.ToGetRoleNames(roles),
		ActiveRole: active.Name,
	}
	accessToken, err := uc.token.SignAccessToken(accessClaims)
	if err != nil {
		return nil, err
	}

	newRefresh, err := uc.token.IssueRefreshToken(user.ID, roles, active.Name)
	if err != nil {
		return nil, err
	}

	return &dto.UserSignInDTO{
		AccessToken:  accessToken,
		ExpiresIn:    time.Now().Add(config.JwtAccessTTL).Unix(),
		RefreshToken: newRefresh,
		TokenType:    config.TokenTypeBearer,
		UserProfileResponse: dto.UserProfileResponse{
			ActiveRole: active.Name,
			Email:      user.Email,
			Name:       user.Name,
			RoleName:   dto.ToGetRoleNames(roles),
		},
	}, nil
}

func (uc *UserUseCase) Logout(ctx context.Context, userID uuid.UUID) error {
	return nil
}

func (uc *UserUseCase) ForgotPassword(ctx context.Context, identifier string) (string, error) {
	getAccount, err := uc.userRepo.FindByEmailOrPhone(ctx, identifier)
	if err != nil {
		if errors.Is(err, domainErrors.ErrUserNotFound) {
			slog.Error("User account not found", "identifier", identifier)
			return "", domainErrors.ErrUserAccountNotExists
		}
		return "", err
	}

	if getAccount.ID == uuid.Nil {
		slog.Error("User account not found", "identifier", identifier, "user_id", "nil")
		return "", domainErrors.ErrUserAccountNotExists
	}

	if err := uc.verificationService.SendPasswordResetCode(ctx, getAccount.ID, getAccount.Phone, getAccount.Email); err != nil {
		if errors.Is(err, domainErrors.ErrSentViaEmail) {
			return getAccount.Email, err
		}
		slog.Error("Failed to send password reset code", "user_id", getAccount.ID, "error", err)
		return "", fmt.Errorf("failed to send password reset code: %w", err)
	}

	return "", nil
}

func (uc *UserUseCase) ResetPassword(ctx context.Context, identifier, code, newPassword string) error {
	user, err := uc.userRepo.FindByEmailOrPhone(ctx, identifier)
	if err != nil {
		if errors.Is(err, domainErrors.ErrUserNotFound) {
			slog.Error("User account not found", "identifier", identifier)
			return domainErrors.ErrUserAccountNotExists
		}
		return err
	}

	if user.ID == uuid.Nil {
		slog.Error("User account not found", "identifier", identifier, "user_id", "nil")
		return domainErrors.ErrUserAccountNotExists
	}

	valid, err := uc.verificationService.VerifyCode(ctx, user.ID, code)
	if err != nil {
		slog.Error("Failed to verify reset code", "user_id", user.ID, "error", err)
		return domainErrors.ErrExpiredCode
	}

	if !valid {
		slog.Error("Invalid reset code", "user_id", user.ID)
		return domainErrors.ErrInvalidCode
	}

	hashedPassword, err := uc.hasher.Hash(newPassword)
	if err != nil {
		slog.Error("Error hashing new password", "user_id", user.ID, "error", err)
		return err
	}

	if err := uc.userRepo.UpdateUserPassword(ctx, user.ID, hashedPassword); err != nil {
		slog.Error("Error updating password", "user_id", user.ID, "error", err)
		return err
	}

	return nil
}

func (uc *UserUseCase) VerifyCode(ctx context.Context, identifier, code string) error {
	user, err := uc.userRepo.FindByEmailOrPhone(ctx, identifier)
	if err != nil {
		if errors.Is(err, domainErrors.ErrUserNotFound) {
			slog.Error("User account not found", "identifier", identifier)
			return domainErrors.ErrUserAccountNotExists
		}
		return err
	}

	if user.ID == uuid.Nil {
		slog.Error("User account not found", "identifier", identifier, "user_id", "nil")
		return domainErrors.ErrUserAccountNotExists
	}

	valid, err := uc.verificationService.VerifyCode(ctx, user.ID, code)
	if err != nil {
		slog.Error("Failed to verify code", "user_id", user.ID, "error", err)
		return domainErrors.ErrExpiredCode
	}

	if !valid {
		slog.Error("Invalid verification code", "user_id", user.ID)
		return domainErrors.ErrInvalidCode
	}

	err = uc.userRepo.MarkAccountVerified(ctx, user.ID)
	if err != nil {
		slog.Error("Error marking account as verified", "user_id", user.ID, "error", err)
		return err
	}

	return nil
}

func (uc *UserUseCase) ResendVerificationCode(ctx context.Context, identifier string) (string, error) {
	user, err := uc.userRepo.FindByEmailOrPhone(ctx, identifier)
	if err != nil {
		return "", err
	}

	if user.AccountVerification.Verified {
		return "", fmt.Errorf("account already verified")
	}

	err = uc.verificationService.ResendCode(ctx, user.ID, user.Phone, user.Email)
	return user.Email, err
}

// FindUserByID retrieves a user's profile by their ID.
func (uc *UserUseCase) FindUserByID(ctx context.Context, userID uuid.UUID) (*dto.UserProfileOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.ID == uuid.Nil {
		slog.Error("User not found", "user_id", userID)
		return nil, domainErrors.ErrUserNotFound
	}

	roles, err := uc.roleRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		slog.Error("Error retrieving user roles", "user_id", user.ID, "error", err)
		return nil, err
	}

	return &dto.UserProfileOutput{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Phone: user.Phone,
		Nif:   user.Nif,
		Address: dto.AddressDTO{
			Country:      user.Address.Country,
			Province:     user.Address.Province,
			Municipality: user.Address.Municipality,
			Neighborhood: user.Address.Neighborhood,
			ZipCode:      user.Address.ZipCode,
		},
		RoleName: dto.ToGetRoleNames(roles),
	}, nil
}

// FindAllUsers retrieves all users with their profiles.
func (uc *UserUseCase) FindAllUsers(ctx context.Context) ([]*dto.UserProfileOutput, error) {
	users, err := uc.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	userProfiles := make([]*dto.UserProfileOutput, 0, len(users))
	for _, user := range users {
		roles, err := uc.roleRepo.GetUserRoles(ctx, user.ID)
		if err != nil {
			slog.Error("Error retrieving user roles", "user_id", user.ID, "error", err)
			return nil, err
		}

		userProfiles = append(userProfiles, &dto.UserProfileOutput{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Phone: user.Phone,
			Nif:   user.Nif,
			Address: dto.AddressDTO{
				Country:      user.Address.Country,
				Province:     user.Address.Province,
				Municipality: user.Address.Municipality,
				Neighborhood: user.Address.Neighborhood,
				ZipCode:      user.Address.ZipCode,
			},
			RoleName: dto.ToGetRoleNames(roles),
		})
	}

	return userProfiles, nil
}

// UpdateUser updates a user's profile information.
func (uc *UserUseCase) UpdateUser(ctx context.Context, userID uuid.UUID, input dto.UpdateUserInput) (*dto.UserProfileOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.ID == uuid.Nil {
		slog.Error("User not found", "user_id", userID)
		return nil, domainErrors.ErrUserNotFound
	}

	if err := user.Update(input.Name, input.Phone, input.Email, input.Address.ToEntityAddress()); err != nil {
		return nil, err
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	roles, err := uc.roleRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		slog.Error("Error retrieving user roles", "user_id", user.ID, "error", err)
		return nil, err
	}

	return &dto.UserProfileOutput{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Phone: user.Phone,
		Nif:   user.Nif,
		Address: dto.AddressDTO{
			Country:      user.Address.Country,
			Province:     user.Address.Province,
			Municipality: user.Address.Municipality,
			Neighborhood: user.Address.Neighborhood,
			ZipCode:      user.Address.ZipCode,
		},
		RoleName: dto.ToGetRoleNames(roles),
	}, nil
}

// ChangePassword changes a user's password.
func (uc *UserUseCase) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil || user.ID == uuid.Nil {
		slog.Error("User not found", "user_id", userID)
		return domainErrors.ErrUserNotFound
	}

	if !uc.hasher.Compare(user.Password, currentPassword) {
		slog.Error("Current password mismatch", "user_id", userID)
		return domainErrors.ErrInvalidCurrentPassword
	}

	hashedPassword, err := uc.hasher.Hash(newPassword)
	if err != nil {
		slog.Error("Error hashing new password", "user_id", userID, "error", err)
		return err
	}

	if err := uc.userRepo.UpdateUserPassword(ctx, user.ID, hashedPassword); err != nil {
		slog.Error("Error updating user password", "user_id", user.ID, "error", err)
		return err
	}

	return nil
}

// DeleteUser deletes a user by their ID.
func (uc *UserUseCase) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return uc.userRepo.Delete(ctx, userID.String())
}

func (uc *UserUseCase) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		slog.Error("User not found", "user_id", userID, "error", err)
		return domainErrors.ErrUserNotFound
	}

	if user == nil || user.ID == uuid.Nil {
		return domainErrors.ErrUserNotFound
	}

	var homeAddress *model.SavedLocation
	var workAddress *model.SavedLocation

	if req.HomeAddress != nil {
		homeAddress = &model.SavedLocation{
			Name:      req.HomeAddress.Name,
			Address:   req.HomeAddress.Address,
			Latitude:  req.HomeAddress.Latitude,
			Longitude: req.HomeAddress.Longitude,
		}
	}

	if req.WorkAddress != nil {
		workAddress = &model.SavedLocation{
			Name:      req.WorkAddress.Name,
			Address:   req.WorkAddress.Address,
			Latitude:  req.WorkAddress.Latitude,
			Longitude: req.WorkAddress.Longitude,
		}
	}

	if err := uc.userRepo.UpdateSavedLocations(ctx, userID, homeAddress, workAddress); err != nil {
		slog.Error("Error updating saved locations", "user_id", userID, "error", err)
		return err
	}

	return nil
}

func (uc *UserUseCase) UpdateDeviceInfo(ctx context.Context, userID uuid.UUID, fcmToken, deviceLanguage string) error {
	if err := uc.userRepo.UpdateUserDeviceInfo(ctx, userID, fcmToken, deviceLanguage); err != nil {
		slog.Error("Failed to update device info", "user_id", userID, "error", err)
		return err
	}
	return nil
}

func (uc *UserUseCase) UpdateNotificationPreferences(ctx context.Context, userID uuid.UUID, deviceID string, pushEnabled, smsEnabled bool) error {
	if userID != uuid.Nil {
		if err := uc.userRepo.UpdateNotificationPreferences(ctx, userID, pushEnabled, smsEnabled); err != nil {
			slog.Error("Failed to update notification preferences for user", "user_id", userID, "error", err)
			return err
		}
	}
	return nil
}

//nolint:nonamedreturns // multiple bool returns need names for clarity
func (uc *UserUseCase) GetNotificationPreferences(ctx context.Context, userID uuid.UUID, deviceID string) (pushEnabled, smsEnabled bool, err error) {
	if userID != uuid.Nil {
		return uc.userRepo.GetNotificationPreferences(ctx, userID)
	}
	return false, false, fmt.Errorf("user ID or device ID required")
}
