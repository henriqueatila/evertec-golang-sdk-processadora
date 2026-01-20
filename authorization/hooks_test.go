package authorization

import (
	"context"
	"testing"
	"time"
)

// testCtxKey is a custom type for context keys to avoid staticcheck SA1029
type testCtxKey string

func TestRequestInfo_Fields(t *testing.T) {
	req := &RequestInfo{
		Path:          "/purchases",
		OperationType: "purchase",
		TransactionID: "txn123",
		AccountID:     "acc456",
		CardID:        "card789",
		Amount:        10000,
		Currency:      "BRL",
		BodySize:      1024,
	}

	if req.Path != "/purchases" {
		t.Errorf("Path = %s, want /purchases", req.Path)
	}
	if req.OperationType != "purchase" {
		t.Errorf("OperationType = %s, want purchase", req.OperationType)
	}
	if req.TransactionID != "txn123" {
		t.Errorf("TransactionID = %s, want txn123", req.TransactionID)
	}
	if req.AccountID != "acc456" {
		t.Errorf("AccountID = %s, want acc456", req.AccountID)
	}
	if req.CardID != "card789" {
		t.Errorf("CardID = %s, want card789", req.CardID)
	}
	if req.Amount != 10000 {
		t.Errorf("Amount = %d, want 10000", req.Amount)
	}
	if req.Currency != "BRL" {
		t.Errorf("Currency = %s, want BRL", req.Currency)
	}
	if req.BodySize != 1024 {
		t.Errorf("BodySize = %d, want 1024", req.BodySize)
	}
}

func TestResponseInfo_Fields(t *testing.T) {
	resp := &ResponseInfo{
		StatusCode:   200,
		ResponseCode: 0,
		Duration:     100 * time.Millisecond,
		Error:        nil,
		Approved:     true,
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode)
	}
	if resp.ResponseCode != 0 {
		t.Errorf("ResponseCode = %d, want 0", resp.ResponseCode)
	}
	if resp.Duration != 100*time.Millisecond {
		t.Errorf("Duration = %v, want 100ms", resp.Duration)
	}
	if resp.Error != nil {
		t.Errorf("Error = %v, want nil", resp.Error)
	}
	if !resp.Approved {
		t.Error("Approved = false, want true")
	}
}

func TestNewHook_BothFunctions(t *testing.T) {
	beforeCalled := false
	afterCalled := false
	var capturedReq *RequestInfo
	var capturedResp *ResponseInfo

	hook := NewHook(
		func(ctx context.Context, req *RequestInfo) context.Context {
			beforeCalled = true
			return context.WithValue(ctx, testCtxKey("test"), "value")
		},
		func(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
			afterCalled = true
			capturedReq = req
			capturedResp = resp
			if ctx.Value(testCtxKey("test")) != "value" {
				t.Error("context value not passed from BeforeAuthorization")
			}
		},
	)

	req := &RequestInfo{Path: "/test", OperationType: "test"}
	resp := &ResponseInfo{StatusCode: 200, Approved: true}

	ctx := hook.BeforeAuthorization(context.Background(), req)
	hook.AfterAuthorization(ctx, req, resp)

	if !beforeCalled {
		t.Error("BeforeAuthorization was not called")
	}
	if !afterCalled {
		t.Error("AfterAuthorization was not called")
	}
	if capturedReq != req {
		t.Error("Request not passed to AfterAuthorization")
	}
	if capturedResp != resp {
		t.Error("Response not passed to AfterAuthorization")
	}
}

func TestNewHook_NilBefore(t *testing.T) {
	afterCalled := false

	hook := NewHook(
		nil,
		func(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
			afterCalled = true
		},
	)

	ctx := context.Background()
	req := &RequestInfo{Path: "/test"}
	resp := &ResponseInfo{StatusCode: 200}

	// Should return the same context when before is nil
	newCtx := hook.BeforeAuthorization(ctx, req)
	if newCtx != ctx {
		t.Error("BeforeAuthorization should return same context when function is nil")
	}

	hook.AfterAuthorization(newCtx, req, resp)
	if !afterCalled {
		t.Error("AfterAuthorization was not called")
	}
}

func TestNewHook_NilAfter(t *testing.T) {
	beforeCalled := false

	hook := NewHook(
		func(ctx context.Context, req *RequestInfo) context.Context {
			beforeCalled = true
			return ctx
		},
		nil,
	)

	ctx := context.Background()
	req := &RequestInfo{Path: "/test"}
	resp := &ResponseInfo{StatusCode: 200}

	newCtx := hook.BeforeAuthorization(ctx, req)
	if !beforeCalled {
		t.Error("BeforeAuthorization was not called")
	}

	// Should not panic when after is nil
	hook.AfterAuthorization(newCtx, req, resp)
}

func TestNewHook_BothNil(t *testing.T) {
	hook := NewHook(nil, nil)

	ctx := context.Background()
	req := &RequestInfo{Path: "/test"}
	resp := &ResponseInfo{StatusCode: 200}

	// Should not panic
	newCtx := hook.BeforeAuthorization(ctx, req)
	if newCtx != ctx {
		t.Error("BeforeAuthorization should return same context when function is nil")
	}

	hook.AfterAuthorization(newCtx, req, resp)
}

func TestNewAfterHook(t *testing.T) {
	afterCalled := false
	var capturedReq *RequestInfo

	hook := NewAfterHook(func(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
		afterCalled = true
		capturedReq = req
	})

	ctx := context.Background()
	req := &RequestInfo{Path: "/test", OperationType: "purchase"}
	resp := &ResponseInfo{StatusCode: 200}

	// Before should return unchanged context
	newCtx := hook.BeforeAuthorization(ctx, req)
	if newCtx != ctx {
		t.Error("BeforeAuthorization should return same context for AfterHook")
	}

	hook.AfterAuthorization(newCtx, req, resp)

	if !afterCalled {
		t.Error("AfterAuthorization was not called")
	}
	if capturedReq != req {
		t.Error("Request not passed correctly")
	}
}

func TestFuncHook_ImplementsInterface(t *testing.T) {
	var _ Hook = &funcHook{}
	var _ Hook = NewHook(nil, nil)
	var _ Hook = NewAfterHook(nil)
}

func TestHook_ContextPropagation(t *testing.T) {
	type ctxKey string
	const key ctxKey = "requestID"

	hook := NewHook(
		func(ctx context.Context, req *RequestInfo) context.Context {
			return context.WithValue(ctx, key, "req-12345")
		},
		func(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
			if ctx.Value(key) != "req-12345" {
				t.Errorf("context value = %v, want req-12345", ctx.Value(key))
			}
		},
	)

	ctx := hook.BeforeAuthorization(context.Background(), &RequestInfo{})
	hook.AfterAuthorization(ctx, &RequestInfo{}, &ResponseInfo{})
}
