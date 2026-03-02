package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// ============================================================================
// CLIENT CREATION TESTS
// ============================================================================

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		config         Config
		wantErr        bool
		wantErrContain string
	}{
		{
			name: "success with all required fields",
			config: Config{
				APIKey:    "test-api-key",
				UserAgent: "TestApp/1.0",
				BaseURL:   "https://api.example.com",
				TLSConfig: &tls.Config{},
			},
			wantErr: false,
		},
		{
			name: "success with defaults (no BaseURL)",
			config: Config{
				APIKey:    "test-api-key",
				UserAgent: "TestApp/1.0",
				TLSConfig: &tls.Config{},
			},
			wantErr: false,
		},
		{
			name: "success with custom timeout",
			config: Config{
				APIKey:    "test-api-key",
				UserAgent: "TestApp/1.0",
				TLSConfig: &tls.Config{},
				Timeout:   60 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "error: missing API key",
			config: Config{
				UserAgent: "TestApp/1.0",
				TLSConfig: &tls.Config{},
			},
			wantErr:        true,
			wantErrContain: "API key is required",
		},
		{
			name: "error: empty API key",
			config: Config{
				APIKey:    "",
				UserAgent: "TestApp/1.0",
				TLSConfig: &tls.Config{},
			},
			wantErr:        true,
			wantErrContain: "API key is required",
		},
		{
			name: "error: missing User-Agent",
			config: Config{
				APIKey:    "test-api-key",
				TLSConfig: &tls.Config{},
			},
			wantErr:        true,
			wantErrContain: "user agent is required",
		},
		{
			name: "error: empty User-Agent",
			config: Config{
				APIKey:    "test-api-key",
				UserAgent: "",
				TLSConfig: &tls.Config{},
			},
			wantErr:        true,
			wantErrContain: "user agent is required",
		},
		{
			name: "error: missing TLS config",
			config: Config{
				APIKey:    "test-api-key",
				UserAgent: "TestApp/1.0",
			},
			wantErr:        true,
			wantErrContain: "TLS config is required",
		},
		{
			name: "error: nil TLS config",
			config: Config{
				APIKey:    "test-api-key",
				UserAgent: "TestApp/1.0",
				TLSConfig: nil,
			},
			wantErr:        true,
			wantErrContain: "TLS config is required",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := New(tt.config)

			if tt.wantErr {
				AssertError(t, err, "New()")
				if tt.wantErrContain != "" {
					AssertErrorContains(t, err, tt.wantErrContain)
				}
				AssertNil(t, client, "client")
			} else {
				AssertNoError(t, err, "New()")
				AssertNotNil(t, client, "client")
			}
		})
	}
}

func TestNew_DefaultValues(t *testing.T) {
	t.Parallel()

	config := Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		TLSConfig: &tls.Config{},
	}

	client, err := New(config)
	AssertNoError(t, err, "New()")
	AssertNotNil(t, client, "client")

	// Verify default BaseURL is set (production)
	if !strings.Contains(client.baseURL, BasePath) {
		t.Errorf("baseURL should contain BasePath %q, got %q", BasePath, client.baseURL)
	}
}

func TestNewWithCertFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		apiKey         string
		userAgent      string
		certFile       string
		keyFile        string
		opts           []Option
		wantErr        bool
		wantErrContain string
	}{
		{
			name:           "error: missing cert file",
			apiKey:         "test-key",
			userAgent:      "test/1.0",
			certFile:       "/nonexistent/cert.pem",
			keyFile:        "/nonexistent/key.pem",
			wantErr:        true,
			wantErrContain: "failed to load TLS config",
		},
		{
			name:           "error: empty API key",
			apiKey:         "",
			userAgent:      "test/1.0",
			certFile:       "/nonexistent/cert.pem",
			keyFile:        "/nonexistent/key.pem",
			wantErr:        true,
			wantErrContain: "",
		},
		{
			name:           "error: empty user agent",
			apiKey:         "test-key",
			userAgent:      "",
			certFile:       "/nonexistent/cert.pem",
			keyFile:        "/nonexistent/key.pem",
			wantErr:        true,
			wantErrContain: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := NewWithCertFiles(tt.apiKey, tt.userAgent, tt.certFile, tt.keyFile, tt.opts...)

			if tt.wantErr {
				AssertError(t, err, "NewWithCertFiles")
				if tt.wantErrContain != "" {
					AssertErrorContains(t, err, tt.wantErrContain)
				}
				AssertNil(t, client, "client")
			} else {
				AssertNoError(t, err, "NewWithCertFiles")
				AssertNotNil(t, client, "client")
			}
		})
	}
}

