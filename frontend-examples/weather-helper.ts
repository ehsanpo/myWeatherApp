// Weather Helper for myWeatherApp
import {
  GetWeather,
  UpdateLocation,
  GetStoredLocation,
} from "../bindings/changeme/weatherservice";

/**
 * Fetch current weather for a location
 * @param location - Location to get weather for (empty string uses stored location)
 * @returns Promise with weather data
 */
export async function fetchWeather(location = "") {
  try {
    const weather = await GetWeather(location);
    return weather;
  } catch (error) {
    console.error("Failed to fetch weather:", error);
    throw error;
  }
}

/**
 * Update the stored location
 * @param location - New location to save
 */
export async function saveLocation(location: string) {
  try {
    await UpdateLocation(location);
    console.log("Location updated:", location);
    return true;
  } catch (error) {
    console.error("Failed to update location:", error);
    return false;
  }
}

/**
 * Get the stored location from config
 */
export async function getLocation() {
  try {
    const location = await GetStoredLocation();
    return location;
  } catch (error) {
    console.error("Failed to get location:", error);
    return "New York";
  }
}

/**
 * Example usage demonstrating weather app functionality
 */
export async function exampleUsage() {
  // Get stored location
  const location = await getLocation();
  console.log("Stored location:", location);

  // Fetch weather for stored location
  const weather = await fetchWeather();
  console.log("Current weather:", weather);

  // Update location
  await saveLocation("London");

  // Fetch weather for new location
  const londonWeather = await fetchWeather("London");
  console.log("London weather:", londonWeather);
}
