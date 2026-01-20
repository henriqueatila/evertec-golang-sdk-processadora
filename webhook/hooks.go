package webhook

import (
	"context"
	"time"
)

// EventInfo contains information about the incoming webhook event.
type EventInfo struct {
	// EventType is the type of event (e.g., DISPUTE_STATUS_CHANGE)
	EventType EventType
	// EventID is the event identifier from Evertec (from X-PaySmart-Event-Id header)
	EventID string
	// BodySize is the request body size in bytes
	BodySize int
	// IsDuplicate indicates if this event was already processed (idempotency)
	IsDuplicate bool
}

// ProcessingInfo contains information about the event processing result.
type ProcessingInfo struct {
	// StatusCode is the HTTP status code returned
	StatusCode int
	// Duration is the total processing duration
	Duration time.Duration
	// Error is the error if processing failed (nil on success)
	Error error
	// Processed indicates if the event was successfully processed
	Processed bool
}

// Hook is an interface for observability hooks.
// Implement this interface to add custom logging, metrics, or tracing
// to your webhook server.
type Hook interface {
	// BeforeEvent is called before processing each webhook event.
	// The returned context is passed to AfterEvent and can be used
	// to store request-scoped data (e.g., trace spans).
	BeforeEvent(ctx context.Context, event *EventInfo) context.Context

	// AfterEvent is called after each event is processed.
	// This is always called, even if processing failed.
	AfterEvent(ctx context.Context, event *EventInfo, result *ProcessingInfo)
}

// BeforeEventFunc is a function type for BeforeEvent hooks.
type BeforeEventFunc func(ctx context.Context, event *EventInfo) context.Context

// AfterEventFunc is a function type for AfterEvent hooks.
type AfterEventFunc func(ctx context.Context, event *EventInfo, result *ProcessingInfo)

// funcHook wraps functions as a Hook implementation.
type funcHook struct {
	before BeforeEventFunc
	after  AfterEventFunc
}

func (f *funcHook) BeforeEvent(ctx context.Context, event *EventInfo) context.Context {
	if f.before != nil {
		return f.before(ctx, event)
	}
	return ctx
}

func (f *funcHook) AfterEvent(ctx context.Context, event *EventInfo, result *ProcessingInfo) {
	if f.after != nil {
		f.after(ctx, event, result)
	}
}

// NewHook creates a Hook from before and after functions.
// Either function can be nil if not needed.
func NewHook(before BeforeEventFunc, after AfterEventFunc) Hook {
	return &funcHook{before: before, after: after}
}

// NewAfterHook creates a Hook that only runs after event processing.
// Useful for simple metrics or logging hooks.
func NewAfterHook(after AfterEventFunc) Hook {
	return &funcHook{after: after}
}
