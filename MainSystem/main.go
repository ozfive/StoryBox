package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kataras/iris/v12"

	// "github.com/kataras/iris/v12/middleware/basicauth"

	"github.com/didip/tollbooth/v6"
	"github.com/iris-contrib/middleware/tollboothic"

	// "github.com/kataras/iris/middleware/basicauth"
	_ "github.com/mattn/go-sqlite3"
	// phatbeat "github.com/ozfive/phatbeat/lib"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	// "time"
)

const (
	readingJSONParamsError = "Error while reading JSON parameters: <b>"
)

type RFID struct {
	ID           int    `json:"id"`
	TagID        string `json:"tagid"`
	UniqueID     string `json:"uniqueid"`
	URL          string `json:"url"`
	PlaylistName string `json:"playlistname"`
}

func main() {
	// playReadySound()

	debug := false

	localIPAddress := "localhost"
	localIPPort := "3001"

	// basicAuthUser := "admin"
	// basicAuthPassword := "password"

	// basicAuthUser, basicAuthPassword

	api := newApp(debug)

	err := api.Run(iris.Addr(localIPAddress+":"+localIPPort), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
	if err != nil {
		log.Println(err.Error())
	}

	/*

	   phatbeat.On(phatbeat.FastForward, true, func(p int) {
	           log.Println("Fast Forward")
	   })

	   phatbeat.Hold(phatbeat.FastForward, false, 2, func(b int) {
	           log.Println("FF Held")
	   })

	   phatbeat.On(phatbeat.PlayPause, false, func(b int) {
	           log.Println("PP")
	   })

	   phatbeat.Hold(phatbeat.PlayPause, false, 2, func(b int) {
	           log.Println("PP held")
	   })

	   phatbeat.On(phatbeat.VolDown, true, func(p int) {
	           log.Println("VolDown")
	   })

	   phatbeat.Hold(phatbeat.VolDown, false, 2, func(b int) {
	           log.Println("VolDown")
	   })

	   phatbeat.On(phatbeat.VolUp, true, func(p int) {
	           log.Println("VolUp")
	   })

	   phatbeat.Hold(phatbeat.VolUp, false, 2, func(b int) {
	           log.Println("VolUp")
	   })

	   phatbeat.On(phatbeat.Rewind, true, func(p int) {
	           log.Println("Rewind")
	   })

	   phatbeat.Hold(phatbeat.Rewind, false, 2, func(b int) {
	           log.Println("Rewind")
	   })

	   phatbeat.On(phatbeat.OnOff, true, func(p int) {
	           log.Println("OnOff")
	   })

	   phatbeat.Hold(phatbeat.OnOff, false, 2, func(b int) {
	           log.Println("OnOff")
	   })

	   defer phatbeat.Clean()

	   select {}

	*/
}

// basicAuthUser string, basicAuthPassword string

func newApp(debug bool) *iris.Application {

	api := iris.New()

	/*
	   authConfig := basicauth.Config{

	       Users:   map[string]string{basicAuthUser: basicAuthPassword},
	       Realm:   "Authorization Required",
	       Expires: time.Duration(1) * time.Minute,
	   }

	   authentication := basicauth.New(authConfig)

	   api.Use(authentication)
	*/

	// Set a rate-limit of 15 seconds to hold off on reloading albums/stories if RFID tag is held over the reader too long.
	limiter := tollbooth.NewLimiter(15, nil)

	database, _ := sql.Open("sqlite3", "./rfids.db")

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS rfid (id INTEGER PRIMARY KEY AUTOINCREMENT, tagid TEXT, uniqueid TEXT, url TEXT, playlistname TEXT)")

	statement.Exec()

	// statement, _ = database.Prepare("INSERT INTO rfid (tagid, uniqueid) VALUES (?, ?)")

	// statement.Exec("167697462420", "Nalini")

	// statement, _ = database.Prepare("UPDATE rfid SET tagid = '167697462420', playlistname = 'remotePlaylist' WHERE uniqueid = 'Nalini'")

	// statement.Exec()

	/*
	   UPDATE rfid
	   SET url = 'https://olm-build-artifacts.sfo2.cdn.digitaloceanspaces.com/track.mp3',
	   playlistname = 'remotePlaylist'
	   WHERE uniqueid = 'Nalini'
	*/
	api.Get("/currentstats/", func(ctx iris.Context) {

		trackName := GetCurrentTrackName()

		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "Current track name",
			"data":        trackName,
		})

	})

	api.Get("/stopcurrentplaylist/", func(ctx iris.Context) {

		StopPlaylist()
		currentPlayState := GetCurrentPlayState()
		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "Current play state",
			"data":        currentPlayState,
		})

	})

	api.Get("/playlistelapsedtime/", func(ctx iris.Context) {

		currentTrackTime := GetCurrentTrackTime()

		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "Current track time",
			"data":        currentTrackTime,
		})

	})

	api.Post("/rfid/", tollboothic.LimitHandler(limiter), func(ctx iris.Context) {

		rfid := new(RFID)

		err := ctx.ReadJSON(&rfid)

		if err != nil {
			if err.Error() == "unexpected end of JSON input" {
				ctx.StatusCode(400)
				ctx.JSON(iris.Map{
					"status_code": 400,
					"message":     "Malformed JSON input.",
				})
			} else if err.Error() == "invalid character '\"' after object key:value pair" {
				ctx.StatusCode(400)
				ctx.JSON(iris.Map{
					"status_code": 400,
					"message":     "Missing comma after object key:value pair in JSON input.",
				})
			}
		} else {

			if rfid.TagID == "" {
				// The TagID parameter in the JSON input is empty. Return a 422 error with appropriate message.
				ctx.StatusCode(422)

				ctx.JSON(iris.Map{
					"status_code": 422,
					"message":     "The TagID value was empty.",
				})

			} else if rfid.UniqueID == "" {
				// The UniqueID parameter in the JSON input is empty. Return a 422 error with appropriate message.

				ctx.StatusCode(422)

				ctx.JSON(iris.Map{
					"status_code": 422,
					"message":     "The UniqueID value was empty.",
				})

			} else {

				var id int = 0
				var tagid string = ""
				var uniqueid string = ""
				var url string = ""
				var playlistname string = ""

				sql := "SELECT id, tagid, uniqueid, url, playlistname FROM rfid WHERE tagid = '" + rfid.TagID + "' AND uniqueid = '" + rfid.UniqueID + "';"

				rows, _ := database.Query(sql)

				for rows.Next() {

					rows.Scan(&id, &tagid, &uniqueid, &url, &playlistname)

					fmt.Println(strconv.Itoa(id) + ": " + tagid + " " + uniqueid + " " + url + " " + playlistname)
				}

				if tagid != "" {
					playAknowledgeSound()
					rfid.ID = id
					rfid.TagID = tagid
					rfid.UniqueID = uniqueid
					rfid.URL = url
					rfid.PlaylistName = playlistname

					b, err := json.Marshal(rfid)

					if err != nil {
						fmt.Println(err)
						return
					}

					data := string(b)

					ctx.JSON(iris.Map{
						"status_code": 200,
						"message":     "INFO: The RFID tag exists. Playlist queued!",
						"data":        data,
					})
					ClearPlaylist()
					CreatePlaylist(rfid.URL, rfid.PlaylistName, ctx)
					LoadPlaylist(rfid.PlaylistName)
					PlayPlaylist(rfid.PlaylistName)

				} else {

					data := ""

					ctx.JSON(iris.Map{
						"status_code": 400,
						"message":     "The RFID tag was not found in the DB.",
						"data":        data,
					})

				}

			}

		}

	})

	// Create an api endpoint to create a new RFID tag.
	api.Post("/rfid/create/", tollboothic.LimitHandler(limiter), func(ctx iris.Context) {

		rfid := new(RFID)

		err := ctx.ReadJSON(&rfid)

		if err != nil {

			if err.Error() == "unexpected end of JSON input" {
				ctx.StatusCode(400)
				ctx.JSON(iris.Map{
					"status_code": 400,
					"message":     "Malformed JSON input.",
				})
			} else if err.Error() == "invalid character '\"' after object key:value pair" {
				ctx.StatusCode(400)
				ctx.JSON(iris.Map{
					"status_code": 400,
					"message":     "Missing comma after object key:value pair in JSON input.",
				})
			}

		} else {

			if rfid.TagID == "" {
				// The TagID parameter in the JSON input is empty. Return a 422 error with appropriate message.
				ctx.StatusCode(422)

				ctx.JSON(iris.Map{
					"status_code": 422,
					"message":     "The TagID value was empty.",
				})

			} else if rfid.UniqueID == "" {
				// The UniqueID parameter in the JSON input is empty. Return a 422 error with appropriate message.
				ctx.StatusCode(422)

				ctx.JSON(iris.Map{
					"status_code": 422,
					"message":     "The UniqueID value was empty.",
				})

			} else if rfid.URL == "" {
				// The URL parameter in the JSON input is empty. Return a 422 error with appropriate message.
				ctx.StatusCode(422)

				ctx.JSON(iris.Map{
					"status_code": 422,
					"message":     "The URL value was empty.",
				})

			} else if rfid.PlaylistName == "" {
				// The PlaylistName parameter in the JSON input is empty. Return a 422 error with appropriate message.
				ctx.StatusCode(422)

				ctx.JSON(iris.Map{
					"status_code": 422,
					"message":     "The PlaylistName value was empty.",
				})

			} else {

				// Create the RFID tag in the database.
				// If the RFID tag already exists in the database, return a 400 error with appropriate message.
				// If the RFID tag does not exist in the database, create the RFID tag in the database.
				// Return a 200 status code with appropriate message.

				var id int = 0
				var tagid string = ""
				var uniqueid string = ""
				var url string = ""
				var playlistname string = ""

				sql := "SELECT id, tagid, uniqueid, url, playlistname FROM rfid WHERE tagid = '" + rfid.TagID + "' AND uniqueid = '" + rfid.UniqueID + "';"

				rows := database.QueryRow(sql)

				err := rows.Scan(&id, &tagid, &uniqueid, &url, &playlistname)

				if err != nil {

					ctx.StatusCode(400)

					ctx.JSON(iris.Map{
						"status_code": 400,
						"message":     "Something went wrong with the database query. Please try again.",
					})

				} else {

					if tagid != "" {
						ctx.StatusCode(400)

						ctx.JSON(iris.Map{
							"status_code": 400,
							"message":     "The RFID tag already exists in the database.",
						})

					} else {

						sql := "INSERT INTO rfid (tagid, uniqueid, url, playlistname) VALUES ('" + rfid.TagID + "', '" + rfid.UniqueID + "', '" + rfid.URL + "', '" + rfid.PlaylistName + "');"

						_, err := database.Exec(sql)

						if err != nil {

							ctx.StatusCode(400)

							ctx.JSON(iris.Map{
								"status_code": 400,
								"message":     "Something went wrong with the database query. Please try again.",
							})

						} else {

							ctx.StatusCode(200)

							ctx.JSON(iris.Map{
								"status_code": 200,
								"message":     "The RFID tag was successfully created in the database.",
							})

						}

					}

				}

			}

		}
	})

	return api
}

