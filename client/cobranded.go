package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// CreateCobranded creates a new cobranded merchant establishment.
func (c *Client) CreateCobranded(ctx context.Context, req *types.MerchantVanRequest) (*types.MerchantVanResponse, error) {
	var resp types.MerchantVanResponse
	if err := c.request(ctx, http.MethodPost, "/cobranded", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListCobrandedParams represents parameters for listing cobranded merchants.
type ListCobrandedParams struct {
	Document           string
	AcquirerID         int
	SubAcquirerCPFCNPJ string
}

// ListCobranded retrieves cobranded merchants based on filters.
func (c *Client) ListCobranded(ctx context.Context, params *ListCobrandedParams) (*types.MerchantVanResponse, error) {
	query := url.Values{}
	if params != nil {
		if params.Document != "" {
			query.Set("document", params.Document)
		}
		if params.AcquirerID > 0 {
			query.Set("acquirerId", strconv.Itoa(params.AcquirerID))
		}
		if params.SubAcquirerCPFCNPJ != "" {
			query.Set("cpfCnpjSubAcquirer", params.SubAcquirerCPFCNPJ)
		}
	}

	path := "/cobranded"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var resp types.MerchantVanResponse
	if err := c.request(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteCobranded removes a cobranded merchant by document.
func (c *Client) DeleteCobranded(ctx context.Context, document string) error {
	return c.request(ctx, http.MethodDelete, fmt.Sprintf("/cobranded/%s", pathParam(document)), nil, nil)
}

// UpdateCobrandedAcquirer updates the acquirer configuration for a cobranded merchant.
func (c *Client) UpdateCobrandedAcquirer(ctx context.Context, document string, acquirerID int, req *types.MerchantVanRequest) error {
	path := fmt.Sprintf("/cobranded/%s/acquirer/%d", pathParam(document), acquirerID)
	return c.request(ctx, http.MethodPut, path, req, nil)
}

// DeleteCobrandedAcquirer removes the acquirer association from a cobranded merchant.
func (c *Client) DeleteCobrandedAcquirer(ctx context.Context, document string, acquirerID int) error {
	path := fmt.Sprintf("/cobranded/%s/acquirer/%d", pathParam(document), acquirerID)
	return c.request(ctx, http.MethodDelete, path, nil, nil)
}

// AddCobrandedSubacquirer adds a subacquirer to a cobranded merchant.
func (c *Client) AddCobrandedSubacquirer(ctx context.Context, document string, acquirerID int, cpfCnpjSubacquirer string, req *types.MerchantVanRequest) error {
	path := fmt.Sprintf("/cobranded/%s/acquirer/%d/cpfcnpjsubacquirer/%s", pathParam(document), acquirerID, pathParam(cpfCnpjSubacquirer))
	return c.request(ctx, http.MethodPut, path, req, nil)
}

// DeleteCobrandedSubacquirer removes a subacquirer from a cobranded merchant.
func (c *Client) DeleteCobrandedSubacquirer(ctx context.Context, document string, acquirerID int, cpfCnpjSubacquirer string) error {
	path := fmt.Sprintf("/cobranded/%s/acquirer/%d/cpfcnpjsubacquirer/%s", pathParam(document), acquirerID, pathParam(cpfCnpjSubacquirer))
	return c.request(ctx, http.MethodDelete, path, nil, nil)
}
