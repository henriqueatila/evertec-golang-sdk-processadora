// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// Balance from OpenAPI spec.
type Balance struct {
	CurrentBalance *Amount `json:"currentBalance,omitempty"`
}

// CreditMattress from OpenAPI spec.
type CreditMattress struct {
	TransactionID          string              `json:"transactionId"`
	CreditReceivedDateTime string              `json:"creditReceivedDateTime"`
	CreditDate             string              `json:"creditDate,omitempty"`
	CreditTime             string              `json:"creditTime,omitempty"`
	Amount                 DebitOrCreditAmount `json:"amount"`
	Status                 string              `json:"status"`
	Cancellation           bool                `json:"cancellation,omitempty"`
}

// CreditListMattress from OpenAPI spec.
type CreditListMattress struct {
	Credits        []CreditMattress `json:"credits,omitempty"`
	HasMoreCredits bool             `json:"hasMoreCredits,omitempty"`
}
