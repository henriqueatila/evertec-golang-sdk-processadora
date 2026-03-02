package client

import "log/slog"

// WithLogger sets a custom slog logger for request logging.
// The logger will receive structured logs for all API requests.
//
// Example:
//
//	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
//	client, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
//	    client.WithLogger(logger),
//	)
func WithLogger(logger *slog.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

// WithDefaultLogger enables logging using slog.Default().
//
// Example:
//
//	client, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
//	    client.WithDefaultLogger(),
//	)
func WithDefaultLogger() Option {
	return func(c *Config) {
		c.Logger = slog.Default()
	}
}

// WithMaxResponseBody sets the maximum response body size in bytes.
// Responses exceeding this limit will be truncated. Default is 10 MB.
func WithMaxResponseBody(size int64) Option {
	return func(c *Config) {
		c.MaxResponseBody = size
	}
}

// WithHooks adds observability hooks to the client.
// Hooks are called before and after each request, useful for
// metrics, tracing, or custom logging.
//
// Example:
//
//	client, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
//	    client.WithHooks(
//	        myMetricsHook,
//	        myTracingHook,
//	    ),
//	)
func WithHooks(hooks ...Hook) Option {
	return func(c *Config) {
		c.Hooks = append(c.Hooks, hooks...)
	}
}

// WithRetry enables retry with exponential backoff for transient errors.
// By default retries up to 3 times on status codes 429, 500, 502, 503, 504.
//
// Example:
//
//	client, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
//	    client.WithRetry(
//	        client.MaxRetries(5),
//	        client.InitialDelay(200*time.Millisecond),
//	    ),
//	)
func WithRetry(opts ...RetryOption) Option {
	return func(c *Config) {
		c.Retry = defaultRetryConfig()
		for _, opt := range opts {
			opt(c.Retry)
		}
	}
}

// WithCircuitBreaker enables the circuit breaker pattern.
// It protects against cascading failures by stopping requests when
// the upstream service is consistently returning 5xx errors.
//
// Example:
//
//	client, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
//	    client.WithCircuitBreaker(
//	        client.FailureThreshold(3),
//	        client.ResetTimeout(30*time.Second),
//	    ),
//	)
func WithCircuitBreaker(opts ...CBOption) Option {
	return func(c *Config) {
		c.CircuitBreaker = defaultCBConfig()
		for _, opt := range opts {
			opt(c.CircuitBreaker)
		}
	}
}

// WithRateLimiter enables client-side rate limiting using a token bucket.
// Prevents flooding the API during retry storms or high-throughput scenarios.
//
// Example:
//
//	client, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
//	    client.WithRateLimiter(
//	        client.RequestsPerSecond(50),
//	        client.BurstSize(100),
//	    ),
//	)
func WithRateLimiter(opts ...RateLimiterOption) Option {
	return func(c *Config) {
		c.RateLimiter = defaultRateLimiterConfig()
		for _, opt := range opts {
			opt(c.RateLimiter)
		}
	}
}

// WithCertRotation enables automatic certificate hot-reload.
// When enabled, the client re-reads the certificate from disk on each TLS handshake,
// allowing certificate renewal without restarting the process.
// Use client.RefreshCertificates() to force a manual reload.
func WithCertRotation() Option {
	return func(c *Config) {
		c.CertRotation = true
	}
}

// WithNoValidation disables automatic request body validation.
// By default, request types implementing the Validatable interface
// are validated before sending. Use this to skip validation.
func WithNoValidation() Option {
	return func(c *Config) {
		c.NoValidation = true
	}
}
