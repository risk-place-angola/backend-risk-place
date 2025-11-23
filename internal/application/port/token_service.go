package port

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type TokenGenerator interface {
	Generate(userID uuid.UUID, role []model.Role) (string, error)
	IssueRefreshToken(userID uuid.UUID, roles []model.Role, activeRole string) (string, error)
	ParseRefresh(tokenStr string) (*jwt.Token, dto.RefreshClaims, error)
	ParseAccess(tokenStr string) (*jwt.Token, dto.AccessClaims, error)
	SignAccessToken(user dto.AccessClaims) (string, error)
	GenerateEmailVerificationToken(userID string) (string, error)
	ValidateEmailVerificationToken(tokenString string) (string, error)
}
