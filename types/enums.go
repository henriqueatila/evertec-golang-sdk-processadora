// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// AccountStatus enum values from OpenAPI spec.
type AccountStatus string

const (
	AccountStatusActive        AccountStatus = "active"
	AccountStatusCanceled      AccountStatus = "canceled"
	AccountStatusClosed        AccountStatus = "closed"
	AccountStatusBlocked       AccountStatus = "blocked"
	AccountStatusEnquadrada    AccountStatus = "enquadrada"
	AccountStatusRequested     AccountStatus = "requested"
	AccountStatusRequestDenied AccountStatus = "request_denied"
)

// TransactionStatus enum values from OpenAPI spec.
type TransactionStatus string

const (
	TransactionStatusReceived          TransactionStatus = "received"
	TransactionStatusPOSted            TransactionStatus = "posted"
	TransactionStatusInDispute         TransactionStatus = "in_dispute"
	TransactionStatusCanceled          TransactionStatus = "canceled"
	TransactionStatusRefunded          TransactionStatus = "refunded"
	TransactionStatusPartiallyRefunded TransactionStatus = "partially_refunded"
)

// TransactionType enum values from OpenAPI spec.
type TransactionType string

const (
	TransactionTypeSale              TransactionType = "sale"
	TransactionTypeWithdrawal        TransactionType = "withdrawal"
	TransactionTypeDebitAdjustment   TransactionType = "debit_adjustment"
	TransactionTypeUndoingLoadCredit TransactionType = "undoing_load_credit"
	TransactionTypeLoadingCredit     TransactionType = "loading_credit"
	TransactionTypeCreditAdjustment  TransactionType = "credit_adjustment"
	TransactionTypeSaleChargeback    TransactionType = "sale_chargeback"
	TransactionTypeInterest          TransactionType = "interest"
	TransactionTypeFine              TransactionType = "fine"
	TransactionTypePayment           TransactionType = "payment"
	TransactionTypeFee               TransactionType = "fee"
)

// TransactionStatusQuery enum values from OpenAPI spec.
type TransactionStatusQuery string

const (
	TransactionStatusQueryReceived          TransactionStatusQuery = "received"
	TransactionStatusQueryPOSted            TransactionStatusQuery = "posted"
	TransactionStatusQueryCanceled          TransactionStatusQuery = "canceled"
	TransactionStatusQueryRefunded          TransactionStatusQuery = "refunded"
	TransactionStatusQueryPartiallyRefunded TransactionStatusQuery = "partially_refunded"
)

// TransactionTypeQuery enum values from OpenAPI spec.
type TransactionTypeQuery string

const (
	TransactionTypeQuerySale             TransactionTypeQuery = "sale"
	TransactionTypeQueryWithdrawal       TransactionTypeQuery = "withdrawal"
	TransactionTypeQueryPreAuthorization TransactionTypeQuery = "pre_authorization"
	TransactionTypeQuerySaleChargeback   TransactionTypeQuery = "sale_chargeback"
	TransactionTypeQueryInterest         TransactionTypeQuery = "interest"
	TransactionTypeQueryFine             TransactionTypeQuery = "fine"
	TransactionTypeQueryPayment          TransactionTypeQuery = "payment"
	TransactionTypeQueryFee              TransactionTypeQuery = "fee"
)

// DisputeStatus enum values from OpenAPI spec.
type DisputeStatus string

const (
	DisputeStatusWaitingForResponseFromIssuer   DisputeStatus = "waiting_for_response_from_issuer"
	DisputeStatusWaitingForResponseFromAcquirer DisputeStatus = "waiting_for_response_from_acquirer"
	DisputeStatusWaitingForResponseFromBrand    DisputeStatus = "waiting_for_response_from_brand"
	DisputeStatusAccepted                       DisputeStatus = "accepted"
	DisputeStatusDenied                         DisputeStatus = "denied"
	DisputeStatusIssuerWon                      DisputeStatus = "issuer_won"
	DisputeStatusIssuerLost                     DisputeStatus = "issuer_lost"
	DisputeStatusCanceledByIssuer               DisputeStatus = "canceled_by_issuer"
)

// DisputeType enum values from OpenAPI spec.
type DisputeType string

