package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"testing"
	"time"
)

// =============================================================================
// Hook Tests
// =============================================================================

func TestNewHook_BothFunctions(t *testing.T) {
	beforeCalled := false
	afterCalled := false

	hook := NewHook(
		func(ctx context.Context, event *EventInfo) context.Context {
			beforeCalled = true
			return ctx
		},
		func(ctx context.Context, event *EventInfo, result *ProcessingInfo) {
			afterCalled = true
		},
	)

	ctx := context.Background()
	event := &EventInfo{EventType: EventTypeDisputeStatusChange}
	result := &ProcessingInfo{StatusCode: 200, Processed: true}

	hook.BeforeEvent(ctx, event)
	hook.AfterEvent(ctx, event, result)

	if !beforeCalled {
		t.Error("BeforeEvent was not called")
	}
	if !afterCalled {
		t.Error("AfterEvent was not called")
	}
}

func TestNewHook_NilBefore(t *testing.T) {
	afterCalled := false

	hook := NewHook(
		nil,
		func(ctx context.Context, event *EventInfo, result *ProcessingInfo) {
			afterCalled = true
		},
	)

	ctx := context.Background()
	event := &EventInfo{EventType: EventTypeDisputeStatusChange}
	result := &ProcessingInfo{StatusCode: 200}

	// Should not panic with nil before
	returnedCtx := hook.BeforeEvent(ctx, event)
	if returnedCtx != ctx {
		t.Error("BeforeEvent should return original context when nil")
	}

	hook.AfterEvent(ctx, event, result)
	if !afterCalled {
		t.Error("AfterEvent was not called")
	}
}

func TestNewHook_NilAfter(t *testing.T) {
	beforeCalled := false

	hook := NewHook(
		func(ctx context.Context, event *EventInfo) context.Context {
			beforeCalled = true
			return ctx
		},
		nil,
	)

	ctx := context.Background()
	event := &EventInfo{EventType: EventTypeDisputeStatusChange}
	result := &ProcessingInfo{StatusCode: 200}

	hook.BeforeEvent(ctx, event)
	if !beforeCalled {
		t.Error("BeforeEvent was not called")
	}

	// Should not panic with nil after
	hook.AfterEvent(ctx, event, result)
}

func TestNewHook_BothNil(t *testing.T) {
	hook := NewHook(nil, nil)

	ctx := context.Background()
	event := &EventInfo{EventType: EventTypeDisputeStatusChange}
	result := &ProcessingInfo{StatusCode: 200}

	// Should not panic
	returnedCtx := hook.BeforeEvent(ctx, event)
	if returnedCtx != ctx {
		t.Error("BeforeEvent should return original context when nil")
	}

	hook.AfterEvent(ctx, event, result)
}

func TestNewAfterHook(t *testing.T) {
	afterCalled := false

	hook := NewAfterHook(func(ctx context.Context, event *EventInfo, result *ProcessingInfo) {
		afterCalled = true
	})

	ctx := context.Background()
	event := &EventInfo{EventType: EventTypeDisputeStatusChange}
	result := &ProcessingInfo{StatusCode: 200}

	// Before should be no-op
	returnedCtx := hook.BeforeEvent(ctx, event)
	if returnedCtx != ctx {
		t.Error("BeforeEvent should return original context")
	}

	hook.AfterEvent(ctx, event, result)
	if !afterCalled {
		t.Error("AfterEvent was not called")
	}
}

func TestFuncHook_ImplementsInterface(t *testing.T) {
	var _ Hook = &funcHook{}
}

func TestHook_ContextPropagation(t *testing.T) {
	type ctxKey string
	key := ctxKey("test-key")

	hook := NewHook(
		func(ctx context.Context, event *EventInfo) context.Context {
			return context.WithValue(ctx, key, "test-value")
		},
		func(ctx context.Context, event *EventInfo, result *ProcessingInfo) {
			if ctx.Value(key) != "test-value" {
				t.Error("Context value not propagated")
			}
		},
	)

	ctx := context.Background()
	event := &EventInfo{EventType: EventTypeDisputeStatusChange}
	result := &ProcessingInfo{StatusCode: 200}

	newCtx := hook.BeforeEvent(ctx, event)
	hook.AfterEvent(newCtx, event, result)
}

// =============================================================================
// LoggingHook Tests
// =============================================================================

func TestNewLoggingHook(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(&bytes.Buffer{}, nil))
	hook := NewLoggingHook(logger)

	if hook == nil {
		t.Fatal("NewLoggingHook returned nil")
	}
	if hook.logger != logger {
		t.Error("Logger not set correctly")
	}
}

