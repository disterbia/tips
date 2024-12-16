package main

import (
	"log"

	"landing-service/core"
	_ "landing-service/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// gRPC 클라이언트 연결 생성
	conn, err := grpc.NewClient("email:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to email service: %v", err)
	}
	defer conn.Close()

	svc := core.NewLandingService(conn)
	kldgaInquireEndpoint := core.KldgaInquireEndpoint(svc)

	app := fiber.New()
	app.Use(logger.New())

	// Swagger 설정
	app.Get("/swagger/*", swagger.HandlerDefault) // Swagger UI 경로 설정
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	// CORS 미들웨어 추가
	app.Use(cors.New())

	app.Post("/kldga/inquire", core.KldgaInquireHandler(kldgaInquireEndpoint))

	log.Fatal(app.Listen(":44500"))

}
