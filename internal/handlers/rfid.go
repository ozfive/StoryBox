package handlers

import (
	"net/http"

	"github.com/didip/tollbooth/v6"
	"github.com/iris-contrib/middleware/tollboothic"
	"github.com/kataras/iris/v12"

	"StoryBox/internal/models"
	"StoryBox/internal/repository"
	"StoryBox/internal/services"
)

func InitRFIDHandlers(app *iris.Application, rfidRepo repository.RFIDRepository, soundService services.SoundService, playlistService services.PlaylistService, limiter *tollbooth.Limiter) {
	r := app.Party("/rfid")
	{
		r.Post("/", tollboothic.LimitHandler(limiter), rfidHandler(rfidRepo, soundService, playlistService))
		r.Post("/create", tollboothic.LimitHandler(limiter), rfidCreateHandler(rfidRepo))
	}
}

func rfidHandler(rfidRepo repository.RFIDRepository, soundService services.SoundService, playlistService services.PlaylistService) iris.Handler {
	return func(ctx iris.Context) {
		var rfid models.RFID
		if err := ctx.ReadJSON(&rfid); err != nil {
			ctx.StatusCode(http.StatusBadRequest)
			ctx.JSON(iris.Map{
				"status_code": 400,
				"message":     "Malformed JSON input.",
			})
			return
		}

		if rfid.TagID == "" || rfid.UniqueID == "" || rfid.URL == "" || rfid.PlaylistName == "" {
			ctx.StatusCode(http.StatusUnprocessableEntity)
			ctx.JSON(iris.Map{
				"status_code": 422,
				"message":     "All fields must be provided.",
			})
			return
		}

		existingRFID, err := rfidRepo.GetByTagAndUniqueID(rfid.TagID, rfid.UniqueID)
		if err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.JSON(iris.Map{
				"status_code": 500,
				"message":     "Database error.",
				"data":        err.Error(),
			})
			return
		}

		if existingRFID != nil {
			// RFID exists, handle accordingly
			ctx.JSON(iris.Map{
				"status_code": 200,
				"message":     "INFO: The RFID tag exists. Playlist queued!",
				"data":        existingRFID,
			})

			// Play notification and manage playlist
			soundService.PlayAcknowledgeNotification()
			playlistService.ClearPlaylist()
			playlistService.LoadPlaylist(existingRFID.PlaylistName)
			playlistService.PlayPlaylist(existingRFID.PlaylistName)
		} else {
			ctx.StatusCode(http.StatusBadRequest)
			ctx.JSON(iris.Map{
				"status_code": 400,
				"message":     "The RFID tag was not found in the DB.",
				"data":        "",
			})
		}
	}
}

func rfidCreateHandler(rfidRepo repository.RFIDRepository) iris.Handler {
	return func(ctx iris.Context) {
		var rfid models.RFID
		if err := ctx.ReadJSON(&rfid); err != nil {
			ctx.StatusCode(http.StatusBadRequest)
			ctx.JSON(iris.Map{
				"status_code": 400,
				"message":     "Malformed JSON input.",
			})
			return
		}

		if rfid.TagID == "" || rfid.UniqueID == "" || rfid.URL == "" || rfid.PlaylistName == "" {
			ctx.StatusCode(http.StatusUnprocessableEntity)
			ctx.JSON(iris.Map{
				"status_code": 422,
				"message":     "All fields must be provided.",
			})
			return
		}

		err := rfidRepo.Create(&rfid)
		if err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.JSON(iris.Map{
				"status_code": 500,
				"message":     "Failed to create RFID tag.",
				"data":        err.Error(),
			})
			return
		}

		ctx.StatusCode(http.StatusOK)
		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "The RFID tag was successfully created in the database.",
			"data":        rfid,
		})
	}
}
