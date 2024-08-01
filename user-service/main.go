package main

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"
	"user-service/core"

	_ "user-service/docs"

	"user-service/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

var ipLimiters = make(map[string]*rate.Limiter)
var ipLimitersMutex sync.Mutex

func getClientIP(c *fiber.Ctx) string {
	if ip := c.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	return c.IP()
}

func RateLimitMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := getClientIP(c)

		ipLimitersMutex.Lock()
		limiter, exists := ipLimiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(time.Hour/10), 10)
			ipLimiters[ip] = limiter
		}
		ipLimitersMutex.Unlock()

		if !limiter.Allow() {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "요청 횟수 초과"})
		}

		return c.Next()
	}
}
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

	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET")
	bucketUrl := os.Getenv("S3_BUCKET_URL")
	s3sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-2"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		log.Println("aws connection error:", err)
	}
	s3svc := s3.New(s3sess)

	svc := core.NewUserService(database, s3svc, bucket, bucketUrl)
	loginEndpoint := core.SnsLoginEndpoint(svc)
	phoneLoginEndpoint := core.PhoneLoginEndpoint(svc)
	sendCodeEndpoint := core.SendCodeEndpoint(svc)
	autoLoginEndpoint := core.AutoLoginEndpoint(svc)
	verifyEndpoint := core.VerifyEndpoint(svc)
	updateEndpoint := core.UpdateUserEndpoint(svc)
	linkEndpoint := core.LinkEndpoint(svc)
	getUserEndpoint := core.GetUserEndpoint(svc)
	removeUserEndpoint := core.RemoveEndpoint(svc)
	getVersionEndpoint := core.GetVersionEndpoint(svc)

	app := fiber.New()
	app.Use(logger.New())
	// Swagger 설정
	app.Get("/swagger/*", swagger.HandlerDefault) // Swagger UI 경로 설정
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	// CORS 미들웨어 추가
	app.Use(cors.New())

	app.Get("/get-user", core.GetUserHandler(getUserEndpoint))
	app.Get("/get-version", core.GetVersionHandeler(getVersionEndpoint))

	app.Post("/sns-login", core.SnsLoginHandler(loginEndpoint))
	app.Post("/phone-login", core.PhoneLoginHandler(phoneLoginEndpoint))
	app.Post("/auto-login", core.AutoLoginHandler(autoLoginEndpoint))
	app.Post("/send-code/:email", RateLimitMiddleware(), core.SendCodeHandler(sendCodeEndpoint))
	app.Post("/verify-code", core.VerifyHandler(verifyEndpoint))
	app.Post("/update-user", core.UpdateUserHandler(updateEndpoint))
	app.Post("/link-email", core.LinkHandler(linkEndpoint))
	app.Post("/remove-user", core.RemoveHandler(removeUserEndpoint))

	log.Fatal(app.Listen(":44409"))

}
