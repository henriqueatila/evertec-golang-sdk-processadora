package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

func TestClient_GetDynamicProvisioningBalance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mockStatus   int
		mockResponse *types.DynamicProvisioningBalance
		wantErr      bool
		validate     func(t *testing.T, resp *types.DynamicProvisioningBalance)
	}{
		{
			name:       "success: get balance",
			mockStatus: http.StatusOK,
			mockResponse: &types.DynamicProvisioningBalance{
				CurrentBalance: &types.Amount{Amount: 1000000, CurrencyCode: 986},
			},
			wantErr: false,
			validate: func(t *testing.T, resp *types.DynamicProvisioningBalance) {
				AssertNotNil(t, resp, "response")
				AssertNotNil(t, resp.CurrentBalance, "CurrentBalance")
				AssertEqual(t, int64(1000000), resp.CurrentBalance.Amount, "CurrentBalance.Amount")
				AssertEqual(t, 986, resp.CurrentBalance.CurrencyCode, "CurrentBalance.CurrencyCode")
			},
		},
		{
			name:       "success: balance without current balance",
			mockStatus: http.StatusOK,
			mockResponse: &types.DynamicProvisioningBalance{
				CurrentBalance: nil,
			},
			wantErr: false,
			validate: func(t *testing.T, resp *types.DynamicProvisioningBalance) {
				AssertNotNil(t, resp, "response")
				AssertNil(t, resp.CurrentBalance, "CurrentBalance")
			},
		},
		{
			name:       "error: unauthorized",
			mockStatus: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "error: server error",
			mockStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/dynamicProvisioningAccount/balance",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetDynamicProvisioningBalance(context.Background())

			if tt.wantErr {
				AssertError(t, err, "expected error")
				return
			}

			AssertNoError(t, err, "unexpected error")
			if tt.validate != nil {
				tt.validate(t, resp)
			}
		})
	}
}

func TestClient_ListDynamicProvisioningCredits(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		params       *ListDynamicProvisioningCreditsParams
		mockStatus   int
		mockResponse *types.DynamicProvisioningCreditList
		wantErr      bool
		validateReq  func(t *testing.T, r *http.Request, body []byte)
		validate     func(t *testing.T, resp *types.DynamicProvisioningCreditList)
	}{
		{
			name: "success: list credits with date range",
			params: &ListDynamicProvisioningCreditsParams{
				BeginningDate: "2024-01-01",
				EndDate:       "2024-12-31",
			},
			mockStatus: http.StatusOK,
			mockResponse: &types.DynamicProvisioningCreditList{
				Credits: []types.DynamicProvisioningCredit{
					{
						TransactionID:          "credit-001",
						Amount:                 types.DebitOrCreditAmount{Amount: 100000, CurrencyCode: 986, DebitOrCredit: "credit"},
						CreditReceivedDateTime: "2024-01-15T10:00:00Z",
						Status:                 "confirmed",
					},
					{
						TransactionID:          "credit-002",
						Amount:                 types.DebitOrCreditAmount{Amount: 200000, CurrencyCode: 986, DebitOrCredit: "credit"},
						CreditReceivedDateTime: "2024-01-14T10:00:00Z",
						Status:                 "confirmed",
					},
				},
				HasMoreCredits: false,
			},
			wantErr: false,
			validateReq: func(t *testing.T, r *http.Request, body []byte) {
				AssertQueryParam(t, r, "beginning_date", "2024-01-01")
				AssertQueryParam(t, r, "end_date", "2024-12-31")
			},
			validate: func(t *testing.T, resp *types.DynamicProvisioningCreditList) {
				AssertNotNil(t, resp, "response")
				AssertEqual(t, 2, len(resp.Credits), "Credits length")
				AssertEqual(t, false, resp.HasMoreCredits, "HasMoreCredits")
				AssertEqual(t, "credit-001", resp.Credits[0].TransactionID, "first credit ID")
				AssertEqual(t, int64(100000), resp.Credits[0].Amount.Amount, "first credit amount")
			},
		},
		{
			name:       "success: empty list",
			params:     &ListDynamicProvisioningCreditsParams{BeginningDate: "2025-01-01", EndDate: "2025-01-31"},
			mockStatus: http.StatusOK,
			mockResponse: &types.DynamicProvisioningCreditList{
				Credits:        []types.DynamicProvisioningCredit{},
				HasMoreCredits: false,
			},
			wantErr: false,
			validate: func(t *testing.T, resp *types.DynamicProvisioningCreditList) {
				AssertEqual(t, 0, len(resp.Credits), "Credits length")
			},
		},
		{
			name:       "success: nil params",
			params:     nil,
			mockStatus: http.StatusOK,
			mockResponse: &types.DynamicProvisioningCreditList{
				Credits:        []types.DynamicProvisioningCredit{},
				HasMoreCredits: false,
			},
			wantErr: false,
		},
		{
			name: "success: has more pages",
			params: &ListDynamicProvisioningCreditsParams{
				BeginningDate: "2024-01-01",
				EndDate:       "2024-12-31",
			},
			mockStatus: http.StatusOK,
			mockResponse: &types.DynamicProvisioningCreditList{
				Credits: []types.DynamicProvisioningCredit{
					{
						TransactionID:          "credit-001",
						Amount:                 types.DebitOrCreditAmount{Amount: 100000, CurrencyCode: 986, DebitOrCredit: "credit"},
						CreditReceivedDateTime: "2024-01-15T10:00:00Z",
						Status:                 "confirmed",
					},
				},
				HasMoreCredits: true,
			},
			wantErr: false,
			validate: func(t *testing.T, resp *types.DynamicProvisioningCreditList) {
				AssertEqual(t, true, resp.HasMoreCredits, "HasMoreCredits")
			},
		},
		{
			name:       "error: bad request",
			params:     &ListDynamicProvisioningCreditsParams{BeginningDate: "invalid-date"},
			mockStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "error: unauthorized",
			params:     &ListDynamicProvisioningCreditsParams{BeginningDate: "2024-01-01", EndDate: "2024-12-31"},
			mockStatus: http.StatusUnauthorized,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:      http.MethodGet,
				Path:        "/dynamicProvisioningAccount/balance/credits",
				Status:      tt.mockStatus,
				Response:    tt.mockResponse,
				ValidateReq: tt.validateReq,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.ListDynamicProvisioningCredits(context.Background(), tt.params)

			if tt.wantErr {
				AssertError(t, err, "expected error")
				return
			}

			AssertNoError(t, err, "unexpected error")
			if tt.validate != nil {
				tt.validate(t, resp)
			}
		})
	}
}

func TestClient_DynamicProvisioning_ContextCancellation(t *testing.T) {
	t.Parallel()

	server := NewMockServer(t, &MockServerConfig{
		Method:   http.MethodGet,
		Path:     "/dynamicProvisioningAccount/balance",
		Status:   http.StatusOK,
		Response: &types.DynamicProvisioningBalance{},
		Delay:    100 * time.Millisecond,
	})
	defer server.Close()

	client := NewTestClient(t, server)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.GetDynamicProvisioningBalance(ctx)
	AssertError(t, err, "expected context cancellation error")
}
