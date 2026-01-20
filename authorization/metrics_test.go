package authorization

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// mockMetricsCollector for testing MetricsHook
type mockMetricsCollector struct {
	mu            sync.Mutex
	calls         []metricsCall
	operationType string
	responseCode  int
	approved      bool
	duration      time.Duration
	err           error
}

type metricsCall struct {
	operationType string
	responseCode  int
	approved      bool
	duration      time.Duration
	err           error
}

func (m *mockMetricsCollector) RecordAuthorization(operationType string, responseCode int, approved bool, duration time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, metricsCall{
		operationType: operationType,
		responseCode:  responseCode,
		approved:      approved,
		duration:      duration,
		err:           err,
	})
	m.operationType = operationType
	m.responseCode = responseCode
	m.approved = approved
	m.duration = duration
	m.err = err
}

func TestMetricsHook_ImplementsInterface(t *testing.T) {
	var _ Hook = &MetricsHook{}
}

func TestNewMetricsHook(t *testing.T) {
	collector := &mockMetricsCollector{}
	hook := NewMetricsHook(collector)

	if hook == nil {
		t.Fatal("NewMetricsHook returned nil")
	}
	if hook.collector != collector {
		t.Error("collector not set correctly")
	}
}

func TestMetricsHook_BeforeAuthorization(t *testing.T) {
	hook := NewMetricsHook(&mockMetricsCollector{})
	ctx := context.Background()
	req := &RequestInfo{Path: "/test"}

	newCtx := hook.BeforeAuthorization(ctx, req)

	// Should return the same context unchanged
	if newCtx != ctx {
		t.Error("BeforeAuthorization should return unchanged context")
	}
}

func TestMetricsHook_AfterAuthorization(t *testing.T) {
	collector := &mockMetricsCollector{}
	hook := NewMetricsHook(collector)

	req := &RequestInfo{
		Path:          "/purchases",
		OperationType: "purchase",
	}
	resp := &ResponseInfo{
		StatusCode:   200,
		ResponseCode: 0,
		Duration:     100 * time.Millisecond,
		Approved:     true,
	}

	hook.AfterAuthorization(context.Background(), req, resp)

	if collector.operationType != "purchase" {
		t.Errorf("operationType = %s, want purchase", collector.operationType)
	}
	if collector.responseCode != 0 {
		t.Errorf("responseCode = %d, want 0", collector.responseCode)
	}
	if !collector.approved {
		t.Error("approved = false, want true")
	}
	if collector.duration != 100*time.Millisecond {
		t.Errorf("duration = %v, want 100ms", collector.duration)
	}
	if collector.err != nil {
		t.Errorf("err = %v, want nil", collector.err)
	}
}

func TestMetricsHook_AfterAuthorization_WithError(t *testing.T) {
	collector := &mockMetricsCollector{}
	hook := NewMetricsHook(collector)

	testErr := errors.New("test error")
	req := &RequestInfo{OperationType: "withdrawal"}
	resp := &ResponseInfo{
		StatusCode: 500,
		Duration:   50 * time.Millisecond,
		Error:      testErr,
		Approved:   false,
	}

	hook.AfterAuthorization(context.Background(), req, resp)

	if collector.err != testErr {
		t.Errorf("err = %v, want %v", collector.err, testErr)
	}
}

func TestNewSimpleMetrics(t *testing.T) {
	metrics := NewSimpleMetrics()

	if metrics == nil {
		t.Fatal("NewSimpleMetrics returned nil")
	}
	if metrics.ByOperation == nil {
		t.Error("ByOperation map not initialized")
	}
	if metrics.ByResponseCode == nil {
		t.Error("ByResponseCode map not initialized")
	}
	if metrics.Authorizations != 0 {
		t.Error("Authorizations should be 0 initially")
	}
}

func TestSimpleMetrics_ImplementsInterface(t *testing.T) {
	var _ MetricsCollector = &SimpleMetrics{}
}

