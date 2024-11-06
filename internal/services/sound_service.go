package services

import (
	"fmt"
	"log"
	"os/exec"
	// "StoryBox/internal/models"
)

type SoundService interface {
	PlayErrorNotification() error
	PlayReadyNotification() error
	PlayAcknowledgeNotification() error
	PlayCustomMessage(message string) error
	PlayLowBatteryNotification(batteryLevel int) error
}

type soundService struct{}

func NewSoundService() SoundService {
	return &soundService{}
}

func (s *soundService) PlayErrorNotification() error {
	return s.playSound("/etc/sound/subtleErrorBell.mp3")
}

func (s *soundService) PlayReadyNotification() error {
	return s.playSound("/etc/sound/ready.mp3")
}

func (s *soundService) PlayAcknowledgeNotification() error {
	return s.playSound("/etc/sound/intuition.mp3")
}

func (s *soundService) PlayCustomMessage(message string) error {
	cmd := "gtts-cli"
	args := []string{message, "|", "mpg123-alsa", "-"}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.Println(fmt.Errorf("failed to play custom message: %v", err))
		return err
	}

	return nil
}

func (s *soundService) PlayLowBatteryNotification(batteryLevel int) error {
	message := fmt.Sprintf("The battery is at %d percent!", batteryLevel)
	return s.PlayCustomMessage(message)
}

func (s *soundService) playSound(filePath string) error {
	cmd := "mpg123-alsa"
	args := []string{filePath}

	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.Println(fmt.Errorf("failed to play sound %s: %v", filePath, err))
		return err
	}

	return nil
}
