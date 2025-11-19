package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	defaultAlertRadiusMins  = 1000
	defaultReportRadiusMins = 500
)

type ProfileVisibility string

const (
	ProfileVisibilityPublic  ProfileVisibility = "public"
	ProfileVisibilityFriends ProfileVisibility = "friends"
	ProfileVisibilityPrivate ProfileVisibility = "private"
)

type SafetySettings struct {
	ID                 uuid.UUID
	UserID             *uuid.UUID // Nullable - either UserID OR (AnonymousSessionID + DeviceID)
	AnonymousSessionID *uuid.UUID // Set for anonymous users
	DeviceID           *string    // Set for anonymous users

	NotificationsEnabled         bool
	NotificationAlertTypes       []string
	NotificationAlertRadiusMins  int
	NotificationReportTypes      []string
	NotificationReportRadiusMins int

	LocationSharingEnabled bool
	LocationHistoryEnabled bool

	ProfileVisibility ProfileVisibility
	AnonymousReports  bool
	ShowOnlineStatus  bool

	AutoAlertsEnabled      bool
	DangerZonesEnabled     bool
	TimeBasedAlertsEnabled bool
	HighRiskStartTime      time.Time
	HighRiskEndTime        time.Time

	NightModeEnabled   bool
	NightModeStartTime time.Time
	NightModeEndTime   time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsAnonymous returns true if these settings belong to an anonymous user
func (s *SafetySettings) IsAnonymous() bool {
	return s.AnonymousSessionID != nil && s.DeviceID != nil
}

// IsAuthenticated returns true if these settings belong to an authenticated user
func (s *SafetySettings) IsAuthenticated() bool {
	return s.UserID != nil
}

// NewSafetySettings creates safety settings for an authenticated user
func NewSafetySettings(userID uuid.UUID) (*SafetySettings, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user_id is required")
	}

	now := time.Now()

	highRiskStart, _ := time.Parse("15:04", "22:00")
	highRiskEnd, _ := time.Parse("15:04", "06:00")
	nightModeStart, _ := time.Parse("15:04", "22:00")
	nightModeEnd, _ := time.Parse("15:04", "06:00")

	return &SafetySettings{
		ID:                           uuid.New(),
		UserID:                       &userID,
		NotificationsEnabled:         true,
		NotificationAlertTypes:       []string{"high", "critical"},
		NotificationAlertRadiusMins:  defaultAlertRadiusMins,
		NotificationReportTypes:      []string{"verified"},
		NotificationReportRadiusMins: defaultReportRadiusMins,
		LocationSharingEnabled:       false,
		LocationHistoryEnabled:       true,
		ProfileVisibility:            ProfileVisibilityPublic,
		AnonymousReports:             false,
		ShowOnlineStatus:             true,
		AutoAlertsEnabled:            false,
		DangerZonesEnabled:           true,
		TimeBasedAlertsEnabled:       false,
		HighRiskStartTime:            highRiskStart,
		HighRiskEndTime:              highRiskEnd,
		NightModeEnabled:             false,
		NightModeStartTime:           nightModeStart,
		NightModeEndTime:             nightModeEnd,
		CreatedAt:                    now,
		UpdatedAt:                    now,
	}, nil
}

// NewAnonymousSafetySettings creates safety settings for an anonymous user
func NewAnonymousSafetySettings(anonymousSessionID uuid.UUID, deviceID string) (*SafetySettings, error) {
	if anonymousSessionID == uuid.Nil {
		return nil, errors.New("anonymous_session_id is required")
	}

	if deviceID == "" {
		return nil, errors.New("device_id is required")
	}

	now := time.Now()

	highRiskStart, _ := time.Parse("15:04", "22:00")
	highRiskEnd, _ := time.Parse("15:04", "06:00")
	nightModeStart, _ := time.Parse("15:04", "22:00")
	nightModeEnd, _ := time.Parse("15:04", "06:00")

	return &SafetySettings{
		ID:                           uuid.New(),
		AnonymousSessionID:           &anonymousSessionID,
		DeviceID:                     &deviceID,
		NotificationsEnabled:         true,
		NotificationAlertTypes:       []string{"high", "critical"},
		NotificationAlertRadiusMins:  defaultAlertRadiusMins,
		NotificationReportTypes:      []string{"verified"},
		NotificationReportRadiusMins: defaultReportRadiusMins,
		LocationSharingEnabled:       false,
		LocationHistoryEnabled:       true,
		ProfileVisibility:            ProfileVisibilityPublic,
		AnonymousReports:             false,
		ShowOnlineStatus:             true,
		AutoAlertsEnabled:            false,
		DangerZonesEnabled:           true,
		TimeBasedAlertsEnabled:       false,
		HighRiskStartTime:            highRiskStart,
		HighRiskEndTime:              highRiskEnd,
		NightModeEnabled:             false,
		NightModeStartTime:           nightModeStart,
		NightModeEndTime:             nightModeEnd,
		CreatedAt:                    now,
		UpdatedAt:                    now,
	}, nil
}

func (s *SafetySettings) Validate() error {
	// Validate that either UserID OR (AnonymousSessionID + DeviceID) is set
	hasUser := s.UserID != nil && *s.UserID != uuid.Nil
	hasAnonymous := s.AnonymousSessionID != nil && *s.AnonymousSessionID != uuid.Nil && s.DeviceID != nil && *s.DeviceID != ""

	if !hasUser && !hasAnonymous {
		return errors.New("either user_id or (anonymous_session_id + device_id) is required")
	}

	if hasUser && hasAnonymous {
		return errors.New("cannot have both user_id and anonymous_session_id set")
	}

	if s.NotificationAlertRadiusMins < 100 || s.NotificationAlertRadiusMins > 10000 {
		return errors.New("notification_alert_radius_mins must be between 100 and 10000")
	}

	if s.NotificationReportRadiusMins < 100 || s.NotificationReportRadiusMins > 10000 {
		return errors.New("notification_report_radius_mins must be between 100 and 10000")
	}

	validVisibilities := map[ProfileVisibility]bool{
		ProfileVisibilityPublic:  true,
		ProfileVisibilityFriends: true,
		ProfileVisibilityPrivate: true,
	}

	if !validVisibilities[s.ProfileVisibility] {
		return errors.New("invalid profile_visibility value")
	}

	return nil
}