func TestNewWithCertFiles_WithOptions(t *testing.T) {
	t.Parallel()

	// Test that options are applied even when cert files don't exist
	// (error happens first due to cert loading)
	_, err := NewWithCertFiles(
		"test-key",
		"test/1.0",
		"/nonexistent/cert.pem",
		"/nonexistent/key.pem",
		WithBaseURL("https://custom.api.com"),
		WithTimeout(60*time.Second),
		WithHomolog(),
	)

	// Should fail due to cert files not existing
	AssertError(t, err, "NewWithCertFiles with options")
	AssertErrorContains(t, err, "failed to load TLS config")
}

// ============================================================================
// FUNCTIONAL OPTIONS TESTS
// ============================================================================

func TestWithBaseURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "custom URL",
			url:      "https://custom.api.com",
			expected: "https://custom.api.com",
		},
		{
			name:     "empty URL",
			url:      "",
			expected: "",
		},
		{
			name:     "URL with trailing slash",
			url:      "https://api.example.com/",
			expected: "https://api.example.com/",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &Config{}
			WithBaseURL(tt.url)(cfg)

			AssertEqual(t, tt.expected, cfg.BaseURL, "BaseURL")
		})
	}
}

func TestWithHomolog(t *testing.T) {
	t.Parallel()

	cfg := &Config{}
	WithHomolog()(cfg)

	AssertEqual(t, BaseURLHomolog, cfg.BaseURL, "BaseURL")
}

func TestWithTimeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		timeout  time.Duration
		expected time.Duration
	}{
		{
			name:     "10 seconds",
			timeout:  10 * time.Second,
			expected: 10 * time.Second,
		},
		{
			name:     "1 minute",
			timeout:  time.Minute,
			expected: time.Minute,
		},
		{
			name:     "zero timeout",
			timeout:  0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &Config{}
			WithTimeout(tt.timeout)(cfg)

			AssertEqual(t, tt.expected, cfg.Timeout, "Timeout")
		})
	}
}

// ============================================================================
// REQUEST HEADERS TESTS
// ============================================================================

func TestClient_RequiredHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		method      string
		checkHeader string
		expected    string
		shouldExist bool
	}{
		{
			name:        "API-Key header present",
			method:      http.MethodGet,
			checkHeader: "API-Key",
			expected:    "test-api-key-12345",
			shouldExist: true,
		},
		{
			name:        "User-Agent header present",
			method:      http.MethodGet,
			checkHeader: "User-Agent",
			expected:    "TestEmissor/1.0.0",
			shouldExist: true,
		},
		{
			name:        "Content-Type for POST",
			method:      http.MethodPost,
			checkHeader: "Content-Type",
			expected:    "application/json",
			shouldExist: true,
		},
		{
			name:        "Accept header present",
			method:      http.MethodGet,
			checkHeader: "Accept",
			expected:    "application/json",
			shouldExist: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var receivedHeader string

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedHeader = r.Header.Get(tt.checkHeader)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]any{})
			}))
			defer server.Close()

			client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})

			switch tt.method {
			case http.MethodGet:
				_, _ = client.GetAccount(context.Background(), "acc123")
			case http.MethodPost:
				_, _ = client.CreateAccount(context.Background(), &types.CreateAccountRequest{})
			}

			if tt.shouldExist && receivedHeader == "" {
				t.Errorf("header %q should be present but was empty", tt.checkHeader)
			}

			if tt.expected != "" && receivedHeader != tt.expected {
				t.Errorf("header %q = %q, want %q", tt.checkHeader, receivedHeader, tt.expected)
			}
		})
	}
}