func TestSimpleMetrics_RecordAuthorization_Approved(t *testing.T) {
	metrics := NewSimpleMetrics()

	metrics.RecordAuthorization("purchase", 0, true, 100*time.Millisecond, nil)

	if metrics.Authorizations != 1 {
		t.Errorf("Authorizations = %d, want 1", metrics.Authorizations)
	}
	if metrics.Approved != 1 {
		t.Errorf("Approved = %d, want 1", metrics.Approved)
	}
	if metrics.Declined != 0 {
		t.Errorf("Declined = %d, want 0", metrics.Declined)
	}
	if metrics.Errors != 0 {
		t.Errorf("Errors = %d, want 0", metrics.Errors)
	}
	if metrics.TotalMs != 100 {
		t.Errorf("TotalMs = %d, want 100", metrics.TotalMs)
	}
}

func TestSimpleMetrics_RecordAuthorization_Declined(t *testing.T) {
	metrics := NewSimpleMetrics()

	metrics.RecordAuthorization("purchase", 51, false, 50*time.Millisecond, nil)

	if metrics.Authorizations != 1 {
		t.Errorf("Authorizations = %d, want 1", metrics.Authorizations)
	}
	if metrics.Approved != 0 {
		t.Errorf("Approved = %d, want 0", metrics.Approved)
	}
	if metrics.Declined != 1 {
		t.Errorf("Declined = %d, want 1", metrics.Declined)
	}
}

func TestSimpleMetrics_RecordAuthorization_WithError(t *testing.T) {
	metrics := NewSimpleMetrics()

	metrics.RecordAuthorization("purchase", 0, false, 10*time.Millisecond, errors.New("test"))

	if metrics.Errors != 1 {
		t.Errorf("Errors = %d, want 1", metrics.Errors)
	}
}

func TestSimpleMetrics_RecordAuthorization_ByOperation(t *testing.T) {
	metrics := NewSimpleMetrics()

	// Record multiple operations
	metrics.RecordAuthorization("purchase", 0, true, 100*time.Millisecond, nil)
	metrics.RecordAuthorization("purchase", 51, false, 50*time.Millisecond, nil)
	metrics.RecordAuthorization("withdrawal", 0, true, 200*time.Millisecond, nil)
	metrics.RecordAuthorization("query", 0, true, 30*time.Millisecond, nil)

	// Check purchase metrics
	purchaseMetrics := metrics.ByOperation["purchase"]
	if purchaseMetrics == nil {
		t.Fatal("purchase metrics not found")
	}
	if purchaseMetrics.Total != 2 {
		t.Errorf("purchase.Total = %d, want 2", purchaseMetrics.Total)
	}
	if purchaseMetrics.Approved != 1 {
		t.Errorf("purchase.Approved = %d, want 1", purchaseMetrics.Approved)
	}
	if purchaseMetrics.Declined != 1 {
		t.Errorf("purchase.Declined = %d, want 1", purchaseMetrics.Declined)
	}
	if purchaseMetrics.TotalMs != 150 {
		t.Errorf("purchase.TotalMs = %d, want 150", purchaseMetrics.TotalMs)
	}

	// Check withdrawal metrics
	withdrawalMetrics := metrics.ByOperation["withdrawal"]
	if withdrawalMetrics == nil {
		t.Fatal("withdrawal metrics not found")
	}
	if withdrawalMetrics.Total != 1 {
		t.Errorf("withdrawal.Total = %d, want 1", withdrawalMetrics.Total)
	}

	// Check query metrics
	queryMetrics := metrics.ByOperation["query"]
	if queryMetrics == nil {
		t.Fatal("query metrics not found")
	}
	if queryMetrics.Total != 1 {
		t.Errorf("query.Total = %d, want 1", queryMetrics.Total)
	}
}

func TestSimpleMetrics_RecordAuthorization_ByResponseCode(t *testing.T) {
	metrics := NewSimpleMetrics()

	metrics.RecordAuthorization("purchase", 0, true, 50*time.Millisecond, nil)
	metrics.RecordAuthorization("purchase", 0, true, 50*time.Millisecond, nil)
	metrics.RecordAuthorization("purchase", 51, false, 50*time.Millisecond, nil)
	metrics.RecordAuthorization("purchase", 14, false, 50*time.Millisecond, nil)

	if metrics.ByResponseCode[0] != 2 {
		t.Errorf("ByResponseCode[0] = %d, want 2", metrics.ByResponseCode[0])
	}
	if metrics.ByResponseCode[51] != 1 {
		t.Errorf("ByResponseCode[51] = %d, want 1", metrics.ByResponseCode[51])
	}
	if metrics.ByResponseCode[14] != 1 {
		t.Errorf("ByResponseCode[14] = %d, want 1", metrics.ByResponseCode[14])
	}
}

