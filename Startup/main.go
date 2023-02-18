package main

import (
	"os"
	"os/exec"
	"log" 
)

func playStartUpSound() {

	cmd := "mpg123-alsa"

	// Final location: /etc/sound/started.mp3
	startupSoundFile := "/home/pi/go/src/Storybox/Startup/started.mp3"

	args := []string{startupSoundFile}
	if err := exec.Command(cmd, args...).Run(); err != nil {

		log.Println(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {

	playStartUpSound()

}