package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/adrg/xdg"
	"github.com/didip/tollbooth/v6"
	"github.com/iris-contrib/middleware/tollboothic"
	"github.com/kataras/iris/v12"
	_ "github.com/mattn/go-sqlite3"
)

type RFID struct {
	ID           int    `json:"id"`
	TagID        string `json:"tagid"`
	UniqueID     string `json:"uniqueid"`
	URL          string `json:"url"`
	PlaylistName string `json:"playlistname"`
}

func main() {

	// Create and configure logger
	err := createLogFile(getLogFilePath("system-errors.log"))
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	debug := false

	localIPAddress := "localhost"
	localIPPort := "3001"

	// basicAuthUser := "admin"
	// basicAuthPassword := "password"

	// basicAuthUser, basicAuthPassword

	api := newApp(debug)

	errs := api.Run(iris.Addr(localIPAddress+":"+localIPPort), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
	if err != nil {
		log.Println(errs.Error())
	}
}

// basicAuthUser string, basicAuthPassword string
// statement, _ = database.Prepare("INSERT INTO rfid (tagid, uniqueid) VALUES (?, ?)")
// statement.Exec("167697462420", "Nalini")
// statement, _ = database.Prepare("UPDATE rfid SET tagid = '167697462420', playlistname = 'remotePlaylist' WHERE uniqueid = 'Nalini'")
// statement.Exec()

/*
	UPDATE rfid
	SET url = '',
	playlistname = 'remotePlaylist'
	WHERE uniqueid = 'Nalini'
*/

func newApp(debug bool) *iris.Application {
	api := iris.New()

	limiter := tollbooth.NewLimiter(15, nil)

	database, _ := sql.Open("sqlite3", "rfids.db")

	createDatabaseTable(database)

	api.Get("/currentstats/", currentStatsHandler)
	api.Get("/stopcurrentplaylist/", stopCurrentPlaylistHandler)
	api.Get("/playlistelapsedtime/", playlistElapsedTimeHandler)
	api.Post("/rfid/", tollboothic.LimitHandler(limiter), rfidHandler(database))
	api.Post("/rfid/create/", tollboothic.LimitHandler(limiter), rfidCreateHandler(database))

	return api
}

func createDatabaseTable(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS rfid (id INTEGER PRIMARY KEY AUTOINCREMENT, tagid TEXT, uniqueid TEXT, url TEXT, playlistname TEXT)")
	statement.Exec()
}

func currentStatsHandler(ctx iris.Context) {
	handleStatsEndpoint(ctx, getCurrentTrackNameWrapper, "track_name", "Failed to get current track name")
	handleStatsEndpoint(ctx, getCurrentTrackTimeWrapper, "elapsed_time", "Failed to get current track time")
	handleStatsEndpoint(ctx, getCurrentVolumeWrapper, "volume", "Failed to get current volume")
}

func getCurrentTrackNameWrapper() (interface{}, error) {
	return GetCurrentTrackName()
}

func getCurrentVolumeWrapper() (interface{}, error) {
	return GetCurrentVolume()
}

func stopCurrentPlaylistHandler(ctx iris.Context) {
	stopPlaylistWrapper := func() (interface{}, error) {
		err := StopPlaylist()
		return nil, err
	}
	getCurrentPlayStateWrapper := func() (interface{}, error) {
		state, err := GetCurrentPlayState()
		return state, err
	}

	handleActionEndpoint(ctx, stopPlaylistWrapper, "Failed to stop playlist")
	handleStatsEndpoint(ctx, getCurrentPlayStateWrapper, "play_state", "Failed to get current play state")
}

func playlistElapsedTimeHandler(ctx iris.Context) {
	handleStatsEndpoint(ctx, getCurrentTrackTimeWrapper, "elapsed_time", "Failed to get current track time")
}

func getCurrentTrackTimeWrapper() (interface{}, error) {
	return GetCurrentTrackTime()
}

func rfidHandler(database *sql.DB) iris.Handler {
	return func(ctx iris.Context) {
		rfid := new(RFID)
		err := ctx.ReadJSON(&rfid)
		handleInputErrors(ctx, err, rfid)
		handleRFIDTag(ctx, database, rfid, false)
	}
}

func rfidCreateHandler(database *sql.DB) iris.Handler {
	return func(ctx iris.Context) {
		rfid := new(RFID)
		err := ctx.ReadJSON(&rfid)
		handleInputErrors(ctx, err, rfid)
		handleRFIDCreation(ctx, database, rfid)
	}
}

func handleStatsEndpoint(ctx iris.Context, getStats func() (interface{}, error), statKey, errMsg string) {
	stats, err := getStats()
	if err != nil {
		ctx.StatusCode(500)
		ctx.JSON(iris.Map{
			"status_code": 500,
			"message":     errMsg,
			"data":        err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"status_code": 200,
		"message":     statKey,
		"data":        stats,
	})
}

func handleActionEndpoint(ctx iris.Context, actionFunc func() (interface{}, error), successMessage string) {
	result, err := actionFunc()
	if err != nil {
		ctx.StatusCode(500)
		ctx.JSON(iris.Map{
			"status_code": 500,
			"message":     fmt.Sprintf("Failed to %s", successMessage),
			"data":        err.Error(),
		})
		return
	}
	ctx.JSON(iris.Map{
		"status_code": 200,
		"message":     successMessage,
		"data":        result,
	})
}

func handleInputErrors(ctx iris.Context, err error, rfid *RFID) bool {
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
		return true
	}

	if rfid.TagID == "" || rfid.UniqueID == "" || rfid.URL == "" || rfid.PlaylistName == "" {
		ctx.StatusCode(422)
		missingField := ""
		if rfid.TagID == "" {
			missingField = "TagID"
		} else if rfid.UniqueID == "" {
			missingField = "UniqueID"
		} else if rfid.URL == "" {
			missingField = "URL"
		} else if rfid.PlaylistName == "" {
			missingField = "PlaylistName"
		}
		ctx.JSON(iris.Map{
			"status_code": 422,
			"message":     fmt.Sprintf("The %s value was empty.", missingField),
		})
		return true
	}

	return false
}

func handleRFIDTag(ctx iris.Context, database *sql.DB, rfid *RFID, existingTag bool) {
	var id int
	var tagid, uniqueid, url, playlistname string

	stmt, err := database.Prepare("SELECT id, tagid, uniqueid, url, playlistname FROM rfid WHERE tagid = ? AND uniqueid = ?;")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(rfid.TagID, rfid.UniqueID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &tagid, &uniqueid, &url, &playlistname)
		if err != nil {
			panic(err)
		}
		fmt.Println(strconv.Itoa(id) + ": " + tagid + " " + uniqueid + " " + url + " " + playlistname)
	}

	if tagid != "" && existingTag {
		playAknowledgeNotification()
		rfid.ID = id
		rfid.TagID = tagid
		rfid.UniqueID = uniqueid
		rfid.URL = url
		rfid.PlaylistName = playlistname

		data, _ := json.Marshal(rfid)

		ctx.JSON(iris.Map{
			"status_code": 200,
			"message":     "INFO: The RFID tag exists. Playlist queued!",
			"data":        string(data),
		})
		ClearPlaylist()
		CreatePlaylist(rfid.URL, rfid.PlaylistName, ctx)
		LoadPlaylist(rfid.PlaylistName)
		PlayPlaylist(rfid.PlaylistName)
	} else if tagid == "" && !existingTag {
		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "The RFID tag was not found in the DB.",
			"data":        "",
		})
	}
}

