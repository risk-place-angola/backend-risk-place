package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type userLocationRepoPG struct {
	db *sql.DB
}

func NewUserLocationRepository(db *sql.DB) repository.UserLocationRepository {
	return &userLocationRepoPG{db: db}
}

func (r *userLocationRepoPG) Upsert(ctx context.Context, location *model.UserLocation) error {
	query := `
		INSERT INTO user_locations (
			id, user_id, device_id, latitude, longitude, speed, heading,
			avatar_id, color, is_anonymous, last_update
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (user_id) DO UPDATE SET
			device_id = EXCLUDED.device_id,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			speed = EXCLUDED.speed,
			heading = EXCLUDED.heading,
			last_update = EXCLUDED.last_update
	`

	slog.Debug("upserting user location",
		slog.String("user_id", location.UserID.String()),
		slog.Float64("lat", location.Latitude),
		slog.Float64("lon", location.Longitude),
		slog.Bool("is_anonymous", location.IsAnonymous))

	_, err := r.db.ExecContext(ctx, query,
		location.ID,
		location.UserID,
		location.DeviceID,
		location.Latitude,
		location.Longitude,
		location.Speed,
		location.Heading,
		location.AvatarID,
		location.Color,
		location.IsAnonymous,
		location.LastUpdate,
	)

	if err != nil {
		slog.Error("failed to upsert user location", slog.Any("error", err))
	}

	return err
}

func (r *userLocationRepoPG) FindByUserID(ctx context.Context, userID uuid.UUID) (*model.UserLocation, error) {
	query := `
		SELECT id, user_id, device_id, latitude, longitude, speed, heading,
		       avatar_id, color, is_anonymous, last_update, created_at
		FROM user_locations
		WHERE user_id = $1
	`
	var loc model.UserLocation
	var deviceID sql.NullString

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&loc.ID,
		&loc.UserID,
		&deviceID,
		&loc.Latitude,
		&loc.Longitude,
		&loc.Speed,
		&loc.Heading,
		&loc.AvatarID,
		&loc.Color,
		&loc.IsAnonymous,
		&loc.LastUpdate,
		&loc.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if deviceID.Valid {
		loc.DeviceID = deviceID.String
	}

	return &loc, nil
}

func (r *userLocationRepoPG) FindNearbyUsers(ctx context.Context, lat, lon, radiusMeters float64, limit int) ([]*model.UserLocation, error) {
	query := `
		SELECT DISTINCT ul.id, ul.user_id, ul.device_id, ul.latitude, ul.longitude, ul.speed, ul.heading,
		       ul.avatar_id, ul.color, ul.is_anonymous, ul.last_update, ul.created_at
		FROM user_locations ul
		LEFT JOIN user_safety_settings s ON (s.user_id = ul.user_id OR s.device_id = ul.device_id)
		WHERE ST_DWithin(
			ul.location,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
			$3
		)
		AND ul.last_update > NOW() - INTERVAL '30 seconds'
		AND (s.id IS NULL OR s.location_sharing_enabled = true)
		AND (s.id IS NULL OR s.show_online_status = true)
		ORDER BY ul.location <-> ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
		LIMIT $4
	`

	rows, err := r.db.QueryContext(ctx, query, lon, lat, radiusMeters, limit)
	if err != nil {
		slog.Error("failed to query nearby users", slog.Any("error", err))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", slog.Any("error", err))
		}
	}()

	var locations []*model.UserLocation
	for rows.Next() {
		var loc model.UserLocation
		var deviceID sql.NullString

		err := rows.Scan(
			&loc.ID,
			&loc.UserID,
			&deviceID,
			&loc.Latitude,
			&loc.Longitude,
			&loc.Speed,
			&loc.Heading,
			&loc.AvatarID,
			&loc.Color,
			&loc.IsAnonymous,
			&loc.LastUpdate,
			&loc.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if deviceID.Valid {
			loc.DeviceID = deviceID.String
		}

		locations = append(locations, &loc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *userLocationRepoPG) DeleteStale(ctx context.Context, thresholdSeconds int) error {
	query := `
		DELETE FROM user_locations
		WHERE last_update < NOW() - INTERVAL '1 second' * $1
	`
	_, err := r.db.ExecContext(ctx, query, thresholdSeconds)
	return err
}

func (r *userLocationRepoPG) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_locations WHERE user_id = $1`
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("user location not found")
	}

	return nil
}
