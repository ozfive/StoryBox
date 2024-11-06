package app

import (
	"log"

	"github.com/didip/tollbooth/v6"
	"github.com/kataras/iris/v12"
	_ "github.com/mattn/go-sqlite3"

	"StoryBox/internal/handlers"
	"StoryBox/internal/repository"
	"StoryBox/internal/services"
)

type Config struct {
	Debug        bool
	Address      string
	Port         string
	LogFilePath  string
	DatabasePath string
}

func NewApp(config *Config) *iris.Application {
	app := iris.New()

	// Rate Limiter
	limiter := tollbooth.NewLimiter(15, nil)

	// Connect to Database
	db, err := repository.ConnectDatabase(config.DatabasePath)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Initialize Repositories
	rfidRepo := repository.NewRFIDRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)

	// Initialize Services
	soundService := services.NewSoundService()
	playlistService := services.NewPlaylistService(playlistRepo, soundService)

	// Initialize Handlers
	handlers.InitRFIDHandlers(app, rfidRepo, soundService, playlistService, limiter)
	handlers.InitPlaylistHandlers(app, playlistService, limiter)
	handlers.InitStatsHandlers(app, soundService)

	return app
}
