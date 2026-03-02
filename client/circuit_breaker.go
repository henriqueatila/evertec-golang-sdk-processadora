package client

import (
	"errors"
	"sync"
	"time"
)

// ErrCircuitOpen is returned when the circuit breaker is open and rejecting requests.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitState represents the state of the circuit breaker.
type CircuitState int

const (
	CircuitClosed   CircuitState = iota // Normal operation
	CircuitOpen                         // Rejecting requests
	CircuitHalfOpen                     // Testing with limited requests
)

// String returns the string representation of a circuit state.
func (s CircuitState) String() string {
	switch s {
	case CircuitClosed:
		return "closed"
	case CircuitOpen:
		return "open"
	case CircuitHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreakerConfig holds circuit breaker configuration.
type CircuitBreakerConfig struct {
	FailureThreshold int           // Consecutive failures before opening (default: 5)
	ResetTimeout     time.Duration // Time to wait before transitioning to half-open (default: 60s)
	HalfOpenMax      int           // Max requests in half-open state (default: 1)
}

// CBOption configures circuit breaker behavior.
type CBOption func(*CircuitBreakerConfig)

func defaultCBConfig() *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		FailureThreshold: 5,
		ResetTimeout:     60 * time.Second,
		HalfOpenMax:      1,
	}
}

// FailureThreshold sets the number of consecutive failures before opening the circuit.
func FailureThreshold(n int) CBOption {
	return func(c *CircuitBreakerConfig) {
		c.FailureThreshold = n
	}
}

// ResetTimeout sets how long the circuit stays open before transitioning to half-open.
func ResetTimeout(d time.Duration) CBOption {
	return func(c *CircuitBreakerConfig) {
		c.ResetTimeout = d
	}
}

// HalfOpenMax sets the maximum number of requests allowed in half-open state.
func HalfOpenMax(n int) CBOption {
	return func(c *CircuitBreakerConfig) {
		c.HalfOpenMax = n
	}
}

// circuitBreaker implements the circuit breaker pattern.
type circuitBreaker struct {
	mu               sync.Mutex
	config           *CircuitBreakerConfig
	state            CircuitState
	consecutiveFails int
	lastFailureTime  time.Time
	halfOpenCount    int
	// metrics counters
	totalRequests    int64
	totalSuccesses   int64
	totalFailures    int64
	totalRejected    int64
	stateTransitions int64
}

func newCircuitBreaker(config *CircuitBreakerConfig) *circuitBreaker {
	return &circuitBreaker{
		config: config,
		state:  CircuitClosed,
	}
}

// allow checks if a request is allowed through the circuit breaker.
// Returns true if the request can proceed, false if the circuit is open.
func (cb *circuitBreaker) allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.totalRequests++

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		// Check if reset timeout has elapsed
		if time.Since(cb.lastFailureTime) >= cb.config.ResetTimeout {
			cb.state = CircuitHalfOpen
			cb.stateTransitions++
			cb.halfOpenCount = 1 // count this probe request
			return true
		}
		cb.totalRejected++
		return false
	case CircuitHalfOpen:
		if cb.halfOpenCount < cb.config.HalfOpenMax {
			cb.halfOpenCount++
			return true
		}
		cb.totalRejected++
		return false
	default:
		cb.totalRejected++
		return false
	}
}

// recordSuccess records a successful request.
func (cb *circuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.totalSuccesses++
	cb.consecutiveFails = 0
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
		cb.stateTransitions++
	}
}

// recordFailure records a failed request (5xx status).
func (cb *circuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.totalFailures++
	cb.consecutiveFails++
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case CircuitClosed:
		if cb.consecutiveFails >= cb.config.FailureThreshold {
			cb.state = CircuitOpen
			cb.stateTransitions++
		}
	case CircuitHalfOpen:
		cb.state = CircuitOpen
		cb.stateTransitions++
	}
}

// State returns the current circuit breaker state.
func (cb *circuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// metrics returns a snapshot of circuit breaker counters.
func (cb *circuitBreaker) metrics() CBMetrics {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return CBMetrics{
		TotalRequests:    cb.totalRequests,
		TotalSuccesses:   cb.totalSuccesses,
		TotalFailures:    cb.totalFailures,
		TotalRejected:    cb.totalRejected,
		StateTransitions: cb.stateTransitions,
	}
}

// health returns the circuit breaker health snapshot.
func (cb *circuitBreaker) health() CBHealth {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return CBHealth{
		State:            cb.state.String(),
		ConsecutiveFails: cb.consecutiveFails,
		ResetTimeout:     cb.config.ResetTimeout,
		Metrics: CBMetrics{
			TotalRequests:    cb.totalRequests,
			TotalSuccesses:   cb.totalSuccesses,
			TotalFailures:    cb.totalFailures,
			TotalRejected:    cb.totalRejected,
			StateTransitions: cb.stateTransitions,
		},
	}
}
