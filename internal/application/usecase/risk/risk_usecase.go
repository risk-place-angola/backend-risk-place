package risk

import (
	"context"

	"github.com/risk-place-angola/backend-risk-place/internal/application/dto"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type RiskUseCase struct {
	riskTypesRepo  repository.RiskTypesRepository
	riskTopicsRepo repository.RiskTopicsRepository
}

func NewRiskUseCase(
	riskTypesRepo repository.RiskTypesRepository,
	riskTopicsRepo repository.RiskTopicsRepository,
) *RiskUseCase {
	return &RiskUseCase{
		riskTypesRepo:  riskTypesRepo,
		riskTopicsRepo: riskTopicsRepo,
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
		response.Data = append(response.Data, dto.RiskTypeResponse{
			ID:            rt.ID,
			Name:          rt.Name,
			Description:   rt.Description,
			DefaultRadius: rt.DefaultRadiusMeters,
			CreatedAt:     rt.CreatedAt,
			UpdatedAt:     rt.UpdatedAt,
		})
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
		response.Data = append(response.Data, dto.RiskTopicResponse{
			ID:          rt.ID,
			RiskTypeID:  rt.RiskTypeID,
			Name:        rt.Name,
			Description: rt.Description,
			CreatedAt:   rt.CreatedAt,
			UpdatedAt:   rt.UpdatedAt,
		})
	}

	return response, nil
}
