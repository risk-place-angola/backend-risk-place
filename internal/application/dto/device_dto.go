package dto

type RegisterDeviceRequest struct {
	DeviceID          string  `json:"device_id"                     validate:"required,min=16"`
	FCMToken          string  `json:"fcm_token,omitempty"`
	Platform          string  `json:"platform"                      validate:"omitempty,oneof=ios android web"`
	Model             string  `json:"model,omitempty"`
	Language          string  `json:"language,omitempty"`
	Latitude          float64 `json:"latitude,omitempty"`
	Longitude         float64 `json:"longitude,omitempty"`
	AlertRadiusMeters int     `json:"alert_radius_meters,omitempty"`
}

type UpdateDeviceLocationRequest struct {
	DeviceID  string  `json:"device_id" validate:"required"`
	Latitude  float64 `json:"latitude"  validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type DeviceResponse struct {
	DeviceID          string  `json:"device_id"`
	FCMToken          string  `json:"fcm_token,omitempty"`
	Platform          string  `json:"platform,omitempty"`
	Latitude          float64 `json:"latitude,omitempty"`
	Longitude         float64 `json:"longitude,omitempty"`
	AlertRadiusMeters int     `json:"alert_radius_meters"`
	Message           string  `json:"message,omitempty"`
}
