package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type anonymousSessionRepoPG struct {
	db *sql.DB
}

func NewAnonymousSessionRepository(db *sql.DB) *anonymousSessionRepoPG {
	return &anonymousSessionRepoPG{db: db}
}

func (r *anonymousSessionRepoPG) Create(ctx context.Context, session *model.AnonymousSession) error {
	query := `
		INSERT INTO anonymous_sessions (
			id, device_id, device_fcm_token, device_platform, device_model,
			latitude, longitude, alert_radius_meters, device_language,
			last_seen, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		session.ID,
		session.DeviceID,
		nullString(session.DeviceFCMToken),
		nullString(session.DevicePlatform),
		nullString(session.DeviceModel),
		nullFloat(session.Latitude),
		nullFloat(session.Longitude),
		session.AlertRadiusMeters,
		session.DeviceLanguage,
		session.LastSeen,
		session.CreatedAt,
		session.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create anonymous session: %w", err)
	}

	return nil
}

func (r *anonymousSessionRepoPG) FindByDeviceID(ctx context.Context, deviceID string) (*model.AnonymousSession, error) {
	query := `
		SELECT id, device_id, device_fcm_token, device_platform, device_model,
		       latitude, longitude, alert_radius_meters, device_language,
		       last_seen, created_at, updated_at
		FROM anonymous_sessions
		WHERE device_id = $1
	`

	var session model.AnonymousSession
	var fcmToken, platform, deviceModel sql.NullString
	var lat, lon sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, deviceID).Scan(
		&session.ID,
		&session.DeviceID,
		&fcmToken,
		&platform,
		&deviceModel,
		&lat,
		&lon,
		&session.AlertRadiusMeters,
		&session.DeviceLanguage,
		&session.LastSeen,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("anonymous session not found")
		}
		return nil, fmt.Errorf("failed to find anonymous session: %w", err)
	}

	session.DeviceFCMToken = fcmToken.String
	session.DevicePlatform = platform.String
	session.DeviceModel = deviceModel.String
	session.Latitude = lat.Float64
	session.Longitude = lon.Float64

	return &session, nil
}

func (r *anonymousSessionRepoPG) Update(ctx context.Context, session *model.AnonymousSession) error {
	query := `
		UPDATE anonymous_sessions
		SET device_fcm_token = $2,
		    device_platform = $3,
		    device_model = $4,
		    latitude = $5,
		    longitude = $6,
		    alert_radius_meters = $7,
		    device_language = $8,
		    last_seen = $9,
		    updated_at = $10
		WHERE device_id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		session.DeviceID,
		nullString(session.DeviceFCMToken),
		nullString(session.DevicePlatform),
		nullString(session.DeviceModel),
		nullFloat(session.Latitude),
		nullFloat(session.Longitude),
		session.AlertRadiusMeters,
		session.DeviceLanguage,
		session.LastSeen,
		session.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update anonymous session: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("anonymous session not found")
	}

	return nil
}

func (r *anonymousSessionRepoPG) UpdateLocation(ctx context.Context, deviceID string, lat, lon float64) error {
	query := `
		UPDATE anonymous_sessions
		SET latitude = $2,
		    longitude = $3,
		    last_seen = NOW(),
		    updated_at = NOW()
		WHERE device_id = $1
	`

	result, err := r.db.ExecContext(ctx, query, deviceID, lat, lon)
	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("anonymous session not found")
	}

	return nil
}

func (r *anonymousSessionRepoPG) UpdateFCMToken(ctx context.Context, deviceID string, fcmToken string) error {
	query := `
		UPDATE anonymous_sessions
		SET device_fcm_token = $2,
		    last_seen = NOW(),
		    updated_at = NOW()
		WHERE device_id = $1
	`

	result, err := r.db.ExecContext(ctx, query, deviceID, fcmToken)
	if err != nil {
		return fmt.Errorf("failed to update FCM token: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("anonymous session not found")
	}

	return nil
}

func (r *anonymousSessionRepoPG) GetFCMTokensInRadius(ctx context.Context, lat, lon, radiusMeters float64) ([]string, error) {
	query := `
		SELECT device_fcm_token
		FROM anonymous_sessions
		WHERE device_fcm_token IS NOT NULL
		  AND latitude IS NOT NULL
		  AND longitude IS NOT NULL
		  AND (
			6371000 * acos(
				cos(radians($1)) * cos(radians(latitude)) *
				cos(radians(longitude) - radians($2)) +
				sin(radians($1)) * sin(radians(latitude))
			)
		  ) <= alert_radius_meters
		  AND (
			6371000 * acos(
				cos(radians($1)) * cos(radians(latitude)) *
				cos(radians(longitude) - radians($2)) +
				sin(radians($1)) * sin(radians(latitude))
			)
		  ) <= $3
	`

	rows, err := r.db.QueryContext(ctx, query, lat, lon, radiusMeters)
	if err != nil {
		return nil, fmt.Errorf("failed to query FCM tokens: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close rows: %w", closeErr)
		}
	}()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, fmt.Errorf("failed to scan FCM token: %w", err)
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return tokens, nil
}

func (r *anonymousSessionRepoPG) Delete(ctx context.Context, deviceID string) error {
	query := `DELETE FROM anonymous_sessions WHERE device_id = $1`

	result, err := r.db.ExecContext(ctx, query, deviceID)
	if err != nil {
		return fmt.Errorf("failed to delete anonymous session: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("anonymous session not found")
	}

	return nil
}

func (r *anonymousSessionRepoPG) CleanupOldSessions(ctx context.Context, daysOld int) error {
	query := `DELETE FROM anonymous_sessions WHERE last_seen < NOW() - INTERVAL '1 day' * $1`

	_, err := r.db.ExecContext(ctx, query, daysOld)
	if err != nil {
		return fmt.Errorf("failed to cleanup old sessions: %w", err)
	}

	return nil
}

func (r *anonymousSessionRepoPG) TouchLastSeen(ctx context.Context, deviceID string) error {
	query := `
		UPDATE anonymous_sessions
		SET last_seen = NOW(),
		    updated_at = NOW()
		WHERE device_id = $1
	`

	result, err := r.db.ExecContext(ctx, query, deviceID)
	if err != nil {
		return fmt.Errorf("failed to touch last seen: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("anonymous session not found")
	}

	return nil
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func nullFloat(f float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: f, Valid: f != 0}
}
