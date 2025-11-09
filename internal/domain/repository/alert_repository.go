package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type AlertRepository interface {
	Create(ctx context.Context, alert *model.Alert) error
	CreateAlertNotification(ctx context.Context, alertID uuid.UUID, userID string) error
}
