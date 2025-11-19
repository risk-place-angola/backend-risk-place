package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type AlertRepository interface {
	Create(ctx context.Context, alert *model.Alert) error
	CreateAlertNotification(ctx context.Context, alertID uuid.UUID, userID string) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Alert, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Alert, error)
	GetSubscribedAlerts(ctx context.Context, userID uuid.UUID) ([]*model.Alert, error)
	Update(ctx context.Context, alert *model.Alert) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
	SubscribeToAlert(ctx context.Context, subscription *model.AlertSubscription) error
	UnsubscribeFromAlert(ctx context.Context, alertID, userID uuid.UUID) error
	IsUserSubscribed(ctx context.Context, alertID, userID uuid.UUID) (bool, error)
}
