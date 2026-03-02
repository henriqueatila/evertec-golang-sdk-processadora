package client

import "time"

// HealthStatus represents the overall health of the client.
type HealthStatus struct {
	// Healthy is true when the client can accept requests.
	Healthy bool `json:"healthy"`
	// CircuitBreaker is the circuit breaker status (nil if not enabled).
	CircuitBreaker *CBHealth `json:"circuit_breaker,omitempty"`
	// RateLimiter is the rate limiter status (nil if not enabled).
	RateLimiter *RateLimiterHealth `json:"rate_limiter,omitempty"`
}

// CBHealth represents circuit breaker health information.
type CBHealth struct {
	State            string        `json:"state"`
	ConsecutiveFails int           `json:"consecutive_fails"`
	ResetTimeout     time.Duration `json:"reset_timeout"`
	Metrics          CBMetrics     `json:"metrics"`
}

// CBMetrics holds circuit breaker counters.
type CBMetrics struct {
	TotalRequests    int64 `json:"total_requests"`
	TotalSuccesses   int64 `json:"total_successes"`
	TotalFailures    int64 `json:"total_failures"`
	TotalRejected    int64 `json:"total_rejected"`
	StateTransitions int64 `json:"state_transitions"`
}

// RateLimiterHealth represents rate limiter health information.
type RateLimiterHealth struct {
	AvailableTokens float64 `json:"available_tokens"`
	BurstSize       int     `json:"burst_size"`
	RPS             float64 `json:"requests_per_second"`
}

// Health returns the current health status of the client.
func (c *Client) Health() HealthStatus {
	status := HealthStatus{Healthy: true}

	if c.cb != nil {
		cbHealth := c.cb.health()
		status.CircuitBreaker = &cbHealth
		if cbHealth.State == "open" {
			status.Healthy = false
		}
	}

	if c.rl != nil {
		status.RateLimiter = &RateLimiterHealth{
			AvailableTokens: c.rl.availableTokens(),
			BurstSize:       c.rl.config.BurstSize,
			RPS:             c.rl.config.RequestsPerSecond,
		}
	}

	return status
}

// CircuitBreakerMetrics returns circuit breaker metrics.
// Returns nil if circuit breaker is not enabled.
func (c *Client) CircuitBreakerMetrics() *CBMetrics {
	if c.cb == nil {
		return nil
	}
	m := c.cb.metrics()
	return &m
}
