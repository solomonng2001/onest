import { useEffect, useState, type FormEvent } from "react";
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

type UENResponse = {
  uen: string;
  valid: boolean;
  format?: string;
  message: string;
};

const WEATHER_API_URL = "http://localhost:8080/api/weather";
const UEN_API_URL = "http://localhost:8080/api/uen/validate";

function App() {
  const [activeTab, setActiveTab] = useState<Tab>("weather");

  const [weather, setWeather] = useState<WeatherResponse | null>(null);
  const [isWeatherLoading, setIsWeatherLoading] = useState(true);
  const [weatherError, setWeatherError] = useState<string | null>(null);

  const [uen, setUen] = useState("");
  const [uenResult, setUenResult] = useState<UENResponse | null>(null);
  const [isUENLoading, setIsUENLoading] = useState(false);
  const [uenError, setUenError] = useState<string | null>(null);

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
        setWeatherError(
          requestError instanceof Error
            ? requestError.message
            : "Unable to load the weather forecast.",
        );
      } finally {
        setIsWeatherLoading(false);
      }
    }

    void loadWeather();
  }, []);

  async function validateUEN(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setUenResult(null);
    setUenError(null);
    setIsUENLoading(true);

    try {
      const response = await fetch(UEN_API_URL, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          uen: uen.trim(),
        }),
      });

      if (!response.ok) {
        throw new Error(
          `UEN request failed with status ${response.status}`,
        );
      }

      const result: UENResponse = await response.json();
      setUenResult(result);
    } catch (requestError) {
      setUenError(
        requestError instanceof Error
          ? requestError.message
          : "Unable to validate the UEN.",
      );
    } finally {
      setIsUENLoading(false);
    }
  }

  return (
    <div className="app-shell">
      <header className="app-header">
        <div className="container">
          <h1>OneST Web Portal</h1>
        </div>
      </header>

      <main className="container main-content">
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
              <p>View the latest forecast across Singapore.</p>

              {weather && (
                <p className="valid-period">
                  Valid period: <strong>{weather.validPeriod}</strong>
                </p>
              )}
            </div>

            {isWeatherLoading && (
              <p className="status">Loading forecast...</p>
            )}

            {weatherError && (
              <p className="status status-error" role="alert">
                {weatherError}
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
              <p>Check whether a Unique Entity Number has a valid format.</p>
            </div>

            <div className="uen-card">
              <form className="uen-form" onSubmit={validateUEN}>
                <label htmlFor="uen">Unique Entity Number</label>

                <div className="uen-input-row">
                  <input
                    id="uen"
                    type="text"
                    value={uen}
                    onChange={(event) => {
                      setUen(event.target.value.toUpperCase());
                      setUenResult(null);
                      setUenError(null);
                    }}
                    placeholder="For example, 200912345N"
                    autoComplete="off"
                    required
                  />

                  <button type="submit" disabled={isUENLoading}>
                    {isUENLoading ? "Validating..." : "Validate UEN"}
                  </button>
                </div>

                <p className="input-hint">
                  Enter the UEN without spaces or symbols.
                </p>
              </form>

              {uenResult && (
                <div
                  className={
                    uenResult.valid
                      ? "uen-result uen-result-valid"
                      : "uen-result uen-result-invalid"
                  }
                  role="status"
                >
                  <p className="result-title">
                    {uenResult.valid ? "Valid UEN format" : "Invalid UEN format"}
                  </p>

                  <p>{uenResult.message}</p>

                  {uenResult.format && (
                    <p className="result-format">
                      Format:{" "}
                      <strong>
                        {uenResult.format.replaceAll("_", " ")}
                      </strong>
                    </p>
                  )}
                </div>
              )}

              {uenError && (
                <div className="uen-result uen-result-invalid" role="alert">
                  <p className="result-title">Unable to validate UEN</p>
                  <p>{uenError}</p>
                </div>
              )}
            </div>
          </section>
        )}
      </main>
    </div>
  );
}

export default App;