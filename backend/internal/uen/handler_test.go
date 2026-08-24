package uen

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newUENTestMux() *http.ServeMux {
	mux := http.NewServeMux()
	NewHandler(NewService()).RegisterRoutes(mux)
	return mux
}

func TestValidateHandlerReturnsValidationResult(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/uen/validate",
		strings.NewReader(`{"uen":"  t26ll1234c  "}`),
	)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	newUENTestMux().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	var got ValidateResponse
	if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	want := ValidateResponse{
		UEN:     "T26LL1234C",
		Valid:   true,
		Format:  FormatOtherEntity,
		Message: validFormatMessage,
	}
	if got != want {
		t.Fatalf("response = %+v, want %+v", got, want)
	}
}

func TestValidateHandlerRejectsInvalidRequestBodies(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		wantError string
	}{
		{
			name:      "empty body",
			body:      "",
			wantError: "request body must be valid JSON containing a uen field",
		},
		{
			name:      "malformed JSON",
			body:      `{"uen":`,
			wantError: "request body must be valid JSON containing a uen field",
		},
		{
			name:      "wrong UEN value type",
			body:      `{"uen":12345678}`,
			wantError: "request body must be valid JSON containing a uen field",
		},
		{
			name:      "unknown field",
			body:      `{"uen":"12345678A","extra":true}`,
			wantError: "request body must be valid JSON containing a uen field",
		},
		{
			name:      "multiple JSON objects",
			body:      `{"uen":"12345678A"} {"uen":"202612345B"}`,
			wantError: "request body must contain exactly one JSON object",
		},
		{
			name:      "trailing non-JSON content",
			body:      `{"uen":"12345678A"} trailing`,
			wantError: "request body must contain exactly one JSON object",
		},
		{
			name:      "body exceeds one MiB limit",
			body:      `{"uen":"` + strings.Repeat("A", (1<<20)+1) + `"}`,
			wantError: "request body must be valid JSON containing a uen field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodPost,
				"/api/uen/validate",
				strings.NewReader(tt.body),
			)
			recorder := httptest.NewRecorder()

			newUENTestMux().ServeHTTP(recorder, request)

			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusBadRequest, recorder.Body.String())
			}

			var got map[string]string
			if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if got["error"] != tt.wantError {
				t.Errorf("error = %q, want %q", got["error"], tt.wantError)
			}
		})
	}
}

func TestValidateHandlerTreatsMissingUENAsRequired(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/uen/validate",
		strings.NewReader(`{}`),
	)
	recorder := httptest.NewRecorder()

	newUENTestMux().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var got ValidateResponse
	if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	want := ValidateResponse{
		UEN:     "",
		Valid:   false,
		Message: "UEN is required",
	}
	if got != want {
		t.Fatalf("response = %+v, want %+v", got, want)
	}
}

func TestValidateRouteRejectsWrongMethod(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/uen/validate", nil)
	recorder := httptest.NewRecorder()

	newUENTestMux().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
}
