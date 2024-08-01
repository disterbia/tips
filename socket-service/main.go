package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

var nc, _ = nats.Connect(nats.DefaultURL)

var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
var mongoClient, _ = mongo.Connect(ctx, clientOptions)
var messagesCollection = mongoClient.Database("chat_db").Collection("messages")

type Client struct {
	conn *websocket.Conn
	room string
	user string
}

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

var clients = make(map[*Client]bool)

func main() {
	app := fiber.New()

	app.Get("/ws/:room/:user", websocket.New(func(c *websocket.Conn) {
		room := c.Params("room")
		user := c.Params("user")
		client := &Client{conn: c, room: room, user: user}
		clients[client] = true

		// 사용자가 WebSocket에 연결되었음을 온라인 상태로 업데이트
		rdb.SAdd(ctx, "online:"+room, user)

		go handleMessages(client)

		// 채팅방 입장 시 캐싱된 히스토리 또는 DB에서 히스토리 불러오기
		go sendChatHistory(client)

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("error: %v", err)
				delete(clients, client)
				// WebSocket 연결이 끊어졌음을 온라인 상태에서 제거
				rdb.SRem(ctx, "online:"+room, user)
				break
			}

			message := Message{
				Room:   client.room,
				User:   client.user,
				Text:   string(msg),
				SentAt: time.Now().Unix(),
				ReadBy: []string{},
			}

			messageBytes, _ := json.Marshal(message)
			rdb.RPush(ctx, "messages:"+client.room, messageBytes)
			nc.Publish("chat."+client.room, messageBytes)
		}
	}))

	log.Println("WebSocket server is running on port 3000")
	log.Fatal(app.Listen(":3000"))
}

func handleMessages(client *Client) {
	pubsub := rdb.Subscribe(ctx, client.room)
	defer pubsub.Close()

	readPubSub := rdb.Subscribe(ctx, "read-status:"+client.room)
	defer readPubSub.Close()

	go func() {
		for {
			msg, err := readPubSub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
			client.conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		}
	}()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		client.conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))

		markMessageAsRead(client.room, client.user, []byte(msg.Payload))
	}
}

func sendChatHistory(client *Client) {
	history, err := rdb.LRange(ctx, "messages:"+client.room, 0, -1).Result()
	if err != nil || len(history) == 0 {
		// Redis에 히스토리가 없으면 DB에서 조회
		filter := bson.M{"room": client.room}
		cursor, err := messagesCollection.Find(ctx, filter)
		if err != nil {
			log.Printf("Error fetching chat history from DB: %v", err)
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var msg Message
			if err := cursor.Decode(&msg); err != nil {
				log.Printf("Error decoding chat history from DB: %v", err)
				return
			}

			// 읽음 상태 조회
			readKey := "read:" + msg.Room + ":" + strconv.FormatInt(msg.SentAt, 10)
			readBy, _ := rdb.SMembers(ctx, readKey).Result()
			msg.ReadBy = readBy

			messageJSON, _ := json.Marshal(msg)
			history = append(history, string(messageJSON))
		}

		// Redis에 캐싱
		for _, msg := range history {
			rdb.RPush(ctx, "messages:"+client.room, msg)
		}
		rdb.Expire(ctx, "messages:"+client.room, 1*time.Hour)
	}

	for _, msg := range history {
		client.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func markMessageAsRead(room, user string, message []byte) {
	var msg Message
	json.Unmarshal(message, &msg)

	readKey := "read:" + room + ":" + strconv.FormatInt(msg.SentAt, 10)
	rdb.SAdd(ctx, readKey, user)
	rdb.Expire(ctx, readKey, 24*time.Hour)

	// 실시간으로 다른 사용자에게 읽음 상태 전송
	readStatus := ReadStatus{
		Room:  room,
		MsgID: msg.SentAt,
		User:  user,
	}
	readStatusBytes, _ := json.Marshal(readStatus)
	rdb.Publish(ctx, "read-status:"+room, readStatusBytes)
}
