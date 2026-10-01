package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type Metrics struct {
	requests   atomic.Uint64
	errors     atomic.Uint64
	inFlight   atomic.Int64
	totalNanos atomic.Uint64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.inFlight.Add(1)
		defer m.inFlight.Add(-1)

		rw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		m.requests.Add(1)
		m.totalNanos.Add(uint64(time.Since(start).Nanoseconds()))
		if rw.status >= http.StatusInternalServerError {
			m.errors.Add(1)
		}
	})
}

func (m *Metrics) Handler(w http.ResponseWriter, _ *http.Request) {
	requests := m.requests.Load()
	average := uint64(0)
	if requests > 0 {
		average = m.totalNanos.Load() / requests
	}
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "# HELP worldbuster_http_requests_total Total HTTP requests handled.
")
	fmt.Fprintf(w, "# TYPE worldbuster_http_requests_total counter
")
	fmt.Fprintf(w, "worldbuster_http_requests_total %d
", requests)
	fmt.Fprintf(w, "# HELP worldbuster_http_errors_total Total HTTP 5xx responses.
")
	fmt.Fprintf(w, "# TYPE worldbuster_http_errors_total counter
")
	fmt.Fprintf(w, "worldbuster_http_errors_total %d
", m.errors.Load())
	fmt.Fprintf(w, "# HELP worldbuster_http_in_flight Current HTTP requests in flight.
")
	fmt.Fprintf(w, "# TYPE worldbuster_http_in_flight gauge
")
	fmt.Fprintf(w, "worldbuster_http_in_flight %d
", m.inFlight.Load())
	fmt.Fprintf(w, "# HELP worldbuster_http_average_latency_nanoseconds Average HTTP latency.
")
	fmt.Fprintf(w, "# TYPE worldbuster_http_average_latency_nanoseconds gauge
")
	fmt.Fprintf(w, "worldbuster_http_average_latency_nanoseconds %d
", average)
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(body []byte) (int, error) {
	return w.ResponseWriter.Write(body)
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = newRequestID()
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}

var requestSequence atomic.Uint64

func newRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) + "-" + strconv.FormatUint(requestSequence.Add(1), 36)
}

type ReadinessCheck func(context.Context) error
