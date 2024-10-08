package main

import (
	"log"
	"os"
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
	sendCodeSignInEndpoint := core.SendCodeForSignInEndpoint(svc)
	sendCodeLoginEndpoint := core.SendCodeForLoginEndpoint(svc)
	autoLoginEndpoint := core.AutoLoginEndpoint(svc)
	verifyEndpoint := core.VerifyEndpoint(svc)
	updateEndpoint := core.UpdateUserEndpoint(svc)
	linkEndpoint := core.LinkEndpoint(svc)
	getUserEndpoint := core.GetUserEndpoint(svc)
	removeUserEndpoint := core.RemoveEndpoint(svc)
	getVersionEndpoint := core.GetVersionEndpoint(svc)
	getPolicesEndpoint := core.GetPolicesEndpoint(svc)

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
	app.Get("/get-police", core.GetPolicesHandeler(getPolicesEndpoint))

	app.Post("/sns-login", core.SnsLoginHandler(loginEndpoint))
	app.Post("/phone-login", core.PhoneLoginHandler(phoneLoginEndpoint))
	app.Post("/auto-login", core.AutoLoginHandler(autoLoginEndpoint))
	app.Post("/send-code-join/:number", core.SendCodeForSignInHandler(sendCodeSignInEndpoint))
	app.Post("/send-code-login/:number", core.SendCodeForLoginHandler(sendCodeLoginEndpoint))
	app.Post("/verify-code", core.VerifyHandler(verifyEndpoint))
	app.Post("/update-user", core.UpdateUserHandler(updateEndpoint))
	app.Post("/link-email", core.LinkHandler(linkEndpoint))
	app.Post("/remove-user", core.RemoveHandler(removeUserEndpoint))

	log.Fatal(app.Listen(":44409"))

}
