package main

import (
	"log"

	"landing-service/core"
	_ "landing-service/docs"

	"github.com/go-redis/redis/v8"
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
	// gRPC 클라이언트 연결 생성
	conn, err := grpc.NewClient("email:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to email service: %v", err)
	}
	defer conn.Close()

	// Redis 클라이언트 설정
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // 비밀번호가 없으면 비워둠
		DB:       0,  // Redis DB 번호
	})

	// 서비스 생성
	svc := core.NewLandingService(conn, redisClient)
	kldgaInquireEndpoint := core.KldgaInquireEndpoint(svc)
	kldgaCompetitionEndpoint := core.KldgaCompetitionEndpoint(svc)
	sendAuthCodeEndpoint := core.SendAuthCodeEndpoint(svc)
	verifyAuthCodeEndpoint := core.VerifyAuthCodeEndpoint(svc)
	adapfitInquireEndpoint := core.AdapfitInqruieEndpoint(svc)

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
	app.Post("/kldga/competition", core.KldgaCompetitionHandler(kldgaCompetitionEndpoint))
	app.Post("/kldga/send-code", core.SendAuthCodeHandler(sendAuthCodeEndpoint))
	app.Post("/kldga/verify-code", core.VerifyAuthCodeHandler(verifyAuthCodeEndpoint))
	app.Post("/adapfit/inquire", core.AdapfitInquireHandler(adapfitInquireEndpoint))

	log.Fatal(app.Listen(":44500"))

}
