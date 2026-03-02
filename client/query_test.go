package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

func TestPathParam(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"simple-id", "simple-id"},
		{"id/with/slashes", "id%2Fwith%2Fslashes"},
		{"id?query=1", "id%3Fquery=1"},
		{"id#fragment", "id%23fragment"},
		{"id with spaces", "id%20with%20spaces"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := pathParam(tt.input)
			if got != tt.want {
				t.Errorf("pathParam(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// Tests for functions that use query building

func TestClient_ListDisputes_WithParams(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/disputes", http.StatusOK, map[string]any{
		"disputes": []map[string]any{},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.ListDisputesRequest{
		DisputeStatus: "OPEN",
		BeginningDate: "2024-01-01",
		EndingDate:    "2024-12-31",
		Limit:         10,
		StartingAfter: "disp100",
	}
	resp, err := client.ListDisputes(context.Background(), req)
	if err != nil {
		t.Fatalf("ListDisputes with params failed: %v", err)
	}
	if resp == nil {
		t.Error("ListDisputes returned nil")
	}
}

func TestClient_ListAccountTransactions_WithParams(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/accounts/acc123/transactions", http.StatusOK, map[string]any{
		"transactions": []map[string]any{},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.ListTransactionsRequest{
		BeginningDate:   "2024-01-01",
		EndingDate:      "2024-12-31",
		Limit:           10,
		StartingAfter:   "tx100",
		TransactionType: "PURCHASE",
	}
	resp, err := client.ListAccountTransactions(context.Background(), "acc123", req)
	if err != nil {
		t.Fatalf("ListAccountTransactions with params failed: %v", err)
	}
	if resp == nil {
		t.Error("ListAccountTransactions returned nil")
	}
}

func TestClient_ListCobranded_WithParams(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/cobranded", http.StatusOK, map[string]any{
		"data": []map[string]any{},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	params := &ListCobrandedParams{
		Document:   "12345678901",
		AcquirerID: 1,
	}
	resp, err := client.ListCobranded(context.Background(), params)
	if err != nil {
		t.Fatalf("ListCobranded with params failed: %v", err)
	}
	if resp == nil {
		t.Error("ListCobranded returned nil")
	}
}

func TestClient_ListInclusiveTransactions_WithParams(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/inclusives", http.StatusOK, map[string]any{
		"transactions": []map[string]any{},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	params := &ListInclusiveTransactionsParams{
		Limit:         10,
		StartingAfter: "txn100",
		BeginningDate: "2024-01-01",
		EndingDate:    "2024-12-31",
	}
	resp, err := client.ListInclusiveTransactions(context.Background(), params)
	if err != nil {
		t.Fatalf("ListInclusiveTransactions with params failed: %v", err)
	}
	if resp == nil {
		t.Error("ListInclusiveTransactions returned nil")
	}
}

// Test SimulateAdvancePayment/RequestAdvancePayment with request body
func TestClient_SimulateAdvancePayment_WithBody(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/accounts/acc123/installmentAdvanceSimulation/txn123", http.StatusOK, map[string]any{
		"simulation": map[string]any{},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.AdvancePaymentRequest{
		TransactionID:         "txn123",
		InstallmentsToAdvance: 3,
	}
	resp, err := client.SimulateAdvancePayment(context.Background(), "acc123", "txn123", req)
	if err != nil {
		t.Fatalf("SimulateAdvancePayment with body failed: %v", err)
	}
	if resp == nil {
		t.Error("SimulateAdvancePayment returned nil")
	}
}

func TestClient_RequestAdvancePayment_WithBody(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/accounts/acc123/installmentAdvanceRequest/txn123", http.StatusOK, map[string]any{
		"resultCode": 0,
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.AdvancePaymentRequest{
		TransactionID:         "txn123",
		InstallmentsToAdvance: 3,
	}
	resp, err := client.RequestAdvancePayment(context.Background(), "acc123", "txn123", req)
	if err != nil {
		t.Fatalf("RequestAdvancePayment with body failed: %v", err)
	}
	if resp == nil {
		t.Error("RequestAdvancePayment returned nil")
	}
}
