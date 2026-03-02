package client

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// ============================================================================
// ACCOUNT OPERATIONS - TABLE-DRIVEN TESTS
// ============================================================================

func TestClient_UpdateAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		accountID   string
		request     *types.UpdateAccountRequest
		mockStatus  int
		wantErr     bool
		validateReq func(t *testing.T, r *http.Request, body []byte)
	}{
		{
			name:      "success: update account owner",
			accountID: "acc-test-001",
			request: &types.UpdateAccountRequest{
				IssuerRequestID: "req123",
				AccountOwner: map[string]any{
					"fullName": "Updated Name",
					"email":    "updated@email.com",
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request, body []byte) {
				AssertEqual(t, http.MethodPut, r.Method, "method")

				var reqBody map[string]any
				_ = json.Unmarshal(body, &reqBody)

				owner, ok := reqBody["accountOwner"].(map[string]any)
				if !ok {
					t.Fatal("accountOwner not found in request body")
				}
				AssertEqual(t, "Updated Name", owner["fullName"], "owner.fullName")
			},
		},
		{
			name:      "success: update billing address",
			accountID: "acc-test-002",
			request: &types.UpdateAccountRequest{
				BillingAddress: map[string]any{
					"addressLine1": "Rua Nova 456",
					"city":         "Rio de Janeiro",
					"state":        "RJ",
					"zipcode":      "20000000",
					"country":      "BR",
					"neighborhood": "Centro",
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request, body []byte) {
				var reqBody map[string]any
				_ = json.Unmarshal(body, &reqBody)

				address, ok := reqBody["billingAddress"].(map[string]any)
				if !ok {
					t.Fatal("billingAddress not found in request body")
				}
				AssertEqual(t, "Rio de Janeiro", address["city"], "address.city")
			},
		},
		{
			name:      "error: account not found",
			accountID: "acc-nonexistent",
			request: &types.UpdateAccountRequest{
				IssuerRequestID: "req456",
			},
			mockStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := NewMockServer(t, &MockServerConfig{
				Status:   tt.mockStatus,
				Response: map[string]any{"resultData": map[string]any{"resultCode": 0}},
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if !strings.Contains(r.URL.Path, tt.accountID) {
						t.Errorf("expected path to contain %s, got %s", tt.accountID, r.URL.Path)
					}
					if tt.validateReq != nil {
						tt.validateReq(t, r, body)
					}
				},
			})
			defer mockServer.Close()

			client := NewTestClient(t, mockServer)
			err := client.UpdateAccount(context.Background(), tt.accountID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "expected error")
			} else {
				AssertNoError(t, err)
			}
		})
	}
}

func TestClient_CreateCredit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		accountID  string
		request    *types.NewCreditRequest
		mockStatus int
		wantErr    bool
	}{
		{
			name:      "success: create payment credit",
			accountID: "acc-test-001",
			request: &types.NewCreditRequest{
				IssuerRequestID: "req123",
				Amount:          types.Amount{Amount: 10000, CurrencyCode: 986},
				Type:            "payment",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "success: create adjustment credit",
			accountID: "acc-test-002",
			request: &types.NewCreditRequest{
				IssuerRequestID: "req456",
				Amount:          types.Amount{Amount: 5000, CurrencyCode: 986},
				Type:            "adjustment",
				Description:     "Refund adjustment",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "error: invalid account",
			accountID: "invalid",
			request: &types.NewCreditRequest{
				Amount: types.Amount{Amount: 1000, CurrencyCode: 986},
			},
			mockStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := NewMockServer(t, &MockServerConfig{
				Status:   tt.mockStatus,
				Response: map[string]any{"resultData": map[string]any{"resultCode": 0}, "transactionId": "tx123"},
			})
			defer mockServer.Close()

			client := NewTestClient(t, mockServer)
			_, err := client.CreateCredit(context.Background(), tt.accountID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "expected error")
			} else {
				AssertNoError(t, err)
			}
		})
	}
}