const (
	DisputeTypeServiceNotOfferedOrProductNotDelivered     DisputeType = "service_not_offered_or_product_not_delivered"
	DisputeTypeRecurringTransactionCancelled              DisputeType = "recurring_transaction_cancelled"
	DisputeTypeProductDefectiveOrDifferingFromDescription DisputeType = "product_defective_or_differing_from_description"
	DisputeTypeMultipleFraudulentTransactions             DisputeType = "multiple_fraudulent_transactions"
	DisputeTypeIllegibleDocument                          DisputeType = "illegible_document"
	DisputeTypeChipResponsabilityTransference             DisputeType = "chip_responsability_transference"
	DisputeTypeAuthorizationDenied                        DisputeType = "authorization_denied"
	DisputeTypeNoAuthorization                            DisputeType = "no_authorization"
	DisputeTypeExpiredCard                                DisputeType = "expired_card"
	DisputeTypeLatePresentation                           DisputeType = "late_presentation"
	DisputeTypeHolderDoesNoRecallTransaction              DisputeType = "holder_does_no_recall_transaction"
	DisputeTypeNonExistingCardNumber                      DisputeType = "non_existing_card_number"
	DisputeTypeIncorrectTransactionValue                  DisputeType = "incorrect_transaction_value"
	DisputeTypeCardPresentFraud                           DisputeType = "card_present_fraud"
	DisputeTypeDuplicatedProcessing                       DisputeType = "duplicated_processing"
	DisputeTypeCardNotPresentFraud                        DisputeType = "card_not_present_fraud"
	DisputeTypeCreditNotProcessed                         DisputeType = "credit_not_processed"
	DisputeTypePaymentByOtherMeans                        DisputeType = "payment_by_other_means"
)

// EntryMode enum values from OpenAPI spec.
type EntryMode string

const (
	EntryModeEmvcontact     EntryMode = "emvContact"
	EntryModeEmvcontactless EntryMode = "emvContactless"
	EntryModeHce            EntryMode = "hce"
	EntryModeQrcode         EntryMode = "qrCode"
	EntryModeEcommerce      EntryMode = "ecommerce"
	EntryModeMagstripe      EntryMode = "magstripe"
	EntryModeManual         EntryMode = "manual"
)

// FraudNotificationStatus enum values from OpenAPI spec.
type FraudNotificationStatus string

const (
	FraudNotificationStatusNotified FraudNotificationStatus = "notified"
	FraudNotificationStatusUndone   FraudNotificationStatus = "undone"
)

// HealthStatus enum values from OpenAPI spec.
type HealthStatus string

const (
	HealthStatusHealthy    HealthStatus = "healthy"
	HealthStatusUnhealthy  HealthStatus = "unhealthy"
	HealthStatusUnstable   HealthStatus = "unstable"
	HealthStatusRecovering HealthStatus = "recovering"
	HealthStatusNoData     HealthStatus = "no_data"
)

// PaymentType enum values from OpenAPI spec.
type PaymentType string

const (
	PaymentTypeA2a PaymentType = "A2A"
	PaymentTypeB2a PaymentType = "B2A"
	PaymentTypeP2p PaymentType = "P2P"
	PaymentTypeCsb PaymentType = "CSB"
	PaymentTypeDsb PaymentType = "DSB"
	PaymentTypeB2b PaymentType = "B2B"
	PaymentTypeM2m PaymentType = "M2M"
)

// Source enum values from OpenAPI spec.
type Source string

const (
	SourcePaysmart Source = "paySmart"
	SourceIssuer   Source = "issuer"
	SourceBrand    Source = "brand"
)

// CardStatus enum values from OpenAPI spec.
type CardStatus string

const (
	CardStatusBlocked       CardStatus = "blocked"
	CardStatusActive        CardStatus = "active"
	CardStatusCancelled     CardStatus = "cancelled"
	CardStatusIssuing       CardStatus = "issuing"
	CardStatusIntransit     CardStatus = "inTransit"
	CardStatusRequested     CardStatus = "requested"
	CardStatusRequestDenied CardStatus = "request_denied"
	CardStatusPurged        CardStatus = "purged"
)

// CardProfile enum values.
type CardProfile string

const (
	CardProfileCredit  CardProfile = "credit"
	CardProfileDebit   CardProfile = "debit"
	CardProfileVoucher CardProfile = "voucher"
	CardProfileFleet   CardProfile = "fleet"
	CardProfileCombo   CardProfile = "combo"
)

// CardholderType enum values.
type CardholderType string

const (
	CardholderTypeMain       CardholderType = "main"
	CardholderTypeAdditional CardholderType = "additional"
)

// Gender enum values.
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// CivilStatus enum values.
type CivilStatus string

const (
	CivilStatusSingle   CivilStatus = "single"
	CivilStatusMarried  CivilStatus = "married"
	CivilStatusDivorced CivilStatus = "divorced"
	CivilStatusWidowed  CivilStatus = "widowed"
	CivilStatusOther    CivilStatus = "other"
)

// DocumentType enum values.
type DocumentType string

const (
	DocumentTypeCpf  DocumentType = "cpf"
	DocumentTypeCnpj DocumentType = "cnpj"
)

