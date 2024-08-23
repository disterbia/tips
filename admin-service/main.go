package main

import (
	"admin-service/core"
	"log"
	"os"

	_ "admin-service/docs"

	"admin-service/model"

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

	svc := core.NewAdminService(database, conn)
	loginEndpoint := core.LoginEndpoint(svc)
	singInEndpoint := core.SignInEndpoint(svc)
	searchHospitalsEndpoint := core.SearchHospitalsEndpoint(svc)
	getPoliciesEndpoint := core.GetPoliciesEndpoint(svc)
	verifyEndpoint := core.VerifyEndpoint(svc)
	sendCodeForSigninEndpoint := core.SendCodeForSignInEndpoint(svc)
	sendCodeForIdEndpoint := core.SendCodeForIdEndpoint(svc)
	sendCodeForPwEndpoint := core.SendCodeForPwEndpoint(svc)
	changwPwEndpoint := core.ChangePwEndpoint(svc)
	findIdEndpoint := core.FindIdEndpoint(svc)
	questionEndpoint := core.QuestionEndpoint(svc)

	app := fiber.New()
	app.Use(logger.New())
	// Swagger 설정
	app.Get("/swagger/*", swagger.HandlerDefault) // Swagger UI 경로 설정
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	// CORS 미들웨어 추가
	app.Use(cors.New())

	app.Get("/search-hospitals", core.SearchHospitalsHandler(searchHospitalsEndpoint))
	app.Get("/get-policies", core.GetPoliciesHandler(getPoliciesEndpoint))

	app.Post("/login", core.LoginHandler(loginEndpoint))
	app.Post("/sign-in", core.SignInHandler(singInEndpoint))
	app.Post("/send-code-sign/:number", core.SendCodeForSignInHandler(sendCodeForSigninEndpoint))
	app.Post("/send-code-id", core.SendCodeForIdHandler(sendCodeForIdEndpoint))
	app.Post("/send-code-pw", core.SendCodeForPwHandler(sendCodeForPwEndpoint))
	app.Post("/verify-code", core.VerifyHandler(verifyEndpoint))
	app.Post("/find-id", core.FindIdHandler(findIdEndpoint))
	app.Post("/change-pw", core.ChangePwHandler(changwPwEndpoint))
	app.Post("/question", core.QuestionHandler(questionEndpoint))

	log.Fatal(app.Listen(":44400"))

}
