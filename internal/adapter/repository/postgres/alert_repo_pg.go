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

	return a.getAlertByIDRowToModel(row), nil
}

func (a alertRepoPG) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Alert, error) {
	rows, err := a.q.GetAlertsByUserID(ctx, uuidToNullUUID(userID))
	if err != nil {
		return nil, err
	}

	alerts := make([]*model.Alert, 0, len(rows))
	for _, row := range rows {
		alerts = append(alerts, a.getAlertsByUserIDRowToModel(row))
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
		alerts = append(alerts, a.getSubscribedAlertsRowToModel(row))
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

// convertToAlert is a generic helper function to convert any alert row type to domain model
func (a alertRepoPG) convertToAlert(
	id uuid.UUID,
	createdBy uuid.NullUUID,
	anonymousSessionID uuid.NullUUID,
	deviceID sql.NullString,
	riskTypeID uuid.UUID,
	riskTopicID uuid.NullUUID,
	message string,
	latitude float64,
	longitude float64,
	province sql.NullString,
	municipality sql.NullString,
	neighborhood sql.NullString,
	address sql.NullString,
	radiusMeters int32,
	severity interface{},
	status interface{},
	createdAt sql.NullTime,
	expiresAt sql.NullTime,
	resolvedAt sql.NullTime,
	riskTypeName sql.NullString,
	riskTypeIconPath sql.NullString,
	riskTopicName sql.NullString,
	riskTopicIconPath sql.NullString,
) *model.Alert {
	var alertStatus model.AlertStatus
	if s, ok := status.(string); ok {
		alertStatus = model.AlertStatus(s)
	}

	var alertSeverity model.Severity
	if sev, ok := severity.(string); ok {
		alertSeverity = model.Severity(sev)
	}

	alert := &model.Alert{
		ID:           id,
		RiskTypeID:   riskTypeID,
		Message:      message,
		Latitude:     latitude,
		Longitude:    longitude,
		RadiusMeters: int(radiusMeters),
		Status:       alertStatus,
		Severity:     alertSeverity,
	}

	if createdBy.Valid {
		cb := createdBy.UUID
		alert.CreatedBy = &cb
	}

	if anonymousSessionID.Valid {
		asid := anonymousSessionID.UUID
		alert.AnonymousSessionID = &asid
	}

	if deviceID.Valid {
		did := deviceID.String
		alert.DeviceID = &did
	}

	if riskTopicID.Valid {
		alert.RiskTopicID = riskTopicID.UUID
	}

	if province.Valid {
		alert.Province = province.String
	}

	if municipality.Valid {
		alert.Municipality = municipality.String
	}

	if neighborhood.Valid {
		alert.Neighborhood = neighborhood.String
	}

	if address.Valid {
		alert.Address = address.String
	}

	if createdAt.Valid {
		alert.CreatedAt = createdAt.Time
	}

	if expiresAt.Valid {
		alert.ExpiresAt = expiresAt.Time
	}

	if resolvedAt.Valid {
		alert.ResolvedAt = resolvedAt.Time
	}

	if riskTypeName.Valid {
		alert.RiskTypeName = riskTypeName.String
	}

	if riskTypeIconPath.Valid {
		alert.RiskTypeIconPath = &riskTypeIconPath.String
	}

	if riskTopicName.Valid {
		alert.RiskTopicName = riskTopicName.String
	}

	if riskTopicIconPath.Valid {
		alert.RiskTopicIconPath = &riskTopicIconPath.String
	}

	return alert
}

// getAlertByIDRowToModel converts GetAlertByIDRow to domain model
func (a alertRepoPG) getAlertByIDRowToModel(row sqlc.GetAlertByIDRow) *model.Alert {
	return a.convertToAlert(
		row.ID,
		row.CreatedBy,
		row.AnonymousSessionID,
		row.DeviceID,
		row.RiskTypeID,
		row.RiskTopicID,
		row.Message,
		row.Latitude,
		row.Longitude,
		row.Province,
		row.Municipality,
		row.Neighborhood,
		row.Address,
		row.RadiusMeters,
		row.Severity,
		row.Status,
		row.CreatedAt,
		row.ExpiresAt,
		row.ResolvedAt,
		row.RiskTypeName,
		row.RiskTypeIconPath,
		row.RiskTopicName,
		row.RiskTopicIconPath,
	)
}

// getAlertsByUserIDRowToModel converts GetAlertsByUserIDRow to domain model
func (a alertRepoPG) getAlertsByUserIDRowToModel(row sqlc.GetAlertsByUserIDRow) *model.Alert {
	return a.convertToAlert(
		row.ID,
		row.CreatedBy,
		row.AnonymousSessionID,
		row.DeviceID,
		row.RiskTypeID,
		row.RiskTopicID,
		row.Message,
		row.Latitude,
		row.Longitude,
		row.Province,
		row.Municipality,
		row.Neighborhood,
		row.Address,
		row.RadiusMeters,
		row.Severity,
		row.Status,
		row.CreatedAt,
		row.ExpiresAt,
		row.ResolvedAt,
		row.RiskTypeName,
		row.RiskTypeIconPath,
		row.RiskTopicName,
		row.RiskTopicIconPath,
	)
}

// getSubscribedAlertsRowToModel converts GetSubscribedAlertsRow to domain model
func (a alertRepoPG) getSubscribedAlertsRowToModel(row sqlc.GetSubscribedAlertsRow) *model.Alert {
	return a.convertToAlert(
		row.ID,
		row.CreatedBy,
		row.AnonymousSessionID,
		row.DeviceID,
		row.RiskTypeID,
		row.RiskTopicID,
		row.Message,
		row.Latitude,
		row.Longitude,
		row.Province,
		row.Municipality,
		row.Neighborhood,
		row.Address,
		row.RadiusMeters,
		row.Severity,
		row.Status,
		row.CreatedAt,
		row.ExpiresAt,
		row.ResolvedAt,
		row.RiskTypeName,
		row.RiskTypeIconPath,
		row.RiskTopicName,
		row.RiskTopicIconPath,
	)
}

func NewAlertRepoPG(db *sql.DB) repository.AlertRepository {
	return &alertRepoPG{
		q: sqlc.New(db),
	}
}
