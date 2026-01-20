// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// Transaction from OpenAPI spec.
type Transaction struct {
	TransactionID                            string                            `json:"transactionId"`
	AccountID                                string                            `json:"accountId"`
	TransactionDateTime                      string                            `json:"transactionDateTime"`
	SettlementDateTime                       string                            `json:"settlementDateTime,omitempty"`
	TransactionType                          TransactionType                   `json:"transactionType"`
	TransactionStatus                        TransactionStatus                 `json:"transactionStatus"`
	TransactionAuthorizationResponse         *TransactionAuthorizationResponse `json:"transactionAuthorizationResponse,omitempty"`
	AcquirerID                               string                            `json:"acquirerId,omitempty"`
	AcquirerTransactionID                    string                            `json:"acquirerTransactionId,omitempty"`
	MerchantID                               string                            `json:"merchantId,omitempty"`
	MerchantName                             string                            `json:"merchantName,omitempty"`
	MerchantDocumentID                       string                            `json:"merchantDocumentId,omitempty"`
	MerchantAddress                          string                            `json:"merchantAddress,omitempty"`
	MerchantCity                             string                            `json:"merchantCity,omitempty"`
	MerchantUf                               string                            `json:"merchantUf,omitempty"`
	CountryCode                              string                            `json:"countryCode,omitempty"`
	MCC                                      string                            `json:"mcc,omitempty"`
	TerminalID                               string                            `json:"terminalId,omitempty"`
	Amount                                   AmountTransaction                 `json:"amount"`
	RemainingAmount                          *Amount                           `json:"remainingAmount,omitempty"`
	RefundedAmount                           *Amount                           `json:"refundedAmount,omitempty"`
	EntryMode                                *EntryMode                        `json:"entryMode,omitempty"`
	Card                                     interface{}                       `json:"card,omitempty"`
	CancellingTransactionID                  string                            `json:"cancellingTransactionId,omitempty"`
	CancellingTransactionIDs                 []string                          `json:"cancellingTransactionIds,omitempty"`
	International                            bool                              `json:"international,omitempty"`
	InternationalTransactionData             interface{}                       `json:"internationalTransactionData,omitempty"`
	AuthorizationAdvice                      bool                              `json:"authorizationAdvice,omitempty"`
	MITAdditionalData                        *MITAdditionalData                `json:"mitAdditionalData,omitempty"`
	AdditionalTerminalData                   *AdditionalTerminalData           `json:"additionalTerminalData,omitempty"`
	TokenPaymentData                         *TokenPaymentData                 `json:"tokenPaymentData,omitempty"`
	InstallmentDetails                       *InstallmentDetails               `json:"installmentDetails,omitempty"`
	PsProductCode                            string                            `json:"psProductCode,omitempty"`
	ProductType                              string                            `json:"productType,omitempty"`
	PsProductName                            string                            `json:"psProductName,omitempty"`
	MerchantZipcode                          string                            `json:"merchantZipcode,omitempty"`
	MerchantPAT                              bool                              `json:"merchantPAT,omitempty"`
	TransactionSource                        *Source                           `json:"transactionSource,omitempty"`
	TransactionDate                          string                            `json:"transactionDate,omitempty"`
	TransactionTime                          string                            `json:"transactionTime,omitempty"`
	Incremental                              bool                              `json:"incremental,omitempty"`
	DisputeID                                string                            `json:"disputeId,omitempty"`
	DisputeStatus                            *DisputeStatus                    `json:"disputeStatus,omitempty"`
	Fees                                     []Fee                             `json:"fees,omitempty"`
	Installments                             []Installment                     `json:"installments,omitempty"`
	UndoData                                 interface{}                       `json:"undoData,omitempty"`
	OriginalTransaction                      string                            `json:"originalTransaction,omitempty"`
	TransactionAuthorizationResponsePaysmart map[string]interface{}            `json:"transactionAuthorizationResponsePaysmart,omitempty"`
	TransferData                             *TransferData                     `json:"transferData,omitempty"`
	FraudData                                *FraudData                        `json:"fraudData,omitempty"`
	SettlementDate                           string                            `json:"settlementDate,omitempty"`
	HCETransaction                           bool                              `json:"hceTransaction,omitempty"`
	CancellationReason                       string                            `json:"cancellationReason,omitempty"`
	PurchaseWithCashbackData                 *PurchaseWithCashbackData         `json:"purchaseWithCashbackData,omitempty"`
	PartialApprovalData                      *PartialApprovalData              `json:"partialApprovalData,omitempty"`
	ISO8583MessageRequest                    *ISO8583MessageData               `json:"ISO8583MessageRequest,omitempty"`
	ISO8583MessageResponse                   *ISO8583MessageData               `json:"ISO8583MessageResponse,omitempty"`
	ELOProductCode                           string                            `json:"EloProductCode,omitempty"`
	AfdTransaction                           string                            `json:"afdTransaction,omitempty"`
	TransitTransaction                       string                            `json:"transitTransaction,omitempty"`
	NRID                                     string                            `json:"nrid,omitempty"`
}

