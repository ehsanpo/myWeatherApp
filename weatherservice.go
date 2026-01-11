package main

import (
	"fmt"
	"math/rand"
	"time"
)

// WeatherService handles weather-related operations
type WeatherService struct {
	app *App
}

// WeatherData represents the weather information
type WeatherData struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feelsLike"`
	Condition   string  `json:"condition"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"windSpeed"`
	Icon        string  `json:"icon"`
	LastUpdated string  `json:"lastUpdated"`
	Forecast    []ForecastDay `json:"forecast"`
}

// ForecastDay represents a single day forecast
type ForecastDay struct {
	Date        string  `json:"date"`
	DayOfWeek   string  `json:"dayOfWeek"`
	MaxTemp     float64 `json:"maxTemp"`
	MinTemp     float64 `json:"minTemp"`
	Condition   string  `json:"condition"`
	Icon        string  `json:"icon"`
}

// NewWeatherService creates a new weather service instance
func NewWeatherService(app *App) *WeatherService {
	return &WeatherService{app: app}
}

// GetWeather fetches weather data for a given location (fake API for now)
func (w *WeatherService) GetWeather(location string) (*WeatherData, error) {
	// Simulate API call delay
	time.Sleep(500 * time.Millisecond)

	// If location is empty, use default from config
	if location == "" {
		config, err := w.app.LoadConfig()
		if err == nil && config.CustomSettings["weatherLocation"] != nil {
			location = config.CustomSettings["weatherLocation"].(string)
		} else {
			location = "New York"
		}
	}

	// Generate fake weather data
	conditions := []string{"Sunny", "Partly Cloudy", "Cloudy", "Rainy", "Stormy", "Snowy", "Foggy"}
	icons := []string{"100", "101", "104", "305", "302", "400", "500"}
	
	rand.Seed(time.Now().UnixNano())
	conditionIndex := rand.Intn(len(conditions))
	
	baseTemp := 15.0 + rand.Float64()*20.0 // 15-35°C
	
	weather := &WeatherData{
		Location:    location,
		Temperature: baseTemp,
		FeelsLike:   baseTemp + (rand.Float64()*4.0 - 2.0),
		Condition:   conditions[conditionIndex],
		Description: fmt.Sprintf("%s weather in %s", conditions[conditionIndex], location),
		Humidity:    40 + rand.Intn(50),
		WindSpeed:   5.0 + rand.Float64()*20.0,
		Icon:        icons[conditionIndex],
		LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
		Forecast:    generateFakeForecast(),
	}

	return weather, nil
}

// generateFakeForecast creates a 5-day forecast
func generateFakeForecast() []ForecastDay {
	forecast := make([]ForecastDay, 5)
	conditions := []string{"Sunny", "Partly Cloudy", "Cloudy", "Rainy", "Stormy"}
	icons := []string{"100", "101", "104", "305", "302"}
	
	for i := 0; i < 5; i++ {
		date := time.Now().AddDate(0, 0, i+1)
		conditionIndex := rand.Intn(len(conditions))
		baseTemp := 15.0 + rand.Float64()*20.0
		
		forecast[i] = ForecastDay{
			Date:      date.Format("2006-01-02"),
			DayOfWeek: date.Format("Monday"),
			MaxTemp:   baseTemp + 5.0,
			MinTemp:   baseTemp - 5.0,
			Condition: conditions[conditionIndex],
			Icon:      icons[conditionIndex],
		}
	}
	
	return forecast
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
	return w.app.SaveConfig(config)
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
