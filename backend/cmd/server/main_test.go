package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	healthHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d; body = %s",
			recorder.Code,
			http.StatusOK,
			recorder.Body.String(),
		)
	}

	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	var got map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if got["status"] != "ok" {
		t.Errorf("status body = %q, want ok", got["status"])
	}
}

func TestCORSMiddlewareAddsHeadersAndForwardsRequest(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusAccepted)
	})

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set("Origin", defaultFrontendURL)
	recorder := httptest.NewRecorder()

	corsMiddleware(next).ServeHTTP(recorder, request)

	if !nextCalled {
		t.Fatal("next handler was not called")
	}
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusAccepted)
	}

	assertCORSHeaders(t, recorder)
}

func TestCORSMiddlewareHandlesPreflightRequest(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodOptions, "/api/uen/validate", nil)
	request.Header.Set("Origin", defaultFrontendURL)
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	recorder := httptest.NewRecorder()

	corsMiddleware(next).ServeHTTP(recorder, request)

	if nextCalled {
		t.Fatal("next handler was called for an OPTIONS preflight request")
	}
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}

	assertCORSHeaders(t, recorder)
}

func assertCORSHeaders(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	tests := []struct {
		header string
		want   string
	}{
		{
			header: "Access-Control-Allow-Origin",
			want:   defaultFrontendURL,
		},
		{
			header: "Access-Control-Allow-Methods",
			want:   "GET, POST, OPTIONS",
		},
		{
			header: "Access-Control-Allow-Headers",
			want:   "Content-Type",
		},
	}

	for _, tt := range tests {
		if got := recorder.Header().Get(tt.header); got != tt.want {
			t.Errorf("%s = %q, want %q", tt.header, got, tt.want)
		}
	}
}
