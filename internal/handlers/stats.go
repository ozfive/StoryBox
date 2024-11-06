package handlers

import (
	"github.com/kataras/iris/v12"

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
		// Implement the logic to retrieve current stats
		// Use services.SoundService as needed
		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "Current Stats",
			"data":        "Some stats data",
		})
	}
}

func stopCurrentPlaylistHandler(soundService services.SoundService) iris.Handler {
	return func(ctx iris.Context) {
		// Implement the logic to stop the current playlist
		// Use services.SoundService as needed
		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "Playlist stopped.",
			"data":        nil,
		})
	}
}

func playlistElapsedTimeHandler(soundService services.SoundService) iris.Handler {
	return func(ctx iris.Context) {
		// Implement the logic to get elapsed time
		// Use services.SoundService as needed
		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "Elapsed Time",
			"data":        "00:03:25",
		})
	}
}
