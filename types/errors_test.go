package types

import (
	"strings"
	"testing"
)

func TestErrorCode_IsSuccess(t *testing.T) {
	if !ErrorCodeSuccess.IsSuccess() {
		t.Error("ErrorCodeSuccess.IsSuccess() should be true")
	}
	if ErrorCodeInvalidCard.IsSuccess() {
		t.Error("ErrorCodeInvalidCard.IsSuccess() should be false")
	}
}

func TestErrorCode_Message(t *testing.T) {
	tests := []struct {
		code    ErrorCode
		wantMsg string
	}{
		{ErrorCodeSuccess, "Comando concluído com sucesso"},
		{ErrorCodeInvalidCard, "Cartão inválido"},
		{ErrorCodeAPIKeyInvalid, "API-Key não enviada ou inválida"},
		{ErrorCodeInternalError, "Erro interno paySmart"},
		{ErrorCode(9999), ""}, // unknown code
	}

	for _, tt := range tests {
		if got := tt.code.Message(); got != tt.wantMsg {
			t.Errorf("ErrorCode(%d).Message() = %q, want %q", tt.code, got, tt.wantMsg)
		}
	}
}

func TestErrorCode_Error(t *testing.T) {
	// Known code
	err := ErrorCodeInvalidCard.Error()
	if !strings.Contains(err, "14") || !strings.Contains(err, "Cartão inválido") {
		t.Errorf("ErrorCode(14).Error() = %q, want code and message", err)
	}

	// Unknown code
	err = ErrorCode(9999).Error()
	if !strings.Contains(err, "9999") {
		t.Errorf("ErrorCode(9999).Error() = %q, want code", err)
	}
}

func TestErrorCode_ImplementsError(t *testing.T) {
	var err error = ErrorCodeInvalidCard
	if err.Error() == "" {
		t.Error("ErrorCode should implement error interface")
	}
}

func TestErrorCode_AllCodesHaveMessages(t *testing.T) {
	codes := []ErrorCode{
		ErrorCodeSuccess, ErrorCodeCommunication,
		ErrorCodeInvalidCard, ErrorCodeCardRequestedEmissionReject, ErrorCodeCardInEmission,
		ErrorCodeInvalidBlockingCode, ErrorCodeBlockingCodeUnlockMismatch, ErrorCodeBlockingCodeLockMismatch,
		ErrorCodeUnlockNotPermitted, ErrorCodeCardAlreadyUnlocked, ErrorCodeCardBlockedHigherSeverity,
		ErrorCodeCardStatusPreventsReissuance, ErrorCodePendingReissuanceExists,
		ErrorCodeInvalidPasswordData, ErrorCodePasswordLocked, ErrorCodeInvalidPassword, ErrorCodeInvalidPasswordCVV,
		ErrorCodeSystemUnavailable,
		ErrorCodeInvalidData, ErrorCodeAccountOrCardBlocked, ErrorCodeInvalidAccount, ErrorCodeMandatoryDataMissing,
		ErrorCodeInvalidFieldFormat, ErrorCodeCardAlreadyAssociated, ErrorCodeAnonymousAccountNoNewCards,
		ErrorCodeAccountCanceled, ErrorCodeValidationFailures, ErrorCodeTransactionNotFound,
		ErrorCodeInvalidTokenizationResponse, ErrorCodeNoVirtualPANAvailable, ErrorCodeCryptographicPANFailed,
		ErrorCodeVirtualCardAnonymousDisallowed, ErrorCodeAnonymousCardPrintingOnly, ErrorCodeQRCodeTokenFailed,
		ErrorCodeFlagAPIKeyError, ErrorCodeInvalidFlagQRCodeResponse, ErrorCodeInvalidPINFormat,
		ErrorCodeVirtualCardDataUnavailable, ErrorCodeQRCodeEncryptionFailed, ErrorCodeDuplicateQRCodeTransaction,
		ErrorCodeAcquirerNotConfiguredQR, ErrorCodeAccountClosed, ErrorCodeCardCanceled,
		ErrorCodeDebitConfirmationFailed, ErrorCodeEmbossingFileNoStatus, ErrorCodeCardNoEmbossingFile,
		ErrorCodePhysicalCardChildProhibited, ErrorCodeAccountChildProductProhibited,
		ErrorCodeVirtualCardParentProductError, ErrorCodeCardholderAssociationRequired,
		ErrorCodeAltBindingKeyMultipleCards, ErrorCodeAltBindingKeySecurityMismatch,
		ErrorCodeAltBindingKeyOrDataRequired, ErrorCodeAltBindingKeyInUse, ErrorCodeNewAltBindingKeyRequired,
		ErrorCodeLinkedPhysicalCardInvalid, ErrorCodePhysicalCardUpdateInapplicable,
		ErrorCodeProductNotRegistered, ErrorCodeIdempotentRequestRunning, ErrorCodeAccountAlreadyExists,
		ErrorCodeInvalidJSONStructure, ErrorCodeAPIKeyInvalid, ErrorCodeTransactionAPIError, ErrorCodeIdempotencyKeyMissing,
		ErrorCodePANHashFailed, ErrorCodeInvoiceExpirationInvalid, ErrorCodeAccountCreationBalanceErr,
		ErrorCodeQRCodeEloConfigNotFound, ErrorCodeNoDisputesFound, ErrorCodeDisputeNotInResubmission,
		ErrorCodeDisputeNotFound, ErrorCodeDisputeInitialTxMissing, ErrorCodeInclusiveTxAlreadyExists,
		ErrorCodePartialValueExceedsTotal, ErrorCodeDisputeReasonNotFound, ErrorCodeTransactionParamsNotFound,
		ErrorCodeDisputedTxNotSettled, ErrorCodeDisputedTxCanceled, ErrorCodeDisputeAlreadyExists,
		ErrorCodeTransactionAPIConfigError, ErrorCodeBufferServiceConfigError, ErrorCodeTokenizationConfigError,
		ErrorCodeIssuerConfigNotFound, ErrorCodePINValidationConfigError, ErrorCodeInternalWriteError, ErrorCodeInternalError,
	}

	for _, code := range codes {
		if code.Message() == "" {
			t.Errorf("ErrorCode(%d) has no message defined", code)
		}
	}
}