// TODO: Add transaction rollback functionality
func handleRFIDCreation(ctx iris.Context, database *sql.DB, rfid *RFID) {
	var id int
	var tagid, uniqueid, url, playlistname string

	stmtSelect, err := database.Prepare("SELECT id, tagid, uniqueid, url, playlistname FROM rfid WHERE tagid = ? AND uniqueid = ?;")
	if err != nil {
		panic(err)
	}
	defer stmtSelect.Close()

	err = stmtSelect.QueryRow(rfid.TagID, rfid.UniqueID).Scan(&id, &tagid, &uniqueid, &url, &playlistname)

	if err == nil && tagid != "" {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "The RFID tag already exists in the database.",
		})
	} else {
		stmtInsert, err := database.Prepare("INSERT INTO rfid (tagid, uniqueid, url, playlistname) VALUES (?, ?, ?, ?);")
		if err != nil {
			panic(err)
		}
		defer stmtInsert.Close()

		_, err = stmtInsert.Exec(rfid.TagID, rfid.UniqueID, rfid.URL, rfid.PlaylistName)

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

// connectToDatabase initializes a single connection to the database.
func connectToDatabase() (database *sql.DB, err error) {
	// Open the SQLite database file
	database, err = sql.Open("sqlite3", "./rfids.db")
	if err != nil {
		return nil, fmt.Errorf("unable to open SQLite database: %w", err)
	}

	// Ping the database to check the connection
	err = database.Ping()
	if err != nil {
		// Close the database before returning an error, since we opened it successfully
		database.Close()
		return nil, fmt.Errorf("unable to establish a connection to the SQLite database: %w", err)
	}

	return database, nil
}

// createLogFile creates a log file at the specified path and sets the logger output to it.
// The function ensures that the directory structure is created and returns an error if anything fails.
func createLogFile(logFilePath string) error {
	if logFilePath == "" {
		return fmt.Errorf("log file path is empty")
	}

	// Create log directory
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Check if the log file already exists
	if _, err := os.Stat(logFilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if log file exists: %w", err)
	}

	// Create log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}

	defer func() {
		if err := logFile.Close(); err != nil {
			log.Printf("failed to close log file: %v", err)
		}
	}()

	// Set logger output
	log.SetOutput(logFile)
	return nil
}

func getLogFilePath(logFileName string) string {
	// Get XDG_DATA_HOME directory
	dataDir := xdg.DataHome

	// Create log file path
	return filepath.Join(dataDir, "storybox", "logs", logFileName)
}

// CreatePlaylist creates a new playlist in the database.
func CreatePlaylist(url string, playlistname string, ctx iris.Context) {
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
}

// DeletePlaylist deletes a playlist from the database.
func DeletePlaylist(url string, playlistname string, ctx iris.Context) {
	database, err := connectToDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer database.Close()

	// Check if the playlist exists in the database.
	var count int
	sqlCheck := "SELECT COUNT(*) FROM playlist WHERE url = ? AND playlistname = ?"
	selectErr := database.QueryRow(sqlCheck, url, playlistname).Scan(&count)

	if selectErr != nil {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "Failed to SELECT playlist " + playlistname + " from the database. Please try again.",
		})
		return
	}

	if count == 0 {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "The playlist " + playlistname + " does not exist in the database.",
		})
		return
	}

	// Delete the playlist from the database.
	sqlDelete := "DELETE FROM playlist WHERE url = ? AND playlistname = ?"
	_, deleteErr := database.Exec(sqlDelete, url, playlistname)

	if deleteErr != nil {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{
			"status_code": 400,
			"message":     "Failed to delete playlist " + playlistname + " from the database. Please try again.",
		})
		return
	}

	// Return a success message.
	ctx.StatusCode(200)
	ctx.JSON(iris.Map{
		"status_code": 200,
		"message":     "The playlist was deleted successfully.",
	})
}

