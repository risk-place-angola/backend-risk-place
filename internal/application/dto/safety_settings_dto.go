package dto

import "time"

type SafetySettingsResponse struct {
	ID     string `example:"550e8400-e29b-41d4-a716-446655440000" json:"id"`
	UserID string `example:"550e8400-e29b-41d4-a716-446655440001" json:"user_id"`

	NotificationsEnabled         bool     `example:"true"          json:"notifications_enabled"`
	NotificationAlertTypes       []string `example:"high,critical" json:"notification_alert_types"`
	NotificationAlertRadiusMins  int      `example:"1000"          json:"notification_alert_radius_mins"`
	NotificationReportTypes      []string `example:"verified"      json:"notification_report_types"`
	NotificationReportRadiusMins int      `example:"500"           json:"notification_report_radius_mins"`

	LocationSharingEnabled bool `example:"false" json:"location_sharing_enabled"`
	LocationHistoryEnabled bool `example:"true"  json:"location_history_enabled"`

	ProfileVisibility string `enums:"public,friends,private" example:"public"          json:"profile_visibility"`
	AnonymousReports  bool   `example:"false"                                          json:"anonymous_reports"`
	ShowOnlineStatus  bool   `example:"true"                                           json:"show_online_status"`

	AutoAlertsEnabled      bool   `example:"false" json:"auto_alerts_enabled"`
	DangerZonesEnabled     bool   `example:"true"  json:"danger_zones_enabled"`
	TimeBasedAlertsEnabled bool   `example:"false" json:"time_based_alerts_enabled"`
	HighRiskStartTime      string `example:"22:00" json:"high_risk_start_time"`
	HighRiskEndTime        string `example:"06:00" json:"high_risk_end_time"`

	NightModeEnabled   bool   `example:"false" json:"night_mode_enabled"`
	NightModeStartTime string `example:"22:00" json:"night_mode_start_time"`
	NightModeEndTime   string `example:"06:00" json:"night_mode_end_time"`

	CreatedAt time.Time `example:"2024-01-15T10:30:00Z" json:"created_at"`
	UpdatedAt time.Time `example:"2024-01-15T10:30:00Z" json:"updated_at"`
}

type UpdateSafetySettingsInput struct {
	NotificationsEnabled         *bool     `json:"notifications_enabled,omitempty"`
	NotificationAlertTypes       *[]string `json:"notification_alert_types,omitempty"`
	NotificationAlertRadiusMins  *int      `json:"notification_alert_radius_mins,omitempty"`
	NotificationReportTypes      *[]string `json:"notification_report_types,omitempty"`
	NotificationReportRadiusMins *int      `json:"notification_report_radius_mins,omitempty"`

	LocationSharingEnabled *bool `json:"location_sharing_enabled,omitempty"`
	LocationHistoryEnabled *bool `json:"location_history_enabled,omitempty"`

	ProfileVisibility *string `enums:"public,friends,private"      json:"profile_visibility,omitempty"`
	AnonymousReports  *bool   `json:"anonymous_reports,omitempty"`
	ShowOnlineStatus  *bool   `json:"show_online_status,omitempty"`

	AutoAlertsEnabled      *bool   `json:"auto_alerts_enabled,omitempty"`
	DangerZonesEnabled     *bool   `json:"danger_zones_enabled,omitempty"`
	TimeBasedAlertsEnabled *bool   `json:"time_based_alerts_enabled,omitempty"`
	HighRiskStartTime      *string `example:"22:00"                            json:"high_risk_start_time,omitempty"`
	HighRiskEndTime        *string `example:"06:00"                            json:"high_risk_end_time,omitempty"`

	NightModeEnabled   *bool   `json:"night_mode_enabled,omitempty"`
	NightModeStartTime *string `example:"22:00"                     json:"night_mode_start_time,omitempty"`
	NightModeEndTime   *string `example:"06:00"                     json:"night_mode_end_time,omitempty"`
}
