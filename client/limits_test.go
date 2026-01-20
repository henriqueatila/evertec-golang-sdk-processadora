package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

func TestClient_UpdateMaxCreditLimits(t *testing.T) {
	server := mockServer(t, http.MethodPut, "/accounts/acc123/maxCreditLimits", http.StatusOK, map[string]interface{}{
		"resultCode":        0,
		"resultDescription": "Success",
		"issuerRequestId":   "req123",
		"psResponseId":      "ps123",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.MaxCreditLimitsRequest{
		TotalLimit: &types.Amount{
			Amount:       100000,
			CurrencyCode: 986,
		},
		WithdrawalLimit: &types.Amount{
			Amount:       50000,
			CurrencyCode: 986,
		},
	}

	resp, err := client.UpdateMaxCreditLimits(context.Background(), "acc123", req)
	if err != nil {
		t.Fatalf("UpdateMaxCreditLimits failed: %v", err)
	}

	if resp == nil {
		t.Fatal("UpdateMaxCreditLimits returned nil response")
	}

	if resp.ResultCode != 0 {
		t.Errorf("Expected resultCode 0, got %d", resp.ResultCode)
	}
}

func TestClient_ChangeUsableCreditLimits(t *testing.T) {
	server := mockServer(t, http.MethodPut, "/accounts/acc123/changeUsableCreditLimits", http.StatusOK, map[string]interface{}{
		"resultCode":        0,
		"resultDescription": "Success",
		"issuerRequestId":   "req456",
		"psResponseId":      "ps456",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.ChangeUsableCreditLimitsRequest{
		NewUsableTotalLimit: 75000,
	}

	resp, err := client.ChangeUsableCreditLimits(context.Background(), "acc123", req)
	if err != nil {
		t.Fatalf("ChangeUsableCreditLimits failed: %v", err)
	}

	if resp == nil {
		t.Fatal("ChangeUsableCreditLimits returned nil response")
	}

	if resp.ResultCode != 0 {
		t.Errorf("Expected resultCode 0, got %d", resp.ResultCode)
	}
}