func TestClient_IdempotencyKey(t *testing.T) {
	t.Parallel()

	t.Run("POST request has IdempotencyKey", func(t *testing.T) {
		t.Parallel()

		var receivedKey string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedKey = r.Header.Get("IdempotencyKey")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{})
		}))
		defer server.Close()

		client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})
		_, _ = client.CreateAccount(context.Background(), &types.CreateAccountRequest{})

		if receivedKey == "" {
			t.Error("IdempotencyKey should be present for POST request")
		}

		// Verify it looks like a UUID
		if len(receivedKey) < 30 {
			t.Errorf("IdempotencyKey should be UUID-like, got %q", receivedKey)
		}
	})

	t.Run("PATCH request has IdempotencyKey", func(t *testing.T) {
		t.Parallel()

		var receivedKey string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedKey = r.Header.Get("IdempotencyKey")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})
		_ = client.UpdateAccount(context.Background(), "acc123", &types.UpdateAccountRequest{})

		if receivedKey == "" {
			t.Error("IdempotencyKey should be present for PATCH request")
		}
	})

	t.Run("GET request has NO IdempotencyKey", func(t *testing.T) {
		t.Parallel()

		var receivedKey string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedKey = r.Header.Get("IdempotencyKey")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{})
		}))
		defer server.Close()

		client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})
		_, _ = client.GetAccount(context.Background(), "acc123")

		if receivedKey != "" {
			t.Errorf("IdempotencyKey should NOT be present for GET request, got %q", receivedKey)
		}
	})

	t.Run("each POST generates unique IdempotencyKey", func(t *testing.T) {
		t.Parallel()

		keys := make(map[string]bool)
		var mu sync.Mutex

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("IdempotencyKey")
			mu.Lock()
			keys[key] = true
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{})
		}))
		defer server.Close()

		client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})

		// Make 10 requests
		for i := 0; i < 10; i++ {
			_, _ = client.CreateAccount(context.Background(), &types.CreateAccountRequest{})
		}

		mu.Lock()
		numUniqueKeys := len(keys)
		mu.Unlock()

		if numUniqueKeys != 10 {
			t.Errorf("expected 10 unique IdempotencyKeys, got %d", numUniqueKeys)
		}
	})

	t.Run("custom IdempotencyKey via context", func(t *testing.T) {
		t.Parallel()

		var receivedKey string
		customKey := "my-custom-idempotency-key-123"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedKey = r.Header.Get("IdempotencyKey")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{})
		}))
		defer server.Close()

		client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})

		// Use custom idempotency key via context
		ctx := WithIdempotencyKey(context.Background(), customKey)
		_, _ = client.CreateAccount(ctx, &types.CreateAccountRequest{})

		if receivedKey != customKey {
			t.Errorf("expected IdempotencyKey %q, got %q", customKey, receivedKey)
		}
	})

	t.Run("custom IdempotencyKey used on retries", func(t *testing.T) {
		t.Parallel()

		var receivedKeys []string
		var mu sync.Mutex
		customKey := "retry-idempotency-key-456"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("IdempotencyKey")
			mu.Lock()
			receivedKeys = append(receivedKeys, key)
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{})
		}))
		defer server.Close()

		client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})

		// Use same context with custom key for multiple "retry" requests
		ctx := WithIdempotencyKey(context.Background(), customKey)
		_, _ = client.CreateAccount(ctx, &types.CreateAccountRequest{})
		_, _ = client.CreateAccount(ctx, &types.CreateAccountRequest{})
		_, _ = client.CreateAccount(ctx, &types.CreateAccountRequest{})

		mu.Lock()
		defer mu.Unlock()

		if len(receivedKeys) != 3 {
			t.Fatalf("expected 3 requests, got %d", len(receivedKeys))
		}

		// All requests should have the same custom key
		for i, key := range receivedKeys {
			if key != customKey {
				t.Errorf("request %d: expected key %q, got %q", i, customKey, key)
			}
		}
	})

	t.Run("empty custom key falls back to auto-generation", func(t *testing.T) {
		t.Parallel()

		var receivedKey string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedKey = r.Header.Get("IdempotencyKey")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{})
		}))
		defer server.Close()

		client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})

		// Use empty custom key - should fall back to auto-generation
		ctx := WithIdempotencyKey(context.Background(), "")
		_, _ = client.CreateAccount(ctx, &types.CreateAccountRequest{})

		// Should have auto-generated UUID-like key
		if receivedKey == "" {
			t.Error("expected auto-generated IdempotencyKey, got empty")
		}
		if len(receivedKey) < 30 {
			t.Errorf("expected UUID-like IdempotencyKey, got %q", receivedKey)
		}
	})
}

