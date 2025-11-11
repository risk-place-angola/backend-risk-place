package postgres

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type ReportPG struct {
	q sqlc.Querier
}

func NewReportRepoPG(db *sql.DB) repository.ReportRepository {
	return &ReportPG{
		q: sqlc.New(db),
	}
}

func dbToModel(r sqlc.Report) *model.Report {
	status, ok := r.Status.(model.ReportStatus)
	if !ok {
		status = model.ReportStatus("pending")
	}

	return &model.Report{
		ID:           r.ID,
		UserID:       r.UserID,
		RiskTypeID:   r.RiskTypeID,
		RiskTopicID:  r.RiskTopicID.UUID,
		Description:  r.Description.String,
		Latitude:     r.Latitude,
		Longitude:    r.Longitude,
		Province:     r.Province.String,
		Municipality: r.Municipality.String,
		Neighborhood: r.Neighborhood.String,
		Address:      r.Address.String,
		ImageURL:     r.ImageUrl.String,
		Status:       status,
		ReviewedBy:   *uuidToPtr(r.ReviewedBy.UUID),
		ResolvedAt:   *timePtr(r.ResolvedAt),
		CreatedAt:    r.CreatedAt.Time,
		UpdatedAt:    r.UpdatedAt.Time,
	}
}

func uuidToPtr(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

func timePtr(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func (r *ReportPG) Create(ctx context.Context, m *model.Report) error {
	return r.q.CreateReport(ctx, sqlc.CreateReportParams{
		UserID:       m.UserID,
		RiskTypeID:   m.RiskTypeID,
		RiskTopicID:  uuid.NullUUID{UUID: m.RiskTopicID, Valid: m.RiskTopicID != uuid.Nil},
		Description:  sqlString(m.Description),
		Latitude:     m.Latitude,
		Longitude:    m.Longitude,
		Province:     sqlString(m.Province),
		Municipality: sqlString(m.Municipality),
		Neighborhood: sqlString(m.Neighborhood),
		Address:      sqlString(m.Address),
		ImageUrl:     sqlString(m.ImageURL),
	})
}

func (r *ReportPG) GetByID(ctx context.Context, id uuid.UUID) (*model.Report, error) {
	item, err := r.q.GetReportByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dbToModel(item), nil
}

func (r *ReportPG) ListByStatus(ctx context.Context, status model.ReportStatus) ([]*model.Report, error) {
	items, err := r.q.ListReportsByStatus(ctx, status)
	if err != nil {
		return nil, err
	}
	return mapSlice(items, dbToModel), nil
}

func (r *ReportPG) ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Report, error) {
	items, err := r.q.ListReportsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mapSlice(items, dbToModel), nil
}

func (r *ReportPG) VerifyReport(ctx context.Context, reportID uuid.UUID, reviewerID uuid.UUID) error {
	return r.q.VerifyReport(ctx, sqlc.VerifyReportParams{
		ID:         reportID,
		ReviewedBy: uuid.NullUUID{UUID: reviewerID, Valid: true},
	})
}

func (r *ReportPG) ResolveReport(ctx context.Context, reportID uuid.UUID) error {
	return r.q.ResolveReport(ctx, reportID)
}

func (r *ReportPG) DeleteReport(ctx context.Context, reportID uuid.UUID) error {
	return r.q.DeleteReport(ctx, reportID)
}

func (r *ReportPG) Update(ctx context.Context, m *model.Report) error {
	// Pode ser criado depois com UpdateCustom
	return sql.ErrNoRows // ainda não implementado no SQLC
}

func (r *ReportPG) CreateReportNotification(ctx context.Context, reportID uuid.UUID, userID uuid.UUID) error {
	return r.q.CreateReportNotification(ctx, sqlc.CreateReportNotificationParams{
		ReferenceID: reportID,
		UserID:      userID,
		Type:        "report",
	})
}

func (r *ReportPG) FindByRadius(ctx context.Context, lat float64, lon float64, radiusMeters float64) ([]*model.Report, error) {
	items, err := r.q.ListReportsNearby(ctx, sqlc.ListReportsNearbyParams{
		Column1: lon,
		Column2: lat,
		Column3: radiusMeters,
	})
	if err != nil {
		slog.Error("failed to find reports by radius", "error", err)
		return nil, err
	}

	return mapSlice(items, dbToModel), nil
}

func sqlString(v string) sql.NullString {
	return sql.NullString{String: v, Valid: v != ""}
}

func mapSlice[T any, R any](in []T, fn func(T) *R) []*R {
	out := make([]*R, 0, len(in))
	for _, v := range in {
		out = append(out, fn(v))
	}
	return out
}
