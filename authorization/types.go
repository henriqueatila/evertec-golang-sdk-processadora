// Package authorization provides types and handlers for the Issuer Authorization API.
// This is the API that the emissor (you) must implement. Evertec calls your endpoints.
// Types are based on: https://paysmart-api.gitlab.io/processadora/PT-br/docs/compra
package authorization

import "github.com/henriqueatila/evertec-golang-sdk-processadora/types"

// ========================================
// Amount Types
// ========================================

// AuthAmount represents an amount in the authorization API (uses currency_code, not currencyCode).
type AuthAmount struct {
	Amount       int64 `json:"amount"`
	CurrencyCode int   `json:"currency_code"`
}

// ========================================
// Common Types
// ========================================

// AuthorizationInfo contains authorization code and description from Evertec.
type AuthorizationInfo struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// AuthCard represents card information in authorization requests.
type AuthCard struct {
	PaysmartID string `json:"paysmart_id"`
	IssuerID   int    `json:"issuer_id,omitempty"`
	PAN        string `json:"pan"`
	PANSeq     string `json:"panseq,omitempty"`
}

// ProcessingCode represents the transaction processing code.
type ProcessingCode struct {
	TipoTransacao          string `json:"tipo_transacao"`
	SourceAccountType      string `json:"source_account_type"`
	DestinationAccountType string `json:"destination_account_type"`
}

// Fee represents a fee associated with a transaction.
type Fee struct {
	Amount *AuthAmount `json:"amount"`
	Type   string      `json:"type"` // iof, markup, boarding_fee, withdrawal_fee, others
}

// Establishment represents merchant/establishment information.
type Establishment struct {
	MCC     string `json:"mcc"`
	Name    string `json:"name"`
	City    string `json:"city,omitempty"`
	Address string `json:"address,omitempty"`
	Zipcode string `json:"zipcode,omitempty"`
	Country string `json:"country,omitempty"`
	CNPJ    string `json:"cnpj,omitempty"`
}

// CancellationReason represents the reason for cancellation.
type CancellationReason struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// OriginalISO8583 represents the original ISO-8583 message data elements.
type OriginalISO8583 struct {
	MTI   string `json:"mti,omitempty"`
	DE002 string `json:"de002,omitempty"` // Primary Account Number
	DE003 string `json:"de003,omitempty"` // Processing Code
	DE004 string `json:"de004,omitempty"` // Transaction Amount
	DE007 string `json:"de007,omitempty"` // Transmission Date/Time
	DE011 string `json:"de011,omitempty"` // System Trace Audit Number
	DE012 string `json:"de012,omitempty"` // Local Transaction Time
	DE013 string `json:"de013,omitempty"` // Local Transaction Date
	DE014 string `json:"de014,omitempty"` // Expiration Date
	DE018 string `json:"de018,omitempty"` // Merchant Category Code
	DE019 string `json:"de019,omitempty"` // Acquiring Country Code
	DE022 string `json:"de022,omitempty"` // POS Entry Mode
	DE024 string `json:"de024,omitempty"` // Network International ID
	DE025 string `json:"de025,omitempty"` // POS Condition Code
	DE032 string `json:"de032,omitempty"` // Acquiring Institution ID
	DE037 string `json:"de037,omitempty"` // Retrieval Reference Number
	DE038 string `json:"de038,omitempty"` // Authorization ID Response
	DE039 string `json:"de039,omitempty"` // Response Code
	DE041 string `json:"de041,omitempty"` // Terminal ID
	DE042 string `json:"de042,omitempty"` // Merchant ID
	DE043 string `json:"de043,omitempty"` // Merchant Name/Location
	DE046 string `json:"de046,omitempty"` // Additional Data
	DE048 string `json:"de048,omitempty"` // Additional Data
	DE049 string `json:"de049,omitempty"` // Currency Code
	DE058 string `json:"de058,omitempty"` // National POS Geographic Data
	DE060 string `json:"de060,omitempty"` // Advice Reason Code
	DE062 string `json:"de062,omitempty"` // Custom Payment Service Fields
	DE127 string `json:"de127,omitempty"` // Network Data
}

