package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

func TestClient_CreateInclusiveTransaction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.InclusiveTransactionRequest
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "success: create inclusive transaction",
			request: &types.InclusiveTransactionRequest{
				AccountID:     "acc123",
				TransactionID: "tx123",
				Code:          "TE10",
				ReasonCode:    "ADJUSTMENT",
				Text:          "Credit adjustment",
				Partial:       false,
				Amount:        types.Amount{Amount: 10000, CurrencyCode: 986},
			},
			mockResponse: map[string]interface{}{
				"inclusiveTransactionId": "inc123",
				"status":                 "COMPLETED",
				"code":                   "TE10",
			},
			mockStatus: http.StatusCreated,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]interface{}) {
				AssertEqual(t, "TE10", body["code"], "code")
				AssertEqual(t, "acc123", body["accountId"], "accountId")
			},
		},
		{
			name: "success: create partial inclusive transaction",
			request: &types.InclusiveTransactionRequest{
				AccountID:     "acc456",
				TransactionID: "tx456",
				Code:          "TE20",
				ReasonCode:    "PARTIAL_REFUND",
				Text:          "Partial refund",
				Partial:       true,
				Amount:        types.Amount{Amount: 5000, CurrencyCode: 986},
			},
			mockResponse: map[string]interface{}{
				"inclusiveTransactionId": "inc456",
				"status":                 "COMPLETED",
				"code":                   "TE20",
			},
			mockStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "error: account not found",
			request: &types.InclusiveTransactionRequest{
				AccountID:     "acc-nonexistent",
				TransactionID: "tx789",
				Code:          "TE10",
				ReasonCode:    "TEST",
				Text:          "Test",
				Amount:        types.Amount{Amount: 1000, CurrencyCode: 986},
			},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
		{
			name: "error: invalid code",
			request: &types.InclusiveTransactionRequest{
				AccountID:     "acc123",
				TransactionID: "tx999",
				Code:          "INVALID",
				ReasonCode:    "TEST",
				Text:          "Test",
				Amount:        types.Amount{Amount: 1000, CurrencyCode: 986},
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
				Path:     "/inclusives",
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
			resp, err := client.CreateInclusiveTransaction(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "CreateInclusiveTransaction")
			} else {
				AssertNoError(t, err, "CreateInclusiveTransaction")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_GetInclusiveTransaction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                   string
		inclusiveTransactionID string
		mockResponse           interface{}
		mockStatus             int
		wantErr                bool
	}{
		{
			name:                   "success: get inclusive transaction",
			inclusiveTransactionID: "inc123",
			mockResponse: map[string]interface{}{
				"inclusiveTransactionId": "inc123",
				"status":                 "COMPLETED",
				"code":                   "TE10",
				"amount":                 map[string]interface{}{"amount": 10000, "currencyCode": 986},
				"accountId":              "acc123",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:                   "success: get pending inclusive transaction",
			inclusiveTransactionID: "inc456",
			mockResponse: map[string]interface{}{
				"inclusiveTransactionId": "inc456",
				"status":                 "PENDING",
				"code":                   "TE20",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:                   "error: inclusive transaction not found",
			inclusiveTransactionID: "inc-nonexistent",
			mockResponse:           nil,
			mockStatus:             http.StatusNotFound,
			wantErr:                true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/inclusives/" + tt.inclusiveTransactionID,
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetInclusiveTransaction(context.Background(), tt.inclusiveTransactionID)

			if tt.wantErr {
				AssertError(t, err, "GetInclusiveTransaction")
			} else {
				AssertNoError(t, err, "GetInclusiveTransaction")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}
