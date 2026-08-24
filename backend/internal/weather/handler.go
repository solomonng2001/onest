package weather

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /api/weather",
		h.getForecasts,
	)
}

func (h *Handler) getForecasts(
	w http.ResponseWriter,
	r *http.Request,
) {
	forecastResponse, err := h.service.GetForecasts(
		r.Context(),
	)
	if err != nil {
		slog.Error(
			"failed to retrieve weather forecasts",
			"error",
			err,
		)

		writeError(
			w,
			http.StatusBadGateway,
			"weather service is currently unavailable",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		forecastResponse,
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

func writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	writeJSON(
		w,
		status,
		map[string]string{
			"error": message,
		},
	)
}
