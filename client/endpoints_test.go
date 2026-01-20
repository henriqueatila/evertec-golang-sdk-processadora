package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// ==================== INSTALLMENTS TESTS ====================

func TestClient_SimulateInstallment(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/accounts/acc123/installmentSimulation", http.StatusOK, map[string]interface{}{
		"options": []map[string]interface{}{
			{"installments": 3, "installmentValue": 3333},
			{"installments": 6, "installmentValue": 1667},
		},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.InstallmentSimulationBody{
		Amount: 10000,
	}

	resp, err := client.SimulateInstallment(context.Background(), "acc123", req)
	if err != nil {
		t.Fatalf("SimulateInstallment failed: %v", err)
	}

	if resp == nil {
		t.Error("SimulateInstallment returned nil response")
	}
}

func TestClient_CreateInstallmentRequest(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/accounts/acc123/installmentRequest", http.StatusOK, map[string]interface{}{
		"resultCode": 0,
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.InstallmentRequest{
		NumberOfInstallments: 6,
	}

	resp, err := client.CreateInstallmentRequest(context.Background(), "acc123", req)
	if err != nil {
		t.Fatalf("CreateInstallmentRequest failed: %v", err)
	}

	if resp == nil {
		t.Error("CreateInstallmentRequest returned nil response")
	}
}

// ==================== VIRTUAL CARDS TESTS ====================

func TestClient_CreateAccountVirtualCard(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/accounts/acc123/virtualcards", http.StatusCreated, map[string]interface{}{
		"vCardId": "vcard123",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.CreateVirtualCardRequest{
		Constraints: &types.VirtualCardConstraints{
			CurrencyCode: "BRL",
			MaxAmount:    "100000",
		},
	}

	resp, err := client.CreateAccountVirtualCard(context.Background(), "acc123", req)
	if err != nil {
		t.Fatalf("CreateAccountVirtualCard failed: %v", err)
	}

	if resp == nil {
		t.Error("CreateAccountVirtualCard returned nil response")
	}
}

func TestClient_ListAccountVirtualCards(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/accounts/acc123/virtualcards/list", http.StatusOK, map[string]interface{}{
		"virtualCards": []map[string]interface{}{
			{"vCardId": "vcard1"},
			{"vCardId": "vcard2"},
		},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.ListAccountVirtualCards(context.Background(), "acc123")
	if err != nil {
		t.Fatalf("ListAccountVirtualCards failed: %v", err)
	}

	if resp == nil {
		t.Error("ListAccountVirtualCards returned nil response")
	}
}

func TestClient_GetVirtualCard(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/virtualcards/vcard123", http.StatusOK, map[string]interface{}{
		"vCardId": "vcard123",
		"status":  "ACTIVE",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.GetVirtualCard(context.Background(), "vcard123")
	if err != nil {
		t.Fatalf("GetVirtualCard failed: %v", err)
	}

	if resp == nil {
		t.Error("GetVirtualCard returned nil response")
	}
}

func TestClient_CancelVirtualCard(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/virtualcards/vcard123/cancel", http.StatusOK, map[string]interface{}{
		"status": "CANCELLED",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.CancelVirtualCardRequest{
		CancellationCode: 1,
		Reason:           "User request",
	}

	resp, err := client.CancelVirtualCard(context.Background(), "vcard123", req)
	if err != nil {
		t.Fatalf("CancelVirtualCard failed: %v", err)
	}

	if resp == nil {
		t.Error("CancelVirtualCard returned nil response")
	}
}

// ==================== HCE TESTS ====================

func TestClient_ProvisionHCE(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/hce/provision", http.StatusOK, map[string]interface{}{
		"resultCode": 0,
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.HCEProvisionRequest{
		Mtv: map[string]interface{}{
			"provData":       "encrypted_data",
			"issuerDeviceId": "device123",
		},
	}

	resp, err := client.ProvisionHCE(context.Background(), req)
	if err != nil {
		t.Fatalf("ProvisionHCE failed: %v", err)
	}

	if resp == nil {
		t.Error("ProvisionHCE returned nil response")
	}
}

func TestClient_CreateHCECard(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/hce/acc123/createHCECard", http.StatusOK, map[string]interface{}{
		"resultCode": 0,
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.CreateHCECardRequest{}

	resp, err := client.CreateHCECard(context.Background(), "acc123", req)
	if err != nil {
		t.Fatalf("CreateHCECard failed: %v", err)
	}

	if resp == nil {
		t.Error("CreateHCECard returned nil response")
	}
}

// ==================== CAMPAIGNS/AGENTS TESTS ====================

func TestClient_CreateCampaign(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/campaigns", http.StatusCreated, map[string]interface{}{
		"campaignId": "camp123",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.CampaignObject{
		CampaignName: "Test Campaign",
		AgentID:      "agent123",
	}

	resp, err := client.CreateCampaign(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateCampaign failed: %v", err)
	}

	if resp == nil {
		t.Error("CreateCampaign returned nil response")
	}
}

func TestClient_ListCampaigns(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/campaigns", http.StatusOK, map[string]interface{}{
		"campaigns": []map[string]interface{}{
			{"campaignId": "camp1"},
			{"campaignId": "camp2"},
		},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.ListCampaigns(context.Background())
	if err != nil {
		t.Fatalf("ListCampaigns failed: %v", err)
	}

	if resp == nil {
		t.Error("ListCampaigns returned nil response")
	}
}

func TestClient_ListAgents(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/agents", http.StatusOK, []map[string]interface{}{
		{"agentId": "agent1"},
		{"agentId": "agent2"},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.ListAgents(context.Background())
	if err != nil {
		t.Fatalf("ListAgents failed: %v", err)
	}

	if resp == nil {
		t.Error("ListAgents returned nil response")
	}
}

// ==================== TRAVEL NOTICE TESTS ====================

func TestClient_CreateTravelNotice(t *testing.T) {
	server := mockServer(t, http.MethodPost, "/accounts/acc123/travelNotice", http.StatusCreated, map[string]interface{}{
		"travelNoticeId": "tn123",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	req := &types.NewTravelNotice{
		BeginDate:    "2024-07-01",
		EndDate:      "2024-07-15",
		CountryCodes: "USA",
		Cards:        []types.Card{},
	}

	resp, err := client.CreateTravelNotice(context.Background(), "acc123", req)
	if err != nil {
		t.Fatalf("CreateTravelNotice failed: %v", err)
	}

	if resp == nil {
		t.Error("CreateTravelNotice returned nil response")
	}
}

// ==================== HEALTH CHECK TESTS ====================

func TestClient_GetStatus(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/status", http.StatusOK, map[string]interface{}{
		"resultCode": 0,
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if resp == nil {
		t.Error("GetStatus returned nil response")
	}
}
