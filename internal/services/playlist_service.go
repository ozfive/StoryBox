package services

import (
	"StoryBox/internal/models"
	"StoryBox/internal/repository"
	"log"

	"github.com/kataras/iris/v12"
	// "StoryBox/internal/models"
	// "StoryBox/internal/repository"
)

type PlaylistService interface {
	CreatePlaylist(url, playlistName string) error
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

func (p *playlistService) CreatePlaylist(url, playlistName string) error {
	database, err := connectToDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer database.Close()

	// Check if the playlist already exists in the database.
	var count int
	sqlCheck := "SELECT COUNT(*) FROM playlist WHERE url = ? AND playlistname = ?"
	err = database.QueryRow(sqlCheck, url, playlistname).Scan(&count)

	if err != nil {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "Failed to SELECT playlist " + playlistname + " from the database. Please try again.",
		})
		return
	}

	if count > 0 {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "The playlist " + playlistname + " already exists in the database.",
		})
		return
	}

	// Insert the new playlist into the database.
	sqlInsert := "INSERT INTO playlist (url, playlistname) VALUES (?, ?)"
	_, err = database.Exec(sqlInsert, url, playlistname)

	if err != nil {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "Failed to INSERT playlist in the database. Please try again.",
		})
		return
	}

	// Return a success message.
	ctx.StatusCode(200)
	ctx.JSON(iris.Map{
		"status_code": 200,
		"message":     "The playlist " + playlistname + " has been created in the database.",
	})
	return nil
}

// Implement other PlaylistService methods similarly...