// TransactionListResult from OpenAPI spec.
type TransactionListResult struct {
	HasMore      bool                     `json:"hasMore,omitempty"`
	Transactions []TransactionQueryResult `json:"transactions,omitempty"`
}

// TransactionQueryResult from OpenAPI spec.
type TransactionQueryResult struct {
	TransactionID                    string                            `json:"transactionId"`
	AccountID                        string                            `json:"accountId"`
	TransactionDateTime              string                            `json:"transactionDateTime"`
	SettlementDateTime               string                            `json:"settlementDateTime,omitempty"`
	TransactionType                  TransactionType                   `json:"transactionType"`
	TransactionStatus                TransactionStatus                 `json:"transactionStatus"`
	TransactionAuthorizationResponse *TransactionAuthorizationResponse `json:"transactionAuthorizationResponse,omitempty"`
	AcquirerID                       string                            `json:"acquirerId,omitempty"`
	AcquirerTransactionID            string                            `json:"acquirerTransactionId,omitempty"`
	MerchantID                       string                            `json:"merchantId,omitempty"`
	MerchantName                     string                            `json:"merchantName,omitempty"`
	MerchantDocumentID               string                            `json:"merchantDocumentId,omitempty"`
	MerchantAddress                  string                            `json:"merchantAddress,omitempty"`
	MerchantCity                     string                            `json:"merchantCity,omitempty"`
	MerchantUf                       string                            `json:"merchantUf,omitempty"`
	MerchantPAT                      bool                              `json:"merchantPAT,omitempty"`
	CountryCode                      string                            `json:"countryCode,omitempty"`
	MCC                              string                            `json:"mcc,omitempty"`
	TerminalID                       string                            `json:"terminalId,omitempty"`
	Amount                           Amount                            `json:"amount"`
	RemainingAmount                  *Amount                           `json:"remainingAmount,omitempty"`
	RefundedAmount                   *Amount                           `json:"refundedAmount,omitempty"`
	EntryMode                        *EntryMode                        `json:"entryMode,omitempty"`
	Card                             interface{}                       `json:"card,omitempty"`
	CancellingTransactionID          string                            `json:"cancellingTransactionId,omitempty"`
	CancellingTransactionIDs         []string                          `json:"cancellingTransactionIds,omitempty"`
	International                    bool                              `json:"international,omitempty"`
	InternationalTransactionData     interface{}                       `json:"internationalTransactionData,omitempty"`
	AuthorizationAdvice              bool                              `json:"authorizationAdvice,omitempty"`
	MITAdditionalData                *MITAdditionalData                `json:"mitAdditionalData,omitempty"`
	AdditionalTerminalData           *AdditionalTerminalData           `json:"additionalTerminalData,omitempty"`
	TokenPaymentData                 *TokenPaymentData                 `json:"tokenPaymentData,omitempty"`
	InstallmentDetails               *InstallmentDetails               `json:"installmentDetails,omitempty"`
	HCETransaction                   bool                              `json:"hceTransaction,omitempty"`
	CancellationReason               string                            `json:"cancellationReason,omitempty"`
	PurchaseWithCashbackData         *PurchaseWithCashbackData         `json:"purchaseWithCashbackData,omitempty"`
	PartialApprovalData              *PartialApprovalData              `json:"partialApprovalData,omitempty"`
	ISO8583MessageRequest            *ISO8583MessageData               `json:"ISO8583MessageRequest,omitempty"`
	ISO8583MessageResponse           *ISO8583MessageData               `json:"ISO8583MessageResponse,omitempty"`
	ELOProductCode                   string                            `json:"EloProductCode,omitempty"`
	AfdTransaction                   string                            `json:"afdTransaction,omitempty"`
	TransitTransaction               string                            `json:"transitTransaction,omitempty"`
	NRID                             string                            `json:"nrid,omitempty"`
}