// FeeType enum values.
type FeeType string

const (
	FeeTypeIof           FeeType = "iof"
	FeeTypeMarkup        FeeType = "markup"
	FeeTypeBoardingFee   FeeType = "boarding_fee"
	FeeTypeWithdrawalFee FeeType = "withdrawal_fee"
	FeeTypeInterest      FeeType = "interest"
	FeeTypeOthers        FeeType = "others"
)

// DebitOrCredit enum values.
type DebitOrCredit string

const (
	DebitOrCreditDebit  DebitOrCredit = "debit"
	DebitOrCreditCredit DebitOrCredit = "credit"
)

// InstallmentFinType enum values.
type InstallmentFinType string

const (
	InstallmentFinTypeInterestfree InstallmentFinType = "interestFree"
	InstallmentFinTypeWithinterest InstallmentFinType = "withInterest"
	InstallmentFinTypeCdc          InstallmentFinType = "cdc"
)

// ProductType enum values.
type ProductType string

const (
	ProductTypePre  ProductType = "pre"
	ProductTypePOSt ProductType = "post"
)

// DisputeStage enum values.
type DisputeStage string

const (
	DisputeStageChargeback     DisputeStage = "chargeback"
	DisputeStageRestatement    DisputeStage = "restatement"
	DisputeStagePreArbitration DisputeStage = "pre-arbitration"
	DisputeStageArbitration    DisputeStage = "arbitration"
	DisputeStageLost           DisputeStage = "lost"
	DisputeStageUndone         DisputeStage = "undone"
	DisputeStageSucessful      DisputeStage = "sucessful"
	DisputeStageWithdrawal     DisputeStage = "withdrawal"
)

// InclusiveTransactionCode enum values.
type InclusiveTransactionCode string

const (
	InclusiveTransactionCodeTe10 InclusiveTransactionCode = "TE10"
	InclusiveTransactionCodeTe20 InclusiveTransactionCode = "TE20"
)

// CreditRequestType enum values.
type CreditRequestType string

const (
	CreditRequestTypePayment    CreditRequestType = "payment"
	CreditRequestTypeAdjustment CreditRequestType = "adjustment"
)

// DebitRequestType enum values.
type DebitRequestType string

const (
	DebitRequestTypeAnnuityfee DebitRequestType = "annuityFee"
	DebitRequestTypeIssuefee   DebitRequestType = "issueFee"
	DebitRequestTypeReissuefee DebitRequestType = "reissueFee"
	DebitRequestTypeRenewalfee DebitRequestType = "renewalFee"
	DebitRequestTypeAdjustment DebitRequestType = "adjustment"
)

// PinFormat enum values.
type PinFormat string

const (
	PinFormatIso0 PinFormat = "ISO-0"
	PinFormatIso1 PinFormat = "ISO-1"
)

// FunctionalityCode enum values.
type FunctionalityCode string

const (
	FunctionalityCodeContactless    FunctionalityCode = "contactless"
	FunctionalityCodeMagneticStripe FunctionalityCode = "magnetic_stripe"
	FunctionalityCodeECommerce      FunctionalityCode = "e_commerce"
	FunctionalityCodeWithdrawal     FunctionalityCode = "withdrawal"
)

// CampaignStatus enum values.
type CampaignStatus string

const (
	CampaignStatusEnabled  CampaignStatus = "enabled"
	CampaignStatusDisabled CampaignStatus = "disabled"
)

// VirtualCardStatus enum values.
type VirtualCardStatus string

const (
	VirtualCardStatusActive   VirtualCardStatus = "active"
	VirtualCardStatusInactive VirtualCardStatus = "inactive"
	VirtualCardStatusExpired  VirtualCardStatus = "expired"
)

// CreditMattressStatus enum values.
type CreditMattressStatus string

const (
	CreditMattressStatusConfirmed CreditMattressStatus = "confirmed"
	CreditMattressStatusCanceled  CreditMattressStatus = "canceled"
	CreditMattressStatusPending   CreditMattressStatus = "pending"
)

// courierId enum values.
type CourierID string

const (
	CourierID000 CourierID = "000"
	CourierID001 CourierID = "001"
	CourierID002 CourierID = "002"
	CourierID003 CourierID = "003"
	CourierID004 CourierID = "004"
	CourierID005 CourierID = "005"
	CourierID006 CourierID = "006"
)

// bureauxId enum values.
type BureauxID string

