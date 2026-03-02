package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultIdempotencyTTL is the default time-to-live for idempotency entries
	// when EnableIdempotencyCleanup is true.
	DefaultIdempotencyTTL = 24 * time.Hour
)

// Handler defines the interface for webhook handlers.
// Implement the methods for the event types you want to handle.
type Handler interface {
	// OnDisputeStatusChange is called when a dispute status changes.
	// EventType: DISPUTE_STATUS_CHANGE
	OnDisputeStatusChange(ctx context.Context, event *DisputeStatusChangeEvent) error

	// OnStatementClosed is called when an invoice/statement closes.
	// EventType: STATEMENT_CLOSED_NOTIFICATION
	// Only for post-paid accounts.
	OnStatementClosed(ctx context.Context, event *StatementClosedEvent) error

	// OnPaymentDue is called when a payment deadline approaches.
	// EventType: PAYMENT_DUE_NOTIFICATION
	// Only for post-paid accounts.
	OnPaymentDue(ctx context.Context, event *PaymentDueEvent) error

	// OnCardStatusChange is called when card status transitions.
	// EventType: CARD_STATUS_CHANGE
	OnCardStatusChange(ctx context.Context, event *CardStatusChangeEvent) error

	// OnDeviceTokenStatus is called when digital wallet token status changes.
	// EventType: DEVICE_TOKEN_STATUS_NOTIFICATION
	OnDeviceTokenStatus(ctx context.Context, event *DeviceTokenStatusEvent) error
}

// idempotencyEntry stores an event ID with its timestamp for TTL cleanup.
type idempotencyEntry struct {
	processedAt time.Time
}

// Server wraps a Handler and provides HTTP handling for webhooks.
//
// The server includes built-in idempotency handling using an in-memory store.
// Note: The default in-memory idempotency store is suitable for development
// and single-instance deployments. For production environments with multiple
// instances or long-running processes, consider implementing your own
// idempotency logic in the Handler using an external store (e.g., Redis,
// database) to track processed event IDs with appropriate TTL.
type Server struct {
	handler        Handler
	webhookSecret  string
	hooks          []Hook
	logger         *slog.Logger
	mu             sync.RWMutex
	processedIDs   map[string]idempotencyEntry
	idempotencyTTL time.Duration
	maxBodySize    int64
	maxEventIDLen  int
	cleanupEnabled bool
	stopCleanup    chan struct{}
	closeOnce      sync.Once
}

// Config holds webhook server configuration.
type Config struct {
	// Handler is the webhook handler implementation.
	Handler Handler

	// WebhookSecret is the secret for verifying webhook signatures.
	// If empty, signature verification is skipped.
	WebhookSecret string

	// Hooks are observability hooks for logging, metrics, tracing, etc. (optional)
	Hooks []Hook

	// Logger is used for internal server logging (e.g., hook panic recovery).
	// If nil, no internal logging is performed.
	Logger *slog.Logger

	// MaxBodySize limits the request body size in bytes.
	// If 0 (default), no limit is applied.
	// Example: 10 * 1024 * 1024 for 10MB limit.
	MaxBodySize int64

	// MaxEventIDLength limits the event ID length.
	// If 0 (default), no limit is applied.
	// Example: 128 for 128 character limit.
	MaxEventIDLength int

	// EnableIdempotencyCleanup enables automatic cleanup of expired idempotency entries.
	// When enabled, a background goroutine periodically removes entries older than IdempotencyTTL.
	// You must call Close() to stop the cleanup goroutine when done.
	// If false (default), idempotency entries are kept indefinitely (suitable for short-lived servers).
	EnableIdempotencyCleanup bool

	// IdempotencyTTL is the time-to-live for idempotency entries.
	// Only used when EnableIdempotencyCleanup is true.
	// Defaults to 24 hours if not set.
	IdempotencyTTL time.Duration
}

