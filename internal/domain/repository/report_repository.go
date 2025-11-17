package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type UpdateLocationParams struct {
	Latitude     float64
	Longitude    float64
	Address      string
	Neighborhood string
	Municipality string
	Province     string
}

type ReportRepository interface {
	Create(ctx context.Context, r *model.Report) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Report, error)
	ListByStatus(ctx context.Context, status model.ReportStatus) ([]*model.Report, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Report, error)
	VerifyReport(ctx context.Context, reportID uuid.UUID, reviewerID uuid.UUID) error
	ResolveReport(ctx context.Context, reportID uuid.UUID) error
	DeleteReport(ctx context.Context, reportID uuid.UUID) error
	UpdateLocation(ctx context.Context, reportID uuid.UUID, params UpdateLocationParams) error
	CreateReportNotification(ctx context.Context, reportID uuid.UUID, userID uuid.UUID) error
	FindByRadius(ctx context.Context, lat float64, lon float64, radiusMeters float64) ([]*model.Report, error)
}
