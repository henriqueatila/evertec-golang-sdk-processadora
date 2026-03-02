package authorization

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestNewLoggingHook(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	hook := NewLoggingHook(logger)

	if hook == nil {
		t.Fatal("NewLoggingHook returned nil")
	}
	if hook.logger != logger {
		t.Error("logger not set correctly")
	}
}

func TestLoggingHook_ImplementsInterface(t *testing.T) {
	var _ Hook = &LoggingHook{}
	var _ Hook = NewLoggingHook(nil)
}

func TestLoggingHook_BeforeAuthorization(t *testing.T) {
	hook := NewLoggingHook(nil)
	ctx := context.Background()
	req := &RequestInfo{Path: "/test"}

	newCtx := hook.BeforeAuthorization(ctx, req)

	// Should return the same context unchanged
	if newCtx != ctx {
		t.Error("BeforeAuthorization should return unchanged context")
	}
}

func TestLoggingHook_AfterAuthorization_NilLogger(t *testing.T) {
	hook := NewLoggingHook(nil)
	req := &RequestInfo{Path: "/test"}
	resp := &ResponseInfo{StatusCode: 200, Approved: true}

	// Should not panic with nil logger
	hook.AfterAuthorization(context.Background(), req, resp)
}

func TestLoggingHook_AfterAuthorization_Approved(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	hook := NewLoggingHook(logger)

	req := &RequestInfo{
		Path:          "/purchases",
		OperationType: "purchase",
		TransactionID: "txn123",
		AccountID:     "acc456",
		CardID:        "card789",
		Amount:        10000,
		Currency:      "BRL",
		BodySize:      512,
	}
	resp := &ResponseInfo{
		StatusCode:   200,
		ResponseCode: 0,
		Duration:     50 * time.Millisecond,
		Approved:     true,
	}

	hook.AfterAuthorization(context.Background(), req, resp)

	// Parse log output
	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log: %v\nLog: %s", err, buf.String())
	}

	// Verify log level is INFO for approved transactions
	if level, ok := logEntry["level"].(string); ok {
		if level != "INFO" {
			t.Errorf("level = %s, want INFO", level)
		}
	}

	// Verify message
	if msg, ok := logEntry["msg"].(string); ok {
		if msg != "evertec_authorization" {
			t.Errorf("msg = %s, want evertec_authorization", msg)
		}
	}

	// Verify required fields
	expectedFields := map[string]any{
		"path":           "/purchases",
		"operation":      "purchase",
		"status":         float64(200),
		"response_code":  float64(0),
		"approved":       true,
		"request_size":   float64(512),
		"transaction_id": "txn123",
		"account_id":     "acc456",
		"card_id":        "card789",
		"amount":         float64(10000),
		"currency":       "BRL",
	}

	for key, expected := range expectedFields {
		if actual, ok := logEntry[key]; ok {
			if actual != expected {
				t.Errorf("%s = %v, want %v", key, actual, expected)
			}
		} else {
			t.Errorf("missing field: %s", key)
		}
	}
}

func TestLoggingHook_AfterAuthorization_Declined(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	hook := NewLoggingHook(logger)

	req := &RequestInfo{
		Path:          "/purchases",
		OperationType: "purchase",
		BodySize:      256,
	}
	resp := &ResponseInfo{
		StatusCode:   200,
		ResponseCode: 51, // Insufficient funds
		Duration:     30 * time.Millisecond,
		Approved:     false,
	}

	hook.AfterAuthorization(context.Background(), req, resp)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log: %v", err)
	}

	// Verify log level is WARN for declined transactions
	if level, ok := logEntry["level"].(string); ok {
		if level != "WARN" {
			t.Errorf("level = %s, want WARN for declined", level)
		}
	}

	if approved, ok := logEntry["approved"].(bool); ok {
		if approved {
			t.Error("approved = true, want false")
		}
	}
}