type Playlist struct {
	ID             int
	URLFromDB      string
	PlaylistNameDB string
	Err            error
}

// GetPlaylist gets a playlist from the database.
func GetPlaylist(url, playlistname string) Playlist {

	sql := "SELECT id, url, playlistname FROM playlist WHERE url = ? AND playlistname = ?;"

	database, err := connectToDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer database.Close()

	row := database.QueryRow(sql, url, playlistname)

	var playlist Playlist

	errs := row.Scan(&playlist.ID, &playlist.URLFromDB, &playlist.PlaylistNameDB)
	if errs != nil {
		playlist.Err = fmt.Errorf("failed to retrieve playlist: %v", errs)
	}

	return playlist
}

type Sound struct {
	File string
	Err  error
}

func playSound(s Sound) {
	cmd := "mpg123-alsa"
	args := []string{s.File}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		s.Err = fmt.Errorf("failed to play sound: %v", err)
		log.Println(s.Err)
	}
}

// playErrorNotification plays the error notification.
func playErrorNotification() Sound {
	sound := Sound{File: "/etc/sound/subtleErrorBell.mp3"}
	playSound(sound)
	return sound
}

// playReadyNotification plays the ready notification.
func playReadyNotification() Sound {
	sound := Sound{File: "/etc/sound/ready.mp3"}
	playSound(sound)
	return sound
}

// playAknowledgeNotification plays the aknowledge notification.
func playAknowledgeNotification() Sound {
	sound := Sound{File: "/etc/sound/intuition.mp3"}
	playSound(sound)
	return sound
}

// playShutdownNotification plays the shutdown notification.
func playShutdownNotification() Sound {
	sound := Sound{File: "/etc/sound/shutdown.mp3"}
	playSound(sound)
	return sound
}