func TestClient_ListAccounts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		params       *types.ListAccountsParams
		mockResponse any
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, r *http.Request)
	}{
		{
			name:   "success: list accounts without params",
			params: nil,
			mockResponse: map[string]any{
				"hasMore": true,
				"accounts": []map[string]any{
					{"accountId": "acc001", "psProductCode": "CREDIT"},
					{"accountId": "acc002", "psProductCode": "DEBIT"},
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: list with limit",
			params: &types.ListAccountsParams{
				Limit: 10,
			},
			mockResponse: map[string]any{
				"hasMore":  false,
				"accounts": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("limit") != "10" {
					t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
				}
			},
		},
		{
			name: "success: list with pagination cursors",
			params: &types.ListAccountsParams{
				StartingAfter: "acc100",
				EndingBefore:  "acc200",
			},
			mockResponse: map[string]any{
				"hasMore":  true,
				"accounts": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("starting_after") != "acc100" {
					t.Errorf("expected starting_after=acc100, got %s", r.URL.Query().Get("starting_after"))
				}
				if r.URL.Query().Get("ending_before") != "acc200" {
					t.Errorf("expected ending_before=acc200, got %s", r.URL.Query().Get("ending_before"))
				}
			},
		},
		{
			name: "success: list with identity document filter",
			params: &types.ListAccountsParams{
				IdentityDocumentNumber: "12345678901",
			},
			mockResponse: map[string]any{
				"hasMore": false,
				"accounts": []map[string]any{
					{"accountId": "acc001"},
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("identity_document_number") != "12345678901" {
					t.Errorf("expected identity_document_number=12345678901, got %s", r.URL.Query().Get("identity_document_number"))
				}
			},
		},
		{
			name: "success: list with full name filter",
			params: &types.ListAccountsParams{
				FullName: "JohnDoe",
			},
			mockResponse: map[string]any{
				"hasMore":  false,
				"accounts": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("full_name") != "JohnDoe" {
					t.Errorf("expected full_name=JohnDoe, got %s", r.URL.Query().Get("full_name"))
				}
			},
		},
		{
			name: "success: list with product code filter",
			params: &types.ListAccountsParams{
				PSProductCode: 100,
			},
			mockResponse: map[string]any{
				"hasMore":  false,
				"accounts": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("ps_product_code") != "100" {
					t.Errorf("expected ps_product_code=100, got %s", r.URL.Query().Get("ps_product_code"))
				}
			},
		},
		{
			name: "success: list with account status filter",
			params: &types.ListAccountsParams{
				AccountStatus: "ACTIVE",
			},
			mockResponse: map[string]any{
				"hasMore":  false,
				"accounts": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("account_status") != "ACTIVE" {
					t.Errorf("expected account_status=ACTIVE, got %s", r.URL.Query().Get("account_status"))
				}
			},
		},
		{
			name: "success: list with issuer account ID filter",
			params: &types.ListAccountsParams{
				IssuerAccountID: "issuer-001",
			},
			mockResponse: map[string]any{
				"hasMore":  false,
				"accounts": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("issuer_id") != "issuer-001" {
					t.Errorf("expected issuer_id=issuer-001, got %s", r.URL.Query().Get("issuer_id"))
				}
			},
		},
		{
			name: "success: list with included since filter",
			params: &types.ListAccountsParams{
				IncludedSince: "2024-01-01",
			},
			mockResponse: map[string]any{
				"hasMore":  false,
				"accounts": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("included_since") != "2024-01-01" {
					t.Errorf("expected included_since=2024-01-01, got %s", r.URL.Query().Get("included_since"))
				}
			},
		},
		{
			name: "success: list with sort filter",
			params: &types.ListAccountsParams{
				Sort: "created_at_desc",
			},
			mockResponse: map[string]any{
				"hasMore":  false,
				"accounts": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("sort") != "created_at_desc" {
					t.Errorf("expected sort=created_at_desc, got %s", r.URL.Query().Get("sort"))
				}
			},
		},
		{
			name: "success: list with first flag",
			params: &types.ListAccountsParams{
				First: true,
			},
			mockResponse: map[string]any{
				"hasMore":  false,
				"accounts": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("first") != "true" {
					t.Errorf("expected first=true, got %s", r.URL.Query().Get("first"))
				}
			},
		},
		{
			name: "success: list with all params",
			params: &types.ListAccountsParams{
				Limit:                  50,
				StartingAfter:          "acc100",
				EndingBefore:           "acc200",
				IdentityDocumentNumber: "12345678901",
				FullName:               "TestUser",
				PSProductCode:          100,
				AccountStatus:          "ACTIVE",
				IssuerAccountID:        "issuer-001",
				IncludedSince:          "2024-01-01",
				Sort:                   "created_at_asc",
				First:                  true,
			},
			mockResponse: map[string]any{
				"hasMore":  false,
				"accounts": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "error: server error",
			params:       nil,
			mockResponse: nil,
			mockStatus:   http.StatusInternalServerError,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/accounts",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						tt.validateReq(t, r)
					}
				},
			})
			defer mockServer.Close()

			client := NewTestClient(t, mockServer)
			resp, err := client.ListAccounts(context.Background(), tt.params)

			if tt.wantErr {
				AssertError(t, err, "ListAccounts")
			} else {
				AssertNoError(t, err, "ListAccounts")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_CancelAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		accountID  string
		request    *types.CancelAccountRequest
		mockStatus int
		wantErr    bool
	}{
		{
			name:      "success: cancel account",
			accountID: "acc-test-001",
			request: &types.CancelAccountRequest{
				CancellationCode: 1,
				Reason:           "Customer request",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "error: account not found",
			accountID: "acc-nonexistent",
			request: &types.CancelAccountRequest{
				CancellationCode: 1,
			},
			mockStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := NewMockServer(t, &MockServerConfig{
				Status:   tt.mockStatus,
				Response: map[string]any{"resultData": map[string]any{"resultCode": 0}},
			})
			defer mockServer.Close()

			client := NewTestClient(t, mockServer)
			err := client.CancelAccount(context.Background(), tt.accountID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "expected error")
			} else {
				AssertNoError(t, err)
			}
		})
	}
}

