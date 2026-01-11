# System Tray

## Overview

System tray support has been added to your application.


## Wails v3 Usage

The system tray is fully initialized in `main.go` with icon and menu:

```go
// At the top of main.go
//go:embed build/appicon.png
var trayIcon []byte

// In main() function
systray := app.SystemTray.New()
systray.SetIcon(trayIcon)
systray.SetLabel("myWeatherApp")

// Add system tray menu
menu := app.NewMenu()
menu.Add("Show Window").OnClick(func(ctx *application.Context) {
	windows := app.Window.GetAll()
	if len(windows) > 0 {
		windows[0].Show()
		windows[0].UnMinimise()
	}
})
menu.AddSeparator()
menu.Add("Quit").OnClick(func(ctx *application.Context) {
	app.Quit()
})
systray.SetMenu(menu)
```

### Adding a Custom Icon

Replace `build/appicon.png` with your own icon file. Supported formats:
- Windows: .ico, .png
- macOS: .png (will be used as template image)
- Linux: .png

### Adding Menu Items

Add more menu items before setting the menu:

```go
menu.Add("Settings").OnClick(func(ctx *application.Context) {
	// Open settings window
})
```


## Features

- Show/Hide window from tray menu
- Quit application
- Custom icon support
- Menu with separators
- Direct menu creation in main.go

## Customization


Edit the menu creation code in `main.go` to:
- Add more menu items
- Change menu item labels
- Add custom click handlers
- Create submenus


## Platform Support

- Windows: Supports icons and menus
- macOS: Supports icons and menus
- Linux: Support varies by desktop environment
