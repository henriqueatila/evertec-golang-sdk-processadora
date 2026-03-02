package client

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestRateLimiterDefaults(t *testing.T) {
	cfg := defaultRateLimiterConfig()
	if cfg.RequestsPerSecond != 10 {
		t.Errorf("expected RPS=10, got %f", cfg.RequestsPerSecond)
	}
	if cfg.BurstSize != 20 {
		t.Errorf("expected BurstSize=20, got %d", cfg.BurstSize)
	}
	if cfg.WaitTimeout != 0 {
		t.Errorf("expected WaitTimeout=0, got %v", cfg.WaitTimeout)
	}
}

func TestRateLimiterOptions(t *testing.T) {
	cfg := defaultRateLimiterConfig()
	RequestsPerSecond(50)(cfg)
	BurstSize(100)(cfg)
	WaitTimeout(5 * time.Second)(cfg)

	if cfg.RequestsPerSecond != 50 {
		t.Errorf("expected RPS=50, got %f", cfg.RequestsPerSecond)
	}
	if cfg.BurstSize != 100 {
		t.Errorf("expected BurstSize=100, got %d", cfg.BurstSize)
	}
	if cfg.WaitTimeout != 5*time.Second {
		t.Errorf("expected WaitTimeout=5s, got %v", cfg.WaitTimeout)
	}
}

func TestRateLimiterAllowsBurst(t *testing.T) {
	rl := newRateLimiter(&RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         5,
	})

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		if err := rl.wait(ctx); err != nil {
			t.Fatalf("request %d should be allowed, got: %v", i, err)
		}
	}

	// 6th request should be rejected (no wait timeout)
	if err := rl.wait(ctx); err != ErrRateLimited {
		t.Errorf("expected ErrRateLimited, got: %v", err)
	}
}

func TestRateLimiterRefills(t *testing.T) {
	rl := newRateLimiter(&RateLimiterConfig{
		RequestsPerSecond: 1000, // fast refill for test
		BurstSize:         1,
	})

	ctx := context.Background()
	// Drain the bucket
	if err := rl.wait(ctx); err != nil {
		t.Fatalf("first request should be allowed: %v", err)
	}

	// Wait for refill
	time.Sleep(5 * time.Millisecond)

	// Should have refilled
	if err := rl.wait(ctx); err != nil {
		t.Errorf("request after refill should be allowed: %v", err)
	}
}

func TestRateLimiterWaitTimeout(t *testing.T) {
	rl := newRateLimiter(&RateLimiterConfig{
		RequestsPerSecond: 1000,
		BurstSize:         1,
		WaitTimeout:       50 * time.Millisecond,
	})

	ctx := context.Background()
	// Drain
	_ = rl.wait(ctx)

	// Should wait and succeed (refill within 50ms at 1000 rps)
	if err := rl.wait(ctx); err != nil {
		t.Errorf("expected success with wait, got: %v", err)
	}
}

func TestRateLimiterWaitTimeoutExpires(t *testing.T) {
	rl := newRateLimiter(&RateLimiterConfig{
		RequestsPerSecond: 0.1, // very slow refill
		BurstSize:         1,
		WaitTimeout:       10 * time.Millisecond,
	})

	ctx := context.Background()
	_ = rl.wait(ctx) // drain

	err := rl.wait(ctx)
	if err != ErrRateLimited {
		t.Errorf("expected ErrRateLimited after timeout, got: %v", err)
	}
}

func TestRateLimiterContextCancellation(t *testing.T) {
	rl := newRateLimiter(&RateLimiterConfig{
		RequestsPerSecond: 0.1,
		BurstSize:         1,
		WaitTimeout:       5 * time.Second,
	})

	_ = rl.wait(context.Background()) // drain

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := rl.wait(ctx)
	if err == nil {
		t.Error("expected error on cancelled context")
	}
}

func TestRateLimiterConcurrentAccess(t *testing.T) {
	rl := newRateLimiter(&RateLimiterConfig{
		RequestsPerSecond: 10000,
		BurstSize:         100,
	})

	var wg sync.WaitGroup
	ctx := context.Background()
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = rl.wait(ctx)
		}()
	}
	wg.Wait()
}

func TestRateLimiterAvailableTokens(t *testing.T) {
	rl := newRateLimiter(&RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         5,
	})

	tokens := rl.availableTokens()
	if tokens != 5 {
		t.Errorf("expected 5 tokens, got %f", tokens)
	}

	_ = rl.wait(context.Background())
	tokens = rl.availableTokens()
	if tokens < 3.9 || tokens > 4.1 {
		t.Errorf("expected ~4 tokens after one use, got %f", tokens)
	}
}