// FraudData contains fraud scoring and recommendation data.
type FraudData struct {
	CreditorFraudScore          int    `json:"creditorFraudScore,omitempty"`
	EloBrandFraudScore          int    `json:"eloBrandFraudScore,omitempty"`
	FraudScorePrimaryReason     int    `json:"fraudScorePrimaryReason,omitempty"`
	FraudScoreSecondaryReason   int    `json:"fraudScoreSecondaryReason,omitempty"`
	FraudScoreTertiaryReason    int    `json:"fraudScoreTertiaryReason,omitempty"`
	FraudDecisionRecommendation string `json:"fraudDecisionRecommendation,omitempty"`
	ScoreOriginIndicator        int    `json:"scoreOriginIndicator,omitempty"`
}

// AdditionalTerminalData contains extended terminal/POS information.
type AdditionalTerminalData struct {
	TerminalType                 string `json:"terminalType,omitempty"`
	PartialApprovalIndicator     string `json:"partialApprovalIndicator,omitempty"` // notSupported, supported, onlyPurchasesSupported, etc.
	TerminalLocationIndicator    string `json:"terminalLocationIndicator,omitempty"`
	CardholderPresenceIndicator  string `json:"cardholderPresenceIndicator,omitempty"`
	CardPresenceIndicator        string `json:"cardPresenceIndicator,omitempty"`
	CardCaptureCapabilityInd     string `json:"cardCaptureCapabilityIndicator,omitempty"`
	TransactionStatusIndicator   string `json:"transactionStatusIndicator,omitempty"`
	TransactionSecurityIndicator string `json:"transactionSecurityIndicator,omitempty"`
	TerminalPOSType              string `json:"terminalPOSType,omitempty"`
	TerminalInputCapability      string `json:"terminalInputCapability,omitempty"`
}

// AuthInstallmentDetails contains installment/parcelamento information.
type AuthInstallmentDetails struct {
	FinType                     string      `json:"finType,omitempty"` // interestFree
	FareAmount                  *AuthAmount `json:"fare_amount,omitempty"`
	InsuranceAmount             *AuthAmount `json:"insurance_amount,omitempty"`
	ThirdPartiesPaymentAmount   *AuthAmount `json:"third_parties_paymnt_amount,omitempty"`
	RecordsPaymentsAmount       *AuthAmount `json:"records_payments_amount,omitempty"`
	IssuerTotalCalculatedAmount *AuthAmount `json:"issuer_total_calculated_amount,omitempty"`
	FirstPaymentDate            string      `json:"first_paymnt_date,omitempty"`
	InstallmentNumber           int         `json:"instalmnt_nbr,omitempty"`
	MonthlyInterestRate         int         `json:"monthly_interest_rate,omitempty"`
	TotalEffectiveCostRateCET   int         `json:"total_effective_cost_rate_cet,omitempty"`
	InstallmentAmount           *AuthAmount `json:"instalmnt_amount,omitempty"`
	AnnualInterestRate          int         `json:"annual_interest_rate,omitempty"`
	InputValue                  *AuthAmount `json:"input_value,omitempty"`
}

// AuthTokenPaymentData contains token payment information (Apple Pay, Google Pay, etc.).
type AuthTokenPaymentData struct {
	FPAN                                       string `json:"fPan,omitempty"`
	RequesterIDToken                           string `json:"requesterIdToken,omitempty"`
	PAN                                        string `json:"pan,omitempty"`
	TokenPSN                                   string `json:"tokenPSN,omitempty"`
	TokenExpirationDate                        string `json:"tokenExpirationDate,omitempty"`
	TokenStatus                                string `json:"tokenStatus,omitempty"`
	TokenCryptogramVerificationResults         string `json:"tokenCryptogramVerificationResults,omitempty"`
	EMVTokenCryptogramVerificationResults      string `json:"EMVTokenCryptogramVerificationResults,omitempty"`
	TokenConstraintsVerificationStatus         string `json:"tokenConstraintsVerificationStatus,omitempty"`
	TransactionDateTimeConstraint              string `json:"transactionDateTimeConstraint,omitempty"`
	TransactionAmountConstraint                string `json:"transactionAmountConstraint,omitempty"`
	UsageConstraint                            string `json:"usageConstraint,omitempty"`
	TokenATCVerificationResults                string `json:"tokenATCVerificationResults,omitempty"`
	CVE2TokenCryptogramVerificationStatus      string `json:"CVE2TokenCryptogramVerificationStatus,omitempty"`
	MerchantVerification                       string `json:"merchantVerification,omitempty"`
	MagstripeTokenCryptogramVerificationResult string `json:"magstripeTokenCryptogramVerificationResults,omitempty"`
	CVE2OutputTokenCryptogramVerificationRes   string `json:"CVE2OutputTokenCryptogramVerificationResults,omitempty"`
}

