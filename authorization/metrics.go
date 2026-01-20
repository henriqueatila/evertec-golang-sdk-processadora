package authorization

import (
	"context"
	"sync"
	"time"
)

// MetricsCollector is an interface for collecting authorization metrics.
// Implement this interface to integrate with your preferred metrics system
// (e.g., Prometheus, DataDog, CloudWatch).
type MetricsCollector interface {
	// RecordAuthorization is called after each authorization with details.
	RecordAuthorization(operationType string, responseCode int, approved bool, duration time.Duration, err error)
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
//	    authCounter    *prometheus.CounterVec
//	    authDuration   *prometheus.HistogramVec
//	    approvalRate   *prometheus.GaugeVec
//	}
//
//	func (p *prometheusCollector) RecordAuthorization(opType string, code int, approved bool, duration time.Duration, err error) {
//	    status := "approved"
//	    if !approved {
//	        status = "declined"
//	    }
//	    p.authCounter.WithLabelValues(opType, status, strconv.Itoa(code)).Inc()
//	    p.authDuration.WithLabelValues(opType).Observe(duration.Seconds())
//	}
//
//	server := authorization.NewServer(handler,
//	    authorization.WithHooks(authorization.NewMetricsHook(&prometheusCollector{...})),
//	)
func NewMetricsHook(collector MetricsCollector) *MetricsHook {
	return &MetricsHook{collector: collector}
}

func (h *MetricsHook) BeforeAuthorization(ctx context.Context, req *RequestInfo) context.Context {
	return ctx
}

func (h *MetricsHook) AfterAuthorization(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
	h.collector.RecordAuthorization(req.OperationType, resp.ResponseCode, resp.Approved, resp.Duration, resp.Error)
}

// SimpleMetrics is a basic in-memory metrics collector for testing and development.
// For production, use a proper metrics system like Prometheus.
type SimpleMetrics struct {
	mu             sync.RWMutex
	Authorizations int64
	Approved       int64
	Declined       int64
	Errors         int64
	TotalMs        int64
	ByOperation    map[string]*OperationMetrics
	ByResponseCode map[int]int64
}

// OperationMetrics holds metrics for a specific operation type.
type OperationMetrics struct {
	Total    int64
	Approved int64
	Declined int64
	TotalMs  int64
}

// NewSimpleMetrics creates a new SimpleMetrics collector.
func NewSimpleMetrics() *SimpleMetrics {
	return &SimpleMetrics{
		ByOperation:    make(map[string]*OperationMetrics),
		ByResponseCode: make(map[int]int64),
	}
}

// RecordAuthorization implements MetricsCollector.
func (m *SimpleMetrics) RecordAuthorization(operationType string, responseCode int, approved bool, duration time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Authorizations++
	m.TotalMs += duration.Milliseconds()
	m.ByResponseCode[responseCode]++

	if approved {
		m.Approved++
	} else {
		m.Declined++
	}

	if err != nil {
		m.Errors++
	}

	// Track per-operation metrics
	opMetrics, ok := m.ByOperation[operationType]
	if !ok {
		opMetrics = &OperationMetrics{}
		m.ByOperation[operationType] = opMetrics
	}
	opMetrics.Total++
	opMetrics.TotalMs += duration.Milliseconds()
	if approved {
		opMetrics.Approved++
	} else {
		opMetrics.Declined++
	}
}

// Stats returns aggregated metrics.
func (m *SimpleMetrics) Stats() (total, approved, declined, errors, avgMs int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.Authorizations > 0 {
		avgMs = m.TotalMs / m.Authorizations
	}
	return m.Authorizations, m.Approved, m.Declined, m.Errors, avgMs
}

// ApprovalRate returns the approval rate as a percentage (0-100).
func (m *SimpleMetrics) ApprovalRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.Authorizations == 0 {
		return 0
	}
	return float64(m.Approved) / float64(m.Authorizations) * 100
}

// Reset clears all collected metrics.
func (m *SimpleMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Authorizations = 0
	m.Approved = 0
	m.Declined = 0
	m.Errors = 0
	m.TotalMs = 0
	m.ByOperation = make(map[string]*OperationMetrics)
	m.ByResponseCode = make(map[int]int64)
}