// dbConn() (db *sql.DB) initializes a single connection to the database.
func dbConn() (database *sql.DB) {

	database, err := sql.Open("sqlite3", "./rfids.db")

	if err != nil {
		panic(err.Error())
	}

	return database
}

// CreatePlaylist creates a new playlist in the database.
func CreatePlaylist(url string, playlistname string, ctx iris.Context) {

	// Create the playlist in the database.
	// If the playlist already exists in the database, return a 400 error with appropriate message.
	// If the playlist does not exist in the database, create the playlist in the database.
	// Return a 200 status code with appropriate message.

	var id int = 0

	sql := "SELECT id, url, playlistname FROM playlist WHERE url = '" + url + "' AND playlistname = '" + playlistname + "';"

	database := dbConn()

	rows := database.QueryRow(sql)

	err := rows.Scan(&id, &url, &playlistname)

	if err != nil {
		ctx.StatusCode(400)

		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "Something went wrong with the database query. Please try again.",
		})

	} else {

		if url != "" {

			sql := "INSERT INTO playlist (url, playlistname) VALUES ('" + url + "', '" + playlistname + "');"

			_, err := database.Exec(sql)

			if err != nil {
				fmt.Println(err)
				return
			}

		} else {

			ctx.StatusCode(400)

			ctx.JSON(iris.Map{
				"status_code": 400,
				"message":     "The playlist already exists in the database.",
			})

		}

	}

}