// TransactionDataResult from OpenAPI spec.
type TransactionDataResult struct {
	TransactionID string  `json:"transactionId,omitempty"`
	AccountID     string  `json:"accountId,omitempty"`
	Amount        *Amount `json:"amount,omitempty"`
	DebitOrCredit string  `json:"debit_or_credit,omitempty"`
	Type          string  `json:"type,omitempty"`
}

// TransactionInfo from OpenAPI spec.
type TransactionInfo struct {
	ID                 string `json:"id,omitempty"`
	Type               string `json:"type,omitempty"`
	ReferenceLabel     string `json:"referenceLabel,omitempty"`
	ApprovalTimestamp  string `json:"approvalTimestamp,omitempty"`
	AuthorizationCode  string `json:"authorizationCode,omitempty"`
	AuthenticationCode string `json:"authenticationCode,omitempty"`
	HostNSU            string `json:"hostNsu,omitempty"`
	TerminalNSU        string `json:"terminalNsu,omitempty"`
	MerchantID         string `json:"merchantId,omitempty"`
	TerminalID         string `json:"terminalId,omitempty"`
	Amount             int64  `json:"amount,omitempty"`
	Status             string `json:"status,omitempty"`
}

// TransactionAuthorizationResponse - Resposta do autorizador para a transação
type TransactionAuthorizationResponse struct {
	Approved          bool        `json:"approved,omitempty"`
	PartiallyApproved bool        `json:"partiallyApproved,omitempty"`
	DenialReason      interface{} `json:"denialReason,omitempty"`
}

// TransactionDenialReason from OpenAPI spec.
type TransactionDenialReason struct {
	DenialCode        string `json:"denialCode,omitempty"`
	DenialDescription string `json:"denialDescription,omitempty"`
}

// TransactionCreatedSuccessfully from OpenAPI spec.
type TransactionCreatedSuccessfully struct {
	ResultData      interface{}           `json:"resultData"`
	TransactionData TransactionDataResult `json:"transactionData"`
}

// AmountFraud from OpenAPI spec.
type AmountFraud struct {
	Amount        map[string]interface{} `json:"amount,omitempty"`
	DebitOrCredit string                 `json:"debit_or_credit,omitempty"`
}

// InstallmentDetails from OpenAPI spec.
type InstallmentDetails struct {
	FinType                     string  `json:"finType"`
	FareAmount                  *Amount `json:"fare_amount,omitempty"`
	InsuranceAmount             *Amount `json:"insurance_amount,omitempty"`
	ThirdPartiesPaymntAmount    *Amount `json:"third_parties_paymnt_amount,omitempty"`
	RecordsPaymentsAmount       *Amount `json:"records_payments_amount,omitempty"`
	IssuerTotalCalculatedAmount *Amount `json:"issuer_total_calculated_amount,omitempty"`
	FirstPaymntDate             string  `json:"first_paymnt_date,omitempty"`
	InstalmntNbr                int     `json:"instalmnt_nbr"`
	MonthlyInterestRate         int     `json:"monthly_interest_rate,omitempty"`
	TotalEffectiveCostRateCET   int     `json:"total_effective_cost_rate_cet,omitempty"`
	InstalmntAmount             Amount  `json:"instalmnt_amount"`
	AnnualInterestRate          int     `json:"annual_interest_rate,omitempty"`
	InputValue                  *Amount `json:"input_value,omitempty"`
}

