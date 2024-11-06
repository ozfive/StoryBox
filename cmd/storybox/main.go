package main

import (
	"StoryBox/internal/app"
	"StoryBox/internal/utils"
	"log"
)

func main() {
	// Load configuration
	config, err := utils.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	err = utils.InitLogger(config.LogFilePath)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize and run the application
	api := app.NewApp(config)
	if err := api.Run(); err != nil {
		log.Fatalf("Application failed to run: %v", err)
	}
}
