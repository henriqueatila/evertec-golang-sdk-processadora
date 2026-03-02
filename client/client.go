// Package client provides the main Evertec Processadora API client.
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/internal/mtls"
	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

const (
	// BaseURLHomolog is the homologation environment URL.
	BaseURLHomolog = "https://api-hml.paysmart.com.br"
	// BaseURLProduction is the production environment URL.
	BaseURLProduction = "https://api.paysmart.com.br"
	// BasePath is the API base path.
	BasePath = "/paySmart/ps-processadora/v1"

	defaultTimeout         = 30 * time.Second
	defaultMaxResponseBody = 10 * 1024 * 1024 // 10 MB
)

// idempotencyKeyCtxKey is the context key for custom idempotency keys.
type idempotencyKeyCtxKey struct{}

// WithIdempotencyKey returns a context with a custom idempotency key.
// Use this to ensure retries of the same request use the same key.
// If not set, the client auto-generates a UUID-v4 for each mutation request.
//
// Example:
//
//	ctx := client.WithIdempotencyKey(ctx, "my-unique-key-123")
//	err := client.CreateAccount(ctx, req)
//	// On retry, use the same context to ensure idempotent behavior
func WithIdempotencyKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, idempotencyKeyCtxKey{}, key)
}

// getIdempotencyKey retrieves the idempotency key from context, if set.
func getIdempotencyKey(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(idempotencyKeyCtxKey{}).(string)
	return key, ok && key != ""
}

// Client is the Evertec Processadora API client.
type Client struct {
	httpClient      *http.Client
	baseURL         string
	apiKey          string
	userAgent       string
	hooks           []Hook
	logger          *slog.Logger
	maxResponseBody int64
	retry           *RetryConfig
	cb              *circuitBreaker
	certRotator     *mtls.CertRotator
	noValidation    bool
	rl              *rateLimiter
}

// Config holds client configuration.
type Config struct {
	// BaseURL is the API base URL (defaults to production)
	BaseURL string
	// APIKey is the API key provided by Evertec
	APIKey string
	// UserAgent is the user agent string (format: {emissor}/{version})
	UserAgent string
	// TLSConfig is the mTLS configuration
	TLSConfig *tls.Config
	// Timeout is the HTTP client timeout (defaults to 30s)
	Timeout time.Duration
	// Logger is the slog logger for request logging (optional)
	Logger *slog.Logger
	// Hooks are observability hooks for metrics, tracing, etc. (optional)
	Hooks []Hook
	// MaxResponseBody is the maximum response body size in bytes (defaults to 10 MB)
	MaxResponseBody int64
	// Retry is the retry configuration for failed requests (optional)
	Retry *RetryConfig
	// CircuitBreaker enables the circuit breaker pattern (optional)
	CircuitBreaker *CircuitBreakerConfig
	// CertRotation enables automatic certificate hot-reload on each TLS handshake
	CertRotation bool
	// NoValidation disables automatic request body validation (default: false)
	NoValidation bool
	// RateLimiter is the rate limiter configuration (optional)
	RateLimiter *RateLimiterConfig
}