// TokenPaymentData - Informações dos Dados Transacionais - Conjunto de dados de Pagamento v
type TokenPaymentData struct {
	FPAN                                         string `json:"fPan,omitempty"`
	RequesterIDToken                             string `json:"requesterIdToken,omitempty"`
	PAN                                          string `json:"pan,omitempty"`
	TokenPSN                                     string `json:"tokenPSN,omitempty"`
	TokenExpirationDate                          string `json:"tokenExpirationDate,omitempty"`
	TokenStatus                                  string `json:"tokenStatus,omitempty"`
	TokenCryptogramVerificationResults           string `json:"tokenCryptogramVerificationResults,omitempty"`
	EMVTokenCryptogramVerificationResults        string `json:"EMVTokenCryptogramVerificationResults,omitempty"`
	TokenConstraintsVerificationStatus           string `json:"tokenConstraintsVerificationStatus,omitempty"`
	TransactionDateTimeConstraint                string `json:"transactionDateTimeConstraint,omitempty"`
	TransactionAmountConstraint                  string `json:"transactionAmountConstraint,omitempty"`
	UsageConstraint                              string `json:"usageConstraint,omitempty"`
	TokenATCVerificationResults                  string `json:"tokenATCVerificationResults,omitempty"`
	CVE2TokenCryptogramVerificationStatus        string `json:"CVE2TokenCryptogramVerificationStatus,omitempty"`
	MerchantVerification                         string `json:"merchantVerification,omitempty"`
	MagstripeTokenCryptogramVerificationResults  string `json:"magstripeTokenCryptogramVerificationResults,omitempty"`
	CVE2OutputTokenCryptogramVerificationResults string `json:"CVE2OutputTokenCryptogramVerificationResults,omitempty"`
}

// MITAdditionalData - Este conjunto de dados é utilizado para passar dados de identificação
type MITAdditionalData struct {
	InstallmentTotalNbr  string `json:"installmentTotalNbr,omitempty"`
	PaymentType          string `json:"paymentType,omitempty"`
	TransactionType      string `json:"transactionType,omitempty"`
	TransactionAmountMIT int64  `json:"transactionAmountMIT,omitempty"`
	UniqueTransactionID  string `json:"uniqueTransactionID,omitempty"`
	TransactionFrequency string `json:"transactionFrequency,omitempty"`
	ValidationIndicator  string `json:"validationIndicator,omitempty"`
	ValidationReference  string `json:"validationReference,omitempty"`
	SequenceIndicator    int    `json:"sequenceIndicator,omitempty"`
}

// AdditionalTerminalData - Contempla os dados adicionais pertencentes ao ponto de venda/saque (po
type AdditionalTerminalData struct {
	TerminalType                   string `json:"terminalType,omitempty"`
	PartialApprovalIndicator       string `json:"partialApprovalIndicator,omitempty"`
	TerminalLocationIndicator      string `json:"terminalLocationIndicator,omitempty"`
	CardholderPresenceIndicator    string `json:"cardholderPresenceIndicator,omitempty"`
	CardPresenceIndicator          string `json:"cardPresenceIndicator,omitempty"`
	CardCaptureCapabilityIndicator string `json:"cardCaptureCapabilityIndicator,omitempty"`
	TransactionStatusIndicator     string `json:"transactionStatusIndicator,omitempty"`
	TransactionSecurityIndicator   string `json:"transactionSecurityIndicator,omitempty"`
	TerminalPOSType                string `json:"terminalPOSType,omitempty"`
	TerminalInputCapability        string `json:"terminalInputCapability,omitempty"`
}

// AdditionalDataInfo from OpenAPI spec.
type AdditionalDataInfo struct {
	ReferenceLabel          string `json:"referenceLabel"`
	PaymentSpecificTemplate string `json:"paymentSpecificTemplate,omitempty"`
}

// NewTransactionRequest from OpenAPI spec.
type NewTransactionRequest struct {
	IssuerTransactionID string          `json:"issuerTransactionId,omitempty"`
	TransactionType     TransactionType `json:"transactionType"`
	SourceAudit         *SourceAudit    `json:"sourceAudit,omitempty"`
	Amount              interface{}     `json:"amount"`
	Reason              string          `json:"reason"`
}

