package client

import (
	"context"
	"crypto/tls"
	"log/slog"
	"os"
	"testing"
)

func TestWithLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		logger *slog.Logger
	}{
		{
			name:   "custom logger",
			logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		},
		{
			name:   "text logger",
			logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		},
		{
			name:   "nil logger",
			logger: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &Config{}
			WithLogger(tt.logger)(cfg)

			if tt.logger != nil {
				if cfg.Logger == nil {
					t.Error("Logger should not be nil when provided")
				}
			}
		})
	}
}

func TestWithDefaultLogger(t *testing.T) {
	t.Parallel()

	cfg := &Config{}
	WithDefaultLogger()(cfg)

	if cfg.Logger == nil {
		t.Error("WithDefaultLogger should set a default logger")
	}
}

func TestWithHooks(t *testing.T) {
	t.Parallel()

	// Create test hooks
	hook1 := &testHook{name: "hook1"}
	hook2 := &testHook{name: "hook2"}

	tests := []struct {
		name      string
		hooks     []Hook
		wantCount int
	}{
		{
			name:      "single hook",
			hooks:     []Hook{hook1},
			wantCount: 1,
		},
		{
			name:      "multiple hooks",
			hooks:     []Hook{hook1, hook2},
			wantCount: 2,
		},
		{
			name:      "empty hooks",
			hooks:     []Hook{},
			wantCount: 0,
		},
		{
			name:      "nil hooks slice",
			hooks:     nil,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &Config{}
			WithHooks(tt.hooks...)(cfg)

			if len(cfg.Hooks) != tt.wantCount {
				t.Errorf("expected %d hooks, got %d", tt.wantCount, len(cfg.Hooks))
			}
		})
	}
}

// testHook is a simple implementation of Hook interface for testing
type testHook struct {
	name        string
	beforeCalls int
	afterCalls  int
}

func (h *testHook) BeforeRequest(ctx context.Context, req *RequestInfo) context.Context {
	h.beforeCalls++
	return ctx
}

func (h *testHook) AfterResponse(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
	h.afterCalls++
}

func TestWithHooks_Integration(t *testing.T) {
	t.Parallel()

	// Verify hooks are properly executed
	hook := &testHook{name: "integration-hook"}

	server := NewMockServer(t, &MockServerConfig{
		Status:   200,
		Response: map[string]interface{}{"accountId": "acc123"},
	})
	defer server.Close()

	// Create client directly with hooks
	clientCfg := Config{
		APIKey:    "test-api-key",
		UserAgent: "test/1.0",
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Hooks:     []Hook{hook},
	}

	client, err := New(clientCfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Make a request
	_, _ = client.GetAccount(context.Background(), "acc123")

	// Verify hooks were called
	if hook.beforeCalls != 1 {
		t.Errorf("expected beforeCalls=1, got %d", hook.beforeCalls)
	}
	if hook.afterCalls != 1 {
		t.Errorf("expected afterCalls=1, got %d", hook.afterCalls)
	}
}