func TestLoggingHook_AfterAuthorization_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	hook := NewLoggingHook(logger)

	req := &RequestInfo{
		Path:          "/purchases",
		OperationType: "purchase",
		BodySize:      128,
	}
	resp := &ResponseInfo{
		StatusCode:   500,
		ResponseCode: 0,
		Duration:     10 * time.Millisecond,
		Error:        errors.New("internal error"),
		Approved:     false,
	}

	hook.AfterAuthorization(context.Background(), req, resp)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log: %v", err)
	}

	// Verify log level is ERROR for errors
	if level, ok := logEntry["level"].(string); ok {
		if level != "ERROR" {
			t.Errorf("level = %s, want ERROR", level)
		}
	}

	// Verify error field is present
	if errMsg, ok := logEntry["error"].(string); ok {
		if errMsg != "internal error" {
			t.Errorf("error = %s, want 'internal error'", errMsg)
		}
	} else {
		t.Error("missing error field in log")
	}
}

func TestLoggingHook_AfterAuthorization_ServerError(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	hook := NewLoggingHook(logger)

	req := &RequestInfo{Path: "/test", OperationType: "test", BodySize: 64}
	resp := &ResponseInfo{
		StatusCode: 503, // Service unavailable
		Duration:   5 * time.Millisecond,
		Approved:   false,
	}

	hook.AfterAuthorization(context.Background(), req, resp)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log: %v", err)
	}

	// Should be ERROR level for 5xx status codes
	if level, ok := logEntry["level"].(string); ok {
		if level != "ERROR" {
			t.Errorf("level = %s, want ERROR for 5xx status", level)
		}
	}
}

func TestLoggingHook_AfterAuthorization_MinimalFields(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	hook := NewLoggingHook(logger)

	// Request with minimal fields (no optional fields)
	req := &RequestInfo{
		Path:          "/status",
		OperationType: "status",
		BodySize:      0,
	}
	resp := &ResponseInfo{
		StatusCode:   200,
		ResponseCode: 0,
		Duration:     1 * time.Millisecond,
		Approved:     true,
	}

	hook.AfterAuthorization(context.Background(), req, resp)

	logStr := buf.String()

	// Should not contain optional fields when empty
	if strings.Contains(logStr, "transaction_id") {
		t.Error("log should not contain transaction_id when empty")
	}
	if strings.Contains(logStr, "account_id") {
		t.Error("log should not contain account_id when empty")
	}
	if strings.Contains(logStr, "card_id") {
		t.Error("log should not contain card_id when empty")
	}
}

func TestLoggingHook_AfterAuthorization_AmountWithoutCurrency(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	hook := NewLoggingHook(logger)

	req := &RequestInfo{
		Path:          "/test",
		OperationType: "test",
		Amount:        5000,
		// Currency is empty
		BodySize: 100,
	}
	resp := &ResponseInfo{StatusCode: 200, Approved: true, Duration: time.Millisecond}

	hook.AfterAuthorization(context.Background(), req, resp)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log: %v", err)
	}

	// Amount should be present
	if _, ok := logEntry["amount"]; !ok {
		t.Error("amount should be present")
	}

	// Currency should not be present when empty
	if _, ok := logEntry["currency"]; ok {
		t.Error("currency should not be present when empty")
	}
}

func TestLoggingHook_AfterAuthorization_ZeroAmount(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	hook := NewLoggingHook(logger)

	req := &RequestInfo{
		Path:          "/queries",
		OperationType: "query",
		Amount:        0, // Zero amount (e.g., balance inquiry)
		Currency:      "BRL",
		BodySize:      100,
	}
	resp := &ResponseInfo{StatusCode: 200, Approved: true, Duration: time.Millisecond}

	hook.AfterAuthorization(context.Background(), req, resp)

	logStr := buf.String()

	// Amount should not be logged when zero
	if strings.Contains(logStr, `"amount"`) {
		t.Error("amount should not be logged when zero")
	}
}
