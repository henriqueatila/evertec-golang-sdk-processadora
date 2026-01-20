package otel

import (
	"context"
	"testing"
	"time"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/client"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func setupTestTracer() (*TracingHook, *tracetest.InMemoryExporter) {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
	)
	hook := NewTracingHookWithProvider(tp)
	return hook, exporter
}

func TestTracingHook_Success(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	reqInfo := &client.RequestInfo{
		Method:         "GET",
		Path:           "/accounts/123",
		BodySize:       0,
		IdempotencyKey: "",
	}
	respInfo := &client.ResponseInfo{
		StatusCode: 200,
		BodySize:   512,
		Duration:   50 * time.Millisecond,
		Error:      nil,
	}

	ctx = hook.BeforeRequest(ctx, reqInfo)
	hook.AfterResponse(ctx, reqInfo, respInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	span := spans[0]

	// Check span name
	if span.Name != "GET /accounts/123" {
		t.Errorf("expected span name 'GET /accounts/123', got %s", span.Name)
	}

	// Check status
	if span.Status.Code != codes.Ok {
		t.Errorf("expected Ok status, got %v", span.Status.Code)
	}

	// Check attributes
	attrs := attributesToMap(span.Attributes)

	if attrs["http.method"] != "GET" {
		t.Errorf("expected http.method GET, got %v", attrs["http.method"])
	}
	if attrs["http.url"] != "/accounts/123" {
		t.Errorf("expected http.url /accounts/123, got %v", attrs["http.url"])
	}
	if attrs["http.status_code"] != int64(200) {
		t.Errorf("expected http.status_code 200, got %v", attrs["http.status_code"])
	}
}

func TestTracingHook_WithIdempotencyKey(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	reqInfo := &client.RequestInfo{
		Method:         "POST",
		Path:           "/accounts",
		BodySize:       256,
		IdempotencyKey: "test-key-123",
	}
	respInfo := &client.ResponseInfo{
		StatusCode: 201,
		BodySize:   128,
		Duration:   100 * time.Millisecond,
	}

	ctx = hook.BeforeRequest(ctx, reqInfo)
	hook.AfterResponse(ctx, reqInfo, respInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	attrs := attributesToMap(spans[0].Attributes)

	if attrs["paysmart.idempotency_key"] != "test-key-123" {
		t.Errorf("expected idempotency_key test-key-123, got %v", attrs["paysmart.idempotency_key"])
	}
}

func TestTracingHook_Error(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	reqInfo := &client.RequestInfo{
		Method: "POST",
		Path:   "/accounts",
	}
	respInfo := &client.ResponseInfo{
		StatusCode: 400,
		Duration:   30 * time.Millisecond,
		Error: &client.APIError{
			StatusCode: 400,
			Code:       "INVALID_REQUEST",
			Message:    "Invalid document",
		},
	}

	ctx = hook.BeforeRequest(ctx, reqInfo)
	hook.AfterResponse(ctx, reqInfo, respInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	span := spans[0]

	// Check error status
	if span.Status.Code != codes.Error {
		t.Errorf("expected Error status, got %v", span.Status.Code)
	}

	// Check error attributes
	attrs := attributesToMap(span.Attributes)
	if attrs["paysmart.error_code"] != "INVALID_REQUEST" {
		t.Errorf("expected error_code INVALID_REQUEST, got %v", attrs["paysmart.error_code"])
	}

	// Check recorded error events
	if len(span.Events) == 0 {
		t.Error("expected error event to be recorded")
	}
}

func TestTracingHook_HTTPError(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	reqInfo := &client.RequestInfo{Method: "GET", Path: "/test"}
	respInfo := &client.ResponseInfo{
		StatusCode: 500,
		Duration:   time.Millisecond,
		Error:      nil, // No error object, just HTTP error
	}

	ctx = hook.BeforeRequest(ctx, reqInfo)
	hook.AfterResponse(ctx, reqInfo, respInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	if spans[0].Status.Code != codes.Error {
		t.Errorf("expected Error status for 500, got %v", spans[0].Status.Code)
	}
}

func TestTracingHook_NoSpanInContext(t *testing.T) {
	hook, _ := setupTestTracer()

	// Call AfterResponse without BeforeRequest (no span in context)
	ctx := context.Background()
	reqInfo := &client.RequestInfo{Method: "GET", Path: "/test"}
	respInfo := &client.ResponseInfo{StatusCode: 200}

	// Should not panic
	hook.AfterResponse(ctx, reqInfo, respInfo)
}

func TestNewTracingHook(t *testing.T) {
	hook := NewTracingHook()
	if hook == nil {
		t.Error("NewTracingHook returned nil")
	}
	if hook.tracer == nil {
		t.Error("tracer is nil")
	}
}

// Helper function to convert attributes to map
func attributesToMap(attrs []attribute.KeyValue) map[string]interface{} {
	result := make(map[string]interface{})
	for _, attr := range attrs {
		result[string(attr.Key)] = attr.Value.AsInterface()
	}
	return result
}