func TestLoggingHook_ImplementsInterface(t *testing.T) {
	var _ Hook = &LoggingHook{}
}

func TestLoggingHook_BeforeEvent(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(&bytes.Buffer{}, nil))
	hook := NewLoggingHook(logger)

	ctx := context.Background()
	event := &EventInfo{EventType: EventTypeDisputeStatusChange}

	// BeforeEvent should return same context (no-op)
	returnedCtx := hook.BeforeEvent(ctx, event)
	if returnedCtx != ctx {
		t.Error("BeforeEvent should return same context")
	}
}

func TestLoggingHook_AfterEvent_NilLogger(t *testing.T) {
	hook := &LoggingHook{logger: nil}

	ctx := context.Background()
	event := &EventInfo{EventType: EventTypeDisputeStatusChange}
	result := &ProcessingInfo{StatusCode: 200, Processed: true}

	// Should not panic with nil logger
	hook.AfterEvent(ctx, event, result)
}

func TestLoggingHook_AfterEvent_Success(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	hook := NewLoggingHook(logger)

	ctx := context.Background()
	event := &EventInfo{
		EventType: EventTypeDisputeStatusChange,
		EventID:   "evt-123",
		BodySize:  256,
	}
	result := &ProcessingInfo{
		StatusCode: 200,
		Duration:   100 * time.Millisecond,
		Processed:  true,
	}

	hook.AfterEvent(ctx, event, result)

	// Parse log output
	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Check log level is INFO for success
	if logEntry["level"] != "INFO" {
		t.Errorf("expected level INFO, got %v", logEntry["level"])
	}

	// Check event_type
	if logEntry["event_type"] != string(EventTypeDisputeStatusChange) {
		t.Errorf("expected event_type %s, got %v", EventTypeDisputeStatusChange, logEntry["event_type"])
	}

	// Check event_id
	if logEntry["event_id"] != "evt-123" {
		t.Errorf("expected event_id evt-123, got %v", logEntry["event_id"])
	}

	// Check status
	if status, ok := logEntry["status"].(float64); !ok || int(status) != 200 {
		t.Errorf("expected status 200, got %v", logEntry["status"])
	}

	// Check processed
	if logEntry["processed"] != true {
		t.Errorf("expected processed true, got %v", logEntry["processed"])
	}
}

func TestLoggingHook_AfterEvent_Duplicate(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	hook := NewLoggingHook(logger)

	ctx := context.Background()
	event := &EventInfo{
		EventType:   EventTypeStatementClosed,
		IsDuplicate: true,
	}
	result := &ProcessingInfo{
		StatusCode: 200,
		Processed:  false, // Duplicate events are not reprocessed
	}

	hook.AfterEvent(ctx, event, result)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Check log level is WARN for not processed
	if logEntry["level"] != "WARN" {
		t.Errorf("expected level WARN for duplicate, got %v", logEntry["level"])
	}

	// Check duplicate flag is present
	if logEntry["duplicate"] != true {
		t.Errorf("expected duplicate true, got %v", logEntry["duplicate"])
	}
}

func TestLoggingHook_AfterEvent_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	hook := NewLoggingHook(logger)

	ctx := context.Background()
	event := &EventInfo{
		EventType: EventTypePaymentDue,
	}
	result := &ProcessingInfo{
		StatusCode: 500,
		Error:      errors.New("handler failed"),
		Processed:  false,
	}

	hook.AfterEvent(ctx, event, result)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Check log level is ERROR
	if logEntry["level"] != "ERROR" {
		t.Errorf("expected level ERROR, got %v", logEntry["level"])
	}

	// Check error is present
	if logEntry["error"] != "handler failed" {
		t.Errorf("expected error 'handler failed', got %v", logEntry["error"])
	}
}

func TestLoggingHook_AfterEvent_ServerError(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	hook := NewLoggingHook(logger)

	ctx := context.Background()
	event := &EventInfo{
		EventType: EventTypeCardStatusChange,
	}
	result := &ProcessingInfo{
		StatusCode: 503,
		Processed:  false,
	}

	hook.AfterEvent(ctx, event, result)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Check log level is ERROR for 5xx
	if logEntry["level"] != "ERROR" {
		t.Errorf("expected level ERROR for 5xx, got %v", logEntry["level"])
	}
}

