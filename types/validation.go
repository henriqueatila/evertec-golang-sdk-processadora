package types

import "errors"

// ErrInvalidAmount is returned when an amount value is invalid.
var ErrInvalidAmount = errors.New("amount cannot be negative")

// ErrInvalidCurrencyCode is returned when a currency code is invalid.
var ErrInvalidCurrencyCode = errors.New("currency code must be positive")

// Validate validates the Amount fields.
func (a Amount) Validate() error {
	if a.Amount < 0 {
		return ErrInvalidAmount
	}
	if a.CurrencyCode <= 0 {
		return ErrInvalidCurrencyCode
	}
	return nil
}

// Validate validates the AmountTransaction fields.
func (a AmountTransaction) Validate() error {
	if a.Amount < 0 {
		return ErrInvalidAmount
	}
	if a.CurrencyCode <= 0 {
		return ErrInvalidCurrencyCode
	}
	return nil
}

// IsValid returns true if the AccountStatus is a valid enum value.
func (s AccountStatus) IsValid() bool {
	switch s {
	case AccountStatusActive, AccountStatusCanceled, AccountStatusClosed,
		AccountStatusBlocked, AccountStatusEnquadrada, AccountStatusRequested,
		AccountStatusRequestDenied:
		return true
	}
	return false
}

// IsValid returns true if the TransactionStatus is a valid enum value.
func (s TransactionStatus) IsValid() bool {
	switch s {
	case TransactionStatusReceived, TransactionStatusPOSted, TransactionStatusInDispute,
		TransactionStatusCanceled, TransactionStatusRefunded, TransactionStatusPartiallyRefunded:
		return true
	}
	return false
}

// IsValid returns true if the TransactionType is a valid enum value.
func (t TransactionType) IsValid() bool {
	switch t {
	case TransactionTypeSale, TransactionTypeWithdrawal, TransactionTypeDebitAdjustment,
		TransactionTypeUndoingLoadCredit, TransactionTypeLoadingCredit, TransactionTypeCreditAdjustment,
		TransactionTypeSaleChargeback, TransactionTypeInterest, TransactionTypeFine,
		TransactionTypePayment, TransactionTypeFee:
		return true
	}
	return false
}

// IsValid returns true if the TransactionStatusQuery is a valid enum value.
func (s TransactionStatusQuery) IsValid() bool {
	switch s {
	case TransactionStatusQueryReceived, TransactionStatusQueryPOSted,
		TransactionStatusQueryCanceled, TransactionStatusQueryRefunded,
		TransactionStatusQueryPartiallyRefunded:
		return true
	}
	return false
}

// IsValid returns true if the TransactionTypeQuery is a valid enum value.
func (t TransactionTypeQuery) IsValid() bool {
	switch t {
	case TransactionTypeQuerySale, TransactionTypeQueryWithdrawal,
		TransactionTypeQueryPreAuthorization, TransactionTypeQuerySaleChargeback,
		TransactionTypeQueryInterest, TransactionTypeQueryFine,
		TransactionTypeQueryPayment, TransactionTypeQueryFee:
		return true
	}
	return false
}

// IsValid returns true if the DisputeStatus is a valid enum value.
func (s DisputeStatus) IsValid() bool {
	switch s {
	case DisputeStatusWaitingForResponseFromIssuer, DisputeStatusWaitingForResponseFromAcquirer,
		DisputeStatusWaitingForResponseFromBrand, DisputeStatusAccepted, DisputeStatusDenied,
		DisputeStatusIssuerWon, DisputeStatusIssuerLost, DisputeStatusCanceledByIssuer:
		return true
	}
	return false
}

