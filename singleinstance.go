package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// lockFile represents the lock file path
var lockFile string

// initSingleInstance initializes the single instance lock
func initSingleInstance() error {
	// Get user's temp directory
	tmpDir := os.TempDir()
	lockFile = filepath.Join(tmpDir, "myWeatherApp.lock")

	// Check if lock file exists
	if _, err := os.Stat(lockFile); err == nil {
		// Lock file exists, check if process is running
		data, err := os.ReadFile(lockFile)
		if err == nil {
			fmt.Printf("Another instance is already running (PID: %s)\n", string(data))
			return fmt.Errorf("application is already running")
		}
	}

	// Create lock file with current PID
	pid := fmt.Sprintf("%d", os.Getpid())
	err := os.WriteFile(lockFile, []byte(pid), 0644)
	if err != nil {
		return fmt.Errorf("failed to create lock file: %w", err)
	}

	return nil
}

// releaseSingleInstance removes the lock file
func releaseSingleInstance() {
	if lockFile != "" {
		os.Remove(lockFile)
	}
}
