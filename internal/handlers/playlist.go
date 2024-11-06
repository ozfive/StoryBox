package handlers

import (
	"net/http"

	"github.com/didip/tollbooth/v6"
	"github.com/iris-contrib/middleware/tollboothic"
	"github.com/kataras/iris/v12"

	"StoryBox/internal/services"
)

func InitPlaylistHandlers(app *iris.Application, playlistService services.PlaylistService, limiter *tollbooth.Limiter) {
	p := app.Party("/playlist")
	{
		p.Post("/create", tollboothic.LimitHandler(limiter), createPlaylistHandler(playlistService))
		p.Post("/delete", tollboothic.LimitHandler(limiter), deletePlaylistHandler(playlistService))
	}
}

func createPlaylistHandler(playlistService services.PlaylistService) iris.Handler {
	return func(ctx iris.Context) {
		var payload struct {
			URL          string `json:"url"`
			PlaylistName string `json:"playlistname"`
		}

		if err := ctx.ReadJSON(&payload); err != nil {
			ctx.StatusCode(http.StatusBadRequest)
			ctx.JSON(iris.Map{
				"status_code": 400,
				"message":     "Malformed JSON input.",
			})
			return
		}

		if payload.URL == "" || payload.PlaylistName == "" {
			ctx.StatusCode(http.StatusUnprocessableEntity)
			ctx.JSON(iris.Map{
				"status_code": 422,
				"message":     "URL and PlaylistName must be provided.",
			})
			return
		}

		if err := playlistService.CreatePlaylist(payload.URL, payload.PlaylistName); err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.JSON(iris.Map{
				"status_code": 500,
				"message":     "Failed to create playlist.",
				"data":        err.Error(),
			})
			return
		}

		ctx.StatusCode(http.StatusOK)
		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "The playlist has been created successfully.",
		})
	}
}

func deletePlaylistHandler(playlistService services.PlaylistService) iris.Handler {
	return func(ctx iris.Context) {
		var payload struct {
			URL          string `json:"url"`
			PlaylistName string `json:"playlistname"`
		}

		if err := ctx.ReadJSON(&payload); err != nil {
			ctx.StatusCode(http.StatusBadRequest)
			ctx.JSON(iris.Map{
				"status_code": 400,
				"message":     "Malformed JSON input.",
			})
			return
		}

		if payload.URL == "" || payload.PlaylistName == "" {
			ctx.StatusCode(http.StatusUnprocessableEntity)
			ctx.JSON(iris.Map{
				"status_code": 422,
				"message":     "URL and PlaylistName must be provided.",
			})
			return
		}

		if err := playlistService.DeletePlaylist(payload.URL, payload.PlaylistName); err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.JSON(iris.Map{
				"status_code": 500,
				"message":     "Failed to delete playlist.",
				"data":        err.Error(),
			})
			return
		}

		ctx.StatusCode(http.StatusOK)
		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "The playlist has been deleted successfully.",
		})
	}
}