// NewServer creates a new webhook server.
//
// The server provides built-in idempotency handling to prevent duplicate
// event processing. Event IDs received via the X-PaySmart-Event-Id header
// are tracked in memory.
//
// If EnableIdempotencyCleanup is true, a background goroutine will periodically
// remove expired entries and you must call Close() when done.
// See Server documentation for production considerations.
func NewServer(cfg Config) *Server {
	ttl := cfg.IdempotencyTTL
	if ttl == 0 {
		ttl = DefaultIdempotencyTTL
	}

	s := &Server{
		handler:        cfg.Handler,
		webhookSecret:  cfg.WebhookSecret,
		hooks:          cfg.Hooks,
		logger:         cfg.Logger,
		processedIDs:   make(map[string]idempotencyEntry),
		idempotencyTTL: ttl,
		maxBodySize:    cfg.MaxBodySize,
		maxEventIDLen:  cfg.MaxEventIDLength,
		cleanupEnabled: cfg.EnableIdempotencyCleanup,
	}

	// Only start cleanup goroutine if explicitly enabled
	if cfg.EnableIdempotencyCleanup {
		s.stopCleanup = make(chan struct{})
		go s.cleanupLoop()
	}

	return s
}

// Close stops the cleanup goroutine and releases resources.
// Only required if EnableIdempotencyCleanup was set to true in Config.
func (s *Server) Close() {
	s.closeOnce.Do(func() {
		if s.cleanupEnabled && s.stopCleanup != nil {
			close(s.stopCleanup)
		}
	})
}

// cleanupLoop periodically removes expired idempotency entries.
func (s *Server) cleanupLoop() {
	ticker := time.NewTicker(s.idempotencyTTL / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupExpired()
		case <-s.stopCleanup:
			return
		}
	}
}