func playLowBatteryNotification(batteryLevel int) Sound {
	sound := generateBatteryMessage(batteryLevel)
	playSound(sound)
	return sound
}

func generateBatteryMessage(batteryLevel int) Sound {
	sound := Sound{File: "batteryMessage.mp3"}
	cmd := "gtts-cli"
	batteryLevelString := strconv.Itoa(batteryLevel)
	message := "The battery is at " + batteryLevelString + " percent!"
	args := []string{"-o", sound.File, message}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		sound.Err = fmt.Errorf("failed to generate battery message: %v", err)
		log.Println(sound.Err)
	}
	return sound
}

/*
func CreatePlaylist(track string, playlist string) error {
	playlistFile := "/var/lib/mpd/playlists/" + playlist + ".m3u"
	f, err := os.Create(playlistFile)
	if err != nil {
		return fmt.Errorf("failed to create playlist file: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(track); err != nil {
		return fmt.Errorf("failed to write track to playlist file: %v", err)
	}

	return nil
}
*/

/*
sudo pip install gTTS
CLI: gtts-cli "come for dinner" | mpg123 -
*/

// Use the Google text to speech engine library at this location
// https://github.com/pndurette/gTTS to play a custom message.
// The function should take a string as an argument and play
// the string as a message.
func playCustomMessage(message string) {

	cmd := "gtts-cli"

	args := []string{message, "|", "mpg123-alsa", "-"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.Println(fmt.Errorf("Failed to play 'Error' notification: %v", err))
	}

}

// playCustomMessageFromGCloud uses the Google text to speech engine to play a custom message.
/*
func playCustomMessageFromGCloud(message string) {
	// Set the environment variable to the path of your JSON key file
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "keyfile.json")

	ctx := context.Background()

	client, err := tts.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create texttospeech client: %v", err)
		return
	}
	defer client.Close()
		// texttospeechpb
		req := &tts.SynthesizeSpeechRequest{
			Input: &tts.SynthesisInput{
				InputSource: &tts.SynthesisInput_Text{Text: message},
			},
			Voice: &tts.VoiceSelectionParams{
				LanguageCode: "en-US",
				SsmlGender:   tts.SsmlVoiceGender_FEMALE,
			},
			AudioConfig: &tts.AudioConfig{
				AudioEncoding: tts.AudioEncoding_MP3,
			},
		}

		resp, err := client.SynthesizeSpeech(ctx, req)
		if err != nil {
			log.Printf("Failed to synthesize speech: %v", err)
			return
		}

	err = ioutil.WriteFile("output.mp3", resp.AudioContent, 0644)
	if err != nil {
		log.Printf("Failed to write audio content: %v", err)
		return
	}

	cmd := "mpg123"
	args := []string{"output.mp3"}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.Printf("Failed to play custom message: %v", err)
	}
}
*/
func ClearPlaylist() error {
	const (
		mpcCmd     = "mpc"
		mpcClear   = "clear"
		mpcUpdate  = "update"
		mpcHostArg = "--host"
		mpcHost    = "alraune22@localhost"
	)

	clearArgs := []string{mpcHostArg, mpcHost, mpcClear}

	if err := exec.Command(mpcCmd, clearArgs...).Run(); err != nil {
		sound := playErrorNotification()
		fmt.Println(sound)
		playCustomMessage("The playlist could not be cleared. Please try again.")
		return fmt.Errorf("failed to clear playlist: %v", err)
	}

	playAknowledgeNotification()

	playCustomMessage("The playlist has been cleared.")

	updateArgs := []string{mpcHostArg, mpcHost, mpcUpdate}

	if err := exec.Command(mpcCmd, updateArgs...).Run(); err != nil {
		log.Println(fmt.Errorf("Failed to update the music database: %v", err))
	}

	return nil
}

func LoadPlaylist(playlist string) error {
	const (
		mpcCmd     = "mpc"
		mpcLoad    = "load"
		mpcHostArg = "--host"
		mpcHost    = "alraune22@localhost"
	)
	if err := ClearPlaylist(); err != nil {
		return fmt.Errorf("failed to clear playlist before loading: %v", err)
	}

	loadArgs := []string{mpcHostArg, mpcHost, mpcLoad, playlist}
	if err := exec.Command(mpcCmd, loadArgs...).Run(); err != nil {
		playErrorNotification()
		playCustomMessage("The playlist could not be loaded. Please try again.")
		return fmt.Errorf("failed to load playlist: %v", err)
	}

	playAknowledgeNotification()
	playCustomMessage("The playlist has been loaded.")
	return nil
}

