package models

type MessageStatus int

const (
	Queued MessageStatus = iota
	Processing
	Failed
)

type Message struct {
	RetryCount int
	ErrorMsg   string
	Data       PushNotificationModel
	Status     MessageStatus
}
