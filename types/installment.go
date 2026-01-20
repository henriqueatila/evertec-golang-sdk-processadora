// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// Installment from OpenAPI spec.
type Installment struct {
	InstallmentDate interface{} `json:"installmentDate"`
	Amount          Amount      `json:"amount"`
}

// InstallmentOption from OpenAPI spec.
type InstallmentOption struct {
	NumInstallments   int    `json:"numInstallments"`
	InstallmentAmount Amount `json:"installmentAmount"`
	TotalAmount       Amount `json:"totalAmount"`
	CET               CET    `json:"CET"`
}

// InstallmentOptions from OpenAPI spec.
type InstallmentOptions struct {
	InstallmentOptions []InstallmentOption `json:"installmentOptions"`
	OriginalAmount     Amount              `json:"originalAmount"`
}

// InstallmentPurchase from OpenAPI spec.
type InstallmentPurchase struct {
	TransactionID          string        `json:"transactionId"`
	CardID                 string        `json:"cardId,omitempty"`
	LastFourDigits         string        `json:"last_four_digits"`
	TransactionDescription string        `json:"transactionDescription"`
	TransactionDate        interface{}   `json:"transactionDate"`
	Installments           []Installment `json:"installments"`
}

// AllInstallmentPurchase from OpenAPI spec.
type AllInstallmentPurchase struct {
	InstallmentPurchases []InstallmentPurchase `json:"installmentPurchases,omitempty"`
}

// InstallmentRequest from OpenAPI spec.
type InstallmentRequest struct {
	NumberOfInstallments int     `json:"numberOfInstallments"`
	InstallmentAmount    *Amount `json:"installmentAmount,omitempty"`
}

// installmentSimulationBody from OpenAPI spec.
type installmentSimulationBody struct {
	Amount          int `json:"amount,omitempty"`
	DownPayment     int `json:"downPayment,omitempty"`
	NumInstallments int `json:"numInstallments,omitempty"`
}

// installmentSimulationResult from OpenAPI spec.
type installmentSimulationResult struct {
	ResultData interface{}      `json:"resultData,omitempty"`
	Data       *installmentData `json:"data,omitempty"`
}

// installmentSimulationResultError from OpenAPI spec.
type installmentSimulationResultError struct {
	ResultData interface{} `json:"resultData,omitempty"`
}

// installmentData from OpenAPI spec.
type installmentData struct {
	TotalAmount          int `json:"totalAmount,omitempty"`
	DownPayment          int `json:"downPayment,omitempty"`
	NumInstallments      int `json:"numInstallments,omitempty"`
	AmountPerInstallment int `json:"amountPerInstallment,omitempty"`
	OriginalAmount       int `json:"originalAmount,omitempty"`
	MonthlyInterest      int `json:"monthlyInterest,omitempty"`
	TotalIOF             int `json:"totalIOF,omitempty"`
}

// AdvancePaymentRequest from OpenAPI spec.
type AdvancePaymentRequest struct {
	TransactionID         string `json:"transactionId"`
	InstallmentsToAdvance int    `json:"installmentsToAdvance"`
}
