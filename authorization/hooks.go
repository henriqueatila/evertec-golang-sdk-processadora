package authorization

import (
	"context"
	"time"
)

// RequestInfo contains information about the incoming authorization request.
type RequestInfo struct {
	// Path is the request path (e.g., /purchases, /withdrawals)
	Path string
	// OperationType is the type of operation (e.g., "purchase", "withdrawal", "query")
	OperationType string
	// TransactionID is the transaction identifier from Evertec
	TransactionID string
	// AccountID is the account identifier if available
	AccountID string
	// CardID is the card identifier if available
	CardID string
	// Amount is the transaction amount in cents (if applicable)
	Amount int64
	// Currency is the currency code (if applicable)
	Currency string
	// BodySize is the request body size in bytes
	BodySize int
}

// ResponseInfo contains information about the authorization response.
type ResponseInfo struct {
	// StatusCode is the HTTP status code returned
	StatusCode int
	// ResponseCode is the Evertec response code (0=approved, others=declined)
	ResponseCode int
	// Duration is the total processing duration
	Duration time.Duration
	// Error is the error if processing failed (nil on success)
	Error error
	// Approved indicates if the authorization was approved
	Approved bool
}

// Hook is an interface for observability hooks.
// Implement this interface to add custom logging, metrics, or tracing
// to your authorization server.
type Hook interface {
	// BeforeAuthorization is called before processing each authorization request.
	// The returned context is passed to AfterAuthorization and can be used
	// to store request-scoped data (e.g., trace spans).
	BeforeAuthorization(ctx context.Context, req *RequestInfo) context.Context

	// AfterAuthorization is called after each authorization response is sent.
	// This is always called, even if processing failed.
	AfterAuthorization(ctx context.Context, req *RequestInfo, resp *ResponseInfo)
}

// BeforeAuthorizationFunc is a function type for BeforeAuthorization hooks.
type BeforeAuthorizationFunc func(ctx context.Context, req *RequestInfo) context.Context

// AfterAuthorizationFunc is a function type for AfterAuthorization hooks.
type AfterAuthorizationFunc func(ctx context.Context, req *RequestInfo, resp *ResponseInfo)

// funcHook wraps functions as a Hook implementation.
type funcHook struct {
	before BeforeAuthorizationFunc
	after  AfterAuthorizationFunc
}

func (f *funcHook) BeforeAuthorization(ctx context.Context, req *RequestInfo) context.Context {
	if f.before != nil {
		return f.before(ctx, req)
	}
	return ctx
}

func (f *funcHook) AfterAuthorization(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
	if f.after != nil {
		f.after(ctx, req, resp)
	}
}

// NewHook creates a Hook from before and after functions.
// Either function can be nil if not needed.
func NewHook(before BeforeAuthorizationFunc, after AfterAuthorizationFunc) Hook {
	return &funcHook{before: before, after: after}
}

// NewAfterHook creates a Hook that only runs after the authorization.
// Useful for simple metrics or logging hooks.
func NewAfterHook(after AfterAuthorizationFunc) Hook {
	return &funcHook{after: after}
}
