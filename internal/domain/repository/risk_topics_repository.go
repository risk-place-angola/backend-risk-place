package repository

import (
	"context"

	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type RiskTopicsRepository interface {
	ListRiskTopics(ctx context.Context, riskTypeID *string) ([]model.RiskTopic, error)
	GetRiskTopicByID(ctx context.Context, id string) (model.RiskTopic, error)
	UpdateRiskTopicIcon(ctx context.Context, id string, iconPath string) error
}
