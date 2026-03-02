// Package types provides SDK-specific types that complement the OpenAPI-generated types.
// These types are used for SDK convenience and for authorization/webhook handlers
// which are not part of the Processadora API spec.
package types

// ============================================================================
// SDK CONVENIENCE TYPES (aliases and wrappers for backward compatibility)
// ============================================================================

// CreateAccountRequest is an alias for NewAccountRequest for SDK convenience.
type CreateAccountRequest = NewAccountRequest

// CreateAccountResponse is an alias for AccountCreatedSuccessfully.
type CreateAccountResponse = AccountCreatedSuccessfully

// ListAccountsResponse is an alias for AccountListResult.
type ListAccountsResponse = AccountListResult

// CardDetails is an alias for Card for backward compatibility.
type CardDetails = Card

// ListCardsResponse is an alias for CardListResult.
type ListCardsResponse = CardListResult

// CreateCreditRequest is an alias for NewCreditRequest.
type CreateCreditRequest = NewCreditRequest

// ScheduledTransaction represents a scheduled transaction entry.
type ScheduledTransaction struct {
	ScheduledTransactionID string `json:"scheduledTransactionId,omitempty"`
	AccountID              string `json:"accountId,omitempty"`
	Amount                 Amount `json:"amount,omitempty"`
	Type                   string `json:"type,omitempty"`
	Status                 string `json:"status,omitempty"`
	ScheduledDate          string `json:"scheduledDate,omitempty"`
	Description            string `json:"description,omitempty"`
	CreatedAt              string `json:"createdAt,omitempty"`
}

// ============================================================================
// LIMITS API TYPES (Pós-Pago)
// ============================================================================

// MaxCreditLimitsRequest represents a request to update max credit limits.
type MaxCreditLimitsRequest struct {
	TotalLimit      *Amount `json:"totalLimit,omitempty"`
	WithdrawalLimit *Amount `json:"withdrawalLimit,omitempty"`
}

// MaxCreditLimitsResponse represents the response from updating max credit limits.
type MaxCreditLimitsResponse struct {
	ResultCode        int    `json:"resultCode,omitempty"`
	ResultDescription string `json:"resultDescription,omitempty"`
	IssuerRequestID   string `json:"issuerRequestId,omitempty"`
	PSResponseID      string `json:"psResponseId,omitempty"`
}

// ChangeUsableCreditLimitsRequest represents a request to change usable credit limits.
type ChangeUsableCreditLimitsRequest struct {
	NewUsableTotalLimit int64 `json:"newUsableTotalLimit"`
}

// ChangeUsableCreditLimitsResponse represents the response from changing usable credit limits.
type ChangeUsableCreditLimitsResponse struct {
	ResultCode        int    `json:"resultCode,omitempty"`
	ResultDescription string `json:"resultDescription,omitempty"`
	IssuerRequestID   string `json:"issuerRequestId,omitempty"`
	PSResponseID      string `json:"psResponseId,omitempty"`
}

// ============================================================================
// ACCOUNT HELPER TYPES
// ============================================================================

// Holder represents an account holder (legacy SDK type).
type Holder struct {
	Name         string `json:"name"`
	Document     string `json:"document"`
	DocumentType string `json:"documentType"` // CPF, CNPJ
	BirthDate    string `json:"birthDate,omitempty"`
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
}

// AccountLimits represents account limits configuration.
type AccountLimits struct {
	CreditLimit      *LimitInfo `json:"creditLimit,omitempty"`
	WithdrawalLimit  *LimitInfo `json:"withdrawalLimit,omitempty"`
	InstallmentLimit *LimitInfo `json:"installmentLimit,omitempty"`
}

// LimitInfo represents limit information.
type LimitInfo struct {
	Total     *Amount `json:"total,omitempty"`
	Available *Amount `json:"available,omitempty"`
	Used      *Amount `json:"used,omitempty"`
	Daily     *Amount `json:"daily,omitempty"`
}

// UpdateLimitsRequest represents a request to update account limits.
type UpdateLimitsRequest struct {
	CreditLimit *Amount `json:"creditLimit,omitempty"`
}

