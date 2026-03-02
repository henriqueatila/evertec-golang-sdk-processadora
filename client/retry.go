package client

import (
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// RetryConfig holds retry configuration.
type RetryConfig struct {
	MaxRetries        int
	InitialDelay      time.Duration
	MaxDelay          time.Duration
	RetryableStatuses map[int]bool
}

// RetryOption configures retry behavior.
type RetryOption func(*RetryConfig)

// defaultRetryConfig returns sensible defaults.
func defaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:   3,
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     30 * time.Second,
		RetryableStatuses: map[int]bool{
			429: true,
			500: true,
			502: true,
			503: true,
			504: true,
		},
	}
}

// MaxRetries sets the maximum number of retry attempts.
func MaxRetries(n int) RetryOption {
	return func(c *RetryConfig) {
		c.MaxRetries = n
	}
}

// InitialDelay sets the initial delay before the first retry.
func InitialDelay(d time.Duration) RetryOption {
	return func(c *RetryConfig) {
		c.InitialDelay = d
	}
}

// MaxDelay sets the maximum delay between retries.
func MaxDelay(d time.Duration) RetryOption {
	return func(c *RetryConfig) {
		c.MaxDelay = d
	}
}

// RetryableStatuses sets which HTTP status codes trigger a retry.
func RetryableStatuses(codes ...int) RetryOption {
	return func(c *RetryConfig) {
		c.RetryableStatuses = make(map[int]bool, len(codes))
		for _, code := range codes {
			c.RetryableStatuses[code] = true
		}
	}
}

// isRetryable checks if a status code should trigger a retry.
func (rc *RetryConfig) isRetryable(statusCode int) bool {
	return rc.RetryableStatuses[statusCode]
}

// isRetryableMethod checks if the HTTP method is safe to retry.
// GET is always retryable. POST/PUT/PATCH are only retryable with an idempotency key.
func isRetryableMethod(method string, hasIdempotencyKey bool) bool {
	if method == http.MethodGet {
		return true
	}
	return hasIdempotencyKey
}

// calculateBackoff computes the delay for a given retry attempt with jitter.
// Uses exponential backoff: initialDelay * 2^attempt, capped at maxDelay.
// Jitter range: 0.5x to 1.5x of the computed delay.
func calculateBackoff(attempt int, initialDelay, maxDelay time.Duration) time.Duration {
	delay := float64(initialDelay) * math.Pow(2, float64(attempt))
	if delay > float64(maxDelay) {
		delay = float64(maxDelay)
	}

	// Apply jitter: 0.5x to 1.5x
	jitter := 0.5 + rand.Float64() // [0.5, 1.5)
	delay *= jitter

	return time.Duration(delay)
}

// parseRetryAfter extracts the delay from a Retry-After header.
// Supports seconds format only (not HTTP-date).
// Returns 0 if the header is missing or unparseable.
func parseRetryAfter(resp *http.Response) time.Duration {
	if resp == nil {
		return 0
	}
	header := resp.Header.Get("Retry-After")
	if header == "" {
		return 0
	}
	seconds, err := strconv.Atoi(header)
	if err != nil || seconds <= 0 {
		return 0
	}
	return time.Duration(seconds) * time.Second
}
