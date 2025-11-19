package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	domainErrors "github.com/risk-place-angola/backend-risk-place/internal/domain/errors"
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
			CreatedBy:    uuidPtrToNullUUID(alert.CreatedBy),
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

func (a alertRepoPG) GetByID(ctx context.Context, id uuid.UUID) (*model.Alert, error) {
	row, err := a.q.GetAlertByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainErrors.ErrAlertNotFound
		}
		return nil, err
	}

	return a.toDomain(row), nil
}

func (a alertRepoPG) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Alert, error) {
	rows, err := a.q.GetAlertsByUserID(ctx, uuidToNullUUID(userID))
	if err != nil {
		return nil, err
	}

	alerts := make([]*model.Alert, 0, len(rows))
	for _, row := range rows {
		alerts = append(alerts, a.toDomain(row))
	}

	return alerts, nil
}

func (a alertRepoPG) GetSubscribedAlerts(ctx context.Context, userID uuid.UUID) ([]*model.Alert, error) {
	rows, err := a.q.GetSubscribedAlerts(ctx, uuidToNullUUID(userID))
	if err != nil {
		return nil, err
	}

	alerts := make([]*model.Alert, 0, len(rows))
	for _, row := range rows {
		alerts = append(alerts, a.toDomain(row))
	}

	return alerts, nil
}

func (a alertRepoPG) Update(ctx context.Context, alert *model.Alert) error {
	if alert.RadiusMeters > math.MaxInt32 || alert.RadiusMeters < math.MinInt32 {
		return fmt.Errorf("radius meters out of range: must be between %d and %d", math.MinInt32, math.MaxInt32)
	}

	radiusMeters := int32(alert.RadiusMeters)

	return a.q.UpdateAlert(ctx, sqlc.UpdateAlertParams{
		ID:           alert.ID,
		Message:      alert.Message,
		Severity:     string(alert.Severity),
		RadiusMeters: radiusMeters,
		CreatedBy:    uuidPtrToNullUUID(alert.CreatedBy),
	})
}

func (a alertRepoPG) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return a.q.DeleteAlert(ctx, sqlc.DeleteAlertParams{
		ID:        id,
		CreatedBy: uuidToNullUUID(userID),
	})
}

func (a alertRepoPG) SubscribeToAlert(ctx context.Context, subscription *model.AlertSubscription) error {
	return a.q.SubscribeToAlert(ctx, sqlc.SubscribeToAlertParams{
		ID:           subscription.ID,
		AlertID:      subscription.AlertID,
		UserID:       uuidPtrToNullUUID(subscription.UserID),
		SubscribedAt: subscription.SubscribedAt,
	})
}

func (a alertRepoPG) UnsubscribeFromAlert(ctx context.Context, alertID, userID uuid.UUID) error {
	return a.q.UnsubscribeFromAlert(ctx, sqlc.UnsubscribeFromAlertParams{
		AlertID: alertID,
		UserID:  uuidToNullUUID(userID),
	})
}

func (a alertRepoPG) IsUserSubscribed(ctx context.Context, alertID, userID uuid.UUID) (bool, error) {
	result, err := a.q.IsUserSubscribed(ctx, sqlc.IsUserSubscribedParams{
		AlertID: alertID,
		UserID:  uuidToNullUUID(userID),
	})
	if err != nil {
		return false, err
	}
	return result, nil
}

func (a alertRepoPG) toDomain(row sqlc.Alert) *model.Alert {
	var status model.AlertStatus
	if s, ok := row.Status.(string); ok {
		status = model.AlertStatus(s)
	}

	var severity model.Severity
	if sev, ok := row.Severity.(string); ok {
		severity = model.Severity(sev)
	}

	alert := &model.Alert{
		ID:           row.ID,
		RiskTypeID:   row.RiskTypeID,
		Message:      row.Message,
		Latitude:     row.Latitude,
		Longitude:    row.Longitude,
		RadiusMeters: int(row.RadiusMeters),
		Status:       status,
		Severity:     severity,
	}

	if row.CreatedBy.Valid {
		createdBy := row.CreatedBy.UUID
		alert.CreatedBy = &createdBy
	}

	if row.RiskTopicID.Valid {
		alert.RiskTopicID = row.RiskTopicID.UUID
	}

	if row.Province.Valid {
		alert.Province = row.Province.String
	}

	if row.Municipality.Valid {
		alert.Municipality = row.Municipality.String
	}

	if row.Neighborhood.Valid {
		alert.Neighborhood = row.Neighborhood.String
	}

	if row.Address.Valid {
		alert.Address = row.Address.String
	}

	if row.CreatedAt.Valid {
		alert.CreatedAt = row.CreatedAt.Time
	}

	if row.ExpiresAt.Valid {
		alert.ExpiresAt = row.ExpiresAt.Time
	}

	if row.ResolvedAt.Valid {
		alert.ResolvedAt = row.ResolvedAt.Time
	}

	return alert
}

func NewAlertRepoPG(db *sql.DB) repository.AlertRepository {
	return &alertRepoPG{
		q: sqlc.New(db),
	}
}
