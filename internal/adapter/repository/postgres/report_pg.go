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

type LocationStore interface {
	UpdateReportLocation(ctx context.Context, reportID string, lat, lon float64) error
	FindReportsInRadius(ctx context.Context, lat, lon float64, radiusMeters float64) ([]string, error)
	RemoveReportLocation(ctx context.Context, reportID string) error
}

type ReportPG struct {
	q             sqlc.Querier
	locationStore LocationStore
}

func NewReportRepoPG(db *sql.DB, locationStore LocationStore) repository.ReportRepository {
	return &ReportPG{
		q:             sqlc.New(db),
		locationStore: locationStore,
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
	// Cria o report no PostgreSQL e obtém o ID gerado
	reportID, err := r.q.CreateReport(ctx, sqlc.CreateReportParams{
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
	if err != nil {
		slog.Error("failed to create report", "error", err)
		return err
	}

	// Salva a localização do report no Redis para buscas geoespaciais rápidas
	if err := r.locationStore.UpdateReportLocation(ctx, reportID.String(), m.Latitude, m.Longitude); err != nil {
		slog.Warn("failed to update report location in redis", "reportID", reportID, "error", err)
		// Não retorna erro, pois o report já foi criado no PostgreSQL
		// A falha no Redis não deve impedir a criação do report
	}

	// Atualiza o ID no model
	m.ID = reportID

	return nil
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
	// Deleta o report do PostgreSQL
	err := r.q.DeleteReport(ctx, reportID)
	if err != nil {
		slog.Error("failed to delete report", "reportID", reportID, "error", err)
		return err
	}

	// Remove a localização do Redis
	if err := r.locationStore.RemoveReportLocation(ctx, reportID.String()); err != nil {
		slog.Warn("failed to remove report location from redis", "reportID", reportID, "error", err)
		// Não retorna erro, pois o report já foi deletado do PostgreSQL
	}

	return nil
}

func (r *ReportPG) UpdateLocation(ctx context.Context, reportID uuid.UUID, params repository.UpdateLocationParams) error {
	// Atualiza a localização no PostgreSQL
	_, err := r.q.UpdateReportLocation(ctx, sqlc.UpdateReportLocationParams{
		ID:        reportID,
		Latitude:  params.Latitude,
		Longitude: params.Longitude,
		Column4:   params.Address,
		Column5:   params.Neighborhood,
		Column6:   params.Municipality,
		Column7:   params.Province,
	})
	if err != nil {
		slog.Error("failed to update report location", "reportID", reportID, "error", err)
		return err
	}

	// Atualiza a localização no Redis para buscas geoespaciais
	if err := r.locationStore.UpdateReportLocation(ctx, reportID.String(), params.Latitude, params.Longitude); err != nil {
		slog.Warn("failed to update report location in redis", "reportID", reportID, "error", err)
		// Não retorna erro, pois a localização já foi atualizada no PostgreSQL
	}

	slog.Info("report location updated successfully", "reportID", reportID)

	return nil
}

func (r *ReportPG) CreateReportNotification(ctx context.Context, reportID uuid.UUID, userID uuid.UUID) error {
	return r.q.CreateReportNotification(ctx, sqlc.CreateReportNotificationParams{
		ReferenceID: reportID,
		UserID:      userID,
		Type:        "report",
	})
}

func (r *ReportPG) FindByRadius(ctx context.Context, lat float64, lon float64, radiusMeters float64) ([]*model.Report, error) {
	// 1. Busca IDs dos reports próximos usando Redis (ultra-rápido)
	reportIDs, err := r.locationStore.FindReportsInRadius(ctx, lat, lon, radiusMeters)
	if err != nil {
		slog.Error("failed to find reports in radius from redis", "error", err)
		return nil, err
	}

	// Se não encontrou nenhum report próximo, retorna lista vazia
	if len(reportIDs) == 0 {
		slog.Info("no reports found in radius")
		return []*model.Report{}, nil
	}

	// 2. Converte strings para UUIDs
	uuids := make([]uuid.UUID, 0, len(reportIDs))
	for _, id := range reportIDs {
		reportUUID, err := uuid.Parse(id)
		if err != nil {
			slog.Warn("invalid report UUID from redis", "id", id, "error", err)
			continue
		}
		uuids = append(uuids, reportUUID)
	}

	// 3. Busca os dados completos dos reports no PostgreSQL
	items, err := r.q.ListReportsByIDs(ctx, uuids)
	if err != nil {
		slog.Error("failed to find reports by ids", "error", err)
		return nil, err
	}

	slog.Info("found reports in radius", "count", len(items))

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
