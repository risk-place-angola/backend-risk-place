package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type EmergencyContactRepository interface {
	GenericRepository[model.EmergencyContact]
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.EmergencyContact, error)
	FindByUserIDAndID(ctx context.Context, userID, contactID uuid.UUID) (*model.EmergencyContact, error)
	FindPriorityByUserID(ctx context.Context, userID uuid.UUID) ([]*model.EmergencyContact, error)
	CountPriorityByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	DeleteByUserIDAndID(ctx context.Context, userID, contactID uuid.UUID) error
}
