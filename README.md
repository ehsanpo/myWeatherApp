# Weather Tray App

A cross-platform desktop weather application built with Wails v3 that shows weather information in a system tray icon.

## Features

- 🌤️ System tray icon with current weather display
- 🌡️ Real-time temperature and weather conditions
- 📍 Configurable location settings
- 🔄 Auto-refresh every 5 minutes
- 📊 5-day weather forecast
- 💨 Wind speed and humidity information
- 🎨 Clean, modern UI with gradient background
- ⚡ Built with React and Wails v3

## Project Structure

```
myWeatherApp/
├── main.go                 # Main application entry point
├── weatherservice.go       # Weather service with fake API
├── config.go              # Configuration management
├── greetservice.go        # (Legacy - can be removed)
├── frontend/
│   ├── src/
│   │   ├── App.jsx        # Main React component
│   │   ├── App.css        # Weather app styles
│   │   └── index.css      # Global styles
│   └── package.json
└── frontend-examples/
    ├── weather-helper.ts  # Weather service helper functions
    └── config-helper.ts   # Config helper functions
```

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

## Weather Service

The `WeatherService` currently uses a fake API that generates random weather data. To integrate a real API:

1. Sign up for a weather API service (e.g., OpenWeatherMap, WeatherAPI, QWeather)
2. Update the `GetWeather` method in `weatherservice.go`
3. Add API key to configuration
4. Parse real API responses into `WeatherData` struct

### WeatherData Structure

```go
type WeatherData struct {
    Location    string
    Temperature float64
    FeelsLike   float64
    Condition   string
    Description string
    Humidity    int
    WindSpeed   float64
    Icon        string
    LastUpdated string
    Forecast    []ForecastDay
}
```

## Development

### Prerequisites

- Go 1.25+
- Node.js 18+
- Wails v3 CLI

### Running the App

```bash
# Development mode with hot reload
wails3 dev

# Build for production
wails3 build
```

### Frontend Development

The frontend is built with:

- React 18
- Vite 5
- qweather-icons for weather icons

To install dependencies:

```bash
cd frontend
npm install
```

## System Tray Features

- **Label**: Shows current location and temperature
- **Menu Items**:
  - Show Weather - Opens the weather window
  - Refresh Weather - Manually updates weather data
  - Quit - Closes the application

## UI Components

### Main Window

- Location editor (click pencil icon)
- Current weather with large icon
- Temperature display (Celsius)
- Feels like temperature
- Humidity and wind speed
- 5-day forecast cards
- Refresh button

### Styling

- Gradient purple background
- Glass-morphism effects
- Responsive design
- Smooth animations

## Future Enhancements

- [ ] Dark/light theme toggle
- [ ] Hourly forecast
- [ ] Weather maps
- [ ] Historical data

## API Integration Example

To integrate OpenWeatherMap API, update `weatherservice.go`:

```go
func (w *WeatherService) GetWeather(location string) (*WeatherData, error) {
    apiKey := "YOUR_API_KEY"
    url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", location, apiKey)

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Parse response and populate WeatherData
    // ...
}
```

## License

This project is built with Wails v3 (https://wails.io)

## Credits

- Icons: qweather-icons package
- Framework: Wails v3
- Frontend: React + Vite
