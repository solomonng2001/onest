package main

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/YOUR_GITHUB_USERNAME/onest/backend/internal/weather"
)

const (
	defaultPort        = "8080"
	defaultFrontendURL = "http://localhost:5173"
)

func main() {
	mux := http.NewServeMux()

	weatherClient := weather.NewClient()
	weatherService := weather.NewService(weatherClient)
	weatherHandler := weather.NewHandler(weatherService)

	weatherHandler.RegisterRoutes(mux)

	mux.HandleFunc(
		"GET /health",
		healthHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           corsMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	slog.Info(
		"starting server",
		"port",
		port,
	)

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func healthHandler(
	w http.ResponseWriter,
	_ *http.Request,
) {
	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "ok",
		},
	)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	value any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Error(
			"failed to encode JSON response",
			"error",
			err,
		)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(
				"Access-Control-Allow-Origin",
				defaultFrontendURL,
			)

			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, OPTIONS",
			)

			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type",
			)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}
