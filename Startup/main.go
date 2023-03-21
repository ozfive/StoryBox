package startup

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	playStartUpSound()
}

func playStartUpSound() {
	cmd := "mpg123-alsa"

	startupSoundFile := "/etc/sound/started.mp3"

	args := []string{startupSoundFile}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		logFile, err := os.OpenFile(filepath.Join(os.Getenv("XDG_DATA_HOME"), "storybox", "logs", "startup-errors.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Failed to open log file:", err)
			os.Exit(1)
		}

		defer logFile.Close()

		logger := log.New(logFile, "", log.LstdFlags)
		logger.Println(fmt.Errorf("failed to play startup sound: %v", err))

		os.Exit(1)
	}
}
