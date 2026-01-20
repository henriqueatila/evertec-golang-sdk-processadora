package webhook

import (
	"context"
	"sync"
	"time"
)

// MetricsCollector is an interface for collecting webhook metrics.
// Implement this interface to integrate with your preferred metrics system
// (e.g., Prometheus, DataDog, CloudWatch).
type MetricsCollector interface {
	// RecordEvent is called after each webhook event with details.
	RecordEvent(eventType EventType, processed bool, duplicate bool, duration time.Duration, err error)
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
//	    eventCounter   *prometheus.CounterVec
//	    eventDuration  *prometheus.HistogramVec
//	}
//
//	func (p *prometheusCollector) RecordEvent(eventType webhook.EventType, processed, duplicate bool, duration time.Duration, err error) {
//	    status := "success"
//	    if err != nil {
//	        status = "error"
//	    } else if duplicate {
//	        status = "duplicate"
//	    }
//	    p.eventCounter.WithLabelValues(string(eventType), status).Inc()
//	    p.eventDuration.WithLabelValues(string(eventType)).Observe(duration.Seconds())
//	}
//
//	server := webhook.NewServer(webhook.Config{
//	    Handler: handler,
//	    Hooks:   []webhook.Hook{webhook.NewMetricsHook(&prometheusCollector{...})},
//	})
func NewMetricsHook(collector MetricsCollector) *MetricsHook {
	return &MetricsHook{collector: collector}
}

func (h *MetricsHook) BeforeEvent(ctx context.Context, event *EventInfo) context.Context {
	return ctx
}

func (h *MetricsHook) AfterEvent(ctx context.Context, event *EventInfo, result *ProcessingInfo) {
	h.collector.RecordEvent(event.EventType, result.Processed, event.IsDuplicate, result.Duration, result.Error)
}

// SimpleMetrics is a basic in-memory metrics collector for testing and development.
// For production, use a proper metrics system like Prometheus.
type SimpleMetrics struct {
	mu          sync.RWMutex
	Events      int64
	Processed   int64
	Duplicates  int64
	Errors      int64
	TotalMs     int64
	ByEventType map[EventType]*EventTypeMetrics
}

// EventTypeMetrics holds metrics for a specific event type.
type EventTypeMetrics struct {
	Total      int64
	Processed  int64
	Duplicates int64
	Errors     int64
	TotalMs    int64
}

// NewSimpleMetrics creates a new SimpleMetrics collector.
func NewSimpleMetrics() *SimpleMetrics {
	return &SimpleMetrics{
		ByEventType: make(map[EventType]*EventTypeMetrics),
	}
}

// RecordEvent implements MetricsCollector.
func (m *SimpleMetrics) RecordEvent(eventType EventType, processed bool, duplicate bool, duration time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Events++
	m.TotalMs += duration.Milliseconds()

	if processed {
		m.Processed++
	}

	if duplicate {
		m.Duplicates++
	}

	if err != nil {
		m.Errors++
	}

	// Track per-event-type metrics
	etMetrics, ok := m.ByEventType[eventType]
	if !ok {
		etMetrics = &EventTypeMetrics{}
		m.ByEventType[eventType] = etMetrics
	}
	etMetrics.Total++
	etMetrics.TotalMs += duration.Milliseconds()
	if processed {
		etMetrics.Processed++
	}
	if duplicate {
		etMetrics.Duplicates++
	}
	if err != nil {
		etMetrics.Errors++
	}
}

// Stats returns aggregated metrics.
func (m *SimpleMetrics) Stats() (total, processed, duplicates, errors, avgMs int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.Events > 0 {
		avgMs = m.TotalMs / m.Events
	}
	return m.Events, m.Processed, m.Duplicates, m.Errors, avgMs
}

// Reset clears all collected metrics.
func (m *SimpleMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Events = 0
	m.Processed = 0
	m.Duplicates = 0
	m.Errors = 0
	m.TotalMs = 0
	m.ByEventType = make(map[EventType]*EventTypeMetrics)
}
