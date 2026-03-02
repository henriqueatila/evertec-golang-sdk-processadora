package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// newRetryTestClient creates a client connected to a mock server with retry config applied.
func newRetryTestClient(t *testing.T, server *httptest.Server, retryCfg *RetryConfig) *Client {
	t.Helper()
	c, err := New(Config{
		APIKey:    "test-api-key-12345",
		UserAgent: "TestEmissor/1.0.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		Retry:     retryCfg,
	})
	if err != nil {
		t.Fatalf("failed to create retry test client: %v", err)
	}
	return c
}

// ============================================================================
// UNIT TESTS: RetryConfig defaults and options
// ============================================================================

func TestRetryConfig_Defaults(t *testing.T) {
	t.Parallel()

	cfg := defaultRetryConfig()

	if cfg.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want 3", cfg.MaxRetries)
	}
	if cfg.InitialDelay != 500*time.Millisecond {
		t.Errorf("InitialDelay = %v, want 500ms", cfg.InitialDelay)
	}
	if cfg.MaxDelay != 30*time.Second {
		t.Errorf("MaxDelay = %v, want 30s", cfg.MaxDelay)
	}

	retryable := []int{429, 500, 502, 503, 504}
	for _, code := range retryable {
		if !cfg.RetryableStatuses[code] {
			t.Errorf("status %d should be retryable by default", code)
		}
	}
}

func TestRetryOptions(t *testing.T) {
	t.Parallel()

	t.Run("MaxRetries", func(t *testing.T) {
		t.Parallel()
		cfg := defaultRetryConfig()
		MaxRetries(10)(cfg)
		if cfg.MaxRetries != 10 {
			t.Errorf("MaxRetries = %d, want 10", cfg.MaxRetries)
		}
	})

	t.Run("InitialDelay", func(t *testing.T) {
		t.Parallel()
		cfg := defaultRetryConfig()
		InitialDelay(200 * time.Millisecond)(cfg)
		if cfg.InitialDelay != 200*time.Millisecond {
			t.Errorf("InitialDelay = %v, want 200ms", cfg.InitialDelay)
		}
	})

	t.Run("MaxDelay", func(t *testing.T) {
		t.Parallel()
		cfg := defaultRetryConfig()
		MaxDelay(10 * time.Second)(cfg)
		if cfg.MaxDelay != 10*time.Second {
			t.Errorf("MaxDelay = %v, want 10s", cfg.MaxDelay)
		}
	})

	t.Run("RetryableStatuses", func(t *testing.T) {
		t.Parallel()
		cfg := defaultRetryConfig()
		RetryableStatuses(503, 429)(cfg)
		if !cfg.RetryableStatuses[503] {
			t.Error("503 should be retryable")
		}
		if !cfg.RetryableStatuses[429] {
			t.Error("429 should be retryable")
		}
		// 500 should no longer be retryable after override
		if cfg.RetryableStatuses[500] {
			t.Error("500 should NOT be retryable after override")
		}
	})
}

// ============================================================================
// UNIT TESTS: isRetryable
// ============================================================================

func TestIsRetryable(t *testing.T) {
	t.Parallel()

	cfg := defaultRetryConfig()

	retryable := []int{429, 500, 502, 503, 504}
	for _, code := range retryable {
		if !cfg.isRetryable(code) {
			t.Errorf("status %d should be retryable", code)
		}
	}

	notRetryable := []int{200, 201, 400, 401, 403, 404, 409, 422}
	for _, code := range notRetryable {
		if cfg.isRetryable(code) {
			t.Errorf("status %d should NOT be retryable", code)
		}
	}
}

// ============================================================================
// UNIT TESTS: isRetryableMethod
// ============================================================================

func TestIsRetryableMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		method            string
		hasIdempotencyKey bool
		want              bool
	}{
		{"GET always retryable", http.MethodGet, false, true},
		{"GET with key also retryable", http.MethodGet, true, true},
		{"POST with idempotency key is retryable", http.MethodPost, true, true},
		{"POST without key is NOT retryable", http.MethodPost, false, false},
		{"PATCH with key is retryable", http.MethodPatch, true, true},
		{"PATCH without key is NOT retryable", http.MethodPatch, false, false},
		{"PUT with key is retryable", http.MethodPut, true, true},
		{"DELETE with key is retryable", http.MethodDelete, true, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := isRetryableMethod(tt.method, tt.hasIdempotencyKey)
			if got != tt.want {
				t.Errorf("isRetryableMethod(%q, %v) = %v, want %v", tt.method, tt.hasIdempotencyKey, got, tt.want)
			}
		})
	}
}

// ============================================================================
// UNIT TESTS: calculateBackoff
// ============================================================================

func TestCalculateBackoff(t *testing.T) {
	t.Parallel()

	initial := 100 * time.Millisecond
	maxDelay := 10 * time.Second

	// Run multiple times to account for jitter randomness
	for i := 0; i < 50; i++ {
		delay0 := calculateBackoff(0, initial, maxDelay)
		delay1 := calculateBackoff(1, initial, maxDelay)
		delay2 := calculateBackoff(2, initial, maxDelay)

		// Attempt 0: base = 100ms, jitter [50ms, 150ms]
		if delay0 < 50*time.Millisecond || delay0 > 150*time.Millisecond {
			t.Errorf("attempt 0 delay %v out of expected range [50ms, 150ms]", delay0)
		}

		// Attempt 1: base = 200ms, jitter [100ms, 300ms]
		if delay1 < 100*time.Millisecond || delay1 > 300*time.Millisecond {
			t.Errorf("attempt 1 delay %v out of expected range [100ms, 300ms]", delay1)
		}

		// Attempt 2: base = 400ms, jitter [200ms, 600ms]
		if delay2 < 200*time.Millisecond || delay2 > 600*time.Millisecond {
			t.Errorf("attempt 2 delay %v out of expected range [200ms, 600ms]", delay2)
		}
	}
}

func TestCalculateBackoff_Cap(t *testing.T) {
	t.Parallel()

	initial := 1 * time.Second
	maxDelay := 5 * time.Second

	// Attempt 10: base = 1024s, capped at 5s. With jitter [2.5s, 7.5s] but capped input means result <= 7.5s.
	// More precisely: delay = min(1s * 2^10, 5s) * jitter = 5s * [0.5, 1.5) = [2.5s, 7.5s).
	for i := 0; i < 20; i++ {
		delay := calculateBackoff(10, initial, maxDelay)
		// The raw cap is applied before jitter, so delay = maxDelay * jitter
		if delay > time.Duration(float64(maxDelay)*1.5)+time.Millisecond {
			t.Errorf("delay %v exceeds maxDelay * 1.5 (%v)", delay, time.Duration(float64(maxDelay)*1.5))
		}
		if delay < time.Duration(float64(maxDelay)*0.5)-time.Millisecond {
			t.Errorf("delay %v below maxDelay * 0.5 (%v)", delay, time.Duration(float64(maxDelay)*0.5))
		}
	}
}

// ============================================================================
// UNIT TESTS: parseRetryAfter
// ============================================================================

func TestParseRetryAfter(t *testing.T) {
	t.Parallel()

	t.Run("nil response", func(t *testing.T) {
		t.Parallel()
		got := parseRetryAfter(nil)
		if got != 0 {
			t.Errorf("parseRetryAfter(nil) = %v, want 0", got)
		}
	})

	t.Run("missing header", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{Header: http.Header{}}
		got := parseRetryAfter(resp)
		if got != 0 {
			t.Errorf("missing header: got %v, want 0", got)
		}
	})

	t.Run("valid seconds", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"5"}}}
		got := parseRetryAfter(resp)
		if got != 5*time.Second {
			t.Errorf("parseRetryAfter = %v, want 5s", got)
		}
	})

	t.Run("invalid value (text)", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"Wed, 21 Oct 2025 07:28:00 GMT"}}}
		got := parseRetryAfter(resp)
		if got != 0 {
			t.Errorf("invalid header: got %v, want 0", got)
		}
	})

	t.Run("zero seconds", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"0"}}}
		got := parseRetryAfter(resp)
		if got != 0 {
			t.Errorf("zero seconds: got %v, want 0", got)
		}
	})

	t.Run("negative seconds", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{Header: http.Header{"Retry-After": []string{"-1"}}}
		got := parseRetryAfter(resp)
		if got != 0 {
			t.Errorf("negative seconds: got %v, want 0", got)
		}
	})
}