// DeletePlaylist deletes a playlist from the database.
func DeletePlaylist(url string, playlistname string, ctx iris.Context) {

	// Delete the playlist from the database.
	// If the playlist does not exist in the database, return a 400 error with appropriate message.
	// If the playlist exists in the database, delete the playlist from the database.
	// Return a 200 status code with appropriate message.

	var id int = 0

	sqlSelect := "SELECT id, url, playlistname FROM playlist WHERE url = '" + url + "' AND playlistname = '" + playlistname + "';"

	database := dbConn()

	rows := database.QueryRow(sqlSelect)

	err := rows.Scan(&id, &url, &playlistname)

	if err != nil {

		ctx.StatusCode(400)

		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "Something went wrong with the database query. Please try again.",
		})

	} else {

		if url != "" {

			sqlDelete := "DELETE FROM playlist WHERE url = '" + url + "' AND playlistname = '" + playlistname + "';"

			// TODO - Delete the playlist from the database.

		} else {

			ctx.StatusCode(400)

			ctx.JSON(iris.Map{
				"status_code": 400,
				"message":     "The playlist does not exist in the database.",
			})

		}

	}

}

// GetPlaylist gets a playlist from the database.
func GetPlaylist(url string, playlistname string, ctx iris.Context) {

	// Get the playlist from the database.
	// If the playlist does not exist in the database, return a 400 error with appropriate message.
	// If the playlist exists in the database, return the playlist from the database.
	// Return a 200 status code with appropriate message.

	var id int = 0

	sql := "SELECT id, url, playlistname FROM playlist WHERE url = '" + url + "' AND playlistname = '" + playlistname + "';"

	database := dbConn()

	rows := database.QueryRow(sql)

	err := rows.Scan(&id, &url, &playlistname)

	if err != nil {

		ctx.StatusCode(400)

		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "Something went wrong with the database query. Please try again.",
		})

	} else {

		if url != "" {

			sql := "SELECT id, url, playlistname FROM playlist WHERE url = '" + url + "' AND playlistname = '" + playlistname + "';"

			rows := database.QueryRow(sql)

			err := rows.Scan(&id, &url, &playlistname)

			if err != nil {
				fmt.Println(err)
				return
			}

		} else {

			ctx.StatusCode(400)

			ctx.JSON(iris.Map{
				"status_code": 400,
				"message":     "The playlist does not exist in the database.",
			})

		}

	}

}

