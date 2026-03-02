package client

import (
	"context"
	"net/http"
	"testing"
)

func TestClient_ListReports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		fileType     ReportFileType
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:     "success: list CSV reports",
			fileType: ReportFileTypeCSV,
			mockResponse: map[string]any{
				"files": []map[string]any{
					{"key": "report_2024_01.csv", "size": 1024, "lastModified": "2024-01-15T10:00:00Z"},
					{"key": "report_2024_02.csv", "size": 2048, "lastModified": "2024-02-15T10:00:00Z"},
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:     "success: list PDF reports",
			fileType: ReportFileTypePDF,
			mockResponse: map[string]any{
				"files": []map[string]any{
					{"key": "statement_2024_01.pdf", "size": 5120},
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:     "success: list TXT reports",
			fileType: ReportFileTypeTXT,
			mockResponse: map[string]any{
				"files": []map[string]any{
					{"key": "log_2024_01.txt", "size": 512},
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:     "success: empty reports list",
			fileType: ReportFileTypeCSV,
			mockResponse: map[string]any{
				"files": []map[string]any{},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/reports/list/" + string(tt.fileType),
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.ListReports(context.Background(), tt.fileType)

			if tt.wantErr {
				AssertError(t, err, "ListReports")
			} else {
				AssertNoError(t, err, "ListReports")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_GetReportTemporaryURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		fileType     ReportFileType
		fileKey      string
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:     "success: get CSV report URL",
			fileType: ReportFileTypeCSV,
			fileKey:  "report_2024_01.csv",
			mockResponse: map[string]any{
				"temporaryUrl": "https://storage.example.com/reports/report_2024_01.csv?token=abc123&expires=1234567890",
				"expiresAt":    "2024-01-15T11:00:00Z",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:     "success: get PDF report URL",
			fileType: ReportFileTypePDF,
			fileKey:  "statement_2024_01.pdf",
			mockResponse: map[string]any{
				"temporaryUrl": "https://storage.example.com/reports/statement_2024_01.pdf?token=xyz789",
				"expiresAt":    "2024-01-15T11:00:00Z",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "error: report not found",
			fileType:     ReportFileTypeCSV,
			fileKey:      "nonexistent_report.csv",
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
				Path:     "/reports/temporaryUrl/" + string(tt.fileType) + "/" + tt.fileKey,
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetReportTemporaryURL(context.Background(), tt.fileType, tt.fileKey)

			if tt.wantErr {
				AssertError(t, err, "GetReportTemporaryURL")
			} else {
				AssertNoError(t, err, "GetReportTemporaryURL")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}
