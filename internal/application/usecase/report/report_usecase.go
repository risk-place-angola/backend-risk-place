package report

import (
	"context"
	"log/slog"
	"time"

	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/event"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type ReportUseCase struct {
	repo            repository.ReportRepository
	riskTypesRepo   repository.RiskTypesRepository
	riskTopicsRepo  repository.RiskTopicsRepository
	locationStore   port.LocationStore
	geoService      port.GeolocationService
	eventDispatcher port.EventDispatcher
}

func NewReportUseCase(
	repo repository.ReportRepository,
	eventDispatcher port.EventDispatcher,
	geoService port.GeolocationService,
	riskTypesRepo repository.RiskTypesRepository,
	riskTopicsRepo repository.RiskTopicsRepository,
	locationStore port.LocationStore,
) *ReportUseCase {
	return &ReportUseCase{
		repo:            repo,
		eventDispatcher: eventDispatcher,
		geoService:      geoService,
		riskTypesRepo:   riskTypesRepo,
		riskTopicsRepo:  riskTopicsRepo,
		locationStore:   locationStore,
	}
}

func (uc *ReportUseCase) Create(ctx context.Context, dto dto.ReportCreate) (*model.Report, error) {
	err := uc.geoService.ValidateCoordinates(dto.Latitude, dto.Longitude)
	if err != nil {
		slog.Error("invalid coordinates", "error", err)
		return nil, err
	}

	riskType, err := uc.riskTypesRepo.GetRiskTypeByID(ctx, dto.RiskTypeID)
	if err != nil {
		slog.Error("failed to get risk type", "error", err)
		return nil, err
	}

	riskTopicID := uuid.MustParse(dto.RiskTopicID)
	riskTopic, err := uc.riskTopicsRepo.GetRiskTopicByID(ctx, dto.RiskTopicID)
	if err != nil {
		slog.Error("failed to get risk topic", "error", err)
		return nil, err
	}

	report := &model.Report{
		ID:           uuid.New(),
		RiskTypeID:   uuid.MustParse(dto.RiskTypeID),
		RiskTopicID:  riskTopicID,
		Description:  dto.Description,
		Province:     dto.Province,
		Municipality: dto.Municipality,
		Neighborhood: dto.Neighborhood,
		Address:      dto.Address,
		ImageURL:     dto.ImageURL,
		UserID:       uuid.MustParse(dto.UserID),
		Latitude:     dto.Latitude,
		Longitude:    dto.Longitude,
		Status:       model.ReportStatusPending,
		IsPrivate:    riskTopic.IsSensitive,
		CreatedAt:    time.Now(),
	}

	err = uc.repo.Create(ctx, report)
	if err != nil {
		slog.Error("failed to create report", "error", err)
		return nil, err
	}

	userIDs, err := uc.locationStore.FindUsersInRadius(
		ctx, report.Latitude, report.Longitude, float64(riskType.DefaultRadiusMeters),
	)
	if err != nil {
		slog.Error("failed to find users in radius", "error", err)
		return nil, err
	}

	uuidUserIDs := make([]uuid.UUID, 0, len(userIDs))
	for _, uid := range userIDs {
		err := uc.repo.CreateReportNotification(ctx, report.ID, uuid.MustParse(uid))
		if err != nil {
			slog.Error("failed to create report notification", "error", err, "user_id", uid)
			continue
		}
		slog.Info("created report notification", "report_id", report.ID, "user_id", uid)
		uuidUserIDs = append(uuidUserIDs, uuid.MustParse(uid))
	}

	uc.eventDispatcher.Dispatch(event.ReportCreatedEvent{
		ReportID:  report.ID,
		UserID:    uuidUserIDs,
		Message:   report.Description,
		Latitude:  report.Latitude,
		Longitude: report.Longitude,
		Radius:    float64(riskType.DefaultRadiusMeters),
		RiskType:  riskType.Name,
	})

	return report, nil
}

func (uc *ReportUseCase) Verify(ctx context.Context, reportID string) error {
	report, err := uc.repo.GetByID(ctx, uuid.MustParse(reportID))
	if err != nil {
		return err
	}

	uc.eventDispatcher.Dispatch(event.ReportVerifiedEvent{
		ReportID: report.ID,
		UserID:   report.UserID,
	})

	return nil
}

