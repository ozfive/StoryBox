package handlers

import (
	"os/exec"

	"github.com/kataras/iris/v12"

	"StoryBox/internal/mpdstats"
	"StoryBox/internal/services"
)

func InitStatsHandlers(app *iris.Application, soundService services.SoundService) {
	s := app.Party("/currentstats")
	{
		s.Get("/", currentStatsHandler(soundService))
	}

	s = app.Party("/stopcurrentplaylist")
	{
		s.Get("/", stopCurrentPlaylistHandler(soundService))
	}

	s = app.Party("/playlistelapsedtime")
	{
		s.Get("/", playlistElapsedTimeHandler(soundService))
	}
}

func currentStatsHandler(soundService services.SoundService) iris.Handler {
	return func(ctx iris.Context) {
		song, err := mpdstats.GetCurrentSongInfo()
		if err != nil {
			soundService.PlayErrorNotification()
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"status_code": 500,
				"message":     "Failed to retrieve current stats.",
				"data":        err.Error(),
			})
			return
		}

		data := map[string]string{
			"current_song": song,
		}
		soundService.PlayReadyNotification()
		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "Current Stats",
			"data":        data,
		})
	}
}

func stopCurrentPlaylistHandler(soundService services.SoundService) iris.Handler {
	return func(ctx iris.Context) {
		cmd := exec.Command("mpc", "stop")
		err := cmd.Run()
		if err != nil {
			soundService.PlayErrorNotification()
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"status_code": 500,
				"message":     "Failed to stop the current playlist.",
				"data":        err.Error(),
			})
			return
		}

		soundService.PlayAcknowledgeNotification()
		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "Playlist stopped.",
			"data":        nil,
		})
	}
}

func playlistElapsedTimeHandler(soundService services.SoundService) iris.Handler {
	return func(ctx iris.Context) {
		elapsed, err := mpdstats.GetElapsedTime()
		if err != nil {
			soundService.PlayErrorNotification()
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"status_code": 500,
				"message":     "Failed to retrieve elapsed time.",
				"data":        err.Error(),
			})
			return
		}

		data := map[string]interface{}{
			"elapsed_time": elapsed,
		}
		soundService.PlayReadyNotification()
		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "Elapsed Time",
			"data":        data,
		})
	}
}