// ============================================================================
// ERROR HANDLING TESTS
// ============================================================================

func TestClient_ErrorHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		statusCode     int
		response       any
		wantStatusCode int
		wantCode       string
		wantMessage    string
	}{
		{
			name:       "400 Bad Request with error body",
			statusCode: http.StatusBadRequest,
			response: map[string]any{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request parameters",
			},
			wantStatusCode: http.StatusBadRequest,
			wantCode:       "INVALID_REQUEST",
			wantMessage:    "Invalid request parameters",
		},
		{
			name:       "401 Unauthorized",
			statusCode: http.StatusUnauthorized,
			response: map[string]any{
				"code":    "UNAUTHORIZED",
				"message": "Invalid API key",
			},
			wantStatusCode: http.StatusUnauthorized,
			wantCode:       "UNAUTHORIZED",
			wantMessage:    "Invalid API key",
		},
		{
			name:       "403 Forbidden",
			statusCode: http.StatusForbidden,
			response: map[string]any{
				"code":    "FORBIDDEN",
				"message": "Access denied",
			},
			wantStatusCode: http.StatusForbidden,
			wantCode:       "FORBIDDEN",
			wantMessage:    "Access denied",
		},
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
			response: map[string]any{
				"code":    "NOT_FOUND",
				"message": "Account not found",
			},
			wantStatusCode: http.StatusNotFound,
			wantCode:       "NOT_FOUND",
			wantMessage:    "Account not found",
		},
		{
			name:       "409 Conflict",
			statusCode: http.StatusConflict,
			response: map[string]any{
				"code":    "CONFLICT",
				"message": "Account already exists",
			},
			wantStatusCode: http.StatusConflict,
			wantCode:       "CONFLICT",
		},
		{
			name:       "422 Unprocessable Entity",
			statusCode: http.StatusUnprocessableEntity,
			response: map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Validation failed",
				"details": []string{"field1 is required", "field2 is invalid"},
			},
			wantStatusCode: http.StatusUnprocessableEntity,
			wantCode:       "VALIDATION_ERROR",
		},
		{
			name:       "429 Rate Limited",
			statusCode: http.StatusTooManyRequests,
			response: map[string]any{
				"code":    "RATE_LIMITED",
				"message": "Too many requests",
			},
			wantStatusCode: http.StatusTooManyRequests,
			wantCode:       "RATE_LIMITED",
		},
		{
			name:           "500 Internal Server Error",
			statusCode:     http.StatusInternalServerError,
			response:       map[string]any{"code": "INTERNAL_ERROR", "message": "Internal error"},
			wantStatusCode: http.StatusInternalServerError,
			wantCode:       "INTERNAL_ERROR",
		},
		{
			name:           "503 Service Unavailable",
			statusCode:     http.StatusServiceUnavailable,
			response:       map[string]any{"code": "SERVICE_UNAVAILABLE", "message": "Service unavailable"},
			wantStatusCode: http.StatusServiceUnavailable,
			wantCode:       "SERVICE_UNAVAILABLE",
		},
		{
			name:           "504 Gateway Timeout",
			statusCode:     http.StatusGatewayTimeout,
			response:       map[string]any{"code": "TIMEOUT", "message": "Gateway timeout"},
			wantStatusCode: http.StatusGatewayTimeout,
		},
		{
			name:           "error with non-JSON response",
			statusCode:     http.StatusBadGateway,
			response:       nil, // Will send plain text
			wantStatusCode: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				if tt.response != nil {
					_ = json.NewEncoder(w).Encode(tt.response)
				} else {
					_, _ = w.Write([]byte("Bad Gateway"))
				}
			}))
			defer server.Close()

			client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})
			_, err := client.GetAccount(context.Background(), "test-acc")

			AssertError(t, err, "GetAccount")

			apiErr := AssertAPIError(t, err, tt.wantStatusCode)

			if tt.wantCode != "" && apiErr.Code != tt.wantCode {
				t.Errorf("APIError.Code = %q, want %q", apiErr.Code, tt.wantCode)
			}

			if tt.wantMessage != "" && apiErr.Message != tt.wantMessage {
				t.Errorf("APIError.Message = %q, want %q", apiErr.Message, tt.wantMessage)
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{
			name: "with code",
			err: &APIError{
				StatusCode: 400,
				Code:       "INVALID_REQUEST",
				Message:    "Bad request",
			},
			expected: "API error 400: [INVALID_REQUEST] Bad request",
		},
		{
			name: "without code",
			err: &APIError{
				StatusCode: 500,
				Message:    "Internal error",
			},
			expected: "API error 500: Internal error",
		},
		{
			name: "empty message",
			err: &APIError{
				StatusCode: 503,
				Code:       "SERVICE_UNAVAILABLE",
			},
			expected: "API error 503: [SERVICE_UNAVAILABLE] ",
		},
		{
			name: "with details",
			err: &APIError{
				StatusCode: 422,
				Code:       "VALIDATION_ERROR",
				Message:    "Validation failed",
				Details:    []string{"field1 is required"},
			},
			expected: "API error 422: [VALIDATION_ERROR] Validation failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.err.Error()
			AssertEqual(t, tt.expected, got, "Error()")
		})
	}
}

