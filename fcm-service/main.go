// /fcm-service/main.go
package main

import (
	"fcm-service/core"
	"fcm-service/model"
	"log"
	"os"

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
	core.StartCentralCronScheduler(database)

	select {}
}
