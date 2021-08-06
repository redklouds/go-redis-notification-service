package interfaces

type IPushNotitificationManager interface {
	PushMessage(Title string, Body string, DeviceId string)

	//PushMesages -> push to an array of device IDs
}
