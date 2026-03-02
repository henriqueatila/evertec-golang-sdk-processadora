package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// mockWebhookHandler implements Handler for testing
type mockWebhookHandler struct {
	BaseHandler
	mu                        sync.Mutex
	disputeStatusChangeCalled bool
	statementClosedCalled     bool
	paymentDueCalled          bool
	cardStatusChangeCalled    bool
	deviceTokenStatusCalled   bool
	lastDisputeEvent          *DisputeStatusChangeEvent
	lastStatementEvent        *StatementClosedEvent
	lastPaymentDueEvent       *PaymentDueEvent
	lastCardStatusEvent       *CardStatusChangeEvent
	lastDeviceTokenEvent      *DeviceTokenStatusEvent
	returnError               error
}

func (m *mockWebhookHandler) OnDisputeStatusChange(_ context.Context, event *DisputeStatusChangeEvent) error {
	m.mu.Lock()
	m.disputeStatusChangeCalled = true
	m.lastDisputeEvent = event
	m.mu.Unlock()
	return m.returnError
}

func (m *mockWebhookHandler) OnStatementClosed(_ context.Context, event *StatementClosedEvent) error {
	m.mu.Lock()
	m.statementClosedCalled = true
	m.lastStatementEvent = event
	m.mu.Unlock()
	return m.returnError
}

func (m *mockWebhookHandler) OnPaymentDue(_ context.Context, event *PaymentDueEvent) error {
	m.mu.Lock()
	m.paymentDueCalled = true
	m.lastPaymentDueEvent = event
	m.mu.Unlock()
	return m.returnError
}

func (m *mockWebhookHandler) OnCardStatusChange(_ context.Context, event *CardStatusChangeEvent) error {
	m.mu.Lock()
	m.cardStatusChangeCalled = true
	m.lastCardStatusEvent = event
	m.mu.Unlock()
	return m.returnError
}

func (m *mockWebhookHandler) OnDeviceTokenStatus(_ context.Context, event *DeviceTokenStatusEvent) error {
	m.mu.Lock()
	m.deviceTokenStatusCalled = true
	m.lastDeviceTokenEvent = event
	m.mu.Unlock()
	return m.returnError
}

// generateSignature creates a valid HMAC signature for testing
func generateSignature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// TestServer_NewServer tests server creation
func TestServer_NewServer(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{
		Handler:       handler,
		WebhookSecret: "test-secret",
	})

	if server == nil {
		t.Fatal("NewServer returned nil")
	}
	if server.handler != handler {
		t.Error("Server handler not set correctly")
	}
	if server.webhookSecret != "test-secret" {
		t.Errorf("webhookSecret = %s, want test-secret", server.webhookSecret)
	}
}

// TestServer_MethodNotAllowed tests wrong HTTP method
func TestServer_MethodNotAllowed(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{Handler: handler})

	req := httptest.NewRequest(http.MethodGet, "/webhook", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

// TestServer_InvalidSignature tests invalid signature rejection
func TestServer_InvalidSignature(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{
		Handler:       handler,
		WebhookSecret: "test-secret",
	})

	body := []byte(`{"eventType":"DISPUTE_STATUS_CHANGE","data":{"disputeId":"disp123"}}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("X-PaySmart-Signature", "sha256=invalidsignature")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

// TestServer_ValidSignature tests valid signature acceptance
func TestServer_ValidSignature(t *testing.T) {
	handler := &mockWebhookHandler{}
	secret := "test-secret"
	server := NewServer(Config{
		Handler:       handler,
		WebhookSecret: secret,
	})

	body := []byte(`{"eventType":"DISPUTE_STATUS_CHANGE","data":{"disputeId":"disp123","newStatus":{"code":1,"status":"OPENED","description":"Dispute opened"}}}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("X-PaySmart-Signature", generateSignature(body, secret))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}
}

