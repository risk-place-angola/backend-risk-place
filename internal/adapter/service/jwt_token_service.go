package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"log/slog"
	"time"
)

type JwtClaims struct {
	cfg config.Config
}

type TokenClaims struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func NewJwtTokenService(cfg config.Config) port.TokenGenerator {
	return &JwtClaims{
		cfg: cfg,
	}
}

func (cl *JwtClaims) SignAccessToken(c dto.AccessClaims) (string, error) {
	c.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    cl.cfg.JWTIssuer,
		Audience:  []string{cl.cfg.JWTAudience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JwtAccessTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenString, err := token.SignedString([]byte(cl.cfg.JWTSecret))
	if err != nil {
		slog.Error("Error signing access token", "error", err)
		return "", err
	}

	return tokenString, nil
}

func (cl *JwtClaims) IssueRefreshToken(userID uuid.UUID, roles []model.Role, activeRole string) (string, error) {
	claims := dto.RefreshClaims{
		Sub:        userID.String(),
		Roles:      roles,
		ActiveRole: activeRole,
		Purpose:    "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cl.cfg.JWTIssuer,
			Audience:  []string{cl.cfg.JWTAudience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.RefreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()), // Not before time
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cl.cfg.JWTSecret))
	if err != nil {
		slog.Error("Error signing refresh token", "error", err)
		return "", err
	}

	return tokenString, nil
}

func (cl *JwtClaims) Generate(userID uuid.UUID, role []model.Role) (string, error) {
	// Generate JWT token with user ID and roles
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"roles":   role,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		"iat":     jwt.NewNumericDate(time.Now()),
		"iss":     cl.cfg.JWTIssuer,
		"aud":     cl.cfg.JWTAudience,
		"jti":     uuid.New().String(),            // Unique identifier for the token
		"nbf":     jwt.NewNumericDate(time.Now()), // Not before time
	})

	tokenString, err := token.SignedString([]byte(cl.cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (cl *JwtClaims) Parse(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return []byte(cl.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return token, nil
}

func (cl *JwtClaims) ParseRefresh(tokenStr string) (*jwt.Token, dto.RefreshClaims, error) {
	var out dto.RefreshClaims
	tok, err := jwt.ParseWithClaims(tokenStr, &out, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(cl.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, out, err
	}
	if !tok.Valid {
		return nil, out, jwt.ErrTokenInvalidClaims
	}
	if out.Purpose != "refresh" {
		return nil, out, jwt.ErrTokenInvalidClaims
	}
	return tok, out, nil
}

func (cl *JwtClaims) ParseAccess(tokenStr string) (*jwt.Token, dto.AccessClaims, error) {
	var out dto.AccessClaims
	tok, err := jwt.ParseWithClaims(tokenStr, &out, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(cl.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, out, err
	}
	if !tok.Valid {
		return nil, out, jwt.ErrTokenInvalidClaims
	}
	return tok, out, nil
}

func (cl *JwtClaims) ExtractUserID(token *jwt.Token) (uuid.UUID, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, jwt.ErrInvalidKeyType
	}

	userIDRaw := claims["sub"]
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		return uuid.Nil, jwt.ErrInvalidKeyType
	}

	return uuid.Parse(userIDStr)
}

func (cl *JwtClaims) GenerateEmailVerificationToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(config.CodeExpirationDuration).Unix(),
		"type":    "email_verification",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cl.cfg.JWTSecret + "emailHS256verification"))
}

func (cl *JwtClaims) ValidateEmailVerificationToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check if the token's signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Error("Invalid signing method for email verification token", "token", tokenString)
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(cl.cfg.JWTSecret + "emailHS256verification"), nil
	})

	if err != nil {
		slog.Error("Error parsing email verification token", "token", tokenString, "error", err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if the token is for email verification
		if claims["type"] != "email_verification" {
			slog.Error("Invalid token type for email verification", "token", tokenString)
			return "", jwt.ErrTokenInvalidClaims
		}
		// Extract the user ID from the claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			slog.Error("Invalid user ID in email verification token", "token", tokenString)
			return "", jwt.ErrTokenMalformed
		}

		// Check if the token has expired
		if exp, ok := claims["exp"].(float64); ok && time.Now().Unix() > int64(exp) {
			slog.Error("Email verification token has expired", "token", tokenString)
			return "", jwt.ErrTokenExpired
		}

		return userID, nil
	}

	slog.Error("Invalid claims in email verification token", "token", tokenString)
	return "", jwt.ErrTokenInvalidClaims
}
