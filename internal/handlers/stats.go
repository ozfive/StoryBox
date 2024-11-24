package handlers

import (
	"os/exec"
	"strings"

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
		cmd := exec.Command("./mpdcurrentsong")
		output, err := cmd.Output()
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"status_code": 500,
				"message":     "Failed to retrieve current stats.",
				"data":        err.Error(),
			})
			return
		}

		data := parseKeyValueOutput(string(output))
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
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"status_code": 500,
				"message":     "Failed to stop the current playlist.",
				"data":        err.Error(),
			})
			return
		}

		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "Playlist stopped.",
			"data":        nil,
		})
	}
}

func playlistElapsedTimeHandler(soundService services.SoundService) iris.Handler {
	return func(ctx iris.Context) {
		cmd := exec.Command("./mpdtime")
		output, err := cmd.Output()
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{
				"status_code": 500,
				"message":     "Failed to retrieve elapsed time.",
				"data":        err.Error(),
			})
			return
		}

		elapsedTime := strings.TrimSpace(string(output))
		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "Elapsed Time",
			"data":        elapsedTime,
		})
	}
}

func parseKeyValueOutput(output string) map[string]string {
	data := make(map[string]string)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) == 2 {
			data[parts[0]] = parts[1]
		}
	}
	return data
}
