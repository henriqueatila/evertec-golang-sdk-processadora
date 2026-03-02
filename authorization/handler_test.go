package authorization

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// mockHandler implements Handler for testing
type mockHandler struct {
	purchaseResponse             *PurchaseResponse
	purchaseError                error
	purchaseCancellationResponse *CancellationResponse
	purchaseCancellationError    error
	queryResponse                *QueryResponse
	queryError                   error
	withdrawalResponse           *WithdrawalResponse
	withdrawalError              error
	withdrawalQueryResponse      *WithdrawalQueryResponse
	withdrawalQueryError         error
	withdrawalCancellationResp   *CancellationResponse
	withdrawalCancellationError  error
	chargebackResponse           *ChargebackResponse
	chargebackError              error
	chargebackCancellationResp   *CancellationResponse
	chargebackCancellationError  error
	transferResponse             *TransferResponse
	transferError                error
	transferCancellationResp     *CancellationResponse
	transferCancellationError    error
	otpChannelResponse           *OTPChannelResponse
	otpChannelError              error
	verifyTransactionResponse    *VerifyTransactionResponse
	verifyTransactionError       error
	xpaysOTPResponse             *XPaysOTPResponse
	xpaysOTPError                error
	customProvisioningResponse   *types.CustomProvisioningDataResponse
	customProvisioningError      error
	statusResponse               *StatusResponse
	statusError                  error
}

func (m *mockHandler) HandlePurchase(_ context.Context, req *PurchaseRequest) (*PurchaseResponse, error) {
	return m.purchaseResponse, m.purchaseError
}

func (m *mockHandler) HandlePurchaseCancel(_ context.Context, req *PurchaseCancellationRequest) (*CancellationResponse, error) {
	return m.purchaseCancellationResponse, m.purchaseCancellationError
}

func (m *mockHandler) HandleQuery(_ context.Context, req *QueryRequest) (*QueryResponse, error) {
	return m.queryResponse, m.queryError
}

func (m *mockHandler) HandleWithdrawal(_ context.Context, req *WithdrawalRequest) (*WithdrawalResponse, error) {
	return m.withdrawalResponse, m.withdrawalError
}

func (m *mockHandler) HandleWithdrawalQuery(_ context.Context, req *WithdrawalQueryRequest) (*WithdrawalQueryResponse, error) {
	return m.withdrawalQueryResponse, m.withdrawalQueryError
}

func (m *mockHandler) HandleWithdrawalCancel(_ context.Context, req *WithdrawalCancellationRequest) (*CancellationResponse, error) {
	return m.withdrawalCancellationResp, m.withdrawalCancellationError
}

func (m *mockHandler) HandleChargeback(_ context.Context, req *ChargebackRequest) (*ChargebackResponse, error) {
	return m.chargebackResponse, m.chargebackError
}

func (m *mockHandler) HandleChargebackCancel(_ context.Context, req *ChargebackCancellationRequest) (*CancellationResponse, error) {
	return m.chargebackCancellationResp, m.chargebackCancellationError
}

func (m *mockHandler) HandleTransfer(_ context.Context, req *TransferRequest) (*TransferResponse, error) {
	return m.transferResponse, m.transferError
}

func (m *mockHandler) HandleTransferCancel(_ context.Context, req *TransferCancellationRequest) (*CancellationResponse, error) {
	return m.transferCancellationResp, m.transferCancellationError
}

func (m *mockHandler) HandleGetOTPChannel(_ context.Context, req *OTPChannelRequest) (*OTPChannelResponse, error) {
	return m.otpChannelResponse, m.otpChannelError
}

func (m *mockHandler) HandleVerifyTransaction(_ context.Context, req *VerifyTransactionRequest) (*VerifyTransactionResponse, error) {
	return m.verifyTransactionResponse, m.verifyTransactionError
}

func (m *mockHandler) HandleXPaysOTP(_ context.Context, req *XPaysOTPRequest) (*XPaysOTPResponse, error) {
	return m.xpaysOTPResponse, m.xpaysOTPError
}

func (m *mockHandler) HandleCustomProvisioningData(_ context.Context, req *CustomProvisioningDataRequest) (*types.CustomProvisioningDataResponse, error) {
	return m.customProvisioningResponse, m.customProvisioningError
}

