package client

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSimpleMetrics_RecordRequest(t *testing.T) {
	m := NewSimpleMetrics()

	// Record successful request
	m.RecordRequest("GET", "/accounts", 200, 50*time.Millisecond, nil)

	requests, errs, avgMs := m.Stats()
	if requests != 1 {
		t.Errorf("expected 1 request, got %d", requests)
	}
	if errs != 0 {
		t.Errorf("expected 0 errors, got %d", errs)
	}
	if avgMs != 50 {
		t.Errorf("expected 50ms avg, got %d", avgMs)
	}
}

func TestSimpleMetrics_RecordErrors(t *testing.T) {
	m := NewSimpleMetrics()

	// Record error response (4xx)
	m.RecordRequest("POST", "/accounts", 400, 30*time.Millisecond, nil)

	// Record error response (5xx)
	m.RecordRequest("GET", "/cards", 500, 100*time.Millisecond, nil)

	// Record network error
	m.RecordRequest("GET", "/test", 0, 10*time.Millisecond, errors.New("connection refused"))

	requests, errs, _ := m.Stats()
	if requests != 3 {
		t.Errorf("expected 3 requests, got %d", requests)
	}
	if errs != 3 {
		t.Errorf("expected 3 errors, got %d", errs)
	}
}

func TestSimpleMetrics_ByMethod(t *testing.T) {
	m := NewSimpleMetrics()

	m.RecordRequest("GET", "/accounts", 200, time.Millisecond, nil)
	m.RecordRequest("GET", "/cards", 200, time.Millisecond, nil)
	m.RecordRequest("POST", "/accounts", 201, time.Millisecond, nil)

	if m.ByMethod["GET"] != 2 {
		t.Errorf("expected 2 GET requests, got %d", m.ByMethod["GET"])
	}
	if m.ByMethod["POST"] != 1 {
		t.Errorf("expected 1 POST request, got %d", m.ByMethod["POST"])
	}
}

func TestSimpleMetrics_ByPath(t *testing.T) {
	m := NewSimpleMetrics()

	m.RecordRequest("GET", "/accounts", 200, time.Millisecond, nil)
	m.RecordRequest("GET", "/accounts", 200, time.Millisecond, nil)
	m.RecordRequest("GET", "/cards", 200, time.Millisecond, nil)

	if m.ByPath["/accounts"] != 2 {
		t.Errorf("expected 2 /accounts requests, got %d", m.ByPath["/accounts"])
	}
	if m.ByPath["/cards"] != 1 {
		t.Errorf("expected 1 /cards request, got %d", m.ByPath["/cards"])
	}
}

func TestSimpleMetrics_ByStatus(t *testing.T) {
	m := NewSimpleMetrics()

	m.RecordRequest("GET", "/test", 200, time.Millisecond, nil)
	m.RecordRequest("GET", "/test", 200, time.Millisecond, nil)
	m.RecordRequest("GET", "/test", 404, time.Millisecond, nil)
	m.RecordRequest("POST", "/test", 500, time.Millisecond, nil)

	if m.ByStatus[200] != 2 {
		t.Errorf("expected 2 status 200, got %d", m.ByStatus[200])
	}
	if m.ByStatus[404] != 1 {
		t.Errorf("expected 1 status 404, got %d", m.ByStatus[404])
	}
	if m.ByStatus[500] != 1 {
		t.Errorf("expected 1 status 500, got %d", m.ByStatus[500])
	}
}

func TestSimpleMetrics_Reset(t *testing.T) {
	m := NewSimpleMetrics()

	m.RecordRequest("GET", "/test", 200, 100*time.Millisecond, nil)
	m.RecordRequest("POST", "/test", 400, 50*time.Millisecond, nil)

	m.Reset()

	requests, errs, avgMs := m.Stats()
	if requests != 0 {
		t.Errorf("expected 0 requests after reset, got %d", requests)
	}
	if errs != 0 {
		t.Errorf("expected 0 errors after reset, got %d", errs)
	}
	if avgMs != 0 {
		t.Errorf("expected 0 avgMs after reset, got %d", avgMs)
	}
	if len(m.ByMethod) != 0 {
		t.Error("expected empty ByMethod after reset")
	}
	if len(m.ByPath) != 0 {
		t.Error("expected empty ByPath after reset")
	}
	if len(m.ByStatus) != 0 {
		t.Error("expected empty ByStatus after reset")
	}
}

func TestSimpleMetrics_AverageMs(t *testing.T) {
	m := NewSimpleMetrics()

	m.RecordRequest("GET", "/test", 200, 100*time.Millisecond, nil)
	m.RecordRequest("GET", "/test", 200, 200*time.Millisecond, nil)
	m.RecordRequest("GET", "/test", 200, 300*time.Millisecond, nil)

	_, _, avgMs := m.Stats()
	if avgMs != 200 {
		t.Errorf("expected 200ms avg, got %d", avgMs)
	}
}

func TestMetricsHook(t *testing.T) {
	m := NewSimpleMetrics()
	hook := NewMetricsHook(m)

	ctx := context.Background()
	reqInfo := &RequestInfo{Method: "GET", Path: "/accounts"}
	respInfo := &ResponseInfo{StatusCode: 200, Duration: 75 * time.Millisecond}

	// BeforeRequest should just return context
	ctx = hook.BeforeRequest(ctx, reqInfo)
	if ctx == nil {
		t.Error("BeforeRequest returned nil context")
	}

	// AfterResponse should record metrics
	hook.AfterResponse(ctx, reqInfo, respInfo)

	requests, _, avgMs := m.Stats()
	if requests != 1 {
		t.Errorf("expected 1 request, got %d", requests)
	}
	if avgMs != 75 {
		t.Errorf("expected 75ms avg, got %d", avgMs)
	}
}

func TestMetricsHook_WithError(t *testing.T) {
	m := NewSimpleMetrics()
	hook := NewMetricsHook(m)

	ctx := context.Background()
	reqInfo := &RequestInfo{Method: "POST", Path: "/accounts"}
	respInfo := &ResponseInfo{
		StatusCode: 400,
		Duration:   25 * time.Millisecond,
		Error:      &APIError{StatusCode: 400, Code: "INVALID"},
	}

	hook.AfterResponse(ctx, reqInfo, respInfo)

	_, errs, _ := m.Stats()
	if errs != 1 {
		t.Errorf("expected 1 error, got %d", errs)
	}
}
