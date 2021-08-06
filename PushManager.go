package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redklouds/go-notification-service-mq/interfaces"
	"github.com/redklouds/go-notification-service-mq/models"
)

type PushManager struct {
	interfaces.IPushNotitificationManager
}

func (mgr *PushManager) PushMessage(Title string, Body string, DeviceId string) (bool, string) /*error*/ {

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
		if res.StatusCode != 200 {
			//log an error
		}
		panic(err)

	}
	defer res.Body.Close()

	resposneBody, _ := ioutil.ReadAll(res.Body)

	var pushResponse models.PushNotificationResponse

	json.Unmarshal(resposneBody, &pushResponse)

	fmt.Printf("%v", pushResponse)

	if pushResponse.Failure >= 1 {
		fmt.Println("Failure occurred")

		//e := errors.New("DSfsdfdss")

		//e1 := fmt.Errorf("Dfsfs")
		//return e
		return false, pushResponse.Results[0].Error
	} else {
		fmt.Println("Success")
		return true, ""
	}
}
