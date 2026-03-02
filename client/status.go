package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// GetStatus retrieves the API status.
func (c *Client) GetStatus(ctx context.Context) (*types.GetStatusResponse, error) {
	var resp types.GetStatusResponse
	if err := c.request(ctx, http.MethodGet, "/status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetHealthStatus retrieves the API health status.
func (c *Client) GetHealthStatus(ctx context.Context) (*types.GetHealthStatusResponse, error) {
	var resp types.GetHealthStatusResponse
	if err := c.request(ctx, http.MethodGet, "/status/health", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListDataprepStatus lists dataprep statuses.
func (c *Client) ListDataprepStatus(ctx context.Context) (*types.ListDataprepStatusResponse, error) {
	var resp types.ListDataprepStatusResponse
	if err := c.request(ctx, http.MethodGet, "/dataprepstatus", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDataprepStatus retrieves the status of a specific dataprep file.
func (c *Client) GetDataprepStatus(ctx context.Context, filename string) (*types.GetDataprepStatusResponse, error) {
	var resp types.GetDataprepStatusResponse
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/dataprepstatus/%s", pathParam(filename)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
