package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/api/option"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

var nc, _ = nats.Connect(nats.DefaultURL)

var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
var mongoClient, _ = mongo.Connect(ctx, clientOptions)
var messagesCollection = mongoClient.Database("chat_db").Collection("messages")

type Message struct {
	Room   string   `json:"room"`
	User   string   `json:"user"`
	Text   string   `json:"text"`
	SentAt int64    `json:"sent_at"`
	ReadBy []string `json:"read_by"`
}

type ReadStatus struct {
	Room  string `json:"room"`
	MsgID int64  `json:"msg_id"`
	User  string `json:"user"`
}

func main() {
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	go subscribeNATS()
	go subscribeReadStatus()
	go batchSaveToDB()

	select {}
}

func subscribeNATS() {
	_, err := nc.Subscribe("chat.*", func(msg *nats.Msg) {
		room := msg.Subject[len("chat."):]
		fmt.Printf("Received message in room %s: %s\n", room, string(msg.Data))

		var message Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			return
		}

		// Redis Pub/Sub을 통해 메시지 전송
		rdb.Publish(ctx, room, string(msg.Data))

		// Redis에 메시지 임시 저장
		rdb.RPush(ctx, "messages:"+room, string(msg.Data))

		// 클라이언트가 채팅방을 보고 있는지 확인 후 FCM 전송
		if !isUserOnline(room) {
			sendFCMNotification(room, string(msg.Data))
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}

func subscribeReadStatus() {
	pubsub := rdb.Subscribe(ctx, "read-status:*")
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}

		var readStatus ReadStatus
		json.Unmarshal([]byte(msg.Payload), &readStatus)

		// Redis에 읽음 상태 업데이트
		readKey := "read:" + readStatus.Room + ":" + fmt.Sprintf("%d", readStatus.MsgID)
		rdb.SAdd(ctx, readKey, readStatus.User)
		rdb.Expire(ctx, readKey, 24*time.Hour)

		// 읽음 상태를 실시간으로 다른 클라이언트에 전송
		readStatusBytes, _ := json.Marshal(readStatus)
		rdb.Publish(ctx, "read-status:"+readStatus.Room, readStatusBytes)
	}
}

func batchSaveToDB() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		rooms, err := rdb.Keys(ctx, "messages:*").Result()
		if err != nil {
			log.Printf("Error fetching rooms from Redis: %v", err)
			continue
		}

		for _, roomKey := range rooms {
			messages, err := rdb.LRange(ctx, roomKey, 0, -1).Result()
			if err != nil {
				log.Printf("Error fetching messages from Redis: %v", err)
				continue
			}

			if len(messages) == 0 {
				continue
			}

			var batch []interface{}
			for _, msgStr := range messages {
				var msg Message
				if err := json.Unmarshal([]byte(msgStr), &msg); err != nil {
					log.Printf("Error unmarshalling message: %v", err)
					continue
				}
				batch = append(batch, msg)
			}

			if len(batch) > 0 {
				_, err := messagesCollection.InsertMany(ctx, batch)
				if err != nil {
					log.Printf("Error inserting messages to MongoDB: %v", err)
				} else {
					// Redis에서 저장된 메시지 삭제
					rdb.Del(ctx, roomKey)
				}
			}
		}
	}
}

func isUserOnline(room string) bool {
	onlineUsers, err := rdb.SMembers(ctx, "online:"+room).Result()
	if err != nil || len(onlineUsers) == 0 {
		return false
	}
	return true
}

func sendFCMNotification(room string, message string) {
	// FCM 초기화
	opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Error initializing app: %v\n", err)
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("Error getting Messaging client: %v\n", err)
	}

	var msg Message
	json.Unmarshal([]byte(message), &msg)

	token := getUserDeviceToken(msg.User) // 사용자 디바이스 토큰 가져오기
	fcmMessage := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: "New Chat Message",
			Body:  msg.Text,
		},
	}
	_, err = client.Send(ctx, fcmMessage)
	if err != nil {
		log.Fatalf("Error sending FCM notification: %v\n", err)
	} else {
		fmt.Println("FCM notification sent successfully!")
	}
}

func getUserDeviceToken(user string) string {
	// 사용자 디바이스 토큰 조회 로직
	return "user_device_token"
}
