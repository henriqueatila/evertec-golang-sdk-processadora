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
