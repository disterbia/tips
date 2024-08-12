package main

import (
	"encoding/json"
	"log"
	"notification-service/core"
	_ "notification-service/docs"
	"notification-service/model"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	dbPath := os.Getenv("DB_PATH")
	database, err := model.NewDB(dbPath)
	if err != nil {
		log.Println("Database connection error:", err)
	}

	// NATS 연결 설정
	nc, err := nats.Connect("nats:4222")
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	notiSvc := core.NewNotificationService(database)

	// NATS 구독 설정을 고루틴으로 실행
	go func() {
		_, err := nc.Subscribe("save-notification", func(m *nats.Msg) {
			var req core.NotificationRequest
			if err := json.Unmarshal(m.Data, &req); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				return
			}

			code, err := notiSvc.SaveNotifications(req)
			if err != nil {
				log.Printf("Error saving notification: %v", err)
				return
			}
			log.Printf("Notification saved with code: %s", code)
		})
		if err != nil {
			log.Fatalf("Error subscribing to save-notification: %v", err)
		}

		// remove-notification 이벤트 구독
		_, err = nc.Subscribe("remove-notification", func(m *nats.Msg) {
			var req core.NotificationRequest
			if err := json.Unmarshal(m.Data, &req); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				return
			}

			code, err := notiSvc.RemoveNotifications(req)
			if err != nil {
				log.Printf("Error removing notification: %v", err)
				return
			}
			log.Printf("Notification removed with code: %s", code)
		})
		if err != nil {
			log.Fatalf("Error subscribing to remove-notification: %v", err)
		}

		// NATS 구독 서버는 고루틴 내에서 계속 실행 상태를 유지
		select {}
	}()

	getMessagesEndpoint := core.GetMessagesEndpoint(notiSvc)
	readMessagesNotisEndpoint := core.ReadAllEndpoint(notiSvc)
	removeMessagesEndpoint := core.RemoveMessagesEndpoint(notiSvc)
	app := fiber.New()
	app.Use(logger.New())

	// Swagger 설정
	app.Get("/swagger/*", swagger.HandlerDefault) // Swagger UI 경로 설정
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	// CORS 미들웨어 추가
	app.Use(cors.New())

	app.Get("/get-messages", core.GetMessagesHandler(getMessagesEndpoint))
	app.Post("/remove-messages", core.RemoveMessagesHandler(removeMessagesEndpoint))
	app.Post("/read-message", core.ReadAllHandler(readMessagesNotisEndpoint))

	log.Fatal(app.Listen(":44406"))
}