func TestAPIError_Is(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		err    *APIError
		target error
		want   bool
	}{
		{
			name:   "match by code",
			err:    &APIError{StatusCode: 400, Code: "INVALID_REQUEST", Message: "Bad request"},
			target: ErrInvalidRequest,
			want:   true,
		},
		{
			name:   "match by status code",
			err:    &APIError{StatusCode: 401, Message: "Unauthorized"},
			target: ErrUnauthorized,
			want:   true,
		},
		{
			name:   "no match different code",
			err:    &APIError{StatusCode: 400, Code: "INVALID_REQUEST", Message: "Bad request"},
			target: ErrCardBlocked,
			want:   false,
		},
		{
			name:   "no match different status",
			err:    &APIError{StatusCode: 400, Message: "Bad request"},
			target: ErrUnauthorized,
			want:   false,
		},
		{
			name:   "match account not found",
			err:    &APIError{StatusCode: 404, Code: "ACCOUNT_NOT_FOUND", Message: "Account not found"},
			target: ErrAccountNotFound,
			want:   true,
		},
		{
			name:   "match card blocked",
			err:    &APIError{StatusCode: 403, Code: "CARD_BLOCKED", Message: "Card is blocked"},
			target: ErrCardBlocked,
			want:   true,
		},
		{
			name:   "match insufficient limit",
			err:    &APIError{StatusCode: 422, Code: "INSUFFICIENT_LIMIT", Message: "Limit exceeded"},
			target: ErrInsufficientLimit,
			want:   true,
		},
		{
			name:   "non-APIError target returns false",
			err:    &APIError{StatusCode: 500, Code: "INTERNAL", Message: "Error"},
			target: errors.New("some error"),
			want:   false,
		},
		{
			name:   "code takes priority over status",
			err:    &APIError{StatusCode: 404, Code: "CARD_NOT_FOUND", Message: "Card not found"},
			target: ErrNotFound, // status 404 but no code
			want:   false,       // code mismatch, doesn't fall back to status
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := errors.Is(tt.err, tt.target)
			if got != tt.want {
				t.Errorf("errors.Is(%v, %v) = %v, want %v", tt.err, tt.target, got, tt.want)
			}
		})
	}
}

