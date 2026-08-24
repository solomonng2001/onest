package weather

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func newWeatherTestMux(client *Client) *http.ServeMux {
	mux := http.NewServeMux()
	NewHandler(NewService(client)).RegisterRoutes(mux)
	return mux
}

func TestWeatherHandlerReturnsForecasts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"code": 0,
			"data": {
				"items": [{
					"valid_period": {"text": "6:00 PM to 8:00 PM"},
					"forecasts": [
						{"area": "Bedok", "forecast": "Cloudy"},
						{"area": "Ang Mo Kio", "forecast": "Light Rain"}
					]
				}]
			}
		}`))
	}))
	defer server.Close()

	request := httptest.NewRequest(http.MethodGet, "/api/weather", nil)
	recorder := httptest.NewRecorder()

	newWeatherTestMux(newWeatherTestClient(server)).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	var got ForecastResponse
	if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	want := ForecastResponse{
		ValidPeriod: "6:00 PM to 8:00 PM",
		Forecasts: []LocationForecast{
			{Location: "Ang Mo Kio", Forecast: "Light Rain"},
			{Location: "Bedok", Forecast: "Cloudy"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("response = %+v, want %+v", got, want)
	}
}

func TestWeatherHandlerReturnsBadGatewayWhenServiceFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	request := httptest.NewRequest(http.MethodGet, "/api/weather", nil)
	recorder := httptest.NewRecorder()

	newWeatherTestMux(newWeatherTestClient(server)).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusBadGateway, recorder.Body.String())
	}

	var got map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got["error"] != "weather service is currently unavailable" {
		t.Errorf("error = %q, want %q", got["error"], "weather service is currently unavailable")
	}
}

func TestWeatherRouteRejectsWrongMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/weather", nil)
	recorder := httptest.NewRecorder()

	newWeatherTestMux(newWeatherTestClient(server)).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
}
