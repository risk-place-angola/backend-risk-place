package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres/sqlc"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type locationSharingRepoPG struct {
	q sqlc.Querier
}

func NewLocationSharingRepository(db *sql.DB) repository.LocationSharingRepository {
	return &locationSharingRepoPG{q: sqlc.New(db)}
}

func (r *locationSharingRepoPG) Save(ctx context.Context, entity *model.LocationSharing) error {
	var userID, anonymousSessionID uuid.NullUUID
	var deviceID sql.NullString

	if entity.UserID != nil {
		userID = uuid.NullUUID{UUID: *entity.UserID, Valid: true}
	}
	if entity.AnonymousSessionID != nil {
		anonymousSessionID = uuid.NullUUID{UUID: *entity.AnonymousSessionID, Valid: true}
	}
	if entity.DeviceID != nil {
		deviceID = sql.NullString{String: *entity.DeviceID, Valid: true}
	}

	// Validate duration to prevent overflow
	if entity.DurationMinutes > 2147483647 || entity.DurationMinutes < -2147483648 {
		return errors.New("duration_minutes out of range for int32")
	}

	return r.q.CreateLocationSharing(ctx, sqlc.CreateLocationSharingParams{
		ID:                 entity.ID,
		UserID:             userID,
		AnonymousSessionID: anonymousSessionID,
		DeviceID:           deviceID,
		OwnerName:          sql.NullString{String: entity.OwnerName, Valid: entity.OwnerName != ""},
		Token:              entity.Token,
		Latitude:           entity.Latitude,
		Longitude:          entity.Longitude,
		DurationMinutes:    int32(entity.DurationMinutes), // #nosec G115 - validated above
		ExpiresAt:          entity.ExpiresAt,
		LastUpdatedAt:      entity.LastUpdatedAt,
		IsActive:           entity.IsActive,
	})
}

func (r *locationSharingRepoPG) Update(ctx context.Context, entity *model.LocationSharing) error {
	return r.q.UpdateLocationSharing(ctx, sqlc.UpdateLocationSharingParams{
		ID:            entity.ID,
		Latitude:      entity.Latitude,
		Longitude:     entity.Longitude,
		LastUpdatedAt: entity.LastUpdatedAt,
		IsActive:      entity.IsActive,
		UpdatedAt:     entity.UpdatedAt,
	})
}

func (r *locationSharingRepoPG) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.DeleteLocationSharing(ctx, uid)
}

func (r *locationSharingRepoPG) FindByID(ctx context.Context, id uuid.UUID) (*model.LocationSharing, error) {
	row, err := r.q.GetLocationSharingByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.rowToModel(row), nil
}

func (r *locationSharingRepoPG) FindAll(ctx context.Context) ([]*model.LocationSharing, error) {
	rows, err := r.q.ListAllLocationSharings(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.LocationSharing, 0, len(rows))
	for _, row := range rows {
		result = append(result, r.rowToModel(row))
	}
	return result, nil
}

func (r *locationSharingRepoPG) FindByToken(ctx context.Context, token string) (*model.LocationSharing, error) {
	row, err := r.q.GetLocationSharingByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return r.rowToModel(row), nil
}

func (r *locationSharingRepoPG) FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*model.LocationSharing, error) {
	rows, err := r.q.ListActiveLocationSharingsByUserID(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	result := make([]*model.LocationSharing, 0, len(rows))
	for _, row := range rows {
		result = append(result, r.rowToModel(row))
	}
	return result, nil
}

func (r *locationSharingRepoPG) FindActiveByDeviceID(ctx context.Context, deviceID string) ([]*model.LocationSharing, error) {
	rows, err := r.q.ListActiveLocationSharingsByDeviceID(ctx, sql.NullString{String: deviceID, Valid: true})
	if err != nil {
		return nil, err
	}

	result := make([]*model.LocationSharing, 0, len(rows))
	for _, row := range rows {
		result = append(result, r.rowToModel(row))
	}
	return result, nil
}

func (r *locationSharingRepoPG) DeactivateExpired(ctx context.Context) error {
	return r.q.DeactivateExpiredLocationSharings(ctx, time.Now())
}

func (r *locationSharingRepoPG) rowToModel(row sqlc.LocationSharing) *model.LocationSharing {
	var userID, anonymousSessionID *uuid.UUID
	var deviceID *string
	var ownerName string

	if row.UserID.Valid {
		userID = &row.UserID.UUID
	}
	if row.AnonymousSessionID.Valid {
		anonymousSessionID = &row.AnonymousSessionID.UUID
	}
	if row.DeviceID.Valid {
		deviceID = &row.DeviceID.String
	}
	if row.OwnerName.Valid {
		ownerName = row.OwnerName.String
	}

	return &model.LocationSharing{
		ID:                 row.ID,
		UserID:             userID,
		AnonymousSessionID: anonymousSessionID,
		DeviceID:           deviceID,
		OwnerName:          ownerName,
		Token:              row.Token,
		Latitude:           row.Latitude,
		Longitude:          row.Longitude,
		DurationMinutes:    int(row.DurationMinutes),
		ExpiresAt:          row.ExpiresAt,
		LastUpdatedAt:      row.LastUpdatedAt,
		IsActive:           row.IsActive,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
}