// MITAdditionalData contains Merchant Initiated Transaction data.
type MITAdditionalData struct {
	InstallmentTotalNbr  string `json:"installmentTotalNbr,omitempty"`
	PaymentType          string `json:"paymentType,omitempty"`
	TransactionType      string `json:"transactionType,omitempty"`
	TransactionAmountMIT int    `json:"transactionAmountMIT,omitempty"`
	UniqueTransactionID  string `json:"uniqueTransactionID,omitempty"`
	TransactionFrequency string `json:"transactionFrequency,omitempty"`
	ValidationIndicator  string `json:"validationIndicator,omitempty"`
	ValidationReference  string `json:"validationReference,omitempty"`
	SequenceIndicator    int    `json:"sequenceIndicator,omitempty"`
}

// PurchaseWithCashbackData contains cashback-related data.
type PurchaseWithCashbackData struct {
	PurchaseOnlyApprovalSupport bool        `json:"purchaseOnlyApprovalSupport,omitempty"`
	PurchaseOnlyAmount          *AuthAmount `json:"purchaseOnlyAmount,omitempty"`
	CashbackAmount              *AuthAmount `json:"cashbackAmount,omitempty"`
}

// ========================================
// Purchase Types
// ========================================

// PurchaseRequest represents an incoming purchase authorization request.
// POST /purchases - Called by Evertec when a purchase needs authorization.
type PurchaseRequest struct {
	// Required identifiers
	PurchaseID    string `json:"purchase_id"`
	AccountID     string `json:"account_id"`
	PSProductCode string `json:"psProductCode"`
	PSProductName string `json:"psProductName,omitempty"`

	// Transaction context
	CountryCode              string `json:"countryCode"`
	Source                   string `json:"source"`                // paySmart, issuer, brand
	ProductType              string `json:"productType,omitempty"` // pre, post
	CallingSystemName        string `json:"callingSystemName,omitempty"`
	PreAuthorization         bool   `json:"preAuthorization,omitempty"`
	IncrementalAuthorization bool   `json:"incrementalAuthorization,omitempty"`

	// Authorization info from Evertec
	Authorization *AuthorizationInfo `json:"authorization,omitempty"`

	// Card and brand info
	Brand string    `json:"brand"` // elo, goodcard, mastercard, visa
	Card  *AuthCard `json:"card"`

	// Amount information
	TotalAmount    *AuthAmount `json:"total_amount"`
	OriginalAmount *AuthAmount `json:"original_amount,omitempty"`
	DollarAmount   *AuthAmount `json:"dollar_amount,omitempty"`
	DollarRealRate string      `json:"dollar_real_rate,omitempty"`
	Spread         string      `json:"spread,omitempty"`

	// Transaction details
	EntryMode            string          `json:"entry_mode"` // chip, ecommerce, contactless, magnetic_stripe, fallback
	ProcessingCode       *ProcessingCode `json:"processing_code,omitempty"`
	HolderValidationMode string          `json:"holder_validation_mode,omitempty"` // online_pin, other

	// Fees and establishment
	Fees          []Fee          `json:"fees,omitempty"`
	Establishment *Establishment `json:"establishment,omitempty"`

	// International transaction flag
	Internacional bool `json:"internacional,omitempty"`

	// ISO-8583 original message
	OriginalISO8583 *OriginalISO8583 `json:"original_iso8583,omitempty"`

	// Optional flags
	ForceAccept         bool   `json:"forceAccept,omitempty"`
	NRID                string `json:"nrid,omitempty"`
	AuthorizationAdvice bool   `json:"authorizationAdvice,omitempty"`

	// Fraud and terminal data
	FraudData              *FraudData              `json:"fraudData,omitempty"`
	AdditionalTerminalData *AdditionalTerminalData `json:"additionalTerminalData,omitempty"`

	// Product and installment details
	EloProductCode     string                  `json:"eloProductCode,omitempty"`
	InstallmentDetails *AuthInstallmentDetails `json:"installmentDetails,omitempty"`

	// MIT and token data
	MITAdditionalData *MITAdditionalData    `json:"mitAdditionalData,omitempty"`
	TokenPaymentData  *AuthTokenPaymentData `json:"tokenPaymentData,omitempty"`

	// Cashback
	PurchaseWithCashbackData *PurchaseWithCashbackData `json:"purchaseWithCashbackData,omitempty"`

	// HCE and AFD
	HCETransaction     bool   `json:"hceTransaction,omitempty"`
	AFDAuthorization   string `json:"afdAuthorization,omitempty"` // initialAuthorization, finalAuthorization
	OriginalPurchaseID string `json:"original_purchase_id,omitempty"`
}

