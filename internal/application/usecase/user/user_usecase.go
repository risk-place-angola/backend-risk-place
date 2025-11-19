package user

import (
	"context"
	stdErrors "errors"
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
)

type UserUseCase struct {
	userRepo     domainrepository.UserRepository
	roleRepo     domainrepository.RoleRepository
	token        port.TokenGenerator
	hasher       port.PasswordHasher
	emailService port.EmailService
	config       *config.Config
}

func NewUserUseCase(
	userRepo domainrepository.UserRepository,
	roleRepo domainrepository.RoleRepository,
	token port.TokenGenerator,
	hasher port.PasswordHasher,
	emailService port.EmailService,
	config *config.Config,
) *UserUseCase {
	return &UserUseCase{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		token:        token,
		hasher:       hasher,
		emailService: emailService,
		config:       config,
	}
}

func (uc *UserUseCase) Signup(ctx context.Context, input dto.RegisterUserInput) (*dto.RegisterUserOutput, error) {
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

	// code := user.GenerateVerificationCode()

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

	/*
			verificationLink := fmt.Sprintf("%s/verify-email?email=%s&code=%s", uc.config.FrontendURL, user.Email, code)
			htmlBody := fmt.Sprintf(`
		        <h2>Bem-vindo(a) à Risk Place Angola!</h2>
		        <p>Para ativar sua conta, clique no link abaixo:</p>
		        <a href="%s">Verificar Email</a>
		        <p>Ou use este código: <strong>%s</strong></p>
		        <p>O código expira em %d minutos.</p>
		    `, verificationLink, code, int(config.CodeExpirationDuration.Minutes()))

			if err := uc.emailService.SendHtml(user.Email, "Verifique seu email", htmlBody); err != nil {
				return nil, fmt.Errorf("erro ao enviar email de verificação: %w", err)
			}
	*/

	return &dto.RegisterUserOutput{ID: user.ID}, nil
}

func (uc *UserUseCase) Login(ctx context.Context, email, rawPassword string) (*dto.UserSignInDTO, error) {
	u, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil || u == nil {
		slog.Error("User not found", "email", email, "error", err)
		return nil, domainErrors.ErrInvalidCredentials
	}

	/*
		if !u.EmailVerification.Verified {
			slog.Error("Email not verified for user", "email", email)
			return nil, domainErrors.ErrEmailNotVerified
		}
	*/

	if !uc.hasher.Compare(u.Password, rawPassword) {
		slog.Error("Password mismatch for user", "email", email)
		return nil, domainErrors.ErrInvalidCredentials
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
		if stdErrors.Is(err, domainErrors.ErrExpiredToken) {
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
	if !user.EmailVerification.Verified {
		return nil, domainErrors.ErrEmailNotVerified
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

func (uc *UserUseCase) ForgotPassword(ctx context.Context, email string) error {
	getAccount, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if stdErrors.Is(err, domainErrors.ErrUserNotFound) {
			slog.Error("User account not found", "email", email)
			return domainErrors.ErrUserAccountNotExists
		}
		return err
	}

	if getAccount.ID == uuid.Nil {
		slog.Error("User account not found", "email", email, "user_id", "nil")
		return domainErrors.ErrUserAccountNotExists
	}

	code := model.GenerateConfirmationCode()
	if err := getAccount.SetGeneratedCode(); err != nil {
		slog.Error("Error setting generated code", "user_id", getAccount.ID, "error", err)
		return err
	}

	err = uc.userRepo.AddCodeToUser(ctx, getAccount.ID, getAccount.EmailVerification.Code, time.Now().Add(config.CodeExpirationDuration))
	if err != nil {
		slog.Error("Error adding code to user", "user_id", getAccount.ID, "error", err)
		return fmt.Errorf("error adding code to user: %w", err)
	}

	if err := uc.emailService.SendEmail(ctx, email, "Password Reset Code",
		fmt.Sprintf("Your password reset code is: %s", code)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (uc *UserUseCase) ResetPassword(ctx context.Context, email, newPassword string) error {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if stdErrors.Is(err, domainErrors.ErrUserNotFound) {
			slog.Error("User account not found", "email", email)
			return domainErrors.ErrUserAccountNotExists
		}
		return err
	}

	if user.ID == uuid.Nil {
		slog.Error("User account not found", "email", email, "user_id", "nil")
		return domainErrors.ErrUserAccountNotExists
	}

	if !user.EmailVerification.CodeVerified {
		slog.Error("Code not verified for user", "user_id", user.ID)
		return domainErrors.ErrInvalidCode
	}

	hashedPassword, err := uc.hasher.Hash(newPassword)
	if err != nil {
		slog.Error("Error hashing new password for user", "user_id", user.ID, "error", err)
		return err
	}

	if err := uc.userRepo.UpdateUserPassword(ctx, user.ID, hashedPassword); err != nil {
		slog.Error("Error updating user password", "user_id", user.ID, "error", err)
		return err
	}

	return nil
}

func (uc *UserUseCase) VerifyCode(ctx context.Context, email, code string) error {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if stdErrors.Is(err, domainErrors.ErrUserNotFound) {
			slog.Error("User account not found", "email", email)
			return domainErrors.ErrUserAccountNotExists
		}
		return err
	}

	if user.ID == uuid.Nil {
		slog.Error("User account not found", "email", email, "user_id", "nil")
		return domainErrors.ErrUserAccountNotExists
	}

	if user.EmailVerification.Code != code {
		slog.Error("Invalid verification code for user", "user_id", user.ID)
		return domainErrors.ErrInvalidCode
	}

	if time.Now().After(user.EmailVerification.ExpiresAt) {
		slog.Error("Verification code expired for user", "user_id", user.ID)
		return domainErrors.ErrExpiredCode
	}

	user.EmailVerification.CodeVerified = true

	err = uc.userRepo.Save(ctx, user)
	if err != nil {
		slog.Error("Error saving user after code verification", "user_id", user.ID, "error", err)
		return err
	}

	return nil
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
