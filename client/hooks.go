package client

import (
	"context"
	"time"
)

// RequestInfo contains information about the outgoing request.
type RequestInfo struct {
	// Method is the HTTP method (GET, POST, etc.)
	Method string
	// Path is the API endpoint path
	Path string
	// BodySize is the request body size in bytes
	BodySize int
	// IdempotencyKey is the idempotency key for mutations
	IdempotencyKey string
}

// ResponseInfo contains information about the response.
type ResponseInfo struct {
	// StatusCode is the HTTP status code
	StatusCode int
	// BodySize is the response body size in bytes
	BodySize int
	// Duration is the total request duration
	Duration time.Duration
	// Error is the error if the request failed (nil on success)
	Error error
}

// Hook is an interface for observability hooks.
// Implement this interface to add custom logging, metrics, or tracing.
type Hook interface {
	// BeforeRequest is called before each request is sent.
	// The returned context is passed to AfterResponse and can be used
	// to store request-scoped data (e.g., trace spans).
	BeforeRequest(ctx context.Context, req *RequestInfo) context.Context

	// AfterResponse is called after each response is received (or on error).
	// This is always called, even if the request failed.
	AfterResponse(ctx context.Context, req *RequestInfo, resp *ResponseInfo)
}

// BeforeRequestFunc is a function type for BeforeRequest hooks.
type BeforeRequestFunc func(ctx context.Context, req *RequestInfo) context.Context

// AfterResponseFunc is a function type for AfterResponse hooks.
type AfterResponseFunc func(ctx context.Context, req *RequestInfo, resp *ResponseInfo)

// funcHook wraps functions as a Hook implementation.
type funcHook struct {
	before BeforeRequestFunc
	after  AfterResponseFunc
}

func (f *funcHook) BeforeRequest(ctx context.Context, req *RequestInfo) context.Context {
	if f.before != nil {
		return f.before(ctx, req)
	}
	return ctx
}

func (f *funcHook) AfterResponse(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
	if f.after != nil {
		f.after(ctx, req, resp)
	}
}

// NewHook creates a Hook from before and after functions.
// Either function can be nil if not needed.
func NewHook(before BeforeRequestFunc, after AfterResponseFunc) Hook {
	return &funcHook{before: before, after: after}
}

// NewAfterHook creates a Hook that only runs after the request.
// Useful for simple metrics or logging hooks.
func NewAfterHook(after AfterResponseFunc) Hook {
	return &funcHook{after: after}
}