// PurchaseResponse represents a purchase authorization response.
// Return this from HandlePurchase to approve/deny the transaction.
type PurchaseResponse struct {
	Message         string      `json:"message"`
	Code            int         `json:"code"` // Response code (00=approved, 51=insufficient funds, etc.)
	AuthorizationID int         `json:"authorization_id,omitempty"`
	Balance         *AuthAmount `json:"balance,omitempty"`

	// Partial approval fields (for cashback scenarios)
	PurchaseOnlyApproval              bool        `json:"purchaseOnlyApproval,omitempty"`
	PurchaseOnlyPartialAmountApproved *AuthAmount `json:"purchaseOnlyPartialAmountApproved,omitempty"`
	CashbackOnlyPartialAmountApproved *AuthAmount `json:"cashbackOnlyPartialAmountApproved,omitempty"`
}

// ========================================
// Purchase Cancellation Types
// ========================================

// PurchaseCancellationRequest represents a purchase cancellation request.
// POST /purchases/cancel - Called by Evertec to cancel a purchase.
type PurchaseCancellationRequest struct {
	CancellationID          string `json:"cancellation_id"`
	OriginalPurchaseID      string `json:"original_purchase_id"`
	AccountID               string `json:"account_id"`
	PSProductCode           string `json:"psProductCode"`
	PSProductName           string `json:"psProductName,omitempty"`
	CountryCode             string `json:"countryCode"`
	Source                  string `json:"source"` // paySmart, issuer, brand
	CallingSystemName       string `json:"callingSystemName,omitempty"`
	OriginalAuthorizationID int    `json:"original_authorization_id"`

	Authorization *AuthorizationInfo `json:"authorization,omitempty"`
	Brand         string             `json:"brand"`
	Card          *AuthCard          `json:"card"`

	OriginalAmount     *AuthAmount         `json:"original_amount"`
	EntryMode          string              `json:"entry_mode,omitempty"`
	CancellationReason *CancellationReason `json:"cancellation_reason,omitempty"`

	ISO8583Message         *OriginalISO8583        `json:"iso8583_message,omitempty"`
	AdditionalTerminalData *AdditionalTerminalData `json:"additionalTerminalData,omitempty"`
	TokenPaymentData       *AuthTokenPaymentData   `json:"tokenPaymentData,omitempty"`

	HCETransaction bool   `json:"hceTransaction,omitempty"`
	EloProductCode string `json:"eloProductCode,omitempty"`
}

// CancellationResponse represents a cancellation response.
type CancellationResponse struct {
	Message         string      `json:"message"`
	Code            int         `json:"code"`
	AuthorizationID int         `json:"authorization_id,omitempty"`
	Balance         *AuthAmount `json:"balance,omitempty"`
}

// ========================================
// Withdrawal Types
// ========================================