// TestServer_NoSecretConfigured tests webhook without secret configured
func TestServer_NoSecretConfigured(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{
		Handler: handler,
		// No secret configured
	})

	body := []byte(`{"eventType":"DISPUTE_STATUS_CHANGE","data":{"disputeId":"disp123"}}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	// No signature header
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	// Should accept without signature when no secret configured
	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}
}

// TestServer_DisputeStatusChangeEvent tests dispute status change event processing
func TestServer_DisputeStatusChangeEvent(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{Handler: handler})

	body := []byte(`{
		"eventType": "DISPUTE_STATUS_CHANGE",
		"data": {
			"disputeId": "disp123",
			"newStatus": {
				"code": 1,
				"status": "OPENED",
				"description": "Dispute opened by cardholder"
			}
		}
	}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	handler.mu.Lock()
	called := handler.disputeStatusChangeCalled
	event := handler.lastDisputeEvent
	handler.mu.Unlock()

	if !called {
		t.Error("OnDisputeStatusChange was not called")
	}
	if event != nil && event.DisputeID != "disp123" {
		t.Errorf("DisputeID = %s, want disp123", event.DisputeID)
	}
}

// TestServer_StatementClosedEvent tests statement closed event processing
func TestServer_StatementClosedEvent(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{Handler: handler})

	body := []byte(`{
		"eventType": "STATEMENT_CLOSED_NOTIFICATION",
		"data": {
			"accountId": "acc123",
			"title": "Fatura de Janeiro 2024"
		}
	}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	handler.mu.Lock()
	called := handler.statementClosedCalled
	event := handler.lastStatementEvent
	handler.mu.Unlock()

	if !called {
		t.Error("OnStatementClosed was not called")
	}
	if event != nil && event.AccountID != "acc123" {
		t.Errorf("AccountID = %s, want acc123", event.AccountID)
	}
}

// TestServer_PaymentDueEvent tests payment due event processing
func TestServer_PaymentDueEvent(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{Handler: handler})

	body := []byte(`{
		"eventType": "PAYMENT_DUE_NOTIFICATION",
		"data": {
			"accountId": "acc456",
			"title": "Vencimento em 3 dias"
		}
	}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	handler.mu.Lock()
	called := handler.paymentDueCalled
	event := handler.lastPaymentDueEvent
	handler.mu.Unlock()

	if !called {
		t.Error("OnPaymentDue was not called")
	}
	if event != nil && event.AccountID != "acc456" {
		t.Errorf("AccountID = %s, want acc456", event.AccountID)
	}
}

// TestServer_CardStatusChangeEvent tests card status change event processing
func TestServer_CardStatusChangeEvent(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{Handler: handler})

	body := []byte(`{
		"eventType": "CARD_STATUS_CHANGE",
		"data": {
			"old_status": {
				"code": 1,
				"status": "ACTIVE",
				"description": "Card is active"
			},
			"new_status": {
				"code": 2,
				"status": "BLOCKED",
				"description": "Card blocked by issuer"
			},
			"eventId": "evt123",
			"cardIdList": ["card001", "card002"]
		}
	}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	handler.mu.Lock()
	called := handler.cardStatusChangeCalled
	event := handler.lastCardStatusEvent
	handler.mu.Unlock()

	if !called {
		t.Error("OnCardStatusChange was not called")
	}
	if event != nil {
		if len(event.CardIDList) != 2 {
			t.Errorf("CardIDList length = %d, want 2", len(event.CardIDList))
		}
		if event.NewStatus != nil && event.NewStatus.Status != "BLOCKED" {
			t.Errorf("NewStatus.Status = %s, want BLOCKED", event.NewStatus.Status)
		}
	}
}

// TestServer_DeviceTokenStatusEvent tests device token status event processing
func TestServer_DeviceTokenStatusEvent(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{Handler: handler})

	body := []byte(`{
		"eventType": "DEVICE_TOKEN_STATUS_NOTIFICATION",
		"data": {
			"previousActivationStatus": "INACTIVE",
			"previousSuspensionStatus": "NOT_SUSPENDED",
			"previousDeploymentStatus": "NOT_DEPLOYED",
			"newActivationStatus": "ACTIVE",
			"newSuspensionStatus": "NOT_SUSPENDED",
			"newDeploymentStatus": "DEPLOYED",
			"deviceTokenId": "token123",
			"cardReferenceId": "cardref123",
			"wallet": "APPLE_PAY",
			"cardId": "card123",
			"accountId": "acc123"
		}
	}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	handler.mu.Lock()
	called := handler.deviceTokenStatusCalled
	event := handler.lastDeviceTokenEvent
	handler.mu.Unlock()

	if !called {
		t.Error("OnDeviceTokenStatus was not called")
	}
	if event != nil {
		if event.Wallet != "APPLE_PAY" {
			t.Errorf("Wallet = %s, want APPLE_PAY", event.Wallet)
		}
		if event.NewActivationStatus != "ACTIVE" {
			t.Errorf("NewActivationStatus = %s, want ACTIVE", event.NewActivationStatus)
		}
	}
}