// PlayPlaylist clears the current playlist, loads the specified playlist, and starts playing it.
// Returns an error if there is an issue with clearing, loading, or playing the playlist.
func PlayPlaylist(playlist string) error {
	// Clear the current playlist.
	cmd := "mpc"
	args := []string{"--host", "alraune22@localhost", "clear"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		return fmt.Errorf("Failed to clear playlist %s: %v", playlist, err)
	}

	// Load the specified playlist
	args = []string{"--host", "alraune22@localhost", "load", playlist}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		return fmt.Errorf("Failed to load playlist %s: %v", playlist, err)
	}

	// Start playing the playlist
	args = []string{"--host", "alraune22@localhost", "play"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		return fmt.Errorf("Failed to play playlist %s: %v", playlist, err)
	}

	log.Println("Playing playlist: ", playlist)
	return nil
}

// PausePlaylist pauses the specified playlist.
func PausePlaylist(playlist string) error {
	// Define the command and arguments to be executed
	cmd := "mpc"
	args := []string{"--host", "alraune22@localhost", "pause"}

	// Execute the command and check for errors
	if err := exec.Command(cmd, args...).Run(); err != nil {
		return fmt.Errorf("Failed to pause playlist: %v", err)
	}

	// Log the success and play the acknowledgement notification
	log.Println("Pausing playlist: ", playlist)
	playAknowledgeNotification()
	playCustomMessage("The playlist has been paused.")

	return nil
}

// StopPlaylist stops the current playlist.
func StopPlaylist() error {
	cmd := "mpc"
	args := []string{"--host", "alraune22@localhost", "stop"}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		return fmt.Errorf("failed to stop playlist: %v", err)
	}

	log.Println("Stopping playlist")
	playAknowledgeNotification()
	playCustomMessage("The playlist has been stopped.")

	return nil
}

func GetCurrentTrackName() (string, error) {
	cmd := "/usr/local/bin/mpdcurrentsong"
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Println(err)
		return "", err
	}
	currentTrackName := strings.TrimSpace(string(out))
	return currentTrackName, nil
}

func GetCurrentPlayState() (string, error) {
	cmdPlayState := "/usr/local/bin/mpdplaystate"

	execCmd := exec.Command(cmdPlayState)

	var out bytes.Buffer
	var stderr bytes.Buffer

	execCmd.Stdout = &out
	execCmd.Stderr = &stderr

	if err := execCmd.Run(); err != nil {
		log.Println(err)
		return "", fmt.Errorf("failed to get current play state: %v", err)
	}

	currentPlayState := strings.TrimSpace(out.String())
	return currentPlayState, nil
}

// GetCurrentTrackTime returns the current elapsed time of the track being played by the MPD server.
// It runs the `mpdtime` command and returns the elapsed time as a string.
// If there is an error with the command, it returns an empty string and an error message.
func GetCurrentTrackTime() (string, error) {
	cmd := "/usr/local/bin/mpdtime"

	execCmd := exec.Command(cmd)

	var out bytes.Buffer
	var stderr bytes.Buffer

	// Set the output stream to the `out` buffer and the error stream to the `stderr` buffer.
	execCmd.Stdout = &out
	execCmd.Stderr = &stderr

	// Run the command and check for errors.
	err := execCmd.Run()
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("Failed to get current track time: %v", err)
	}

	// Trim the newline character from the output and return the elapsed time as a string.
	elapsedTime := strings.Trim(out.String(), "\n")
	return elapsedTime, nil
}

func GetCurrentVolume() (int, error) {
	cmd := "mpc"
	args := []string{"--host", "alraune22@localhost", "volume"}

	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return 0, fmt.Errorf("Failed to get current volume: %v", err)
	}

	// The output of the "mpc volume" command is in the format "volume: N%", where N is the volume level.
	// We extract the volume level using a regular expression.
	re := regexp.MustCompile(`volume:\s+(\d+)%`)
	match := re.FindStringSubmatch(string(output))
	if len(match) < 2 {
		return 0, errors.New("Failed to extract volume level from command output")
	}

	volume, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("Failed to convert volume to integer: %v", err)
	}

	return volume, nil
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
