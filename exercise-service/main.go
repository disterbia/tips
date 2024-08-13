// /exercise-service/main.go
package main

import (
	"exercise-service/core"
	_ "exercise-service/docs"
	"exercise-service/model"
	"log"
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
		return
	}

	dbPath := os.Getenv("DB_PATH")
	database, err := model.NewDB(dbPath)
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}

	// NATS 설정
	nc, err := nats.Connect("nats:4222")
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	defer nc.Close()
	svc := core.NewExerciseService(database, nc)

	saveExerciseEndpoint := core.SaveExerciseEndpoint(svc)
	getExpectsEndpoint := core.GetExpectsEndpoint(svc)
	removeExercisesEndpoint := core.RemoveExerciseEndpoint(svc)
	doExerciseEndpoint := core.DoExerciseEndpoint(svc)
	getExercisesEndpoint := core.GetExercisesEndpoint(svc)
	getProjectsEndpoint := core.GetProjectsEndpoint(svc)
	getVideosEndpoint := core.GetVideosEndpoint(svc)

	app := fiber.New()
	app.Use(logger.New())

	// Swagger 설정
	app.Get("/swagger/*", swagger.HandlerDefault) // Swagger UI 경로 설정
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	// CORS 미들웨어 추가
	app.Use(cors.New())

	app.Post("/save-exercise", core.SaveExerciseHandler(saveExerciseEndpoint))
	app.Post("/remove-exercises", core.RemoveExercisesHandler(removeExercisesEndpoint))
	app.Post("/do-exercise", core.DoExerciseHandler(doExerciseEndpoint))
	app.Get("/get-takens", core.GetExpectsHandler(getExpectsEndpoint))
	app.Get("/get-exercises", core.GetExercisesHandler(getExercisesEndpoint))
	app.Get("/get-projects", core.GetProjectsHandler(getProjectsEndpoint))
	app.Get("/get-videos", core.GetVideosHandler(getVideosEndpoint))

	log.Fatal(app.Listen(":44405"))
}
