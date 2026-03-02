// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// FraudNotification from OpenAPI spec.
type FraudNotification struct {
	CardHolderData          *FraudNotificationCardholderData `json:"cardHolderData,omitempty"`
	NotificationDate        string                           `json:"notificationDate,omitempty"`
	FraudType               string                           `json:"fraudType,omitempty"`
	FraudNotificationStatus *FraudNotificationStatus         `json:"fraudNotificationStatus,omitempty"`
	AccountID               string                           `json:"accountId,omitempty"`
	TransactionID           string                           `json:"transactionId,omitempty"`
	Transaction             *TransactionFraud                `json:"transaction,omitempty"`
	AuditHistory            []SourceAudit                    `json:"audit_history,omitempty"`
}

// FraudNotificationRequest from OpenAPI spec.
type FraudNotificationRequest struct {
	IssuerFraudID  string                           `json:"issuerFraudId,omitempty"`
	AccountID      string                           `json:"accountId"`
	TransactionID  string                           `json:"transactionId"`
	FraudType      string                           `json:"fraudType"`
	CardholderData *FraudNotificationCardholderData `json:"cardholderData,omitempty"`
	SourceAudit    *SourceAudit                     `json:"sourceAudit,omitempty"`
}

// FraudNotificationCreatedSuccessfully from OpenAPI spec.
type FraudNotificationCreatedSuccessfully struct {
	ResultData any `json:"resultData,omitempty"`
}

// FraudListResult from OpenAPI spec.
type FraudListResult struct {
	HasMore            bool                `json:"hasMore,omitempty"`
	FraudNotifications []FraudNotification `json:"fraudNotifications,omitempty"`
}

// FraudNotificationCardholderData from OpenAPI spec.
type FraudNotificationCardholderData struct {
	Name    string `json:"name,omitempty"`
	Zipcode string `json:"zipcode"`
	City    string `json:"city"`
	State   string `json:"state"`
}

// FraudNotificationUndoRequest from OpenAPI spec.
type FraudNotificationUndoRequest struct {
	AccountID     string       `json:"accountId"`
	TransactionID string       `json:"transactionId"`
	SourceAudit   *SourceAudit `json:"sourceAudit,omitempty"`
}

// FraudNotificationUndoCreatedSuccessfully from OpenAPI spec.
type FraudNotificationUndoCreatedSuccessfully struct {
	ResultData any `json:"resultData,omitempty"`
}

// TransactionFraud from OpenAPI spec.
type TransactionFraud struct {
	AcquirerTransactionID            string                            `json:"acquirerTransactionId,omitempty"`
	TransactionAuthorizationResponse *TransactionAuthorizationResponse `json:"transactionAuthorizationResponse,omitempty"`
	PsProductName                    string                            `json:"psProductName,omitempty"`
	MerchantName                     string                            `json:"merchantName,omitempty"`
	DisputeID                        string                            `json:"disputeId,omitempty"`
	TransactionType                  TransactionType                   `json:"transactionType"`
	TransactionSource                string                            `json:"transactionSource,omitempty"`
	Fees                             []Fee                             `json:"fees,omitempty"`
	TerminalID                       string                            `json:"terminalId,omitempty"`
	TransactionDate                  string                            `json:"transactionDate,omitempty"`
	TransactionTime                  string                            `json:"transactionTime,omitempty"`
	Incremental                      bool                              `json:"incremental,omitempty"`
	International                    bool                              `json:"international,omitempty"`
	AccountID                        string                            `json:"accountId"`
	CountryCode                      string                            `json:"countryCode,omitempty"`
	EntryMode                        *EntryMode                        `json:"entryMode,omitempty"`
	AcquirerID                       string                            `json:"acquirerId,omitempty"`
	MerchantZipcode                  string                            `json:"merchantZipcode,omitempty"`
	TransactionStatus                TransactionStatus                 `json:"transactionStatus"`
	MerchantAddress                  string                            `json:"merchantAddress,omitempty"`
	Card                             any                       `json:"card,omitempty"`
	TransactionID                    string                            `json:"transactionId"`
	MCC                              string                            `json:"mcc,omitempty"`
	MerchantCity                     string                            `json:"merchantCity,omitempty"`
	FraudNotificationStatus          *FraudNotificationStatus          `json:"fraudNotificationStatus,omitempty"`
	PsProductCode                    string                            `json:"psProductCode,omitempty"`
	Amount                           AmountFraud                       `json:"amount"`
	RemainingAmount                  *Amount                           `json:"remainingAmount,omitempty"`
	RefundedAmount                   *Amount                           `json:"refundedAmount,omitempty"`
	MerchantID                       string                            `json:"merchantId,omitempty"`
	SettlementDateTime               string                            `json:"settlementDateTime,omitempty"`
	MerchantDocumentID               string                            `json:"merchantDocumentId,omitempty"`
	MerchantUf                       string                            `json:"merchantUf,omitempty"`
	CancellingTransactionID          string                            `json:"cancellingTransactionId,omitempty"`
	CancellingTransactionIDs         []string                          `json:"cancellingTransactionIds,omitempty"`
	InternationalTransactionData     any                       `json:"internationalTransactionData,omitempty"`
	AuthorizationAdvice              bool                              `json:"authorizationAdvice,omitempty"`
	MITAdditionalData                *MITAdditionalData                `json:"mitAdditionalData,omitempty"`
	AdditionalTerminalData           *AdditionalTerminalData           `json:"additionalTerminalData,omitempty"`
}
