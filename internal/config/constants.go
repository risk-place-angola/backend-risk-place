package config

import "time"

const (
	RefreshTTL      = 14 * 24 * time.Hour
	JwtAccessTTL    = time.Hour * 1
	TokenTypeBearer = "Bearer"

	CodeExpirationDuration = time.Minute * 30

	ReadTimeout  = 15 * time.Second
	WriteTimeout = 15 * time.Second
	IdleTimeout  = 60 * time.Second
)
