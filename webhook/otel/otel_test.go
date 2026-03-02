package otel

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/webhook"
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

func TestTracingHook_DisputeStatusChange(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	eventInfo := &webhook.EventInfo{
		EventType:   webhook.EventTypeDisputeStatusChange,
		EventID:     "evt-123",
		BodySize:    512,
		IsDuplicate: false,
	}
	processingInfo := &webhook.ProcessingInfo{
		StatusCode: 200,
		Duration:   50 * time.Millisecond,
		Processed:  true,
	}

	ctx = hook.BeforeEvent(ctx, eventInfo)
	hook.AfterEvent(ctx, eventInfo, processingInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	span := spans[0]

	// Check span name
	expectedName := "evertec.webhook.DISPUTE_STATUS_CHANGE"
	if span.Name != expectedName {
		t.Errorf("expected span name %q, got %q", expectedName, span.Name)
	}

	// Check status is OK
	if span.Status.Code != codes.Ok {
		t.Errorf("expected Ok status, got %v", span.Status.Code)
	}

	// Check attributes
	attrs := attributesToMap(span.Attributes)

	if attrs["evertec.event_type"] != "DISPUTE_STATUS_CHANGE" {
		t.Errorf("expected event_type 'DISPUTE_STATUS_CHANGE', got %v", attrs["evertec.event_type"])
	}
	if attrs["evertec.event_id"] != "evt-123" {
		t.Errorf("expected event_id 'evt-123', got %v", attrs["evertec.event_id"])
	}
	if attrs["http.request_content_length"] != int64(512) {
		t.Errorf("expected http.request_content_length 512, got %v", attrs["http.request_content_length"])
	}
	if attrs["http.status_code"] != int64(200) {
		t.Errorf("expected http.status_code 200, got %v", attrs["http.status_code"])
	}
	if attrs["evertec.processed"] != true {
		t.Errorf("expected processed true, got %v", attrs["evertec.processed"])
	}
}

func TestTracingHook_StatementClosed(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	eventInfo := &webhook.EventInfo{
		EventType:   webhook.EventTypeStatementClosed,
		EventID:     "evt-stmt-456",
		BodySize:    300,
		IsDuplicate: false,
	}
	processingInfo := &webhook.ProcessingInfo{
		StatusCode: 200,
		Duration:   30 * time.Millisecond,
		Processed:  true,
	}

	ctx = hook.BeforeEvent(ctx, eventInfo)
	hook.AfterEvent(ctx, eventInfo, processingInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	if spans[0].Name != "evertec.webhook.STATEMENT_CLOSED_NOTIFICATION" {
		t.Errorf("expected span name 'evertec.webhook.STATEMENT_CLOSED_NOTIFICATION', got %q", spans[0].Name)
	}
}

func TestTracingHook_PaymentDue(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	eventInfo := &webhook.EventInfo{
		EventType:   webhook.EventTypePaymentDue,
		EventID:     "evt-pay-789",
		BodySize:    250,
		IsDuplicate: false,
	}
	processingInfo := &webhook.ProcessingInfo{
		StatusCode: 200,
		Duration:   25 * time.Millisecond,
		Processed:  true,
	}

	ctx = hook.BeforeEvent(ctx, eventInfo)
	hook.AfterEvent(ctx, eventInfo, processingInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	if spans[0].Name != "evertec.webhook.PAYMENT_DUE_NOTIFICATION" {
		t.Errorf("expected span name 'evertec.webhook.PAYMENT_DUE_NOTIFICATION', got %q", spans[0].Name)
	}
}

func TestTracingHook_CardStatusChange(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	eventInfo := &webhook.EventInfo{
		EventType:   webhook.EventTypeCardStatusChange,
		EventID:     "evt-card-101",
		BodySize:    180,
		IsDuplicate: false,
	}
	processingInfo := &webhook.ProcessingInfo{
		StatusCode: 200,
		Duration:   20 * time.Millisecond,
		Processed:  true,
	}

	ctx = hook.BeforeEvent(ctx, eventInfo)
	hook.AfterEvent(ctx, eventInfo, processingInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	if spans[0].Name != "evertec.webhook.CARD_STATUS_CHANGE" {
		t.Errorf("expected span name 'evertec.webhook.CARD_STATUS_CHANGE', got %q", spans[0].Name)
	}
}

func TestTracingHook_DeviceTokenStatus(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	eventInfo := &webhook.EventInfo{
		EventType:   webhook.EventTypeDeviceTokenStatus,
		EventID:     "evt-token-202",
		BodySize:    400,
		IsDuplicate: false,
	}
	processingInfo := &webhook.ProcessingInfo{
		StatusCode: 200,
		Duration:   35 * time.Millisecond,
		Processed:  true,
	}

	ctx = hook.BeforeEvent(ctx, eventInfo)
	hook.AfterEvent(ctx, eventInfo, processingInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	if spans[0].Name != "evertec.webhook.DEVICE_TOKEN_STATUS_NOTIFICATION" {
		t.Errorf("expected span name 'evertec.webhook.DEVICE_TOKEN_STATUS_NOTIFICATION', got %q", spans[0].Name)
	}
}

func TestTracingHook_DuplicateEvent(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	eventInfo := &webhook.EventInfo{
		EventType:   webhook.EventTypeDisputeStatusChange,
		EventID:     "evt-dup-123",
		BodySize:    256,
		IsDuplicate: true,
	}
	processingInfo := &webhook.ProcessingInfo{
		StatusCode: 200,
		Duration:   5 * time.Millisecond,
		Processed:  false, // Not processed because duplicate
	}

	ctx = hook.BeforeEvent(ctx, eventInfo)
	hook.AfterEvent(ctx, eventInfo, processingInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	span := spans[0]

	// Should be OK status
	if span.Status.Code != codes.Ok {
		t.Errorf("expected Ok status for skipped, got %v", span.Status.Code)
	}

	attrs := attributesToMap(span.Attributes)
	if attrs["evertec.duplicate"] != true {
		t.Errorf("expected duplicate true, got %v", attrs["evertec.duplicate"])
	}
	if attrs["evertec.processed"] != false {
		t.Errorf("expected processed false, got %v", attrs["evertec.processed"])
	}
}

func TestTracingHook_WithError(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	eventInfo := &webhook.EventInfo{
		EventType:   webhook.EventTypeDisputeStatusChange,
		EventID:     "evt-err-123",
		BodySize:    128,
		IsDuplicate: false,
	}
	testErr := errors.New("handler error: failed to process dispute")
	processingInfo := &webhook.ProcessingInfo{
		StatusCode: 500,
		Duration:   10 * time.Millisecond,
		Error:      testErr,
		Processed:  false,
	}

	ctx = hook.BeforeEvent(ctx, eventInfo)
	hook.AfterEvent(ctx, eventInfo, processingInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	span := spans[0]

	// Error should set Error status
	if span.Status.Code != codes.Error {
		t.Errorf("expected Error status, got %v", span.Status.Code)
	}

	// Check error event was recorded
	if len(span.Events) == 0 {
		t.Error("expected error event to be recorded")
	}
}

func TestTracingHook_EmptyEventType(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	eventInfo := &webhook.EventInfo{
		EventType:   "", // Empty event type
		EventID:     "evt-empty",
		BodySize:    100,
		IsDuplicate: false,
	}
	processingInfo := &webhook.ProcessingInfo{
		StatusCode: 200,
		Duration:   15 * time.Millisecond,
		Processed:  true,
	}

	ctx = hook.BeforeEvent(ctx, eventInfo)
	hook.AfterEvent(ctx, eventInfo, processingInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	// Should use default span name without event type
	if spans[0].Name != "evertec.webhook" {
		t.Errorf("expected span name 'evertec.webhook', got %q", spans[0].Name)
	}
}

func TestTracingHook_NoEventID(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	eventInfo := &webhook.EventInfo{
		EventType:   webhook.EventTypeDisputeStatusChange,
		EventID:     "", // Empty event ID
		BodySize:    100,
		IsDuplicate: false,
	}
	processingInfo := &webhook.ProcessingInfo{
		StatusCode: 200,
		Duration:   15 * time.Millisecond,
		Processed:  true,
	}

	ctx = hook.BeforeEvent(ctx, eventInfo)
	hook.AfterEvent(ctx, eventInfo, processingInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	attrs := attributesToMap(spans[0].Attributes)

	// event_id should not be set for empty string
	if _, exists := attrs["evertec.event_id"]; exists {
		t.Error("expected event_id to not be set for empty string")
	}
}

func TestTracingHook_NoSpanInContext(t *testing.T) {
	hook, _ := setupTestTracer()

	// Call AfterEvent without BeforeEvent (no span in context)
	ctx := context.Background()
	eventInfo := &webhook.EventInfo{EventType: webhook.EventTypeDisputeStatusChange}
	processingInfo := &webhook.ProcessingInfo{StatusCode: 200}

	// Should not panic
	hook.AfterEvent(ctx, eventInfo, processingInfo)
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

func TestNewTracingHookWithProvider(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
	)

	hook := NewTracingHookWithProvider(tp)
	if hook == nil {
		t.Error("NewTracingHookWithProvider returned nil")
	}
	if hook.tracer == nil {
		t.Error("tracer is nil")
	}
}

func TestTracingHook_ImplementsInterface(t *testing.T) {
	hook := NewTracingHook()

	// Verify hook implements webhook.Hook interface
	var _ webhook.Hook = hook
}

func TestTracingHook_SpanAttributes(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	eventInfo := &webhook.EventInfo{
		EventType:   webhook.EventTypeDisputeStatusChange,
		EventID:     "evt-attr-test",
		BodySize:    200,
		IsDuplicate: false,
	}
	processingInfo := &webhook.ProcessingInfo{
		StatusCode: 200,
		Duration:   25 * time.Millisecond,
		Processed:  true,
	}

	ctx = hook.BeforeEvent(ctx, eventInfo)
	hook.AfterEvent(ctx, eventInfo, processingInfo)

	spans := exporter.GetSpans()
	attrs := attributesToMap(spans[0].Attributes)

	// Check standard attributes
	if attrs["rpc.system"] != "evertec" {
		t.Errorf("expected rpc.system 'evertec', got %v", attrs["rpc.system"])
	}
	if attrs["rpc.service"] != "webhook" {
		t.Errorf("expected rpc.service 'webhook', got %v", attrs["rpc.service"])
	}
	if attrs["http.request_content_length"] != int64(200) {
		t.Errorf("expected http.request_content_length 200, got %v", attrs["http.request_content_length"])
	}
}

func TestTracingHook_DurationTracking(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	eventInfo := &webhook.EventInfo{
		EventType:   webhook.EventTypeDisputeStatusChange,
		EventID:     "evt-duration",
		BodySize:    100,
		IsDuplicate: false,
	}
	processingInfo := &webhook.ProcessingInfo{
		StatusCode: 200,
		Duration:   150 * time.Millisecond,
		Processed:  true,
	}

	ctx = hook.BeforeEvent(ctx, eventInfo)
	hook.AfterEvent(ctx, eventInfo, processingInfo)

	spans := exporter.GetSpans()
	attrs := attributesToMap(spans[0].Attributes)

	if attrs["evertec.duration_ms"] != int64(150) {
		t.Errorf("expected duration_ms 150, got %v", attrs["evertec.duration_ms"])
	}
}

// Helper function to convert attributes to map
func attributesToMap(attrs []attribute.KeyValue) map[string]any {
	result := make(map[string]any)
	for _, attr := range attrs {
		result[string(attr.Key)] = attr.Value.AsInterface()
	}
	return result
}
