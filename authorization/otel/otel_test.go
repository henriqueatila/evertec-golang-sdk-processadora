package otel

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/authorization"
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

func TestTracingHook_PurchaseApproved(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	reqInfo := &authorization.RequestInfo{
		Path:          "/purchases",
		OperationType: "purchase",
		TransactionID: "txn-123",
		AccountID:     "acc-456",
		CardID:        "card-789",
		Amount:        10000,
		Currency:      "BRL",
		BodySize:      512,
	}
	respInfo := &authorization.ResponseInfo{
		StatusCode:   200,
		ResponseCode: 0,
		Duration:     50 * time.Millisecond,
		Approved:     true,
	}

	ctx = hook.BeforeAuthorization(ctx, reqInfo)
	hook.AfterAuthorization(ctx, reqInfo, respInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	span := spans[0]

	// Check span name
	expectedName := "evertec.authorization.purchase"
	if span.Name != expectedName {
		t.Errorf("expected span name %q, got %q", expectedName, span.Name)
	}

	// Check status is OK
	if span.Status.Code != codes.Ok {
		t.Errorf("expected Ok status, got %v", span.Status.Code)
	}

	// Check attributes
	attrs := attributesToMap(span.Attributes)

	if attrs["evertec.operation"] != "purchase" {
		t.Errorf("expected operation 'purchase', got %v", attrs["evertec.operation"])
	}
	if attrs["http.route"] != "/purchases" {
		t.Errorf("expected http.route '/purchases', got %v", attrs["http.route"])
	}
	if attrs["evertec.transaction_id"] != "txn-123" {
		t.Errorf("expected transaction_id 'txn-123', got %v", attrs["evertec.transaction_id"])
	}
	if attrs["evertec.account_id"] != "acc-456" {
		t.Errorf("expected account_id 'acc-456', got %v", attrs["evertec.account_id"])
	}
	if attrs["evertec.card_id"] != "card-789" {
		t.Errorf("expected card_id 'card-789', got %v", attrs["evertec.card_id"])
	}
	if attrs["evertec.amount"] != int64(10000) {
		t.Errorf("expected amount 10000, got %v", attrs["evertec.amount"])
	}
	if attrs["evertec.currency"] != "BRL" {
		t.Errorf("expected currency 'BRL', got %v", attrs["evertec.currency"])
	}
	if attrs["http.status_code"] != int64(200) {
		t.Errorf("expected http.status_code 200, got %v", attrs["http.status_code"])
	}
	if attrs["evertec.response_code"] != int64(0) {
		t.Errorf("expected response_code 0, got %v", attrs["evertec.response_code"])
	}
	if attrs["evertec.approved"] != true {
		t.Errorf("expected approved true, got %v", attrs["evertec.approved"])
	}
}

func TestTracingHook_PurchaseDeclined(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	reqInfo := &authorization.RequestInfo{
		Path:          "/purchases",
		OperationType: "purchase",
		TransactionID: "txn-declined",
		Amount:        50000,
		Currency:      "BRL",
		BodySize:      256,
	}
	respInfo := &authorization.ResponseInfo{
		StatusCode:   200,
		ResponseCode: 51, // Insufficient funds
		Duration:     30 * time.Millisecond,
		Approved:     false,
	}

	ctx = hook.BeforeAuthorization(ctx, reqInfo)
	hook.AfterAuthorization(ctx, reqInfo, respInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	span := spans[0]

	// Declined should still be OK status (business logic, not error)
	if span.Status.Code != codes.Ok {
		t.Errorf("expected Ok status for declined, got %v", span.Status.Code)
	}

	attrs := attributesToMap(span.Attributes)
	if attrs["evertec.approved"] != false {
		t.Errorf("expected approved false, got %v", attrs["evertec.approved"])
	}
	if attrs["evertec.response_code"] != int64(51) {
		t.Errorf("expected response_code 51, got %v", attrs["evertec.response_code"])
	}
}

func TestTracingHook_WithError(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	reqInfo := &authorization.RequestInfo{
		Path:          "/purchases",
		OperationType: "purchase",
		TransactionID: "txn-error",
		BodySize:      128,
	}
	testErr := errors.New("handler error: database connection failed")
	respInfo := &authorization.ResponseInfo{
		StatusCode:   500,
		ResponseCode: 96,
		Duration:     10 * time.Millisecond,
		Error:        testErr,
		Approved:     false,
	}

	ctx = hook.BeforeAuthorization(ctx, reqInfo)
	hook.AfterAuthorization(ctx, reqInfo, respInfo)

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

func TestTracingHook_Withdrawal(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	reqInfo := &authorization.RequestInfo{
		Path:          "/withdrawals",
		OperationType: "withdrawal",
		TransactionID: "txn-wd-123",
		Amount:        20000,
		Currency:      "USD",
		BodySize:      300,
	}
	respInfo := &authorization.ResponseInfo{
		StatusCode:   200,
		ResponseCode: 0,
		Duration:     45 * time.Millisecond,
		Approved:     true,
	}

	ctx = hook.BeforeAuthorization(ctx, reqInfo)
	hook.AfterAuthorization(ctx, reqInfo, respInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	if spans[0].Name != "evertec.authorization.withdrawal" {
		t.Errorf("expected span name 'evertec.authorization.withdrawal', got %q", spans[0].Name)
	}
}

func TestTracingHook_Query(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	reqInfo := &authorization.RequestInfo{
		Path:          "/queries",
		OperationType: "query",
		TransactionID: "txn-q-123",
		BodySize:      100,
	}
	respInfo := &authorization.ResponseInfo{
		StatusCode:   200,
		ResponseCode: 0,
		Duration:     20 * time.Millisecond,
		Approved:     true,
	}

	ctx = hook.BeforeAuthorization(ctx, reqInfo)
	hook.AfterAuthorization(ctx, reqInfo, respInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	if spans[0].Name != "evertec.authorization.query" {
		t.Errorf("expected span name 'evertec.authorization.query', got %q", spans[0].Name)
	}
}

func TestTracingHook_MinimalRequest(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	// Minimal request info - only required fields
	reqInfo := &authorization.RequestInfo{
		Path:          "/purchases",
		OperationType: "purchase",
		BodySize:      50,
	}
	respInfo := &authorization.ResponseInfo{
		StatusCode:   200,
		ResponseCode: 0,
		Duration:     15 * time.Millisecond,
		Approved:     true,
	}

	ctx = hook.BeforeAuthorization(ctx, reqInfo)
	hook.AfterAuthorization(ctx, reqInfo, respInfo)

	spans := exporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	attrs := attributesToMap(spans[0].Attributes)

	// Optional fields should not be present
	if _, exists := attrs["evertec.transaction_id"]; exists {
		t.Error("expected transaction_id to not be set for empty string")
	}
	if _, exists := attrs["evertec.account_id"]; exists {
		t.Error("expected account_id to not be set for empty string")
	}
	if _, exists := attrs["evertec.card_id"]; exists {
		t.Error("expected card_id to not be set for empty string")
	}
	if _, exists := attrs["evertec.amount"]; exists {
		t.Error("expected amount to not be set for zero value")
	}
}

func TestTracingHook_NoSpanInContext(t *testing.T) {
	hook, _ := setupTestTracer()

	// Call AfterAuthorization without BeforeAuthorization (no span in context)
	ctx := context.Background()
	reqInfo := &authorization.RequestInfo{Path: "/test", OperationType: "test"}
	respInfo := &authorization.ResponseInfo{StatusCode: 200}

	// Should not panic
	hook.AfterAuthorization(ctx, reqInfo, respInfo)
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

	// Verify hook implements authorization.Hook interface
	var _ authorization.Hook = hook
}

func TestTracingHook_SpanAttributes(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	reqInfo := &authorization.RequestInfo{
		Path:          "/transfers",
		OperationType: "transfer",
		BodySize:      200,
	}
	respInfo := &authorization.ResponseInfo{
		StatusCode:   200,
		ResponseCode: 0,
		Duration:     25 * time.Millisecond,
		Approved:     true,
	}

	ctx = hook.BeforeAuthorization(ctx, reqInfo)
	hook.AfterAuthorization(ctx, reqInfo, respInfo)

	spans := exporter.GetSpans()
	attrs := attributesToMap(spans[0].Attributes)

	// Check standard attributes
	if attrs["rpc.system"] != "evertec" {
		t.Errorf("expected rpc.system 'evertec', got %v", attrs["rpc.system"])
	}
	if attrs["rpc.service"] != "authorization" {
		t.Errorf("expected rpc.service 'authorization', got %v", attrs["rpc.service"])
	}
	if attrs["http.request_content_length"] != int64(200) {
		t.Errorf("expected http.request_content_length 200, got %v", attrs["http.request_content_length"])
	}
}

func TestTracingHook_DurationTracking(t *testing.T) {
	hook, exporter := setupTestTracer()

	ctx := context.Background()
	reqInfo := &authorization.RequestInfo{
		Path:          "/purchases",
		OperationType: "purchase",
		BodySize:      100,
	}
	respInfo := &authorization.ResponseInfo{
		StatusCode:   200,
		ResponseCode: 0,
		Duration:     150 * time.Millisecond,
		Approved:     true,
	}

	ctx = hook.BeforeAuthorization(ctx, reqInfo)
	hook.AfterAuthorization(ctx, reqInfo, respInfo)

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
