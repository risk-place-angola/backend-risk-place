package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type RiskTypePG struct {
	q sqlc.Querier
}

func (r *RiskTypePG) CreateRiskType(ctx context.Context, name string, description string, defaultRadiusMeters *int) error {
	return r.q.CreateRiskType(ctx, sqlc.CreateRiskTypeParams{
		Name:                name,
		Description:         sql.NullString{String: description, Valid: description != ""},
		DefaultRadiusMeters: sql.NullInt32{Int32: int32(*defaultRadiusMeters), Valid: defaultRadiusMeters != nil},
	})
}

func (r *RiskTypePG) ListRiskTypes(ctx context.Context) ([]model.RiskType, error) {
	riskTypes, err := r.q.ListRiskTypes(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.RiskType, 0, len(riskTypes))
	for _, rt := range riskTypes {
		result = append(result, model.RiskType{
			ID:                  rt.ID,
			Name:                rt.Name,
			Description:         rt.Description.String,
			DefaultRadiusMeters: int(rt.DefaultRadiusMeters.Int32),
		})
	}

	return result, nil
}

func (r *RiskTypePG) GetRiskTypeByID(ctx context.Context, id string) (model.RiskType, error) {
	rt, err := r.q.GetRiskTypeByID(ctx, uuid.MustParse(id))
	if err != nil {
		return model.RiskType{}, err
	}

	return model.RiskType{
		ID:                  rt.ID,
		Name:                rt.Name,
		Description:         rt.Description.String,
		DefaultRadiusMeters: int(rt.DefaultRadiusMeters.Int32),
	}, nil
}

func (r *RiskTypePG) UpdateRiskType(ctx context.Context, id string, name string, description string, defaultRadiusMeters int) error {
	return r.q.UpdateRiskType(ctx, sqlc.UpdateRiskTypeParams{
		ID:                  uuid.MustParse(id),
		Name:                name,
		Description:         sql.NullString{String: description, Valid: description != ""},
		DefaultRadiusMeters: sql.NullInt32{Int32: int32(defaultRadiusMeters), Valid: true},
	})
}

func (r *RiskTypePG) DeleteRiskType(ctx context.Context, id string) error {
	return r.q.DeleteRiskType(ctx, uuid.MustParse(id))
}

func NewRiskTypeRepoPG(db *sql.DB) repository.RiskTypesRepository {
	return &RiskTypePG{
		q: sqlc.New(db),
	}
}