func TestLoggingHook_AfterEvent_NoEventID(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	hook := NewLoggingHook(logger)

	ctx := context.Background()
	event := &EventInfo{
		EventType: EventTypeDeviceTokenStatus,
		EventID:   "", // No event ID
	}
	result := &ProcessingInfo{
		StatusCode: 200,
		Processed:  true,
	}

	hook.AfterEvent(ctx, event, result)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// event_id should not be present when empty
	if _, ok := logEntry["event_id"]; ok {
		t.Error("event_id should not be present when empty")
	}
}

// =============================================================================
// MetricsHook Tests
// =============================================================================

func TestNewMetricsHook(t *testing.T) {
	collector := NewSimpleMetrics()
	hook := NewMetricsHook(collector)

	if hook == nil {
		t.Fatal("NewMetricsHook returned nil")
	}
	if hook.collector != collector {
		t.Error("Collector not set correctly")
	}
}

func TestMetricsHook_ImplementsInterface(t *testing.T) {
	var _ Hook = &MetricsHook{}
}

func TestMetricsHook_BeforeEvent(t *testing.T) {
	collector := NewSimpleMetrics()
	hook := NewMetricsHook(collector)

	ctx := context.Background()
	event := &EventInfo{EventType: EventTypeDisputeStatusChange}

	// BeforeEvent should return same context (no-op)
	returnedCtx := hook.BeforeEvent(ctx, event)
	if returnedCtx != ctx {
		t.Error("BeforeEvent should return same context")
	}
}

func TestMetricsHook_AfterEvent(t *testing.T) {
	collector := NewSimpleMetrics()
	hook := NewMetricsHook(collector)

	ctx := context.Background()
	event := &EventInfo{
		EventType:   EventTypeDisputeStatusChange,
		IsDuplicate: false,
	}
	result := &ProcessingInfo{
		StatusCode: 200,
		Duration:   50 * time.Millisecond,
		Processed:  true,
	}

	hook.AfterEvent(ctx, event, result)

	total, processed, duplicates, errs, _ := collector.Stats()
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if processed != 1 {
		t.Errorf("expected processed 1, got %d", processed)
	}
	if duplicates != 0 {
		t.Errorf("expected duplicates 0, got %d", duplicates)
	}
	if errs != 0 {
		t.Errorf("expected errors 0, got %d", errs)
	}
}

func TestMetricsHook_AfterEvent_WithError(t *testing.T) {
	collector := NewSimpleMetrics()
	hook := NewMetricsHook(collector)

	ctx := context.Background()
	event := &EventInfo{EventType: EventTypePaymentDue}
	result := &ProcessingInfo{
		StatusCode: 500,
		Duration:   100 * time.Millisecond,
		Error:      errors.New("processing failed"),
		Processed:  false,
	}

	hook.AfterEvent(ctx, event, result)

	total, processed, duplicates, errs, _ := collector.Stats()
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if processed != 0 {
		t.Errorf("expected processed 0, got %d", processed)
	}
	if duplicates != 0 {
		t.Errorf("expected duplicates 0, got %d", duplicates)
	}
	if errs != 1 {
		t.Errorf("expected errors 1, got %d", errs)
	}
}

func TestMetricsHook_AfterEvent_Duplicate(t *testing.T) {
	collector := NewSimpleMetrics()
	hook := NewMetricsHook(collector)

	ctx := context.Background()
	event := &EventInfo{
		EventType:   EventTypeStatementClosed,
		IsDuplicate: true,
	}
	result := &ProcessingInfo{
		StatusCode: 200,
		Duration:   10 * time.Millisecond,
		Processed:  false,
	}

	hook.AfterEvent(ctx, event, result)

	total, processed, duplicates, errs, _ := collector.Stats()
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if processed != 0 {
		t.Errorf("expected processed 0 for duplicate, got %d", processed)
	}
	if duplicates != 1 {
		t.Errorf("expected duplicates 1, got %d", duplicates)
	}
	if errs != 0 {
		t.Errorf("expected errors 0, got %d", errs)
	}
}

// =============================================================================
// SimpleMetrics Tests
// =============================================================================

func TestNewSimpleMetrics(t *testing.T) {
	metrics := NewSimpleMetrics()

	if metrics == nil {
		t.Fatal("NewSimpleMetrics returned nil")
	}
	if metrics.ByEventType == nil {
		t.Error("ByEventType map not initialized")
	}
}

func TestSimpleMetrics_ImplementsInterface(t *testing.T) {
	var _ MetricsCollector = &SimpleMetrics{}
}