// ListAccountsParams represents query parameters for listing accounts.
type ListAccountsParams struct {
	Limit                  int    `json:"limit,omitempty"`
	StartingAfter          string `json:"starting_after,omitempty"`
	EndingBefore           string `json:"ending_before,omitempty"`
	IdentityDocumentNumber string `json:"identity_document_number,omitempty"`
	FullName               string `json:"full_name,omitempty"`
	PSProductCode          int    `json:"ps_product_code,omitempty"`
	AccountStatus          string `json:"account_status,omitempty"`
	IssuerAccountID        string `json:"issuer_id,omitempty"`
	IncludedSince          string `json:"included_since,omitempty"`
	Sort                   string `json:"sort,omitempty"`
	First                  bool   `json:"first,omitempty"`
}

// CreditVerificationStatus represents credit verification status.
type CreditVerificationStatus struct {
	Status      string `json:"status"` // PENDING, APPROVED, REJECTED
	VerifiedAt  string `json:"verifiedAt,omitempty"`
	Reason      string `json:"reason,omitempty"`
	CreditScore int    `json:"creditScore,omitempty"`
}

// CreditTransaction represents a credit transaction.
type CreditTransaction struct {
	TransactionID string `json:"transactionId,omitempty"`
	AccountID     string `json:"accountId,omitempty"`
	Amount        Amount `json:"amount,omitempty"`
	Status        string `json:"status,omitempty"`
	Description   string `json:"description,omitempty"`
	Reference     string `json:"reference,omitempty"`
	CreatedAt     string `json:"createdAt,omitempty"`
}

// ListCreditsResponse represents the response from listing credits.
type ListCreditsResponse struct {
	Data       []CreditTransaction `json:"data,omitempty"`
	HasMore    bool                `json:"hasMore,omitempty"`
	TotalCount int                 `json:"totalCount,omitempty"`
}

// ============================================================================
// AUTHORIZATION SERVER TYPES (not in Processadora API)
// These types are used for real-time transaction authorization callbacks.
// ============================================================================

// Establishment represents a merchant establishment.
type Establishment struct {
	Name       string `json:"name,omitempty"`
	MCC        string `json:"mcc,omitempty"`
	MerchantID string `json:"merchantId,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
}

// ProcessingCode represents ISO 8583 processing code.
type ProcessingCode string

const (
	ProcessingCodePurchase   ProcessingCode = "00"
	ProcessingCodeWithdrawal ProcessingCode = "01"
	ProcessingCodeCashback   ProcessingCode = "09"
	ProcessingCodeRefund     ProcessingCode = "20"
	ProcessingCodePayment    ProcessingCode = "28"
)

// Wallet represents a digital wallet type.
type Wallet string

const (
	WalletApplePay   Wallet = "APPLE_PAY"
	WalletGooglePay  Wallet = "GOOGLE_PAY"
	WalletSamsungPay Wallet = "SAMSUNG_PAY"
)

// TokenStatus represents the status of a wallet token.
type TokenStatus string

const (
	TokenStatusInactive  TokenStatus = "INACTIVE"
	TokenStatusActive    TokenStatus = "ACTIVE"
	TokenStatusSuspended TokenStatus = "SUSPENDED"
	TokenStatusDeleted   TokenStatus = "DELETED"
)

// CustomProvisioningDataResponse represents custom data for wallet display.
type CustomProvisioningDataResponse struct {
	CardArt               string `json:"cardArt"`
	IssuerName            string `json:"issuerName"`
	ShortDescription      string `json:"shortDescription"`
	TermsAndConditionsURL string `json:"termsAndConditionsUrl"`
}

// Brand represents a card brand.
type Brand string

const (
	BrandMastercard Brand = "MASTERCARD"
	BrandVisa       Brand = "VISA"
	BrandElo        Brand = "ELO"
	BrandHipercard  Brand = "HIPERCARD"
	BrandAmex       Brand = "AMEX"
)

// Error represents an API error response.
type Error struct {
	Code    string   `json:"code,omitempty"`
	Message string   `json:"message,omitempty"`
	Details []string `json:"details,omitempty"`
}

// Pagination represents pagination information.
type Pagination struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"pageSize,omitempty"`
	TotalItems int `json:"totalItems,omitempty"`
	TotalPages int `json:"totalPages,omitempty"`
}

// ============================================================================
// CARD HELPER TYPES
// ============================================================================

// PhysicalCardRequest represents a request to create a physical card.
type PhysicalCardRequest = NewCardRequest

// PhysicalCardResponse represents the response from creating a physical card.
type PhysicalCardResponse = NewCardrequestCreatedSuccessfully

// NewCardResponse is an alias for NewCardrequestCreatedSuccessfully.
type NewCardResponse = NewCardrequestCreatedSuccessfully

