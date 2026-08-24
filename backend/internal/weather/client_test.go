package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newWeatherTestClient(server *httptest.Server) *Client {
	return &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}
}

func TestClientGetForecastSuccess(t *testing.T) {
	type requestDetails struct {
		method string
		accept string
	}
	requestReceived := make(chan requestDetails, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived <- requestDetails{
			method: r.Method,
			accept: r.Header.Get("Accept"),
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"code": 0,
			"data": {
				"items": [{
					"valid_period": {"text": "6:00 PM to 8:00 PM"},
					"forecasts": [
						{"area": "Bedok", "forecast": "Cloudy"}
					]
				}]
			}
		}`))
	}))
	defer server.Close()

	got, err := newWeatherTestClient(server).GetForecast(context.Background())
	if err != nil {
		t.Fatalf("GetForecast() error = %v", err)
	}

	details := <-requestReceived
	if details.method != http.MethodGet {
		t.Errorf("request method = %q, want GET", details.method)
	}
	if details.accept != "application/json" {
		t.Errorf("Accept header = %q, want application/json", details.accept)
	}
	if got.Code != 0 {
		t.Errorf("Code = %d, want 0", got.Code)
	}
	if len(got.Data.Items) != 1 {
		t.Fatalf("items length = %d, want 1", len(got.Data.Items))
	}
	if got.Data.Items[0].ValidPeriod.Text != "6:00 PM to 8:00 PM" {
		t.Errorf("valid period = %q, want %q", got.Data.Items[0].ValidPeriod.Text, "6:00 PM to 8:00 PM")
	}
	if len(got.Data.Items[0].Forecasts) != 1 {
		t.Fatalf("forecasts length = %d, want 1", len(got.Data.Items[0].Forecasts))
	}
	if got.Data.Items[0].Forecasts[0].Area != "Bedok" {
		t.Errorf("area = %q, want Bedok", got.Data.Items[0].Forecasts[0].Area)
	}
}

func TestClientGetForecastFailures(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantError  string
	}{
		{
			name:       "non-200 HTTP response",
			statusCode: http.StatusServiceUnavailable,
			body:       `{"error":"unavailable"}`,
			wantError:  "weather API returned status 503",
		},
		{
			name:       "malformed JSON",
			statusCode: http.StatusOK,
			body:       `{"code":`,
			wantError:  "decode weather response",
		},
		{
			name:       "API error with message",
			statusCode: http.StatusOK,
			body:       `{"code":42,"errorMsg":"invalid request"}`,
			wantError:  "weather API error: invalid request",
		},
		{
			name:       "API error without message",
			statusCode: http.StatusOK,
			body:       `{"code":42}`,
			wantError:  "weather API returned an unsuccessful response",
		},
		{
			name:       "no forecast items",
			statusCode: http.StatusOK,
			body:       `{"code":0,"data":{"items":[]}}`,
			wantError:  "weather API returned no forecast items",
		},
		{
			name:       "no location forecasts",
			statusCode: http.StatusOK,
			body:       `{"code":0,"data":{"items":[{"valid_period":{"text":"6:00 PM to 8:00 PM"},"forecasts":[]}]}}`,
			wantError:  "weather API returned no location forecasts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			got, err := newWeatherTestClient(server).GetForecast(context.Background())
			if err == nil {
				t.Fatalf("GetForecast() error = nil, want error containing %q; response = %+v", tt.wantError, got)
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Errorf("error = %q, want it to contain %q", err.Error(), tt.wantError)
			}
		})
	}
}

func TestClientGetForecastRejectsInvalidBaseURL(t *testing.T) {
	client := &Client{
		baseURL:    "://invalid-url",
		httpClient: http.DefaultClient,
	}

	got, err := client.GetForecast(context.Background())
	if err == nil {
		t.Fatalf("GetForecast() error = nil, want request-creation error; response = %+v", got)
	}
	if !strings.Contains(err.Error(), "create weather request") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "create weather request")
	}
}
