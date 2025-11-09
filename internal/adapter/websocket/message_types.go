package websocket

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type UpdateLocationPayload struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type AlertNotification struct {
	AlertID   string  `json:"alert_id"`
	Message   string  `json:"message"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
}

type ReportNotification struct {
	ReportID  string  `json:"report_id"`
	Message   string  `json:"message"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
