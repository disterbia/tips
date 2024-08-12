package main

import (
	"log"
	"medicine-service/core"

	_ "medicine-service/docs"
	"medicine-service/model"
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
	}
	dbPath := os.Getenv("DB_PATH")
	database, err := model.NewDB(dbPath)
	if err != nil {
		log.Println("Database connection error:", err)
	}

	// NATS 설정
	nc, err := nats.Connect("nats:4222")
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	svc := core.NewMedicineService(database, nc)

	saveEndpoint := core.SaveEndpoint(svc)
	removeEndpoint := core.RemoveEndpoint(svc)
	getTakensEndpoint := core.GetExpectsEndpoint(svc)
	getMedicinesEndpoint := core.GetMedicinesEndpoint(svc)
	takeEndpoint := core.TakeEndpoint(svc)
	unTakeEndpoint := core.UnTakeEndpoint(svc)
	searchEndpoint := core.SearchsEndpoint(svc)

	app := fiber.New()
	app.Use(logger.New())

	// Swagger 설정
	app.Get("/swagger/*", swagger.HandlerDefault) // Swagger UI 경로 설정
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	// CORS 미들웨어 추가
	app.Use(cors.New())

	app.Post("/save-medicine", core.SaveHandler(saveEndpoint))
	app.Post("/remove-medicine/:id", core.RemoveHandler(removeEndpoint))
	app.Post("/take-medicine", core.TakeHandler(takeEndpoint))
	app.Post("/untake-medicine/:id", core.UnTakeHandler(unTakeEndpoint))
	app.Get("/get-takens", core.GetExpectsHandler(getTakensEndpoint))
	app.Get("/get-medicines", core.GetMedicinesHandler(getMedicinesEndpoint))
	app.Get("/search-medicines", core.SearchHandler(searchEndpoint))

	log.Fatal(app.Listen(":44407"))

}