func TestSentinelErrors(t *testing.T) {
	t.Parallel()

	// Verify all sentinel errors are properly defined
	sentinels := map[string]*APIError{
		"ErrAccountNotFound":      ErrAccountNotFound,
		"ErrAccountAlreadyExists": ErrAccountAlreadyExists,
		"ErrAccountBlocked":       ErrAccountBlocked,
		"ErrCardNotFound":         ErrCardNotFound,
		"ErrCardBlocked":          ErrCardBlocked,
		"ErrCardExpired":          ErrCardExpired,
		"ErrCardCancelled":        ErrCardCancelled,
		"ErrInvalidDocument":      ErrInvalidDocument,
		"ErrInvalidRequest":       ErrInvalidRequest,
		"ErrInvalidPIN":           ErrInvalidPIN,
		"ErrInsufficientBalance":  ErrInsufficientBalance,
		"ErrInsufficientLimit":    ErrInsufficientLimit,
		"ErrUnauthorized":         ErrUnauthorized,
		"ErrForbidden":            ErrForbidden,
		"ErrNotFound":             ErrNotFound,
		"ErrConflict":             ErrConflict,
		"ErrTooManyRequests":      ErrTooManyRequests,
		"ErrInternalServer":       ErrInternalServer,
	}

	for name, sentinel := range sentinels {
		if sentinel == nil {
			t.Errorf("%s is nil", name)
		}
	}

	// Verify code-based sentinels have codes
	codeBasedSentinels := []*APIError{
		ErrAccountNotFound, ErrAccountAlreadyExists, ErrAccountBlocked,
		ErrCardNotFound, ErrCardBlocked, ErrCardExpired, ErrCardCancelled,
		ErrInvalidDocument, ErrInvalidRequest, ErrInvalidPIN,
		ErrInsufficientBalance, ErrInsufficientLimit,
	}
	for _, sentinel := range codeBasedSentinels {
		if sentinel.Code == "" {
			t.Errorf("code-based sentinel has empty code: %+v", sentinel)
		}
	}

	// Verify status-based sentinels have status codes
	statusBasedSentinels := map[*APIError]int{
		ErrUnauthorized:    401,
		ErrForbidden:       403,
		ErrNotFound:        404,
		ErrConflict:        409,
		ErrTooManyRequests: 429,
		ErrInternalServer:  500,
	}
	for sentinel, expectedStatus := range statusBasedSentinels {
		if sentinel.StatusCode != expectedStatus {
			t.Errorf("status sentinel %+v has wrong status, want %d", sentinel, expectedStatus)
		}
	}
}

// ============================================================================
// CONTEXT HANDLING TESTS
// ============================================================================

func TestClient_ContextCancellation(t *testing.T) {
	t.Parallel()

	t.Run("cancelled context before request", func(t *testing.T) {
		t.Parallel()

		server := NewMockServer(t, &MockServerConfig{
			Method:   http.MethodGet,
			Path:     "/accounts/acc123",
			Status:   http.StatusOK,
			Response: map[string]any{},
		})
		defer server.Close()

		client := NewTestClient(t, server)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := client.GetAccount(ctx, "acc123")
		AssertError(t, err, "GetAccount with cancelled context")
		AssertErrorContains(t, err, "context canceled")
	})

	t.Run("context timeout during slow request", func(t *testing.T) {
		t.Parallel()

		server := NewMockServer(t, &MockServerConfig{
			Method:   http.MethodGet,
			Path:     "/accounts/acc123",
			Status:   http.StatusOK,
			Response: map[string]any{},
			Delay:    500 * time.Millisecond,
		})
		defer server.Close()

		client := NewTestClient(t, server, TestClientConfig{
			APIKey:    "test-key",
			UserAgent: "test/1.0",
			Timeout:   50 * time.Millisecond,
		})

		ctx := context.Background()
		_, err := client.GetAccount(ctx, "acc123")
		AssertError(t, err, "GetAccount with timeout")
	})
}

// ============================================================================
// CONCURRENCY TESTS
// ============================================================================