// WithdrawalRequest represents a withdrawal authorization request.
// POST /withdrawals - Called by Evertec when a withdrawal needs authorization.
type WithdrawalRequest struct {
	WithdrawalID      string `json:"withdrawal_id"`
	AccountID         string `json:"account_id"`
	PSProductCode     string `json:"psProductCode"`
	PSProductName     string `json:"psProductName,omitempty"`
	CountryCode       string `json:"countryCode"`
	Source            string `json:"source"` // paySmart, issuer, brand
	CallingSystemName string `json:"callingSystemName,omitempty"`

	Authorization *AuthorizationInfo `json:"authorization,omitempty"`
	Brand         string             `json:"brand"`
	Card          *AuthCard          `json:"card"`

	TotalAmount    *AuthAmount `json:"total_amount"`
	OriginalAmount *AuthAmount `json:"original_amount,omitempty"`
	DollarAmount   *AuthAmount `json:"dollar_amount,omitempty"`
	DollarRealRate string      `json:"dollar_real_rate,omitempty"`
	Spread         string      `json:"spread,omitempty"`

	EntryMode            string          `json:"entry_mode"`
	ProcessingCode       *ProcessingCode `json:"processing_code,omitempty"`
	HolderValidationMode string          `json:"holder_validation_mode,omitempty"`

	Fees          []Fee          `json:"fees,omitempty"`
	Establishment *Establishment `json:"establishment,omitempty"`

	Internacional   bool             `json:"internacional,omitempty"`
	OriginalISO8583 *OriginalISO8583 `json:"original_iso8583,omitempty"`

	FraudData              *FraudData              `json:"fraudData,omitempty"`
	AdditionalTerminalData *AdditionalTerminalData `json:"additionalTerminalData,omitempty"`
	TokenPaymentData       *AuthTokenPaymentData   `json:"tokenPaymentData,omitempty"`

	HCETransaction      bool   `json:"hceTransaction,omitempty"`
	EloProductCode      string `json:"eloProductCode,omitempty"`
	ForceAccept         bool   `json:"forceAccept,omitempty"`
	AuthorizationAdvice bool   `json:"authorizationAdvice,omitempty"`
}

// WithdrawalResponse represents a withdrawal authorization response.
type WithdrawalResponse struct {
	Message         string      `json:"message"`
	Code            int         `json:"code"`
	AuthorizationID int         `json:"authorization_id,omitempty"`
	Balance         *AuthAmount `json:"balance,omitempty"`
}

// WithdrawalCancellationRequest represents a withdrawal cancellation request.
// POST /withdrawals/cancel - Called by Evertec to cancel a withdrawal.
type WithdrawalCancellationRequest struct {
	CancellationID          string `json:"cancellation_id"`
	OriginalWithdrawalID    string `json:"original_withdrawal_id"`
	AccountID               string `json:"account_id"`
	PSProductCode           string `json:"psProductCode"`
	PSProductName           string `json:"psProductName,omitempty"`
	CountryCode             string `json:"countryCode"`
	Source                  string `json:"source"`
	CallingSystemName       string `json:"callingSystemName,omitempty"`
	OriginalAuthorizationID int    `json:"original_authorization_id"`

	Authorization *AuthorizationInfo `json:"authorization,omitempty"`
	Brand         string             `json:"brand"`
	Card          *AuthCard          `json:"card"`

	OriginalAmount     *AuthAmount         `json:"original_amount"`
	EntryMode          string              `json:"entry_mode,omitempty"`
	CancellationReason *CancellationReason `json:"cancellation_reason,omitempty"`

	ISO8583Message         *OriginalISO8583        `json:"iso8583_message,omitempty"`
	AdditionalTerminalData *AdditionalTerminalData `json:"additionalTerminalData,omitempty"`
	TokenPaymentData       *AuthTokenPaymentData   `json:"tokenPaymentData,omitempty"`

	HCETransaction bool   `json:"hceTransaction,omitempty"`
	EloProductCode string `json:"eloProductCode,omitempty"`
}

// ========================================
// Query Types
// ========================================

// QueryRequest represents a balance query request.
// POST /queries - Called by Evertec to query account balance.
type QueryRequest struct {
	QueryID           string `json:"query_id"`
	AccountID         string `json:"account_id"`
	PSProductCode     string `json:"psProductCode"`
	CountryCode       string `json:"countryCode"`
	Source            string `json:"source"`
	CallingSystemName string `json:"callingSystemName,omitempty"`

	Authorization *AuthorizationInfo `json:"authorization,omitempty"`
	Brand         string             `json:"brand"`
	Card          *AuthCard          `json:"card"`

	EntryMode       string           `json:"entry_mode,omitempty"`
	OriginalISO8583 *OriginalISO8583 `json:"original_iso8583,omitempty"`
}

