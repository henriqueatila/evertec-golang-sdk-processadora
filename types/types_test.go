package types

import (
	"encoding/json"
	"testing"
)

// TestAmount_JSON tests Amount serialization/deserialization
func TestAmount_JSON(t *testing.T) {
	tests := []struct {
		name    string
		amount  Amount
		jsonStr string
	}{
		{
			name:    "basic amount",
			amount:  Amount{Amount: 10000, CurrencyCode: 986},
			jsonStr: `{"amount":10000,"currencyCode":986}`,
		},
		{
			name:    "zero amount",
			amount:  Amount{Amount: 0, CurrencyCode: 986},
			jsonStr: `{"amount":0,"currencyCode":986}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Marshal
			data, err := json.Marshal(tt.amount)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if string(data) != tt.jsonStr {
				t.Errorf("Marshal = %s, want %s", data, tt.jsonStr)
			}

			// Test Unmarshal
			var decoded Amount
			if err := json.Unmarshal([]byte(tt.jsonStr), &decoded); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if decoded != tt.amount {
				t.Errorf("Unmarshal = %+v, want %+v", decoded, tt.amount)
			}
		})
	}
}

// TestCard_JSON tests Card serialization
func TestCard_JSON(t *testing.T) {
	card := Card{
		CardID:      "card123",
		AccountID:   "acc123",
		Last4Digits: "1234",
		Status:      CardStatusActive,
	}

	data, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Card
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.CardID != card.CardID {
		t.Errorf("CardID = %s, want %s", decoded.CardID, card.CardID)
	}
	if decoded.Last4Digits != card.Last4Digits {
		t.Errorf("Last4Digits = %s, want %s", decoded.Last4Digits, card.Last4Digits)
	}
}

// TestEntryMode_Values tests EntryMode constants
func TestEntryMode_Values(t *testing.T) {
	tests := []struct {
		mode EntryMode
		want string
	}{
		{EntryModeEmvcontact, "emvContact"},
		{EntryModeEmvcontactless, "emvContactless"},
		{EntryModeHce, "hce"},
		{EntryModeQrcode, "qrCode"},
		{EntryModeEcommerce, "ecommerce"},
		{EntryModeMagstripe, "magstripe"},
		{EntryModeManual, "manual"},
	}

	for _, tt := range tests {
		if string(tt.mode) != tt.want {
			t.Errorf("EntryMode = %s, want %s", tt.mode, tt.want)
		}
	}
}

// TestBrand_Values tests Brand constants
func TestBrand_Values(t *testing.T) {
	tests := []struct {
		brand Brand
		want  string
	}{
		{BrandMastercard, "MASTERCARD"},
		{BrandVisa, "VISA"},
		{BrandElo, "ELO"},
		{BrandHipercard, "HIPERCARD"},
		{BrandAmex, "AMEX"},
	}

	for _, tt := range tests {
		if string(tt.brand) != tt.want {
			t.Errorf("Brand = %s, want %s", tt.brand, tt.want)
		}
	}
}

// TestTransactionStatus_Values tests TransactionStatus constants
func TestTransactionStatus_Values(t *testing.T) {
	tests := []struct {
		status TransactionStatus
		want   string
	}{
		{TransactionStatusReceived, "received"},
		{TransactionStatusPOSted, "posted"},
		{TransactionStatusInDispute, "in_dispute"},
		{TransactionStatusCanceled, "canceled"},
		{TransactionStatusRefunded, "refunded"},
		{TransactionStatusPartiallyRefunded, "partially_refunded"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.want {
			t.Errorf("TransactionStatus = %s, want %s", tt.status, tt.want)
		}
	}
}

// TestResultData_JSON tests ResultData serialization
func TestResultData_JSON(t *testing.T) {
	result := ResultData{
		ResultCode:        0,
		ResultDescription: "Success",
		IssuerRequestID:   "req123",
		PsResponseID:      "ps123",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded ResultData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.ResultCode != result.ResultCode {
		t.Errorf("ResultCode = %d, want %d", decoded.ResultCode, result.ResultCode)
	}
	if decoded.ResultDescription != result.ResultDescription {
		t.Errorf("ResultDescription = %s, want %s", decoded.ResultDescription, result.ResultDescription)
	}
}

// TestTokenPaymentData_JSON tests TokenPaymentData serialization
func TestTokenPaymentData_JSON(t *testing.T) {
	data := TokenPaymentData{
		RequesterIDToken: "req123",
		PAN:              "4111111111111111",
		TokenPSN:         "001",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded TokenPaymentData
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.RequesterIDToken != data.RequesterIDToken {
		t.Errorf("RequesterIDToken = %s, want %s", decoded.RequesterIDToken, data.RequesterIDToken)
	}
}

// TestInstallmentDetails_JSON tests InstallmentDetails serialization
func TestInstallmentDetails_JSON(t *testing.T) {
	details := InstallmentDetails{
		FinType: "withInterest",
	}

	data, err := json.Marshal(details)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded InstallmentDetails
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.FinType != details.FinType {
		t.Errorf("FinType = %s, want %s", decoded.FinType, details.FinType)
	}
}

// TestHolder_JSON tests Holder serialization
func TestHolder_JSON(t *testing.T) {
	holder := Holder{
		Name:         "John Doe",
		Document:     "12345678900",
		DocumentType: "CPF",
		Email:        "john@example.com",
	}

	data, err := json.Marshal(holder)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Holder
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Name != holder.Name {
		t.Errorf("Name = %s, want %s", decoded.Name, holder.Name)
	}
	if decoded.Document != holder.Document {
		t.Errorf("Document = %s, want %s", decoded.Document, holder.Document)
	}
}

// TestAddress_JSON tests Address serialization
func TestAddress_JSON(t *testing.T) {
	addr := Address{
		AddressLine1: "Rua Test 123",
		AddressLine2: "Apt 456",
		Neighborhood: "Centro",
		City:         "São Paulo",
		State:        "SP",
		Zipcode:      "01000-000",
		Country:      "BR",
	}

	data, err := json.Marshal(addr)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Address
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.City != addr.City {
		t.Errorf("City = %s, want %s", decoded.City, addr.City)
	}
}

// TestEstablishment_JSON tests Establishment serialization
func TestEstablishment_JSON(t *testing.T) {
	est := Establishment{
		Name:       "Test Store",
		MCC:        "5411",
		MerchantID: "merchant123",
		City:       "São Paulo",
		Country:    "BR",
	}

	data, err := json.Marshal(est)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Establishment
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.MCC != est.MCC {
		t.Errorf("MCC = %s, want %s", decoded.MCC, est.MCC)
	}
}

// TestFee_JSON tests Fee serialization
func TestFee_JSON(t *testing.T) {
	fee := Fee{
		Type:   "interchange",
		Amount: Amount{Amount: 150, CurrencyCode: 986},
	}

	data, err := json.Marshal(fee)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Fee
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Amount.Amount != fee.Amount.Amount {
		t.Errorf("Amount = %d, want %d", decoded.Amount.Amount, fee.Amount.Amount)
	}
}

// TestWallet_Values tests Wallet constants
func TestWallet_Values(t *testing.T) {
	tests := []struct {
		wallet Wallet
		want   string
	}{
		{WalletApplePay, "APPLE_PAY"},
		{WalletGooglePay, "GOOGLE_PAY"},
		{WalletSamsungPay, "SAMSUNG_PAY"},
	}

	for _, tt := range tests {
		if string(tt.wallet) != tt.want {
			t.Errorf("Wallet = %s, want %s", tt.wallet, tt.want)
		}
	}
}

// TestTokenStatus_Values tests TokenStatus constants
func TestTokenStatus_Values(t *testing.T) {
	tests := []struct {
		status TokenStatus
		want   string
	}{
		{TokenStatusActive, "ACTIVE"},
		{TokenStatusInactive, "INACTIVE"},
		{TokenStatusSuspended, "SUSPENDED"},
		{TokenStatusDeleted, "DELETED"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.want {
			t.Errorf("TokenStatus = %s, want %s", tt.status, tt.want)
		}
	}
}

// TestProcessingCode_Values tests ProcessingCode constants
func TestProcessingCode_Values(t *testing.T) {
	tests := []struct {
		code ProcessingCode
		want string
	}{
		{ProcessingCodePurchase, "00"},
		{ProcessingCodeWithdrawal, "01"},
		{ProcessingCodeCashback, "09"},
		{ProcessingCodeRefund, "20"},
	}

	for _, tt := range tests {
		if string(tt.code) != tt.want {
			t.Errorf("ProcessingCode = %s, want %s", tt.code, tt.want)
		}
	}
}

// TestAccountStatus_Values tests AccountStatus constants
func TestAccountStatus_Values(t *testing.T) {
	tests := []struct {
		status AccountStatus
		want   string
	}{
		{AccountStatusActive, "active"},
		{AccountStatusCanceled, "canceled"},
		{AccountStatusClosed, "closed"},
		{AccountStatusBlocked, "blocked"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.want {
			t.Errorf("AccountStatus = %s, want %s", tt.status, tt.want)
		}
	}
}

// TestCardStatus_Values tests CardStatus constants
func TestCardStatus_Values(t *testing.T) {
	tests := []struct {
		status CardStatus
		want   string
	}{
		{CardStatusActive, "active"},
		{CardStatusBlocked, "blocked"},
		{CardStatusCancelled, "cancelled"},
		{CardStatusIssuing, "issuing"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.want {
			t.Errorf("CardStatus = %s, want %s", tt.status, tt.want)
		}
	}
}
