package service

import (
	"context"

	"github.com/google/uuid"
)

type VerificationService interface {
	SendCode(ctx context.Context, userID uuid.UUID, phone, email string) error
	SendPasswordResetCode(ctx context.Context, userID uuid.UUID, phone, email string) error
	VerifyCode(ctx context.Context, userID uuid.UUID, code string) (bool, error)
	ResendCode(ctx context.Context, userID uuid.UUID, phone, email string) error
	ResendPasswordResetCode(ctx context.Context, userID uuid.UUID, phone, email string) error
}