// NewCreditRequest from OpenAPI spec.
type NewCreditRequest struct {
	IssuerRequestID      string       `json:"issuerRequestId,omitempty"`
	Description          string       `json:"description,omitempty"`
	InclusionDate        string       `json:"inclusion_date,omitempty"`
	EffectivePaymentDate string       `json:"effectivePaymentDate,omitempty"`
	Amount               interface{}  `json:"amount"`
	CreditType           int          `json:"credit_type,omitempty"`
	Type                 string       `json:"type,omitempty"`
	SourceAudit          *SourceAudit `json:"sourceAudit,omitempty"`
}

// NewDebitRequest from OpenAPI spec.
type NewDebitRequest struct {
	IssuerRequestID string       `json:"issuerRequestId,omitempty"`
	Description     string       `json:"description,omitempty"`
	Amount          interface{}  `json:"amount"`
	Type            string       `json:"type"`
	SourceAudit     *SourceAudit `json:"sourceAudit,omitempty"`
}

// Credit from OpenAPI spec.
type Credit struct {
	CreditID      string      `json:"creditId"`
	InclusionDate string      `json:"inclusion_date"`
	Amount        interface{} `json:"amount"`
	CreditType    int         `json:"credit_type"`
}

// CreditCreatedSuccessfully from OpenAPI spec.
type CreditCreatedSuccessfully struct {
	ResultData interface{} `json:"resultData"`
	Credit     Credit      `json:"credit"`
}

// CreditListResult from OpenAPI spec.
type CreditListResult struct {
	HasMore bool     `json:"hasMore,omitempty"`
	Credits []Credit `json:"credits,omitempty"`
}

// ScheduledTransactionsResult from OpenAPI spec.
type ScheduledTransactionsResult struct {
	ResultData interface{}                `json:"resultData,omitempty"`
	Data       *ScheduledTransactionsList `json:"data,omitempty"`
}

// ScheduledTransactionsList from OpenAPI spec.
type ScheduledTransactionsList struct {
	ScheduledTransactions map[string]interface{} `json:"scheduled_transactions"`
}

// ScheduledTransactionsErrorResult from OpenAPI spec.
type ScheduledTransactionsErrorResult struct {
	ResultData interface{} `json:"resultData,omitempty"`
}

// FraudData - A presença deste campo indica que um Score de Fraude é entregue ao Emi
type FraudData struct {
	CreditorFraudScore          int    `json:"creditorFraudScore,omitempty"`
	ELOBrandFraudScore          int    `json:"eloBrandFraudScore,omitempty"`
	FraudScorePrimaryReason     int    `json:"fraudScorePrimaryReason,omitempty"`
	FraudScoreSecondaryReason   int    `json:"fraudScoreSecondaryReason,omitempty"`
	FraudScoreTertiaryReason    int    `json:"fraudScoreTertiaryReason,omitempty"`
	FraudDecisionRecommendation string `json:"fraudDecisionRecommendation,omitempty"`
	ScoreOriginIndicator        int    `json:"scoreOriginIndicator,omitempty"`
}

