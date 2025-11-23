package dto

type UpdateDeviceInfoRequest struct {
	DeviceFCMToken string `json:"device_fcm_token,omitempty"`
	DeviceLanguage string `json:"device_language,omitempty" validate:"omitempty,oneof=pt en"`
}

type NotificationPreferencesRequest struct {
	PushEnabled bool `json:"push_enabled"`
	SMSEnabled  bool `json:"sms_enabled"`
}

type NotificationPreferencesResponse struct {
	PushEnabled bool `json:"push_enabled"`
	SMSEnabled  bool `json:"sms_enabled"`
}
