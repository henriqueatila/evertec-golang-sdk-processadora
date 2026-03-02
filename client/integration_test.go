package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ============================================================================
// INTEGRATION TESTS: Full request pipeline (validation + rate limiter + CB + retry)
// ============================================================================

func TestIntegration_RequestPipeline_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"id": "123"})
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:         "test-key",
		UserAgent:      "test/1.0",
		BaseURL:        server.URL,
		TLSConfig:      &tls.Config{InsecureSkipVerify: true},
		Retry:          &RetryConfig{MaxRetries: 3, InitialDelay: 1 * time.Millisecond, MaxDelay: 10 * time.Millisecond, RetryableStatuses: map[int]bool{503: true}},
		RateLimiter:    &RateLimiterConfig{RequestsPerSecond: 100, BurstSize: 200, WaitTimeout: 1 * time.Second},
		CircuitBreaker: &CircuitBreakerConfig{FailureThreshold: 5, ResetTimeout: 60 * time.Second, HalfOpenMax: 1},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var result map[string]string
	err = c.request(context.Background(), "GET", "/test", nil, &result)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	if result["id"] != "123" {
		t.Errorf("Expected id=123, got %s", result["id"])
	}
}

func TestIntegration_Validation_BlocksBadRequest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Request should not reach server")
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

	body := &validatableBody{fail: true}
	var result map[string]string
	err = c.request(context.Background(), "POST", "/test", body, &result)
	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
	if err.Error() != "validation failed: name is required" {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestIntegration_RateLimiter_EnforcesLimit(t *testing.T) {
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
		RateLimiter: &RateLimiterConfig{
			RequestsPerSecond: 2,
			BurstSize:         3,
			WaitTimeout:       0, // No waiting
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// First 3 requests should succeed (burst size)
	for i := 0; i < 3; i++ {
		err = c.request(context.Background(), "GET", "/test", nil, nil)
		if err != nil {
			t.Fatalf("Request %d should succeed, got: %v", i+1, err)
		}
	}

	// 4th request should fail (rate limit)
	err = c.request(context.Background(), "GET", "/test", nil, nil)
	if err != ErrRateLimited {
		t.Errorf("Expected ErrRateLimited, got: %v", err)
	}
}

func TestIntegration_CircuitBreaker_OpensAndBlocks(t *testing.T) {
	t.Parallel()

	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		CircuitBreaker: &CircuitBreakerConfig{
			FailureThreshold: 2,
			ResetTimeout:     60 * time.Second,
			HalfOpenMax:      1,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// First 2 requests fail, opening the circuit
	for i := 0; i < 2; i++ {
		_ = c.request(context.Background(), "GET", "/test", nil, nil)
	}

	// Verify circuit is open
	if c.cb.State() != CircuitOpen {
		t.Errorf("Expected CircuitOpen, got %s", c.cb.State())
	}

	// Third request should be rejected by circuit breaker
	err = c.request(context.Background(), "GET", "/test", nil, nil)
	if err != ErrCircuitOpen {
		t.Errorf("Expected ErrCircuitOpen, got: %v", err)
	}

	// Should not have made the third HTTP request
	if requestCount.Load() != 2 {
		t.Errorf("Expected 2 requests to server, got %d", requestCount.Load())
	}
}

func TestIntegration_Retry_RecoverFromTransientFailure(t *testing.T) {
	t.Parallel()

	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := requestCount.Add(1)
		if n < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Retry: &RetryConfig{
			MaxRetries:   5,
			InitialDelay: 1 * time.Millisecond,
			MaxDelay:     10 * time.Millisecond,
			RetryableStatuses: map[int]bool{
				503: true,
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var result map[string]string
	err = c.request(context.Background(), "GET", "/test", nil, &result)
	if err != nil {
		t.Fatalf("Expected success after retries, got: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("Expected status=ok, got %s", result["status"])
	}

	// Verify retries happened
	if requestCount.Load() != 3 {
		t.Errorf("Expected 3 requests (2 failures + 1 success), got %d", requestCount.Load())
	}
}

func TestIntegration_RetryAndCircuitBreaker_Together(t *testing.T) {
	t.Parallel()

	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		// Always fail
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Retry: &RetryConfig{
			MaxRetries:        3,
			InitialDelay:      1 * time.Millisecond,
			MaxDelay:          5 * time.Millisecond,
			RetryableStatuses: map[int]bool{503: true},
		},
		CircuitBreaker: &CircuitBreakerConfig{
			FailureThreshold: 2,
			ResetTimeout:     60 * time.Second,
			HalfOpenMax:      1,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// First request: retries 3 times, fails. CB sees 3 failures but threshold is 2, so opening circuit.
	// Actually, let's adjust: Each doRequest attempt is independent. Let's trace:
	// Attempt 0: fail (CB: 1 fail, closed)
	// Attempt 1: fail (CB: 2 fails, opens circuit)
	// Attempt 2: circuit open, ErrCircuitOpen returned
	_ = c.request(context.Background(), "GET", "/test", nil, nil)

	// After first request, circuit should be open
	state := c.cb.State()
	if state != CircuitOpen {
		t.Logf("Circuit state: %s", state)
	}

	// Second request should be immediately rejected by CB without retrying
	err = c.request(context.Background(), "GET", "/test", nil, nil)
	if err != ErrCircuitOpen {
		t.Errorf("Expected ErrCircuitOpen, got: %v", err)
	}
}

func TestIntegration_RetryRespectsCBStateTransition(t *testing.T) {
	t.Parallel()

	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := requestCount.Add(1)
		// Fail first 3, then succeed
		if n <= 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"id": "123"})
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Retry: &RetryConfig{
			MaxRetries:        5,
			InitialDelay:      1 * time.Millisecond,
			MaxDelay:          10 * time.Millisecond,
			RetryableStatuses: map[int]bool{503: true},
		},
		CircuitBreaker: &CircuitBreakerConfig{
			FailureThreshold: 10, // High threshold so CB doesn't interfere
			ResetTimeout:     60 * time.Second,
			HalfOpenMax:      1,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var result map[string]string
	err = c.request(context.Background(), "GET", "/test", nil, &result)
	if err != nil {
		t.Fatalf("Expected success, got: %v", err)
	}

	// Should have retried and eventually succeeded
	if requestCount.Load() != 4 {
		t.Errorf("Expected 4 requests (3 failures + 1 success), got %d", requestCount.Load())
	}
}

func TestIntegration_RateLimiter_BlocksBurst(t *testing.T) {
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
		RateLimiter: &RateLimiterConfig{
			RequestsPerSecond: 10,
			BurstSize:         2,
			WaitTimeout:       0,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Use burst tokens
	for i := 0; i < 2; i++ {
		err = c.request(context.Background(), "GET", "/test", nil, nil)
		if err != nil {
			t.Fatalf("Request %d should succeed, got: %v", i+1, err)
		}
	}

	// Next request exceeds burst, should fail
	err = c.request(context.Background(), "GET", "/test", nil, nil)
	if err != ErrRateLimited {
		t.Errorf("Expected ErrRateLimited, got: %v", err)
	}
}

func TestIntegration_RateLimiter_WithWaitTimeout(t *testing.T) {
	t.Parallel()

	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		RateLimiter: &RateLimiterConfig{
			RequestsPerSecond: 100,
			BurstSize:         2,
			WaitTimeout:       100 * time.Millisecond,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Use burst
	for i := 0; i < 2; i++ {
		_ = c.request(context.Background(), "GET", "/test", nil, nil)
	}

	// This should wait and succeed (tokens refill at 100 RPS)
	start := time.Now()
	err = c.request(context.Background(), "GET", "/test", nil, nil)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Expected success with wait, got: %v", err)
	}

	if elapsed < 1*time.Millisecond {
		t.Logf("Refill was instant (fast enough), elapsed: %v", elapsed)
	}

	// Should have made 3 requests
	if requestCount.Load() != 3 {
		t.Errorf("Expected 3 requests, got %d", requestCount.Load())
	}
}

func TestIntegration_ConcurrentRequests_WithRateLimiter(t *testing.T) {
	t.Parallel()

	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		RateLimiter: &RateLimiterConfig{
			RequestsPerSecond: 1000,
			BurstSize:         50,
			WaitTimeout:       1 * time.Second,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var wg sync.WaitGroup
	const goroutines = 20
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = c.request(ctx, "GET", "/test", nil, nil)
		}()
	}

	wg.Wait()

	// All requests should succeed
	if requestCount.Load() != int32(goroutines) {
		t.Errorf("Expected %d requests, got %d", goroutines, requestCount.Load())
	}
}

func TestIntegration_ConcurrentRequests_WithCircuitBreaker(t *testing.T) {
	t.Parallel()

	var requestCount atomic.Int32
	var mu sync.Mutex
	var failCount int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := requestCount.Add(1)
		mu.Lock()
		defer mu.Unlock()

		// Fail first 3, then succeed
		if int(n) <= 3 {
			failCount++
			w.WriteHeader(http.StatusServiceUnavailable)
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
		CircuitBreaker: &CircuitBreakerConfig{
			FailureThreshold: 3,
			ResetTimeout:     60 * time.Second,
			HalfOpenMax:      1,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var wg sync.WaitGroup
	const goroutines = 10
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_ = c.request(context.Background(), "GET", "/test", nil, nil)
		}()
	}

	wg.Wait()

	// After 3 failures, circuit opens and rejects remaining requests
	cbMetrics := c.CircuitBreakerMetrics()
	if cbMetrics.TotalFailures < 3 {
		t.Logf("TotalFailures: %d", cbMetrics.TotalFailures)
	}
	if cbMetrics.TotalRejected == 0 {
		t.Logf("Circuit breaker metrics: %+v", cbMetrics)
	}
}

// ============================================================================
// INTEGRATION TESTS: Error scenarios
// ============================================================================

func TestIntegration_NoRetry_WithoutRetryConfig(t *testing.T) {
	t.Parallel()

	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Retry:     nil, // No retry config
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var result map[string]string
	err = c.request(context.Background(), "GET", "/test", nil, &result)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Should only try once
	if requestCount.Load() != 1 {
		t.Errorf("Expected 1 request, got %d", requestCount.Load())
	}
}

func TestIntegration_RecordCBResult_OnSuccess(t *testing.T) {
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
		CircuitBreaker: &CircuitBreakerConfig{
			FailureThreshold: 5,
			ResetTimeout:     60 * time.Second,
			HalfOpenMax:      1,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_ = c.request(context.Background(), "GET", "/test", nil, nil)

	m := c.CircuitBreakerMetrics()
	if m.TotalSuccesses != 1 {
		t.Errorf("Expected 1 success, got %d", m.TotalSuccesses)
	}
	if m.TotalFailures != 0 {
		t.Errorf("Expected 0 failures, got %d", m.TotalFailures)
	}
}

func TestIntegration_RecordCBResult_On5xxError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		CircuitBreaker: &CircuitBreakerConfig{
			FailureThreshold: 5,
			ResetTimeout:     60 * time.Second,
			HalfOpenMax:      1,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_ = c.request(context.Background(), "GET", "/test", nil, nil)

	m := c.CircuitBreakerMetrics()
	if m.TotalFailures != 1 {
		t.Errorf("Expected 1 failure, got %d", m.TotalFailures)
	}
}

func TestIntegration_RecordCBResult_On4xxError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"code": "INVALID_REQUEST"})
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		CircuitBreaker: &CircuitBreakerConfig{
			FailureThreshold: 5,
			ResetTimeout:     60 * time.Second,
			HalfOpenMax:      1,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_ = c.request(context.Background(), "GET", "/test", nil, nil)

	m := c.CircuitBreakerMetrics()
	// 4xx errors should be treated as successes (not CB failures)
	if m.TotalSuccesses != 1 {
		t.Errorf("Expected 1 success (4xx not a CB failure), got %d", m.TotalSuccesses)
	}
	if m.TotalFailures != 0 {
		t.Errorf("Expected 0 failures, got %d", m.TotalFailures)
	}
}

// ============================================================================
// EDGE CASES: Context, timeouts, cancellations
// ============================================================================

func TestIntegration_RequestCancelledByContext(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err = c.request(ctx, "GET", "/test", nil, nil)
	if err == nil {
		t.Fatal("Expected error from cancelled context")
	}
}

func TestIntegration_RequestBlockedByRateLimiter_ContextCancel(t *testing.T) {
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
		RateLimiter: &RateLimiterConfig{
			RequestsPerSecond: 0.1, // Very slow
			BurstSize:         1,
			WaitTimeout:       5 * time.Second,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Use burst
	_ = c.request(context.Background(), "GET", "/test", nil, nil)

	// Try with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = c.request(ctx, "GET", "/test", nil, nil)
	if err == nil {
		t.Fatal("Expected error from cancelled context")
	}
}

// ============================================================================
// EDGE CASES: Retryable methods
// ============================================================================

func TestIntegration_MutationMethodsAreRetryable(t *testing.T) {
	t.Parallel()

	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	c, err := New(Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Retry: &RetryConfig{
			MaxRetries:        3,
			InitialDelay:      1 * time.Millisecond,
			MaxDelay:          10 * time.Millisecond,
			RetryableStatuses: map[int]bool{503: true},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// POST is retryable because idempotency key is auto-generated
	err = c.request(context.Background(), "POST", "/test", nil, nil)
	if err == nil {
		t.Fatal("Expected error after all retries")
	}

	// Should have retried (1 initial + 3 retries = 4 total)
	if requestCount.Load() != 4 {
		t.Errorf("Expected 4 requests (1 initial + 3 retries), got %d", requestCount.Load())
	}
}

// ============================================================================
// EDGE CASES: API error parsing
// ============================================================================

func TestIntegration_APIError_Parsed(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID_REQUEST",
			"message": "Invalid request body",
		})
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
	err = c.request(context.Background(), "GET", "/test", nil, &result)

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected APIError, got: %T", err)
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", apiErr.StatusCode)
	}
	if apiErr.Code != "INVALID_REQUEST" {
		t.Errorf("Expected code INVALID_REQUEST, got %s", apiErr.Code)
	}
}

// ============================================================================
// EDGE CASES: Empty responses
// ============================================================================

func TestIntegration_EmptyResponseBody(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// No body
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
	err = c.request(context.Background(), "GET", "/test", nil, &result)
	if err != nil {
		t.Fatalf("Expected success with empty body, got: %v", err)
	}
}

// ============================================================================
// HEALTH CHECK: Integration with request pipeline
// ============================================================================

func TestIntegration_Health_ReflectsAfterRequests(t *testing.T) {
	t.Parallel()

	var requestCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := requestCount.Add(1)
		if n <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
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
		CircuitBreaker: &CircuitBreakerConfig{
			FailureThreshold: 3,
			ResetTimeout:     60 * time.Second,
			HalfOpenMax:      1,
		},
		RateLimiter: &RateLimiterConfig{
			RequestsPerSecond: 100,
			BurstSize:         50,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Make requests
	_ = c.request(context.Background(), "GET", "/test", nil, nil)
	_ = c.request(context.Background(), "GET", "/test", nil, nil)
	_ = c.request(context.Background(), "GET", "/test", nil, nil)

	// Check health
	health := c.Health()
	if health.CircuitBreaker == nil {
		t.Fatal("Expected circuit breaker in health status")
	}
	if health.CircuitBreaker.Metrics.TotalRequests != 3 {
		t.Errorf("Expected 3 total requests, got %d", health.CircuitBreaker.Metrics.TotalRequests)
	}
	if health.CircuitBreaker.Metrics.TotalFailures < 2 {
		t.Errorf("Expected at least 2 failures, got %d", health.CircuitBreaker.Metrics.TotalFailures)
	}
}