// QueryResponse represents a balance query response.
type QueryResponse struct {
	Message         string      `json:"message"`
	Code            int         `json:"code"`
	AuthorizationID int         `json:"authorization_id,omitempty"`
	Balance         *AuthAmount `json:"balance,omitempty"`
}

// WithdrawalQueryRequest represents a withdrawal query request.
// POST /withdrawalQueries - Called by Evertec to query withdrawal limits.
type WithdrawalQueryRequest struct {
	QueryID           string `json:"query_id"`
	AccountID         string `json:"account_id"`
	PSProductCode     string `json:"psProductCode"`
	CountryCode       string `json:"countryCode"`
	Source            string `json:"source"`
	CallingSystemName string `json:"callingSystemName,omitempty"`

	Authorization *AuthorizationInfo `json:"authorization,omitempty"`
	Brand         string             `json:"brand"`
	Card          *AuthCard          `json:"card"`

	EntryMode       string           `json:"entry_mode,omitempty"`
	OriginalISO8583 *OriginalISO8583 `json:"original_iso8583,omitempty"`
}

// WithdrawalQueryResponse represents a withdrawal query response.
type WithdrawalQueryResponse struct {
	Message             string      `json:"message"`
	Code                int         `json:"code"`
	AuthorizationID     int         `json:"authorization_id,omitempty"`
	MaxWithdrawalAmount *AuthAmount `json:"maxWithdrawalAmount,omitempty"`
	DailyLimit          *AuthAmount `json:"dailyLimit,omitempty"`
	RemainingDaily      *AuthAmount `json:"remainingDaily,omitempty"`
}

// ========================================
// Chargeback Types
// ========================================

// ChargebackRequest represents a chargeback request.
// POST /chargebacks - Called by Evertec when a chargeback is initiated.
type ChargebackRequest struct {
	ChargebackID            string `json:"chargeback_id"`
	OriginalPurchaseID      string `json:"original_purchase_id"`
	AccountID               string `json:"account_id"`
	PSProductCode           string `json:"psProductCode"`
	CountryCode             string `json:"countryCode"`
	Source                  string `json:"source"`
	CallingSystemName       string `json:"callingSystemName,omitempty"`
	OriginalAuthorizationID int    `json:"original_authorization_id"`

	Authorization *AuthorizationInfo `json:"authorization,omitempty"`
	Brand         string             `json:"brand"`
	Card          *AuthCard          `json:"card"`

	ChargebackAmount *AuthAmount `json:"chargeback_amount"`
	OriginalAmount   *AuthAmount `json:"original_amount,omitempty"`
	Reason           string      `json:"reason,omitempty"`
	DisputeID        string      `json:"disputeId,omitempty"`
}

// ChargebackResponse represents a chargeback response.
type ChargebackResponse struct {
	Message         string      `json:"message"`
	Code            int         `json:"code"`
	AuthorizationID int         `json:"authorization_id,omitempty"`
	Balance         *AuthAmount `json:"balance,omitempty"`
}

// ChargebackCancellationRequest represents a chargeback cancellation request.
type ChargebackCancellationRequest struct {
	CancellationID          string `json:"cancellation_id"`
	OriginalChargebackID    string `json:"original_chargeback_id"`
	AccountID               string `json:"account_id"`
	PSProductCode           string `json:"psProductCode"`
	CountryCode             string `json:"countryCode"`
	Source                  string `json:"source"`
	CallingSystemName       string `json:"callingSystemName,omitempty"`
	OriginalAuthorizationID int    `json:"original_authorization_id"`

	Authorization      *AuthorizationInfo  `json:"authorization,omitempty"`
	Brand              string              `json:"brand"`
	Card               *AuthCard           `json:"card"`
	OriginalAmount     *AuthAmount         `json:"original_amount"`
	CancellationReason *CancellationReason `json:"cancellation_reason,omitempty"`
}

