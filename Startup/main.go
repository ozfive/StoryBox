package startup

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adrg/xdg"
)

func main() {
	// Create and configure startup logger file.
	err := createLogFile(getLogFilePath("startup-errors.log"))
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	// Play the startup sound.
	playStartUpSound()
}

func playStartUpSound() {
	cmd := "mpg123-alsa"
	startupSoundFile := "/etc/sound/started.mp3"
	args := []string{startupSoundFile}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.Println(fmt.Errorf("failed to play startup sound: %v", err))
		os.Exit(1)
	}
}

func createLogFile(logFilePath string) error {
	// Create log directory
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Create log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer logFile.Close()

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
