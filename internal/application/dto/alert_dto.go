package dto

type Alert struct {
	RiskTypeID       string  `json:"risk_type_id"`
	RiskTypeIconURL  *string `json:"risk_type_icon_url,omitempty"`
	RiskTopicID      string  `json:"risk_topic_id"`
	RiskTopicIconURL *string `json:"risk_topic_icon_url,omitempty"`
	UserID           string  `json:"user_id,omitempty"`
	Message          string  `json:"message"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Radius           float64 `json:"radius"`
	Severity         string  `json:"severity"`
}
