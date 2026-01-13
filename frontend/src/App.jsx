import { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [weather, setWeather] = useState(null);
  const [location, setLocation] = useState('');
  const [editingLocation, setEditingLocation] = useState(false);
  const [newLocation, setNewLocation] = useState('');
  const [loading, setLoading] = useState(true);
  const [weatherIcons, setWeatherIcons] = useState({});

  // Load SVG icons
  const loadIcon = async (iconCode) => {
    try {
      const response = await fetch(
        new URL(`../node_modules/qweather-icons/icons/${iconCode}.svg`, import.meta.url).href
      );
      const svgText = await response.text();
      return svgText;
    } catch {
      return null;
    }
  };

  const handleMinimize = async () => {
    try {
      const { HideWindow } = await import('../bindings/changeme/app');
      await HideWindow();
    } catch (error) {
      console.error('Failed to hide window:', error);
    }
  };

  // Import weather service methods dynamically
  const loadWeather = async () => {
    try {
      setLoading(true);
      const { GetWeather } = await import('../bindings/changeme/weatherservice');
      const data = await GetWeather(location);
      setWeather(data);

      // Load icons for current weather and forecast
      const icons = {};
      icons[data.icon] = await loadIcon(data.icon);
      for (const day of data.forecast) {
        if (!icons[day.icon]) {
          icons[day.icon] = await loadIcon(day.icon);
        }
      }
      setWeatherIcons(icons);

      setLoading(false);
    } catch (error) {
      console.error('Failed to load weather:', error);
      setLoading(false);
    }
  };

  const handleUpdateLocation = async () => {
    if (!newLocation.trim()) return;

    try {
      const { UpdateLocation } = await import('../bindings/changeme/weatherservice');
      await UpdateLocation(newLocation);
      setLocation(newLocation);
      setEditingLocation(false);
      loadWeather();
    } catch (error) {
      console.error('Failed to update location:', error);
    }
  };

  const getStoredLocation = async () => {
    try {
      const { GetStoredLocation } = await import('../bindings/changeme/weatherservice');
      const stored = await GetStoredLocation();
      setLocation(stored);
    } catch (error) {
      console.error('Failed to get stored location:', error);
      setLocation('New York');
    }
  };

  useEffect(() => {
    getStoredLocation();
  }, []);

  useEffect(() => {
    if (location) {
      loadWeather();
      // Refresh weather every 5 minutes
      const interval = setInterval(loadWeather, 300000);
      return () => clearInterval(interval);
    }
  }, [location]);

  if (loading) {
    return (
      <div className="weather-app loading">
        <div className="spinner"></div>
        <p>Loading weather...</p>
      </div>
    );
  }

  if (!weather) {
    return (
      <div className="weather-app error">
        <p>Unable to load weather data</p>
      </div>
    );
  }

  return (
    <div className="weather-app">
      <div className="window-controls">
        <button className="minimize-btn" onClick={handleMinimize} title="Hide to tray">
          ✕
        </button>
      </div>
      <div className="weather-header">
        <div className="location-section">
          {!editingLocation ? (
            <>
              <h2>{weather.location}</h2>
              <button
                className="edit-btn"
                onClick={() => {
                  setEditingLocation(true);
                  setNewLocation(weather.location);
                }}
              >
                ✏️
              </button>
            </>
          ) : (
            <div className="location-edit">
              <input
                type="text"
                value={newLocation}
                onChange={(e) => setNewLocation(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleUpdateLocation()}
                autoFocus
              />
              <button onClick={handleUpdateLocation}>✓</button>
              <button onClick={() => setEditingLocation(false)}>✗</button>
            </div>
          )}
        </div>
        <div className="last-updated">
          Updated: {new Date(weather.lastUpdated).toLocaleTimeString()}
        </div>
      </div>

      <div className="current-weather">
        <div
          className="weather-icon-large"
          dangerouslySetInnerHTML={{ __html: weatherIcons[weather.icon] || '' }}
        />
        <div className="temperature">
          <span className="temp-value">{Math.round(weather.temperature)}°</span>
          <span className="temp-unit">C</span>
        </div>
        <div className="condition">{weather.condition}</div>
        <div className="feels-like">Feels like {Math.round(weather.feelsLike)}°C</div>
      </div>

      <div className="weather-details">
        <div className="detail-item">
          <span className="detail-label">Humidity</span>
          <span className="detail-value">{weather.humidity}%</span>
        </div>
        <div className="detail-item">
          <span className="detail-label">Wind Speed</span>
          <span className="detail-value">{weather.windSpeed.toFixed(1)} km/h</span>
        </div>
      </div>

      <div className="forecast">
        <h3>5-Day Forecast</h3>
        <div className="forecast-list">
          {weather.forecast.map((day, index) => (
            <div key={index} className="forecast-item">
              <div className="forecast-day">{day.dayOfWeek.substring(0, 3)}</div>
              <div
                className="weather-icon-small"
                dangerouslySetInnerHTML={{ __html: weatherIcons[day.icon] || '' }}
              />
              <div className="forecast-temps">
                <span className="temp-max">{Math.round(day.maxTemp)}°</span>
                <span className="temp-min">{Math.round(day.minTemp)}°</span>
              </div>
            </div>
          ))}
        </div>
      </div>

      <button className="refresh-btn" onClick={loadWeather}>
        🔄 Refresh
      </button>
    </div>
  );
}

export default App;