func (uc *ReportUseCase) Resolve(ctx context.Context, reportID, moderatorID string) error {
	report, err := uc.repo.GetByID(ctx, uuid.MustParse(reportID))
	if err != nil {
		return err
	}

	riskType, _ := uc.riskTypesRepo.GetRiskTypeByID(ctx, report.RiskTypeID.String())

	userIDs, _ := uc.locationStore.FindUsersInRadius(
		ctx, report.Latitude, report.Longitude,
		float64(riskType.DefaultRadiusMeters),
	)

	uc.eventDispatcher.Dispatch(event.ReportResolvedEvent{
		ReportID: report.ID,
		Message:  "Situação foi resolvida",
		UserIDs:  userIDs,
	})

	return nil
}

func (uc *ReportUseCase) List(ctx context.Context, params dto.ListReportsQueryParams) (*dto.ListReportsResponse, error) {
	const (
		defaultPage  = 1
		defaultLimit = 20
		maxLimit     = 100
	)

	// Validate and set defaults
	if params.Page <= 0 {
		params.Page = defaultPage
	}
	if params.Limit <= 0 {
		params.Limit = defaultLimit
	}
	if params.Limit > maxLimit {
		params.Limit = maxLimit
	}
	if params.Order == "" {
		params.Order = "desc"
	}
	if params.Sort == "" {
		params.Sort = "created_at"
	}

	// Call repository
	reports, total, err := uc.repo.ListWithPagination(ctx, repository.ListReportsParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Status: params.Status,
		Sort:   params.Sort,
		Order:  params.Order,
	})
	if err != nil {
		slog.Error("failed to list reports with pagination", "error", err)
		return nil, err
	}

	// Convert to DTOs
	reportDTOs := make([]dto.ReportDTO, 0, len(reports))
	for _, report := range reports {
		reportDTOs = append(reportDTOs, dto.ReportToDTO(report))
	}

	// Calculate pagination metadata
	totalPages := (total + params.Limit - 1) / params.Limit
	hasMore := params.Page < totalPages
	hasPrevious := params.Page > 1

	response := &dto.ListReportsResponse{
		Reports: reportDTOs,
		Pagination: dto.PaginationMetadata{
			Page:        params.Page,
			Limit:       params.Limit,
			Total:       total,
			TotalPages:  totalPages,
			HasMore:     hasMore,
			HasPrevious: hasPrevious,
		},
	}

	return response, nil
}

func (uc *ReportUseCase) ListNearby(ctx context.Context, lat, lon, radius float64) ([]*model.Report, error) {
	err := uc.geoService.ValidateCoordinates(lat, lon)
	if err != nil {
		return nil, err
	}
	return uc.repo.FindByRadius(ctx, lat, lon, radius)
}

func (uc *ReportUseCase) ListNearbyWithDistance(ctx context.Context, params dto.NearbyReportsQueryParams) (*dto.NearbyReportsResponse, error) {
	// Validate coordinates
	err := uc.geoService.ValidateCoordinates(params.Latitude, params.Longitude)
	if err != nil {
		slog.Error("invalid coordinates", "error", err)
		return nil, err
	}

	// Set defaults
	if params.Radius <= 0 {
		params.Radius = 500 // 500 meters default
	}
	if params.Limit <= 0 {
		params.Limit = 50
	}

	// Call repository with distance calculation
	reportsWithDist, err := uc.repo.FindByRadiusWithDistance(ctx, params.Latitude, params.Longitude, params.Radius, params.Limit)
	if err != nil {
		slog.Error("failed to find reports nearby with distance", "error", err)
		return nil, err
	}

	// Convert to DTOs
	reportDTOs := make([]dto.ReportWithDistance, 0, len(reportsWithDist))
	for _, rwd := range reportsWithDist {
		reportDTOs = append(reportDTOs, dto.ReportToDTOWithDistance(rwd.Report, rwd.Distance))
	}

	response := &dto.NearbyReportsResponse{
		Reports: reportDTOs,
	}

	return response, nil
}

func (uc *ReportUseCase) UpdateLocation(ctx context.Context, reportID string, req dto.UpdateReportLocationRequest) error {
	id, err := uuid.Parse(reportID)
	if err != nil {
		slog.Error("invalid report ID", "reportID", reportID, "error", err)
		return err
	}

	err = uc.geoService.ValidateCoordinates(req.Latitude, req.Longitude)
	if err != nil {
		slog.Error("invalid coordinates", "error", err)
		return err
	}

	report, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.Error("failed to get report", "reportID", reportID, "error", err)
		return err
	}

	err = uc.repo.UpdateLocation(ctx, id, repository.UpdateLocationParams{
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Address:      req.Address,
		Neighborhood: req.Neighborhood,
		Municipality: req.Municipality,
		Province:     req.Province,
	})
	if err != nil {
		slog.Error("failed to update report location", "reportID", reportID, "error", err)
		return err
	}

	slog.Info("report location updated", "reportID", reportID, "userID", report.UserID)

	return nil
}