func TestClient_BlockAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		accountID   string
		request     *types.BlockAccountRequest
		mockStatus  int
		wantErr     bool
		validateReq func(t *testing.T, r *http.Request, body []byte)
	}{
		{
			name:      "success: block account",
			accountID: "acc-test-001",
			request: &types.BlockAccountRequest{
				BlockCode: 1,
				Reason:    "Fraud suspected",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request, body []byte) {
				AssertEqual(t, http.MethodPost, r.Method, "method")
				var reqBody map[string]any
				_ = json.Unmarshal(body, &reqBody)
				AssertEqual(t, "Fraud suspected", reqBody["reason"], "reason")
			},
		},
		{
			name:      "error: account not found",
			accountID: "acc-nonexistent",
			request: &types.BlockAccountRequest{
				BlockCode: 1,
			},
			mockStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/accounts/" + tt.accountID + "/blockAccount",
				Status:   tt.mockStatus,
				Response: map[string]any{"resultData": map[string]any{"resultCode": 0}},
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						tt.validateReq(t, r, body)
					}
				},
			})
			defer mockServer.Close()

			client := NewTestClient(t, mockServer)
			err := client.BlockAccount(context.Background(), tt.accountID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "expected error")
			} else {
				AssertNoError(t, err)
			}
		})
	}
}

