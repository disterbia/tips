package main

import (
	"log"
	"os"

	"inquire-service/core"
	_ "inquire-service/docs"
	"inquire-service/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// gRPC 클라이언트 연결 생성
	conn, err := grpc.NewClient("email:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to email service: %v", err)
	}
	defer conn.Close()

	inquireSvc := core.NewInquireService(database, conn)
	adminLoginEndpoint := core.AdminLoginEndpoint(inquireSvc)
	answerEndpoint := core.AnswerEndpoint(inquireSvc)
	sendEndpoint := core.SendEndpoint(inquireSvc)
	getEndpoint := core.GetEndpoint(inquireSvc)
	allEndpoint := core.GetAllEndpoint(inquireSvc)
	removeInquireEndpoint := core.RemoveInquireEndpoint(inquireSvc)
	reomoveReplyEndpoint := core.RemoveReplyEndpoint(inquireSvc)

	app := fiber.New()
	app.Use(logger.New())

	// Swagger 설정
	app.Get("/swagger/*", swagger.HandlerDefault) // Swagger UI 경로 설정
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	// CORS 미들웨어 추가
	app.Use(cors.New())

	app.Post("/login", core.AdminLoginHandler(adminLoginEndpoint))
	app.Post("/inquire-reply", core.AnswerHandler(answerEndpoint))
	app.Post("/send-inquire", core.SendHandler(sendEndpoint))
	app.Post("/remove-inquire/:id", core.RemoveInquireHandler(removeInquireEndpoint))
	app.Post("/remove-reply/:id", core.RemoveReplyHandler(reomoveReplyEndpoint))
	app.Get("/get-inquires", core.GetHandler(getEndpoint))
	app.Get("/all-inquires", core.GetAllHandler(allEndpoint))

	log.Fatal(app.Listen(":44406"))

}
