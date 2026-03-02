package client

import (
	"testing"
	"time"
)

func TestHealthNoCBNoRL(t *testing.T) {
	c := &Client{}
	h := c.Health()
	if !h.Healthy {
		t.Error("expected healthy without CB/RL")
	}
	if h.CircuitBreaker != nil {
		t.Error("expected nil circuit breaker")
	}
	if h.RateLimiter != nil {
		t.Error("expected nil rate limiter")
	}
}

func TestHealthWithCBClosed(t *testing.T) {
	c := &Client{
		cb: newCircuitBreaker(defaultCBConfig()),
	}
	h := c.Health()
	if !h.Healthy {
		t.Error("expected healthy when CB closed")
	}
	if h.CircuitBreaker == nil {
		t.Fatal("expected circuit breaker info")
	}
	if h.CircuitBreaker.State != "closed" {
		t.Errorf("expected state=closed, got %s", h.CircuitBreaker.State)
	}
}

func TestHealthWithCBOpen(t *testing.T) {
	cb := newCircuitBreaker(&CircuitBreakerConfig{
		FailureThreshold: 2,
		ResetTimeout:     60 * time.Second,
		HalfOpenMax:      1,
	})
	// Force open
	cb.recordFailure()
	cb.recordFailure()

	c := &Client{cb: cb}
	h := c.Health()
	if h.Healthy {
		t.Error("expected unhealthy when CB open")
	}
	if h.CircuitBreaker.State != "open" {
		t.Errorf("expected state=open, got %s", h.CircuitBreaker.State)
	}
}

func TestHealthWithRateLimiter(t *testing.T) {
	c := &Client{
		rl: newRateLimiter(&RateLimiterConfig{
			RequestsPerSecond: 50,
			BurstSize:         100,
		}),
	}
	h := c.Health()
	if !h.Healthy {
		t.Error("expected healthy")
	}
	if h.RateLimiter == nil {
		t.Fatal("expected rate limiter info")
	}
	if h.RateLimiter.RPS != 50 {
		t.Errorf("expected RPS=50, got %f", h.RateLimiter.RPS)
	}
	if h.RateLimiter.BurstSize != 100 {
		t.Errorf("expected BurstSize=100, got %d", h.RateLimiter.BurstSize)
	}
}

func TestCircuitBreakerMetricsNilWhenDisabled(t *testing.T) {
	c := &Client{}
	if m := c.CircuitBreakerMetrics(); m != nil {
		t.Error("expected nil metrics when CB disabled")
	}
}

func TestCircuitBreakerMetricsCounters(t *testing.T) {
	cb := newCircuitBreaker(&CircuitBreakerConfig{
		FailureThreshold: 5,
		ResetTimeout:     60 * time.Second,
		HalfOpenMax:      1,
	})

	// Simulate traffic
	cb.allow()
	cb.recordSuccess()
	cb.allow()
	cb.recordSuccess()
	cb.allow()
	cb.recordFailure()

	c := &Client{cb: cb}
	m := c.CircuitBreakerMetrics()

	if m.TotalRequests != 3 {
		t.Errorf("expected 3 requests, got %d", m.TotalRequests)
	}
	if m.TotalSuccesses != 2 {
		t.Errorf("expected 2 successes, got %d", m.TotalSuccesses)
	}
	if m.TotalFailures != 1 {
		t.Errorf("expected 1 failure, got %d", m.TotalFailures)
	}
	if m.TotalRejected != 0 {
		t.Errorf("expected 0 rejected, got %d", m.TotalRejected)
	}
}

func TestCircuitBreakerMetricsRejected(t *testing.T) {
	cb := newCircuitBreaker(&CircuitBreakerConfig{
		FailureThreshold: 1,
		ResetTimeout:     60 * time.Second,
		HalfOpenMax:      1,
	})

	cb.allow()
	cb.recordFailure() // opens circuit

	cb.allow() // rejected

	c := &Client{cb: cb}
	m := c.CircuitBreakerMetrics()

	if m.TotalRejected != 1 {
		t.Errorf("expected 1 rejected, got %d", m.TotalRejected)
	}
	if m.StateTransitions != 1 {
		t.Errorf("expected 1 state transition, got %d", m.StateTransitions)
	}
}

func TestCircuitBreakerMetricsStateTransitions(t *testing.T) {
	cb := newCircuitBreaker(&CircuitBreakerConfig{
		FailureThreshold: 1,
		ResetTimeout:     1 * time.Millisecond,
		HalfOpenMax:      1,
	})

	// Closed -> Open (1 transition)
	cb.allow()
	cb.recordFailure()

	time.Sleep(5 * time.Millisecond)

	// Open -> HalfOpen (2nd transition) -> Closed (3rd transition)
	cb.allow()
	cb.recordSuccess()

	m := cb.metrics()
	if m.StateTransitions != 3 {
		t.Errorf("expected 3 transitions (closed->open->half-open->closed), got %d", m.StateTransitions)
	}
}
