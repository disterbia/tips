// /emotion-service/main.go
package main

import (
	"check-service/core"
	_ "check-service/docs"
	"check-service/model"
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

	svc := core.NewCheckService(database)

	getSampleVideos := core.GetSampleVideosEndpoint(svc)
	getFaceScoresEndpoint := core.GetFaceInfoEndpoint(svc)
	getScoresEndpoint := core.GetTapBlinkScoreEndpoint(svc)
	saveFaceScoreEndpoint := core.SaveFaceInfoEndpoint(svc)
	saveScoreEndpoint := core.SaveTapBlinkScoreEndpoint(svc)

	app := fiber.New()
	app.Use(logger.New())

	// Swagger 설정
	app.Get("/swagger/*", swagger.HandlerDefault) // Swagger UI 경로 설정
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	// CORS 미들웨어 추가
	app.Use(cors.New())
	app.Get("/get-videos", core.GetSampleVideosHandler(getSampleVideos))
	app.Get("/get-face-infos", core.GetFaceInfosHandler(getFaceScoresEndpoint))
	app.Get("/get-scores", core.GetScoresHandler(getScoresEndpoint))
	app.Post("/save-face-info", core.SaveFaceInfoHandler(saveFaceScoreEndpoint))
	app.Post("/save-score", core.SaveScoreHandler(saveScoreEndpoint))

	log.Fatal(app.Listen(":44411"))
	// router.RunTLS(":8080", "cert.pem", "key.pem")

}
