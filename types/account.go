// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// Account from OpenAPI spec.
type Account struct {
	AccountID           string         `json:"accountId,omitempty"`
	PsProductCode       string         `json:"psProductCode"`
	PsProductName       string         `json:"psProductName,omitempty"`
	IssuerAccountID     string         `json:"issuerAccountId,omitempty"`
	AccountOwner        interface{}    `json:"accountOwner"`
	CreditInfo          *CreditInfo    `json:"creditInfo,omitempty"`
	Cards               []Card         `json:"cards,omitempty"`
	Status              *AccountStatus `json:"status,omitempty"`
	BillingAddress      *Address       `json:"billingAddress,omitempty"`
	CardDeliveryAddress *Address       `json:"cardDeliveryAddress,omitempty"`
	InclusionDate       string         `json:"inclusion_date,omitempty"`
}

// AccountOwnerData from OpenAPI spec.
type AccountOwnerData struct {
	FullName                    string                        `json:"fullName"`
	IDentityDocumentNumber      string                        `json:"identityDocumentNumber"`
	OtherIDentityDocumentNumber *PersonalIdentityDocumentInfo `json:"otherIdentityDocumentNumber,omitempty"`
	MotherName                  string                        `json:"motherName,omitempty"`
	BirthDate                   string                        `json:"birthDate,omitempty"`
	ContactInformation          ContactInformation            `json:"contactInformation"`
	OccupationInfo              *OccupationInfo               `json:"occupationInfo,omitempty"`
}

// AccountOwnerDataNoDoc from OpenAPI spec.
type AccountOwnerDataNoDoc struct {
	FullName           string              `json:"fullName,omitempty"`
	MotherName         string              `json:"motherName,omitempty"`
	ContactInformation *ContactInformation `json:"contactInformation,omitempty"`
	OccupationInfo     *OccupationInfo     `json:"occupationInfo,omitempty"`
}

// AccountBalance from OpenAPI spec (Limits API).
// GET /accounts/{accountId}/balance
// Reference: https://paysmart-api.gitlab.io/processadora/PT-br/docs/limits-api
// Note: Field "payamentDue" spelling matches official API (typo in original spec)
type AccountBalance struct {
	AccountID                    string  `json:"accountId"`
	Balance                      *Amount `json:"balance,omitempty"`
	WithdrawalBalance            *Amount `json:"withdrawalBalance,omitempty"`
	CreditLimit                  *Amount `json:"creditLimit,omitempty"`
	UsableCreditLimit            *Amount `json:"usableCreditLimit,omitempty"`
	WithdrawalCreditLimit        *Amount `json:"withdrawalCreditLimit,omitempty"`
	CurrentCreditLimit           *Amount `json:"currentCreditLimit,omitempty"`
	CurrentUsableCreditLimit     *Amount `json:"currentUsableCreditLimit,omitempty"`
	CurrentWithdrawalCreditLimit *Amount `json:"currentWithdrawalCreditLimit,omitempty"`
	PayamentDue                  *Amount `json:"payamentDue,omitempty"` // Note: API typo preserved
	QueryDate                    string  `json:"query_date,omitempty"`
}

// AccountListResult from OpenAPI spec.
type AccountListResult struct {
	HasMore  bool      `json:"hasMore,omitempty"`
	Accounts []Account `json:"accounts,omitempty"`
}

// AccountObject from OpenAPI spec.
type AccountObject struct {
	ProductCode             string  `json:"productCode,omitempty"`
	CampaignID              string  `json:"campaignId,omitempty"`
	AgentID                 string  `json:"agentId,omitempty"`
	LatePeriod              string  `json:"latePeriod,omitempty"`
	LateStartPeriod         string  `json:"lateStartPeriod,omitempty"`
	LastPaymentDate         string  `json:"lastPaymentDate,omitempty"`
	AccountID               string  `json:"accountId,omitempty"`
	FullName                string  `json:"fullName,omitempty"`
	DocID                   string  `json:"docId,omitempty"`
	Address                 string  `json:"address,omitempty"`
	Phone                   string  `json:"phone,omitempty"`
	Email                   string  `json:"email,omitempty"`
	PaymentDue              string  `json:"paymentDue,omitempty"`
	MinimumPayment          float64 `json:"minimumPayment,omitempty"`
	CurrentStatementBalance float64 `json:"currentStatementBalance,omitempty"`
	PendingMinimumPayment   float64 `json:"pendingMinimumPayment,omitempty"`
	FinancedBalance         float64 `json:"financedBalance,omitempty"`
	CurrentOnDebitbalance   float64 `json:"currentOnDebitbalance,omitempty"`
	TotalOnDebitBalance     float64 `json:"totalOnDebitBalance,omitempty"`
}

// AccountBlocking from OpenAPI spec.
type AccountBlocking struct {
	BlockingCode    int    `json:"blockingCode,omitempty"`
	BlockedAt       string `json:"blockedAt,omitempty"`
	BlockReason     string `json:"blockReason,omitempty"`
	UnblockedAt     string `json:"unblockedAt,omitempty"`
	UnblockReason   string `json:"unblockReason,omitempty"`
	Description     string `json:"description,omitempty"`
	Category        string `json:"category,omitempty"`
	ResponsibleTeam string `json:"responsibleTeam,omitempty"`
}

// AccountBlockingConfig from OpenAPI spec.
type AccountBlockingConfig struct {
	Enable       bool   `json:"enable,omitempty"`
	BlockingCode int    `json:"blockingCode,omitempty"`
	Reason       string `json:"reason,omitempty"`
	BlockedAt    string `json:"blockedAt,omitempty"`
}

