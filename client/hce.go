package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// ProvisionHCE provisions a Host Card Emulation (HCE) card.
func (c *Client) ProvisionHCE(ctx context.Context, req *types.HCEProvisionRequest) (*types.HCEProvisionSuccessfully, error) {
	var resp types.HCEProvisionSuccessfully
	if err := c.request(ctx, http.MethodPost, "/hce/provision", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateHCECard creates an HCE card for an account.
func (c *Client) CreateHCECard(ctx context.Context, accountID string, req *types.CreateHCECardRequest) (*types.CreateHCECardSuccessfully, error) {
	var resp types.CreateHCECardSuccessfully
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/hce/%s/createHCECard", accountID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnprovisionHCE unprovisions an HCE card.
func (c *Client) UnprovisionHCE(ctx context.Context, issuerDeviceID string, req *types.UnprovisionRequest) (*types.UnprovisionSuccessfully, error) {
	var resp types.UnprovisionSuccessfully
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/hce/unprovision/%s", issuerDeviceID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
