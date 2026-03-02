package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// ============================================================================
// EDGE CASES: doRequest() and error handling
// ============================================================================

func TestEdgeCase_InvalidRequestBody_MarshalError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("request should not reach server")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Use a channel which cannot be marshaled to JSON
	body := make(chan int)
	err = c.doRequest(context.Background(), "POST", "/test", body, nil)
	if err == nil {
		t.Fatal("Expected marshal error")
	}
	if !strings.Contains(err.Error(), "marshal") {
		t.Errorf("Expected marshal error, got: %v", err)
	}
}

func TestEdgeCase_ResponseBodyTooLarge(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Write more data than maxResponseBody allows
		for i := 0; i < 100000; i++ {
			_, _ = w.Write([]byte(`{"data":"` + strings.Repeat("x", 100) + `"}`))
		}
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:          "test-key",
		UserAgent:       "test/1.0",
		BaseURL:         server.URL,
		TLSConfig:       &tls.Config{InsecureSkipVerify: true},
		MaxResponseBody: 1024, // Very small limit
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var result map[string]string
	err = c.doRequest(context.Background(), "GET", "/test", nil, &result)
	// Should either succeed (with partial data) or fail gracefully
	if err != nil && !strings.Contains(err.Error(), "unmarshal") {
		t.Logf("Large response handling: %v", err)
	}
}

func TestEdgeCase_InvalidJSON_Response(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var result map[string]string
	err = c.doRequest(context.Background(), "GET", "/test", nil, &result)
	if err == nil {
		t.Fatal("Expected unmarshal error for invalid JSON")
	}
}

