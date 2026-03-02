package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

func TestClient_CreateFraudNotification(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.FraudNotificationRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, body map[string]any)
	}{
		{
			name: "success: create fraud notification",
			request: &types.FraudNotificationRequest{
				AccountID:     "acc123",
				TransactionID: "tx123",
				FraudType:     "CARD_NOT_PRESENT",
			},
			mockResponse: map[string]any{
				"fraudNotificationId": "fn123",
				"status":              "CREATED",
				"accountId":           "acc123",
				"transactionId":       "tx123",
			},
			mockStatus: http.StatusCreated,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]any) {
				AssertEqual(t, "acc123", body["accountId"], "accountId")
				AssertEqual(t, "tx123", body["transactionId"], "transactionId")
			},
		},
		{
			name: "success: create fraud notification with issuer fraud id",
			request: &types.FraudNotificationRequest{
				IssuerFraudID: "issuer-fraud-001",
				AccountID:     "acc456",
				TransactionID: "tx456",
				FraudType:     "LOST_CARD",
			},
			mockResponse: map[string]any{
				"fraudNotificationId": "fn456",
				"status":              "CREATED",
			},
			mockStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "error: account not found",
			request: &types.FraudNotificationRequest{
				AccountID:     "acc-nonexistent",
				TransactionID: "tx789",
				FraudType:     "STOLEN_CARD",
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
				Path:     "/frauds/notification",
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
			resp, err := client.CreateFraudNotification(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "CreateFraudNotification")
			} else {
				AssertNoError(t, err, "CreateFraudNotification")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_ListFraudNotifications(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mockResponse any
		mockStatus   int
		wantErr      bool
		wantLen      int
	}{
		{
			name: "success: list fraud notifications",
			mockResponse: []map[string]any{
				{"fraudNotificationId": "fn001", "status": "CREATED", "fraudType": "CARD_NOT_PRESENT"},
				{"fraudNotificationId": "fn002", "status": "INVESTIGATING", "fraudType": "STOLEN_CARD"},
				{"fraudNotificationId": "fn003", "status": "RESOLVED", "fraudType": "LOST_CARD"},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			wantLen:    3,
		},
		{
			name:         "success: empty list",
			mockResponse: []map[string]any{},
			mockStatus:   http.StatusOK,
			wantErr:      false,
			wantLen:      0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/frauds/notification",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.ListFraudNotifications(context.Background())

			if tt.wantErr {
				AssertError(t, err, "ListFraudNotifications")
			} else {
				AssertNoError(t, err, "ListFraudNotifications")
				AssertEqual(t, tt.wantLen, len(resp), "notifications count")
			}
		})
	}
}
