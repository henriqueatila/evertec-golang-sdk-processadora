package client

import (
	"context"
	"net/http"
	"testing"
)

func TestClient_GetHealthStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name: "success: get health status - all healthy",
			mockResponse: map[string]any{
				"status": "healthy",
				"services": map[string]any{
					"database":     "healthy",
					"cache":        "healthy",
					"messageQueue": "healthy",
					"externalAPI":  "healthy",
				},
				"timestamp": "2024-01-15T10:30:00Z",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: get health status - degraded",
			mockResponse: map[string]any{
				"status": "degraded",
				"services": map[string]any{
					"database":     "healthy",
					"cache":        "degraded",
					"messageQueue": "healthy",
				},
				"timestamp": "2024-01-15T10:30:00Z",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: get health status - unhealthy",
			mockResponse: map[string]any{
				"status": "unhealthy",
				"services": map[string]any{
					"database": "unhealthy",
					"cache":    "healthy",
				},
				"timestamp": "2024-01-15T10:30:00Z",
			},
			mockStatus: http.StatusServiceUnavailable,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/status/health",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetHealthStatus(context.Background())

			if tt.wantErr {
				AssertError(t, err, "GetHealthStatus")
			} else {
				AssertNoError(t, err, "GetHealthStatus")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}
