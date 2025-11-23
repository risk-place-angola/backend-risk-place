package dto

type MyAlertResponse struct {
	ID               string  `json:"id"`
	RiskTypeName     string  `json:"risk_type_name"`
	RiskTypeIconURL  *string `json:"risk_type_icon_url,omitempty"`
	RiskTopicName    string  `json:"risk_topic_name,omitempty"`
	RiskTopicIconURL *string `json:"risk_topic_icon_url,omitempty"`
	Message          string  `json:"message"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Province         string  `json:"province,omitempty"`
	Municipality     string  `json:"municipality,omitempty"`
	Neighborhood     string  `json:"neighborhood,omitempty"`
	Address          string  `json:"address,omitempty"`
	RadiusMeters     int     `json:"radius_meters"`
	Status           string  `json:"status"`
	Severity         string  `json:"severity"`
	IsSubscribed     bool    `json:"is_subscribed,omitempty"`
	Subscribers      int     `json:"subscribers,omitempty"`
	CreatedAt        string  `json:"created_at"`
	ExpiresAt        string  `json:"expires_at,omitempty"`
	ResolvedAt       string  `json:"resolved_at,omitempty"`
}

type UpdateAlertInput struct {
	Message      string `json:"message"       validate:"required"`
	Severity     string `json:"severity"      validate:"required,oneof=low medium high critical"`
	RadiusMeters int    `json:"radius_meters" validate:"required,min=100,max=10000"`
}

type AlertSubscriptionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
