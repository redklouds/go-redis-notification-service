package main

import (
	"bytes"
	"os"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/redklouds/go-notification-service-mq/interfaces"
	"github.com/redklouds/go-notification-service-mq/models"
)

type PushManager struct {
	interfaces.IPushNotitificationManager
}

func (mgr *PushManager) PushMessage(Title string, Body string, DeviceId string) {

	url := "https://fcm.googleapis.com/fcm/send"

	client := &http.Client{}
	pushSerKey := os.Getenv("pushServerKey")
	serverKey := pushSerKey

	payload := models.PushNotificationModel{
		To: DeviceId,
		Notification: models.NotificationData{
			Body:  Body,
			Title: Title,
		},
	}

	json_data, err := json.Marshal(payload)

	if err != nil {
		fmt.Println("Error serializing push payload: ", err)
	}

	r, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(json_data))

	r.Header.Add("Authorization", "key="+serverKey)
	r.Header.Add("Content-Type", "application/json")

	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	resposneBody, _ := ioutil.ReadAll(res.Body)

	fmt.Println(resposneBody)

	var pushResponse models.PushNotificationResponse

	json.Unmarshal(resposneBody, &pushResponse)

	fmt.Printf("%v", pushResponse)

	if pushResponse.Failure >= 1 {
		fmt.Println("Failure occurred")
	} else {
		fmt.Println("Success")
	}

}
