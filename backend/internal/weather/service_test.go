package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestServiceGetForecastsMapsAndSortsForecasts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"code": 0,
			"data": {
				"items": [{
					"valid_period": {"text": "6:00 PM to 8:00 PM"},
					"forecasts": [
						{"area": "Changi", "forecast": "Light Rain"},
						{"area": "Ang Mo Kio", "forecast": "Cloudy"},
						{"area": "Bedok", "forecast": "Partly Cloudy"}
					]
				}]
			}
		}`))
	}))
	defer server.Close()

	service := NewService(newWeatherTestClient(server))
	got, err := service.GetForecasts(context.Background())
	if err != nil {
		t.Fatalf("GetForecasts() error = %v", err)
	}

	want := &ForecastResponse{
		ValidPeriod: "6:00 PM to 8:00 PM",
		Forecasts: []LocationForecast{
			{Location: "Ang Mo Kio", Forecast: "Cloudy"},
			{Location: "Bedok", Forecast: "Partly Cloudy"},
			{Location: "Changi", Forecast: "Light Rain"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetForecasts() = %+v, want %+v", got, want)
	}
}

func TestServiceGetForecastsPropagatesClientError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	service := NewService(newWeatherTestClient(server))
	got, err := service.GetForecasts(context.Background())

	if err == nil {
		t.Fatalf("GetForecasts() error = nil, want error; response = %+v", got)
	}
	if got != nil {
		t.Errorf("GetForecasts() response = %+v, want nil", got)
	}
	if !strings.Contains(err.Error(), "weather API returned status 502") {
		t.Errorf("error = %q, want upstream status error", err.Error())
	}
}
