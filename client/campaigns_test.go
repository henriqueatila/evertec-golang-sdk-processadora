package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

func TestClient_GetCampaign(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/campaigns/camp123", http.StatusOK, map[string]interface{}{
		"campaignId": "camp123",
		"name":       "Test Campaign",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.GetCampaign(context.Background(), "camp123")
	if err != nil {
		t.Fatalf("GetCampaign failed: %v", err)
	}
	if resp == nil {
		t.Error("GetCampaign returned nil")
	}
}

func TestClient_UpdateCampaign(t *testing.T) {
	server := mockServer(t, http.MethodPut, "/campaigns/camp123", http.StatusOK, map[string]interface{}{
		"campaignId": "camp123",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.UpdateCampaign(context.Background(), "camp123", &types.CampaignObject{})
	if err != nil {
		t.Fatalf("UpdateCampaign failed: %v", err)
	}
	if resp == nil {
		t.Error("UpdateCampaign returned nil")
	}
}

func TestClient_UpdateCampaignStatus(t *testing.T) {
	server := mockServer(t, http.MethodPut, "/campaigns/camp123/status", http.StatusOK, nil)
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	err := client.UpdateCampaignStatus(context.Background(), "camp123", types.CampaignStatusEnabled)
	if err != nil {
		t.Fatalf("UpdateCampaignStatus failed: %v", err)
	}
}

func TestClient_GetCampaignAccounts(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/campaigns/camp123/accounts", http.StatusOK, map[string]interface{}{
		"accounts": []map[string]interface{}{
			{"accountId": "acc1"},
		},
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.GetCampaignAccounts(context.Background(), "camp123")
	if err != nil {
		t.Fatalf("GetCampaignAccounts failed: %v", err)
	}
	if resp == nil {
		t.Error("GetCampaignAccounts returned nil")
	}
}

func TestClient_GetAgent(t *testing.T) {
	server := mockServer(t, http.MethodGet, "/agents/agent123", http.StatusOK, map[string]interface{}{
		"agentId": "agent123",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.GetAgent(context.Background(), "agent123")
	if err != nil {
		t.Fatalf("GetAgent failed: %v", err)
	}
	if resp == nil {
		t.Error("GetAgent returned nil")
	}
}

func TestClient_UpdateAgent(t *testing.T) {
	server := mockServer(t, http.MethodPut, "/agents/agent123", http.StatusOK, map[string]interface{}{
		"agentId": "agent123",
	})
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	resp, err := client.UpdateAgent(context.Background(), "agent123", &types.AgentObject{})
	if err != nil {
		t.Fatalf("UpdateAgent failed: %v", err)
	}
	if resp == nil {
		t.Error("UpdateAgent returned nil")
	}
}

func TestClient_UpdateAgentStatus(t *testing.T) {
	server := mockServer(t, http.MethodPut, "/agents/agent123/status", http.StatusOK, nil)
	defer server.Close()

	client, _ := New(testConfig(server.URL))

	err := client.UpdateAgentStatus(context.Background(), "agent123", "ACTIVE")
	if err != nil {
		t.Fatalf("UpdateAgentStatus failed: %v", err)
	}
}

func TestClient_CreateAgent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.AgentObject
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
	}{
		{
			name: "success: create agent",
			request: &types.AgentObject{
				DocID:       "12345678000199",
				CompanyName: "Test Company",
				TradingName: "Test Agent",
				Address:     "Rua Test, 123",
				Phone:       "11999999999",
				Email:       "test@example.com",
				ContactName: "John Doe",
			},
			mockResponse: map[string]interface{}{
				"agentId":     "agent123",
				"companyName": "Test Company",
			},
			mockStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "success: create agent with minimal fields",
			request: &types.AgentObject{
				DocID:       "98765432000199",
				CompanyName: "Minimal Company",
				TradingName: "Minimal Agent",
				Address:     "Rua Minimal, 1",
				Phone:       "11888888888",
				Email:       "minimal@example.com",
				ContactName: "Jane Doe",
			},
			mockResponse: map[string]interface{}{
				"agentId": "agent456",
			},
			mockStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "error: bad request",
			request: &types.AgentObject{
				DocID: "",
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
				Path:     "/agents",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.CreateAgent(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "CreateAgent")
			} else {
				AssertNoError(t, err, "CreateAgent")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}
