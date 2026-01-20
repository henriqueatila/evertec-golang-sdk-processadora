// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// InclusiveTransaction from OpenAPI spec.
type InclusiveTransaction struct {
	InclusiveTransactionID string `json:"inclusiveTransactionId,omitempty"`
	AccountID              string `json:"accountId,omitempty"`
	TransactionID          string `json:"transactionId,omitempty"`
	Code                   string `json:"code,omitempty"`
	ReasonCode             string `json:"reasonCode,omitempty"`
	Text                   string `json:"text,omitempty"`
	CreatedAt              string `json:"createdAt,omitempty"`
}

// InclusiveTransactionRequest from OpenAPI spec.
type InclusiveTransactionRequest struct {
	AccountID     string      `json:"accountId"`
	TransactionID string      `json:"transactionId"`
	Code          string      `json:"code"`
	ReasonCode    string      `json:"reasonCode"`
	Text          string      `json:"text"`
	Partial       bool        `json:"partial"`
	Amount        interface{} `json:"amount,omitempty"`
}

// InclusiveTransactionCreationSuccess from OpenAPI spec.
type InclusiveTransactionCreationSuccess struct {
	ResultData           interface{} `json:"resultData,omitempty"`
	InclusiveTransaction interface{} `json:"inclusiveTransaction,omitempty"`
}

// InclusiveTransactionUndoSuccess from OpenAPI spec.
type InclusiveTransactionUndoSuccess struct {
	ResultData interface{} `json:"resultData,omitempty"`
}

// InclusiveTransactionsListResult from OpenAPI spec.
type InclusiveTransactionsListResult struct {
	HasMore               bool                   `json:"hasMore,omitempty"`
	InclusiveTransactions []InclusiveTransaction `json:"inclusiveTransactions,omitempty"`
}

// UndoInclusiveTransactionRequest from OpenAPI spec.
type UndoInclusiveTransactionRequest struct {
	InclusiveTransactionID string `json:"inclusiveTransactionId"`
}

// UndoData - Dados sobre o estorno.
type UndoData struct {
	PartialUndo     bool   `json:"partialUndo,omitempty"`
	UndoDescription string `json:"undoDescription,omitempty"`
}