// TransferData from OpenAPI spec.
type TransferData struct {
	PaymentType                    PaymentType `json:"paymentType"`
	UniqueReferenceNumber          string      `json:"uniqueReferenceNumber,omitempty"`
	SendersName                    string      `json:"sendersName,omitempty"`
	SendersAddress                 string      `json:"sendersAddress,omitempty"`
	SendersCity                    string      `json:"sendersCity,omitempty"`
	CountryOrStateCodeIfUS         string      `json:"countryOrStateCodeIfUS,omitempty"`
	CardholderZipcode              string      `json:"cardholderZipcode,omitempty"`
	CardholderIDentificationNumber string      `json:"cardholderIdentificationNumber,omitempty"`
	OriginOfFunds                  string      `json:"originOfFunds,omitempty"`
	BirthDateOfSender              string      `json:"birthDateOfSender,omitempty"`
	RecipientsName                 string      `json:"recipientsName,omitempty"`
	AdditionalTransferData         string      `json:"additionalTransferData,omitempty"`
	RecipientCode                  string      `json:"recipientCode,omitempty"`
	FundSenderEmail                string      `json:"fundSenderEmail,omitempty"`
	FundRecipientEmail             string      `json:"fundRecipientEmail,omitempty"`
	FundSenderPhone                string      `json:"fundSenderPhone,omitempty"`
	FundRecipientPhone             string      `json:"fundRecipientPhone,omitempty"`
	DeviceID                       string      `json:"deviceID,omitempty"`
	CardholderCPFOrCNPJ            string      `json:"cardholderCpfOrCnpj,omitempty"`
	BINOrigin                      string      `json:"BINOrigin,omitempty"`
	BINDestination                 string      `json:"BINDestination,omitempty"`
	OriginCardLast4digits          string      `json:"originCardLast4digits,omitempty"`
	DestinationCardLast4digits     string      `json:"destinationCardLast4digits,omitempty"`
}

// PurchaseWithCashbackData - Informações dos Dados de compra com devolução de dinheiro(Cashback).
type PurchaseWithCashbackData struct {
	PurchaseOnlyApprovalSupport bool   `json:"purchaseOnlyApprovalSupport"`
	PurchaseOnlyAmount          Amount `json:"purchaseOnlyAmount"`
	CashbackAmount              Amount `json:"cashbackAmount"`
	PurchaseOnlyApproval        bool   `json:"purchaseOnlyApproval"`
}

// PartialApprovalData - Dados de aprovação parcial para compra e compra com devolução de dinhe
type PartialApprovalData struct {
	TotalAmountRequested              Amount  `json:"totalAmountRequested"`
	PurchaseOnlyPartialAmountApproved Amount  `json:"purchaseOnlyPartialAmountApproved"`
	CashbackOnlyPartialAmountApproved *Amount `json:"cashbackOnlyPartialAmountApproved,omitempty"`
}

