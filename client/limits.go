package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// UpdateMaxCreditLimits updates the maximum credit limits for an account.
// This endpoint is part of the Limits API (Pós-Pago).
func (c *Client) UpdateMaxCreditLimits(ctx context.Context, accountID string, req *types.MaxCreditLimitsRequest) (*types.MaxCreditLimitsResponse, error) {
	var resp types.MaxCreditLimitsResponse
	if err := c.request(ctx, http.MethodPut, fmt.Sprintf("/accounts/%s/maxCreditLimits", pathParam(accountID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ChangeUsableCreditLimits changes the usable credit limits for an account.
// This endpoint is part of the Limits API (Pós-Pago).
func (c *Client) ChangeUsableCreditLimits(ctx context.Context, accountID string, req *types.ChangeUsableCreditLimitsRequest) (*types.ChangeUsableCreditLimitsResponse, error) {
	var resp types.ChangeUsableCreditLimitsResponse
	if err := c.request(ctx, http.MethodPut, fmt.Sprintf("/accounts/%s/changeUsableCreditLimits", pathParam(accountID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