// IsValid returns true if the DisputeType is a valid enum value.
func (t DisputeType) IsValid() bool {
	switch t {
	case DisputeTypeServiceNotOfferedOrProductNotDelivered, DisputeTypeRecurringTransactionCancelled,
		DisputeTypeProductDefectiveOrDifferingFromDescription, DisputeTypeMultipleFraudulentTransactions,
		DisputeTypeIllegibleDocument, DisputeTypeChipResponsabilityTransference,
		DisputeTypeAuthorizationDenied, DisputeTypeNoAuthorization, DisputeTypeExpiredCard,
		DisputeTypeLatePresentation, DisputeTypeHolderDoesNoRecallTransaction,
		DisputeTypeNonExistingCardNumber, DisputeTypeIncorrectTransactionValue,
		DisputeTypeCardPresentFraud, DisputeTypeDuplicatedProcessing,
		DisputeTypeCardNotPresentFraud, DisputeTypeCreditNotProcessed, DisputeTypePaymentByOtherMeans:
		return true
	}
	return false
}

// IsValid returns true if the EntryMode is a valid enum value.
func (e EntryMode) IsValid() bool {
	switch e {
	case EntryModeEmvcontact, EntryModeEmvcontactless, EntryModeHce,
		EntryModeQrcode, EntryModeEcommerce, EntryModeMagstripe, EntryModeManual:
		return true
	}
	return false
}

// IsValid returns true if the FraudNotificationStatus is a valid enum value.
func (s FraudNotificationStatus) IsValid() bool {
	switch s {
	case FraudNotificationStatusNotified, FraudNotificationStatusUndone:
		return true
	}
	return false
}

// IsValid returns true if the HealthStatus is a valid enum value.
func (s HealthStatus) IsValid() bool {
	switch s {
	case HealthStatusHealthy, HealthStatusUnhealthy, HealthStatusUnstable,
		HealthStatusRecovering, HealthStatusNoData:
		return true
	}
	return false
}

// IsValid returns true if the PaymentType is a valid enum value.
func (p PaymentType) IsValid() bool {
	switch p {
	case PaymentTypeA2a, PaymentTypeB2a, PaymentTypeP2p,
		PaymentTypeCsb, PaymentTypeDsb, PaymentTypeB2b, PaymentTypeM2m:
		return true
	}
	return false
}

// IsValid returns true if the Source is a valid enum value.
func (s Source) IsValid() bool {
	switch s {
	case SourcePaysmart, SourceIssuer, SourceBrand:
		return true
	}
	return false
}

// IsValid returns true if the CardStatus is a valid enum value.
func (s CardStatus) IsValid() bool {
	switch s {
	case CardStatusBlocked, CardStatusActive, CardStatusCancelled,
		CardStatusIssuing, CardStatusIntransit, CardStatusRequested,
		CardStatusRequestDenied, CardStatusPurged:
		return true
	}
	return false
}

// IsValid returns true if the CardProfile is a valid enum value.
func (p CardProfile) IsValid() bool {
	switch p {
	case CardProfileCredit, CardProfileDebit, CardProfileVoucher,
		CardProfileFleet, CardProfileCombo:
		return true
	}
	return false
}

// IsValid returns true if the CardholderType is a valid enum value.
func (t CardholderType) IsValid() bool {
	switch t {
	case CardholderTypeMain, CardholderTypeAdditional:
		return true
	}
	return false
}

// IsValid returns true if the Gender is a valid enum value.
func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderOther:
		return true
	}
	return false
}

// IsValid returns true if the CivilStatus is a valid enum value.
func (c CivilStatus) IsValid() bool {
	switch c {
	case CivilStatusSingle, CivilStatusMarried, CivilStatusDivorced,
		CivilStatusWidowed, CivilStatusOther:
		return true
	}
	return false
}

// IsValid returns true if the DocumentType is a valid enum value.
func (d DocumentType) IsValid() bool {
	switch d {
	case DocumentTypeCpf, DocumentTypeCnpj:
		return true
	}
	return false
}

// IsValid returns true if the FeeType is a valid enum value.
func (f FeeType) IsValid() bool {
	switch f {
	case FeeTypeIof, FeeTypeMarkup, FeeTypeBoardingFee,
		FeeTypeWithdrawalFee, FeeTypeInterest, FeeTypeOthers:
		return true
	}
	return false
}