// cleanupExpired removes idempotency entries older than the TTL.
func (s *Server) cleanupExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-s.idempotencyTTL)
	for id, entry := range s.processedIDs {
		if entry.processedAt.Before(cutoff) {
			delete(s.processedIDs, id)
		}
	}
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.handler == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "handler not configured"})
		return
	}

	start := time.Now()
	ctx := r.Context()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Apply request body size limit if configured
	if s.maxBodySize > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, s.maxBodySize)
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// Get event ID from header
	eventID := r.Header.Get("X-PaySmart-Event-Id")

	// Validate event ID length if configured
	if s.maxEventIDLen > 0 && len(eventID) > s.maxEventIDLen {
		http.Error(w, "Event ID too long", http.StatusBadRequest)
		return
	}

	// Verify signature if secret is configured
	if s.webhookSecret != "" {
		signature := r.Header.Get("X-PaySmart-Signature")
		if !s.verifySignature(body, signature) {
			eventInfo := &EventInfo{
				EventID:  eventID,
				BodySize: len(body),
			}
			ctx = s.executeBeforeHooks(ctx, eventInfo)
			s.executeAfterHooks(ctx, eventInfo, &ProcessingInfo{
				StatusCode: http.StatusUnauthorized,
				Duration:   time.Since(start),
				Error:      fmt.Errorf("invalid signature"),
			})
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}

	// Check idempotency (thread-safe with double-checked locking)
	skipProcessing := false
	if eventID != "" {
		// Fast path: try to detect duplicate with read lock
		s.mu.RLock()
		entry, exists := s.processedIDs[eventID]
		if exists {
			if s.cleanupEnabled {
				skipProcessing = time.Since(entry.processedAt) < s.idempotencyTTL
			} else {
				skipProcessing = true // No TTL, all entries are valid
			}
		}
		s.mu.RUnlock()

		if skipProcessing {
			// Already processed, return success
			eventInfo := &EventInfo{
				EventID:     eventID,
				BodySize:    len(body),
				IsDuplicate: true,
			}
			ctx = s.executeBeforeHooks(ctx, eventInfo)
			s.executeAfterHooks(ctx, eventInfo, &ProcessingInfo{
				StatusCode: http.StatusOK,
				Duration:   time.Since(start),
				Processed:  true,
			})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(WebhookResponse{Received: true})
			return
		}

		// Slow path: acquire write lock and double-check before marking as processed
		s.mu.Lock()
		entry, exists = s.processedIDs[eventID]
		if exists && (!s.cleanupEnabled || time.Since(entry.processedAt) < s.idempotencyTTL) {
			// Another request processed this event while we were waiting for the lock
			skipProcessing = true
		}
		s.mu.Unlock()

		if skipProcessing {
			// Duplicate detected under lock, return success
			eventInfo := &EventInfo{
				EventID:     eventID,
				BodySize:    len(body),
				IsDuplicate: true,
			}
			ctx = s.executeBeforeHooks(ctx, eventInfo)
			s.executeAfterHooks(ctx, eventInfo, &ProcessingInfo{
				StatusCode: http.StatusOK,
				Duration:   time.Since(start),
				Processed:  true,
			})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(WebhookResponse{Received: true})
			return
		}
	}

	// Parse event type first
	var eventWrapper struct {
		EventType EventType       `json:"eventType"`
		EventID   string          `json:"eventId,omitempty"`
		Data      json.RawMessage `json:"data,omitempty"`
	}
	if err := json.Unmarshal(body, &eventWrapper); err != nil {
		eventInfo := &EventInfo{
			EventID:  eventID,
			BodySize: len(body),
		}
		ctx = s.executeBeforeHooks(ctx, eventInfo)
		s.executeAfterHooks(ctx, eventInfo, &ProcessingInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		http.Error(w, "Failed to parse event", http.StatusBadRequest)
		return
	}

	// Create event info for hooks
	eventInfo := &EventInfo{
		EventType: eventWrapper.EventType,
		EventID:   eventID,
		BodySize:  len(body),
	}

	// Execute before hooks
	ctx = s.executeBeforeHooks(ctx, eventInfo)

	// Process event based on type
	// Errors are intentionally ignored - we return success to prevent retries.
	// The error is passed to hooks for logging/metrics.
	processErr := s.processEvent(ctx, eventWrapper.EventType, eventWrapper.Data)

	// Mark as processed (thread-safe with double-checked locking)
	// Only mark if we got past the idempotency checks (eventID != "" and no duplicate detected)
	if eventID != "" {
		s.mu.Lock()
		// Double-check: another goroutine might have processed this event
		entry, exists := s.processedIDs[eventID]
		if exists && (!s.cleanupEnabled || time.Since(entry.processedAt) < s.idempotencyTTL) {
			// Already processed, do not overwrite
			s.mu.Unlock()
		} else {
			// Mark as processed
			s.processedIDs[eventID] = idempotencyEntry{processedAt: time.Now()}
			s.mu.Unlock()
		}
	}

	// Execute after hooks
	s.executeAfterHooks(ctx, eventInfo, &ProcessingInfo{
		StatusCode: http.StatusOK,
		Duration:   time.Since(start),
		Error:      processErr,
		Processed:  processErr == nil,
	})

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(WebhookResponse{Received: true})
}

// executeBeforeHooks runs all before hooks and returns the modified context.
func (s *Server) executeBeforeHooks(ctx context.Context, event *EventInfo) context.Context {
	for _, hook := range s.hooks {
		func() {
			defer func() {
				if r := recover(); r != nil && s.logger != nil {
					s.logger.Error("hook panic recovered",
						slog.Any("panic", r),
						slog.String("hook_phase", "BeforeEvent"))
				}
			}()
			ctx = hook.BeforeEvent(ctx, event)
		}()
	}
	return ctx
}

// executeAfterHooks runs all after hooks.
func (s *Server) executeAfterHooks(ctx context.Context, event *EventInfo, result *ProcessingInfo) {
	for _, hook := range s.hooks {
		func() {
			defer func() {
				if r := recover(); r != nil && s.logger != nil {
					s.logger.Error("hook panic recovered",
						slog.Any("panic", r),
						slog.String("hook_phase", "AfterEvent"))
				}
			}()
			hook.AfterEvent(ctx, event, result)
		}()
	}
}

func (s *Server) verifySignature(payload []byte, signature string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	mac := hmac.New(sha256.New, []byte(s.webhookSecret))
	mac.Write(payload)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

func (s *Server) processEvent(ctx context.Context, eventType EventType, data json.RawMessage) error {
	switch eventType {
	case EventTypeDisputeStatusChange:
		var event DisputeStatusChangeEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("failed to parse dispute status change event: %w", err)
		}
		return s.handler.OnDisputeStatusChange(ctx, &event)

	case EventTypeStatementClosed:
		var event StatementClosedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("failed to parse statement closed event: %w", err)
		}
		return s.handler.OnStatementClosed(ctx, &event)

	case EventTypePaymentDue:
		var event PaymentDueEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("failed to parse payment due event: %w", err)
		}
		return s.handler.OnPaymentDue(ctx, &event)

	case EventTypeCardStatusChange:
		var event CardStatusChangeEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("failed to parse card status change event: %w", err)
		}
		return s.handler.OnCardStatusChange(ctx, &event)

	case EventTypeDeviceTokenStatus:
		var event DeviceTokenStatusEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("failed to parse device token status event: %w", err)
		}
		return s.handler.OnDeviceTokenStatus(ctx, &event)

	default:
		return fmt.Errorf("unknown event type: %s", eventType)
	}
}

