// /emotion-service/main.go
package main

import (
	"emotion-service/core"
	"emotion-service/model"
	"log"
	"os"

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
		return
	}

	dbPath := os.Getenv("DB_PATH")
	database, err := model.NewDB(dbPath)
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}

	svc := core.NewEmotionService(database)

	saveEmotionEndpoint := core.SaveEmotionEndpoint(svc)
	getEmotionsEndpoint := core.GetEmotionsEndpoint(svc)

	app := fiber.New()
	app.Use(logger.New())

	// Swagger 설정
	app.Get("/swagger/*", swagger.HandlerDefault) // Swagger UI 경로 설정
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	// CORS 미들웨어 추가
	app.Use(cors.New())
	app.Post("/save-emotion", core.SaveEmotionHandler(saveEmotionEndpoint))
	app.Get("/get-emotions", core.GetEmotionsHandler(getEmotionsEndpoint))

	log.Fatal(app.Listen(":44408"))
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