func TestSimpleMetrics_Stats(t *testing.T) {
	metrics := NewSimpleMetrics()

	// Empty metrics
	total, approved, declined, errs, avgMs := metrics.Stats()
	if total != 0 || approved != 0 || declined != 0 || errs != 0 || avgMs != 0 {
		t.Error("Stats should return all zeros for empty metrics")
	}

	// Add some data
	metrics.RecordAuthorization("purchase", 0, true, 100*time.Millisecond, nil)
	metrics.RecordAuthorization("purchase", 0, true, 200*time.Millisecond, nil)
	metrics.RecordAuthorization("purchase", 51, false, 50*time.Millisecond, nil)
	metrics.RecordAuthorization("purchase", 0, false, 50*time.Millisecond, errors.New("err"))

	total, approved, declined, errs, avgMs = metrics.Stats()

	if total != 4 {
		t.Errorf("total = %d, want 4", total)
	}
	if approved != 2 {
		t.Errorf("approved = %d, want 2", approved)
	}
	if declined != 2 {
		t.Errorf("declined = %d, want 2", declined)
	}
	if errs != 1 {
		t.Errorf("errors = %d, want 1", errs)
	}
	// Average: (100+200+50+50)/4 = 100
	if avgMs != 100 {
		t.Errorf("avgMs = %d, want 100", avgMs)
	}
}

func TestSimpleMetrics_ApprovalRate(t *testing.T) {
	metrics := NewSimpleMetrics()

	// Empty metrics should return 0
	if rate := metrics.ApprovalRate(); rate != 0 {
		t.Errorf("ApprovalRate = %f, want 0 for empty metrics", rate)
	}

	// Add 3 approved, 1 declined = 75% approval rate
	metrics.RecordAuthorization("purchase", 0, true, 50*time.Millisecond, nil)
	metrics.RecordAuthorization("purchase", 0, true, 50*time.Millisecond, nil)
	metrics.RecordAuthorization("purchase", 0, true, 50*time.Millisecond, nil)
	metrics.RecordAuthorization("purchase", 51, false, 50*time.Millisecond, nil)

	rate := metrics.ApprovalRate()
	if rate != 75.0 {
		t.Errorf("ApprovalRate = %f, want 75.0", rate)
	}
}

func TestSimpleMetrics_ApprovalRate_100Percent(t *testing.T) {
	metrics := NewSimpleMetrics()

	metrics.RecordAuthorization("purchase", 0, true, 50*time.Millisecond, nil)
	metrics.RecordAuthorization("withdrawal", 0, true, 50*time.Millisecond, nil)

	rate := metrics.ApprovalRate()
	if rate != 100.0 {
		t.Errorf("ApprovalRate = %f, want 100.0", rate)
	}
}

func TestSimpleMetrics_ApprovalRate_0Percent(t *testing.T) {
	metrics := NewSimpleMetrics()

	metrics.RecordAuthorization("purchase", 51, false, 50*time.Millisecond, nil)
	metrics.RecordAuthorization("purchase", 14, false, 50*time.Millisecond, nil)

	rate := metrics.ApprovalRate()
	if rate != 0.0 {
		t.Errorf("ApprovalRate = %f, want 0.0", rate)
	}
}

