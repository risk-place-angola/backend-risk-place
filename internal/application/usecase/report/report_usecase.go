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
	locationStore   port.LocationStore
	geoService      port.GeolocationService
	eventDispatcher port.EventDispatcher
}

func NewReportUseCase(
	repo repository.ReportRepository,
	eventDispatcher port.EventDispatcher,
	geoService port.GeolocationService,
	riskTypesRepo repository.RiskTypesRepository,
	locationStore port.LocationStore,
) *ReportUseCase {
	return &ReportUseCase{
		repo:            repo,
		eventDispatcher: eventDispatcher,
		geoService:      geoService,
		riskTypesRepo:   riskTypesRepo,
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

	report := &model.Report{
		ID:           uuid.New(),
		RiskTypeID:   uuid.MustParse(dto.RiskTypeID),
		RiskTopicID:  uuid.MustParse(dto.RiskTopicID),
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
			slog.Error("failed to create report notification", "error", err)
			return nil, err
		}
		slog.Info("created report notification", "report_id", report.ID)
		uuidUserIDs = append(uuidUserIDs, uuid.MustParse(uid))
	}

	uc.eventDispatcher.Dispatch(event.ReportCreatedEvent{
		ReportID:  report.ID,
		UserID:    uuidUserIDs,
		Message:   report.Description,
		Latitude:  report.Latitude,
		Longitude: report.Longitude,
		Radius:    float64(riskType.DefaultRadiusMeters),
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

func (uc *ReportUseCase) ListNearby(ctx context.Context, lat, lon, radius float64) ([]*model.Report, error) {
	err := uc.geoService.ValidateCoordinates(lat, lon)
	if err != nil {
		return nil, err
	}
	return uc.repo.FindByRadius(ctx, lat, lon, radius)
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
