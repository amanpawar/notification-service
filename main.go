package main

import (
	"log"
	"notification-service/internal/app"
	"notification-service/internal/config"
)

func main() {
	cfg := config.NewConfig()
	application := app.NewApp(cfg)

	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