// Create a function called playErrorSound() that plays the error sound.
func playErrorSound() {

	cmd := "mpg123-alsa"

	// Final location: /etc/sound/error.mp3
	errorSoundFile := "subtleErrorBell.mp3"

	args := []string{errorSoundFile}
	if err := exec.Command(cmd, args...).Run(); err != nil {

		log.Println(os.Stderr, err)

	}

}

func playReadySound() {

	cmd := "mpg123-alsa"

	// Final location: /etc/sound/started.mp3
	startupSoundFile := "ready.mp3"

	args := []string{startupSoundFile}
	if err := exec.Command(cmd, args...).Run(); err != nil {

		log.Println(os.Stderr, err)

	}

}

func playAknowledgeSound() {

	cmd := "mpg123-alsa"

	// Final location: /etc/sound/started.mp3
	startupSoundFile := "intuition.mp3"

	args := []string{startupSoundFile}
	if err := exec.Command(cmd, args...).Run(); err != nil {

		log.Println(os.Stderr, err)

	}

}

func playShutdownSound() {

	cmd := "mpg123-alsa"

	// Final location: /etc/sound/started.mp3
	startupSoundFile := "shutdown.mp3"

	args := []string{startupSoundFile}
	if err := exec.Command(cmd, args...).Run(); err != nil {

		log.Println(os.Stderr, err)

	}

}