// AccountBlockingSummary from OpenAPI spec.
type AccountBlockingSummary struct {
	AccountStatus          *AccountStatus         `json:"accountStatus,omitempty"`
	Authorization          *AccountBlockingConfig `json:"authorization,omitempty"`
	Credits                *AccountBlockingConfig `json:"credits,omitempty"`
	Debits                 *AccountBlockingConfig `json:"debits,omitempty"`
	CardReissuing          *AccountBlockingConfig `json:"cardReissuing,omitempty"`
	OverLimitAuthorization *AccountBlockingConfig `json:"overLimitAuthorization,omitempty"`
	ActiveBlockingCodes    []int                  `json:"activeBlockingCodes,omitempty"`
}

// NewAccountRequest from OpenAPI spec.
type NewAccountRequest struct {
	IssuerRequestID       string       `json:"issuerRequestId,omitempty"`
	PsProductCode         string       `json:"psProductCode"`
	IssuerAccountID       string       `json:"issuerAccountId,omitempty"`
	AccountOwner          interface{}  `json:"accountOwner"`
	CreditInfo            *CreditInfo  `json:"creditInfo,omitempty"`
	BillingAddress        interface{}  `json:"billingAddress"`
	CardDeliveryAddress   interface{}  `json:"cardDeliveryAddress"`
	RequestingCompanyInfo interface{}  `json:"requestingCompanyInfo,omitempty"`
	BankAccount           interface{}  `json:"bankAccount,omitempty"`
	SourceAudit           *SourceAudit `json:"sourceAudit,omitempty"`
}

// AccountCreatedSuccessfully from OpenAPI spec.
type AccountCreatedSuccessfully struct {
	ResultData interface{} `json:"resultData"`
	Account    Account     `json:"account"`
}

// UpdateAccountRequest from OpenAPI spec.
type UpdateAccountRequest struct {
	IssuerRequestID       string       `json:"issuerRequestId,omitempty"`
	IssuerAccountID       string       `json:"issuerAccountId,omitempty"`
	PsProductCode         string       `json:"psProductCode,omitempty"`
	AccountOwner          interface{}  `json:"accountOwner,omitempty"`
	CreditInfo            *CreditInfo  `json:"creditInfo,omitempty"`
	BillingAddress        interface{}  `json:"billingAddress,omitempty"`
	CardDeliveryAddress   interface{}  `json:"cardDeliveryAddress,omitempty"`
	RequestingCompanyInfo interface{}  `json:"requestingCompanyInfo,omitempty"`
	SourceAudit           *SourceAudit `json:"sourceAudit,omitempty"`
}

// BlockAccountRequest from OpenAPI spec.
type BlockAccountRequest struct {
	BlockCode int    `json:"blockCode"`
	Reason    string `json:"reason,omitempty"`
}

// BlockAccountResult from OpenAPI spec.
type BlockAccountResult struct {
	ResultData      interface{}             `json:"resultData,omitempty"`
	AppliedBlocking *AccountBlocking        `json:"appliedBlocking,omitempty"`
	BlockingSummary *AccountBlockingSummary `json:"blockingSummary,omitempty"`
	BlockedAccount  *Account                `json:"blockedAccount,omitempty"`
}

// UnblockAccountRequest from OpenAPI spec.
type UnblockAccountRequest struct {
	UnblockCode int    `json:"unblockCode"`
	Reason      string `json:"reason,omitempty"`
}

// UnblockAccountResult from OpenAPI spec.
type UnblockAccountResult struct {
	ResultData       interface{}             `json:"resultData,omitempty"`
	RemovedBlocking  *AccountBlocking        `json:"removedBlocking,omitempty"`
	BlockingSummary  *AccountBlockingSummary `json:"blockingSummary,omitempty"`
	UnblockedAccount *Account                `json:"unblockedAccount,omitempty"`
}

// CancelAccountRequest from OpenAPI spec.
type CancelAccountRequest struct {
	CancellationCode int          `json:"cancellationCode"`
	Reason           string       `json:"reason,omitempty"`
	SourceAudit      *SourceAudit `json:"sourceAudit,omitempty"`
}

// AccountCancelledSuccessfully from OpenAPI spec.
type AccountCancelledSuccessfully struct {
	ResultData          interface{}             `json:"resultData,omitempty"`
	AppliedCancellation *AccountBlocking        `json:"appliedCancellation,omitempty"`
	BlockingSummary     *AccountBlockingSummary `json:"blockingSummary,omitempty"`
	CancelledAccount    *Account                `json:"cancelledAccount,omitempty"`
}

// CreditInfo from OpenAPI spec.
type CreditInfo struct {
	CreditLimit     interface{} `json:"creditLimit,omitempty"`
	WithdrawalLimit interface{} `json:"withdrawalLimit,omitempty"`
	PaymentDue      int         `json:"paymentDue,omitempty"`
}

// GetCreditAnalysisResponse from OpenAPI spec.
type GetCreditAnalysisResponse struct {
	ResultData           interface{} `json:"resultData,omitempty"`
	AccountID            string      `json:"accountId,omitempty"`
	Status               string      `json:"status,omitempty"`
	RequestedLimit       int         `json:"requestedLimit,omitempty"`
	MaximumApprovedLimit int         `json:"maximumApprovedLimit,omitempty"`
}