// ============================================================================
// INTEGRATION TESTS
// ============================================================================

func TestRetryIntegration(t *testing.T) {
	t.Parallel()

	// Server fails with 503 for the first 2 requests, then succeeds.
	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := requestCount.Add(1)
		if n <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"accountId":"acc-123"}`))
	}))
	defer server.Close()

	retryCfg := &RetryConfig{
		MaxRetries:   3,
		InitialDelay: 1 * time.Millisecond,
		MaxDelay:     10 * time.Millisecond,
		RetryableStatuses: map[int]bool{
			503: true,
		},
	}
	c := newRetryTestClient(t, server, retryCfg)

	_, err := c.GetAccount(context.Background(), "acc-123")
	if err != nil {
		t.Fatalf("expected success after retries, got error: %v", err)
	}

	if int(requestCount.Load()) != 3 {
		t.Errorf("expected 3 total requests (2 failures + 1 success), got %d", requestCount.Load())
	}
}

func TestRetryWithContextCancellation(t *testing.T) {
	t.Parallel()

	// Server always returns 503 to keep triggering retries.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	retryCfg := &RetryConfig{
		MaxRetries:        10,
		InitialDelay:      50 * time.Millisecond, // long enough for context to cancel
		MaxDelay:          1 * time.Second,
		RetryableStatuses: map[int]bool{503: true},
	}
	c := newRetryTestClient(t, server, retryCfg)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := c.GetAccount(ctx, "acc-123")
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected error due to context cancellation, got nil")
	}

	// Should exit well before MaxRetries * InitialDelay (500ms)
	if elapsed > 500*time.Millisecond {
		t.Errorf("context cancellation took too long: %v (expected < 500ms)", elapsed)
	}
}

func TestRetryWith429AndRetryAfter(t *testing.T) {
	t.Parallel()

	// Server returns 429 once then succeeds.
	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := requestCount.Add(1)
		if n == 1 {
			w.Header().Set("Retry-After", "1") // 1 second (not used by current retry loop, but header is present)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"accountId":"acc-456"}`))
	}))
	defer server.Close()

	retryCfg := &RetryConfig{
		MaxRetries:        2,
		InitialDelay:      1 * time.Millisecond,
		MaxDelay:          10 * time.Millisecond,
		RetryableStatuses: map[int]bool{429: true},
	}
	c := newRetryTestClient(t, server, retryCfg)

	_, err := c.GetAccount(context.Background(), "acc-456")
	if err != nil {
		t.Fatalf("expected success after retry on 429, got: %v", err)
	}

	if int(requestCount.Load()) != 2 {
		t.Errorf("expected 2 requests (1 failure + 1 success), got %d", requestCount.Load())
	}
}

// ============================================================================
// WithRetry option test
// ============================================================================

func TestWithRetry_Option(t *testing.T) {
	t.Parallel()

	cfg := &Config{}
	WithRetry()(cfg)

	if cfg.Retry == nil {
		t.Fatal("WithRetry() should set Retry config")
	}
	if cfg.Retry.MaxRetries != 3 {
		t.Errorf("default MaxRetries = %d, want 3", cfg.Retry.MaxRetries)
	}

	// With custom options
	cfg2 := &Config{}
	WithRetry(MaxRetries(5), InitialDelay(100*time.Millisecond))(cfg2)

	if cfg2.Retry.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d, want 5", cfg2.Retry.MaxRetries)
	}
	if cfg2.Retry.InitialDelay != 100*time.Millisecond {
		t.Errorf("InitialDelay = %v, want 100ms", cfg2.Retry.InitialDelay)
	}
}
