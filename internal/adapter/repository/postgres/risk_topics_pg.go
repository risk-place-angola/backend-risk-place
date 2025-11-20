package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type RiskTopicPG struct {
	q sqlc.Querier
}

func (r *RiskTopicPG) ListRiskTopics(ctx context.Context, riskTypeID *string) ([]model.RiskTopic, error) {
	var riskTopics []sqlc.RiskTopic
	var err error

	if riskTypeID != nil {
		parsedID, parseErr := uuid.Parse(*riskTypeID)
		if parseErr != nil {
			return nil, parseErr
		}
		riskTopics, err = r.q.ListRiskTopicsByType(ctx, parsedID)
	} else {
		riskTopics, err = r.q.ListRiskTopics(ctx)
	}

	if err != nil {
		return nil, err
	}

	result := make([]model.RiskTopic, 0, len(riskTopics))
	for _, rt := range riskTopics {
		var description *string
		if rt.Description.Valid {
			description = &rt.Description.String
		}
		var iconPath *string
		if rt.IconPath.Valid {
			iconPath = &rt.IconPath.String
		}

		result = append(result, model.RiskTopic{
			ID:          rt.ID,
			RiskTypeID:  rt.RiskTypeID,
			Name:        rt.Name,
			Description: description,
			IconPath:    iconPath,
			CreatedAt:   rt.CreatedAt.Time,
			UpdatedAt:   rt.UpdatedAt.Time,
		})
	}

	return result, nil
}

func (r *RiskTopicPG) GetRiskTopicByID(ctx context.Context, id string) (model.RiskTopic, error) {
	rt, err := r.q.GetRiskTopicByID(ctx, uuid.MustParse(id))
	if err != nil {
		return model.RiskTopic{}, err
	}

	var description *string
	if rt.Description.Valid {
		description = &rt.Description.String
	}
	var iconPath *string
	if rt.IconPath.Valid {
		iconPath = &rt.IconPath.String
	}

	return model.RiskTopic{
		ID:          rt.ID,
		RiskTypeID:  rt.RiskTypeID,
		Name:        rt.Name,
		Description: description,
		IconPath:    iconPath,
		CreatedAt:   rt.CreatedAt.Time,
		UpdatedAt:   rt.UpdatedAt.Time,
	}, nil
}

func (r *RiskTopicPG) UpdateRiskTopicIcon(ctx context.Context, id string, iconPath string) error {
	return r.q.UpdateRiskTopicIcon(ctx, sqlc.UpdateRiskTopicIconParams{
		ID:       uuid.MustParse(id),
		IconPath: sql.NullString{String: iconPath, Valid: true},
	})
}

func NewRiskTopicRepoPG(db *sql.DB) repository.RiskTopicsRepository {
	return &RiskTopicPG{
		q: sqlc.New(db),
	}
}
