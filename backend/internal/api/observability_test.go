package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsMiddlewarePreservesFlusher(t *testing.T) {
	metrics := NewMetrics()
	handler := metrics.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("metrics middleware must preserve http.Flusher")
		}
		flusher.Flush()
		w.WriteHeader(http.StatusNoContent)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/events/stream", nil))
	if metrics.requests.Load() != 1 {
		t.Fatalf("requests = %d, want 1", metrics.requests.Load())
	}
}

func TestRequestIDMiddlewareAddsID(t *testing.T) {
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := w.Header().Get("X-Request-ID"); got == "" {
			t.Fatal("request ID header missing")
		}
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Header().Get("X-Request-ID") == "" {
		t.Fatal("response request ID missing")
	}
}