func TestSimpleMetrics_RecordEvent_Processed(t *testing.T) {
	metrics := NewSimpleMetrics()

	metrics.RecordEvent(EventTypeDisputeStatusChange, true, false, 50*time.Millisecond, nil)

	if metrics.Events != 1 {
		t.Errorf("expected Events 1, got %d", metrics.Events)
	}
	if metrics.Processed != 1 {
		t.Errorf("expected Processed 1, got %d", metrics.Processed)
	}
	if metrics.Duplicates != 0 {
		t.Errorf("expected Duplicates 0, got %d", metrics.Duplicates)
	}
	if metrics.Errors != 0 {
		t.Errorf("expected Errors 0, got %d", metrics.Errors)
	}
	if metrics.TotalMs != 50 {
		t.Errorf("expected TotalMs 50, got %d", metrics.TotalMs)
	}
}

func TestSimpleMetrics_RecordEvent_Duplicate(t *testing.T) {
	metrics := NewSimpleMetrics()

	metrics.RecordEvent(EventTypeCardStatusChange, false, true, 10*time.Millisecond, nil)

	if metrics.Duplicates != 1 {
		t.Errorf("expected Duplicates 1, got %d", metrics.Duplicates)
	}
	if metrics.Processed != 0 {
		t.Errorf("expected Processed 0, got %d", metrics.Processed)
	}
}

func TestSimpleMetrics_RecordEvent_Error(t *testing.T) {
	metrics := NewSimpleMetrics()

	metrics.RecordEvent(EventTypePaymentDue, false, false, 100*time.Millisecond, errors.New("failed"))

	if metrics.Errors != 1 {
		t.Errorf("expected Errors 1, got %d", metrics.Errors)
	}
	if metrics.Processed != 0 {
		t.Errorf("expected Processed 0, got %d", metrics.Processed)
	}
}

func TestSimpleMetrics_RecordEvent_ByEventType(t *testing.T) {
	metrics := NewSimpleMetrics()

	// Record multiple events of different types
	metrics.RecordEvent(EventTypeDisputeStatusChange, true, false, 50*time.Millisecond, nil)
	metrics.RecordEvent(EventTypeDisputeStatusChange, true, false, 30*time.Millisecond, nil)
	metrics.RecordEvent(EventTypeStatementClosed, true, false, 40*time.Millisecond, nil)
	metrics.RecordEvent(EventTypePaymentDue, false, false, 100*time.Millisecond, errors.New("failed"))

	// Check dispute metrics
	disputeMetrics := metrics.ByEventType[EventTypeDisputeStatusChange]
	if disputeMetrics == nil {
		t.Fatal("Expected dispute metrics to exist")
	}
	if disputeMetrics.Total != 2 {
		t.Errorf("expected dispute Total 2, got %d", disputeMetrics.Total)
	}
	if disputeMetrics.Processed != 2 {
		t.Errorf("expected dispute Processed 2, got %d", disputeMetrics.Processed)
	}
	if disputeMetrics.TotalMs != 80 {
		t.Errorf("expected dispute TotalMs 80, got %d", disputeMetrics.TotalMs)
	}

	// Check statement metrics
	statementMetrics := metrics.ByEventType[EventTypeStatementClosed]
	if statementMetrics == nil {
		t.Fatal("Expected statement metrics to exist")
	}
	if statementMetrics.Total != 1 {
		t.Errorf("expected statement Total 1, got %d", statementMetrics.Total)
	}

	// Check payment metrics (with error)
	paymentMetrics := metrics.ByEventType[EventTypePaymentDue]
	if paymentMetrics == nil {
		t.Fatal("Expected payment metrics to exist")
	}
	if paymentMetrics.Errors != 1 {
		t.Errorf("expected payment Errors 1, got %d", paymentMetrics.Errors)
	}
}

func TestSimpleMetrics_Stats(t *testing.T) {
	metrics := NewSimpleMetrics()

	// Record various events
	metrics.RecordEvent(EventTypeDisputeStatusChange, true, false, 100*time.Millisecond, nil)
	metrics.RecordEvent(EventTypeStatementClosed, false, true, 50*time.Millisecond, nil)
	metrics.RecordEvent(EventTypePaymentDue, false, false, 150*time.Millisecond, errors.New("err"))

	total, processed, duplicates, errs, avgMs := metrics.Stats()

	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if processed != 1 {
		t.Errorf("expected processed 1, got %d", processed)
	}
	if duplicates != 1 {
		t.Errorf("expected duplicates 1, got %d", duplicates)
	}
	if errs != 1 {
		t.Errorf("expected errors 1, got %d", errs)
	}
	if avgMs != 100 { // (100+50+150)/3 = 100
		t.Errorf("expected avgMs 100, got %d", avgMs)
	}
}

