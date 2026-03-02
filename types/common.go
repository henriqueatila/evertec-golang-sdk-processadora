// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// Amount from OpenAPI spec.
type Amount struct {
	Amount       int64 `json:"amount"`
	CurrencyCode int   `json:"currencyCode"`
}

// AmountTransaction from OpenAPI spec.
type AmountTransaction struct {
	Amount        int64  `json:"amount"`
	CurrencyCode  int    `json:"currencyCode"`
	DebitOrCredit string `json:"debit_or_credit,omitempty"`
}

// Address from OpenAPI spec.
type Address struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3,omitempty"`
	Reference    string `json:"reference,omitempty"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zipcode      string `json:"zipcode"`
	Country      string `json:"country,omitempty"`
}

// AddressWithRecipient from OpenAPI spec.
type AddressWithRecipient struct {
	Recipient    string `json:"recipient,omitempty"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3,omitempty"`
	Reference    string `json:"reference,omitempty"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zipcode      string `json:"zipcode"`
	Country      string `json:"country,omitempty"`
}

// Fee from OpenAPI spec.
type Fee struct {
	Amount Amount `json:"amount"`
	Type   string `json:"type"`
}

// ResultData from OpenAPI spec.
type ResultData struct {
	ResultCode        int    `json:"resultCode"`
	ResultDescription string `json:"resultDescription"`
	IssuerRequestID   string `json:"issuerRequestId,omitempty"`
	PsResponseID      string `json:"psResponseId"`
}

// ContactInformation from OpenAPI spec.
type ContactInformation struct {
	PersonalPhoneNumber1  string `json:"personalPhoneNumber1"`
	PersonalPhoneNumber2  string `json:"personalPhoneNumber2,omitempty"`
	ComercialPhoneNumber1 string `json:"comercialPhoneNumber1,omitempty"`
	Email                 string `json:"email"`
}

// PersonalIdentityDocumentInfo from OpenAPI spec.
type PersonalIdentityDocumentInfo struct {
	IdentityDocumentNumber string `json:"identityDocumentNumber"`
	State                  string `json:"state"`
	IssuedBy               string `json:"issuedBy"`
}

// OccupationInfo from OpenAPI spec.
type OccupationInfo struct {
	CompanyName    string         `json:"companyName,omitempty"`
	JobTitle       string         `json:"jobTitle,omitempty"`
	CurrentJobTime int            `json:"currentJobTime,omitempty"`
	Income         map[string]any `json:"income,omitempty"`
	IncomeRange    string         `json:"incomeRange,omitempty"`
}

// CompanyInfo from OpenAPI spec.
type CompanyInfo struct {
	IdentityDocumentNumber string `json:"identityDocumentNumber,omitempty"`
	Name                   string `json:"name,omitempty"`
	DepartmentCode         string `json:"departmentCode,omitempty"`
	ContactPerson          string `json:"contactPerson,omitempty"`
}

// BankAccountInfo from OpenAPI spec.
type BankAccountInfo struct {
	BankID        string `json:"bankID,omitempty"`
	BankName      string `json:"bankName,omitempty"`
	AgencyCode    int    `json:"agencyCode,omitempty"`
	AccountNumber int    `json:"accountNumber,omitempty"`
}

// SourceAudit - Informações de auditoria para registro de informações. São de extrema
type SourceAudit struct {
	OperatorID string `json:"operatorId,omitempty"`
	ProcessID  string `json:"processId,omitempty"`
}

// DebitOrCreditAmount from OpenAPI spec.
type DebitOrCreditAmount struct {
	Amount        int64  `json:"amount"`
	CurrencyCode  int    `json:"currencyCode"`
	DebitOrCredit string `json:"debit_or_credit"`
}

// CET from OpenAPI spec.
type CET struct {
	MonthlyInterest string `json:"monthlyInterest"`
	YearlyInterest  string `json:"yearlyInterest"`
	IOF             string `json:"IOF"`
	DailyIOF        string `json:"DailyIOF,omitempty"`
	CET             string `json:"CET"`
}

// CETLateness from OpenAPI spec.
type CETLateness struct {
	MonthlyInterest string `json:"monthlyInterest"`
	LatenessFee     string `json:"latenessFee"`
	YearlyInterest  string `json:"yearlyInterest"`
	IOF             string `json:"IOF"`
	DailyIOF        string `json:"DailyIOF,omitempty"`
	CET             string `json:"CET"`
}

// FeesAndCET from OpenAPI spec.
type FeesAndCET struct {
	RevolvingCredit               CET         `json:"revolvingCredit"`
	StatementInstallmentFinancing *CET        `json:"statementInstallmentFinancing,omitempty"`
	CashWithdrawal                *CET        `json:"cashWithdrawal,omitempty"`
	LatenessInterest              CETLateness `json:"latenessInterest"`
}

// DailyFeeEntry from OpenAPI spec.
type DailyFeeEntry struct {
	Day           string         `json:"day,omitempty"`
	Financed      int            `json:"financed,omitempty"`
	DailyIOF      map[string]any `json:"daily_iof,omitempty"`
	DailyInterest map[string]any `json:"daily_interest,omitempty"`
	ExpectedFees  map[string]any `json:"expected_fees,omitempty"`
	Lateness      map[string]any `json:"lateness,omitempty"`
	RefundFees    map[string]any `json:"refund_fees,omitempty"`
}

// ProductsDTO from OpenAPI spec.
type ProductsDTO struct {
	ProductID                   int      `json:"productId,omitempty"`
	ClosedArrangementEntryModes []string `json:"closedArrangementEntryModes,omitempty"`
}

// BrandDTO from OpenAPI spec.
type BrandDTO struct {
	BrandID   int    `json:"brandId,omitempty"`
	BrandName string `json:"brandName,omitempty"`
}

// BrandErrorDTO from OpenAPI spec.
type BrandErrorDTO struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

// DisputeReason is a string type for dispute reasons.
type DisputeReason string

// MerchantAccountInfo is merchant account information array.
type MerchantAccountInfo []map[string]any

// QrCodeTransactionInfo is QR code transaction info array.
type QrCodeTransactionInfo []map[string]any