// BaseHandler provides a default implementation of Handler.
// Embed this in your handler to only implement the methods you need.
type BaseHandler struct{}

func (h *BaseHandler) OnDisputeStatusChange(_ context.Context, _ *DisputeStatusChangeEvent) error {
	return nil
}

func (h *BaseHandler) OnStatementClosed(_ context.Context, _ *StatementClosedEvent) error {
	return nil
}

func (h *BaseHandler) OnPaymentDue(_ context.Context, _ *PaymentDueEvent) error {
	return nil
}

func (h *BaseHandler) OnCardStatusChange(_ context.Context, _ *CardStatusChangeEvent) error {
	return nil
}

func (h *BaseHandler) OnDeviceTokenStatus(_ context.Context, _ *DeviceTokenStatusEvent) error {
	return nil
}

// VerifySignature verifies a webhook signature manually.
func VerifySignature(payload []byte, signature, secret string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

// ParseEvent parses a webhook event from JSON.
// Returns the event type and raw data for further processing.
func ParseEvent(data []byte) (EventType, json.RawMessage, error) {
	var event struct {
		EventType EventType       `json:"eventType"`
		Data      json.RawMessage `json:"data,omitempty"`
	}
	if err := json.Unmarshal(data, &event); err != nil {
		return "", nil, fmt.Errorf("failed to parse event: %w", err)
	}
	return event.EventType, event.Data, nil
}

// ParseDisputeStatusChange parses a DISPUTE_STATUS_CHANGE event.
func ParseDisputeStatusChange(data json.RawMessage) (*DisputeStatusChangeEvent, error) {
	var event DisputeStatusChangeEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to parse dispute status change: %w", err)
	}
	return &event, nil
}

// ParseStatementClosed parses a STATEMENT_CLOSED_NOTIFICATION event.
func ParseStatementClosed(data json.RawMessage) (*StatementClosedEvent, error) {
	var event StatementClosedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to parse statement closed: %w", err)
	}
	return &event, nil
}

// ParsePaymentDue parses a PAYMENT_DUE_NOTIFICATION event.
func ParsePaymentDue(data json.RawMessage) (*PaymentDueEvent, error) {
	var event PaymentDueEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to parse payment due: %w", err)
	}
	return &event, nil
}

// ParseCardStatusChange parses a CARD_STATUS_CHANGE event.
func ParseCardStatusChange(data json.RawMessage) (*CardStatusChangeEvent, error) {
	var event CardStatusChangeEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to parse card status change: %w", err)
	}
	return &event, nil
}

// ParseDeviceTokenStatus parses a DEVICE_TOKEN_STATUS_NOTIFICATION event.
func ParseDeviceTokenStatus(data json.RawMessage) (*DeviceTokenStatusEvent, error) {
	var event DeviceTokenStatusEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to parse device token status: %w", err)
	}
	return &event, nil
}