// ISO8583MessageData from OpenAPI spec.
type ISO8583MessageData struct {
	Mti   string `json:"mti,omitempty"`
	De002 string `json:"de002,omitempty"`
	De003 string `json:"de003,omitempty"`
	De004 string `json:"de004,omitempty"`
	De005 string `json:"de005,omitempty"`
	De006 string `json:"de006,omitempty"`
	De007 string `json:"de007,omitempty"`
	De008 string `json:"de008,omitempty"`
	De009 string `json:"de009,omitempty"`
	De010 string `json:"de010,omitempty"`
	De011 string `json:"de011,omitempty"`
	De012 string `json:"de012,omitempty"`
	De013 string `json:"de013,omitempty"`
	De014 string `json:"de014,omitempty"`
	De015 string `json:"de015,omitempty"`
	De016 string `json:"de016,omitempty"`
	De017 string `json:"de017,omitempty"`
	De018 string `json:"de018,omitempty"`
	De019 string `json:"de019,omitempty"`
	De020 string `json:"de020,omitempty"`
	De021 string `json:"de021,omitempty"`
	De022 string `json:"de022,omitempty"`
	De023 string `json:"de023,omitempty"`
	De024 string `json:"de024,omitempty"`
	De025 string `json:"de025,omitempty"`
	De026 string `json:"de026,omitempty"`
	De027 string `json:"de027,omitempty"`
	De028 string `json:"de028,omitempty"`
	De029 string `json:"de029,omitempty"`
	De030 string `json:"de030,omitempty"`
	De031 string `json:"de031,omitempty"`
	De032 string `json:"de032,omitempty"`
	De033 string `json:"de033,omitempty"`
	De034 string `json:"de034,omitempty"`
	De035 string `json:"de035,omitempty"`
	De036 string `json:"de036,omitempty"`
	De037 string `json:"de037,omitempty"`
	De038 string `json:"de038,omitempty"`
	De039 string `json:"de039,omitempty"`
	De040 string `json:"de040,omitempty"`
	De041 string `json:"de041,omitempty"`
	De042 string `json:"de042,omitempty"`
	De043 string `json:"de043,omitempty"`
	De044 string `json:"de044,omitempty"`
	De045 string `json:"de045,omitempty"`
	De046 string `json:"de046,omitempty"`
	De047 string `json:"de047,omitempty"`
	De048 string `json:"de048,omitempty"`
	De049 string `json:"de049,omitempty"`
	De050 string `json:"de050,omitempty"`
	De051 string `json:"de051,omitempty"`
	De052 string `json:"de052,omitempty"`
	De053 string `json:"de053,omitempty"`
	De054 string `json:"de054,omitempty"`
	De055 string `json:"de055,omitempty"`
	De056 string `json:"de056,omitempty"`
	De057 string `json:"de057,omitempty"`
	De058 string `json:"de058,omitempty"`
	De059 string `json:"de059,omitempty"`
	De060 string `json:"de060,omitempty"`
	De061 string `json:"de061,omitempty"`
	De062 string `json:"de062,omitempty"`
	De063 string `json:"de063,omitempty"`
	De064 string `json:"de064,omitempty"`
	De065 string `json:"de065,omitempty"`
	De066 string `json:"de066,omitempty"`
	De067 string `json:"de067,omitempty"`
	De068 string `json:"de068,omitempty"`
	De069 string `json:"de069,omitempty"`
	De070 string `json:"de070,omitempty"`
	De071 string `json:"de071,omitempty"`
	De072 string `json:"de072,omitempty"`
	De073 string `json:"de073,omitempty"`
	De074 string `json:"de074,omitempty"`
	De075 string `json:"de075,omitempty"`
	De076 string `json:"de076,omitempty"`
	De077 string `json:"de077,omitempty"`
	De078 string `json:"de078,omitempty"`
	De079 string `json:"de079,omitempty"`
	De080 string `json:"de080,omitempty"`
	De081 string `json:"de081,omitempty"`
	De082 string `json:"de082,omitempty"`
	De083 string `json:"de083,omitempty"`
	De084 string `json:"de084,omitempty"`
	De085 string `json:"de085,omitempty"`
	De086 string `json:"de086,omitempty"`
	De087 string `json:"de087,omitempty"`
	De088 string `json:"de088,omitempty"`
	De089 string `json:"de089,omitempty"`
	De090 string `json:"de090,omitempty"`
	De091 string `json:"de091,omitempty"`
	De092 string `json:"de092,omitempty"`
	De093 string `json:"de093,omitempty"`
	De094 string `json:"de094,omitempty"`
	De095 string `json:"de095,omitempty"`
	De096 string `json:"de096,omitempty"`
	De097 string `json:"de097,omitempty"`
	De098 string `json:"de098,omitempty"`
	De099 string `json:"de099,omitempty"`
	De100 string `json:"de100,omitempty"`
	De101 string `json:"de101,omitempty"`
	De102 string `json:"de102,omitempty"`
	De103 string `json:"de103,omitempty"`
	De104 string `json:"de104,omitempty"`
	De105 string `json:"de105,omitempty"`
	De106 string `json:"de106,omitempty"`
	De107 string `json:"de107,omitempty"`
	De108 string `json:"de108,omitempty"`
	De109 string `json:"de109,omitempty"`
	De110 string `json:"de110,omitempty"`
	De111 string `json:"de111,omitempty"`
	De112 string `json:"de112,omitempty"`
	De113 string `json:"de113,omitempty"`
	De114 string `json:"de114,omitempty"`
	De115 string `json:"de115,omitempty"`
	De116 string `json:"de116,omitempty"`
	De117 string `json:"de117,omitempty"`
	De118 string `json:"de118,omitempty"`
	De119 string `json:"de119,omitempty"`
	De120 string `json:"de120,omitempty"`
	De121 string `json:"de121,omitempty"`
	De122 string `json:"de122,omitempty"`
	De123 string `json:"de123,omitempty"`
	De124 string `json:"de124,omitempty"`
	De125 string `json:"de125,omitempty"`
	De126 string `json:"de126,omitempty"`
	De127 string `json:"de127,omitempty"`
}