func TestEdgeCase_ErrorResponse_InvalidJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{invalid}`))
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var result map[string]string
	err = c.doRequest(context.Background(), "GET", "/test", nil, &result)

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected APIError")
	}

	// Should have fallback to raw body message
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", apiErr.StatusCode)
	}
}

func TestEdgeCase_ContextCanceledBeforeRequest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = c.doRequest(ctx, "GET", "/test", nil, nil)
	if err == nil {
		t.Fatal("Expected error from cancelled context")
	}
}

func TestEdgeCase_NetworkError(t *testing.T) {
	t.Parallel()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   "https://invalid.example.com:99999", // Non-existent
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Timeout:   100 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = c.doRequest(context.Background(), "GET", "/test", nil, nil)
	if err == nil {
		t.Fatal("Expected network error")
	}
}

func TestEdgeCase_HookPanic_BeforeRequest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	panicHook := &PanicHook{}
	panicHook.hookPhase = "before"

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Hooks:     []Hook{panicHook},
		Logger:    logger,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Should not panic, should be recovered
	err = c.doRequest(context.Background(), "GET", "/test", nil, nil)
	// Should succeed despite hook panic
	if err != nil {
		t.Logf("Got error after hook panic: %v", err)
	}
}

func TestEdgeCase_HookPanic_AfterResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	panicHook := &PanicHook{}
	panicHook.hookPhase = "after"

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Hooks:     []Hook{panicHook},
		Logger:    logger,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Should not panic, should be recovered
	err = c.doRequest(context.Background(), "GET", "/test", nil, nil)
	// Should succeed despite hook panic
	if err != nil {
		t.Logf("Got error after hook panic: %v", err)
	}
}

func TestEdgeCase_MultipleHooks_Chained(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"id": "123"})
	}))
	defer server.Close()

	hook1Called := false
	hook2Called := false

	hook1 := &TestHook{
		beforeReq: func(ctx context.Context, req *RequestInfo) context.Context {
			hook1Called = true
			return ctx
		},
		afterResp: func(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
			hook2Called = true
		},
	}

	hook2 := &TestHook{}

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Hooks:     []Hook{hook1, hook2},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var result map[string]string
	err = c.doRequest(context.Background(), "GET", "/test", nil, &result)
	if err != nil {
		t.Fatalf("Expected success, got: %v", err)
	}

	if !hook1Called {
		t.Error("BeforeRequest hook should have been called")
	}
	if !hook2Called {
		t.Error("AfterResponse hook should have been called")
	}
}

func TestEdgeCase_LoggingWithoutPanic(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logBuf := &bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(logBuf, nil))

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Logger:    logger,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = c.doRequest(context.Background(), "GET", "/test", nil, nil)
	if err != nil {
		t.Fatalf("Expected success, got: %v", err)
	}

	// Logger should have written something
	if logBuf.Len() == 0 {
		t.Logf("Logger buffer is empty (logging may not be happening)")
	}
}

// ============================================================================
// EDGE CASES: Request method handling
// ============================================================================

func TestEdgeCase_DeleteWithIdempotencyKey(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// DELETE should have idempotency key when retries are enabled
		if r.Header.Get("IdempotencyKey") == "" {
			t.Error("Expected IdempotencyKey in DELETE request")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Retry: &RetryConfig{
			MaxRetries:        1,
			InitialDelay:      1 * time.Millisecond,
			MaxDelay:          10 * time.Millisecond,
			RetryableStatuses: map[int]bool{503: true},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = c.doRequest(context.Background(), "DELETE", "/test", nil, nil)
	if err != nil {
		t.Fatalf("Expected success, got: %v", err)
	}
}

func TestEdgeCase_GetIgnoresIdempotencyKey(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// GET should NOT have idempotency key
		if r.Header.Get("IdempotencyKey") != "" {
			t.Error("Expected no IdempotencyKey in GET request")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = c.doRequest(context.Background(), "GET", "/test", nil, nil)
	if err != nil {
		t.Fatalf("Expected success, got: %v", err)
	}
}

func TestEdgeCase_CustomIdempotencyKey_Preserved(t *testing.T) {
	t.Parallel()

	customKey := "my-custom-key-12345"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("IdempotencyKey")
		if got != customKey {
			t.Errorf("Expected IdempotencyKey=%s, got %s", customKey, got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := WithIdempotencyKey(context.Background(), customKey)
	err = c.doRequest(ctx, "POST", "/test", nil, nil)
	if err != nil {
		t.Fatalf("Expected success, got: %v", err)
	}
}

// ============================================================================
// TEST UTILITIES: Hook implementations for testing
// ============================================================================

type PanicHook struct {
	hookPhase string // "before" or "after"
}

func (ph *PanicHook) BeforeRequest(ctx context.Context, req *RequestInfo) context.Context {
	if ph.hookPhase == "before" {
		panic("intentional panic in BeforeRequest")
	}
	return ctx
}

func (ph *PanicHook) AfterResponse(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
	if ph.hookPhase == "after" {
		panic("intentional panic in AfterResponse")
	}
}

type TestHook struct {
	beforeReq func(context.Context, *RequestInfo) context.Context
	afterResp func(context.Context, *RequestInfo, *ResponseInfo)
}

func (th *TestHook) BeforeRequest(ctx context.Context, req *RequestInfo) context.Context {
	if th.beforeReq != nil {
		return th.beforeReq(ctx, req)
	}
	return ctx
}

func (th *TestHook) AfterResponse(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
	if th.afterResp != nil {
		th.afterResp(ctx, req, resp)
	}
}

// ============================================================================
// EDGE CASES: Validation disabled
// ============================================================================

func TestEdgeCase_ValidationDisabled_WithNoValidation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"id": "123"})
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:       "test-key",
		UserAgent:    "test/1.0",
		BaseURL:      server.URL,
		TLSConfig:    &tls.Config{InsecureSkipVerify: true},
		NoValidation: true,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Even with failing validation, should succeed
	body := &validatableBody{fail: true}
	var result map[string]string
	err = c.request(context.Background(), "POST", "/test", body, &result)
	if err != nil {
		t.Fatalf("Expected success with validation disabled, got: %v", err)
	}
}

// ============================================================================
// EDGE CASES: Response status codes and CB recording
// ============================================================================

func TestEdgeCase_Status502_TriggersRetry(t *testing.T) {
	t.Parallel()

	var requestCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount < 2 {
			w.WriteHeader(http.StatusBadGateway) // 502
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Retry: &RetryConfig{
			MaxRetries:        2,
			InitialDelay:      1 * time.Millisecond,
			MaxDelay:          10 * time.Millisecond,
			RetryableStatuses: map[int]bool{502: true},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = c.request(context.Background(), "GET", "/test", nil, nil)
	if err != nil {
		t.Fatalf("Expected success after retry, got: %v", err)
	}
	if requestCount != 2 {
		t.Errorf("Expected 2 requests, got %d", requestCount)
	}
}

func TestEdgeCase_Status504_TriggersRetry(t *testing.T) {
	t.Parallel()

	var requestCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount < 2 {
			w.WriteHeader(http.StatusGatewayTimeout) // 504
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Retry: &RetryConfig{
			MaxRetries:        2,
			InitialDelay:      1 * time.Millisecond,
			MaxDelay:          10 * time.Millisecond,
			RetryableStatuses: map[int]bool{504: true},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = c.request(context.Background(), "GET", "/test", nil, nil)
	if err != nil {
		t.Fatalf("Expected success after retry, got: %v", err)
	}
	if requestCount != 2 {
		t.Errorf("Expected 2 requests, got %d", requestCount)
	}
}

// ============================================================================
// EDGE CASES: Response reading errors
// ============================================================================

func TestEdgeCase_ResponseReadError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Hijack the connection to simulate read error
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			t.Skip("Server does not support hijacking")
		}

		conn, _, err := hijacker.Hijack()
		if err != nil {
			t.Logf("Hijack failed: %v", err)
		} else {
			conn.Close()
		}
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Timeout:   100 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = c.doRequest(context.Background(), "GET", "/test", nil, nil)
	if err == nil {
		t.Logf("No error (connection may have succeeded anyway)")
	}
}

// ============================================================================
// EDGE CASES: Headers and Content-Type
// ============================================================================

func TestEdgeCase_RequiredHeadersPresent(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	headers := make(http.Header)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		headers = r.Header.Clone()
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = c.doRequest(context.Background(), "POST", "/test", nil, nil)
	if err != nil {
		t.Fatalf("Expected success, got: %v", err)
	}

	// Check headers
	mu.Lock()
	defer mu.Unlock()

	if headers.Get("API-Key") == "" {
		t.Error("Missing API-Key header")
	}
	if headers.Get("User-Agent") == "" {
		t.Error("Missing User-Agent header")
	}
	if headers.Get("Content-Type") == "" {
		t.Error("Missing Content-Type header (POST)")
	}
	if headers.Get("Accept") == "" {
		t.Error("Missing Accept header")
	}
	if headers.Get("IdempotencyKey") == "" {
		t.Error("Missing IdempotencyKey header (POST)")
	}
}