func TestSimpleMetrics_Stats_Empty(t *testing.T) {
	metrics := NewSimpleMetrics()

	total, processed, duplicates, errs, avgMs := metrics.Stats()

	if total != 0 || processed != 0 || duplicates != 0 || errs != 0 || avgMs != 0 {
		t.Error("Expected all stats to be 0 for empty metrics")
	}
}

func TestSimpleMetrics_Reset(t *testing.T) {
	metrics := NewSimpleMetrics()

	// Record some events
	metrics.RecordEvent(EventTypeDisputeStatusChange, true, false, 100*time.Millisecond, nil)
	metrics.RecordEvent(EventTypeStatementClosed, true, false, 50*time.Millisecond, nil)

	// Verify we have data
	if metrics.Events != 2 {
		t.Errorf("expected Events 2 before reset, got %d", metrics.Events)
	}

	// Reset
	metrics.Reset()

	// Verify all cleared
	if metrics.Events != 0 {
		t.Errorf("expected Events 0 after reset, got %d", metrics.Events)
	}
	if metrics.Processed != 0 {
		t.Errorf("expected Processed 0 after reset, got %d", metrics.Processed)
	}
	if metrics.Duplicates != 0 {
		t.Errorf("expected Duplicates 0 after reset, got %d", metrics.Duplicates)
	}
	if metrics.Errors != 0 {
		t.Errorf("expected Errors 0 after reset, got %d", metrics.Errors)
	}
	if metrics.TotalMs != 0 {
		t.Errorf("expected TotalMs 0 after reset, got %d", metrics.TotalMs)
	}
	if len(metrics.ByEventType) != 0 {
		t.Errorf("expected empty ByEventType after reset, got %d entries", len(metrics.ByEventType))
	}
}

func TestSimpleMetrics_Concurrency(t *testing.T) {
	metrics := NewSimpleMetrics()
	var wg sync.WaitGroup

	// Simulate concurrent recording
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			eventType := EventTypeDisputeStatusChange
			if i%2 == 0 {
				eventType = EventTypeStatementClosed
			}
			metrics.RecordEvent(eventType, true, false, time.Duration(i)*time.Millisecond, nil)
		}(i)
	}

	wg.Wait()

	total, processed, _, _, _ := metrics.Stats()
	if total != 100 {
		t.Errorf("expected total 100, got %d", total)
	}
	if processed != 100 {
		t.Errorf("expected processed 100, got %d", processed)
	}
}

func TestSimpleMetrics_ConcurrentStatsAndRecord(t *testing.T) {
	metrics := NewSimpleMetrics()
	var wg sync.WaitGroup

	// Concurrent recording
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			metrics.RecordEvent(EventTypeDisputeStatusChange, true, false, 10*time.Millisecond, nil)
		}()
	}

	// Concurrent reading
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			metrics.Stats()
		}()
	}

	wg.Wait()

	// Should not panic and have correct count
	total, _, _, _, _ := metrics.Stats()
	if total != 50 {
		t.Errorf("expected total 50, got %d", total)
	}
}

// =============================================================================
// EventInfo and ProcessingInfo Tests
// =============================================================================

func TestEventInfo_Fields(t *testing.T) {
	event := EventInfo{
		EventType:   EventTypeDisputeStatusChange,
		EventID:     "evt-123",
		BodySize:    512,
		IsDuplicate: true,
	}

	if event.EventType != EventTypeDisputeStatusChange {
		t.Errorf("expected EventType %s, got %s", EventTypeDisputeStatusChange, event.EventType)
	}
	if event.EventID != "evt-123" {
		t.Errorf("expected EventID evt-123, got %s", event.EventID)
	}
	if event.BodySize != 512 {
		t.Errorf("expected BodySize 512, got %d", event.BodySize)
	}
	if !event.IsDuplicate {
		t.Error("expected IsDuplicate true")
	}
}

func TestProcessingInfo_Fields(t *testing.T) {
	err := errors.New("test error")
	info := ProcessingInfo{
		StatusCode: 500,
		Duration:   100 * time.Millisecond,
		Error:      err,
		Processed:  false,
	}

	if info.StatusCode != 500 {
		t.Errorf("expected StatusCode 500, got %d", info.StatusCode)
	}
	if info.Duration != 100*time.Millisecond {
		t.Errorf("expected Duration 100ms, got %v", info.Duration)
	}
	if info.Error != err {
		t.Errorf("expected Error %v, got %v", err, info.Error)
	}
	if info.Processed {
		t.Error("expected Processed false")
	}
}