// BlockCardRequest is an alias for CardBlockRequest.
type BlockCardRequest = CardBlockRequest

// UnblockCardRequest is an alias for CardUnblockRequest.
type UnblockCardRequest = CardUnblockRequest

// ChangePinRequest is an alias for PinChangeRequest.
type ChangePinRequest = PinChangeRequest

// VerifyPinRequest is an alias for PinValidateRequest.
type VerifyPinRequest = PinValidateRequest

// VerifyPinResponse is an alias for PinValidateSuccessfully.
type VerifyPinResponse = PinValidateSuccessfully

// SensitiveCardData represents encrypted sensitive card data.
type SensitiveCardData struct {
	EncryptedPAN        string `json:"encryptedPan,omitempty"`
	EncryptedCVV        string `json:"encryptedCvv,omitempty"`
	EncryptedExpiration string `json:"encryptedExpiration,omitempty"`
}

// CreateDebitRequest is an alias for NewDebitRequest.
type CreateDebitRequest = NewDebitRequest

// CardLimits represents card-level limits.
type CardLimits struct {
	Limits []CardLimit `json:"limits,omitempty"`
}

// CardLimit represents a single limit configuration.
type CardLimit struct {
	Type             string  `json:"type,omitempty"`
	MaxAmount        *Amount `json:"maxAmount,omitempty"`
	UsedAmount       *Amount `json:"usedAmount,omitempty"`
	MaxTransactions  int     `json:"maxTransactions,omitempty"`
	UsedTransactions int     `json:"usedTransactions,omitempty"`
}

// UpdateCardLimitsRequest represents a request to update card limits.
type UpdateCardLimitsRequest struct {
	Limits []CardLimit `json:"limits,omitempty"`
}

// MCCRestrictions represents MCC restrictions for a card.
type MCCRestrictions struct {
	Mode        string   `json:"mode,omitempty"` // BLOCKLIST, ALLOWLIST
	BlockedMCCs []string `json:"blockedMccs,omitempty"`
	AllowedMCCs []string `json:"allowedMccs,omitempty"`
}

// GeoRestrictions represents geographic restrictions for a card.
type GeoRestrictions struct {
	InternationalAllowed bool     `json:"internationalAllowed,omitempty"`
	AllowedCountries     []string `json:"allowedCountries,omitempty"`
	BlockedCountries     []string `json:"blockedCountries,omitempty"`
}

// ForceAcceptConfig represents force accept configuration.
type ForceAcceptConfig struct {
	Scenarios       []string `json:"scenarios,omitempty"`
	Duration        string   `json:"duration,omitempty"` // e.g., "1H", "24H"
	MaxTransactions int      `json:"maxTransactions,omitempty"`
	MaxAmount       *Amount  `json:"maxAmount,omitempty"`
}

// ForceAcceptResponse represents force accept status.
type ForceAcceptResponse struct {
	ForceAcceptID         string   `json:"forceAcceptId,omitempty"`
	Enabled               bool     `json:"enabled,omitempty"`
	Scenarios             []string `json:"scenarios,omitempty"`
	ExpiresAt             string   `json:"expiresAt,omitempty"`
	RemainingTransactions int      `json:"remainingTransactions,omitempty"`
	UsedTransactions      int      `json:"usedTransactions,omitempty"`
}

// ListCardsParams represents query parameters for listing cards.
// Uses snake_case parameter names as per API spec.
type ListCardsParams struct {
	Limit                  int    `json:"limit,omitempty"`                    // Maximum number of cards to return
	StartingAfter          string `json:"starting_after,omitempty"`           // Cursor for pagination (card ID)
	EndingBefore           string `json:"ending_before,omitempty"`            // Cursor for pagination (card ID)
	IssuerCardID           string `json:"issuer_id,omitempty"`                // Filter: issuer's card identifier
	IdentityDocumentNumber string `json:"identity_document_number,omitempty"` // Filter: cardholder document number
	PANLast4Digits         string `json:"pan_last_4_digits,omitempty"`        // Filter: last 4 digits of PAN
	IssuedOnOrAfterDate    string `json:"issued_on_or_after_date,omitempty"`  // Filter: issuance date (dd/MM/yyyy)
	PSProductCode          int    `json:"ps_product_code,omitempty"`          // Filter: product code
	LinkID                 string `json:"link_id,omitempty"`                  // Filter: link ID between cards
	AlternativeBindingKey  string `json:"alternative_binding_key,omitempty"`  // Filter: alternative binding key
}

