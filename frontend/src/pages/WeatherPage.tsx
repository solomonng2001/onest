import { useEffect, useState } from "react";
import "./WeatherPage.css";

type Forecast = {
  location: string;
  forecast: string;
};

type WeatherResponse = {
  validPeriod: string;
  forecasts: Forecast[];
};

const WEATHER_API_URL = "http://localhost:8080/api/weather";

function WeatherPage() {
  const [weather, setWeather] = useState<WeatherResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function loadWeather() {
      setIsLoading(true);
      setError(null);

      try {
        const response = await fetch(WEATHER_API_URL);

        if (!response.ok) {
          throw new Error(
            `Weather request failed with status ${response.status}`,
          );
        }

        const result: WeatherResponse = await response.json();

        if (result.forecasts.length === 0) {
          throw new Error("No weather forecast is available.");
        }

        setWeather(result);
      } catch (requestError) {
        setError(
          requestError instanceof Error
            ? requestError.message
            : "Unable to load the weather forecast.",
        );
      } finally {
        setIsLoading(false);
      }
    }

    void loadWeather();
  }, []);

  return (
    <section
      id="weather-panel"
      className="page"
      role="tabpanel"
      aria-labelledby="weather-tab"
    >
      <div className="page-heading">
        <h2>Two-hour weather forecast</h2>
        <p>View the latest forecast across Singapore.</p>

        {weather && (
          <p className="valid-period">
            Valid period: <strong>{weather.validPeriod}</strong>
          </p>
        )}
      </div>

      {isLoading && (
        <p className="status" role="status">
          Loading forecast...
        </p>
      )}

      {error && (
        <p className="status status-error" role="alert">
          {error}
        </p>
      )}

      {weather && (
        <div className="table-wrapper">
          <table>
            <thead>
              <tr>
                <th scope="col">Location</th>
                <th scope="col">Forecast</th>
              </tr>
            </thead>

            <tbody>
              {weather.forecasts.map((item) => (
                <tr key={item.location}>
                  <td>{item.location}</td>
                  <td>{item.forecast}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}

export default WeatherPage;
