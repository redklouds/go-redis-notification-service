package models

type PushNotificationModel struct {
	To           string           `json:"to"`
	Notification NotificationData `json:"notification"`
}

type NotificationData struct {
	Body  string `json:"body"`
	Title string `json:"title"`
}
