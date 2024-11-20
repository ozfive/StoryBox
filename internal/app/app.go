package app

import (
	"log"
	"time"

	"github.com/didip/tollbooth/v6/limiter"
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

	// Set Expirable Options for Rate Limiter
	// newLimiter := limiter.New(&limiter.ExpirableOptions{
	// 	DefaultExpirationTTL: time.Minute,
	// 	MaxExpire:            1000,
	// })

	// Rate Limiter
	newLimiter := *limiter.New(&limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Minute})

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
	handlers.InitRFIDHandlers(app, rfidRepo, soundService, playlistService, &newLimiter) // Pass the correct limiter type
	handlers.InitPlaylistHandlers(app, playlistService, &newLimiter)
	handlers.InitStatsHandlers(app, soundService)

	return app
}
