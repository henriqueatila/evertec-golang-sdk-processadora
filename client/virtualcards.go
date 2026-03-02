package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// CreateAccountVirtualCard creates a virtual card for an account.
func (c *Client) CreateAccountVirtualCard(ctx context.Context, accountID string, req *types.CreateVirtualCardRequest) (*types.CreateVirtualCardResponse, error) {
	var resp types.CreateVirtualCardResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/virtualcards", pathParam(accountID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAccountVirtualCards lists all virtual cards for an account.
func (c *Client) ListAccountVirtualCards(ctx context.Context, accountID string) (*types.ListVirtualCardsResponse, error) {
	var resp types.ListVirtualCardsResponse
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/accounts/%s/virtualcards/list", pathParam(accountID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateCardVirtualCard creates a virtual card for a physical card.
func (c *Client) CreateCardVirtualCard(ctx context.Context, cardID string, req *types.CreateVirtualCardRequest) (*types.CreateVirtualCardResponse, error) {
	var resp types.CreateVirtualCardResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/cards/%s/virtualcards", pathParam(cardID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListCardVirtualCards lists all virtual cards for a physical card.
func (c *Client) ListCardVirtualCards(ctx context.Context, cardID string) (*types.ListVirtualCardsResponse, error) {
	var resp types.ListVirtualCardsResponse
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/cards/%s/virtualcards/list", pathParam(cardID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVirtualCard retrieves details of a specific virtual card.
func (c *Client) GetVirtualCard(ctx context.Context, vCardID string) (*types.GetVirtualCardResponse, error) {
	var resp types.GetVirtualCardResponse
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/virtualcards/%s", pathParam(vCardID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelVirtualCard cancels a virtual card.
func (c *Client) CancelVirtualCard(ctx context.Context, vCardID string, req *types.CancelVirtualCardRequest) (*types.CancelVirtualCardResponse, error) {
	var resp types.CancelVirtualCardResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/virtualcards/%s/cancel", pathParam(vCardID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ModifyVirtualCardCVV modifies the CVV of a virtual card.
func (c *Client) ModifyVirtualCardCVV(ctx context.Context, vCardID string) (*types.ModifyVirtualCardCVVResponse, error) {
	var resp types.ModifyVirtualCardCVVResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/virtualcards/%s/modify/cvv", pathParam(vCardID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