func TestClient_UnblockAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		accountID  string
		request    *types.UnblockAccountRequest
		mockStatus int
		wantErr    bool
	}{
		{
			name:      "success: unblock account",
			accountID: "acc-test-001",
			request: &types.UnblockAccountRequest{
				Reason: "Fraud investigation cleared",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "success: unblock without reason",
			accountID:  "acc-test-002",
			request:    &types.UnblockAccountRequest{},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "error: account not found",
			accountID: "acc-nonexistent",
			request: &types.UnblockAccountRequest{
				Reason: "Test",
			},
			mockStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/accounts/" + tt.accountID + "/unblockAccount",
				Status:   tt.mockStatus,
				Response: map[string]any{"resultData": map[string]any{"resultCode": 0}},
			})
			defer mockServer.Close()

			client := NewTestClient(t, mockServer)
			err := client.UnblockAccount(context.Background(), tt.accountID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "expected error")
			} else {
				AssertNoError(t, err)
			}
		})
	}
}

func TestClient_GetAccountBalance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		accountID    string
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:      "success: get account balance",
			accountID: "acc-test-001",
			mockResponse: map[string]any{
				"available": map[string]any{"amount": 500000, "currencyCode": 986},
				"current":   map[string]any{"amount": 450000, "currencyCode": 986},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "success: zero balance",
			accountID: "acc-test-002",
			mockResponse: map[string]any{
				"available": map[string]any{"amount": 0, "currencyCode": 986},
				"current":   map[string]any{"amount": 0, "currencyCode": 986},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "error: account not found",
			accountID:    "acc-nonexistent",
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/accounts/" + tt.accountID + "/balance",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer mockServer.Close()

			client := NewTestClient(t, mockServer)
			resp, err := client.GetAccountBalance(context.Background(), tt.accountID)

			if tt.wantErr {
				AssertError(t, err, "expected error")
			} else {
				AssertNoError(t, err)
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_ListAccountCards(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		accountID    string
		mockResponse any
		mockStatus   int
		wantErr      bool
		wantLen      int
	}{
		{
			name:      "success: list account cards",
			accountID: "acc-test-001",
			mockResponse: []map[string]any{
				{"cardId": "card001", "status": "ACTIVE", "last4Digits": "1234"},
				{"cardId": "card002", "status": "BLOCKED", "last4Digits": "5678"},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			wantLen:    2,
		},
		{
			name:         "success: no cards",
			accountID:    "acc-test-002",
			mockResponse: []map[string]any{},
			mockStatus:   http.StatusOK,
			wantErr:      false,
			wantLen:      0,
		},
		{
			name:         "error: account not found",
			accountID:    "acc-nonexistent",
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/accounts/" + tt.accountID + "/cards",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer mockServer.Close()

			client := NewTestClient(t, mockServer)
			resp, err := client.ListAccountCards(context.Background(), tt.accountID)

			if tt.wantErr {
				AssertError(t, err, "expected error")
			} else {
				AssertNoError(t, err)
				AssertEqual(t, tt.wantLen, len(resp), "cards count")
			}
		})
	}
}

func TestClient_GetCreditVerificationStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		accountID    string
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:      "success: get credit verification status approved",
			accountID: "acc-test-001",
			mockResponse: map[string]any{
				"status":   "APPROVED",
				"decision": "AUTO_APPROVED",
				"score":    850,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "success: get credit verification status pending",
			accountID: "acc-test-002",
			mockResponse: map[string]any{
				"status": "PENDING",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "error: account not found",
			accountID:    "acc-nonexistent",
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/accounts/" + tt.accountID + "/creditVerificationStatus",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer mockServer.Close()

			client := NewTestClient(t, mockServer)
			resp, err := client.GetCreditVerificationStatus(context.Background(), tt.accountID)

			if tt.wantErr {
				AssertError(t, err, "expected error")
			} else {
				AssertNoError(t, err)
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_ListScheduledTransactions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		accountID    string
		mockResponse any
		mockStatus   int
		wantErr      bool
		wantLen      int
	}{
		{
			name:      "success: list scheduled transactions",
			accountID: "acc-test-001",
			mockResponse: []map[string]any{
				{"transactionId": "tx001", "amount": map[string]any{"amount": 10000}, "scheduledDate": "2024-02-01"},
				{"transactionId": "tx002", "amount": map[string]any{"amount": 20000}, "scheduledDate": "2024-02-15"},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			wantLen:    2,
		},
		{
			name:         "success: no scheduled transactions",
			accountID:    "acc-test-002",
			mockResponse: []map[string]any{},
			mockStatus:   http.StatusOK,
			wantErr:      false,
			wantLen:      0,
		},
		{
			name:         "error: account not found",
			accountID:    "acc-nonexistent",
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/accounts/" + tt.accountID + "/findScheduledTransactions",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer mockServer.Close()

			client := NewTestClient(t, mockServer)
			resp, err := client.ListScheduledTransactions(context.Background(), tt.accountID)

			if tt.wantErr {
				AssertError(t, err, "expected error")
			} else {
				AssertNoError(t, err)
				AssertEqual(t, tt.wantLen, len(resp), "scheduled transactions count")
			}
		})
	}
}

func TestClient_ListCredits(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		accountID    string
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:      "success: list credits",
			accountID: "acc-test-001",
			mockResponse: map[string]any{
				"credits": []map[string]any{
					{"creditId": "cr001", "amount": map[string]any{"amount": 10000}, "type": "payment"},
					{"creditId": "cr002", "amount": map[string]any{"amount": 5000}, "type": "refund"},
				},
				"hasMore": false,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "success: empty credits list",
			accountID: "acc-test-002",
			mockResponse: map[string]any{
				"credits": []map[string]any{},
				"hasMore": false,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "error: account not found",
			accountID:    "acc-nonexistent",
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/accounts/" + tt.accountID + "/transactions/credits",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer mockServer.Close()

			client := NewTestClient(t, mockServer)
			resp, err := client.ListCredits(context.Background(), tt.accountID)

			if tt.wantErr {
				AssertError(t, err, "expected error")
			} else {
				AssertNoError(t, err)
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_CreateDebit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		accountID   string
		request     *types.CreateDebitRequest
		mockStatus  int
		wantErr     bool
		validateReq func(t *testing.T, body map[string]any)
	}{
		{
			name:      "success: create debit transaction",
			accountID: "acc-test-001",
			request: &types.CreateDebitRequest{
				IssuerRequestID: "req123",
				Amount:          types.Amount{Amount: 10000, CurrencyCode: 986},
				Type:            "fee",
				Description:     "Monthly fee",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]any) {
				AssertEqual(t, "fee", body["type"], "type")
				AssertEqual(t, "Monthly fee", body["description"], "description")
			},
		},
		{
			name:      "success: create debit with minimum fields",
			accountID: "acc-test-002",
			request: &types.CreateDebitRequest{
				Amount: types.Amount{Amount: 5000, CurrencyCode: 986},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "error: insufficient balance",
			accountID: "acc-test-003",
			request: &types.CreateDebitRequest{
				Amount: types.Amount{Amount: 999999999, CurrencyCode: 986},
			},
			mockStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:      "error: account not found",
			accountID: "acc-nonexistent",
			request: &types.CreateDebitRequest{
				Amount: types.Amount{Amount: 1000, CurrencyCode: 986},
			},
			mockStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockServer := NewMockServer(t, &MockServerConfig{
				Method: http.MethodPost,
				Path:   "/accounts/" + tt.accountID + "/transactions/debits",
				Status: tt.mockStatus,
				Response: map[string]any{
					"transactionId": "tx123",
					"status":        "COMPLETED",
				},
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]any
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer mockServer.Close()

			client := NewTestClient(t, mockServer)
			resp, err := client.CreateDebit(context.Background(), tt.accountID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "expected error")
			} else {
				AssertNoError(t, err)
				AssertNotNil(t, resp, "response")
			}
		})
	}
}
