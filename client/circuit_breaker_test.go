package client

import (
	"sync"
	"testing"
	"time"
)

func TestCircuitBreaker_ClosedState(t *testing.T) {
	cb := newCircuitBreaker(defaultCBConfig())

	if !cb.allow() {
		t.Fatal("expected allow() to return true in closed state")
	}
	if cb.State() != CircuitClosed {
		t.Fatalf("expected state closed, got %s", cb.State())
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	cfg := defaultCBConfig()
	cfg.FailureThreshold = 3
	cb := newCircuitBreaker(cfg)

	for i := 0; i < 3; i++ {
		cb.recordFailure()
	}

	if cb.State() != CircuitOpen {
		t.Fatalf("expected state open after %d failures, got %s", cfg.FailureThreshold, cb.State())
	}
}

func TestCircuitBreaker_RejectsWhenOpen(t *testing.T) {
	cfg := defaultCBConfig()
	cfg.FailureThreshold = 1
	cfg.ResetTimeout = 10 * time.Minute // keep it open
	cb := newCircuitBreaker(cfg)

	cb.recordFailure()

	if cb.allow() {
		t.Fatal("expected allow() to return false when circuit is open")
	}
}

func TestCircuitBreaker_TransitionsToHalfOpen(t *testing.T) {
	cfg := defaultCBConfig()
	cfg.FailureThreshold = 1
	cfg.ResetTimeout = 1 * time.Millisecond
	cb := newCircuitBreaker(cfg)

	cb.recordFailure()
	if cb.State() != CircuitOpen {
		t.Fatalf("expected open state, got %s", cb.State())
	}

	// Wait for reset timeout to elapse
	time.Sleep(5 * time.Millisecond)

	if !cb.allow() {
		t.Fatal("expected allow() to return true after reset timeout")
	}
	if cb.State() != CircuitHalfOpen {
		t.Fatalf("expected half-open state after timeout, got %s", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenSuccess(t *testing.T) {
	cfg := defaultCBConfig()
	cfg.FailureThreshold = 1
	cfg.ResetTimeout = 1 * time.Millisecond
	cb := newCircuitBreaker(cfg)

	cb.recordFailure()
	time.Sleep(5 * time.Millisecond)
	cb.allow() // transition to half-open

	cb.recordSuccess()

	if cb.State() != CircuitClosed {
		t.Fatalf("expected closed state after success in half-open, got %s", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenFailure(t *testing.T) {
	cfg := defaultCBConfig()
	cfg.FailureThreshold = 1
	cfg.ResetTimeout = 1 * time.Millisecond
	cb := newCircuitBreaker(cfg)

	cb.recordFailure()
	time.Sleep(5 * time.Millisecond)
	cb.allow() // transition to half-open

	cb.recordFailure()

	if cb.State() != CircuitOpen {
		t.Fatalf("expected open state after failure in half-open, got %s", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenLimit(t *testing.T) {
	cfg := defaultCBConfig()
	cfg.FailureThreshold = 1
	cfg.ResetTimeout = 1 * time.Millisecond
	cfg.HalfOpenMax = 2
	cb := newCircuitBreaker(cfg)

	cb.recordFailure()
	time.Sleep(5 * time.Millisecond)

	// First call transitions to half-open and allows
	if !cb.allow() {
		t.Fatal("expected first allow() in half-open to return true")
	}
	// Second call still within limit
	if !cb.allow() {
		t.Fatal("expected second allow() in half-open to return true (within HalfOpenMax)")
	}
	// Third call exceeds limit
	if cb.allow() {
		t.Fatal("expected third allow() to return false (exceeds HalfOpenMax)")
	}
}

func TestCircuitBreaker_SuccessResetsFailCount(t *testing.T) {
	cfg := defaultCBConfig()
	cfg.FailureThreshold = 3
	cb := newCircuitBreaker(cfg)

	cb.recordFailure()
	cb.recordFailure()
	cb.recordSuccess() // should reset consecutive fails

	// Two more failures should not open (need 3 consecutive)
	cb.recordFailure()
	cb.recordFailure()

	if cb.State() != CircuitClosed {
		t.Fatalf("expected closed state (reset should have cleared fail count), got %s", cb.State())
	}

	// One more failure reaches threshold
	cb.recordFailure()
	if cb.State() != CircuitOpen {
		t.Fatalf("expected open state after 3 consecutive failures post-reset, got %s", cb.State())
	}
}

func TestCircuitState_String(t *testing.T) {
	tests := []struct {
		state CircuitState
		want  string
	}{
		{CircuitClosed, "closed"},
		{CircuitOpen, "open"},
		{CircuitHalfOpen, "half-open"},
		{CircuitState(99), "unknown"},
	}

	for _, tt := range tests {
		got := tt.state.String()
		if got != tt.want {
			t.Errorf("CircuitState(%d).String() = %q, want %q", tt.state, got, tt.want)
		}
	}
}

func TestCBOptions(t *testing.T) {
	cfg := defaultCBConfig()

	FailureThreshold(10)(cfg)
	if cfg.FailureThreshold != 10 {
		t.Errorf("FailureThreshold option: got %d, want 10", cfg.FailureThreshold)
	}

	ResetTimeout(30 * time.Second)(cfg)
	if cfg.ResetTimeout != 30*time.Second {
		t.Errorf("ResetTimeout option: got %v, want 30s", cfg.ResetTimeout)
	}

	HalfOpenMax(3)(cfg)
	if cfg.HalfOpenMax != 3 {
		t.Errorf("HalfOpenMax option: got %d, want 3", cfg.HalfOpenMax)
	}
}

func TestErrCircuitOpen(t *testing.T) {
	if ErrCircuitOpen == nil {
		t.Fatal("ErrCircuitOpen should not be nil")
	}
	if ErrCircuitOpen.Error() != "circuit breaker is open" {
		t.Errorf("ErrCircuitOpen.Error() = %q, want %q", ErrCircuitOpen.Error(), "circuit breaker is open")
	}
}

func TestCircuitBreaker_Concurrent(t *testing.T) {
	cfg := defaultCBConfig()
	cfg.FailureThreshold = 5
	cfg.ResetTimeout = 1 * time.Millisecond
	cb := newCircuitBreaker(cfg)

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			cb.allow()
			if i%3 == 0 {
				cb.recordFailure()
			} else {
				cb.recordSuccess()
			}
			_ = cb.State()
		}(i)
	}

	wg.Wait()
	// If we reach here without a race condition panic, concurrent access is safe.
	// Run with -race flag for full detection.
}