// AnonymousCardRequest is an alias for NewAnonymousCardRequest.
type AnonymousCardRequest = NewAnonymousCardRequest

// AnonymousCardResponse is an alias for NewAnonymousCardRequestCreatedSuccessfully.
type AnonymousCardResponse = NewAnonymousCardRequestCreatedSuccessfully

// BindAnonymousCardRequest is an alias for AssociateAnonymousCardRequest.
type BindAnonymousCardRequest = AssociateAnonymousCardRequest

// BlockAndReissueCardRequest is an alias for CardBlockAndReissueRequest.
type BlockAndReissueCardRequest = CardBlockAndReissueRequest

// ReissueCardRequest is an alias for CardReissueRequest.
type ReissueCardRequest = CardReissueRequest

// CardFunctions represents card functions configuration.
type CardFunctions struct {
	ContactlessEnabled   bool `json:"contactlessEnabled,omitempty"`
	EcommerceEnabled     bool `json:"ecommerceEnabled,omitempty"`
	ATMEnabled           bool `json:"atmEnabled,omitempty"`
	InternationalEnabled bool `json:"internationalEnabled,omitempty"`
}

// CardDataprepStatus represents the dataprep status of a card.
type CardDataprepStatus struct {
	Status       string `json:"status,omitempty"` // PENDING, COMPLETED, FAILED
	LastUpdated  string `json:"lastUpdated,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// ============================================================================
// DISPUTE HELPER TYPES
// ============================================================================

// ListDisputesRequest represents query parameters for listing disputes.
// Uses snake_case parameter names as per API spec.
type ListDisputesRequest struct {
	Limit         int    `json:"limit,omitempty"`          // Maximum number of disputes to return
	StartingAfter string `json:"starting_after,omitempty"` // Cursor for pagination (dispute ID)
	EndingBefore  string `json:"ending_before,omitempty"`  // Cursor for pagination (dispute ID)
	DisputeCode   string `json:"dispute_reason,omitempty"` // Filter: dispute reason code
	DisputeStatus string `json:"dispute_status,omitempty"` // Filter: dispute status
	BeginningDate string `json:"beginning_date,omitempty"` // Filter: start date (YYYY-MM-DD)
	EndingDate    string `json:"ending_date,omitempty"`    // Filter: end date (YYYY-MM-DD)
}

// ListDisputesResponse is an alias for DisputeListResult.
type ListDisputesResponse = DisputeListResult

// CreateDisputeRequest is an alias for DisputeRequest.
type CreateDisputeRequest = DisputeRequest

// CreateDisputeResponse is an alias for DisputeCreatedSuccessfully.
type CreateDisputeResponse = DisputeCreatedSuccessfully

// CancelDisputeRequest represents a request to cancel a dispute.
type CancelDisputeRequest struct {
	Reason string `json:"reason,omitempty"`
}

// UpdateDisputeRequest represents a request to update a dispute.
type UpdateDisputeRequest struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// ============================================================================
// HCE HELPER TYPES
// ============================================================================

// CreateHCECardSuccessfully is an alias for CreateHCECardSuccesfully (typo in API).
type CreateHCECardSuccessfully = CreateHCECardSuccesfully

// UnprovisionSuccessfully is an alias for UnprovisionSuccesfully (typo in API).
type UnprovisionSuccessfully = UnprovisionSuccesfully

// ============================================================================
// DYNAMIC PROVISIONING HELPER TYPES
// ============================================================================

// DynamicProvisioningBalance is an alias for Balance.
type DynamicProvisioningBalance = Balance

// DynamicProvisioningCreditList is an alias for CreditListMattress.
type DynamicProvisioningCreditList = CreditListMattress

// DynamicProvisioningCredit is an alias for CreditMattress.
type DynamicProvisioningCredit = CreditMattress

// ============================================================================
// INSTALLMENT HELPER TYPES
// ============================================================================

// InstallmentSimulationBody is an alias for installmentSimulationBody.
type InstallmentSimulationBody = installmentSimulationBody

// InstallmentSimulationResult is an alias for installmentSimulationResult.
type InstallmentSimulationResult = installmentSimulationResult

// CreateInstallmentResponse represents the response from creating an installment.
type CreateInstallmentResponse struct {
	Success       bool   `json:"success,omitempty"`
	TransactionID string `json:"transactionId,omitempty"`
	Message       string `json:"message,omitempty"`
}

// ============================================================================
// STATEMENT HELPER TYPES
// ============================================================================

// ClosedStatementsRequest represents query parameters for getting closed statements.
type ClosedStatementsRequest struct {
	ClosingDateQuery string `json:"closing_date_query,omitempty"` // Format: dd/MM/yyyy
	StartingAfter    string `json:"starting_after,omitempty"`     // Cursor for pagination
}

// ClosedStatementsResponse represents the response from getting closed statements.
type ClosedStatementsResponse = StatementList

// ============================================================================
// REPORT HELPER TYPES
// ============================================================================

// ReportListResult represents the result of listing reports.
type ReportListResult struct {
	Reports []string `json:"reports,omitempty"`
	HasMore bool     `json:"hasMore,omitempty"`
}

// ReportTemporaryUrl represents a temporary URL for a report.
type ReportTemporaryUrl struct {
	URL       string `json:"url,omitempty"`
	ExpiresAt string `json:"expiresAt,omitempty"`
}

// ============================================================================
// QRCODE/PAYMENT HELPER TYPES
// ============================================================================

// SimplePaymentResponse represents the response from a simple payment.
type SimplePaymentResponse struct {
	TransactionID string `json:"transactionId,omitempty"`
	Status        string `json:"status,omitempty"`
	Message       string `json:"message,omitempty"`
}

// TransactionCallbackRequest represents a transaction callback request.
type TransactionCallbackRequest struct {
	TransactionID string `json:"transactionId,omitempty"`
	Status        string `json:"status,omitempty"`
	ResponseCode  string `json:"responseCode,omitempty"`
	Message       string `json:"message,omitempty"`
}

// TransactionCallbackResponse represents the response from a transaction callback.
type TransactionCallbackResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

// ============================================================================
// TRANSACTION HELPER TYPES
// ============================================================================

// ListTransactionsRequest represents query parameters for listing transactions per account.
// Uses snake_case parameter names as per API spec.
type ListTransactionsRequest struct {
	Limit                 int    `json:"limit,omitempty"`                   // Maximum number of transactions to return
	StartingAfter         string `json:"starting_after,omitempty"`          // Cursor for pagination (transaction ID)
	EndingBefore          string `json:"ending_before,omitempty"`           // Cursor for pagination (transaction ID)
	BeginningDate         string `json:"beginning_date,omitempty"`          // Filter: start date (YYYY-MM-DD)
	EndingDate            string `json:"ending_date,omitempty"`             // Filter: end date (YYYY-MM-DD)
	TransactionType       string `json:"transaction_type,omitempty"`        // Filter: transaction type
	TransactionStatus     string `json:"transaction_status,omitempty"`      // Filter: transaction status
	TransactionApproved   *bool  `json:"transaction_approved,omitempty"`    // Filter: approved transactions
	TransactionDenialCode string `json:"transaction_denial_code,omitempty"` // Filter: denial code
	MinimumAmount         int64  `json:"minimum_amount,omitempty"`          // Filter: minimum amount
	MaxAmount             int64  `json:"max_amount,omitempty"`              // Filter: maximum amount
	TransactionEntryMode  string `json:"transaction_mode,omitempty"`        // Filter: entry mode
	Ordinated             string `json:"ordinated,omitempty"`               // Sort order
}

// ListAllTransactionsRequest represents query parameters for GET /transactions endpoint.
type ListAllTransactionsRequest struct {
	Limit                 int    `json:"limit,omitempty"`                   // Maximum number of transactions to return
	StartingAfter         string `json:"starting_after,omitempty"`          // Cursor for pagination (transaction ID)
	EndingBefore          string `json:"ending_before,omitempty"`           // Cursor for pagination (transaction ID)
	BeginningDate         string `json:"beginning_date,omitempty"`          // Filter: start date (YYYY-MM-DD)
	EndingDate            string `json:"ending_date,omitempty"`             // Filter: end date (YYYY-MM-DD)
	TransactionType       string `json:"transaction_type,omitempty"`        // Filter: transaction type
	TransactionStatus     string `json:"transaction_status,omitempty"`      // Filter: transaction status
	TransactionApproved   *bool  `json:"transaction_approved,omitempty"`    // Filter: approved transactions
	TransactionDenialCode string `json:"transaction_denial_code,omitempty"` // Filter: denial code
	MinimumAmount         int64  `json:"minimum_amount,omitempty"`          // Filter: minimum amount
	MaxAmount             int64  `json:"max_amount,omitempty"`              // Filter: maximum amount
	TransactionEntryMode  string `json:"transaction_mode,omitempty"`        // Filter: entry mode
}

// ListTransactionsResponse is an alias for TransactionListResult.
type ListTransactionsResponse = TransactionListResult

// CreditAdjustmentRequest represents a request to create a credit adjustment.
type CreditAdjustmentRequest struct {
	IssuerRequestID string       `json:"issuerRequestId,omitempty"`
	Description     string       `json:"description,omitempty"`
	Amount          Amount       `json:"amount"`
	Type            string       `json:"type,omitempty"` // adjustment
	SourceAudit     *SourceAudit `json:"sourceAudit,omitempty"`
}

// AdjustmentResponse represents the response from an adjustment operation.
type AdjustmentResponse struct {
	ResultData    any `json:"resultData,omitempty"`
	TransactionID string      `json:"transactionId,omitempty"`
}

// ============================================================================
// STATUS HELPER TYPES
// ============================================================================

// GetStatusResponse is an alias for ResultData.
type GetStatusResponse = ResultData

// GetHealthStatusResponse is an alias for HealthCheckResponse.
type GetHealthStatusResponse = HealthCheckResponse

// ListDataprepStatusResponse is an alias for DataprepStatusListResult.
type ListDataprepStatusResponse = DataprepStatusListResult

// GetDataprepStatusResponse is an alias for DataprepStatusResult.
type GetDataprepStatusResponse = DataprepStatusResult

// RegisterPaymentResponse represents the response from registering a payment.
type RegisterPaymentResponse struct {
	ResultData    any `json:"resultData,omitempty"`
	TransactionID string      `json:"transactionId,omitempty"`
	Status        string      `json:"status,omitempty"`
}

// DebitAdjustmentRequest represents a request to create a debit adjustment.
type DebitAdjustmentRequest struct {
	IssuerRequestID string       `json:"issuerRequestId,omitempty"`
	Description     string       `json:"description,omitempty"`
	Amount          Amount       `json:"amount"`
	Type            string       `json:"type,omitempty"` // adjustment
	SourceAudit     *SourceAudit `json:"sourceAudit,omitempty"`
}

// ============================================================================
// VIRTUAL CARD HELPER TYPES
// ============================================================================

// ModifyVirtualCardCVVResponse is an alias for ModifyVirtualCardResponse.
type ModifyVirtualCardCVVResponse = ModifyVirtualCardResponse

// ============================================================================
// XPAYS / DEVICE TOKEN TYPES (Based on Official PaySmart xPays API OpenAPI Spec)
// Reference: https://paysmart-api.gitlab.io/processadora/PT-br/docs/xpays-api
// ============================================================================

// IDVMethod represents the identity verification method for token activation.
// Reference: xpays-api.json - idvMethod enum
type IDVMethod string

const (
	IDVMethodApp      IDVMethod = "App"
	IDVMethodOTPEmail IDVMethod = "OTPEmail"
	IDVMethodOTPSMS   IDVMethod = "OTPSMS"
)

// TerminateReason represents the reason for terminating a device token.
// Reference: xpays-api.json - terminateReason enum
type TerminateReason string

const (
	TerminateReasonDeleted TerminateReason = "DELETED"
	TerminateReasonRevoked TerminateReason = "REVOKED"
	TerminateReasonExpired TerminateReason = "EXPIRED"
)

// DeviceToken represents a device token for mobile payments (xPays).
// Reference: xpays-api.json - DeviceToken schema
// All fields match the official OpenAPI specification.
type DeviceToken struct {
	// Core identifiers
	DeviceTokenID     string `json:"deviceTokenId,omitempty"`
	CardReferenceID   string `json:"cardReferenceId,omitempty"`
	WalletDeviceToken string `json:"walletDeviceTokenId,omitempty"`

	// Wallet info
	Wallet     string `json:"wallet,omitempty"` // google, apple
	DPANLast4  string `json:"dpanLast4Digits,omitempty"`
	DeviceType string `json:"deviceType,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`

	// Deployment status
	DeploymentStatus string `json:"deploymentStatus,omitempty"`
	DeploymentTs     string `json:"deploymentTs,omitempty"`

	// Suspension status
	SuspensionStatus string `json:"suspensionStatus,omitempty"`
	SuspensionTs     string `json:"suspensionTs,omitempty"`

	// Activation status
	ActivationStatus string `json:"activationStatus,omitempty"`
	ActivationTs     string `json:"activationTs,omitempty"`

	// Identity verification
	IDVMethod string `json:"idvMethod,omitempty"`
	IDVStatus string `json:"idvStatus,omitempty"`
	IDVTs     string `json:"idvTs,omitempty"`

	// Timestamps
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// DeviceTokensResponse represents the response from listing device tokens.
// GET /cards/{cardId}/deviceTokens
type DeviceTokensResponse struct {
	DeviceTokens []DeviceToken `json:"deviceTokens,omitempty"`
}

// MigrateTokensRequest represents a request to migrate device tokens.
// POST /cards/{cardId}/deviceTokens/migrate
// Reference: xpays-api.json - Simple cardId field as per spec
type MigrateTokensRequest struct {
	CardID string `json:"cardId"`
}

// MigrateTokensResponse represents the response from migrating device tokens.
type MigrateTokensResponse struct {
	ResultData *ResultData `json:"resultData,omitempty"`
}

// SuspendResumeRequest represents a request to suspend or resume device tokens.
// POST /cards/{cardId}/deviceTokens/suspendOrResume
// POST /deviceTokens/{deviceTokenId}/suspendOrResume
// Reference: xpays-api.json - Uses boolean suspend, not action enum
type SuspendResumeRequest struct {
	Suspend           bool   `json:"suspend"`
	ReasonDescription string `json:"reasonDescription,omitempty"`
}

// SuspendResumeResponse represents the response from suspending/resuming device tokens.
type SuspendResumeResponse struct {
	DeviceTokens []DeviceToken `json:"deviceTokens,omitempty"`
}

// TerminateTokensRequest represents a request to terminate device tokens for a card.
// POST /cards/{cardId}/deviceTokens/terminate
// Reference: xpays-api.json - Uses enum reason
type TerminateTokensRequest struct {
	Reason            TerminateReason `json:"reason"`
	ReasonDescription string          `json:"reasonDescription,omitempty"`
}

// TerminateTokensResponse represents the response from terminating device tokens.
type TerminateTokensResponse struct {
	DeviceTokens []DeviceToken `json:"deviceTokens,omitempty"`
}

// TerminateTokenRequest represents a request to terminate a single device token.
// POST /deviceTokens/{deviceTokenId}/terminate
type TerminateTokenRequest struct {
	Reason            TerminateReason `json:"reason"`
	ReasonDescription string          `json:"reasonDescription,omitempty"`
}

// TerminateTokenResponse represents the response from terminating a device token.
type TerminateTokenResponse struct {
	DeviceToken *DeviceToken `json:"deviceToken,omitempty"`
}

// ActivateTokenRequest represents a request to activate a device token.
// POST /deviceTokens/{deviceTokenId}/activate
// Reference: xpays-api.json - Uses idvMethod enum
type ActivateTokenRequest struct {
	IDVMethod         IDVMethod `json:"idvMethod"`
	ReasonDescription string    `json:"reasonDescription,omitempty"`
}

// ActivateTokenResponse represents the response from activating a device token.
type ActivateTokenResponse struct {
	DeviceToken *DeviceToken `json:"deviceToken,omitempty"`
}

// EncryptedCardRequest represents a request to get encrypted card data.
// POST /cards/{cardId}/encryptedCard
// Reference: Official PaySmart xPays API - Uses base64-encoded strings
type EncryptedCardRequest struct {
	Nonce          string   `json:"nonce"`
	NonceSignature string   `json:"nonceSignature"`
	Certificates   []string `json:"certificates"`
}

// EncryptedCardResponse represents the response with encrypted card data.
// Reference: Official PaySmart xPays API - Returns base64-encoded strings
type EncryptedCardResponse struct {
	EphemeralPublicKey string `json:"ephemeralPublicKey,omitempty"`
	EncryptedPassData  string `json:"encryptedPassData,omitempty"`
	ActivationData     string `json:"activationData,omitempty"`
}

// OpaqueCardResponse represents the response with opaque card data.
// GET /cards/{cardId}/opaqueCard
// Reference: Official PaySmart xPays API - sender string + cardDescriptor base64 string
type OpaqueCardResponse struct {
	Sender         string `json:"sender,omitempty"`
	CardDescriptor string `json:"cardDescriptor,omitempty"`
}