// Need a power management board for this functionality MoPi 2 or ...
func playLowBatterySound(batteryLevel int) {

	cmd := "gtts-cli"
	batteryLevelString := strconv.Itoa(batteryLevel)
	// Final location: /etc/sound/started.mp3
	batteryMessage := "\"The battery is at, " + batteryLevelString + " percent!\""

	args := []string{batteryMessage, "|", "mpg123-alsa", "-"}
	if err := exec.Command(cmd, args...).Run(); err != nil {

		log.Println(os.Stderr, err)

	}
}

/*
func CreatePlaylist(track string, playlist string) {

	playlistFile := "/var/lib/mpd/playlists/" + playlist + ".m3u"

	f, err := os.Create(playlistFile)
	if err != nil {
		log.Println(err)
		return
	}
	l, err := f.WriteString(track)
	if err != nil {
		log.Println(err)
		f.Close()
		return
	}
	log.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		log.Println(err)
		return
	}

}
*/

/*
sudo pip install gTTS
CLI: gtts-cli "come for dinner" | mpg123 -
*/

// Use the Google text to speech engine library at this location https://github.com/pndurette/gTTS to play a custom message.
// The function should take a string as an argument and play the string as a message.
func playCustomMessage(message string) {

	cmd := "gtts-cli"

	args := []string{message, "|", "mpg123-alsa", "-"}
	if err := exec.Command(cmd, args...).Run(); err != nil {

		log.Println(os.Stderr, err)

	}

}

func ClearPlaylist() {

	cmdClear := "mpc"

	argsClear := []string{"--host", "alraune22@localhost", "clear"}

	if err := exec.Command(cmdClear, argsClear...).Run(); err != nil {

		log.Println(os.Stderr, err)

		playErrorSound()
		playCustomMessage("The playlist could not be cleared. Please try again.")

	} else {

		playAknowledgeSound()
		playCustomMessage("The playlist has been cleared.")

	}

	cmdUpdate := "mpc"

	argsUpdate := []string{"--host", "alraune22@localhost", "update"}

	if err := exec.Command(cmdUpdate, argsUpdate...).Run(); err != nil {

		log.Println(os.Stderr, err)

	}

}

func LoadPlaylist(playlist string) {

	ClearPlaylist()

	cmdLoad := "mpc"

	argsLoad := []string{"--host", "alraune22@localhost", "load", playlist}

	if err := exec.Command(cmdLoad, argsLoad...).Run(); err != nil {

		log.Println(os.Stderr, err)

		playErrorSound()
		playCustomMessage("The playlist could not be loaded. Please try again.")

	} else {

		playAknowledgeSound()
		playCustomMessage("The playlist has been loaded.")
	}

}

func PlayPlaylist(playlist string) {

	cmd := "mpc"

	args := []string{"--host", "alraune22@localhost", "play"}

	if err := exec.Command(cmd, args...).Run(); err != nil {

		log.Println(os.Stderr, err)

		playErrorSound()
		playCustomMessage("The playlist could not be played. Please try again.")

	} else {

		log.Println("Playing playlist: ", playlist)

	}

}

