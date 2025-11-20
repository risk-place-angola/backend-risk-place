package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type RiskTypePG struct {
	q sqlc.Querier
}

func (r *RiskTypePG) CreateRiskType(ctx context.Context, name string, description string, defaultRadiusMeters *int) error {
	var radiusInt32 sql.NullInt32

	if defaultRadiusMeters != nil {
		if *defaultRadiusMeters > math.MaxInt32 || *defaultRadiusMeters < math.MinInt32 {
			return fmt.Errorf("default radius meters out of range: must be between %d and %d", math.MinInt32, math.MaxInt32)
		}
		radiusInt32 = sql.NullInt32{Int32: int32(*defaultRadiusMeters), Valid: true} // #nosec G115
	}

	return r.q.CreateRiskType(ctx, sqlc.CreateRiskTypeParams{
		Name:                name,
		Description:         sql.NullString{String: description, Valid: description != ""},
		DefaultRadiusMeters: radiusInt32,
	})
}

func (r *RiskTypePG) ListRiskTypes(ctx context.Context) ([]model.RiskType, error) {
	riskTypes, err := r.q.ListRiskTypes(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.RiskType, 0, len(riskTypes))
	for _, rt := range riskTypes {
		var iconPath *string
		if rt.IconPath.Valid {
			iconPath = &rt.IconPath.String
		}
		result = append(result, model.RiskType{
			ID:                  rt.ID,
			Name:                rt.Name,
			Description:         rt.Description.String,
			IconPath:            iconPath,
			DefaultRadiusMeters: int(rt.DefaultRadiusMeters.Int32),
			CreatedAt:           rt.CreatedAt.Time,
			UpdatedAt:           rt.UpdatedAt.Time,
		})
	}

	return result, nil
}

func (r *RiskTypePG) GetRiskTypeByID(ctx context.Context, id string) (model.RiskType, error) {
	rt, err := r.q.GetRiskTypeByID(ctx, uuid.MustParse(id))
	if err != nil {
		return model.RiskType{}, err
	}

	var iconPath *string
	if rt.IconPath.Valid {
		iconPath = &rt.IconPath.String
	}

	return model.RiskType{
		ID:                  rt.ID,
		Name:                rt.Name,
		Description:         rt.Description.String,
		IconPath:            iconPath,
		DefaultRadiusMeters: int(rt.DefaultRadiusMeters.Int32),
		CreatedAt:           rt.CreatedAt.Time,
		UpdatedAt:           rt.UpdatedAt.Time,
	}, nil
}

func (r *RiskTypePG) UpdateRiskType(ctx context.Context, id string, name string, description string, defaultRadiusMeters int) error {
	if defaultRadiusMeters > math.MaxInt32 || defaultRadiusMeters < math.MinInt32 {
		return fmt.Errorf("default radius meters out of range: must be between %d and %d", math.MinInt32, math.MaxInt32)
	}
	radiusInt32 := int32(defaultRadiusMeters) // #nosec G115

	return r.q.UpdateRiskType(ctx, sqlc.UpdateRiskTypeParams{
		ID:                  uuid.MustParse(id),
		Name:                name,
		Description:         sql.NullString{String: description, Valid: description != ""},
		DefaultRadiusMeters: sql.NullInt32{Int32: radiusInt32, Valid: true},
	})
}

func (r *RiskTypePG) UpdateRiskTypeIcon(ctx context.Context, id string, iconPath string) error {
	return r.q.UpdateRiskTypeIcon(ctx, sqlc.UpdateRiskTypeIconParams{
		ID:       uuid.MustParse(id),
		IconPath: sql.NullString{String: iconPath, Valid: true},
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
