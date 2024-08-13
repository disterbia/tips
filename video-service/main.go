// /video-service/main.go
package main

import (
	"log"
	"os"
	"video-service/core"
	_ "video-service/docs"
	"video-service/model"

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

	svc := core.NewAdminVideoService(database)

	getVimeoLevel1sEndpoint := core.GetVimeoLevel1sEndpoint(svc)
	getVimeoLevel2sEndpoint := core.GetVimeoLevel2sEndpoint(svc)
	saveEndpoint := core.SaveEndpoint(svc)

	app := fiber.New()
	app.Use(logger.New())

	// Swagger 설정
	app.Get("/swagger/*", swagger.HandlerDefault) // Swagger UI 경로 설정
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	// CORS 미들웨어 추가
	app.Use(cors.New())

	app.Get("/get-items", core.GetVimeoLevel1sHandler(getVimeoLevel1sEndpoint))
	app.Get("/get-videos/:id", core.GetVimeoLevel2sHandler(getVimeoLevel2sEndpoint))
	app.Post("/save-videos", core.SaveHandler(saveEndpoint))

	log.Fatal(app.Listen(":44410"))

}