const (
	BureauxID0001 BureauxID = "0001"
	BureauxID0002 BureauxID = "0002"
	BureauxID0003 BureauxID = "0003"
	BureauxID0004 BureauxID = "0004"
	BureauxID0005 BureauxID = "0005"
	BureauxID0006 BureauxID = "0006"
	BureauxID0007 BureauxID = "0007"
	BureauxID0008 BureauxID = "0008"
	BureauxID0009 BureauxID = "0009"
)

// disputeCode type (string).
type DisputeCode string

// ============================================================================
// ACCOUNT BLOCK/UNBLOCK/CANCELLATION CODES
// Reference: https://paysmart-api.gitlab.io/processadora/PT-br/docs/gerenciamento-contas
// ============================================================================

// AccountBlockCode represents valid block codes for accounts.
// POST /accounts/{accountId}/blockAccount
type AccountBlockCode int

const (
	// Account block codes from official API documentation
	AccountBlockCodeIssuerRequest       AccountBlockCode = 10 // Bloqueio por solicitação do emissor
	AccountBlockCodeCustomerRequest     AccountBlockCode = 11 // Bloqueio por solicitação do cliente
	AccountBlockCodeFraudSuspect        AccountBlockCode = 12 // Bloqueio por suspeita de fraude
	AccountBlockCodeAddressReturn       AccountBlockCode = 16 // Bloqueio por devolução de endereço
	AccountBlockCodeDelinquency         AccountBlockCode = 18 // Bloqueio por inadimplência
	AccountBlockCodeProtectiveService   AccountBlockCode = 20 // Bloqueio por serviço de proteção
	AccountBlockCodeCreditBureau        AccountBlockCode = 21 // Bloqueio por informação do bureau de crédito
	AccountBlockCodeJudicialOrder       AccountBlockCode = 22 // Bloqueio por ordem judicial
	AccountBlockCodeTemporaryInactivity AccountBlockCode = 23 // Bloqueio por inatividade temporária
	AccountBlockCodeExcessiveDeclines   AccountBlockCode = 24 // Bloqueio por excesso de negativas
	AccountBlockCodeLost                AccountBlockCode = 30 // Bloqueio por perda
	AccountBlockCodeStolen              AccountBlockCode = 70 // Bloqueio por roubo
	AccountBlockCodeIssuerDecision      AccountBlockCode = 80 // Bloqueio por decisão do emissor
	AccountBlockCodeIssuerDecision2     AccountBlockCode = 81 // Bloqueio por decisão do emissor (2)
	AccountBlockCodeFraudConfirmed      AccountBlockCode = 82 // Bloqueio por fraude confirmada
	AccountBlockCodeSystemBlock         AccountBlockCode = 90 // Bloqueio do sistema
	AccountBlockCodeSystemBlock2        AccountBlockCode = 91 // Bloqueio do sistema (2)
	AccountBlockCodeRegulatory          AccountBlockCode = 94 // Bloqueio regulatório
	AccountBlockCodeCompliance          AccountBlockCode = 95 // Bloqueio por compliance
)

// AccountUnblockCode represents valid unblock codes for accounts.
// POST /accounts/{accountId}/unblockAccount
type AccountUnblockCode int

const (
	AccountUnblockCodeStandard AccountUnblockCode = 0 // Desbloqueio padrão
)

// AccountCancellationCode represents valid cancellation codes for accounts.
// POST /accounts/{accountId}/cancelAccount
type AccountCancellationCode int

const (
	AccountCancellationCodeIssuerDecision  AccountCancellationCode = 80 // Cancelamento por decisão do emissor
	AccountCancellationCodeCustomerRequest AccountCancellationCode = 81 // Cancelamento por solicitação do cliente
)

// CardBlockCode represents valid block codes for cards.
// POST /cards/{cardId}/blockCardRequest
type CardBlockCode int

const (
	// Card block codes mirror account block codes where applicable
	CardBlockCodeIssuerRequest   CardBlockCode = 10 // Bloqueio por solicitação do emissor
	CardBlockCodeCustomerRequest CardBlockCode = 11 // Bloqueio por solicitação do cliente
	CardBlockCodeFraudSuspect    CardBlockCode = 12 // Bloqueio por suspeita de fraude
	CardBlockCodeDelinquency     CardBlockCode = 18 // Bloqueio por inadimplência
	CardBlockCodeLost            CardBlockCode = 30 // Bloqueio por perda
	CardBlockCodeStolen          CardBlockCode = 70 // Bloqueio por roubo
	CardBlockCodeFraudConfirmed  CardBlockCode = 82 // Bloqueio por fraude confirmada
)

// CardUnblockCode represents valid unblock codes for cards.
// POST /cards/{cardId}/unblockCardRequest
type CardUnblockCode int

const (
	CardUnblockCodeStandard CardUnblockCode = 0 // Desbloqueio padrão
)
