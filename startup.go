package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"golang.org/x/sys/windows/registry"
)

// EnableStartup enables the app to launch on system startup
func (a *App) EnableStartup() error {
	switch runtime.GOOS {
	case "windows":
		return a.enableStartupWindows()
	case "darwin":
		return a.enableStartupMacOS()
	case "linux":
		return a.enableStartupLinux()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// DisableStartup disables the app from launching on system startup
func (a *App) DisableStartup() error {
	switch runtime.GOOS {
	case "windows":
		return a.disableStartupWindows()
	case "darwin":
		return a.disableStartupMacOS()
	case "linux":
		return a.disableStartupLinux()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// IsStartupEnabled checks if startup is enabled
func (a *App) IsStartupEnabled() (bool, error) {
	switch runtime.GOOS {
	case "windows":
		return a.isStartupEnabledWindows()
	case "darwin":
		return a.isStartupEnabledMacOS()
	case "linux":
		return a.isStartupEnabledLinux()
	default:
		return false, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// Windows implementation
func (a *App) enableStartupWindows() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	
	return key.SetStringValue("myWeatherApp", exePath)
}

func (a *App) disableStartupWindows() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	
	return key.DeleteValue("myWeatherApp")
}

func (a *App) isStartupEnabledWindows() (bool, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false, nil
	}
	defer key.Close()
	
	_, _, err = key.GetStringValue("myWeatherApp")
	return err == nil, nil
}

// macOS implementation
func (a *App) enableStartupMacOS() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	launchAgentsDir := filepath.Join(homeDir, "Library/LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return err
	}
	
	plistPath := filepath.Join(launchAgentsDir, "com.myWeatherApp.plist")
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.myWeatherApp</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
</dict>
</plist>`, exePath)
	
	return os.WriteFile(plistPath, []byte(plistContent), 0644)
}

func (a *App) disableStartupMacOS() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	plistPath := filepath.Join(homeDir, "Library/LaunchAgents/com.myWeatherApp.plist")
	return os.Remove(plistPath)
}

func (a *App) isStartupEnabledMacOS() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	plistPath := filepath.Join(homeDir, "Library/LaunchAgents/com.myWeatherApp.plist")
	_, err = os.Stat(plistPath)
	return err == nil, nil
}

// Linux implementation
func (a *App) enableStartupLinux() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	autostartDir := filepath.Join(homeDir, ".config/autostart")
	if err := os.MkdirAll(autostartDir, 0755); err != nil {
		return err
	}
	
	desktopPath := filepath.Join(autostartDir, "myWeatherApp.desktop")
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	
	desktopContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=My Weather App 2
Exec=%s
Terminal=false
X-GNOME-Autostart-enabled=true`, exePath)
	
	return os.WriteFile(desktopPath, []byte(desktopContent), 0644)
}

func (a *App) disableStartupLinux() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	autostartPath := filepath.Join(homeDir, ".config/autostart/myWeatherApp.desktop")
	return os.Remove(autostartPath)
}

func (a *App) isStartupEnabledLinux() (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	autostartPath := filepath.Join(homeDir, ".config/autostart/myWeatherApp.desktop")
	_, err = os.Stat(autostartPath)
	return err == nil, nil
}