func PausePlaylist(playlist string) {

	cmd := "mpc"

	args := []string{"--host", "alraune22@localhost", "pause"}

	if err := exec.Command(cmd, args...).Run(); err != nil {

		log.Println(os.Stderr, err)
		playErrorSound()
		playCustomMessage("The playlist could not be paused. Please try again.")

	} else {

		log.Println("Pausing playlist: ", playlist)
		playAknowledgeSound()
		playCustomMessage("The playlist has been paused.")
	}

}

func StopPlaylist() {

	cmd := "mpc"

	args := []string{"--host", "alraune22@localhost", "stop"}

	if err := exec.Command(cmd, args...).Run(); err != nil {

		log.Println(os.Stderr, err)
		playErrorSound()
		playCustomMessage("The playlist could not be stopped. Please try again.")

	} else {

		log.Println("Stopping playlist")
		playAknowledgeSound()
		playCustomMessage("The playlist has been stopped.")

	}

}

func GetCurrentTrackName() (currentTrackName string) {

	currentTrackName = ""

	cmdGetCurrentTrackName := "bin/mpdcurrentsong"

	execCmd := exec.Command(cmdGetCurrentTrackName)

	var out bytes.Buffer
	var stderr bytes.Buffer

	execCmd.Stdout = &out
	execCmd.Stderr = &stderr

	err := execCmd.Run()

	if err != nil {
		log.Println(os.Stderr, err)
		currentTrackName = "ERROR"

		playErrorSound()
		playCustomMessage("The current track could not be retrieved. Please try again.")

	} else {
		currentTrackName = strings.Trim(out.String(), "\n")
	}

	return currentTrackName

}

func GetCurrentPlayState() (currentPlayState string) {

	cmdPlayState := "bin/mpdplaystate"

	execCmd := exec.Command(cmdPlayState)

	var out bytes.Buffer
	var stderr bytes.Buffer

	execCmd.Stdout = &out
	execCmd.Stderr = &stderr

	err := execCmd.Run()

	if err != nil {
		log.Println(os.Stderr, err)
		currentPlayState = "ERROR"
	} else {
		// 0 = State unknown, 1 = STATE STOP, 2 = STATE PLAY, 3 = STATE PAUSE,
		currentPlayState = strings.Trim(out.String(), "\n")

	}

	return currentPlayState
}

func GetCurrentTrackTime() (elapsedTime string) {

	cmdGetCurrentTrack := "bin/mpdtime"

	execCmd := exec.Command(cmdGetCurrentTrack)

	var out bytes.Buffer
	var stderr bytes.Buffer

	execCmd.Stdout = &out
	execCmd.Stderr = &stderr

	err := execCmd.Run()

	if err != nil {
		log.Println(os.Stderr, err)
		elapsedTime = "ERROR"
	} else {

		elapsedTime = strings.Trim(out.String(), "\n")

	}

	return elapsedTime
}

func runShutDownSequence() {

}

// Get preferred outbound ip of this machine
func GetOutboundIP() (string, error) {

	connection, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}

	defer connection.Close()

	localAddr := connection.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

var SdNotifyNoSocket = errors.New("No socket")

// SdNotify sends a message to the init daemon. It is common to ignore the error.
func SdNotify(state string) error {
	socketAddr := &net.UnixAddr{
		Name: os.Getenv("NOTIFY_SOCKET"),
		Net:  "unixgram",
	}

	if socketAddr.Name == "" {
		return SdNotifyNoSocket
	}

	conn, err := net.DialUnix(socketAddr.Net, nil, socketAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(state))
	return err
}

func HandleError(errorConst string, errors error, ctx iris.Context) {

	_, err := ctx.HTML(errorConst + errors.Error() + "</b>")
	if err != nil {
		log.Println(err.Error())
	}
}
