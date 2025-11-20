package websocket

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type UpdateLocationPayload struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed"`
	Heading   float64 `json:"heading"`
	Radius    float64 `json:"radius"`
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

type NearbyUserResponse struct {
	UserID    string  `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	AvatarID  string  `json:"avatar_id"`
	Color     string  `json:"color"`
	Speed     float64 `json:"speed"`
	Heading   float64 `json:"heading"`
}

type NearbyUsersData struct {
	Users      []NearbyUserResponse `json:"users"`
	Radius     float64              `json:"radius"`
	TotalCount int                  `json:"total_count"`
}
