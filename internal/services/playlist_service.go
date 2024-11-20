package services

import (
	"StoryBox/internal/models"
	"StoryBox/internal/repository"
	"StoryBox/internal/utils"
	"log"

	"github.com/kataras/iris/v12"
)

type PlaylistService interface {
	CreatePlaylist(ctx iris.Context, url, playlistName string) error
	DeletePlaylist(url, playlistName string) error
	GetPlaylist(url, playlistName string) (*models.Playlist, error)
	ClearPlaylist() error
	LoadPlaylist(playlistName string) error
	PlayPlaylist(playlistName string) error
	PausePlaylist(playlistName string) error
	StopPlaylist() error
}

type playlistService struct {
	repo         repository.PlaylistRepository
	soundService SoundService
}

func NewPlaylistService(repo repository.PlaylistRepository, soundService SoundService) PlaylistService {
	return &playlistService{
		repo:         repo,
		soundService: soundService,
	}
}

func (p *playlistService) CreatePlaylist(ctx iris.Context, url, playlistName string) error {
	database, err := repository.ConnectDatabase(utils.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer database.Close()

	// Check if the playlist already exists in the database.
	var count int
	sqlCheck := "SELECT COUNT(*) FROM playlist WHERE url = ? AND playlistname = ?"
	err = database.QueryRow(sqlCheck, url, playlistName).Scan(&count)

	if err != nil {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "Failed to SELECT playlist " + playlistName + " from the database. Please try again.",
		})
		return err
	}

	if count > 0 {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "The playlist " + playlistName + " already exists in the database.",
		})
		return nil
	}

	// Add code to create the playlist in the database here.

	return nil
}

func (p *playlistService) ClearPlaylist() error {
	// Implement the ClearPlaylist method here.
	return nil
}

// Implement other methods of PlaylistService interface here.
