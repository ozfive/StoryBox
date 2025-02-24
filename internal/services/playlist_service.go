package services

import (
	"StoryBox/internal/models"
	"StoryBox/internal/repository"
	"fmt"

	"github.com/kataras/iris/v12"
)

type PlaylistService interface {
	CreatePlaylist(ctx iris.Context, url, playlistName string) error
	DeletePlaylist(url, playlistName string) error
	GetPlaylist(url, playlistName string) (*models.Playlist, error)
	ClearPlaylist(playlistName string) error
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

func (p *playlistService) GetPlaylist(url, playlistName string) (*models.Playlist, error) {
	if url == "" || playlistName == "" {
		return nil, fmt.Errorf("url and playlistName must be provided")
	}
	playlist, err := p.repo.Get(url, playlistName)
	if err != nil {
		return nil, err
	}
	return playlist, nil
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

func (p *playlistService) ClearPlaylist(playlistName string) error {
	return p.repo.Clear(playlistName)
}

func (p *playlistService) DeletePlaylist(url, playlistName string) error {
	return p.repo.Delete(url, playlistName)
}

func (p *playlistService) LoadPlaylist(playlistName string) error {
	return p.repo.Load(playlistName)
}

func (p *playlistService) PausePlaylist(playlistName string) error {
	return p.repo.Pause(playlistName)
}

func (p *playlistService) PlayPlaylist(playlistName string) error {
	return p.repo.Play(playlistName)
}

func (p *playlistService) StopPlaylist() error {
	return p.repo.Stop()
}
