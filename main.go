package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v7"
)

var (
	RETRY_THREADSHOLD = 3
	DEAD_LETTER_CHAN  = "DeadLetterQueue"
)

type PushNotificationModel struct {
	To           string           `json:"to"`
	Notification NotificationData `json:"notification"`
}

type NotificationData struct {
	Body  string `json:"body"`
	Title string `json:"title"`
}

type Message struct {
	RetryCount int
	ErrorMsg   string
	Data       PushNotificationModel
}

func pushMessage(notificationPayload *PushNotificationModel) (bool, string) {
	//check if the message has reached its re try threshold and place in dead letter queue if so

	//else use BLPLUSH or ROPLPUBSH to Dequeue from the ActiveJobQueue
	//move to the ProcessingJobQueue,

	//if the job is successful
	//either drop into MOngoDB hard store, or just remove the message from the process queue
	//if the job not ificaiton fails
	//update the retry count and move it back to activeJobQueue
	pushManager := &PushManager{}

	success, errorMsg := pushManager.PushMessage(
		notificationPayload.Notification.Title,
		notificationPayload.Notification.Body,
		notificationPayload.To)

	return success, errorMsg

}
func main() {
	data := PushNotificationModel{
		To: "dkfjdlskfjlk234",
		Notification: NotificationData{
			Body:  "HEYS",
			Title: "FUCKOFF",
		},
	}
	pushMessage(&data)
	redisUri := os.Getenv("redisUri")
	opt, _ := redis.ParseURL(redisUri)
	c := redis.NewClient(opt) /*&redis.Options{
		Addr: "localhost:6379",
	}*/
	const key = "myJobQueue"

	workChan := make(chan string)
	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal)
	for i := 1; i < 9; i++ {
		go func(id int, c chan string, r *redis.Client) {
			for {
				work, ok := <-c
				//BEFORE doing work , lets put this item into this goroutines OWN backup channel/list

				//PROCESS the message, if its successcode OK
				//REM all items from the back up list

				//IF FAIL
				//Increment the retry in the message, and requeue it into the jobQueue
				//if retry has hit 3 times on the current increment , send it to the dead letter Queue hard store to analyze
				if !ok {
					fmt.Println("Channel Closed~ exiting GoRoutine-", id)
					break
				}

				backupChan := fmt.Sprintf("backupChan-%v", id)
				r.RPush(backupChan, work) //push to the backup Channel

				fmt.Printf("GoRoutine-%d prcoessing %v\n", id, work)
				//deserialize the object
				var msg Message

				err := json.Unmarshal([]byte(work), &msg)
				if err != nil {
					fmt.Errorf(err.Error())
					panic(err)
				}

				//attempt to send the notification here
				isSuccess, errMsg := pushMessage(&msg.Data)

				if !isSuccess {
					//error sending push notification

					//check for retry count
					msg.RetryCount++
					msg.ErrorMsg = errMsg
					if msg.RetryCount > RETRY_THREADSHOLD {
						//send to dead letter queue too many retries
						//r.BRPopLPush(backupChan, DEAD_LETTER_CHAN, 10*time.Second)
						jsonStr, _ := json.Marshal(msg)
						r.RPush(DEAD_LETTER_CHAN, jsonStr)

					} else {
						jsonStr, _ := json.Marshal(msg)
						r.RPush(key, jsonStr)
						//r.BRPopLPush(backupChan, key, 10*time.Second)
					}
					r.Del(backupChan)
				} else {
					//successful push notificaiton remove from the temp store
					r.Del(backupChan)
				}

				fmt.Printf("Payload %v", msg)
			}

		}(i, workChan, c)
	}

	fmt.Println("Waiting for jobs on jobQueue: ", key)
	go func() {
		for {
			result, err := c.BLPop(0*time.Second, key).Result()

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Executing job: ", result[1])
			workChan <- result[1] //BY PRICIPAL ALWAYS have PUB close the channel, never allow anyone else esp SUB to close channels
		}
	}()
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	fmt.Println("Exit signal recieved")
	close(workChan)

	/*
		go func() {

			for i := 1; i < 3; i++ {
				go func(id int) {
					fmt.Printf("gorutine ID: %d working\n", id)
					for {
						go func(id int) {
							processQueue := fmt.Sprintf("%d-processQueue", id)
							timeout := 0 * time.Second
							result, err := c.BRPopLPush(key, processQueue, timeout).Result()
							if err != nil {
								fmt.Errorf(err.Error())
								panic(err)
							}
							fmt.Printf("go-routine %v", result)
						}(id)
					}
				}(i)

			}
		}()

	*/
	/*
			for {
				result, err := c.BLPop(0*time.Second, key).Result()

				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Executing job: ", result[1])
			}
		}()

	*/

	// block for ever, used for testing only
	//select {}

	//i want a thread pool going multiple go routines listing and readying my nigga
}

func PushIntoDeadLetterQueue(r *redis.Client, msg Message) {
	jsonStr, _ := json.Marshal(msg)
	r.RPush(DEAD_LETTER_CHAN, jsonStr)
}
