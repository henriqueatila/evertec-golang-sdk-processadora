package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

func TestClient_ListDisputes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.ListDisputesRequest
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
	}{
		{
			name:    "success: list disputes without filters",
			request: nil,
			mockResponse: map[string]interface{}{
				"hasMore": true,
				"disputes": []map[string]interface{}{
					{"disputeId": "disp001", "disputeStatus": "OPEN"},
					{"disputeId": "disp002", "disputeStatus": "RESOLVED"},
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: list with status filter",
			request: &types.ListDisputesRequest{
				DisputeStatus: "OPEN",
			},
			mockResponse: map[string]interface{}{
				"hasMore": false,
				"disputes": []map[string]interface{}{
					{"disputeId": "disp001", "disputeStatus": "OPEN"},
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: list with all filters",
			request: &types.ListDisputesRequest{
				DisputeStatus: "RESOLVED",
				BeginningDate: "2024-01-01",
				EndingDate:    "2024-12-31",
				Limit:         20,
				StartingAfter: "disp100",
			},
			mockResponse: map[string]interface{}{
				"hasMore":  false,
				"disputes": []map[string]interface{}{},
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
				Path:     "/disputes",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.ListDisputes(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "ListDisputes")
			} else {
				AssertNoError(t, err, "ListDisputes")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_CreateDispute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.CreateDisputeRequest
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
	}{
		{
			name: "success: create dispute",
			request: &types.CreateDisputeRequest{
				TransactionID:      "tx123",
				AccountID:          "acc123",
				IssuerDisputeID:    "issuer-disp-001",
				DisputeCode:        "FRAUD",
				DisputeTextMessage: "Unauthorized transaction",
			},
			mockResponse: map[string]interface{}{
				"resultData": map[string]interface{}{
					"resultCode":        0,
					"resultDescription": "Disputa criada com sucesso!",
					"issuerRequestId":   "issuer-disp-001",
					"psResponseId":      "ps-resp-001",
				},
				"dispute": map[string]interface{}{
					"disputeId":     "disp001",
					"disputeStatus": "OPEN",
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "error: bad request",
			request:      &types.CreateDisputeRequest{},
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
				Path:     "/disputes",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.CreateDispute(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "CreateDispute")
			} else {
				AssertNoError(t, err, "CreateDispute")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_GetDispute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		disputeID    string
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
	}{
		{
			name:      "success: get dispute",
			disputeID: "disp123",
			mockResponse: map[string]interface{}{
				"disputeId":     "disp123",
				"disputeStatus": "OPEN",
				"disputeType":   "chargeback",
				"currentStage":  "chargeback",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "error: dispute not found",
			disputeID:    "nonexistent",
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
				Path:     "/disputes/" + tt.disputeID,
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetDispute(context.Background(), tt.disputeID)

			if tt.wantErr {
				AssertError(t, err, "GetDispute")
			} else {
				AssertNoError(t, err, "GetDispute")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_CancelDispute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		disputeID  string
		request    *types.CancelDisputeRequest
		mockStatus int
		wantErr    bool
	}{
		{
			name:      "success: cancel dispute",
			disputeID: "disp123",
			request: &types.CancelDisputeRequest{
				Reason: "Customer request",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "error: dispute not found",
			disputeID:  "nonexistent",
			request:    &types.CancelDisputeRequest{},
			mockStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/disputes/" + tt.disputeID + "/undo",
				Status:   tt.mockStatus,
				Response: nil,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			err := client.CancelDispute(context.Background(), tt.disputeID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "CancelDispute")
			} else {
				AssertNoError(t, err, "CancelDispute")
			}
		})
	}
}

func TestClient_AttachDisputeDocument(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		disputeID    string
		request      *types.DisputeDocumentRequest
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
	}{
		{
			name:      "success: attach document",
			disputeID: "disp123",
			request: &types.DisputeDocumentRequest{
				IssuerDisputeDocumentID: "doc-001",
				Document: types.DisputeDocument{
					DocumentName: "evidence.pdf",
					Document:     "base64encodeddata",
				},
			},
			mockResponse: map[string]interface{}{
				"resultData": map[string]interface{}{
					"resultCode":        0,
					"resultDescription": "Success",
				},
				"dispute": map[string]interface{}{
					"disputeId": "disp123",
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "success: attach document with source audit",
			disputeID: "disp456",
			request: &types.DisputeDocumentRequest{
				IssuerDisputeDocumentID: "doc-002",
				Document: types.DisputeDocument{
					DocumentName:        "evidence.png",
					DocumentDescription: "Screenshot evidence",
					Document:            "base64encodedimagedata",
				},
				SourceAudit: &types.SourceAudit{
					OperatorID: "user123",
				},
			},
			mockResponse: map[string]interface{}{
				"resultData": map[string]interface{}{
					"resultCode": 0,
				},
				"dispute": map[string]interface{}{
					"disputeId": "disp456",
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "error: dispute not found",
			disputeID: "nonexistent",
			request: &types.DisputeDocumentRequest{
				Document: types.DisputeDocument{
					DocumentName: "evidence.pdf",
					Document:     "data",
				},
			},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
		{
			name:      "error: bad request",
			disputeID: "disp123",
			request: &types.DisputeDocumentRequest{
				Document: types.DisputeDocument{},
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
				Path:     "/disputes/" + tt.disputeID + "/documents",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.AttachDisputeDocument(context.Background(), tt.disputeID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "AttachDisputeDocument")
			} else {
				AssertNoError(t, err, "AttachDisputeDocument")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_RespondToDispute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		disputeID    string
		request      *types.DisputeResponseRequest
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
	}{
		{
			name:      "success: accept dispute ruling",
			disputeID: "disp123",
			request: &types.DisputeResponseRequest{
				IssuerDisputeResponseID:    "resp-001",
				Accept:                     true,
				DisputeResponseTextMessage: "Accepted",
				WillAddDocuments:           false,
			},
			mockResponse: map[string]interface{}{
				"resultData": map[string]interface{}{
					"resultCode":        0,
					"resultDescription": "Success",
				},
				"dispute": map[string]interface{}{
					"disputeId": "disp123",
					"status":    "ACCEPTED",
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "success: reject and continue to next phase",
			disputeID: "disp456",
			request: &types.DisputeResponseRequest{
				IssuerDisputeResponseID:    "resp-002",
				Accept:                     false,
				DisputeResponseTextMessage: "Rejected, continuing to arbitration",
				WillAddDocuments:           true,
				SourceAudit: &types.SourceAudit{
					OperatorID: "user123",
				},
			},
			mockResponse: map[string]interface{}{
				"resultData": map[string]interface{}{
					"resultCode": 0,
				},
				"dispute": map[string]interface{}{
					"disputeId": "disp456",
					"status":    "ARBITRATION",
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:      "error: dispute not found",
			disputeID: "nonexistent",
			request: &types.DisputeResponseRequest{
				Accept:                     true,
				DisputeResponseTextMessage: "Test",
			},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
		{
			name:      "error: invalid response",
			disputeID: "disp123",
			request: &types.DisputeResponseRequest{
				Accept:                     false,
				DisputeResponseTextMessage: "",
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
				Path:     "/disputes/" + tt.disputeID + "/response",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.RespondToDispute(context.Background(), tt.disputeID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "RespondToDispute")
			} else {
				AssertNoError(t, err, "RespondToDispute")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}
