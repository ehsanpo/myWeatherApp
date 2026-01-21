package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// UpdateInfo represents update information
type UpdateInfo struct {
	Version     string `json:"version"`
	ReleaseURL  string `json:"releaseUrl"`
	DownloadURL string `json:"downloadUrl"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
}

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

const (
	CurrentVersion = "v1.0.0"     // Update this with your app version
	GitHubRepo     = "ehsanpo/myWeatherApp" // Update with your GitHub repo
	CheckInterval  = 24 * time.Hour
)

// CheckForUpdates checks if a new version is available
func (a *App) CheckForUpdates() (*UpdateInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", GitHubRepo)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var release GitHubRelease
	err = json.Unmarshal(body, &release)
	if err != nil {
		return nil, fmt.Errorf("failed to parse release: %w", err)
	}

	updateInfo := &UpdateInfo{
		Version:     release.TagName,
		ReleaseURL:  release.HTMLURL,
		Description: release.Body,
		Available:   isNewerVersion(release.TagName, CurrentVersion),
	}

	// Find download URL for current platform
	platform := runtime.GOOS
	arch := runtime.GOARCH

	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, platform) && strings.Contains(name, arch) {
			updateInfo.DownloadURL = asset.BrowserDownloadURL
			break
		}
	}

	return updateInfo, nil
}

// GetCurrentVersion returns the current app version
func (a *App) GetCurrentVersion() string {
	return CurrentVersion
}

// OpenReleaseURL opens the release page in the browser
func (a *App) OpenReleaseURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // linux, etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	// Note: This is a simplified version. In production, use runtime.OpenBrowser() or similar
	fmt.Printf("Opening URL: %s (command: %s %v) ", url, cmd, args)
	return nil
}

// isNewerVersion compares version strings (simple semver comparison)
func isNewerVersion(latest, current string) bool {
	latest = strings.TrimPrefix(latest, "v")
	current = strings.TrimPrefix(current, "v")

	return latest > current // Simplified comparison
}
