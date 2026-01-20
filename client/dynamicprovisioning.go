package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// ListDynamicProvisioningCreditsParams represents parameters for listing dynamic provisioning credits.
type ListDynamicProvisioningCreditsParams struct {
	BeginningDate string // Required: Start date (YYYY-MM-DD format)
	EndDate       string // Required: End date (YYYY-MM-DD format)
}

// GetDynamicProvisioningBalance retrieves the balance of the dynamic provisioning account (conta colchão).
// GET /dynamicProvisioningAccount/balance
func (c *Client) GetDynamicProvisioningBalance(ctx context.Context) (*types.DynamicProvisioningBalance, error) {
	var resp types.DynamicProvisioningBalance
	if err := c.request(ctx, http.MethodGet, "/dynamicProvisioningAccount/balance", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListDynamicProvisioningCredits retrieves the list of credits/deposits made to the dynamic provisioning account.
// GET /dynamicProvisioningAccount/balance/credits
func (c *Client) ListDynamicProvisioningCredits(ctx context.Context, params *ListDynamicProvisioningCreditsParams) (*types.DynamicProvisioningCreditList, error) {
	query := url.Values{}
	if params != nil {
		if params.BeginningDate != "" {
			query.Set("beginning_date", params.BeginningDate)
		}
		if params.EndDate != "" {
			query.Set("end_date", params.EndDate)
		}
	}

	path := "/dynamicProvisioningAccount/balance/credits"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var resp types.DynamicProvisioningCreditList
	if err := c.request(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