func (m *mockHandler) HandleStatus(_ context.Context) (*StatusResponse, error) {
	return m.statusResponse, m.statusError
}

// TestServer_NewServer tests server creation
func TestServer_NewServer(t *testing.T) {
	handler := &mockHandler{}
	server := NewServer(handler)

	if server == nil {
		t.Fatal("NewServer returned nil")
	}
	if server.handler != handler {
		t.Error("Server handler not set correctly")
	}
}

// TestServer_Status tests /status endpoint
func TestServer_Status(t *testing.T) {
	handler := &mockHandler{
		statusResponse: &StatusResponse{
			Status:    "OK",
			Timestamp: "2024-01-01T00:00:00Z",
		},
	}
	server := NewServer(handler)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp StatusResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != "OK" {
		t.Errorf("Status = %s, want OK", resp.Status)
	}
}

// TestServer_Status_MethodNotAllowed tests wrong method on /status
func TestServer_Status_MethodNotAllowed(t *testing.T) {
	handler := &mockHandler{}
	server := NewServer(handler)

	req := httptest.NewRequest(http.MethodPost, "/status", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

// TestServer_Purchase_Approved tests approved purchase
func TestServer_Purchase_Approved(t *testing.T) {
	handler := &mockHandler{
		purchaseResponse: &PurchaseResponse{
			Message:         "Transação aprovada",
			Code:            0, // 00 = Approved
			AuthorizationID: 123456,
			Balance: &AuthAmount{
				Amount:       100000,
				CurrencyCode: 986,
			},
		},
	}
	server := NewServer(handler)

	purchaseReq := PurchaseRequest{
		PurchaseID: "tx123",
		AccountID:  "acc123",
		Card: &AuthCard{
			PaysmartID: "card123",
			PAN:        "************1234",
		},
		TotalAmount: &AuthAmount{
			Amount:       10000,
			CurrencyCode: 986,
		},
	}

	body, _ := json.Marshal(purchaseReq)
	req := httptest.NewRequest(http.MethodPost, "/purchases", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp PurchaseResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0 (approved)", resp.Code)
	}
	if resp.AuthorizationID != 123456 {
		t.Errorf("AuthorizationID = %d, want 123456", resp.AuthorizationID)
	}
}

// TestServer_Purchase_Denied tests denied purchase with response codes
func TestServer_Purchase_Denied(t *testing.T) {
	tests := []struct {
		name           string
		code           int
		expectedStatus int
	}{
		{"insufficient funds", 51, 400},
		{"exceeds limit", 61, 400},
		{"invalid card", 14, 404},
		{"expired card", 54, 404},
		{"fraud suspect", 59, 459},
		{"restricted card", 62, 412},
		{"invalid MCC", 57, 483},
		{"system error", 96, 503},
		{"unknown code", 99, 499},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &mockHandler{
				purchaseResponse: &PurchaseResponse{
					Message: "Declined",
					Code:    tt.code,
				},
			}
			server := NewServer(handler)

			purchaseReq := PurchaseRequest{
				PurchaseID:  "tx123",
				TotalAmount: &AuthAmount{Amount: 10000, CurrencyCode: 986},
			}

			body, _ := json.Marshal(purchaseReq)
			req := httptest.NewRequest(http.MethodPost, "/purchases", bytes.NewReader(body))
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Status code = %d, want %d for code %d", w.Code, tt.expectedStatus, tt.code)
			}
		})
	}
}

