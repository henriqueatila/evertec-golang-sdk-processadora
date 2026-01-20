package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// ReportFileType represents a type of report file.
type ReportFileType string

const (
	ReportFileTypeCSV ReportFileType = "csv"
	ReportFileTypePDF ReportFileType = "pdf"
	ReportFileTypeTXT ReportFileType = "txt"
)

// ListReports lists available reports of a specific file type.
func (c *Client) ListReports(ctx context.Context, fileType ReportFileType) (*types.ReportListResult, error) {
	var resp types.ReportListResult
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/reports/list/%s", fileType), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetReportTemporaryURL retrieves a temporary URL for downloading a report.
func (c *Client) GetReportTemporaryURL(ctx context.Context, fileType ReportFileType, fileKey string) (*types.ReportTemporaryUrl, error) {
	var resp types.ReportTemporaryUrl
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/reports/temporaryUrl/%s/%s", fileType, fileKey), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
