package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

func TestClient_GetClosedStatements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.ClosedStatementsRequest
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, r *http.Request)
	}{
		{
			name:    "success: get closed statements without params",
			request: nil,
			mockResponse: map[string]interface{}{
				"statements": []map[string]interface{}{
					{
						"statementId":   "stmt001",
						"accountId":     "acc001",
						"closingDate":   "2024-01-15",
						"dueDate":       "2024-02-10",
						"totalAmount":   15000,
						"minimumAmount": 1500,
					},
					{
						"statementId":   "stmt002",
						"accountId":     "acc002",
						"closingDate":   "2024-01-15",
						"dueDate":       "2024-02-10",
						"totalAmount":   25000,
						"minimumAmount": 2500,
					},
				},
				"hasMore": true,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: get closed statements with closing date query",
			request: &types.ClosedStatementsRequest{
				ClosingDateQuery: "2024-01-15",
			},
			mockResponse: map[string]interface{}{
				"statements": []map[string]interface{}{
					{
						"statementId": "stmt001",
						"closingDate": "2024-01-15",
					},
				},
				"hasMore": false,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("closing_date_query") != "2024-01-15" {
					t.Errorf("expected closing_date_query=2024-01-15, got %s", r.URL.Query().Get("closing_date_query"))
				}
			},
		},
		{
			name: "success: get closed statements with starting after cursor",
			request: &types.ClosedStatementsRequest{
				StartingAfter: "stmt100",
			},
			mockResponse: map[string]interface{}{
				"statements": []map[string]interface{}{
					{
						"statementId": "stmt101",
						"closingDate": "2024-01-15",
					},
				},
				"hasMore": true,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("starting_after") != "stmt100" {
					t.Errorf("expected starting_after=stmt100, got %s", r.URL.Query().Get("starting_after"))
				}
			},
		},
		{
			name: "success: get closed statements with all params",
			request: &types.ClosedStatementsRequest{
				ClosingDateQuery: "2024-01-15",
				StartingAfter:    "stmt100",
			},
			mockResponse: map[string]interface{}{
				"statements": []map[string]interface{}{},
				"hasMore":    false,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("closing_date_query") != "2024-01-15" {
					t.Errorf("expected closing_date_query=2024-01-15, got %s", r.URL.Query().Get("closing_date_query"))
				}
				if r.URL.Query().Get("starting_after") != "stmt100" {
					t.Errorf("expected starting_after=stmt100, got %s", r.URL.Query().Get("starting_after"))
				}
			},
		},
		{
			name: "success: empty statements list",
			request: &types.ClosedStatementsRequest{
				ClosingDateQuery: "2020-01-01",
			},
			mockResponse: map[string]interface{}{
				"statements": []map[string]interface{}{},
				"hasMore":    false,
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
		{
			name:         "error: bad request",
			request:      &types.ClosedStatementsRequest{ClosingDateQuery: "invalid-date"},
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
				Method:   http.MethodGet,
				Path:     "/accounts/statements/closed",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						tt.validateReq(t, r)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetClosedStatements(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "GetClosedStatements")
			} else {
				AssertNoError(t, err, "GetClosedStatements")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}
