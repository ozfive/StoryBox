package services

import (
	"StoryBox/internal/models"
	"StoryBox/internal/repository"

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
	err := p.repo.Create(url, playlistName)
	if err != nil {
		ctx.StatusCode(500)
		ctx.JSON(iris.Map{
			"status_code": 500,
			"message":     "Failed to create playlist.",
			"data":        err.Error(),
		})
		return err
	}

	ctx.StatusCode(200)
	ctx.JSON(iris.Map{
		"status_code": 200,
		"message":     "The playlist has been created successfully.",
	})

	return nil
}

func (p *playlistService) ClearPlaylist() error {
	// TODO: Implement the ClearPlaylist method here.
	return nil
}

func (p *playlistService) DeletePlaylist(url, playlistName string) error {
	// TODO: Implement the DeletePlaylist method here.
	return nil
}

func (p *playlistService) GetPlaylist(url, playlistName string) (*models.Playlist, error) {
	// TODO: Implement the GetPlaylist method here.
	return nil, nil
}

func (p *playlistService) LoadPlaylist(playlistName string) error {
	// TODO: Implement the LoadPlaylist method here.
	return nil
}

func (p *playlistService) PausePlaylist(playlistName string) error {
	// TODO: Implement the PausePlaylist functionality here.
	return nil
}

func (p *playlistService) PlayPlaylist(playlistName string) error {
	// TODO: Implement the PlayPlaylist functionality here.
	return nil
}

func (p *playlistService) StopPlaylist() error {
	// TODO: Implement the StopPlaylist functionality here.
	return nil
}