func TestSimpleMetrics_Reset(t *testing.T) {
	metrics := NewSimpleMetrics()

	// Add some data
	metrics.RecordAuthorization("purchase", 0, true, 100*time.Millisecond, nil)
	metrics.RecordAuthorization("purchase", 51, false, 50*time.Millisecond, nil)
	metrics.RecordAuthorization("withdrawal", 0, true, 200*time.Millisecond, errors.New("err"))

	// Verify data exists
	if metrics.Authorizations != 3 {
		t.Fatalf("Expected 3 authorizations before reset")
	}

	// Reset
	metrics.Reset()

	// Verify all cleared
	if metrics.Authorizations != 0 {
		t.Errorf("Authorizations = %d, want 0 after reset", metrics.Authorizations)
	}
	if metrics.Approved != 0 {
		t.Errorf("Approved = %d, want 0 after reset", metrics.Approved)
	}
	if metrics.Declined != 0 {
		t.Errorf("Declined = %d, want 0 after reset", metrics.Declined)
	}
	if metrics.Errors != 0 {
		t.Errorf("Errors = %d, want 0 after reset", metrics.Errors)
	}
	if metrics.TotalMs != 0 {
		t.Errorf("TotalMs = %d, want 0 after reset", metrics.TotalMs)
	}
	if len(metrics.ByOperation) != 0 {
		t.Errorf("ByOperation has %d entries, want 0 after reset", len(metrics.ByOperation))
	}
	if len(metrics.ByResponseCode) != 0 {
		t.Errorf("ByResponseCode has %d entries, want 0 after reset", len(metrics.ByResponseCode))
	}
}

func TestSimpleMetrics_Concurrency(t *testing.T) {
	metrics := NewSimpleMetrics()
	var wg sync.WaitGroup

	// Simulate concurrent authorization recordings
	numGoroutines := 100
	numRecordsPerGoroutine := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numRecordsPerGoroutine; j++ {
				approved := j%2 == 0
				var err error
				if j%10 == 0 {
					err = errors.New("test error")
				}
				metrics.RecordAuthorization("purchase", j%5, approved, time.Duration(j)*time.Millisecond, err)
			}
		}(i)
	}

	wg.Wait()

	// Verify totals
	expectedTotal := int64(numGoroutines * numRecordsPerGoroutine)
	total, approved, declined, errs, _ := metrics.Stats()

	if total != expectedTotal {
		t.Errorf("total = %d, want %d", total, expectedTotal)
	}

	// Half should be approved (even j values)
	expectedApproved := expectedTotal / 2
	if approved != expectedApproved {
		t.Errorf("approved = %d, want %d", approved, expectedApproved)
	}

	expectedDeclined := expectedTotal / 2
	if declined != expectedDeclined {
		t.Errorf("declined = %d, want %d", declined, expectedDeclined)
	}

	// Every 10th record has an error
	expectedErrors := expectedTotal / 10
	if errs != expectedErrors {
		t.Errorf("errors = %d, want %d", errs, expectedErrors)
	}
}

func TestSimpleMetrics_Stats_Concurrency(t *testing.T) {
	metrics := NewSimpleMetrics()
	var wg sync.WaitGroup

	// Concurrent writes
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			metrics.RecordAuthorization("purchase", 0, true, time.Millisecond, nil)
		}
	}()

	// Concurrent reads
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			metrics.Stats()
			metrics.ApprovalRate()
		}
	}()

	wg.Wait()
	// If we get here without panic, concurrency is handled correctly
}

func TestSimpleMetrics_Reset_Concurrency(t *testing.T) {
	metrics := NewSimpleMetrics()
	var wg sync.WaitGroup

	// Concurrent writes and resets
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			metrics.RecordAuthorization("purchase", 0, true, time.Millisecond, nil)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			metrics.Reset()
			time.Sleep(time.Microsecond)
		}
	}()

	wg.Wait()
	// If we get here without panic, concurrency is handled correctly
}

func TestOperationMetrics_Fields(t *testing.T) {
	op := &OperationMetrics{
		Total:    100,
		Approved: 80,
		Declined: 20,
		TotalMs:  5000,
	}

	if op.Total != 100 {
		t.Errorf("Total = %d, want 100", op.Total)
	}
	if op.Approved != 80 {
		t.Errorf("Approved = %d, want 80", op.Approved)
	}
	if op.Declined != 20 {
		t.Errorf("Declined = %d, want 20", op.Declined)
	}
	if op.TotalMs != 5000 {
		t.Errorf("TotalMs = %d, want 5000", op.TotalMs)
	}
}
