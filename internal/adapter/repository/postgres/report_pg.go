package postgres

import (
	"context"
	"database/sql"
	"log/slog"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type ReportPG struct {
	q             sqlc.Querier
	locationStore port.LocationStore
}

func NewReportRepoPG(db *sql.DB, locationStore port.LocationStore) repository.ReportRepository {
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

	var reviewedBy uuid.UUID
	if r.ReviewedBy.Valid {
		reviewedBy = r.ReviewedBy.UUID
	}

	var resolvedAt time.Time
	if r.ResolvedAt.Valid {
		resolvedAt = r.ResolvedAt.Time
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
		ReviewedBy:   reviewedBy,
		ResolvedAt:   resolvedAt,
		CreatedAt:    r.CreatedAt.Time,
		UpdatedAt:    r.UpdatedAt.Time,
	}
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
	if r.locationStore != nil {
		if err := r.locationStore.UpdateReportLocation(ctx, reportID.String(), m.Latitude, m.Longitude); err != nil {
			slog.Warn("failed to update report location in redis", "reportID", reportID, "error", err)
			// Não retorna erro, pois o report já foi criado no PostgreSQL
			// A falha no Redis não deve impedir a criação do report
		}
	} else {
		slog.Warn("location store is nil, skipping redis update", "reportID", reportID)
	}

	// Atualiza o ID no model
	m.ID = reportID

	return nil
}

func (r *ReportPG) GetByID(ctx context.Context, id uuid.UUID) (*model.Report, error) {
	item, err := r.q.GetReportByID(ctx, id)
	if err != nil {
		slog.Error("failed to get report by id", "reportID", id, "error", err)
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
	if r.locationStore != nil {
		if err := r.locationStore.UpdateReportLocation(ctx, reportID.String(), params.Latitude, params.Longitude); err != nil {
			slog.Warn("failed to update report location in redis", "reportID", reportID, "error", err)
			// Não retorna erro, pois a localização já foi atualizada no PostgreSQL
		}
	} else {
		slog.Warn("location store is nil, skipping redis update", "reportID", reportID)
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

func (r *ReportPG) ListWithPagination(ctx context.Context, params repository.ListReportsParams) ([]*model.Report, int, error) {
	const (
		defaultLimit = 20
		maxLimit     = 100
	)

	// Validate and set defaults
	if params.Limit <= 0 {
		params.Limit = defaultLimit
	}
	if params.Limit > maxLimit {
		params.Limit = maxLimit
	}
	if params.Order != "asc" && params.Order != "desc" {
		params.Order = "desc"
	}

	// Calculate offset
	offset := (params.Page - 1) * params.Limit

	// Get total count
	var statusFilter sql.NullString
	if params.Status != "" {
		statusFilter = sql.NullString{String: params.Status, Valid: true}
	}

	total, err := r.q.CountReports(ctx, statusFilter)
	if err != nil {
		slog.Error("failed to count reports", "error", err)
		return nil, 0, err
	}

	// Get paginated results
	// #nosec G115 -- params.Limit and offset are validated to be within safe bounds
	items, err := r.q.ListReportsWithPagination(ctx, sqlc.ListReportsWithPaginationParams{
		Column1: params.Order,
		Limit:   int32(params.Limit),
		Offset:  int32(offset),
		Status:  statusFilter,
	})
	if err != nil {
		slog.Error("failed to list reports with pagination", "error", err)
		return nil, 0, err
	}

	reports := mapSlice(items, dbToModel)
	return reports, int(total), nil
}

func (r *ReportPG) FindByRadiusWithDistance(ctx context.Context, lat float64, lon float64, radiusMeters float64, limit int) ([]repository.ReportWithDistance, error) {
	// 1. Busca IDs dos reports próximos usando Redis com distâncias
	type reportDistance struct {
		id       string
		distance float64
	}

	// Busca reports próximos do Redis
	reportIDs, err := r.locationStore.FindReportsInRadius(ctx, lat, lon, radiusMeters)
	if err != nil {
		slog.Error("failed to find reports in radius from redis", "error", err)
		return nil, err
	}

	// Se não encontrou nenhum report próximo, retorna lista vazia
	if len(reportIDs) == 0 {
		slog.Info("no reports found in radius")
		return []repository.ReportWithDistance{}, nil
	}

	// 2. Calcula distâncias para cada report (usando fórmula de Haversine)
	reportsWithDist := make([]reportDistance, 0, len(reportIDs))
	for _, id := range reportIDs {
		// Por enquanto, vamos buscar o report do PostgreSQL para obter as coordenadas
		// e calcular a distância usando a fórmula de Haversine
		reportUUID, err := uuid.Parse(id)
		if err != nil {
			slog.Warn("invalid report UUID from redis", "id", id, "error", err)
			continue
		}

		// Busca report do PostgreSQL
		reportData, err := r.q.GetReportByID(ctx, reportUUID)
		if err != nil {
			slog.Warn("failed to get report", "id", id, "error", err)
			continue
		}

		// Calcula distância usando Haversine
		distance := haversineDistance(lat, lon, reportData.Latitude, reportData.Longitude)

		reportsWithDist = append(reportsWithDist, reportDistance{
			id:       id,
			distance: distance,
		})
	}

	// 3. Ordena por distância (mais próximo primeiro)
	// Implementação simples de bubble sort (para poucos items é suficiente)
	for i := range len(reportsWithDist) - 1 {
		for j := range len(reportsWithDist) - i - 1 {
			if reportsWithDist[j].distance > reportsWithDist[j+1].distance {
				reportsWithDist[j], reportsWithDist[j+1] = reportsWithDist[j+1], reportsWithDist[j]
			}
		}
	}

	// 4. Aplica limite se especificado
	if limit > 0 && limit < len(reportsWithDist) {
		reportsWithDist = reportsWithDist[:limit]
	}

	// 5. Converte para UUIDs
	uuids := make([]uuid.UUID, 0, len(reportsWithDist))
	for _, rd := range reportsWithDist {
		reportUUID, _ := uuid.Parse(rd.id)
		uuids = append(uuids, reportUUID)
	}

	// 6. Busca os dados completos dos reports no PostgreSQL
	items, err := r.q.ListReportsByIDs(ctx, uuids)
	if err != nil {
		slog.Error("failed to find reports by ids", "error", err)
		return nil, err
	}

	// 7. Monta resultado final com distâncias
	result := make([]repository.ReportWithDistance, 0, len(items))
	distanceMap := make(map[string]float64)
	for _, rd := range reportsWithDist {
		distanceMap[rd.id] = rd.distance
	}

	for _, item := range items {
		distance := distanceMap[item.ID.String()]
		result = append(result, repository.ReportWithDistance{
			Report:   dbToModel(item),
			Distance: distance,
		})
	}

	slog.Info("found reports in radius with distances", "count", len(result))

	return result, nil
}

// haversineDistance calcula a distância entre dois pontos usando a fórmula de Haversine
// Retorna a distância em metros
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const (
		earthRadiusMeters = 6371000 // Raio da Terra em metros
		degreesToRadians  = math.Pi / 180
		half              = 0.5
		two               = 2.0
	)

	// Converte graus para radianos
	lat1Rad := lat1 * degreesToRadians
	lat2Rad := lat2 * degreesToRadians
	deltaLat := (lat2 - lat1) * degreesToRadians
	deltaLon := (lon2 - lon1) * degreesToRadians

	// Fórmula de Haversine
	sinDeltaLatHalf := math.Sin(deltaLat * half)
	sinDeltaLonHalf := math.Sin(deltaLon * half)

	a := sinDeltaLatHalf*sinDeltaLatHalf +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			sinDeltaLonHalf*sinDeltaLonHalf
	c := two * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeters * c
}

func mapSlice[T any, R any](in []T, fn func(T) *R) []*R {
	out := make([]*R, 0, len(in))
	for _, v := range in {
		out = append(out, fn(v))
	}
	return out
}
