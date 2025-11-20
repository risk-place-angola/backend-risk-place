package postgres

import (
	"context"
	"database/sql"
	"log/slog"
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

// Common conversion logic
func convertToReport(
	id, userID, riskTypeID uuid.UUID,
	riskTypeName, riskTypeIconPath sql.NullString,
	riskTopicID uuid.NullUUID,
	riskTopicName, riskTopicIconPath sql.NullString,
	description sql.NullString,
	latitude, longitude float64,
	province, municipality, neighborhood, address, imageURL sql.NullString,
	status interface{},
	reviewedBy uuid.NullUUID,
	resolvedAt sql.NullTime,
	verificationCount, rejectionCount sql.NullInt32,
	expiresAt, createdAt, updatedAt sql.NullTime,
) *model.Report {
	reportStatus, ok := status.(model.ReportStatus)
	if !ok {
		reportStatus = model.ReportStatus("pending")
	}

	var reviewedByID uuid.UUID
	if reviewedBy.Valid {
		reviewedByID = reviewedBy.UUID
	}

	var resolvedAtTime time.Time
	if resolvedAt.Valid {
		resolvedAtTime = resolvedAt.Time
	}

	var expiresAtTime *time.Time
	if expiresAt.Valid {
		t := expiresAt.Time
		expiresAtTime = &t
	}

	var riskTypeIcon *string
	if riskTypeIconPath.Valid {
		riskTypeIcon = &riskTypeIconPath.String
	}

	var riskTopicIcon *string
	if riskTopicIconPath.Valid {
		riskTopicIcon = &riskTopicIconPath.String
	}

	return &model.Report{
		ID:                id,
		UserID:            userID,
		RiskTypeID:        riskTypeID,
		RiskTypeName:      riskTypeName.String,
		RiskTypeIconPath:  riskTypeIcon,
		RiskTopicID:       riskTopicID.UUID,
		RiskTopicName:     riskTopicName.String,
		RiskTopicIconPath: riskTopicIcon,
		Description:       description.String,
		Latitude:          latitude,
		Longitude:         longitude,
		Province:          province.String,
		Municipality:      municipality.String,
		Neighborhood:      neighborhood.String,
		Address:           address.String,
		ImageURL:          imageURL.String,
		Status:            reportStatus,
		ReviewedBy:        reviewedByID,
		ResolvedAt:        resolvedAtTime,
		VerificationCount: int(verificationCount.Int32),
		RejectionCount:    int(rejectionCount.Int32),
		ExpiresAt:         expiresAtTime,
		CreatedAt:         createdAt.Time,
		UpdatedAt:         updatedAt.Time,
	}
}

func getReportByIDRowToModel(r sqlc.GetReportByIDRow) *model.Report {
	return convertToReport(
		r.ID, r.UserID, r.RiskTypeID,
		r.RiskTypeName, r.RiskTypeIconPath,
		r.RiskTopicID,
		r.RiskTopicName, r.RiskTopicIconPath,
		r.Description,
		r.Latitude, r.Longitude,
		r.Province, r.Municipality, r.Neighborhood, r.Address, r.ImageUrl,
		r.Status,
		r.ReviewedBy,
		r.ResolvedAt,
		r.VerificationCount, r.RejectionCount,
		r.ExpiresAt, r.CreatedAt, r.UpdatedAt,
	)
}

func listReportsByStatusRowToModel(r sqlc.ListReportsByStatusRow) *model.Report {
	return convertToReport(
		r.ID, r.UserID, r.RiskTypeID,
		r.RiskTypeName, r.RiskTypeIconPath,
		r.RiskTopicID,
		r.RiskTopicName, r.RiskTopicIconPath,
		r.Description,
		r.Latitude, r.Longitude,
		r.Province, r.Municipality, r.Neighborhood, r.Address, r.ImageUrl,
		r.Status,
		r.ReviewedBy,
		r.ResolvedAt,
		r.VerificationCount, r.RejectionCount,
		r.ExpiresAt, r.CreatedAt, r.UpdatedAt,
	)
}

func listReportsByUserRowToModel(r sqlc.ListReportsByUserRow) *model.Report {
	return convertToReport(
		r.ID, r.UserID, r.RiskTypeID,
		r.RiskTypeName, r.RiskTypeIconPath,
		r.RiskTopicID,
		r.RiskTopicName, r.RiskTopicIconPath,
		r.Description,
		r.Latitude, r.Longitude,
		r.Province, r.Municipality, r.Neighborhood, r.Address, r.ImageUrl,
		r.Status,
		r.ReviewedBy,
		r.ResolvedAt,
		r.VerificationCount, r.RejectionCount,
		r.ExpiresAt, r.CreatedAt, r.UpdatedAt,
	)
}

func listReportsByIDsRowToModel(r sqlc.ListReportsByIDsRow) *model.Report {
	return convertToReport(
		r.ID, r.UserID, r.RiskTypeID,
		r.RiskTypeName, r.RiskTypeIconPath,
		r.RiskTopicID,
		r.RiskTopicName, r.RiskTopicIconPath,
		r.Description,
		r.Latitude, r.Longitude,
		r.Province, r.Municipality, r.Neighborhood, r.Address, r.ImageUrl,
		r.Status,
		r.ReviewedBy,
		r.ResolvedAt,
		r.VerificationCount, r.RejectionCount,
		r.ExpiresAt, r.CreatedAt, r.UpdatedAt,
	)
}

func listReportsWithPaginationRowToModel(r sqlc.ListReportsWithPaginationRow) *model.Report {
	return convertToReport(
		r.ID, r.UserID, r.RiskTypeID,
		r.RiskTypeName, r.RiskTypeIconPath,
		r.RiskTopicID,
		r.RiskTopicName, r.RiskTopicIconPath,
		r.Description,
		r.Latitude, r.Longitude,
		r.Province, r.Municipality, r.Neighborhood, r.Address, r.ImageUrl,
		r.Status,
		r.ReviewedBy,
		r.ResolvedAt,
		r.VerificationCount, r.RejectionCount,
		r.ExpiresAt, r.CreatedAt, r.UpdatedAt,
	)
}

func (r *ReportPG) Create(ctx context.Context, m *model.Report) error {
	// Creates o report no PostgreSQL e obtém o ID gerado
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
	return getReportByIDRowToModel(item), nil
}

func (r *ReportPG) ListByStatus(ctx context.Context, status model.ReportStatus) ([]*model.Report, error) {
	items, err := r.q.ListReportsByStatus(ctx, status)
	if err != nil {
		return nil, err
	}
	return mapSlice(items, listReportsByStatusRowToModel), nil
}

func (r *ReportPG) ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Report, error) {
	items, err := r.q.ListReportsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mapSlice(items, listReportsByUserRowToModel), nil
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

	return mapSlice(items, listReportsByIDsRowToModel), nil
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

	reports := mapSlice(items, listReportsWithPaginationRowToModel)
	return reports, int(total), nil
}

