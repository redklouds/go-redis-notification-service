package models

type PushNotificationResponse struct {
	MultiCastID int64        `json:"multicast_id"`
	Success     int          `json:"success"`
	Failure     int          `json:"failure"`
	CanonicalID int          `json:"canonical_ids"`
	Results     []ResultInfo `json:"results"`
}

type ResultInfo struct {
	MessageID string `json:"message_id"`
	Error     string `json:"error"`
}
