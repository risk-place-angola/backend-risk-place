package repository

import (
	"context"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type RiskTypesRepository interface {
	CreateRiskType(ctx context.Context, name string, description string, defaultRadiusMeters *int) error
	ListRiskTypes(ctx context.Context) ([]model.RiskType, error)
	GetRiskTypeByID(ctx context.Context, id string) (model.RiskType, error)
	UpdateRiskType(ctx context.Context, id string, name string, description string, defaultRadiusMeters int) error
	DeleteRiskType(ctx context.Context, id string) error
}
