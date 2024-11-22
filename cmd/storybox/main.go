package main

import (
	"StoryBox/internal/app"
	"StoryBox/internal/utils"
	"log"

	"github.com/kataras/iris/v12"
)

func main() {
	// Load configuration
	config, err := utils.LoadConfiguration("config/config.yaml")
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

	runAddress := iris.Addr(config.Address + ":" + config.Port)

	// Run the application
	if err := api.Run(runAddress); err != nil {
		log.Fatalf("Application failed to run: %v", err)
	}
}
