package middleware

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"os"
	"strings"
	"time"
)

type CustomClaims struct {
	Username string
	jwt.StandardClaims
}

func AuthMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(jwtConfig())
}

func jwtConfig() echojwt.Config {
	signingKey := os.Getenv("JWT_SECRET")
	return echojwt.Config{
		TokenLookup: "header:Authorization",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			auth = strings.Split(auth, "Bearer ")[1]
			keyFunc := func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != "HS256" {
					return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
				}
				return []byte(signingKey), nil
			}

			token, err := jwt.Parse(auth, keyFunc)
			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, errors.New("invalid token")
			}
			return token, nil
		},
	}
}

func WebsocketAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get("Authorization")
		if ok, err := IsValidToken(authHeader); !ok || err != nil {
			return echo.NewHTTPError(401, "Unauthorized")
		}
		return next(ctx)
	}
}

func NewAuthToken(username string) (string, error) {
	var signingKey string = os.Getenv("JWT_SECRET")
	claims := CustomClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", nil
	}
	return t, nil
}

func IsValidToken(token string) (bool, error) {
	var signingKey string = os.Getenv("JWT_SECRET")
	token = strings.Split(token, "Bearer ")[1]
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
