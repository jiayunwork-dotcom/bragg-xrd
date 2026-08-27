package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func postJSON(t *testing.T, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(data))
	rec := httptest.NewRecorder()
	Routes().ServeHTTP(rec, req)
	return rec
}

func TestBraggEndpoint(t *testing.T) {
	rec := postJSON(t, "/api/bragg", map[string]interface{}{"lambda": 1.5406, "d": 2.087, "n": 1})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
}

func TestPowderEndpoint(t *testing.T) {
	rec := postJSON(t, "/api/powder", map[string]interface{}{
		"lambda": 1.5406, "a": 3.615, "lattice": "fcc",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; body=%s", rec.Code, rec.Body.String())
	}
}

func TestInvalidBraggReturns400(t *testing.T) {
	rec := postJSON(t, "/api/bragg", map[string]interface{}{"lambda": 0, "d": 1, "n": 1})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestInvalidPowderReturns400(t *testing.T) {
	rec := postJSON(t, "/api/powder", map[string]interface{}{
		"lambda": 1.5406, "a": 3.615, "lattice": "hex",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/bragg", nil)
	rec := httptest.NewRecorder()
	Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}
