package dto

import "github.com/google/uuid"

type CreateEmergencyContactInput struct {
	Name       string `json:"name"        validate:"required"`
	Phone      string `json:"phone"       validate:"required"`
	Relation   string `json:"relation"    validate:"required,oneof=family friend colleague neighbor other"`
	IsPriority bool   `json:"is_priority"`
}

type UpdateEmergencyContactInput struct {
	Name       string `json:"name"        validate:"required"`
	Phone      string `json:"phone"       validate:"required"`
	Relation   string `json:"relation"    validate:"required,oneof=family friend colleague neighbor other"`
	IsPriority bool   `json:"is_priority"`
}

type EmergencyContactResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Relation   string `json:"relation"`
	IsPriority bool   `json:"is_priority"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type EmergencyAlertInput struct {
	Latitude  float64 `json:"latitude"  validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	Message   string  `json:"message"`
}

type EmergencyAlertResponse struct {
	Success          bool     `json:"success"`
	ContactsNotified int      `json:"contacts_notified"`
	NotifiedContacts []string `json:"notified_contacts"`
	Message          string   `json:"message"`
}

type EmergencyAlertNotification struct {
	UserName  string  `json:"user_name"`
	UserPhone string  `json:"user_phone"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Message   string  `json:"message"`
	Timestamp string  `json:"timestamp"`
	MapLink   string  `json:"map_link"`
}

func ToEmergencyContactResponse(userID uuid.UUID, contact interface{}) EmergencyContactResponse {
	v, ok := contact.(map[string]interface{})
	if !ok {
		return EmergencyContactResponse{UserID: userID.String()}
	}

	getString := func(key string) string {
		if val, ok := v[key].(string); ok {
			return val
		}
		return ""
	}

	getBool := func(key string) bool {
		if val, ok := v[key].(bool); ok {
			return val
		}
		return false
	}

	return EmergencyContactResponse{
		ID:         getString("id"),
		UserID:     userID.String(),
		Name:       getString("name"),
		Phone:      getString("phone"),
		Relation:   getString("relation"),
		IsPriority: getBool("is_priority"),
		CreatedAt:  getString("created_at"),
		UpdatedAt:  getString("updated_at"),
	}
}