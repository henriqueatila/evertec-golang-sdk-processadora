package client

import (
	"context"
	"sync"
	"time"
)

// MetricsCollector is an interface for collecting request metrics.
// Implement this interface to integrate with your preferred metrics system
// (e.g., Prometheus, DataDog, CloudWatch).
type MetricsCollector interface {
	// RecordRequest is called after each API request with request details.
	RecordRequest(method, path string, statusCode int, duration time.Duration, err error)
}

// MetricsHook implements Hook for collecting metrics.
type MetricsHook struct {
	collector MetricsCollector
}

// NewMetricsHook creates a new metrics hook with the given collector.
//
// Example with a custom Prometheus collector:
//
//	type prometheusCollector struct {
//	    requestCounter  *prometheus.CounterVec
//	    requestDuration *prometheus.HistogramVec
//	}
//
//	func (p *prometheusCollector) RecordRequest(method, path string, statusCode int, duration time.Duration, err error) {
//	    p.requestCounter.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
//	    p.requestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
//	}
//
//	client, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
//	    client.WithHooks(client.NewMetricsHook(&prometheusCollector{...})),
//	)
func NewMetricsHook(collector MetricsCollector) *MetricsHook {
	return &MetricsHook{collector: collector}
}

func (h *MetricsHook) BeforeRequest(ctx context.Context, req *RequestInfo) context.Context {
	return ctx
}

func (h *MetricsHook) AfterResponse(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
	h.collector.RecordRequest(req.Method, req.Path, resp.StatusCode, resp.Duration, resp.Error)
}

// SimpleMetrics is a basic in-memory metrics collector for testing and development.
// For production, use a proper metrics system like Prometheus.
type SimpleMetrics struct {
	mu       sync.RWMutex
	Requests int64
	Errors   int64
	TotalMs  int64
	ByMethod map[string]int64
	ByPath   map[string]int64
	ByStatus map[int]int64
}

// NewSimpleMetrics creates a new SimpleMetrics collector.
func NewSimpleMetrics() *SimpleMetrics {
	return &SimpleMetrics{
		ByMethod: make(map[string]int64),
		ByPath:   make(map[string]int64),
		ByStatus: make(map[int]int64),
	}
}

// RecordRequest implements MetricsCollector.
func (m *SimpleMetrics) RecordRequest(method, path string, statusCode int, duration time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Requests++
	m.TotalMs += duration.Milliseconds()
	m.ByMethod[method]++
	m.ByPath[path]++
	m.ByStatus[statusCode]++

	if err != nil || statusCode >= 400 {
		m.Errors++
	}
}

// Stats returns aggregated metrics.
func (m *SimpleMetrics) Stats() (requests, errors, avgMs int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.Requests > 0 {
		avgMs = m.TotalMs / m.Requests
	}
	return m.Requests, m.Errors, avgMs
}

// Reset clears all collected metrics.
func (m *SimpleMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Requests = 0
	m.Errors = 0
	m.TotalMs = 0
	m.ByMethod = make(map[string]int64)
	m.ByPath = make(map[string]int64)
	m.ByStatus = make(map[int]int64)
}
