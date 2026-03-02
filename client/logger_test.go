package client

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestLogRequest_NoLogger(t *testing.T) {
	c := &Client{} // No logger set

	// Should not panic
	c.logRequest(context.Background(), &RequestInfo{}, &ResponseInfo{})
}

func TestLogRequest_Success(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	c := &Client{logger: logger}

	reqInfo := &RequestInfo{
		Method:         "GET",
		Path:           "/accounts/123",
		BodySize:       0,
		IdempotencyKey: "",
	}
	respInfo := &ResponseInfo{
		StatusCode: 200,
		BodySize:   512,
		Duration:   45 * time.Millisecond,
		Error:      nil,
	}

	c.logRequest(context.Background(), reqInfo, respInfo)

	logOutput := buf.String()

	// Parse JSON log
	var logEntry map[string]any
	if err := json.Unmarshal([]byte(logOutput), &logEntry); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	// Verify log level (INFO for success)
	if level, ok := logEntry["level"].(string); !ok || level != "INFO" {
		t.Errorf("expected INFO level, got %v", logEntry["level"])
	}

	// Verify method
	if method, ok := logEntry["method"].(string); !ok || method != "GET" {
		t.Errorf("expected method GET, got %v", logEntry["method"])
	}

	// Verify path
	if path, ok := logEntry["path"].(string); !ok || path != "/accounts/123" {
		t.Errorf("expected path /accounts/123, got %v", logEntry["path"])
	}

	// Verify status
	if status, ok := logEntry["status"].(float64); !ok || status != 200 {
		t.Errorf("expected status 200, got %v", logEntry["status"])
	}

	// Verify response_size
	if respSize, ok := logEntry["response_size"].(float64); !ok || respSize != 512 {
		t.Errorf("expected response_size 512, got %v", logEntry["response_size"])
	}

	// Should NOT have idempotency_key (empty)
	if _, ok := logEntry["idempotency_key"]; ok {
		t.Error("expected no idempotency_key for GET request")
	}
}

func TestLogRequest_WithIdempotencyKey(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	c := &Client{logger: logger}

	reqInfo := &RequestInfo{
		Method:         "POST",
		Path:           "/accounts",
		BodySize:       256,
		IdempotencyKey: "test-key-123",
	}
	respInfo := &ResponseInfo{
		StatusCode: 201,
		BodySize:   128,
		Duration:   100 * time.Millisecond,
	}

	c.logRequest(context.Background(), reqInfo, respInfo)

	logOutput := buf.String()

	if !strings.Contains(logOutput, "test-key-123") {
		t.Error("expected idempotency_key in log output")
	}
}

func TestLogRequest_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	c := &Client{logger: logger}

	reqInfo := &RequestInfo{
		Method: "POST",
		Path:   "/accounts",
	}
	respInfo := &ResponseInfo{
		StatusCode: 400,
		Duration:   30 * time.Millisecond,
		Error: &APIError{
			StatusCode: 400,
			Code:       "INVALID_REQUEST",
			Message:    "Invalid document",
		},
	}

	c.logRequest(context.Background(), reqInfo, respInfo)

	logOutput := buf.String()

	// Parse JSON log
	var logEntry map[string]any
	if err := json.Unmarshal([]byte(logOutput), &logEntry); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	// Verify log level (ERROR for failures)
	if level, ok := logEntry["level"].(string); !ok || level != "ERROR" {
		t.Errorf("expected ERROR level, got %v", logEntry["level"])
	}

	// Verify error_code
	if code, ok := logEntry["error_code"].(string); !ok || code != "INVALID_REQUEST" {
		t.Errorf("expected error_code INVALID_REQUEST, got %v", logEntry["error_code"])
	}

	// Verify error_message
	if msg, ok := logEntry["error_message"].(string); !ok || msg != "Invalid document" {
		t.Errorf("expected error_message 'Invalid document', got %v", logEntry["error_message"])
	}
}

func TestLogRequest_ErrorStatus(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	c := &Client{logger: logger}

	// Test 4xx error without APIError
	reqInfo := &RequestInfo{Method: "GET", Path: "/test"}
	respInfo := &ResponseInfo{StatusCode: 404, Duration: time.Millisecond}

	c.logRequest(context.Background(), reqInfo, respInfo)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	// Should be ERROR level for 4xx
	if level := logEntry["level"]; level != "ERROR" {
		t.Errorf("expected ERROR level for 404, got %v", level)
	}
}

func TestLogRequest_5xxError(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	c := &Client{logger: logger}

	reqInfo := &RequestInfo{Method: "GET", Path: "/test"}
	respInfo := &ResponseInfo{StatusCode: 500, Duration: time.Millisecond}

	c.logRequest(context.Background(), reqInfo, respInfo)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("failed to parse log output: %v", err)
	}

	if level := logEntry["level"]; level != "ERROR" {
		t.Errorf("expected ERROR level for 500, got %v", level)
	}
}