func (r *ReportPG) FindByRadiusWithDistance(ctx context.Context, lat float64, lon float64, radiusMeters float64, limit int) ([]repository.ReportWithDistance, error) {
	// 1. Busca reports com distâncias JÁ CALCULADAS e ORDENADAS
	// O Redis usa algoritmo otimizado em C e retorna resultados já ordenados por distância
	geoResults, err := r.locationStore.FindReportsInRadiusWithDistance(ctx, lat, lon, radiusMeters)
	if err != nil {
		slog.Error("failed to find reports with distance from redis", "error", err)
		return nil, err
	}

	// Se não encontrou nenhum report próximo, retorna lista vazia
	if len(geoResults) == 0 {
		slog.Info("no reports found in radius")
		return []repository.ReportWithDistance{}, nil
	}

	// 2. Aplica limite ANTES de buscar no PostgreSQL (economiza queries)
	if limit > 0 && limit < len(geoResults) {
		geoResults = geoResults[:limit]
	}

	// 3. Converte IDs para UUIDs e cria mapa de distâncias
	uuids := make([]uuid.UUID, 0, len(geoResults))
	distanceMap := make(map[string]float64, len(geoResults))

	for _, gr := range geoResults {
		reportUUID, err := uuid.Parse(gr.Member)
		if err != nil {
			slog.Warn("invalid report UUID from redis", "id", gr.Member, "error", err)
			continue
		}
		uuids = append(uuids, reportUUID)
		distanceMap[gr.Member] = gr.Distance
	}

	// 4. Busca dados completos em UMA ÚNICA query batch no PostgreSQL
	items, err := r.q.ListReportsByIDs(ctx, uuids)
	if err != nil {
		slog.Error("failed to find reports by ids", "error", err)
		return nil, err
	}

	// 5. Monta resultado final mantendo a ordem do Redis (já ordenada por distância)
	// Creates um mapa para lookup rápido O(1)
	reportMap := make(map[string]*model.Report, len(items))
	for _, item := range items {
		reportMap[item.ID.String()] = listReportsByIDsRowToModel(item)
	}

	// Reconstrói na ordem correta (ordem do Redis)
	result := make([]repository.ReportWithDistance, 0, len(geoResults))
	for _, gr := range geoResults {
		if report, exists := reportMap[gr.Member]; exists {
			result = append(result, repository.ReportWithDistance{
				Report:   report,
				Distance: gr.Distance,
			})
		}
	}

	slog.Info("found reports in radius with distances", "count", len(result))
	return result, nil
}

func mapSlice[T any, R any](in []T, fn func(T) *R) []*R {
	out := make([]*R, 0, len(in))
	for _, v := range in {
		out = append(out, fn(v))
	}
	return out
}