// ========================================
// Transfer Types
// ========================================

// TransferRequest represents a transfer authorization request.
// POST /transfers - Called by Evertec when a P2P transfer needs authorization.
type TransferRequest struct {
	TransferID        string `json:"transfer_id"`
	AccountID         string `json:"account_id"`
	PSProductCode     string `json:"psProductCode"`
	CountryCode       string `json:"countryCode"`
	Source            string `json:"source"`
	CallingSystemName string `json:"callingSystemName,omitempty"`

	Authorization *AuthorizationInfo `json:"authorization,omitempty"`
	Brand         string             `json:"brand"`

	SourceCard      *AuthCard `json:"sourceCard"`
	DestinationCard *AuthCard `json:"destinationCard,omitempty"`

	TotalAmount *AuthAmount `json:"total_amount"`
	Fees        []Fee       `json:"fees,omitempty"`

	EntryMode       string           `json:"entry_mode,omitempty"`
	OriginalISO8583 *OriginalISO8583 `json:"original_iso8583,omitempty"`
}

// TransferResponse represents a transfer authorization response.
type TransferResponse struct {
	Message         string      `json:"message"`
	Code            int         `json:"code"`
	AuthorizationID int         `json:"authorization_id,omitempty"`
	Balance         *AuthAmount `json:"balance,omitempty"`
}

// TransferCancellationRequest represents a transfer cancellation request.
type TransferCancellationRequest struct {
	CancellationID          string `json:"cancellation_id"`
	OriginalTransferID      string `json:"original_transfer_id"`
	AccountID               string `json:"account_id"`
	PSProductCode           string `json:"psProductCode"`
	CountryCode             string `json:"countryCode"`
	Source                  string `json:"source"`
	CallingSystemName       string `json:"callingSystemName,omitempty"`
	OriginalAuthorizationID int    `json:"original_authorization_id"`

	Authorization      *AuthorizationInfo  `json:"authorization,omitempty"`
	Brand              string              `json:"brand"`
	SourceCard         *AuthCard           `json:"sourceCard"`
	OriginalAmount     *AuthAmount         `json:"original_amount"`
	CancellationReason *CancellationReason `json:"cancellation_reason,omitempty"`
}

// ========================================
// 3DS / ACS Types
// ========================================

// OTPChannelRequest represents a request to get OTP channel.
// POST /acs/getOTPChannel - Called by Evertec for 3DS authentication.
type OTPChannelRequest struct {
	CardID        string `json:"cardId"`
	TransactionID string `json:"transactionId"`
	AccountID     string `json:"account_id,omitempty"`
}

// OTPChannelResponse represents the OTP channel response.
type OTPChannelResponse struct {
	Channel           string `json:"channel"` // SMS, EMAIL, PUSH
	MaskedDestination string `json:"maskedDestination"`
}

// VerifyTransactionRequest represents a 3DS verification request.
// POST /acs/verifyTransaction - Called by Evertec to verify 3DS OTP.
type VerifyTransactionRequest struct {
	CardID        string      `json:"cardId"`
	TransactionID string      `json:"transactionId"`
	AccountID     string      `json:"account_id,omitempty"`
	Amount        *AuthAmount `json:"amount,omitempty"`
	OTP           string      `json:"otp,omitempty"`
}

// VerifyTransactionResponse represents a 3DS verification response.
type VerifyTransactionResponse struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message,omitempty"`
	Code     int    `json:"code,omitempty"`
}

// ========================================
// xPays Types (Apple Pay, Google Pay, Samsung Pay)
// ========================================

// XPaysOTPRequest represents an xPays OTP request.
// POST /xpays/otp - Called by Evertec for wallet provisioning OTP.
type XPaysOTPRequest struct {
	CardID    string       `json:"cardId"`
	AccountID string       `json:"account_id,omitempty"`
	Wallet    types.Wallet `json:"wallet"`
	Channel   string       `json:"channel"` // SMS, EMAIL
}

// XPaysOTPResponse represents an xPays OTP response.
type XPaysOTPResponse struct {
	Status            string `json:"status"` // SENT, FAILED
	MaskedDestination string `json:"maskedDestination"`
	ExpiresAt         string `json:"expiresAt"`
	Code              int    `json:"code,omitempty"`
	Message           string `json:"message,omitempty"`
}