func TestClient_ConcurrentRequests(t *testing.T) {
	t.Parallel()

	server := NewMockServer(t, &MockServerConfig{
		Method:   http.MethodGet,
		Path:     "", // Accept any path for concurrent test
		Status:   http.StatusOK,
		Response: NewAccountBuilder().Build(),
	})
	defer server.Close()

	client := NewTestClient(t, server)

	runner := NewConcurrentTestRunner(t)

	// Run 100 concurrent requests
	runner.Run(100, func() error {
		_, err := client.GetAccount(context.Background(), "acc123")
		return err
	})

	runner.Wait()

	// Verify all requests were received
	if server.RequestCount() != 100 {
		t.Errorf("expected 100 requests, got %d", server.RequestCount())
	}
}

func TestClient_ConcurrentDifferentEndpoints(t *testing.T) {
	t.Parallel()

	var accountRequests, cardRequests int
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		if strings.Contains(r.URL.Path, "/accounts") {
			accountRequests++
		}
		if strings.Contains(r.URL.Path, "/cards") {
			cardRequests++
		}
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{})
	}))
	defer server.Close()

	client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})

	var wg sync.WaitGroup
	errCh := make(chan error, 200)

	// 100 account requests
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := client.GetAccount(context.Background(), "acc123"); err != nil {
				errCh <- err
			}
		}()
	}

	// 100 card requests
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := client.GetCard(context.Background(), "card123"); err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("unexpected error: %v", err)
	}

	mu.Lock()
	if accountRequests != 100 {
		t.Errorf("expected 100 account requests, got %d", accountRequests)
	}
	if cardRequests != 100 {
		t.Errorf("expected 100 card requests, got %d", cardRequests)
	}
	mu.Unlock()
}

// ============================================================================
// EDGE CASE TESTS
// ============================================================================

func TestClient_EdgeCaseAccountIDs(t *testing.T) {
	t.Parallel()

	for _, id := range EdgeCaseIDs() {
		id := id
		t.Run("ID: "+id, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "", // Accept any path
				Status:   http.StatusOK,
				Response: NewAccountBuilder().WithID(id).Build(),
			})
			defer server.Close()

			client := NewTestClient(t, server)

			// Empty ID might cause issues - we just verify no panic
			_, _ = client.GetAccount(context.Background(), id)
			// We don't necessarily expect success, just no panic
		})
	}
}

func TestClient_EmptyResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Empty body
	}))
	defer server.Close()

	client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})

	// Should not panic on empty response
	_, err := client.GetAccount(context.Background(), "acc123")

	// Empty response might succeed (empty struct) or fail - just verify no panic
	_ = err
}

func TestClient_MalformedJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		response string
	}{
		{"incomplete JSON", `{"accountId": "acc123"`},
		{"invalid JSON", `not json at all`},
		{"empty object", `{}`},
		{"null response", `null`},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})

			// Should handle malformed JSON gracefully
			_, err := client.GetAccount(context.Background(), "acc123")

			// We just verify no panic - error behavior depends on JSON content
			_ = err
		})
	}
}

// ============================================================================
// REQUEST BODY VALIDATION TESTS
// ============================================================================