func (r *ReportPG) AddVote(ctx context.Context, vote *model.ReportVote) error {
	if vote.AnonymousSessionID != nil {
		return r.q.AddAnonymousReportVote(ctx, sqlc.AddAnonymousReportVoteParams{
			ReportID:           vote.ReportID,
			AnonymousSessionID: uuidPtrToNullUUID(vote.AnonymousSessionID),
			VoteType:           string(vote.VoteType),
		})
	}

	return r.q.AddUserReportVote(ctx, sqlc.AddUserReportVoteParams{
		ReportID: vote.ReportID,
		UserID:   uuidPtrToNullUUID(vote.UserID),
		VoteType: string(vote.VoteType),
	})
}

func (r *ReportPG) RemoveVote(ctx context.Context, reportID, userID uuid.UUID) error {
	return r.q.RemoveUserVote(ctx, sqlc.RemoveUserVoteParams{
		ReportID: reportID,
		UserID:   uuidToNullUUID(userID),
	})
}

func (r *ReportPG) RemoveAnonymousVote(ctx context.Context, reportID, sessionID uuid.UUID) error {
	return r.q.RemoveAnonymousVote(ctx, sqlc.RemoveAnonymousVoteParams{
		ReportID:           reportID,
		AnonymousSessionID: uuidToNullUUID(sessionID),
	})
}

func (r *ReportPG) GetUserVote(ctx context.Context, reportID, userID uuid.UUID) (*model.ReportVote, error) {
	vote, err := r.q.GetUserVote(ctx, sqlc.GetUserVoteParams{
		ReportID: reportID,
		UserID:   uuidToNullUUID(userID),
	})
	if err != nil {
		return nil, err
	}

	return &model.ReportVote{
		ID:       vote.ID,
		ReportID: vote.ReportID,
		UserID:   nullUUIDToPtr(vote.UserID),
		VoteType: model.VoteType(vote.VoteType),
	}, nil
}

func (r *ReportPG) GetAnonymousVote(ctx context.Context, reportID, sessionID uuid.UUID) (*model.ReportVote, error) {
	vote, err := r.q.GetAnonymousVote(ctx, sqlc.GetAnonymousVoteParams{
		ReportID:           reportID,
		AnonymousSessionID: uuidToNullUUID(sessionID),
	})
	if err != nil {
		return nil, err
	}

	return &model.ReportVote{
		ID:                 vote.ID,
		ReportID:           vote.ReportID,
		AnonymousSessionID: nullUUIDToPtr(vote.AnonymousSessionID),
		VoteType:           model.VoteType(vote.VoteType),
	}, nil
}

func (r *ReportPG) UpdateVerificationCounts(ctx context.Context, reportID uuid.UUID, upvotes, downvotes int) error {
	slog.Info("updating verification counts", "reportID", reportID, "upvotes", upvotes, "downvotes", downvotes)
	return nil
}

func (r *ReportPG) FindDuplicates(ctx context.Context, lat, lon float64, riskTypeID uuid.UUID, radiusMeters float64, since time.Time) ([]*model.Report, error) {
	items, err := r.q.FindDuplicateReports(ctx, sqlc.FindDuplicateReportsParams{
		RiskTypeID:    riskTypeID,
		CreatedAt:     sql.NullTime{Time: since, Valid: true},
		StMakepoint:   lat,
		StMakepoint_2: lon,
		StDwithin:     radiusMeters,
	})
	if err != nil {
		return nil, err
	}

	reports := make([]*model.Report, 0, len(items))
	for _, item := range items {
		reports = append(reports, findDuplicateReportsRowToModel(item))
	}
	return reports, nil
}

func findDuplicateReportsRowToModel(r sqlc.FindDuplicateReportsRow) *model.Report {
	return convertToReport(
		r.ID, r.UserID, r.RiskTypeID,
		r.RiskTypeName, r.RiskTypeIconPath,
		r.RiskTopicID,
		r.RiskTopicName, r.RiskTopicIconPath,
		r.Description,
		r.Latitude, r.Longitude,
		r.Province, r.Municipality, r.Neighborhood, r.Address, r.ImageUrl,
		r.Status,
		r.ReviewedBy,
		r.ResolvedAt,
		r.VerificationCount, r.RejectionCount,
		r.ExpiresAt, r.CreatedAt, r.UpdatedAt,
	)
}

func (r *ReportPG) ExpireOldReports(ctx context.Context, before time.Time) error {
	return r.q.ExpireOldReports(ctx, sql.NullTime{Time: before, Valid: true})
}

func (r *ReportPG) UpdateTrustScore(ctx context.Context, userID uuid.UUID, score int) error {
	slog.Info("updating trust score", "userID", userID, "score", score)
	return nil
}

func (r *ReportPG) IncrementReportsSubmitted(ctx context.Context, userID uuid.UUID) error {
	return r.q.IncrementUserReportsSubmitted(ctx, userID)
}

func (r *ReportPG) IncrementReportsVerified(ctx context.Context, userID uuid.UUID) error {
	return r.q.IncrementUserReportsVerified(ctx, userID)
}
