package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

func TestClient_CreateCobranded(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.MerchantVanRequest
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "success: create cobranded merchant",
			request: &types.MerchantVanRequest{
				AcquirerID:      1,
				Document:        "12345678000199",
				DocumentType:    "CNPJ",
				LegalName:       "Test Merchant LTDA",
				FantasyName:     "Test Merchant",
				MerchantVanCode: "VAN001",
			},
			mockResponse: map[string]interface{}{
				"document":  "12345678000199",
				"legalName": "Test Merchant LTDA",
				"status":    "ACTIVE",
			},
			mockStatus: http.StatusCreated,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]interface{}) {
				AssertEqual(t, "12345678000199", body["document"], "document")
				AssertEqual(t, "Test Merchant LTDA", body["legalName"], "legalName")
			},
		},
		{
			name: "success: create cobranded with products",
			request: &types.MerchantVanRequest{
				AcquirerID:      2,
				Document:        "98765432000199",
				DocumentType:    "CNPJ",
				LegalName:       "Another Merchant SA",
				MerchantVanCode: "VAN002",
			},
			mockResponse: map[string]interface{}{
				"document":  "98765432000199",
				"legalName": "Another Merchant SA",
				"status":    "ACTIVE",
			},
			mockStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "error: duplicate document",
			request: &types.MerchantVanRequest{
				AcquirerID:      1,
				Document:        "12345678000199",
				DocumentType:    "CNPJ",
				LegalName:       "Duplicate Merchant",
				MerchantVanCode: "VAN003",
			},
			mockResponse: nil,
			mockStatus:   http.StatusConflict,
			wantErr:      true,
		},
		{
			name: "error: invalid document",
			request: &types.MerchantVanRequest{
				AcquirerID:      1,
				Document:        "invalid",
				DocumentType:    "CNPJ",
				LegalName:       "Invalid Merchant",
				MerchantVanCode: "VAN004",
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
				Path:     "/cobranded",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]interface{}
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.CreateCobranded(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "CreateCobranded")
			} else {
				AssertNoError(t, err, "CreateCobranded")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_DeleteCobranded(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		document   string
		mockStatus int
		wantErr    bool
	}{
		{
			name:       "success: delete cobranded merchant",
			document:   "12345678000199",
			mockStatus: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "success: delete another cobranded",
			document:   "98765432000199",
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "error: merchant not found",
			document:   "00000000000000",
			mockStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodDelete,
				Path:     "/cobranded/" + tt.document,
				Status:   tt.mockStatus,
				Response: nil,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			err := client.DeleteCobranded(context.Background(), tt.document)

			if tt.wantErr {
				AssertError(t, err, "DeleteCobranded")
			} else {
				AssertNoError(t, err, "DeleteCobranded")
			}
		})
	}
}
