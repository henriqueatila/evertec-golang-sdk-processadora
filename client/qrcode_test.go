package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

func TestClient_ParseQRCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.ParseQrCodeParams
		mockResponse any
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, body map[string]any)
	}{
		{
			name: "success: parse QR code",
			request: &types.ParseQrCodeParams{
				QRCode: "00020126580014br.gov.bcb.pix0136123e4567-e12b-12d1-a456-426614174000",
			},
			mockResponse: map[string]any{
				"type":     "PIX",
				"merchant": map[string]any{"name": "Test Merchant"},
				"amount":   map[string]any{"amount": 10000, "currencyCode": 986},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]any) {
				AssertNotNil(t, body["qrCode"], "qrCode")
			},
		},
		{
			name: "success: parse static QR code",
			request: &types.ParseQrCodeParams{
				QRCode: "static-qr-code-data",
			},
			mockResponse: map[string]any{
				"type":     "STATIC",
				"merchant": map[string]any{"name": "Static Merchant"},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "error: invalid QR code",
			request: &types.ParseQrCodeParams{
				QRCode: "invalid-qr-code",
			},
			mockResponse: nil,
			mockStatus:   http.StatusBadRequest,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/qrcode/parse",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]any
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.ParseQRCode(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "ParseQRCode")
			} else {
				AssertNoError(t, err, "ParseQRCode")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_SimplePayment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.SimplePaymentCardParams
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name: "success: simple payment",
			request: &types.SimplePaymentCardParams{
				VCardID: "vcard123",
				QRCode:  "00020126580014br.gov.bcb.pix0136123e4567-e12b-12d1-a456-426614174000",
			},
			mockResponse: map[string]any{
				"transactionId":     "tx123",
				"status":            "APPROVED",
				"authorizationCode": "AUTH123",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: payment with different card",
			request: &types.SimplePaymentCardParams{
				VCardID: "vcard456",
				QRCode:  "static-qr-code-data",
			},
			mockResponse: map[string]any{
				"transactionId":     "tx456",
				"status":            "APPROVED",
				"authorizationCode": "AUTH456",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "error: insufficient funds",
			request: &types.SimplePaymentCardParams{
				VCardID: "vcard-low-balance",
				QRCode:  "qr-code-high-amount",
			},
			mockResponse: nil,
			mockStatus:   http.StatusPaymentRequired,
			wantErr:      true,
		},
		{
			name: "error: card blocked",
			request: &types.SimplePaymentCardParams{
				VCardID: "vcard-blocked",
				QRCode:  "qr-code-data",
			},
			mockResponse: nil,
			mockStatus:   http.StatusForbidden,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/qrcode/simplePayment",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.SimplePayment(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "SimplePayment")
			} else {
				AssertNoError(t, err, "SimplePayment")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_TransactionCallback(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		transactionID string
		request       *types.TransactionCallbackRequest
		mockResponse  any
		mockStatus    int
		wantErr       bool
	}{
		{
			name:          "success: transaction callback",
			transactionID: "tx123",
			request: &types.TransactionCallbackRequest{
				Status:       "COMPLETED",
				ResponseCode: "00",
			},
			mockResponse: map[string]any{
				"transactionId": "tx123",
				"status":        "COMPLETED",
				"processed":     true,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:          "success: callback with failure",
			transactionID: "tx456",
			request: &types.TransactionCallbackRequest{
				Status:       "FAILED",
				ResponseCode: "51",
				Message:      "Insufficient funds",
			},
			mockResponse: map[string]any{
				"transactionId": "tx456",
				"status":        "FAILED",
				"processed":     true,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:          "error: transaction not found",
			transactionID: "tx-nonexistent",
			request: &types.TransactionCallbackRequest{
				Status: "COMPLETED",
			},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/callbacks/transactions/" + tt.transactionID,
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.TransactionCallback(context.Background(), tt.transactionID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "TransactionCallback")
			} else {
				AssertNoError(t, err, "TransactionCallback")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}
