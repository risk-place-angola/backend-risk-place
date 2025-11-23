package repository

import (
	"context"
	"time"

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

type ListReportsParams struct {
	Page   int
	Limit  int
	Status string
	Sort   string
	Order  string
}

type ReportWithDistance struct {
	Report   *model.Report
	Distance float64
}

// ReportRepository defines the repository interface for report operations.
// This interface has more than 10 methods because it handles complete CRUD operations
// plus specialized queries (pagination, geospatial, notifications).
//
//nolint:interfacebloat // Repository pattern requires comprehensive interface
type ReportRepository interface {
	Create(ctx context.Context, r *model.Report) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Report, error)
	ListByStatus(ctx context.Context, status model.ReportStatus) ([]*model.Report, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Report, error)
	ListWithPagination(ctx context.Context, params ListReportsParams) ([]*model.Report, int, error)
	VerifyReport(ctx context.Context, reportID uuid.UUID, reviewerID uuid.UUID) error
	ResolveReport(ctx context.Context, reportID uuid.UUID) error
	DeleteReport(ctx context.Context, reportID uuid.UUID) error
	UpdateLocation(ctx context.Context, reportID uuid.UUID, params UpdateLocationParams) error
	CreateReportNotification(ctx context.Context, reportID uuid.UUID, userID uuid.UUID) error
	FindByRadius(ctx context.Context, lat float64, lon float64, radiusMeters float64) ([]*model.Report, error)
	FindByRadiusWithDistance(ctx context.Context, lat float64, lon float64, radiusMeters float64, limit int) ([]ReportWithDistance, error)

	AddVote(ctx context.Context, vote *model.ReportVote) error
	RemoveVote(ctx context.Context, reportID, userID uuid.UUID) error
	RemoveAnonymousVote(ctx context.Context, reportID, sessionID uuid.UUID) error
	GetUserVote(ctx context.Context, reportID, userID uuid.UUID) (*model.ReportVote, error)
	GetAnonymousVote(ctx context.Context, reportID, sessionID uuid.UUID) (*model.ReportVote, error)
	UpdateVerificationCounts(ctx context.Context, reportID uuid.UUID, upvotes, downvotes int) error
	FindDuplicates(ctx context.Context, lat, lon float64, riskTypeID uuid.UUID, radiusMeters float64, since time.Time) ([]*model.Report, error)
	ExpireOldReports(ctx context.Context, before time.Time) error
	UpdateTrustScore(ctx context.Context, userID uuid.UUID, score int) error
	IncrementReportsSubmitted(ctx context.Context, userID uuid.UUID) error
	IncrementReportsVerified(ctx context.Context, userID uuid.UUID) error
}
