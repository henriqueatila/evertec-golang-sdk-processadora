package client

import (
	"context"
	"testing"
	"time"
)

func TestNewHook(t *testing.T) {
	beforeCalled := false
	afterCalled := false

	hook := NewHook(
		func(ctx context.Context, req *RequestInfo) context.Context {
			beforeCalled = true
			return ctx
		},
		func(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
			afterCalled = true
		},
	)

	ctx := context.Background()
	reqInfo := &RequestInfo{Method: "GET", Path: "/test"}
	respInfo := &ResponseInfo{StatusCode: 200, Duration: time.Millisecond}

	ctx = hook.BeforeRequest(ctx, reqInfo)
	hook.AfterResponse(ctx, reqInfo, respInfo)

	if !beforeCalled {
		t.Error("BeforeRequest was not called")
	}
	if !afterCalled {
		t.Error("AfterResponse was not called")
	}
}

func TestNewHook_NilFunctions(t *testing.T) {
	hook := NewHook(nil, nil)

	ctx := context.Background()
	reqInfo := &RequestInfo{Method: "GET", Path: "/test"}
	respInfo := &ResponseInfo{StatusCode: 200}

	// Should not panic
	ctx = hook.BeforeRequest(ctx, reqInfo)
	hook.AfterResponse(ctx, reqInfo, respInfo)

	if ctx == nil {
		t.Error("BeforeRequest returned nil context")
	}
}

func TestNewAfterHook(t *testing.T) {
	afterCalled := false

	hook := NewAfterHook(func(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
		afterCalled = true
	})

	ctx := context.Background()
	reqInfo := &RequestInfo{Method: "POST", Path: "/accounts"}
	respInfo := &ResponseInfo{StatusCode: 201}

	// BeforeRequest should just return ctx
	ctx = hook.BeforeRequest(ctx, reqInfo)
	if ctx == nil {
		t.Error("BeforeRequest returned nil context")
	}

	hook.AfterResponse(ctx, reqInfo, respInfo)

	if !afterCalled {
		t.Error("AfterResponse was not called")
	}
}

func TestRequestInfo(t *testing.T) {
	req := &RequestInfo{
		Method:         "POST",
		Path:           "/accounts",
		BodySize:       256,
		IdempotencyKey: "test-key-123",
	}

	if req.Method != "POST" {
		t.Errorf("expected Method POST, got %s", req.Method)
	}
	if req.Path != "/accounts" {
		t.Errorf("expected Path /accounts, got %s", req.Path)
	}
	if req.BodySize != 256 {
		t.Errorf("expected BodySize 256, got %d", req.BodySize)
	}
	if req.IdempotencyKey != "test-key-123" {
		t.Errorf("expected IdempotencyKey test-key-123, got %s", req.IdempotencyKey)
	}
}

func TestResponseInfo(t *testing.T) {
	apiErr := &APIError{StatusCode: 400, Code: "INVALID", Message: "Invalid request"}
	resp := &ResponseInfo{
		StatusCode: 400,
		BodySize:   128,
		Duration:   50 * time.Millisecond,
		Error:      apiErr,
	}

	if resp.StatusCode != 400 {
		t.Errorf("expected StatusCode 400, got %d", resp.StatusCode)
	}
	if resp.BodySize != 128 {
		t.Errorf("expected BodySize 128, got %d", resp.BodySize)
	}
	if resp.Duration != 50*time.Millisecond {
		t.Errorf("expected Duration 50ms, got %v", resp.Duration)
	}
	if resp.Error == nil {
		t.Error("expected Error to be set")
	}
}

func TestHookContextPropagation(t *testing.T) {
	type ctxKey string
	const testKey ctxKey = "test-value"

	hook := NewHook(
		func(ctx context.Context, req *RequestInfo) context.Context {
			return context.WithValue(ctx, testKey, "hello")
		},
		func(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
			if v := ctx.Value(testKey); v != "hello" {
				t.Errorf("expected context value 'hello', got %v", v)
			}
		},
	)

	ctx := context.Background()
	reqInfo := &RequestInfo{Method: "GET", Path: "/test"}
	respInfo := &ResponseInfo{StatusCode: 200}

	ctx = hook.BeforeRequest(ctx, reqInfo)
	hook.AfterResponse(ctx, reqInfo, respInfo)
}
