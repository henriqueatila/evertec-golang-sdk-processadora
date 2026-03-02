package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

func TestClient_GetTransaction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		transactionID string
		mockResponse  any
		mockStatus    int
		wantErr       bool
	}{
		{
			name:          "success: get transaction",
			transactionID: "tx123",
			mockResponse: map[string]any{
				"transactionId": "tx123",
				"status":        "COMPLETED",
				"type":          "PURCHASE",
				"amount":        map[string]any{"amount": 10000, "currencyCode": 986},
				"merchant": map[string]any{
					"name": "Test Merchant",
					"mcc":  "5411",
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:          "success: get pending transaction",
			transactionID: "tx456",
			mockResponse: map[string]any{
				"transactionId": "tx456",
				"status":        "PENDING",
				"type":          "AUTHORIZATION",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:          "error: transaction not found",
			transactionID: "tx-nonexistent",
			mockResponse:  nil,
			mockStatus:    http.StatusNotFound,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/transactions/" + tt.transactionID,
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetTransaction(context.Background(), tt.transactionID)

			if tt.wantErr {
				AssertError(t, err, "GetTransaction")
			} else {
				AssertNoError(t, err, "GetTransaction")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_ListAllTransactions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.ListAllTransactionsRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:    "success: list all transactions without filters",
			request: nil,
			mockResponse: map[string]any{
				"hasMore": true,
				"transactions": []map[string]any{
					{"transactionId": "tx001", "transactionStatus": "COMPLETED"},
					{"transactionId": "tx002", "transactionStatus": "PENDING"},
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: list with limit filter",
			request: &types.ListAllTransactionsRequest{
				Limit: 10,
			},
			mockResponse: map[string]any{
				"hasMore": false,
				"transactions": []map[string]any{
					{"transactionId": "tx001"},
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: list with date filters",
			request: &types.ListAllTransactionsRequest{
				BeginningDate: "2024-01-01",
				EndingDate:    "2024-12-31",
			},
			mockResponse: map[string]any{
				"hasMore":      false,
				"transactions": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: list with pagination cursors",
			request: &types.ListAllTransactionsRequest{
				StartingAfter: "tx100",
				EndingBefore:  "tx200",
			},
			mockResponse: map[string]any{
				"hasMore":      false,
				"transactions": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: list with all filters",
			request: &types.ListAllTransactionsRequest{
				Limit:                 10,
				StartingAfter:         "tx100",
				BeginningDate:         "2024-01-01",
				EndingDate:            "2024-12-31",
				TransactionType:       "PURCHASE",
				TransactionStatus:     "COMPLETED",
				TransactionApproved:   boolPtr(true),
				TransactionDenialCode: "00",
				MinimumAmount:         1000,
				MaxAmount:             100000,
				TransactionEntryMode:  "CHIP",
			},
			mockResponse: map[string]any{
				"hasMore": false,
				"transactions": []map[string]any{
					{"transactionId": "tx001"},
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "error: server error",
			request:      nil,
			mockResponse: nil,
			mockStatus:   http.StatusInternalServerError,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/transactions",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.ListAllTransactions(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "ListAllTransactions")
			} else {
				AssertNoError(t, err, "ListAllTransactions")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_ListAccountTransactions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		accountID    string
		request      *types.ListTransactionsRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:      "success: list account transactions",
			accountID: "acc123",
			request:   nil,
			mockResponse: map[string]any{
				"hasMore": true,
				"transactions": []map[string]any{
					{"transactionId": "tx001"},
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "success: list with filters",
			accountID: "acc123",
			request: &types.ListTransactionsRequest{
				BeginningDate:     "2024-01-01",
				EndingDate:        "2024-12-31",
				TransactionType:   "PURCHASE",
				TransactionStatus: "COMPLETED",
				Limit:             20,
				StartingAfter:     "tx100",
			},
			mockResponse: map[string]any{
				"hasMore":      false,
				"transactions": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "error: account not found",
			accountID:    "nonexistent",
			request:      nil,
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
				Method:   http.MethodGet,
				Path:     "/accounts/" + tt.accountID + "/transactions",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.ListAccountTransactions(context.Background(), tt.accountID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "ListAccountTransactions")
			} else {
				AssertNoError(t, err, "ListAccountTransactions")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

// boolPtr returns a pointer to a bool value.
func boolPtr(b bool) *bool {
	return &b
}
