// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// Dispute from OpenAPI spec.
type Dispute struct {
	DisputeID              string         `json:"disputeId,omitempty"`
	DisputeRequest         DisputeRequest `json:"disputeRequest"`
	DisputeDate            string         `json:"disputeDate"`
	DisputeMaxDateResponse any    `json:"disputeMaxDateResponse,omitempty"`
	DisputeType            DisputeType    `json:"disputeType"`
	DisputeStatus          DisputeStatus  `json:"disputeStatus"`
	CurrentStage           string         `json:"currentStage"`
	FraudType              string         `json:"fraudType,omitempty"`
	DisputeHistory         []DisputeEvent `json:"dispute_history,omitempty"`
	Transaction            *Transaction   `json:"transaction,omitempty"`
}

// DisputeListResult from OpenAPI spec.
type DisputeListResult struct {
	HasMore  bool      `json:"hasMore,omitempty"`
	Disputes []Dispute `json:"disputes,omitempty"`
}

// DisputeRequest from OpenAPI spec.
type DisputeRequest struct {
	IssuerDisputeID    string       `json:"issuerDisputeId,omitempty"`
	AccountID          string       `json:"accountId"`
	TransactionID      string       `json:"transactionId"`
	DisputeCode        string       `json:"disputeCode"`
	DisputeTextMessage string       `json:"disputeTextMessage"`
	FraudType          string       `json:"fraudType,omitempty"`
	Partial            bool         `json:"partial,omitempty"`
	AmountDisputed     any  `json:"amount_disputed,omitempty"`
	WillAddDocuments   bool         `json:"willAddDocuments,omitempty"`
	SourceAudit        *SourceAudit `json:"sourceAudit,omitempty"`
}

// DisputeCreatedSuccessfully from OpenAPI spec.
type DisputeCreatedSuccessfully struct {
	ResultData any `json:"resultData,omitempty"`
	Dispute    *Dispute    `json:"dispute,omitempty"`
}

// DisputeInfo from OpenAPI spec.
type DisputeInfo struct {
	DisputeType   DisputeRequest `json:"disputeType"`
	DisputeReason DisputeReason  `json:"disputeReason"`
	DisputeStatus DisputeStatus  `json:"disputeStatus"`
}

// DisputeEvent from OpenAPI spec.
type DisputeEvent struct {
	SentFrom    string      `json:"sent_from,omitempty"`
	DateSent    any `json:"date_sent,omitempty"`
	Description string      `json:"description,omitempty"`
}

// DisputeDocument from OpenAPI spec.
type DisputeDocument struct {
	DocumentName        string `json:"documentName"`
	DocumentDescription string `json:"documentDescription,omitempty"`
	Document            string `json:"document"`
}

// DisputeDocumentRequest from OpenAPI spec.
type DisputeDocumentRequest struct {
	IssuerDisputeDocumentID string          `json:"issuerDisputeDocumentId,omitempty"`
	Document                DisputeDocument `json:"document"`
	SourceAudit             *SourceAudit    `json:"sourceAudit,omitempty"`
}

// DisputeDocumentCreatedSuccessfully from OpenAPI spec.
type DisputeDocumentCreatedSuccessfully struct {
	ResultData any `json:"resultData,omitempty"`
	Dispute    *Dispute    `json:"dispute,omitempty"`
}

// DisputeResponseRequest from OpenAPI spec.
type DisputeResponseRequest struct {
	IssuerDisputeResponseID    string       `json:"issuerDisputeResponseId,omitempty"`
	Accept                     bool         `json:"accept"`
	DisputeResponseTextMessage string       `json:"disputeResponseTextMessage"`
	WillAddDocuments           bool         `json:"willAddDocuments,omitempty"`
	SourceAudit                *SourceAudit `json:"sourceAudit,omitempty"`
}

// DisputeResponseCreatedSuccessfully from OpenAPI spec.
type DisputeResponseCreatedSuccessfully struct {
	ResultData any `json:"resultData,omitempty"`
	Dispute    *Dispute    `json:"dispute,omitempty"`
}

// DisputeReversalRequest from OpenAPI spec.
type DisputeReversalRequest struct {
	IssuerDisputeReversalID string       `json:"issuerDisputeReversalId,omitempty"`
	TextMessage             string       `json:"textMessage,omitempty"`
	SourceAudit             *SourceAudit `json:"sourceAudit,omitempty"`
}

// DisputeReversalCreatedSuccessfully from OpenAPI spec.
type DisputeReversalCreatedSuccessfully struct {
	ResultData any `json:"resultData,omitempty"`
	Dispute    *Dispute    `json:"dispute,omitempty"`
}

// UpdateDisputeStatusRequest from OpenAPI spec.
type UpdateDisputeStatusRequest struct {
	NewStatus                   string       `json:"newStatus"`
	IssuerUpdateDisputeStatusID string       `json:"issuerUpdateDisputeStatusId,omitempty"`
	SourceAudit                 *SourceAudit `json:"sourceAudit,omitempty"`
}

// DisputeStatusUpdatedSuccessfully from OpenAPI spec.
type DisputeStatusUpdatedSuccessfully struct {
	ResultData any `json:"resultData,omitempty"`
}
