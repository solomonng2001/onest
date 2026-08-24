import { useEffect, useState } from "react";
import "./App.css";

type Tab = "weather" | "uen";

type Forecast = {
  location: string;
  forecast: string;
};

type WeatherResponse = {
  validPeriod: string;
  forecasts: Forecast[];
};

const WEATHER_API_URL = "http://localhost:8080/api/weather";

function App() {
  const [activeTab, setActiveTab] = useState<Tab>("weather");
  const [weather, setWeather] = useState<WeatherResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function loadWeather() {
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
    <div className="app-shell">
      <header className="app-header">
        <div className="container">
          <h1>OneST Web Portal</h1>
        </div>
      </header>

      <main className="container">
        <nav className="tabs" aria-label="OneST services">
          <button
            type="button"
            role="tab"
            aria-selected={activeTab === "weather"}
            className={
              activeTab === "weather" ? "tab tab-active" : "tab"
            }
            onClick={() => setActiveTab("weather")}
          >
            Weather forecast
          </button>

          <button
            type="button"
            role="tab"
            aria-selected={activeTab === "uen"}
            className={activeTab === "uen" ? "tab tab-active" : "tab"}
            onClick={() => setActiveTab("uen")}
          >
            UEN validation
          </button>
        </nav>

        {activeTab === "weather" ? (
          <section className="page" role="tabpanel">
            <div className="page-heading">
              <h2>Two-hour weather forecast</h2>

              {weather && (
                <p>
                  Valid period: <strong>{weather.validPeriod}</strong>
                </p>
              )}
            </div>

            {isLoading && (
              <p className="status">Loading forecast...</p>
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
        ) : (
          <section className="page" role="tabpanel">
            <div className="page-heading">
              <h2>UEN validation</h2>
              <p>This service will be added next.</p>
            </div>
          </section>
        )}
      </main>
    </div>
  );
}

export default App;