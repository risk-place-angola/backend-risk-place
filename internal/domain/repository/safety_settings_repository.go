package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type SafetySettingsRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.SafetySettings, error)
	Upsert(ctx context.Context, settings *model.SafetySettings) error
}
