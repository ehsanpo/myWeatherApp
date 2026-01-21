package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// WeatherService handles weather-related operations
type WeatherService struct {
	app            *App
	trayUpdateFunc func(*WeatherData)
}

// WeatherData represents the weather information
type WeatherData struct {
	Location    string        `json:"location"`
	Temperature float64       `json:"temperature"`
	FeelsLike   float64       `json:"feelsLike"`
	Condition   string        `json:"condition"`
	Description string        `json:"description"`
	Humidity    int           `json:"humidity"`
	WindSpeed   float64       `json:"windSpeed"`
	Icon        string        `json:"icon"`
	LastUpdated string        `json:"lastUpdated"`
	Forecast    []ForecastDay `json:"forecast"`
}

// ForecastDay represents a single day forecast
type ForecastDay struct {
	Date      string  `json:"date"`
	DayOfWeek string  `json:"dayOfWeek"`
	MaxTemp   float64 `json:"maxTemp"`
	MinTemp   float64 `json:"minTemp"`
	Condition string  `json:"condition"`
	Icon      string  `json:"icon"`
}

// NewWeatherService creates a new weather service instance
func NewWeatherService(app *App) *WeatherService {
	return &WeatherService{app: app}
}

// SetTrayUpdateFunc sets the function to update the tray icon
func (w *WeatherService) SetTrayUpdateFunc(updateFunc func(*WeatherData)) {
	w.trayUpdateFunc = updateFunc
}

// GeocodingResult represents geocoding API response
type GeocodingResult struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country   string  `json:"country"`
	} `json:"results"`
}

// OpenMeteoResponse represents the Open-Meteo API response
type OpenMeteoResponse struct {
	Current struct {
		Temperature      float64 `json:"temperature_2m"`
		RelativeHumidity int     `json:"relative_humidity_2m"`
		ApparentTemp     float64 `json:"apparent_temperature"`
		WindSpeed        float64 `json:"wind_speed_10m"`
		WeatherCode      int     `json:"weather_code"`
	} `json:"current"`
	Daily struct {
		Time        []string  `json:"time"`
		TempMax     []float64 `json:"temperature_2m_max"`
		TempMin     []float64 `json:"temperature_2m_min"`
		WeatherCode []int     `json:"weather_code"`
	} `json:"daily"`
}

// getCoordinates geocodes a location name to coordinates
func (w *WeatherService) getCoordinates(location string) (float64, float64, error) {
	baseURL := "https://geocoding-api.open-meteo.com/v1/search"
	params := url.Values{}
	params.Add("name", location)
	params.Add("count", "1")
	params.Add("language", "en")
	params.Add("format", "json")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	var result GeocodingResult
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, 0, err
	}

	if len(result.Results) == 0 {
		return 0, 0, fmt.Errorf("location not found")
	}

	return result.Results[0].Latitude, result.Results[0].Longitude, nil
}

// weatherCodeToCondition converts Open-Meteo weather code to condition string and icon
func weatherCodeToCondition(code int) (string, string) {
	switch code {
	case 0:
		return "Clear Sky", "100"
	case 1, 2, 3:
		return "Partly Cloudy", "101"
	case 45, 48:
		return "Foggy", "500"
	case 51, 53, 55:
		return "Drizzle", "300"
	case 61, 63, 65:
		return "Rainy", "305"
	case 66, 67:
		return "Freezing Rain", "313"
	case 71, 73, 75:
		return "Snowy", "400"
	case 77:
		return "Snow Grains", "400"
	case 80, 81, 82:
		return "Rain Showers", "309"
	case 85, 86:
		return "Snow Showers", "404"
	case 95:
		return "Thunderstorm", "302"
	case 96, 99:
		return "Thunderstorm with Hail", "302"
	default:
		return "Unknown", "999"
	}
}

// GetWeather fetches weather data for a given location from Open-Meteo API
func (w *WeatherService) GetWeather(location string) (*WeatherData, error) {
	// If location is empty, use default from config
	if location == "" {
		config, err := w.app.LoadConfig()
		if err == nil && config.CustomSettings["weatherLocation"] != nil {
			location = config.CustomSettings["weatherLocation"].(string)
		} else {
			location = "New York"
		}
	}

	// Get coordinates for location
	lat, lon, err := w.getCoordinates(location)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode location: %w", err)
	}

	// Fetch weather from Open-Meteo
	baseURL := "https://api.open-meteo.com/v1/forecast"
	params := url.Values{}
	params.Add("latitude", fmt.Sprintf("%.4f", lat))
	params.Add("longitude", fmt.Sprintf("%.4f", lon))
	params.Add("current", "temperature_2m,relative_humidity_2m,apparent_temperature,weather_code,wind_speed_10m")
	params.Add("daily", "weather_code,temperature_2m_max,temperature_2m_min")
	params.Add("timezone", "auto")
	params.Add("forecast_days", "6")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp OpenMeteoResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse weather data: %w", err)
	}

	// Convert to our weather data structure
	condition, icon := weatherCodeToCondition(apiResp.Current.WeatherCode)

	weather := &WeatherData{
		Location:    location,
		Temperature: apiResp.Current.Temperature,
		FeelsLike:   apiResp.Current.ApparentTemp,
		Condition:   condition,
		Description: fmt.Sprintf("%s in %s", condition, location),
		Humidity:    apiResp.Current.RelativeHumidity,
		WindSpeed:   apiResp.Current.WindSpeed,
		Icon:        icon,
		LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
		Forecast:    make([]ForecastDay, 0),
	}

	// Build forecast (skip today, get next 5 days)
	for i := 1; i < len(apiResp.Daily.Time) && i <= 5; i++ {
		date, _ := time.Parse("2006-01-02", apiResp.Daily.Time[i])
		condition, icon := weatherCodeToCondition(apiResp.Daily.WeatherCode[i])

		forecast := ForecastDay{
			Date:      apiResp.Daily.Time[i],
			DayOfWeek: date.Format("Monday"),
			MaxTemp:   apiResp.Daily.TempMax[i],
			MinTemp:   apiResp.Daily.TempMin[i],
			Condition: condition,
			Icon:      icon,
		}
		weather.Forecast = append(weather.Forecast, forecast)
	}

	return weather, nil
}

// UpdateLocation updates the weather location in config
func (w *WeatherService) UpdateLocation(location string) error {
	config, err := w.app.LoadConfig()
	if err != nil {
		return err
	}

	if config.CustomSettings == nil {
		config.CustomSettings = make(map[string]interface{})
	}

	config.CustomSettings["weatherLocation"] = location
	err = w.app.SaveConfig(config)
	if err != nil {
		return err
	}

	// Update tray icon with weather for the new location
	weather, err := w.GetWeather(location)
	if err != nil {
		return err
	}

	if w.trayUpdateFunc != nil {
		w.trayUpdateFunc(weather)
	}

	return nil
}

// RefreshWeather refreshes the weather data and updates tray icon
func (w *WeatherService) RefreshWeather(location string) (*WeatherData, error) {
	weather, err := w.GetWeather(location)
	if err != nil {
		return nil, err
	}

	// Update tray icon directly
	if w.trayUpdateFunc != nil {
		w.trayUpdateFunc(weather)
	}

	return weather, nil
}

// GetStoredLocation retrieves the stored weather location
func (w *WeatherService) GetStoredLocation() (string, error) {
	config, err := w.app.LoadConfig()
	if err != nil {
		return "New York", err
	}

	if config.CustomSettings["weatherLocation"] != nil {
		return config.CustomSettings["weatherLocation"].(string), nil
	}

	return "New York", nil
}