// TestServer_Idempotency tests idempotent event handling
func TestServer_Idempotency(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{Handler: handler})

	body := []byte(`{"eventType":"DISPUTE_STATUS_CHANGE","data":{"disputeId":"disp123"}}`)

	// First request
	req1 := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req1.Header.Set("X-PaySmart-Event-Id", "evt-idempotent")
	w1 := httptest.NewRecorder()
	server.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First request: Status code = %d, want %d", w1.Code, http.StatusOK)
	}

	// Second request with same event ID (duplicate)
	req2 := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req2.Header.Set("X-PaySmart-Event-Id", "evt-idempotent")
	w2 := httptest.NewRecorder()
	server.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Second request: Status code = %d, want %d", w2.Code, http.StatusOK)
	}

	// Response should indicate received
	var resp WebhookResponse
	_ = json.Unmarshal(w2.Body.Bytes(), &resp)
	if !resp.Received {
		t.Error("Response.Received = false, want true")
	}
}

// TestServer_InvalidJSON tests invalid JSON handling
func TestServer_InvalidJSON(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{Handler: handler})

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestServer_ResponseFormat tests response format
func TestServer_ResponseFormat(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{Handler: handler})

	body := []byte(`{"eventType":"DISPUTE_STATUS_CHANGE","data":{"disputeId":"disp123"}}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	// Check Content-Type header
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %s, want application/json", contentType)
	}

	// Check response body
	var resp WebhookResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !resp.Received {
		t.Error("Response.Received = false, want true")
	}
}

// TestVerifySignature tests standalone signature verification function
func TestVerifySignature(t *testing.T) {
	secret := "test-secret"
	payload := []byte(`{"eventType":"DISPUTE_STATUS_CHANGE","data":{}}`)

	tests := []struct {
		name      string
		signature string
		want      bool
	}{
		{
			name:      "valid signature",
			signature: generateSignature(payload, secret),
			want:      true,
		},
		{
			name:      "invalid signature",
			signature: "sha256=invalidsignature",
			want:      false,
		},
		{
			name:      "missing prefix",
			signature: "invalidsignature",
			want:      false,
		},
		{
			name:      "empty signature",
			signature: "",
			want:      false,
		},
		{
			name:      "wrong secret",
			signature: generateSignature(payload, "wrong-secret"),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VerifySignature(payload, tt.signature, secret)
			if got != tt.want {
				t.Errorf("VerifySignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseEvent tests event parsing function
func TestParseEvent(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		wantErr  bool
		wantType EventType
	}{
		{
			name:     "valid dispute status change event",
			data:     `{"eventType":"DISPUTE_STATUS_CHANGE","data":{"disputeId":"disp123"}}`,
			wantErr:  false,
			wantType: EventTypeDisputeStatusChange,
		},
		{
			name:     "valid card status change event",
			data:     `{"eventType":"CARD_STATUS_CHANGE","data":{"cardIdList":["card123"]}}`,
			wantErr:  false,
			wantType: EventTypeCardStatusChange,
		},
		{
			name:     "valid statement closed event",
			data:     `{"eventType":"STATEMENT_CLOSED_NOTIFICATION","data":{"accountId":"acc123"}}`,
			wantErr:  false,
			wantType: EventTypeStatementClosed,
		},
		{
			name:     "valid payment due event",
			data:     `{"eventType":"PAYMENT_DUE_NOTIFICATION","data":{"accountId":"acc123"}}`,
			wantErr:  false,
			wantType: EventTypePaymentDue,
		},
		{
			name:     "valid device token status event",
			data:     `{"eventType":"DEVICE_TOKEN_STATUS_NOTIFICATION","data":{"wallet":"APPLE_PAY"}}`,
			wantErr:  false,
			wantType: EventTypeDeviceTokenStatus,
		},
		{
			name:    "invalid JSON",
			data:    "invalid json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventType, _, err := ParseEvent([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && eventType != tt.wantType {
				t.Errorf("EventType = %s, want %s", eventType, tt.wantType)
			}
		})
	}
}

// TestEventType_Values tests EventType constants
func TestEventType_Values(t *testing.T) {
	tests := []struct {
		eventType EventType
		want      string
	}{
		{EventTypeDisputeStatusChange, "DISPUTE_STATUS_CHANGE"},
		{EventTypeStatementClosed, "STATEMENT_CLOSED_NOTIFICATION"},
		{EventTypePaymentDue, "PAYMENT_DUE_NOTIFICATION"},
		{EventTypeCardStatusChange, "CARD_STATUS_CHANGE"},
		{EventTypeDeviceTokenStatus, "DEVICE_TOKEN_STATUS_NOTIFICATION"},
	}

	for _, tt := range tests {
		if string(tt.eventType) != tt.want {
			t.Errorf("EventType = %s, want %s", tt.eventType, tt.want)
		}
	}
}

// TestBaseHandler tests default BaseHandler implementation
func TestBaseHandler(t *testing.T) {
	handler := &BaseHandler{}
	ctx := context.Background()

	// All methods should return nil (no-op)
	if err := handler.OnDisputeStatusChange(ctx, &DisputeStatusChangeEvent{}); err != nil {
		t.Errorf("OnDisputeStatusChange() error = %v, want nil", err)
	}
	if err := handler.OnStatementClosed(ctx, &StatementClosedEvent{}); err != nil {
		t.Errorf("OnStatementClosed() error = %v, want nil", err)
	}
	if err := handler.OnPaymentDue(ctx, &PaymentDueEvent{}); err != nil {
		t.Errorf("OnPaymentDue() error = %v, want nil", err)
	}
	if err := handler.OnCardStatusChange(ctx, &CardStatusChangeEvent{}); err != nil {
		t.Errorf("OnCardStatusChange() error = %v, want nil", err)
	}
	if err := handler.OnDeviceTokenStatus(ctx, &DeviceTokenStatusEvent{}); err != nil {
		t.Errorf("OnDeviceTokenStatus() error = %v, want nil", err)
	}
}

// TestDisputeStatusChangeEvent_Fields tests DisputeStatusChangeEvent struct
func TestDisputeStatusChangeEvent_Fields(t *testing.T) {
	event := DisputeStatusChangeEvent{
		DisputeID: "disp123",
		NewStatus: &StatusInfo{
			Code:        1,
			Status:      "OPENED",
			Description: "Dispute opened",
		},
	}

	if event.DisputeID != "disp123" {
		t.Errorf("DisputeID = %s, want disp123", event.DisputeID)
	}
	if event.NewStatus.Code != 1 {
		t.Errorf("NewStatus.Code = %d, want 1", event.NewStatus.Code)
	}
}

// TestCardStatusChangeEvent_Fields tests CardStatusChangeEvent struct
func TestCardStatusChangeEvent_Fields(t *testing.T) {
	event := CardStatusChangeEvent{
		OldStatus: &StatusInfo{
			Code:        1,
			Status:      "ACTIVE",
			Description: "Card active",
		},
		NewStatus: &StatusInfo{
			Code:        2,
			Status:      "BLOCKED",
			Description: "Card blocked",
		},
		EventID:    "evt456",
		CardIDList: []string{"card001", "card002"},
	}

	if len(event.CardIDList) != 2 {
		t.Errorf("CardIDList length = %d, want 2", len(event.CardIDList))
	}
	if event.NewStatus.Status != "BLOCKED" {
		t.Errorf("NewStatus.Status = %s, want BLOCKED", event.NewStatus.Status)
	}
}

// TestDeviceTokenStatusEvent_Fields tests DeviceTokenStatusEvent struct
func TestDeviceTokenStatusEvent_Fields(t *testing.T) {
	event := DeviceTokenStatusEvent{
		PreviousActivationStatus: "INACTIVE",
		NewActivationStatus:      "ACTIVE",
		DeviceTokenID:            "token123",
		CardReferenceID:          "cardref123",
		Wallet:                   "GOOGLE_PAY",
		CardID:                   "card123",
		AccountID:                "acc123",
	}

	if event.Wallet != "GOOGLE_PAY" {
		t.Errorf("Wallet = %s, want GOOGLE_PAY", event.Wallet)
	}
	if event.NewActivationStatus != "ACTIVE" {
		t.Errorf("NewActivationStatus = %s, want ACTIVE", event.NewActivationStatus)
	}
}

// TestWebhookResponse_JSON tests WebhookResponse serialization
func TestWebhookResponse_JSON(t *testing.T) {
	resp := WebhookResponse{Received: true}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	expected := `{"received":true}`
	if string(data) != expected {
		t.Errorf("JSON = %s, want %s", data, expected)
	}
}

// TestParseDisputeStatusChange tests parsing DISPUTE_STATUS_CHANGE event
func TestParseDisputeStatusChange(t *testing.T) {
	data := []byte(`{"disputeId":"disp123","newStatus":{"code":1,"status":"OPENED","description":"Opened"}}`)

	event, err := ParseDisputeStatusChange(data)
	if err != nil {
		t.Fatalf("ParseDisputeStatusChange() error = %v", err)
	}

	if event.DisputeID != "disp123" {
		t.Errorf("DisputeID = %s, want disp123", event.DisputeID)
	}
}

// TestParseCardStatusChange tests parsing CARD_STATUS_CHANGE event
func TestParseCardStatusChange(t *testing.T) {
	data := []byte(`{"old_status":{"code":1,"status":"ACTIVE","description":"Active"},"new_status":{"code":2,"status":"BLOCKED","description":"Blocked"},"eventId":"evt123","cardIdList":["card001"]}`)

	event, err := ParseCardStatusChange(data)
	if err != nil {
		t.Fatalf("ParseCardStatusChange() error = %v", err)
	}

	if event.EventID != "evt123" {
		t.Errorf("EventID = %s, want evt123", event.EventID)
	}
	if len(event.CardIDList) != 1 {
		t.Errorf("CardIDList length = %d, want 1", len(event.CardIDList))
	}
}

// TestParseStatementClosed tests parsing STATEMENT_CLOSED_NOTIFICATION event
func TestParseStatementClosed(t *testing.T) {
	data := []byte(`{"accountId":"acc123","title":"Fatura Janeiro"}`)

	event, err := ParseStatementClosed(data)
	if err != nil {
		t.Fatalf("ParseStatementClosed() error = %v", err)
	}

	if event.AccountID != "acc123" {
		t.Errorf("AccountID = %s, want acc123", event.AccountID)
	}
	if event.Title != "Fatura Janeiro" {
		t.Errorf("Title = %s, want Fatura Janeiro", event.Title)
	}
}

// TestParsePaymentDue tests parsing PAYMENT_DUE_NOTIFICATION event
func TestParsePaymentDue(t *testing.T) {
	data := []byte(`{"accountId":"acc456","title":"Vencimento em 3 dias"}`)

	event, err := ParsePaymentDue(data)
	if err != nil {
		t.Fatalf("ParsePaymentDue() error = %v", err)
	}

	if event.AccountID != "acc456" {
		t.Errorf("AccountID = %s, want acc456", event.AccountID)
	}
}

// TestParseDeviceTokenStatus tests parsing DEVICE_TOKEN_STATUS_NOTIFICATION event
func TestParseDeviceTokenStatus(t *testing.T) {
	data := []byte(`{"previousActivationStatus":"INACTIVE","newActivationStatus":"ACTIVE","deviceTokenId":"token123","wallet":"SAMSUNG_PAY","cardId":"card123","accountId":"acc123"}`)

	event, err := ParseDeviceTokenStatus(data)
	if err != nil {
		t.Fatalf("ParseDeviceTokenStatus() error = %v", err)
	}

	if event.Wallet != "SAMSUNG_PAY" {
		t.Errorf("Wallet = %s, want SAMSUNG_PAY", event.Wallet)
	}
	if event.NewActivationStatus != "ACTIVE" {
		t.Errorf("NewActivationStatus = %s, want ACTIVE", event.NewActivationStatus)
	}
}

// TestStatusInfo_JSON tests StatusInfo serialization
func TestStatusInfo_JSON(t *testing.T) {
	status := StatusInfo{
		Code:        1,
		Status:      "ACTIVE",
		Description: "Card is active",
	}

	data, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var parsed StatusInfo
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if parsed.Code != 1 {
		t.Errorf("Code = %d, want 1", parsed.Code)
	}
	if parsed.Status != "ACTIVE" {
		t.Errorf("Status = %s, want ACTIVE", parsed.Status)
	}
}

// TestServer_IdempotencyCleanup tests idempotency TTL cleanup functionality
func TestServer_IdempotencyCleanup(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{
		Handler:                  handler,
		EnableIdempotencyCleanup: true,
		IdempotencyTTL:           50 * time.Millisecond,
	})
	defer server.Close()

	// Add an event to idempotency store
	body := []byte(`{"eventType":"DISPUTE_STATUS_CHANGE","data":{"disputeId":"disp123"}}`)
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("X-PaySmart-Event-Id", "evt-cleanup-test")
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("First request failed: %d", w.Code)
	}

	// Verify it's in the store
	server.mu.RLock()
	_, exists := server.processedIDs["evt-cleanup-test"]
	server.mu.RUnlock()
	if !exists {
		t.Error("Event ID should be in processedIDs")
	}

	// Wait for TTL to expire plus cleanup interval
	time.Sleep(100 * time.Millisecond)

	// Manually trigger cleanup to ensure it runs
	server.cleanupExpired()

	// Verify it's been cleaned up
	server.mu.RLock()
	_, exists = server.processedIDs["evt-cleanup-test"]
	server.mu.RUnlock()
	if exists {
		t.Error("Event ID should have been cleaned up")
	}
}

// TestServer_Close tests graceful server shutdown
func TestServer_Close(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{
		Handler:                  handler,
		EnableIdempotencyCleanup: true,
		IdempotencyTTL:           1 * time.Hour,
	})

	// Close should not panic
	server.Close()

	// Multiple closes should be safe (idempotent)
	server.Close()
}

// TestServer_CloseWithoutCleanup tests Close when cleanup is disabled
func TestServer_CloseWithoutCleanup(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{
		Handler: handler,
		// EnableIdempotencyCleanup is false by default
	})

	// Close should not panic even when cleanup is disabled
	server.Close()
}

// TestServer_MaxBodySize tests request body size limit
func TestServer_MaxBodySize(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{
		Handler:     handler,
		MaxBodySize: 50, // Very small limit for testing
	})

	// Create a body larger than the limit
	largeBody := make([]byte, 100)
	for i := range largeBody {
		largeBody[i] = 'a'
	}

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(largeBody))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	// Should fail due to body size limit
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestServer_MaxEventIDLength tests event ID length validation
func TestServer_MaxEventIDLength(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{
		Handler:          handler,
		MaxEventIDLength: 10, // Short limit for testing
	})

	body := []byte(`{"eventType":"DISPUTE_STATUS_CHANGE","data":{"disputeId":"disp123"}}`)
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("X-PaySmart-Event-Id", "this-is-a-very-long-event-id-that-exceeds-limit")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestServer_DefaultIdempotencyTTL tests default TTL value
func TestServer_DefaultIdempotencyTTL(t *testing.T) {
	handler := &mockWebhookHandler{}
	server := NewServer(Config{
		Handler:                  handler,
		EnableIdempotencyCleanup: true,
		// IdempotencyTTL not set - should use default
	})
	defer server.Close()

	if server.idempotencyTTL != DefaultIdempotencyTTL {
		t.Errorf("idempotencyTTL = %v, want %v", server.idempotencyTTL, DefaultIdempotencyTTL)
	}
}
