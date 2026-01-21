# Weather Tray App

A cross-platform desktop weather application built with Wails v3 that displays real-time weather information in a system tray icon.

## Features

- 🌤️ System tray icon with current weather display
- 🌡️ Real-time temperature and weather conditions
- 📍 Configurable location settings
- 🔄 Auto-refresh every 5 minutes
- 📊 5-day weather forecast
- 💨 Wind speed and humidity information
- 🎨 Clean, modern UI with gradient background
- ⚡ Built with React and Wails v3

## Installation

### Prerequisites

- Go 1.25+
- Node.js 18+
- Wails v3 CLI

### Building from Source

```bash
# Clone the repository
git clone https://github.com/ehsanpo/myWeatherApp.git
cd myWeatherApp

# Install frontend dependencies
cd frontend
npm install
cd ..

# Build the application
wails3 build
```

The built application will be available in the `build/bin` directory.

## Usage

1. Run the application executable
2. The app will appear in your system tray
3. Click the tray icon to open the weather window
4. Use the tray menu to show/hide the window or quit the app

### System Tray Features

- **Label**: Shows current location and temperature
- **Menu Items**:
  - Show Weather - Opens the weather window
  - Refresh Weather - Manually updates weather data
  - Quit - Closes the application

## Configuration

The app stores configuration in `~/.myWeatherApp/config.json`:

```json
{
  "theme": "light",
  "language": "en",
  "windowWidth": 400,
  "windowHeight": 600,
  "customSettings": {
    "weatherLocation": "New York",
    "updateInterval": 300,
    "temperatureUnit": "celsius"
  }
}
```

## Project Structure

```
myWeatherApp/
├── main.go                 # Main application entry point
├── weatherservice.go       # Weather service
├── config.go               # Configuration management
├── frontend/
│   ├── src/
│   │   ├── App.jsx        # Main React component
│   │   ├── App.css        # Weather app styles
│   │   └── index.css      # Global styles
│   └── package.json
└── build/                  # Build output directory
```

## Weather Service

The application uses a weather service that provides current conditions and forecasts. The service is designed to be easily extensible for integration with real weather APIs.

## License

This project is built with Wails v3 (https://wails.io)

## Credits

- Icons: qweather-icons package
- Framework: Wails v3
- Frontend: React + Vite
