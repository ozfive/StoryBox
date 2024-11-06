package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	// "github.com/adrg/xdg"
)

func InitLogger(logFilePath string) error {
	if logFilePath == "" {
		return fmt.Errorf("log file path is empty")
	}

	// Create log directory
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Set logger output
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return nil
}
