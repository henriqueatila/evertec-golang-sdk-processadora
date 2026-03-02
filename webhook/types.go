// Package webhook provides types and handlers for EventHub webhooks.
// EventHub pushes real-time notifications from Evertec to your endpoint.
// Reference: https://paysmart-api.gitlab.io/processadora/PT-br/docs/tipos-e-configs
package webhook

// ========================================
// Event Types
// ========================================

// EventType represents the type of webhook event from EventHub.
type EventType string

const (
	// EventTypeDisputeStatusChange notifies about dispute status changes.
	EventTypeDisputeStatusChange EventType = "DISPUTE_STATUS_CHANGE"

	// EventTypeStatementClosed alerts when invoices close (post-paid accounts).
	EventTypeStatementClosed EventType = "STATEMENT_CLOSED_NOTIFICATION"

	// EventTypePaymentDue notifies about invoice expiration deadlines (post-paid accounts).
	EventTypePaymentDue EventType = "PAYMENT_DUE_NOTIFICATION"

	// EventTypeCardStatusChange alerts when card status transitions occur.
	EventTypeCardStatusChange EventType = "CARD_STATUS_CHANGE"

	// EventTypeDeviceTokenStatus notifies when digital wallet token status changes.
	EventTypeDeviceTokenStatus EventType = "DEVICE_TOKEN_STATUS_NOTIFICATION"
)

// ========================================
// Common Types
// ========================================

// StatusInfo represents a status object with code, status, and description.
// Used in CARD_STATUS_CHANGE and DISPUTE_STATUS_CHANGE events.
type StatusInfo struct {
	Code        int    `json:"code"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// ========================================
// Event Wrapper
// ========================================

// Event represents a webhook event from EventHub.
// The Data field contains event-specific payload based on EventType.
type Event struct {
	EventType EventType   `json:"eventType"`
	EventID   string      `json:"eventId,omitempty"`
	Data      any `json:"data,omitempty"`
}

// ========================================
// DISPUTE_STATUS_CHANGE
// ========================================

// DisputeStatusChangeEvent represents a dispute status change notification.
// EventType: DISPUTE_STATUS_CHANGE
type DisputeStatusChangeEvent struct {
	DisputeID string      `json:"disputeId"`
	NewStatus *StatusInfo `json:"newStatus"`
}

// ========================================
// STATEMENT_CLOSED_NOTIFICATION
// ========================================

// StatementClosedEvent represents a statement/invoice closed notification.
// EventType: STATEMENT_CLOSED_NOTIFICATION
// Only for post-paid accounts.
type StatementClosedEvent struct {
	AccountID string `json:"accountId"`
	Title     string `json:"title"`
}

// ========================================
// PAYMENT_DUE_NOTIFICATION
// ========================================

// PaymentDueEvent represents a payment due notification.
// EventType: PAYMENT_DUE_NOTIFICATION
// Only for post-paid accounts.
type PaymentDueEvent struct {
	AccountID string `json:"accountId"`
	Title     string `json:"title"`
}

// ========================================
// CARD_STATUS_CHANGE
// ========================================

// CardStatusChangeEvent represents a card status change notification.
// EventType: CARD_STATUS_CHANGE
type CardStatusChangeEvent struct {
	OldStatus  *StatusInfo `json:"old_status"`
	NewStatus  *StatusInfo `json:"new_status"`
	EventID    string      `json:"eventId"`
	CardIDList []string    `json:"cardIdList"`
}

// ========================================
// DEVICE_TOKEN_STATUS_NOTIFICATION
// ========================================

// DeviceTokenStatusEvent represents a digital wallet token status change.
// EventType: DEVICE_TOKEN_STATUS_NOTIFICATION
// Notifies when Apple Pay, Google Pay, Samsung Pay token status changes.
type DeviceTokenStatusEvent struct {
	// Previous status fields
	PreviousActivationStatus string `json:"previousActivationStatus"`
	PreviousSuspensionStatus string `json:"previousSuspensionStatus"`
	PreviousDeploymentStatus string `json:"previousDeploymentStatus"`

	// New status fields
	NewActivationStatus string `json:"newActivationStatus"`
	NewSuspensionStatus string `json:"newSuspensionStatus"`
	NewDeploymentStatus string `json:"newDeploymentStatus"`

	// Token and card identifiers
	DeviceTokenID   string `json:"deviceTokenId"`
	CardReferenceID string `json:"cardReferenceId"`
	Wallet          string `json:"wallet"` // APPLE_PAY, GOOGLE_PAY, SAMSUNG_PAY
	CardID          string `json:"cardId"`
	AccountID       string `json:"accountId"`
}

// ========================================
// Response
// ========================================

// WebhookResponse represents the expected response to a webhook.
type WebhookResponse struct {
	Received bool `json:"received"`
}