// New creates a new Evertec client.
func New(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	if cfg.UserAgent == "" {
		return nil, fmt.Errorf("user agent is required")
	}
	if cfg.TLSConfig == nil {
		return nil, fmt.Errorf("TLS config is required for mTLS")
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = BaseURLProduction
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	maxResponseBody := cfg.MaxResponseBody
	if maxResponseBody <= 0 {
		maxResponseBody = defaultMaxResponseBody
	}

	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig:     cfg.TLSConfig,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	client := &Client{
		httpClient:      httpClient,
		baseURL:         baseURL + BasePath,
		apiKey:          cfg.APIKey,
		userAgent:       cfg.UserAgent,
		hooks:           cfg.Hooks,
		logger:          cfg.Logger,
		maxResponseBody: maxResponseBody,
		retry:           cfg.Retry,
		noValidation:    cfg.NoValidation,
	}

	if cfg.CircuitBreaker != nil {
		client.cb = newCircuitBreaker(cfg.CircuitBreaker)
	}
	if cfg.RateLimiter != nil {
		client.rl = newRateLimiter(cfg.RateLimiter)
	}

	return client, nil
}

// NewWithCertFiles creates a new client loading certificates from files.
// When WithCertRotation() option is provided, the client re-reads the certificate
// from disk on each TLS handshake, enabling hot-reload without restart.
func NewWithCertFiles(apiKey, userAgent, certFile, keyFile string, opts ...Option) (*Client, error) {
	cfg := Config{
		APIKey:    apiKey,
		UserAgent: userAgent,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.CertRotation {
		rotator, err := mtls.NewCertRotator(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create cert rotator: %w", err)
		}

		cfg.TLSConfig = &tls.Config{
			GetClientCertificate: rotator.GetClientCertificate,
			MinVersion:           tls.VersionTLS12,
		}

		c, err := New(cfg)
		if err != nil {
			return nil, err
		}
		c.certRotator = rotator
		return c, nil
	}

	tlsConfig, err := mtls.LoadTLSConfig(mtls.Config{
		CertFile: certFile,
		KeyFile:  keyFile,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS config: %w", err)
	}

	cfg.TLSConfig = tlsConfig
	return New(cfg)
}

// RefreshCertificates forces a reload of the mTLS certificates from disk.
// Only works when certificate rotation is enabled via WithCertRotation().
func (c *Client) RefreshCertificates() error {
	if c.certRotator == nil {
		return fmt.Errorf("certificate rotation not enabled")
	}
	return c.certRotator.Refresh()
}

// Option is a functional option for configuring the client.
type Option func(*Config)

// WithBaseURL sets a custom base URL.
func WithBaseURL(url string) Option {
	return func(c *Config) {
		c.BaseURL = url
	}
}

// WithTimeout sets a custom timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithHomolog configures the client for the homologation environment.
func WithHomolog() Option {
	return func(c *Config) {
		c.BaseURL = BaseURLHomolog
	}
}

// request performs an HTTP request, wrapping doRequest with retry and circuit breaker.
func (c *Client) request(ctx context.Context, method, path string, body any, result any) error {
	// Auto-validate request body if it implements Validatable
	if !c.noValidation && body != nil {
		if v, ok := body.(Validatable); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}
	}

	// Rate limiter check
	if c.rl != nil {
		if err := c.rl.wait(ctx); err != nil {
			return err
		}
	}

	// Circuit breaker check
	if c.cb != nil {
		if !c.cb.allow() {
			return ErrCircuitOpen
		}
	}

	// If no retry configured, do a single request (with CB recording)
	if c.retry == nil {
		err := c.doRequest(ctx, method, path, body, result)
		if c.cb != nil {
			c.recordCBResult(err)
		}
		return err
	}

	// Mutations always carry an idempotency key so they are safe to retry
	canRetry := isRetryableMethod(method, true)

	var lastErr error
	for attempt := 0; attempt <= c.retry.MaxRetries; attempt++ {
		if attempt > 0 {
			// Re-check circuit breaker before each retry
			if c.cb != nil {
				if !c.cb.allow() {
					return ErrCircuitOpen
				}
			}

			delay := calculateBackoff(attempt-1, c.retry.InitialDelay, c.retry.MaxDelay)

			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during retry: %w", ctx.Err())
			case <-time.After(delay):
			}
		}

		err := c.doRequest(ctx, method, path, body, result)
		if c.cb != nil {
			c.recordCBResult(err)
		}
		if err == nil {
			return nil
		}

		lastErr = err

		if !canRetry {
			return err
		}

		var apiErr *APIError
		if errors.As(err, &apiErr) && c.retry.isRetryable(apiErr.StatusCode) {
			continue
		}

		// Not a retryable error
		return err
	}

	return lastErr
}

// recordCBResult records success or failure on the circuit breaker.
func (c *Client) recordCBResult(err error) {
	if err == nil {
		c.cb.recordSuccess()
		return
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode >= 500 {
		c.cb.recordFailure()
	} else {
		c.cb.recordSuccess()
	}
}

// doRequest performs a single HTTP request with required headers.
func (c *Client) doRequest(ctx context.Context, method, path string, body any, result any) error {
	start := time.Now()

	// Check for context cancellation early to avoid unnecessary work
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error before request: %w", err)
	}

	// Prepare request body
	var bodyReader io.Reader
	var bodySize int
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
		bodySize = len(bodyBytes)
	}

	// Get or generate idempotency key for mutations
	var idempotencyKey string
	if method != http.MethodGet {
		if key, ok := getIdempotencyKey(ctx); ok {
			idempotencyKey = key
		} else {
			idempotencyKey = uuid.New().String()
		}
	}

	// Prepare request info for hooks
	reqInfo := &RequestInfo{
		Method:         method,
		Path:           path,
		BodySize:       bodySize,
		IdempotencyKey: idempotencyKey,
	}

	// Execute before hooks with panic recovery
	for _, hook := range c.hooks {
		func() {
			defer func() {
				if r := recover(); r != nil && c.logger != nil {
					c.logger.Error("hook panic recovered",
						slog.Any("panic", r),
						slog.String("hook_phase", "BeforeRequest"),
						slog.String("path", path))
				}
			}()
			ctx = hook.BeforeRequest(ctx, reqInfo)
		}()
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		c.executeAfterHooks(ctx, reqInfo, &ResponseInfo{Error: err, Duration: time.Since(start)})
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("API-Key", c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add IdempotencyKey for mutations
	if idempotencyKey != "" {
		req.Header.Set("IdempotencyKey", idempotencyKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.executeAfterHooks(ctx, reqInfo, &ResponseInfo{Error: err, Duration: time.Since(start)})
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, c.maxResponseBody))
	if err != nil {
		c.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: resp.StatusCode,
			Error:      err,
			Duration:   time.Since(start),
		})
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var apiErr error
	if resp.StatusCode >= 400 {
		var errResp types.Error
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			apiErr = &APIError{
				StatusCode: resp.StatusCode,
				Message:    string(respBody),
			}
		} else {
			apiErr = &APIError{
				StatusCode: resp.StatusCode,
				Code:       errResp.Code,
				Message:    errResp.Message,
				Details:    errResp.Details,
			}
		}
	}

	// Prepare response info
	respInfo := &ResponseInfo{
		StatusCode: resp.StatusCode,
		BodySize:   len(respBody),
		Duration:   time.Since(start),
		Error:      apiErr,
	}

	// Execute after hooks and logging
	c.executeAfterHooks(ctx, reqInfo, respInfo)

	if apiErr != nil {
		return apiErr
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// executeAfterHooks runs all after hooks and logging with panic recovery.
func (c *Client) executeAfterHooks(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
	// Execute hooks with panic recovery
	for _, hook := range c.hooks {
		func() {
			defer func() {
				if r := recover(); r != nil && c.logger != nil {
					c.logger.Error("hook panic recovered",
						slog.Any("panic", r),
						slog.String("hook_phase", "AfterResponse"),
						slog.String("path", req.Path))
				}
			}()
			hook.AfterResponse(ctx, req, resp)
		}()
	}

	// Execute logging
	c.logRequest(ctx, req, resp)
}

// APIError represents an API error.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Details    []string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("API error %d: [%s] %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// Is implements errors.Is matching by error code.
// This allows using errors.Is(err, client.ErrAccountNotFound) for error checking.
// Matching priority: Code (if target has Code set) > StatusCode (if target has no Code).
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	if !ok {
		return false
	}

	// If target has Code, match by Code only
	if t.Code != "" {
		return e.Code == t.Code
	}

	// If target has no Code but has StatusCode, match by StatusCode
	// only if the error also has no Code (pure status matching)
	if t.StatusCode != 0 && e.Code == "" {
		return e.StatusCode == t.StatusCode
	}

	return false
}

// Sentinel errors for common API error codes.
// Use with errors.Is(err, client.ErrAccountNotFound) for type-safe error checking.
var (
	// ErrAccountNotFound indicates the account was not found.
	ErrAccountNotFound = &APIError{Code: "ACCOUNT_NOT_FOUND"}
	// ErrAccountAlreadyExists indicates the account already exists.
	ErrAccountAlreadyExists = &APIError{Code: "ACCOUNT_ALREADY_EXISTS"}
	// ErrAccountBlocked indicates the account is blocked.
	ErrAccountBlocked = &APIError{Code: "ACCOUNT_BLOCKED"}

	// ErrCardNotFound indicates the card was not found.
	ErrCardNotFound = &APIError{Code: "CARD_NOT_FOUND"}
	// ErrCardBlocked indicates the card is blocked.
	ErrCardBlocked = &APIError{Code: "CARD_BLOCKED"}
	// ErrCardExpired indicates the card has expired.
	ErrCardExpired = &APIError{Code: "CARD_EXPIRED"}
	// ErrCardCancelled indicates the card is cancelled.
	ErrCardCancelled = &APIError{Code: "CARD_CANCELLED"}

	// ErrInvalidDocument indicates an invalid CPF/CNPJ document.
	ErrInvalidDocument = &APIError{Code: "INVALID_DOCUMENT"}
	// ErrInvalidRequest indicates a malformed request.
	ErrInvalidRequest = &APIError{Code: "INVALID_REQUEST"}
	// ErrInvalidPIN indicates an invalid PIN.
	ErrInvalidPIN = &APIError{Code: "INVALID_PIN"}

	// ErrInsufficientBalance indicates insufficient account balance.
	ErrInsufficientBalance = &APIError{Code: "INSUFFICIENT_BALANCE"}
	// ErrInsufficientLimit indicates insufficient credit limit.
	ErrInsufficientLimit = &APIError{Code: "INSUFFICIENT_LIMIT"}

	// ErrUnauthorized indicates authentication failure (401).
	ErrUnauthorized = &APIError{StatusCode: 401}
	// ErrForbidden indicates authorization failure (403).
	ErrForbidden = &APIError{StatusCode: 403}
	// ErrNotFound indicates resource not found (404).
	ErrNotFound = &APIError{StatusCode: 404}
	// ErrConflict indicates a conflict (409).
	ErrConflict = &APIError{StatusCode: 409}
	// ErrTooManyRequests indicates rate limiting (429).
	ErrTooManyRequests = &APIError{StatusCode: 429}
	// ErrInternalServer indicates a server error (500).
	ErrInternalServer = &APIError{StatusCode: 500}
)
