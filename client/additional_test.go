package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// Client options tests

func TestClient_WithTimeout(t *testing.T) {
	opt := WithTimeout(30 * time.Second)
	// Verify option doesn't panic
	if opt == nil {
		t.Error("WithTimeout returned nil")
	}
}

// Travel notice tests

func TestClient_ListTravelNotices(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/accounts/acc123/travelNotice", http.StatusOK, []map[string]interface{}{
		{"noticeId": "notice1"},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.ListTravelNotices(context.Background(), "acc123")
	if err != nil {
		t.Fatalf("ListTravelNotices failed: %v", err)
	}
	if resp == nil {
		t.Error("ListTravelNotices returned nil")
	}
}

func TestClient_GetTravelNotice(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/accounts/acc123/travelNotice/notice1", http.StatusOK, map[string]interface{}{
		"noticeId": "notice1",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.GetTravelNotice(context.Background(), "acc123", "notice1")
	if err != nil {
		t.Fatalf("GetTravelNotice failed: %v", err)
	}
	if resp == nil {
		t.Error("GetTravelNotice returned nil")
	}
}

func TestClient_UpdateTravelNotice(t *testing.T) {
	server := mockServer(t, http.MethodPut, "/accounts/acc123/travelNotice/notice1", http.StatusOK, map[string]interface{}{
		"noticeId": "notice1",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.UpdateTravelNotice(context.Background(), "acc123", "notice1", &types.UpdateTravelNotice{})
	if err != nil {
		t.Fatalf("UpdateTravelNotice failed: %v", err)
	}
	if resp == nil {
		t.Error("UpdateTravelNotice returned nil")
	}
}

// Virtual card tests

func TestClient_CreateCardVirtualCard(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/cards/card123/virtualcards", http.StatusCreated, map[string]interface{}{
		"virtualCardId": "vcard123",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.CreateCardVirtualCard(context.Background(), "card123", &types.CreateVirtualCardRequest{})
	if err != nil {
		t.Fatalf("CreateCardVirtualCard failed: %v", err)
	}
	if resp == nil {
		t.Error("CreateCardVirtualCard returned nil")
	}
}

func TestClient_ListCardVirtualCards(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/cards/card123/virtualcards/list", http.StatusOK, map[string]interface{}{
		"virtualCards": []map[string]interface{}{
			{"virtualCardId": "vcard1"},
		},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.ListCardVirtualCards(context.Background(), "card123")
	if err != nil {
		t.Fatalf("ListCardVirtualCards failed: %v", err)
	}
	if resp == nil {
		t.Error("ListCardVirtualCards returned nil")
	}
}

func TestClient_ModifyVirtualCardCVV(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/virtualcards/vcard123/modify/cvv", http.StatusOK, map[string]interface{}{
		"newCvv": "123",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.ModifyVirtualCardCVV(context.Background(), "vcard123")
	if err != nil {
		t.Fatalf("ModifyVirtualCardCVV failed: %v", err)
	}
	if resp == nil {
		t.Error("ModifyVirtualCardCVV returned nil")
	}
}

// XPays tests

func TestClient_SuspendOrResumeDeviceTokensByCard(t *testing.T) {
	// API returns array of DeviceToken directly per official spec
	server := mockServer(t, http.MethodPost, "/cards/card123/deviceTokens/suspendOrResume", http.StatusOK, []map[string]interface{}{
		{"deviceTokenId": "token-001", "suspensionStatus": "SUSPENDED"},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.SuspendOrResumeDeviceTokensByCard(context.Background(), "card123", &types.SuspendResumeRequest{})
	if err != nil {
		t.Fatalf("SuspendOrResumeDeviceTokensByCard failed: %v", err)
	}
	if resp == nil {
		t.Error("SuspendOrResumeDeviceTokensByCard returned nil")
	}
}

func TestClient_TerminateDeviceTokensByCard(t *testing.T) {
	// API returns array of DeviceToken directly per official spec
	server := mockServer(t, http.MethodPost, "/cards/card123/deviceTokens/terminate", http.StatusOK, []map[string]interface{}{
		{"deviceTokenId": "token-001", "activationStatus": "TERMINATED"},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.TerminateDeviceTokensByCard(context.Background(), "card123", &types.TerminateTokensRequest{})
	if err != nil {
		t.Fatalf("TerminateDeviceTokensByCard failed: %v", err)
	}
	if resp == nil {
		t.Error("TerminateDeviceTokensByCard returned nil")
	}
}

// Installments tests

func TestClient_GetInstallmentOptions(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/accounts/acc123/installmentSimulation", http.StatusOK, map[string]interface{}{
		"options": []map[string]interface{}{},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.GetInstallmentOptions(context.Background(), "acc123")
	if err != nil {
		t.Fatalf("GetInstallmentOptions failed: %v", err)
	}
	if resp == nil {
		t.Error("GetInstallmentOptions returned nil")
	}
}

func TestClient_SimulateAdvancePayment(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/accounts/acc123/installmentAdvanceSimulation/txn123", http.StatusOK, map[string]interface{}{
		"simulation": map[string]interface{}{},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.SimulateAdvancePayment(context.Background(), "acc123", "txn123", &types.AdvancePaymentRequest{})
	if err != nil {
		t.Fatalf("SimulateAdvancePayment failed: %v", err)
	}
	if resp == nil {
		t.Error("SimulateAdvancePayment returned nil")
	}
}

func TestClient_RequestAdvancePayment(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/accounts/acc123/installmentAdvanceRequest/txn123", http.StatusOK, map[string]interface{}{
		"resultCode": 0,
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.RequestAdvancePayment(context.Background(), "acc123", "txn123", &types.AdvancePaymentRequest{})
	if err != nil {
		t.Fatalf("RequestAdvancePayment failed: %v", err)
	}
	if resp == nil {
		t.Error("RequestAdvancePayment returned nil")
	}
}

// Status tests

func TestClient_ListDataprepStatus(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/dataprepstatus", http.StatusOK, map[string]interface{}{
		"data": []map[string]interface{}{
			{"status": "READY"},
		},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.ListDataprepStatus(context.Background())
	if err != nil {
		t.Fatalf("ListDataprepStatus failed: %v", err)
	}
	if resp == nil {
		t.Error("ListDataprepStatus returned nil")
	}
}

func TestClient_GetDataprepStatus(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/dataprepstatus/batch123", http.StatusOK, map[string]interface{}{
		"status": "READY",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.GetDataprepStatus(context.Background(), "batch123")
	if err != nil {
		t.Fatalf("GetDataprepStatus failed: %v", err)
	}
	if resp == nil {
		t.Error("GetDataprepStatus returned nil")
	}
}

// Fraud tests

func TestClient_GetFraudNotification(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/frauds/notification/acc123/txn123", http.StatusOK, map[string]interface{}{
		"notificationId": "fraud123",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.GetFraudNotification(context.Background(), "acc123", "txn123")
	if err != nil {
		t.Fatalf("GetFraudNotification failed: %v", err)
	}
	if resp == nil {
		t.Error("GetFraudNotification returned nil")
	}
}

func TestClient_UndoFraudNotification(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/frauds/notification/undo", http.StatusOK, nil)
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	err := client.UndoFraudNotification(context.Background(), &types.FraudNotificationUndoRequest{})
	if err != nil {
		t.Fatalf("UndoFraudNotification failed: %v", err)
	}
}

// HCE tests

func TestClient_UnprovisionHCE(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/hce/unprovision/device123", http.StatusOK, map[string]interface{}{
		"success": true,
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.UnprovisionHCE(context.Background(), "device123", &types.UnprovisionRequest{})
	if err != nil {
		t.Fatalf("UnprovisionHCE failed: %v", err)
	}
	if resp == nil {
		t.Error("UnprovisionHCE returned nil")
	}
}

// Inclusives tests

func TestClient_ListInclusiveTransactions(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/inclusives", http.StatusOK, map[string]interface{}{
		"transactions": []map[string]interface{}{
			{"transactionId": "txn123"},
		},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.ListInclusiveTransactions(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListInclusiveTransactions failed: %v", err)
	}
	if resp == nil {
		t.Error("ListInclusiveTransactions returned nil")
	}
}

func TestClient_UndoInclusiveTransaction(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/inclusives/txn123/undo", http.StatusOK, map[string]interface{}{
		"success": true,
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.UndoInclusiveTransaction(context.Background(), "txn123", &types.UndoInclusiveTransactionRequest{})
	if err != nil {
		t.Fatalf("UndoInclusiveTransaction failed: %v", err)
	}
	if resp == nil {
		t.Error("UndoInclusiveTransaction returned nil")
	}
}

// Cobranded tests

func TestClient_ListCobranded(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/cobranded", http.StatusOK, map[string]interface{}{
		"data": []map[string]interface{}{
			{"cobrandedId": "cob123"},
		},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.ListCobranded(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListCobranded failed: %v", err)
	}
	if resp == nil {
		t.Error("ListCobranded returned nil")
	}
}

func TestClient_UpdateCobrandedAcquirer(t *testing.T) {
	server := mockServer(t, http.MethodPut, "/cobranded/doc123/acquirer/1", http.StatusOK, nil)
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	err := client.UpdateCobrandedAcquirer(context.Background(), "doc123", 1, &types.MerchantVanRequest{})
	if err != nil {
		t.Fatalf("UpdateCobrandedAcquirer failed: %v", err)
	}
}

func TestClient_DeleteCobrandedAcquirer(t *testing.T) {
	server := mockServer(t, http.MethodDelete, "/cobranded/doc123/acquirer/1", http.StatusOK, nil)
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	err := client.DeleteCobrandedAcquirer(context.Background(), "doc123", 1)
	if err != nil {
		t.Fatalf("DeleteCobrandedAcquirer failed: %v", err)
	}
}

func TestClient_AddCobrandedSubacquirer(t *testing.T) {
	server := mockServer(t, http.MethodPut, "/cobranded/doc123/acquirer/1/cpfcnpjsubacquirer/cpf456", http.StatusCreated, nil)
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	err := client.AddCobrandedSubacquirer(context.Background(), "doc123", 1, "cpf456", &types.MerchantVanRequest{})
	if err != nil {
		t.Fatalf("AddCobrandedSubacquirer failed: %v", err)
	}
}

func TestClient_DeleteCobrandedSubacquirer(t *testing.T) {
	server := mockServer(t, http.MethodDelete, "/cobranded/doc123/acquirer/1/cpfcnpjsubacquirer/cpf456", http.StatusOK, nil)
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	err := client.DeleteCobrandedSubacquirer(context.Background(), "doc123", 1, "cpf456")
	if err != nil {
		t.Fatalf("DeleteCobrandedSubacquirer failed: %v", err)
	}
}
