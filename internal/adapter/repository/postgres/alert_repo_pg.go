package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type alertRepoPG struct {
	q sqlc.Querier
}

func (a alertRepoPG) Create(ctx context.Context, alert *model.Alert) error {
	if alert.RadiusMeters > math.MaxInt32 || alert.RadiusMeters < math.MinInt32 {
		return fmt.Errorf("radius meters out of range: must be between %d and %d", math.MinInt32, math.MaxInt32)
	}

	radiusMeters := int32(alert.RadiusMeters) // #nosec G115

	return a.q.CreateAlert(ctx,
		sqlc.CreateAlertParams{
			ID:           alert.ID,
			CreatedBy:    uuid.NullUUID{UUID: alert.CreatedBy, Valid: true},
			RiskTypeID:   alert.RiskTypeID,
			RiskTopicID:  uuid.NullUUID{UUID: alert.RiskTopicID, Valid: true},
			Message:      alert.Message,
			Latitude:     alert.Latitude,
			Longitude:    alert.Longitude,
			Province:     sql.NullString{String: alert.Province, Valid: true},
			Municipality: sql.NullString{String: alert.Municipality, Valid: true},
			Neighborhood: sql.NullString{String: alert.Neighborhood, Valid: true},
			Address:      sql.NullString{String: alert.Address, Valid: true},
			Severity:     string(alert.Severity),
			RadiusMeters: radiusMeters,
			ExpiresAt:    sql.NullTime{Time: alert.ExpiresAt, Valid: alert.ExpiresAt != time.Time{}},
		})
}

func (a alertRepoPG) CreateAlertNotification(ctx context.Context, alertID uuid.UUID, userID string) error {
	return a.q.CreateAlertNotification(ctx,
		sqlc.CreateAlertNotificationParams{
			ReferenceID: alertID,
			UserID:      uuid.MustParse(userID),
			Type:        "alert",
		})
}

func NewAlertRepoPG(db *sql.DB) repository.AlertRepository {
	return &alertRepoPG{
		q: sqlc.New(db),
	}
}
