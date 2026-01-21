package main

import (
	//"os"
	"embed"
	_ "embed"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Register custom events
func init() {
	application.RegisterEvent[*WeatherData]("trayIconUpdate")
}

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var trayIcon []byte

// App struct to hold application state and provide utility methods
type App struct {
	mainWindow *application.WebviewWindow
}

// HideWindow hides the main window
func (a *App) HideWindow() {
	if a.mainWindow != nil {
		a.mainWindow.Hide()
	}
}

// MinimizeWindow minimizes the main window
func (a *App) MinimizeWindow() {
	if a.mainWindow != nil {
		a.mainWindow.Minimise()
	}
}

// PositionWindowNearTray positions the window near the system tray
func (a *App) PositionWindowNearTray() {
	if a.mainWindow == nil {
		return
	}

	// Get primary screen dimensions
	screen, err := a.mainWindow.GetScreen()
	if err != nil || screen == nil {
		log.Printf("Failed to get screen: %v", err)
		return
	}

	windowWidth := 400
	windowHeight := 450
	padding := 10

	var x, y int

	if runtime.GOOS == "windows" {
		// Windows: Position at bottom-right above taskbar
		x = screen.Size.Width - windowWidth - padding
		y = screen.Size.Height - windowHeight - padding - 40 // Extra space for taskbar
	} else if runtime.GOOS == "darwin" {
		// macOS: Position at top-right below menu bar
		x = screen.Size.Width - windowWidth - padding
		y = padding + 25 // Space for menu bar
	} else {
		// Linux/other: Position at top-right
		x = screen.Size.Width - windowWidth - padding
		y = padding
	}

	a.mainWindow.SetPosition(x, y)
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	// Create app instance for methods
	appInstance := &App{}
	weatherService := NewWeatherService(appInstance)

	app := application.New(application.Options{
		Name:        "myWeatherApp",
		Description: "A weather app with system tray",
		Services: []application.Service{
			application.NewService(weatherService),
			application.NewService(appInstance),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	// Create system tray
	systray := app.SystemTray.New()

	// Create a function to update tray icon that can be called from weather service
	updateTrayIconFunc := func(weather *WeatherData) {
		log.Printf("Updating tray icon: Location=%s, Temperature=%.2f°C, Condition=%s",
			weather.Location, weather.Temperature, weather.Condition)

		iconData, err := generateTrayIconWithWeather(weather)
		if err == nil {
			systray.SetIcon(iconData)
		} else {
			log.Printf("Failed to generate tray icon: %v", err)
		}
		systray.SetLabel(fmt.Sprintf("%s: %.0f°C - %s", weather.Location, weather.Temperature, weather.Condition))
	}

	// Pass the update function to the weather service
	weatherService.SetTrayUpdateFunc(updateTrayIconFunc)

	// Function to update tray icon with current weather
	updateTrayIcon := func() {
		weather, err := weatherService.GetWeather("")
		if err != nil {
			log.Printf("Failed to get weather: %v", err)
			return
		}

		log.Printf("Weather data: Location=%s, Temperature=%.2f°C, Condition=%s",
			weather.Location, weather.Temperature, weather.Condition)

		updateTrayIconFunc(weather)
	}

	// Set initial icon
	updateTrayIcon()

	// Update tray icon periodically
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			updateTrayIcon()
		}
	}()

	// Create a new window with the necessary options.
	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Weather App",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		Width:            400,
		Height:           450,
		MinWidth:         400,
		MinHeight:        450,
		MaxWidth:         400,
		MaxHeight:        450,
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		URL:              "/",
		Hidden:           true, // Start hidden
		Frameless:        true,
	})
	// Store window reference in app instance
	appInstance.mainWindow = mainWindow

	// Show and position window after a delay
	go func() {
		time.Sleep(300 * time.Millisecond)
		log.Println("Positioning and showing window...")
		appInstance.PositionWindowNearTray()
		mainWindow.Show()
		log.Println("Window shown and positioned")
	}()

	// Add system tray menu
	menu := app.NewMenu()
	menu.Add("Show Weather").OnClick(func(ctx *application.Context) {
		mainWindow.Show()
		mainWindow.UnMinimise()
		mainWindow.Focus()
		// Position after showing the window
		appInstance.PositionWindowNearTray()
	})
	menu.AddSeparator()
	menu.Add("Refresh Weather").OnClick(func(ctx *application.Context) {
		updateTrayIcon()
	})
	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	systray.SetMenu(menu)

	// Run the application. This blocks until the application has been exited.
	// Initialize single instance lock
	//if err := initSingleInstance(); err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//defer releaseSingleInstance()

	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