func TestClient_CreateAccount_RequestBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		request     *types.CreateAccountRequest
		validateReq func(t *testing.T, body map[string]any)
	}{
		{
			name: "full request",
			request: &types.CreateAccountRequest{
				PsProductCode: "CREDIT_CARD",
				AccountOwner: map[string]any{
					"fullName":               "Test User",
					"identityDocumentNumber": "12345678901",
				},
				BillingAddress: &types.Address{
					AddressLine1: "Rua Test 123",
					AddressLine2: "Apt 1",
					City:         "Sao Paulo",
					State:        "SP",
					Zipcode:      "01234567",
					Country:      "BR",
					Neighborhood: "Centro",
				},
				CardDeliveryAddress: map[string]any{
					"addressLine1": "Rua Test 123",
					"city":         "Sao Paulo",
				},
			},
			validateReq: func(t *testing.T, body map[string]any) {
				AssertEqual(t, "CREDIT_CARD", body["psProductCode"], "psProductCode")

				owner, ok := body["accountOwner"].(map[string]any)
				if !ok {
					t.Fatal("accountOwner not found in request")
				}
				AssertEqual(t, "Test User", owner["fullName"], "accountOwner.fullName")

				address, ok := body["billingAddress"].(map[string]any)
				if !ok {
					t.Fatal("billingAddress not found in request")
				}
				AssertEqual(t, "Rua Test 123", address["addressLine1"], "billingAddress.addressLine1")
			},
		},
		{
			name: "minimal request",
			request: &types.CreateAccountRequest{
				PsProductCode:       "DEBIT_CARD",
				AccountOwner:        map[string]any{"fullName": "Test"},
				BillingAddress:      map[string]any{"city": "SP"},
				CardDeliveryAddress: map[string]any{"city": "SP"},
			},
			validateReq: func(t *testing.T, body map[string]any) {
				AssertEqual(t, "DEBIT_CARD", body["psProductCode"], "psProductCode")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/accounts",
				Status:   http.StatusCreated,
				Response: map[string]any{"accountId": "acc123"},
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					var reqBody map[string]any
					if err := json.Unmarshal(body, &reqBody); err != nil {
						t.Fatalf("failed to unmarshal request body: %v", err)
					}
					tt.validateReq(t, reqBody)
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			_, err := client.CreateAccount(context.Background(), tt.request)

			AssertNoError(t, err, "CreateAccount")
		})
	}
}

// ============================================================================
// RESPONSE VALIDATION TESTS
// ============================================================================

func TestClient_GetAccount_ResponseValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mockResponse any
		validate     func(t *testing.T, resp *types.Account)
	}{
		{
			name: "full response",
			mockResponse: map[string]any{
				"accountId":     "acc-12345",
				"psProductCode": "CREDIT_CARD",
				"accountOwner": map[string]any{
					"fullName":               "John Doe",
					"identityDocumentNumber": "12345678901",
				},
			},
			validate: func(t *testing.T, resp *types.Account) {
				AssertEqual(t, "acc-12345", resp.AccountID, "AccountID")
				AssertEqual(t, "CREDIT_CARD", resp.PsProductCode, "PsProductCode")
				AssertNotNil(t, resp.AccountOwner, "AccountOwner")
			},
		},
		{
			name: "minimal response",
			mockResponse: map[string]any{
				"accountId":     "acc-minimal",
				"psProductCode": "DEBIT_CARD",
				"accountOwner":  map[string]any{"fullName": "Test"},
			},
			validate: func(t *testing.T, resp *types.Account) {
				AssertEqual(t, "acc-minimal", resp.AccountID, "AccountID")
				AssertEqual(t, "DEBIT_CARD", resp.PsProductCode, "PsProductCode")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "", // Accept any path
				Status:   http.StatusOK,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetAccount(context.Background(), "test-acc")

			AssertNoError(t, err, "GetAccount")
			AssertNotNil(t, resp, "response")
			tt.validate(t, resp)
		})
	}
}

// ============================================================================
// BENCHMARKS
// ============================================================================

func BenchmarkClient_GetAccount(b *testing.B) {
	server := NewBenchmarkMockServer(NewAccountBuilder().Build())
	defer server.Close()

	client := NewBenchmarkClient(b, server.URL)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetAccount(ctx, "acc123")
	}
}

func BenchmarkClient_CreateAccount(b *testing.B) {
	server := NewBenchmarkMockServer(map[string]any{"accountId": "acc123"})
	defer server.Close()

	client := NewBenchmarkClient(b, server.URL)
	ctx := context.Background()
	req := &types.CreateAccountRequest{
		PsProductCode:       "CREDIT",
		AccountOwner:        map[string]any{"fullName": "Benchmark User"},
		BillingAddress:      map[string]any{"city": "SP"},
		CardDeliveryAddress: map[string]any{"city": "SP"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.CreateAccount(ctx, req)
	}
}

func BenchmarkClient_ConcurrentRequests(b *testing.B) {
	server := NewBenchmarkMockServer(NewAccountBuilder().Build())
	defer server.Close()

	client := NewBenchmarkClient(b, server.URL)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.GetAccount(ctx, "acc123")
		}
	})
}