// TestServer_Purchase_InvalidJSON tests invalid JSON handling
func TestServer_Purchase_InvalidJSON(t *testing.T) {
	handler := &mockHandler{}
	server := NewServer(handler)

	req := httptest.NewRequest(http.MethodPost, "/purchases", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestServer_PurchaseCancel tests purchase cancellation
func TestServer_PurchaseCancel(t *testing.T) {
	handler := &mockHandler{
		purchaseCancellationResponse: &CancellationResponse{
			Message:         "Cancelamento aprovado",
			Code:            0,
			AuthorizationID: 654321,
			Balance: &AuthAmount{
				Amount:       110000,
				CurrencyCode: 986,
			},
		},
	}
	server := NewServer(handler)

	cancelReq := PurchaseCancellationRequest{
		CancellationID:     "cancel123",
		OriginalPurchaseID: "tx123",
		AccountID:          "acc123",
		OriginalAmount:     &AuthAmount{Amount: 10000, CurrencyCode: 986},
	}

	body, _ := json.Marshal(cancelReq)
	req := httptest.NewRequest(http.MethodPost, "/purchases/cancel", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp CancellationResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0 (approved)", resp.Code)
	}
}

// TestServer_Query tests balance query
func TestServer_Query(t *testing.T) {
	handler := &mockHandler{
		queryResponse: &QueryResponse{
			Message: "Consulta aprovada",
			Code:    0,
			Balance: &AuthAmount{Amount: 100000, CurrencyCode: 986},
		},
	}
	server := NewServer(handler)

	queryReq := QueryRequest{
		QueryID:   "query123",
		AccountID: "acc123",
		Card:      &AuthCard{PaysmartID: "card123"},
	}

	body, _ := json.Marshal(queryReq)
	req := httptest.NewRequest(http.MethodPost, "/queries", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}
}

// TestServer_Withdrawal tests withdrawal authorization
func TestServer_Withdrawal(t *testing.T) {
	handler := &mockHandler{
		withdrawalResponse: &WithdrawalResponse{
			Message:         "Saque aprovado",
			Code:            0,
			AuthorizationID: 654321,
			Balance:         &AuthAmount{Amount: 90000, CurrencyCode: 986},
		},
	}
	server := NewServer(handler)

	withdrawalReq := WithdrawalRequest{
		WithdrawalID: "tx456",
		AccountID:    "acc123",
		Card:         &AuthCard{PaysmartID: "card123"},
		TotalAmount:  &AuthAmount{Amount: 10000, CurrencyCode: 986},
	}

	body, _ := json.Marshal(withdrawalReq)
	req := httptest.NewRequest(http.MethodPost, "/withdrawals", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp WithdrawalResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0", resp.Code)
	}
}

// TestServer_WithdrawalQuery tests withdrawal query
func TestServer_WithdrawalQuery(t *testing.T) {
	handler := &mockHandler{
		withdrawalQueryResponse: &WithdrawalQueryResponse{
			Message:             "Consulta aprovada",
			Code:                0,
			MaxWithdrawalAmount: &AuthAmount{Amount: 100000, CurrencyCode: 986},
			DailyLimit:          &AuthAmount{Amount: 500000, CurrencyCode: 986},
			RemainingDaily:      &AuthAmount{Amount: 400000, CurrencyCode: 986},
		},
	}
	server := NewServer(handler)

	queryReq := WithdrawalQueryRequest{
		QueryID:   "query456",
		AccountID: "acc123",
		Card:      &AuthCard{PaysmartID: "card123"},
	}

	body, _ := json.Marshal(queryReq)
	req := httptest.NewRequest(http.MethodPost, "/withdrawalQueries", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}
}

// TestServer_Chargeback tests chargeback handling
func TestServer_Chargeback(t *testing.T) {
	handler := &mockHandler{
		chargebackResponse: &ChargebackResponse{
			Message:         "Chargeback processado",
			Code:            0,
			AuthorizationID: 789012,
		},
	}
	server := NewServer(handler)

	chargebackReq := ChargebackRequest{
		ChargebackID:       "cb123",
		OriginalPurchaseID: "tx123",
		AccountID:          "acc123",
		ChargebackAmount:   &AuthAmount{Amount: 10000, CurrencyCode: 986},
		Reason:             "FRAUD",
		DisputeID:          "disp123",
	}

	body, _ := json.Marshal(chargebackReq)
	req := httptest.NewRequest(http.MethodPost, "/chargebacks", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp ChargebackResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Code != 0 {
		t.Errorf("Code = %d, want 0", resp.Code)
	}
}

// TestServer_Transfer tests transfer authorization
func TestServer_Transfer(t *testing.T) {
	handler := &mockHandler{
		transferResponse: &TransferResponse{
			Message:         "Transferência aprovada",
			Code:            0,
			AuthorizationID: 789012,
			Balance:         &AuthAmount{Amount: 80000, CurrencyCode: 986},
		},
	}
	server := NewServer(handler)

	transferReq := TransferRequest{
		TransferID:  "tx789",
		AccountID:   "acc123",
		SourceCard:  &AuthCard{PaysmartID: "card123"},
		TotalAmount: &AuthAmount{Amount: 20000, CurrencyCode: 986},
	}

	body, _ := json.Marshal(transferReq)
	req := httptest.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}
}

// TestServer_GetOTPChannel tests 3DS OTP channel request
func TestServer_GetOTPChannel(t *testing.T) {
	handler := &mockHandler{
		otpChannelResponse: &OTPChannelResponse{
			Channel:           "SMS",
			MaskedDestination: "+55 ** *****-1234",
		},
	}
	server := NewServer(handler)

	otpReq := OTPChannelRequest{
		CardID:        "card123",
		TransactionID: "tx123",
	}

	body, _ := json.Marshal(otpReq)
	req := httptest.NewRequest(http.MethodPost, "/acs/getOTPChannel", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp OTPChannelResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Channel != "SMS" {
		t.Errorf("Channel = %s, want SMS", resp.Channel)
	}
}

// TestServer_VerifyTransaction tests 3DS verification
func TestServer_VerifyTransaction(t *testing.T) {
	handler := &mockHandler{
		verifyTransactionResponse: &VerifyTransactionResponse{
			Verified: true,
			Message:  "OTP verified successfully",
			Code:     0,
		},
	}
	server := NewServer(handler)

	verifyReq := VerifyTransactionRequest{
		CardID:        "card123",
		TransactionID: "tx123",
		Amount:        &AuthAmount{Amount: 10000, CurrencyCode: 986},
		OTP:           "123456",
	}

	body, _ := json.Marshal(verifyReq)
	req := httptest.NewRequest(http.MethodPost, "/acs/verifyTransaction", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp VerifyTransactionResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if !resp.Verified {
		t.Error("Verified = false, want true")
	}
}

// TestServer_XPaysOTP tests xPays OTP request
func TestServer_XPaysOTP(t *testing.T) {
	handler := &mockHandler{
		xpaysOTPResponse: &XPaysOTPResponse{
			Status:            "SENT",
			MaskedDestination: "j***@example.com",
			ExpiresAt:         "2024-01-01T00:05:00Z",
			Code:              0,
		},
	}
	server := NewServer(handler)

	xpaysReq := XPaysOTPRequest{
		CardID:  "card123",
		Wallet:  types.WalletApplePay,
		Channel: "EMAIL",
	}

	body, _ := json.Marshal(xpaysReq)
	req := httptest.NewRequest(http.MethodPost, "/xpays/otp", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	var resp XPaysOTPResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Status != "SENT" {
		t.Errorf("Status = %s, want SENT", resp.Status)
	}
}

// TestServer_CustomProvisioningData tests xPays custom provisioning
func TestServer_CustomProvisioningData(t *testing.T) {
	handler := &mockHandler{
		customProvisioningResponse: &types.CustomProvisioningDataResponse{
			CardArt:          "https://example.com/card.png",
			IssuerName:       "Test Bank",
			ShortDescription: "Test Card",
		},
	}
	server := NewServer(handler)

	provReq := CustomProvisioningDataRequest{
		CardID: "card123",
		Wallet: types.WalletGooglePay,
	}

	body, _ := json.Marshal(provReq)
	req := httptest.NewRequest(http.MethodPost, "/xpays/customProvisioningData", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}
}

// TestServer_NotFound tests unknown endpoint
func TestServer_NotFound(t *testing.T) {
	handler := &mockHandler{}
	server := NewServer(handler)

	req := httptest.NewRequest(http.MethodPost, "/unknown", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// TestServer_AllCancelEndpoints tests all cancellation endpoints
func TestServer_AllCancelEndpoints(t *testing.T) {
	handler := &mockHandler{
		purchaseCancellationResponse: &CancellationResponse{
			Message: "Cancelled",
			Code:    0,
		},
		withdrawalCancellationResp: &CancellationResponse{
			Message: "Cancelled",
			Code:    0,
		},
		chargebackCancellationResp: &CancellationResponse{
			Message: "Cancelled",
			Code:    0,
		},
		transferCancellationResp: &CancellationResponse{
			Message: "Cancelled",
			Code:    0,
		},
	}
	server := NewServer(handler)

	t.Run("/purchases/cancel", func(t *testing.T) {
		cancelReq := PurchaseCancellationRequest{
			CancellationID:     "cancel1",
			OriginalPurchaseID: "tx123",
		}
		body, _ := json.Marshal(cancelReq)
		req := httptest.NewRequest(http.MethodPost, "/purchases/cancel", bytes.NewReader(body))
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("/withdrawals/cancel", func(t *testing.T) {
		cancelReq := WithdrawalCancellationRequest{
			CancellationID:       "cancel2",
			OriginalWithdrawalID: "tx456",
		}
		body, _ := json.Marshal(cancelReq)
		req := httptest.NewRequest(http.MethodPost, "/withdrawals/cancel", bytes.NewReader(body))
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("/chargebacks/cancel", func(t *testing.T) {
		cancelReq := ChargebackCancellationRequest{
			CancellationID:       "cancel3",
			OriginalChargebackID: "cb123",
		}
		body, _ := json.Marshal(cancelReq)
		req := httptest.NewRequest(http.MethodPost, "/chargebacks/cancel", bytes.NewReader(body))
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("/transfers/cancel", func(t *testing.T) {
		cancelReq := TransferCancellationRequest{
			CancellationID:     "cancel4",
			OriginalTransferID: "tx789",
		}
		body, _ := json.Marshal(cancelReq)
		req := httptest.NewRequest(http.MethodPost, "/transfers/cancel", bytes.NewReader(body))
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
		}
	})
}

// TestServer_ContentTypeHeader tests Content-Type header is set
func TestServer_ContentTypeHeader(t *testing.T) {
	handler := &mockHandler{
		statusResponse: &StatusResponse{Status: "OK"},
	}
	server := NewServer(handler)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %s, want application/json", contentType)
	}
}

// TestMapResponseCodeToHTTP tests response code mapping
func TestMapResponseCodeToHTTP(t *testing.T) {
	tests := []struct {
		code int
		want int
	}{
		{0, 200},  // Approved
		{51, 400}, // Insufficient funds
		{61, 400}, // Exceeds limit
		{14, 404}, // Invalid card
		{54, 404}, // Expired card
		{59, 459}, // Fraud suspect
		{62, 412}, // Restricted card
		{57, 483}, // Invalid MCC
		{96, 503}, // System error
		{99, 499}, // Unknown denial code
	}

	for _, tt := range tests {
		got := mapResponseCodeToHTTP(tt.code)
		if got != tt.want {
			t.Errorf("mapResponseCodeToHTTP(%d) = %d, want %d", tt.code, got, tt.want)
		}
	}
}

// TestResponseCode_IsApproved tests IsApproved method
func TestResponseCode_IsApproved(t *testing.T) {
	tests := []struct {
		code ResponseCode
		want bool
	}{
		{ResponseCodeApproved, true},
		{ResponseCodeInsufficientFunds, false},
		{ResponseCodeFraudSuspect, false},
		{ResponseCodeSystemFailure, false},
	}

	for _, tt := range tests {
		if got := tt.code.IsApproved(); got != tt.want {
			t.Errorf("ResponseCode(%d).IsApproved() = %v, want %v", tt.code, got, tt.want)
		}
	}
}

// TestResponseCode_String tests String method
func TestResponseCode_String(t *testing.T) {
	if s := ResponseCodeApproved.String(); s != "00" {
		t.Errorf("ResponseCodeApproved.String() = %s, want 00", s)
	}
	if s := ResponseCodeInsufficientFunds.String(); s != "" {
		t.Errorf("ResponseCodeInsufficientFunds.String() = %s, want empty", s)
	}
}

// TestServer_HandlerErrors tests error handling when handlers return errors
func TestServer_HandlerErrors(t *testing.T) {
	t.Run("PurchaseError", func(t *testing.T) {
		handler := &mockHandler{
			purchaseError: fmt.Errorf("internal error"),
		}
		server := NewServer(handler)

		purchaseReq := PurchaseRequest{PurchaseID: "tx123"}
		body, _ := json.Marshal(purchaseReq)
		req := httptest.NewRequest(http.MethodPost, "/purchases", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("QueryError", func(t *testing.T) {
		handler := &mockHandler{
			queryError: fmt.Errorf("query error"),
		}
		server := NewServer(handler)

		queryReq := QueryRequest{QueryID: "query123"}
		body, _ := json.Marshal(queryReq)
		req := httptest.NewRequest(http.MethodPost, "/queries", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("WithdrawalError", func(t *testing.T) {
		handler := &mockHandler{
			withdrawalError: fmt.Errorf("withdrawal error"),
		}
		server := NewServer(handler)

		withdrawalReq := WithdrawalRequest{WithdrawalID: "tx123"}
		body, _ := json.Marshal(withdrawalReq)
		req := httptest.NewRequest(http.MethodPost, "/withdrawals", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("WithdrawalQueryError", func(t *testing.T) {
		handler := &mockHandler{
			withdrawalQueryError: fmt.Errorf("withdrawal query error"),
		}
		server := NewServer(handler)

		queryReq := WithdrawalQueryRequest{}
		body, _ := json.Marshal(queryReq)
		req := httptest.NewRequest(http.MethodPost, "/withdrawalQueries", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("ChargebackError", func(t *testing.T) {
		handler := &mockHandler{
			chargebackError: fmt.Errorf("chargeback error"),
		}
		server := NewServer(handler)

		chargebackReq := ChargebackRequest{ChargebackID: "cb123"}
		body, _ := json.Marshal(chargebackReq)
		req := httptest.NewRequest(http.MethodPost, "/chargebacks", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("TransferError", func(t *testing.T) {
		handler := &mockHandler{
			transferError: fmt.Errorf("transfer error"),
		}
		server := NewServer(handler)

		transferReq := TransferRequest{TransferID: "tx123"}
		body, _ := json.Marshal(transferReq)
		req := httptest.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("OTPChannelError", func(t *testing.T) {
		handler := &mockHandler{
			otpChannelError: fmt.Errorf("otp channel error"),
		}
		server := NewServer(handler)

		otpReq := OTPChannelRequest{CardID: "card123"}
		body, _ := json.Marshal(otpReq)
		req := httptest.NewRequest(http.MethodPost, "/acs/getOTPChannel", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("VerifyTransactionError", func(t *testing.T) {
		handler := &mockHandler{
			verifyTransactionError: fmt.Errorf("verify error"),
		}
		server := NewServer(handler)

		verifyReq := VerifyTransactionRequest{CardID: "card123"}
		body, _ := json.Marshal(verifyReq)
		req := httptest.NewRequest(http.MethodPost, "/acs/verifyTransaction", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("XPaysOTPError", func(t *testing.T) {
		handler := &mockHandler{
			xpaysOTPError: fmt.Errorf("xpays otp error"),
		}
		server := NewServer(handler)

		xpaysReq := XPaysOTPRequest{CardID: "card123"}
		body, _ := json.Marshal(xpaysReq)
		req := httptest.NewRequest(http.MethodPost, "/xpays/otp", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("CustomProvisioningError", func(t *testing.T) {
		handler := &mockHandler{
			customProvisioningError: fmt.Errorf("provisioning error"),
		}
		server := NewServer(handler)

		provReq := CustomProvisioningDataRequest{CardID: "card123"}
		body, _ := json.Marshal(provReq)
		req := httptest.NewRequest(http.MethodPost, "/xpays/customProvisioningData", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("StatusError", func(t *testing.T) {
		handler := &mockHandler{
			statusError: fmt.Errorf("status error"),
		}
		server := NewServer(handler)

		req := httptest.NewRequest(http.MethodGet, "/status", nil)
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		// Status errors return 503 (Service Unavailable)
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusServiceUnavailable)
		}
	})

	t.Run("PurchaseCancellationError", func(t *testing.T) {
		handler := &mockHandler{
			purchaseCancellationError: fmt.Errorf("cancellation error"),
		}
		server := NewServer(handler)

		cancelReq := PurchaseCancellationRequest{OriginalPurchaseID: "tx123"}
		body, _ := json.Marshal(cancelReq)
		req := httptest.NewRequest(http.MethodPost, "/purchases/cancel", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})
}

// TestServer_InvalidJSON_AllEndpoints tests invalid JSON for all POST endpoints
func TestServer_InvalidJSON_AllEndpoints(t *testing.T) {
	endpoints := []string{
		"/purchases",
		"/purchases/cancel",
		"/queries",
		"/withdrawals",
		"/withdrawalQueries",
		"/withdrawals/cancel",
		"/chargebacks",
		"/chargebacks/cancel",
		"/transfers",
		"/transfers/cancel",
		"/acs/getOTPChannel",
		"/acs/verifyTransaction",
		"/xpays/otp",
		"/xpays/customProvisioningData",
	}

	handler := &mockHandler{}
	server := NewServer(handler)

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewReader([]byte("invalid json")))
			w := httptest.NewRecorder()

			server.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("%s: Status code = %d, want %d", endpoint, w.Code, http.StatusBadRequest)
			}
		})
	}
}

// TestServer_Withdrawal_Denied tests denied withdrawal
func TestServer_Withdrawal_Denied(t *testing.T) {
	handler := &mockHandler{
		withdrawalResponse: &WithdrawalResponse{
			Message: "Saldo insuficiente",
			Code:    51,
		},
	}
	server := NewServer(handler)

	withdrawalReq := WithdrawalRequest{WithdrawalID: "tx123"}
	body, _ := json.Marshal(withdrawalReq)
	req := httptest.NewRequest(http.MethodPost, "/withdrawals", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Status code = %d, want 400", w.Code)
	}
}

// TestServer_Transfer_Denied tests denied transfer
func TestServer_Transfer_Denied(t *testing.T) {
	handler := &mockHandler{
		transferResponse: &TransferResponse{
			Message: "Suspeita de fraude",
			Code:    59,
		},
	}
	server := NewServer(handler)

	transferReq := TransferRequest{TransferID: "tx123"}
	body, _ := json.Marshal(transferReq)
	req := httptest.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != 459 {
		t.Errorf("Status code = %d, want 459", w.Code)
	}
}

// TestAuthTypes tests the new authorization-specific types
func TestAuthTypes(t *testing.T) {
	t.Run("AuthAmount", func(t *testing.T) {
		amt := &AuthAmount{
			Amount:       10000,
			CurrencyCode: 986,
		}
		if amt.Amount != 10000 {
			t.Errorf("Amount = %d, want 10000", amt.Amount)
		}
		if amt.CurrencyCode != 986 {
			t.Errorf("CurrencyCode = %d, want 986", amt.CurrencyCode)
		}
	})

	t.Run("AuthCard", func(t *testing.T) {
		card := &AuthCard{
			PaysmartID: "card-uuid",
			IssuerID:   123,
			PAN:        "************1234",
			PANSeq:     "001",
		}
		if card.PaysmartID != "card-uuid" {
			t.Errorf("PaysmartID = %s, want card-uuid", card.PaysmartID)
		}
	})

	t.Run("ProcessingCode", func(t *testing.T) {
		pc := &ProcessingCode{
			TipoTransacao:          "PURCHASE",
			SourceAccountType:      "CREDIT_CARD_ACCOUNT",
			DestinationAccountType: "NOT_APPLICABLE",
		}
		if pc.TipoTransacao != "PURCHASE" {
			t.Errorf("TipoTransacao = %s, want PURCHASE", pc.TipoTransacao)
		}
	})

	t.Run("CancellationReason", func(t *testing.T) {
		reason := &CancellationReason{
			Code:        "001",
			Description: "Timeout",
		}
		if reason.Code != "001" {
			t.Errorf("Code = %s, want 001", reason.Code)
		}
	})

	t.Run("FraudData", func(t *testing.T) {
		fd := &FraudData{
			CreditorFraudScore:          75,
			FraudDecisionRecommendation: "REVIEW",
		}
		if fd.CreditorFraudScore != 75 {
			t.Errorf("CreditorFraudScore = %d, want 75", fd.CreditorFraudScore)
		}
	})
}

// TestServer_WithOptions tests server creation with options
func TestServer_WithOptions(t *testing.T) {
	handler := &mockHandler{
		statusResponse: &StatusResponse{Status: "healthy"},
	}

	// Create custom hook for testing
	hookCalled := false
	testHook := NewHook(
		func(ctx context.Context, req *RequestInfo) context.Context { return ctx },
		func(ctx context.Context, req *RequestInfo, resp *ResponseInfo) { hookCalled = true },
	)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	server := NewServer(handler,
		WithHooks(testHook),
		WithLogger(logger),
		WithMaxBodySize(10*1024*1024),
		WithPanicRecovery(),
	)

	if server == nil {
		t.Fatal("NewServer returned nil")
	}
	if server.maxBodySize != 10*1024*1024 {
		t.Errorf("maxBodySize = %d, want 10MB", server.maxBodySize)
	}
	if !server.panicRecovery {
		t.Error("panicRecovery should be enabled")
	}

	// Test that hooks work
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if !hookCalled {
		t.Error("Hook was not called")
	}
}

// TestServer_MaxBodySize tests request body size limit
func TestServer_MaxBodySize(t *testing.T) {
	handler := &mockHandler{
		purchaseResponse: &PurchaseResponse{Code: 0},
	}
	server := NewServer(handler, WithMaxBodySize(50)) // 50 bytes limit

	// Create a request body larger than limit
	largeBody := make([]byte, 100)
	for i := range largeBody {
		largeBody[i] = 'a'
	}

	req := httptest.NewRequest(http.MethodPost, "/purchases", bytes.NewReader(largeBody))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	// Should fail due to body size limit
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestServer_PanicRecovery tests panic recovery middleware
func TestServer_PanicRecovery(t *testing.T) {
	// Handler that panics
	panicHandler := &mockHandler{}
	panicHandler.purchaseResponse = nil
	panicHandler.purchaseError = nil

	// Swap to a handler that will panic
	server := NewServer(&panicingHandler{}, WithPanicRecovery())

	purchaseReq := PurchaseRequest{PurchaseID: "tx123"}
	body, _ := json.Marshal(purchaseReq)
	req := httptest.NewRequest(http.MethodPost, "/purchases", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Should not panic, should return 500
	server.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// panicingHandler is a handler that panics for testing
type panicingHandler struct{}

func (p *panicingHandler) HandlePurchase(_ context.Context, req *PurchaseRequest) (*PurchaseResponse, error) {
	panic("test panic in purchase handler")
}
func (p *panicingHandler) HandlePurchaseCancel(_ context.Context, req *PurchaseCancellationRequest) (*CancellationResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleQuery(_ context.Context, req *QueryRequest) (*QueryResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleWithdrawal(_ context.Context, req *WithdrawalRequest) (*WithdrawalResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleWithdrawalQuery(_ context.Context, req *WithdrawalQueryRequest) (*WithdrawalQueryResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleWithdrawalCancel(_ context.Context, req *WithdrawalCancellationRequest) (*CancellationResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleChargeback(_ context.Context, req *ChargebackRequest) (*ChargebackResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleChargebackCancel(_ context.Context, req *ChargebackCancellationRequest) (*CancellationResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleTransfer(_ context.Context, req *TransferRequest) (*TransferResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleTransferCancel(_ context.Context, req *TransferCancellationRequest) (*CancellationResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleGetOTPChannel(_ context.Context, req *OTPChannelRequest) (*OTPChannelResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleVerifyTransaction(_ context.Context, req *VerifyTransactionRequest) (*VerifyTransactionResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleXPaysOTP(_ context.Context, req *XPaysOTPRequest) (*XPaysOTPResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleCustomProvisioningData(_ context.Context, req *CustomProvisioningDataRequest) (*types.CustomProvisioningDataResponse, error) {
	return nil, nil
}
func (p *panicingHandler) HandleStatus(_ context.Context) (*StatusResponse, error) {
	return nil, nil
}

// TestServer_WithLogger tests logger option
func TestServer_WithLogger(t *testing.T) {
	handler := &mockHandler{
		purchaseResponse: &PurchaseResponse{Code: 0, Message: "Approved"},
	}

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	server := NewServer(handler, WithLogger(logger))

	purchaseReq := PurchaseRequest{PurchaseID: "tx123"}
	body, _ := json.Marshal(purchaseReq)
	req := httptest.NewRequest(http.MethodPost, "/purchases", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}
}
