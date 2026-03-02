package client

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ErrRateLimited is returned when the rate limiter rejects a request.
var ErrRateLimited = fmt.Errorf("rate limited: too many requests")

// RateLimiterConfig holds rate limiter configuration.
type RateLimiterConfig struct {
	// RequestsPerSecond is the sustained rate limit (default: 10)
	RequestsPerSecond float64
	// BurstSize is the maximum burst above the sustained rate (default: 20)
	BurstSize int
	// WaitTimeout is max time to wait for a token before returning ErrRateLimited.
	// Zero means never wait (reject immediately). Default: 0.
	WaitTimeout time.Duration
}

// RateLimiterOption configures rate limiter behavior.
type RateLimiterOption func(*RateLimiterConfig)

func defaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         20,
		WaitTimeout:       0,
	}
}

// RequestsPerSecond sets the sustained request rate.
func RequestsPerSecond(rps float64) RateLimiterOption {
	return func(c *RateLimiterConfig) {
		c.RequestsPerSecond = rps
	}
}

// BurstSize sets the maximum burst size.
func BurstSize(n int) RateLimiterOption {
	return func(c *RateLimiterConfig) {
		c.BurstSize = n
	}
}

// WaitTimeout sets the maximum time to wait for a token.
// Zero means reject immediately when no tokens available.
func WaitTimeout(d time.Duration) RateLimiterOption {
	return func(c *RateLimiterConfig) {
		c.WaitTimeout = d
	}
}

// rateLimiter implements a token bucket rate limiter.
type rateLimiter struct {
	mu        sync.Mutex
	config    *RateLimiterConfig
	tokens    float64
	lastRefil time.Time
}

func newRateLimiter(config *RateLimiterConfig) *rateLimiter {
	return &rateLimiter{
		config:    config,
		tokens:    float64(config.BurstSize),
		lastRefil: time.Now(),
	}
}

// wait blocks until a token is available or ctx/timeout expires.
// Returns nil if a token was acquired, ErrRateLimited otherwise.
func (rl *rateLimiter) wait(ctx context.Context) error {
	// Fast path: try to acquire immediately
	if rl.tryAcquire() {
		return nil
	}

	if rl.config.WaitTimeout == 0 {
		return ErrRateLimited
	}

	// Slow path: wait for a token
	deadline := time.Now().Add(rl.config.WaitTimeout)
	ticker := time.NewTicker(time.Duration(float64(time.Second) / rl.config.RequestsPerSecond / 2))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled waiting for rate limit: %w", ctx.Err())
		case <-ticker.C:
			if time.Now().After(deadline) {
				return ErrRateLimited
			}
			if rl.tryAcquire() {
				return nil
			}
		}
	}
}

// tryAcquire attempts to take a token. Returns true if successful.
func (rl *rateLimiter) tryAcquire() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()

	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

// refill adds tokens based on elapsed time. Must be called under lock.
func (rl *rateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefil).Seconds()
	rl.tokens += elapsed * rl.config.RequestsPerSecond

	max := float64(rl.config.BurstSize)
	if rl.tokens > max {
		rl.tokens = max
	}
	rl.lastRefil = now
}

// availableTokens returns the current token count (for metrics).
func (rl *rateLimiter) availableTokens() float64 {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.refill()
	return rl.tokens
}
