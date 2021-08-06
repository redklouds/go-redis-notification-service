package models

type Message struct {
	RetryCount int
	ErrorMsg   string
	Data       PushNotificationModel
}