// IsValid returns true if the DebitOrCredit is a valid enum value.
func (d DebitOrCredit) IsValid() bool {
	switch d {
	case DebitOrCreditDebit, DebitOrCreditCredit:
		return true
	}
	return false
}

// IsValid returns true if the InstallmentFinType is a valid enum value.
func (i InstallmentFinType) IsValid() bool {
	switch i {
	case InstallmentFinTypeInterestfree, InstallmentFinTypeWithinterest, InstallmentFinTypeCdc:
		return true
	}
	return false
}

// IsValid returns true if the ProductType is a valid enum value.
func (p ProductType) IsValid() bool {
	switch p {
	case ProductTypePre, ProductTypePOSt:
		return true
	}
	return false
}

// IsValid returns true if the DisputeStage is a valid enum value.
func (d DisputeStage) IsValid() bool {
	switch d {
	case DisputeStageChargeback, DisputeStageRestatement, DisputeStagePreArbitration,
		DisputeStageArbitration, DisputeStageLost, DisputeStageUndone,
		DisputeStageSucessful, DisputeStageWithdrawal:
		return true
	}
	return false
}

// IsValid returns true if the InclusiveTransactionCode is a valid enum value.
func (i InclusiveTransactionCode) IsValid() bool {
	switch i {
	case InclusiveTransactionCodeTe10, InclusiveTransactionCodeTe20:
		return true
	}
	return false
}

// IsValid returns true if the CreditRequestType is a valid enum value.
func (c CreditRequestType) IsValid() bool {
	switch c {
	case CreditRequestTypePayment, CreditRequestTypeAdjustment:
		return true
	}
	return false
}

// IsValid returns true if the DebitRequestType is a valid enum value.
func (d DebitRequestType) IsValid() bool {
	switch d {
	case DebitRequestTypeAnnuityfee, DebitRequestTypeIssuefee,
		DebitRequestTypeReissuefee, DebitRequestTypeRenewalfee, DebitRequestTypeAdjustment:
		return true
	}
	return false
}

// IsValid returns true if the PinFormat is a valid enum value.
func (p PinFormat) IsValid() bool {
	switch p {
	case PinFormatIso0, PinFormatIso1:
		return true
	}
	return false
}

// IsValid returns true if the FunctionalityCode is a valid enum value.
func (f FunctionalityCode) IsValid() bool {
	switch f {
	case FunctionalityCodeContactless, FunctionalityCodeMagneticStripe,
		FunctionalityCodeECommerce, FunctionalityCodeWithdrawal:
		return true
	}
	return false
}

// IsValid returns true if the CampaignStatus is a valid enum value.
func (c CampaignStatus) IsValid() bool {
	switch c {
	case CampaignStatusEnabled, CampaignStatusDisabled:
		return true
	}
	return false
}

// IsValid returns true if the VirtualCardStatus is a valid enum value.
func (v VirtualCardStatus) IsValid() bool {
	switch v {
	case VirtualCardStatusActive, VirtualCardStatusInactive, VirtualCardStatusExpired:
		return true
	}
	return false
}

// IsValid returns true if the CreditMattressStatus is a valid enum value.
func (c CreditMattressStatus) IsValid() bool {
	switch c {
	case CreditMattressStatusConfirmed, CreditMattressStatusCanceled, CreditMattressStatusPending:
		return true
	}
	return false
}

// IsValid returns true if the CourierID is a valid enum value.
func (c CourierID) IsValid() bool {
	switch c {
	case CourierID000, CourierID001, CourierID002, CourierID003,
		CourierID004, CourierID005, CourierID006:
		return true
	}
	return false
}

// IsValid returns true if the BureauxID is a valid enum value.
func (b BureauxID) IsValid() bool {
	switch b {
	case BureauxID0001, BureauxID0002, BureauxID0003, BureauxID0004,
		BureauxID0005, BureauxID0006, BureauxID0007, BureauxID0008, BureauxID0009:
		return true
	}
	return false
}
