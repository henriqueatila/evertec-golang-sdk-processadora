package types

import (
	"errors"
	"testing"
)

func TestAmountValidate(t *testing.T) {
	tests := []struct {
		name    string
		amount  Amount
		wantErr error
	}{
		{
			name:    "valid amount",
			amount:  Amount{Amount: 1000, CurrencyCode: 986},
			wantErr: nil,
		},
		{
			name:    "zero amount",
			amount:  Amount{Amount: 0, CurrencyCode: 986},
			wantErr: nil,
		},
		{
			name:    "negative amount",
			amount:  Amount{Amount: -100, CurrencyCode: 986},
			wantErr: ErrInvalidAmount,
		},
		{
			name:    "zero currency code",
			amount:  Amount{Amount: 1000, CurrencyCode: 0},
			wantErr: ErrInvalidCurrencyCode,
		},
		{
			name:    "negative currency code",
			amount:  Amount{Amount: 1000, CurrencyCode: -1},
			wantErr: ErrInvalidCurrencyCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.amount.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Amount.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAmountTransactionValidate(t *testing.T) {
	tests := []struct {
		name    string
		amount  AmountTransaction
		wantErr error
	}{
		{
			name:    "valid amount",
			amount:  AmountTransaction{Amount: 1000, CurrencyCode: 986},
			wantErr: nil,
		},
		{
			name:    "negative amount",
			amount:  AmountTransaction{Amount: -100, CurrencyCode: 986},
			wantErr: ErrInvalidAmount,
		},
		{
			name:    "zero currency code",
			amount:  AmountTransaction{Amount: 1000, CurrencyCode: 0},
			wantErr: ErrInvalidCurrencyCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.amount.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AmountTransaction.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccountStatusIsValid(t *testing.T) {
	validStatuses := []AccountStatus{
		AccountStatusActive,
		AccountStatusCanceled,
		AccountStatusClosed,
		AccountStatusBlocked,
		AccountStatusEnquadrada,
		AccountStatusRequested,
		AccountStatusRequestDenied,
	}

	for _, status := range validStatuses {
		if !status.IsValid() {
			t.Errorf("AccountStatus(%q).IsValid() = false, want true", status)
		}
	}

	invalidStatuses := []AccountStatus{
		"invalid",
		"",
		"ACTIVE",
		"unknown_status",
	}

	for _, status := range invalidStatuses {
		if status.IsValid() {
			t.Errorf("AccountStatus(%q).IsValid() = true, want false", status)
		}
	}
}

func TestTransactionStatusIsValid(t *testing.T) {
	validStatuses := []TransactionStatus{
		TransactionStatusReceived,
		TransactionStatusPOSted,
		TransactionStatusInDispute,
		TransactionStatusCanceled,
		TransactionStatusRefunded,
		TransactionStatusPartiallyRefunded,
	}

	for _, status := range validStatuses {
		if !status.IsValid() {
			t.Errorf("TransactionStatus(%q).IsValid() = false, want true", status)
		}
	}

	if TransactionStatus("invalid").IsValid() {
		t.Error("TransactionStatus(\"invalid\").IsValid() = true, want false")
	}
}

func TestCardStatusIsValid(t *testing.T) {
	validStatuses := []CardStatus{
		CardStatusBlocked,
		CardStatusActive,
		CardStatusCancelled,
		CardStatusIssuing,
		CardStatusIntransit,
		CardStatusRequested,
		CardStatusRequestDenied,
		CardStatusPurged,
	}

	for _, status := range validStatuses {
		if !status.IsValid() {
			t.Errorf("CardStatus(%q).IsValid() = false, want true", status)
		}
	}

	if CardStatus("invalid").IsValid() {
		t.Error("CardStatus(\"invalid\").IsValid() = true, want false")
	}
}

func TestDisputeStatusIsValid(t *testing.T) {
	validStatuses := []DisputeStatus{
		DisputeStatusWaitingForResponseFromIssuer,
		DisputeStatusWaitingForResponseFromAcquirer,
		DisputeStatusWaitingForResponseFromBrand,
		DisputeStatusAccepted,
		DisputeStatusDenied,
		DisputeStatusIssuerWon,
		DisputeStatusIssuerLost,
		DisputeStatusCanceledByIssuer,
	}

	for _, status := range validStatuses {
		if !status.IsValid() {
			t.Errorf("DisputeStatus(%q).IsValid() = false, want true", status)
		}
	}

	if DisputeStatus("invalid").IsValid() {
		t.Error("DisputeStatus(\"invalid\").IsValid() = true, want false")
	}
}

func TestEntryModeIsValid(t *testing.T) {
	validModes := []EntryMode{
		EntryModeEmvcontact,
		EntryModeEmvcontactless,
		EntryModeHce,
		EntryModeQrcode,
		EntryModeEcommerce,
		EntryModeMagstripe,
		EntryModeManual,
	}

	for _, mode := range validModes {
		if !mode.IsValid() {
			t.Errorf("EntryMode(%q).IsValid() = false, want true", mode)
		}
	}

	if EntryMode("invalid").IsValid() {
		t.Error("EntryMode(\"invalid\").IsValid() = true, want false")
	}
}

func TestGenderIsValid(t *testing.T) {
	validGenders := []Gender{
		GenderMale,
		GenderFemale,
		GenderOther,
	}

	for _, gender := range validGenders {
		if !gender.IsValid() {
			t.Errorf("Gender(%q).IsValid() = false, want true", gender)
		}
	}

	if Gender("invalid").IsValid() {
		t.Error("Gender(\"invalid\").IsValid() = true, want false")
	}
}

func TestDocumentTypeIsValid(t *testing.T) {
	validTypes := []DocumentType{
		DocumentTypeCpf,
		DocumentTypeCnpj,
	}

	for _, dt := range validTypes {
		if !dt.IsValid() {
			t.Errorf("DocumentType(%q).IsValid() = false, want true", dt)
		}
	}

	if DocumentType("rg").IsValid() {
		t.Error("DocumentType(\"rg\").IsValid() = true, want false")
	}
}

func TestPinFormatIsValid(t *testing.T) {
	validFormats := []PinFormat{
		PinFormatIso0,
		PinFormatIso1,
	}

	for _, format := range validFormats {
		if !format.IsValid() {
			t.Errorf("PinFormat(%q).IsValid() = false, want true", format)
		}
	}

	if PinFormat("ISO-2").IsValid() {
		t.Error("PinFormat(\"ISO-2\").IsValid() = true, want false")
	}
}

func TestCardProfileIsValid(t *testing.T) {
	validProfiles := []CardProfile{
		CardProfileCredit,
		CardProfileDebit,
		CardProfileVoucher,
		CardProfileFleet,
		CardProfileCombo,
	}

	for _, profile := range validProfiles {
		if !profile.IsValid() {
			t.Errorf("CardProfile(%q).IsValid() = false, want true", profile)
		}
	}

	if CardProfile("prepaid").IsValid() {
		t.Error("CardProfile(\"prepaid\").IsValid() = true, want false")
	}
}

func TestProductTypeIsValid(t *testing.T) {
	if !ProductTypePre.IsValid() {
		t.Error("ProductTypePre.IsValid() = false, want true")
	}
	if !ProductTypePOSt.IsValid() {
		t.Error("ProductTypePOSt.IsValid() = false, want true")
	}
	if ProductType("hybrid").IsValid() {
		t.Error("ProductType(\"hybrid\").IsValid() = true, want false")
	}
}

func TestCourierIDIsValid(t *testing.T) {
	validIDs := []CourierID{
		CourierID000,
		CourierID001,
		CourierID002,
		CourierID003,
		CourierID004,
		CourierID005,
		CourierID006,
	}

	for _, id := range validIDs {
		if !id.IsValid() {
			t.Errorf("CourierID(%q).IsValid() = false, want true", id)
		}
	}

	if CourierID("007").IsValid() {
		t.Error("CourierID(\"007\").IsValid() = true, want false")
	}
}

func TestBureauxIDIsValid(t *testing.T) {
	validIDs := []BureauxID{
		BureauxID0001,
		BureauxID0002,
		BureauxID0003,
		BureauxID0004,
		BureauxID0005,
		BureauxID0006,
		BureauxID0007,
		BureauxID0008,
		BureauxID0009,
	}

	for _, id := range validIDs {
		if !id.IsValid() {
			t.Errorf("BureauxID(%q).IsValid() = false, want true", id)
		}
	}

	if BureauxID("0010").IsValid() {
		t.Error("BureauxID(\"0010\").IsValid() = true, want false")
	}
}

func TestTransactionTypeIsValid(t *testing.T) {
	validTypes := []TransactionType{
		TransactionTypeSale,
		TransactionTypeWithdrawal,
		TransactionTypeDebitAdjustment,
		TransactionTypeUndoingLoadCredit,
		TransactionTypeLoadingCredit,
		TransactionTypeCreditAdjustment,
		TransactionTypeSaleChargeback,
		TransactionTypeInterest,
		TransactionTypeFine,
		TransactionTypePayment,
		TransactionTypeFee,
	}

	for _, tt := range validTypes {
		if !tt.IsValid() {
			t.Errorf("TransactionType(%q).IsValid() = false, want true", tt)
		}
	}

	if TransactionType("invalid").IsValid() {
		t.Error("TransactionType(\"invalid\").IsValid() = true, want false")
	}
}

func TestTransactionStatusQueryIsValid(t *testing.T) {
	validStatuses := []TransactionStatusQuery{
		TransactionStatusQueryReceived,
		TransactionStatusQueryPOSted,
		TransactionStatusQueryCanceled,
		TransactionStatusQueryRefunded,
		TransactionStatusQueryPartiallyRefunded,
	}

	for _, status := range validStatuses {
		if !status.IsValid() {
			t.Errorf("TransactionStatusQuery(%q).IsValid() = false, want true", status)
		}
	}

	if TransactionStatusQuery("invalid").IsValid() {
		t.Error("TransactionStatusQuery(\"invalid\").IsValid() = true, want false")
	}
}

func TestTransactionTypeQueryIsValid(t *testing.T) {
	validTypes := []TransactionTypeQuery{
		TransactionTypeQuerySale,
		TransactionTypeQueryWithdrawal,
		TransactionTypeQueryPreAuthorization,
		TransactionTypeQuerySaleChargeback,
		TransactionTypeQueryInterest,
		TransactionTypeQueryFine,
		TransactionTypeQueryPayment,
		TransactionTypeQueryFee,
	}

	for _, tt := range validTypes {
		if !tt.IsValid() {
			t.Errorf("TransactionTypeQuery(%q).IsValid() = false, want true", tt)
		}
	}

	if TransactionTypeQuery("invalid").IsValid() {
		t.Error("TransactionTypeQuery(\"invalid\").IsValid() = true, want false")
	}
}

func TestDisputeTypeIsValid(t *testing.T) {
	validTypes := []DisputeType{
		DisputeTypeServiceNotOfferedOrProductNotDelivered,
		DisputeTypeRecurringTransactionCancelled,
		DisputeTypeProductDefectiveOrDifferingFromDescription,
		DisputeTypeMultipleFraudulentTransactions,
		DisputeTypeIllegibleDocument,
		DisputeTypeChipResponsabilityTransference,
		DisputeTypeAuthorizationDenied,
		DisputeTypeNoAuthorization,
		DisputeTypeExpiredCard,
		DisputeTypeLatePresentation,
		DisputeTypeHolderDoesNoRecallTransaction,
		DisputeTypeNonExistingCardNumber,
		DisputeTypeIncorrectTransactionValue,
		DisputeTypeCardPresentFraud,
		DisputeTypeDuplicatedProcessing,
		DisputeTypeCardNotPresentFraud,
		DisputeTypeCreditNotProcessed,
		DisputeTypePaymentByOtherMeans,
	}

	for _, dt := range validTypes {
		if !dt.IsValid() {
			t.Errorf("DisputeType(%q).IsValid() = false, want true", dt)
		}
	}

	if DisputeType("invalid").IsValid() {
		t.Error("DisputeType(\"invalid\").IsValid() = true, want false")
	}
}

func TestFraudNotificationStatusIsValid(t *testing.T) {
	validStatuses := []FraudNotificationStatus{
		FraudNotificationStatusNotified,
		FraudNotificationStatusUndone,
	}

	for _, status := range validStatuses {
		if !status.IsValid() {
			t.Errorf("FraudNotificationStatus(%q).IsValid() = false, want true", status)
		}
	}

	if FraudNotificationStatus("invalid").IsValid() {
		t.Error("FraudNotificationStatus(\"invalid\").IsValid() = true, want false")
	}
}

func TestHealthStatusIsValid(t *testing.T) {
	validStatuses := []HealthStatus{
		HealthStatusHealthy,
		HealthStatusUnhealthy,
		HealthStatusUnstable,
		HealthStatusRecovering,
		HealthStatusNoData,
	}

	for _, status := range validStatuses {
		if !status.IsValid() {
			t.Errorf("HealthStatus(%q).IsValid() = false, want true", status)
		}
	}

	if HealthStatus("invalid").IsValid() {
		t.Error("HealthStatus(\"invalid\").IsValid() = true, want false")
	}
}

func TestPaymentTypeIsValid(t *testing.T) {
	validTypes := []PaymentType{
		PaymentTypeA2a,
		PaymentTypeB2a,
		PaymentTypeP2p,
		PaymentTypeCsb,
		PaymentTypeDsb,
		PaymentTypeB2b,
		PaymentTypeM2m,
	}

	for _, pt := range validTypes {
		if !pt.IsValid() {
			t.Errorf("PaymentType(%q).IsValid() = false, want true", pt)
		}
	}

	if PaymentType("invalid").IsValid() {
		t.Error("PaymentType(\"invalid\").IsValid() = true, want false")
	}
}

func TestSourceIsValid(t *testing.T) {
	validSources := []Source{
		SourcePaysmart,
		SourceIssuer,
		SourceBrand,
	}

	for _, src := range validSources {
		if !src.IsValid() {
			t.Errorf("Source(%q).IsValid() = false, want true", src)
		}
	}

	if Source("invalid").IsValid() {
		t.Error("Source(\"invalid\").IsValid() = true, want false")
	}
}

func TestCardholderTypeIsValid(t *testing.T) {
	validTypes := []CardholderType{
		CardholderTypeMain,
		CardholderTypeAdditional,
	}

	for _, ct := range validTypes {
		if !ct.IsValid() {
			t.Errorf("CardholderType(%q).IsValid() = false, want true", ct)
		}
	}

	if CardholderType("invalid").IsValid() {
		t.Error("CardholderType(\"invalid\").IsValid() = true, want false")
	}
}

func TestCivilStatusIsValid(t *testing.T) {
	validStatuses := []CivilStatus{
		CivilStatusSingle,
		CivilStatusMarried,
		CivilStatusDivorced,
		CivilStatusWidowed,
		CivilStatusOther,
	}

	for _, status := range validStatuses {
		if !status.IsValid() {
			t.Errorf("CivilStatus(%q).IsValid() = false, want true", status)
		}
	}

	if CivilStatus("invalid").IsValid() {
		t.Error("CivilStatus(\"invalid\").IsValid() = true, want false")
	}
}

func TestFeeTypeIsValid(t *testing.T) {
	validTypes := []FeeType{
		FeeTypeIof,
		FeeTypeMarkup,
		FeeTypeBoardingFee,
		FeeTypeWithdrawalFee,
		FeeTypeInterest,
		FeeTypeOthers,
	}

	for _, ft := range validTypes {
		if !ft.IsValid() {
			t.Errorf("FeeType(%q).IsValid() = false, want true", ft)
		}
	}

	if FeeType("invalid").IsValid() {
		t.Error("FeeType(\"invalid\").IsValid() = true, want false")
	}
}

func TestDebitOrCreditIsValid(t *testing.T) {
	if !DebitOrCreditDebit.IsValid() {
		t.Error("DebitOrCreditDebit.IsValid() = false, want true")
	}
	if !DebitOrCreditCredit.IsValid() {
		t.Error("DebitOrCreditCredit.IsValid() = false, want true")
	}
	if DebitOrCredit("invalid").IsValid() {
		t.Error("DebitOrCredit(\"invalid\").IsValid() = true, want false")
	}
}

func TestInstallmentFinTypeIsValid(t *testing.T) {
	validTypes := []InstallmentFinType{
		InstallmentFinTypeInterestfree,
		InstallmentFinTypeWithinterest,
		InstallmentFinTypeCdc,
	}

	for _, ift := range validTypes {
		if !ift.IsValid() {
			t.Errorf("InstallmentFinType(%q).IsValid() = false, want true", ift)
		}
	}

	if InstallmentFinType("invalid").IsValid() {
		t.Error("InstallmentFinType(\"invalid\").IsValid() = true, want false")
	}
}

func TestDisputeStageIsValid(t *testing.T) {
	validStages := []DisputeStage{
		DisputeStageChargeback,
		DisputeStageRestatement,
		DisputeStagePreArbitration,
		DisputeStageArbitration,
		DisputeStageLost,
		DisputeStageUndone,
		DisputeStageSucessful,
		DisputeStageWithdrawal,
	}

	for _, stage := range validStages {
		if !stage.IsValid() {
			t.Errorf("DisputeStage(%q).IsValid() = false, want true", stage)
		}
	}

	if DisputeStage("invalid").IsValid() {
		t.Error("DisputeStage(\"invalid\").IsValid() = true, want false")
	}
}

func TestInclusiveTransactionCodeIsValid(t *testing.T) {
	if !InclusiveTransactionCodeTe10.IsValid() {
		t.Error("InclusiveTransactionCodeTe10.IsValid() = false, want true")
	}
	if !InclusiveTransactionCodeTe20.IsValid() {
		t.Error("InclusiveTransactionCodeTe20.IsValid() = false, want true")
	}
	if InclusiveTransactionCode("TE30").IsValid() {
		t.Error("InclusiveTransactionCode(\"TE30\").IsValid() = true, want false")
	}
}

func TestCreditRequestTypeIsValid(t *testing.T) {
	if !CreditRequestTypePayment.IsValid() {
		t.Error("CreditRequestTypePayment.IsValid() = false, want true")
	}
	if !CreditRequestTypeAdjustment.IsValid() {
		t.Error("CreditRequestTypeAdjustment.IsValid() = false, want true")
	}
	if CreditRequestType("invalid").IsValid() {
		t.Error("CreditRequestType(\"invalid\").IsValid() = true, want false")
	}
}

func TestDebitRequestTypeIsValid(t *testing.T) {
	validTypes := []DebitRequestType{
		DebitRequestTypeAnnuityfee,
		DebitRequestTypeIssuefee,
		DebitRequestTypeReissuefee,
		DebitRequestTypeRenewalfee,
		DebitRequestTypeAdjustment,
	}

	for _, drt := range validTypes {
		if !drt.IsValid() {
			t.Errorf("DebitRequestType(%q).IsValid() = false, want true", drt)
		}
	}

	if DebitRequestType("invalid").IsValid() {
		t.Error("DebitRequestType(\"invalid\").IsValid() = true, want false")
	}
}

func TestFunctionalityCodeIsValid(t *testing.T) {
	validCodes := []FunctionalityCode{
		FunctionalityCodeContactless,
		FunctionalityCodeMagneticStripe,
		FunctionalityCodeECommerce,
		FunctionalityCodeWithdrawal,
	}

	for _, fc := range validCodes {
		if !fc.IsValid() {
			t.Errorf("FunctionalityCode(%q).IsValid() = false, want true", fc)
		}
	}

	if FunctionalityCode("invalid").IsValid() {
		t.Error("FunctionalityCode(\"invalid\").IsValid() = true, want false")
	}
}

func TestCampaignStatusIsValid(t *testing.T) {
	if !CampaignStatusEnabled.IsValid() {
		t.Error("CampaignStatusEnabled.IsValid() = false, want true")
	}
	if !CampaignStatusDisabled.IsValid() {
		t.Error("CampaignStatusDisabled.IsValid() = false, want true")
	}
	if CampaignStatus("invalid").IsValid() {
		t.Error("CampaignStatus(\"invalid\").IsValid() = true, want false")
	}
}

func TestVirtualCardStatusIsValid(t *testing.T) {
	validStatuses := []VirtualCardStatus{
		VirtualCardStatusActive,
		VirtualCardStatusInactive,
		VirtualCardStatusExpired,
	}

	for _, status := range validStatuses {
		if !status.IsValid() {
			t.Errorf("VirtualCardStatus(%q).IsValid() = false, want true", status)
		}
	}

	if VirtualCardStatus("invalid").IsValid() {
		t.Error("VirtualCardStatus(\"invalid\").IsValid() = true, want false")
	}
}

func TestCreditMattressStatusIsValid(t *testing.T) {
	validStatuses := []CreditMattressStatus{
		CreditMattressStatusConfirmed,
		CreditMattressStatusCanceled,
		CreditMattressStatusPending,
	}

	for _, status := range validStatuses {
		if !status.IsValid() {
			t.Errorf("CreditMattressStatus(%q).IsValid() = false, want true", status)
		}
	}

	if CreditMattressStatus("invalid").IsValid() {
		t.Error("CreditMattressStatus(\"invalid\").IsValid() = true, want false")
	}
}
