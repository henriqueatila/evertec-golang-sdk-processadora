// Package validation provides tests to validate SDK types against OpenAPI schemas.
package validation

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// TestAccount_SchemaCompliance validates Account struct against OpenAPI schema
func TestAccount_SchemaCompliance(t *testing.T) {
	status := types.AccountStatusActive
	account := types.Account{
		AccountID:     "acc123",
		Status:        &status,
		PsProductCode: "CREDIT_CARD",
	}

	// Test JSON serialization
	data, err := json.Marshal(account)
	if err != nil {
		t.Fatalf("Failed to marshal Account: %v", err)
	}

	// Verify JSON contains expected fields
	var jsonMap map[string]any
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	expectedFields := []string{"accountId", "status", "psProductCode"}
	for _, field := range expectedFields {
		if _, ok := jsonMap[field]; !ok {
			t.Errorf("Missing expected field: %s", field)
		}
	}
}

// TestCard_SchemaCompliance validates Card struct against OpenAPI schema
func TestCard_SchemaCompliance(t *testing.T) {
	card := types.Card{
		CardID:         "card123",
		AccountID:      "acc123",
		Last4Digits:    "1234",
		Status:         types.CardStatusActive,
		ExpirationDate: "12/2025",
	}

	data, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("Failed to marshal Card: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	// cardId should serialize to "cardId" not "card_id"
	if _, ok := jsonMap["cardId"]; !ok {
		t.Error("Card.CardID should serialize as 'cardId'")
	}

	if _, ok := jsonMap["last4Digits"]; !ok {
		t.Error("Card.Last4Digits should serialize as 'last4Digits'")
	}
}

// TestTransaction_SchemaCompliance validates Transaction struct
func TestTransaction_SchemaCompliance(t *testing.T) {
	tx := types.Transaction{
		TransactionID: "tx123",
		Amount:        types.AmountTransaction{Amount: 10000, CurrencyCode: 986},
	}

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("Failed to marshal Transaction: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	requiredFields := []string{"transactionId", "amount"}
	for _, field := range requiredFields {
		if _, ok := jsonMap[field]; !ok {
			t.Errorf("Missing required field: %s", field)
		}
	}

	// Verify amount structure
	if amount, ok := jsonMap["amount"].(map[string]any); ok {
		if _, ok := amount["amount"]; !ok {
			t.Error("Amount missing 'amount' field")
		}
		if _, ok := amount["currencyCode"]; !ok {
			t.Error("Amount missing 'currencyCode' field")
		}
	}
}

// TestTravelNotice_SchemaCompliance validates TravelNotice struct
func TestTravelNotice_SchemaCompliance(t *testing.T) {
	notice := types.TravelNotice{
		CountryCodes: "US",
		BeginDate:    "2024-01-01",
		EndDate:      "2024-01-15",
	}

	data, err := json.Marshal(notice)
	if err != nil {
		t.Fatalf("Failed to marshal TravelNotice: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	// Verify date format (should be YYYY-MM-DD)
	if bd, ok := jsonMap["beginDate"].(string); ok {
		if !isValidDateFormat(bd) {
			t.Errorf("beginDate format invalid: %s (expected YYYY-MM-DD)", bd)
		}
	}
}

// TestNewTravelNotice_SchemaCompliance validates NewTravelNotice request struct
func TestNewTravelNotice_SchemaCompliance(t *testing.T) {
	req := types.NewTravelNotice{
		CountryCodes: "US,CA,MX",
		BeginDate:    "2024-01-01",
		EndDate:      "2024-01-15",
		Cards:        []types.Card{{CardID: "card123"}},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal NewTravelNotice: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	// Required fields per OpenAPI
	requiredFields := []string{"countryCodes", "beginDate", "endDate", "cards"}
	for _, field := range requiredFields {
		if _, ok := jsonMap[field]; !ok {
			t.Errorf("Missing required field: %s", field)
		}
	}
}

// TestDispute_SchemaCompliance validates Dispute struct
func TestDispute_SchemaCompliance(t *testing.T) {
	dispute := types.Dispute{
		DisputeID:     "disp123",
		DisputeStatus: types.DisputeStatusAccepted,
		DisputeType:   types.DisputeTypeCardNotPresentFraud,
		DisputeDate:   "2024-01-01",
		CurrentStage:  "chargeback",
	}

	data, err := json.Marshal(dispute)
	if err != nil {
		t.Fatalf("Failed to marshal Dispute: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	if _, ok := jsonMap["disputeId"]; !ok {
		t.Error("Dispute.DisputeID should serialize as 'disputeId'")
	}
	if _, ok := jsonMap["disputeStatus"]; !ok {
		t.Error("Dispute should have 'disputeStatus' field")
	}
	if _, ok := jsonMap["disputeType"]; !ok {
		t.Error("Dispute should have 'disputeType' field")
	}
}

// TestInclusiveTransaction_SchemaCompliance validates InclusiveTransaction struct
func TestInclusiveTransaction_SchemaCompliance(t *testing.T) {
	inclusive := types.InclusiveTransaction{
		InclusiveTransactionID: "inc123",
		AccountID:              "acc123",
		Code:                   "TE10",
	}

	data, err := json.Marshal(inclusive)
	if err != nil {
		t.Fatalf("Failed to marshal InclusiveTransaction: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	// Verify code is TE10 or TE20
	if c, ok := jsonMap["code"].(string); ok {
		if c != "TE10" && c != "TE20" {
			t.Errorf("code should be TE10 or TE20, got: %s", c)
		}
	}
}

// TestDeviceToken_SchemaCompliance validates DeviceToken struct
func TestDeviceToken_SchemaCompliance(t *testing.T) {
	token := types.DeviceToken{
		DeviceTokenID: "token123",
		Wallet:        "APPLE_PAY",
	}

	data, err := json.Marshal(token)
	if err != nil {
		t.Fatalf("Failed to marshal DeviceToken: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	if _, ok := jsonMap["deviceTokenId"]; !ok {
		t.Error("DeviceToken.DeviceTokenID should serialize as 'deviceTokenId'")
	}

	// Wallet should be valid wallet type
	if w, ok := jsonMap["wallet"].(string); ok {
		validWallets := []string{"APPLE_PAY", "GOOGLE_PAY", "SAMSUNG_PAY"}
		valid := false
		for _, vw := range validWallets {
			if w == vw {
				valid = true
				break
			}
		}
		if !valid {
			t.Errorf("wallet should be a valid wallet type, got: %s", w)
		}
	}
}

// TestResultData_SchemaCompliance validates ResultData struct
func TestResultData_SchemaCompliance(t *testing.T) {
	result := types.ResultData{
		ResultCode:        0,
		ResultDescription: "Success",
		IssuerRequestID:   "req123",
		PsResponseID:      "ps123",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal ResultData: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	// resultCode should be integer
	if _, ok := jsonMap["resultCode"].(float64); !ok {
		t.Error("resultCode should be a number")
	}

	// psResponseId should be present
	if _, ok := jsonMap["psResponseId"]; !ok {
		t.Error("psResponseId should be present")
	}
}

// TestAmount_SchemaCompliance validates Amount struct
func TestAmount_SchemaCompliance(t *testing.T) {
	amount := types.Amount{
		Amount:       10050, // R$ 100,50 in centavos
		CurrencyCode: 986,   // BRL ISO 4217
	}

	data, err := json.Marshal(amount)
	if err != nil {
		t.Fatalf("Failed to marshal Amount: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	// amount should be integer (centavos)
	if v, ok := jsonMap["amount"].(float64); ok {
		if v != 10050 {
			t.Errorf("amount = %v, want 10050", v)
		}
	} else {
		t.Error("amount should be a number")
	}

	// currencyCode should be ISO 4217 numeric
	if c, ok := jsonMap["currencyCode"].(float64); ok {
		if c != 986 {
			t.Errorf("currencyCode = %v, want 986 (BRL)", c)
		}
	}
}

// TestVirtualCard_SchemaCompliance validates virtual card types
func TestVirtualCard_SchemaCompliance(t *testing.T) {
	vcard := types.VirtualCardDescriptor{
		VCardID:     "vcard123",
		VPAN:        "4111111111111111",
		VDateExp:    "12/2025",
		VCVV:        "123",
		VCardholder: "JOHN DOE",
	}

	data, err := json.Marshal(vcard)
	if err != nil {
		t.Fatalf("Failed to marshal VirtualCardDescriptor: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	// Verify all required fields
	requiredFields := []string{"vCardId", "vPan", "vDateExp", "vCvv", "vCardholder"}
	for _, field := range requiredFields {
		if _, ok := jsonMap[field]; !ok {
			t.Errorf("Missing field: %s", field)
		}
	}

	// vCvv should be exactly 3 digits
	if cvv, ok := jsonMap["vCvv"].(string); ok {
		if len(cvv) != 3 {
			t.Errorf("vCvv should be 3 digits, got: %s", cvv)
		}
	}
}

// TestCampaign_SchemaCompliance validates Campaign struct
func TestCampaign_SchemaCompliance(t *testing.T) {
	campaign := types.CampaignObject{
		AgentID:      "agent123",
		ProductCode:  "PRODUCT1",
		CampaignName: "Test Campaign",
	}

	data, err := json.Marshal(campaign)
	if err != nil {
		t.Fatalf("Failed to marshal CampaignObject: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	// Required fields
	requiredFields := []string{"agentId", "productCode", "campaignName"}
	for _, field := range requiredFields {
		if _, ok := jsonMap[field]; !ok {
			t.Errorf("Missing required field: %s", field)
		}
	}
}

// TestParseQrCodeParams_SchemaCompliance validates QR Code parsing params
func TestParseQrCodeParams_SchemaCompliance(t *testing.T) {
	params := types.ParseQrCodeParams{
		QRCode: "00020126580014br.gov.bcb.pix0136a629532e-7693-4846-852d-1bbff817b5a8520400005303986540510.005802BR5913Test Merchant6008BRASILIA62290525202012345678901234567890163045B13",
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal ParseQrCodeParams: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	// qrCode is required
	if _, ok := jsonMap["qrCode"]; !ok {
		t.Error("qrCode field is required")
	}

	// QR code should start with "0002" (EMV standard)
	if qr, ok := jsonMap["qrCode"].(string); ok {
		if !strings.HasPrefix(qr, "0002") {
			t.Log("Warning: QR code doesn't follow EMV standard (should start with 0002)")
		}
	}
}

// TestHCEProvision_SchemaCompliance validates HCE provisioning types
func TestHCEProvision_SchemaCompliance(t *testing.T) {
	req := types.HCEProvisionRequest{
		Mtv: map[string]any{
			"prov_data":        "encrypted_data_here",
			"issuer_device_id": "device456",
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal HCEProvisionRequest: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	// Verify MTV structure exists
	if mtv, ok := jsonMap["mtv"].(map[string]any); ok {
		if _, ok := mtv["prov_data"]; !ok {
			t.Error("mtv.prov_data is required for HCE provisioning")
		}
		if _, ok := mtv["issuer_device_id"]; !ok {
			t.Error("mtv.issuer_device_id is required for HCE provisioning")
		}
	}
}

// TestCobranded_SchemaCompliance validates cobranded merchant types
func TestCobranded_SchemaCompliance(t *testing.T) {
	req := types.MerchantVanRequest{
		AcquirerID:      1,
		Document:        "12345678901234",
		DocumentType:    "cnpj",
		LegalName:       "Test Company LTDA",
		MerchantVanCode: "MERCHANT001",
		Products: []types.ProductsDTO{
			{ProductID: 1, ClosedArrangementEntryModes: []string{"CHIP", "CONTACTLESS"}},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal MerchantVanRequest: %v", err)
	}

	var jsonMap map[string]any
	_ = json.Unmarshal(data, &jsonMap)

	// Required fields
	requiredFields := []string{"acquirerId", "document", "documentType", "legalName", "merchantVanCode", "products"}
	for _, field := range requiredFields {
		if _, ok := jsonMap[field]; !ok {
			t.Errorf("Missing required field: %s", field)
		}
	}

	// documentType should be cpf or cnpj
	if dt, ok := jsonMap["documentType"].(string); ok {
		if dt != "cpf" && dt != "cnpj" {
			t.Errorf("documentType should be 'cpf' or 'cnpj', got: %s", dt)
		}
	}
}

// TestEnumValues_SchemaCompliance validates enum constants match OpenAPI
func TestEnumValues_SchemaCompliance(t *testing.T) {
	// Brand values
	brands := []types.Brand{
		types.BrandMastercard,
		types.BrandVisa,
		types.BrandElo,
	}
	for _, b := range brands {
		if string(b) == "" {
			t.Errorf("Brand value is empty")
		}
	}

	// Entry modes
	entryModes := []types.EntryMode{
		types.EntryModeEmvcontact,
		types.EntryModeEmvcontactless,
		types.EntryModeMagstripe,
		types.EntryModeEcommerce,
	}
	for _, em := range entryModes {
		if string(em) == "" {
			t.Errorf("EntryMode value is empty")
		}
	}

	// Transaction status
	statuses := []types.TransactionStatus{
		types.TransactionStatusReceived,
		types.TransactionStatusPOSted,
		types.TransactionStatusCanceled,
		types.TransactionStatusRefunded,
	}
	for _, s := range statuses {
		if string(s) == "" {
			t.Errorf("TransactionStatus value is empty")
		}
	}
}

// TestJSONTagsMatchOpenAPI validates that JSON tags match OpenAPI field names
func TestJSONTagsMatchOpenAPI(t *testing.T) {
	// Test Account struct
	accountType := reflect.TypeOf(types.Account{})
	expectedTags := map[string]string{
		"AccountID":     "accountId",
		"Status":        "status",
		"PsProductCode": "psProductCode",
	}
	validateJSONTags(t, accountType, expectedTags)

	// Test Card struct
	cardType := reflect.TypeOf(types.Card{})
	cardExpectedTags := map[string]string{
		"CardID":      "cardId",
		"Last4Digits": "last4Digits",
		"Status":      "status",
	}
	validateJSONTags(t, cardType, cardExpectedTags)

	// Test Amount struct
	amountType := reflect.TypeOf(types.Amount{})
	amountExpectedTags := map[string]string{
		"Amount":       "amount",
		"CurrencyCode": "currencyCode",
	}
	validateJSONTags(t, amountType, amountExpectedTags)
}

// Helper functions

func validateJSONTags(t *testing.T, typ reflect.Type, expected map[string]string) {
	for fieldName, expectedTag := range expected {
		field, ok := typ.FieldByName(fieldName)
		if !ok {
			t.Errorf("Field %s not found in %s", fieldName, typ.Name())
			continue
		}
		tag := field.Tag.Get("json")
		// Remove omitempty suffix for comparison
		tag = strings.Split(tag, ",")[0]
		if tag != expectedTag {
			t.Errorf("%s.%s JSON tag = %s, want %s", typ.Name(), fieldName, tag, expectedTag)
		}
	}
}

func isValidDateFormat(date string) bool {
	// YYYY-MM-DD format
	if len(date) != 10 {
		return false
	}
	if date[4] != '-' || date[7] != '-' {
		return false
	}
	return true
}
