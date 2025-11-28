package risk

import (
	"context"

	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type RiskUseCase struct {
	riskTypesRepo  repository.RiskTypesRepository
	riskTopicsRepo repository.RiskTopicsRepository
	storageService port.StorageService
}

func NewRiskUseCase(
	riskTypesRepo repository.RiskTypesRepository,
	riskTopicsRepo repository.RiskTopicsRepository,
	storageService port.StorageService,
) *RiskUseCase {
	return &RiskUseCase{
		riskTypesRepo:  riskTypesRepo,
		riskTopicsRepo: riskTopicsRepo,
		storageService: storageService,
	}
}

// ListRiskTypes retrieves all risk types from the database
func (uc *RiskUseCase) ListRiskTypes(ctx context.Context) (*dto.RiskTypesListResponse, error) {
	riskTypes, err := uc.riskTypesRepo.ListRiskTypes(ctx)
	if err != nil {
		return nil, err
	}

	response := &dto.RiskTypesListResponse{
		Data: make([]dto.RiskTypeResponse, 0, len(riskTypes)),
	}

	for _, rt := range riskTypes {
		var iconURL *string
		if rt.IconPath != nil {
			url := uc.storageService.GetURL(*rt.IconPath)
			iconURL = &url
		}
		response.Data = append(response.Data, dto.RiskTypeResponse{
			ID:            rt.ID,
			Name:          rt.Name,
			Description:   rt.Description,
			IconURL:       iconURL,
			DefaultRadius: rt.DefaultRadiusMeters,
			IsEnabled:     rt.IsEnabled,
			CreatedAt:     rt.CreatedAt,
			UpdatedAt:     rt.UpdatedAt,
		})
	}

	return response, nil
}

// GetRiskType retrieves a specific risk type by ID
func (uc *RiskUseCase) GetRiskType(ctx context.Context, id string) (*dto.RiskTypeResponse, error) {
	riskType, err := uc.riskTypesRepo.GetRiskTypeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var iconURL *string
	if riskType.IconPath != nil {
		url := uc.storageService.GetURL(*riskType.IconPath)
		iconURL = &url
	}

	response := &dto.RiskTypeResponse{
		ID:            riskType.ID,
		Name:          riskType.Name,
		Description:   riskType.Description,
		IconURL:       iconURL,
		DefaultRadius: riskType.DefaultRadiusMeters,
		IsEnabled:     riskType.IsEnabled,
		CreatedAt:     riskType.CreatedAt,
		UpdatedAt:     riskType.UpdatedAt,
	}

	return response, nil
}

// ListRiskTopics retrieves risk topics, optionally filtered by risk type ID
func (uc *RiskUseCase) ListRiskTopics(ctx context.Context, riskTypeID *string) (*dto.RiskTopicsListResponse, error) {
	riskTopics, err := uc.riskTopicsRepo.ListRiskTopics(ctx, riskTypeID)
	if err != nil {
		return nil, err
	}

	response := &dto.RiskTopicsListResponse{
		Data: make([]dto.RiskTopicResponse, 0, len(riskTopics)),
	}

	for _, rt := range riskTopics {
		var iconURL *string
		if rt.IconPath != nil {
			url := uc.storageService.GetURL(*rt.IconPath)
			iconURL = &url
		}
		response.Data = append(response.Data, dto.RiskTopicResponse{
			ID:          rt.ID,
			RiskTypeID:  rt.RiskTypeID,
			Name:        rt.Name,
			Description: rt.Description,
			IconURL:     iconURL,
			CreatedAt:   rt.CreatedAt,
			UpdatedAt:   rt.UpdatedAt,
		})
	}

	return response, nil
}

// GetRiskTopic retrieves a specific risk topic by ID
func (uc *RiskUseCase) GetRiskTopic(ctx context.Context, id string) (*dto.RiskTopicResponse, error) {
	riskTopic, err := uc.riskTopicsRepo.GetRiskTopicByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var iconURL *string
	if riskTopic.IconPath != nil {
		url := uc.storageService.GetURL(*riskTopic.IconPath)
		iconURL = &url
	}

	response := &dto.RiskTopicResponse{
		ID:          riskTopic.ID,
		RiskTypeID:  riskTopic.RiskTypeID,
		Name:        riskTopic.Name,
		Description: riskTopic.Description,
		IconURL:     iconURL,
		CreatedAt:   riskTopic.CreatedAt,
		UpdatedAt:   riskTopic.UpdatedAt,
	}

	return response, nil
}

func (uc *RiskUseCase) UpdateRiskTypeIcon(ctx context.Context, id string, iconPath string) error {
	return uc.riskTypesRepo.UpdateRiskTypeIcon(ctx, id, iconPath)
}

func (uc *RiskUseCase) UpdateRiskTopicIcon(ctx context.Context, id string, iconPath string) error {
	return uc.riskTopicsRepo.UpdateRiskTopicIcon(ctx, id, iconPath)
}

func (uc *RiskUseCase) UpdateRiskTypeIsEnabled(ctx context.Context, id string, isEnabled bool) error {
	return uc.riskTypesRepo.UpdateRiskTypeIsEnabled(ctx, id, isEnabled)
}