// CustomProvisioningDataRequest represents a custom provisioning data request.
// POST /xpays/customProvisioningData - Called by Evertec for custom wallet data.
type CustomProvisioningDataRequest struct {
	CardID    string       `json:"cardId"`
	AccountID string       `json:"account_id,omitempty"`
	Wallet    types.Wallet `json:"wallet"`
}

// ========================================
// Status Types
// ========================================

// StatusResponse represents a health check response.
// GET /status - Called by Evertec to check if your server is healthy.
type StatusResponse struct {
	Status    string `json:"status"` // OK, DOWN
	Timestamp string `json:"timestamp"`
	Message   string `json:"message,omitempty"`
}

// ========================================
// Response Codes
// ========================================

// ResponseCode represents standard response codes for authorization.
// These codes are returned in the `code` field of responses.
type ResponseCode int

const (
	// Approved
	ResponseCodeApproved ResponseCode = 0 // "00" - Transação aprovada

	// Generic denials
	ResponseCodeGenericDenial  ResponseCode = 5  // Generic denial - retry
	ResponseCodeGenericError   ResponseCode = 6  // Generic message error
	ResponseCodeSecurityBlock  ResponseCode = 7  // Card blocked for security
	ResponseCodeDenied         ResponseCode = 10 // Transaction denied
	ResponseCodeRetry          ResponseCode = 12 // Retry the transaction
	ResponseCodeInvalidAmount  ResponseCode = 13 // Invalid transaction amount
	ResponseCodeInvalidCard    ResponseCode = 14 // Invalid card number
	ResponseCodeIssuerNotFound ResponseCode = 19 // Issuer not found
	ResponseCodeInvalidCancel  ResponseCode = 21 // Invalid cancellation
	ResponseCodeInvalidInstall ResponseCode = 23 // Invalid installment amount
	ResponseCodeInvalidAccount ResponseCode = 25 // Invalid account number

	// Funds and limits
	ResponseCodeInsufficientFunds ResponseCode = 51 // Insufficient funds
	ResponseCodeExceedsLimit      ResponseCode = 61 // Daily limit exceeded
	ResponseCodeTempBlock         ResponseCode = 62 // Temporary collection block
	ResponseCodeAmountLimit       ResponseCode = 64 // Amount limit violation
	ResponseCodePurchaseLimit     ResponseCode = 65 // Purchase limit exceeded

	// Card status
	ResponseCodeLostCard       ResponseCode = 41 // Lost card
	ResponseCodeStolenCard     ResponseCode = 43 // Stolen card
	ResponseCodeExpiredCard    ResponseCode = 54 // Expired card
	ResponseCodeInvalidMCC     ResponseCode = 57 // Transaction not allowed for card
	ResponseCodeInvalidMCCCode ResponseCode = 58 // Invalid MCC
	ResponseCodeRestrictedCard ResponseCode = 78 // Card blocking

	// Fraud
	ResponseCodeFraudSuspect ResponseCode = 59 // Fraud suspicion

	// System errors
	ResponseCodeSystemError   ResponseCode = 91 // System error
	ResponseCodeDuplicateTx   ResponseCode = 94 // Duplicate transaction
	ResponseCodeSystemFailure ResponseCode = 96 // System failure
)

// IsApproved returns true if the response code indicates approval.
func (c ResponseCode) IsApproved() bool {
	return c == ResponseCodeApproved
}

// String returns the two-digit string representation of the code.
func (c ResponseCode) String() string {
	if c == ResponseCodeApproved {
		return "00"
	}
	return ""
}

// ========================================
// Legacy/Backward Compatibility
// ========================================

// CancellationRequest is deprecated. Use PurchaseCancellationRequest, WithdrawalCancellationRequest, etc.
// Kept for backward compatibility.
type CancellationRequest = PurchaseCancellationRequest

// BalanceInfo represents balance information (legacy).
type BalanceInfo struct {
	Available *types.Amount `json:"available"`
	Blocked   *types.Amount `json:"blocked,omitempty"`
	Currency  string        `json:"currency"`
}
