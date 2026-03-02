package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// CreateTravelNotice creates a travel notice for an account.
func (c *Client) CreateTravelNotice(ctx context.Context, accountID string, req *types.NewTravelNotice) (*types.TravelNoticeCreatedSuccessfully, error) {
	var resp types.TravelNoticeCreatedSuccessfully
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/travelNotice", pathParam(accountID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListTravelNotices lists all travel notices for an account.
func (c *Client) ListTravelNotices(ctx context.Context, accountID string) ([]types.TravelNotice, error) {
	var resp []types.TravelNotice
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/accounts/%s/travelNotice", pathParam(accountID)), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTravelNotice retrieves a specific travel notice.
func (c *Client) GetTravelNotice(ctx context.Context, accountID, travelNoticeID string) (*types.TravelNotice, error) {
	var resp types.TravelNotice
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/accounts/%s/travelNotice/%s", pathParam(accountID), pathParam(travelNoticeID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateTravelNotice updates a travel notice.
func (c *Client) UpdateTravelNotice(ctx context.Context, accountID, travelNoticeID string, req *types.UpdateTravelNotice) (*types.TravelNoticeCreatedSuccessfully, error) {
	var resp types.TravelNoticeCreatedSuccessfully
	if err := c.request(ctx, http.MethodPut, fmt.Sprintf("/accounts/%s/travelNotice/%s", pathParam(accountID), pathParam(travelNoticeID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
