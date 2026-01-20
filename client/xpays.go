package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// GetDeviceTokensFromCard retrieves all device tokens associated with a card.
// GET /cards/{cardId}/deviceTokens
// Returns array of DeviceToken directly per official API spec.
func (c *Client) GetDeviceTokensFromCard(ctx context.Context, cardID string) ([]types.DeviceToken, error) {
	var resp []types.DeviceToken
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/cards/%s/deviceTokens", cardID), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// MigrateDeviceTokens migrates device tokens to a new card.
// POST /cards/{cardId}/deviceTokens/migrate
func (c *Client) MigrateDeviceTokens(ctx context.Context, cardID string, req *types.MigrateTokensRequest) (*types.MigrateTokensResponse, error) {
	var resp types.MigrateTokensResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/cards/%s/deviceTokens/migrate", cardID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SuspendOrResumeDeviceTokensByCard suspends or resumes all device tokens for a card.
// POST /cards/{cardId}/deviceTokens/suspendOrResume
// Returns array of DeviceToken directly per official API spec.
func (c *Client) SuspendOrResumeDeviceTokensByCard(ctx context.Context, cardID string, req *types.SuspendResumeRequest) ([]types.DeviceToken, error) {
	var resp []types.DeviceToken
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/cards/%s/deviceTokens/suspendOrResume", cardID), req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// TerminateDeviceTokensByCard permanently disables all device tokens for a card.
// POST /cards/{cardId}/deviceTokens/terminate
// Returns array of DeviceToken directly per official API spec.
func (c *Client) TerminateDeviceTokensByCard(ctx context.Context, cardID string, req *types.TerminateTokensRequest) ([]types.DeviceToken, error) {
	var resp []types.DeviceToken
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/cards/%s/deviceTokens/terminate", cardID), req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetEncryptedCard obtains encrypted card data for provisioning.
// POST /cards/{cardId}/encryptedCard
func (c *Client) GetEncryptedCard(ctx context.Context, cardID string, req *types.EncryptedCardRequest) (*types.EncryptedCardResponse, error) {
	var resp types.EncryptedCardResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/cards/%s/encryptedCard", cardID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOpaqueCard retrieves encrypted card data for a card.
// GET /cards/{cardId}/opaqueCard
func (c *Client) GetOpaqueCard(ctx context.Context, cardID string) (*types.OpaqueCardResponse, error) {
	var resp types.OpaqueCardResponse
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/cards/%s/opaqueCard", cardID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ActivateDeviceToken activates a device token that is not yet active.
// POST /deviceTokens/{deviceTokenId}/activate
// Returns DeviceToken directly per official API spec.
func (c *Client) ActivateDeviceToken(ctx context.Context, deviceTokenID string, req *types.ActivateTokenRequest) (*types.DeviceToken, error) {
	var resp types.DeviceToken
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/deviceTokens/%s/activate", deviceTokenID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SuspendOrResumeDeviceToken suspends or resumes a specific device token.
// POST /deviceTokens/{deviceTokenId}/suspendOrResume
// Returns DeviceToken directly per official API spec.
func (c *Client) SuspendOrResumeDeviceToken(ctx context.Context, deviceTokenID string, req *types.SuspendResumeRequest) (*types.DeviceToken, error) {
	var resp types.DeviceToken
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/deviceTokens/%s/suspendOrResume", deviceTokenID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TerminateDeviceToken permanently disables a specific device token.
// POST /deviceTokens/{deviceTokenId}/terminate
// Returns DeviceToken directly per official API spec.
func (c *Client) TerminateDeviceToken(ctx context.Context, deviceTokenID string, req *types.TerminateTokenRequest) (*types.DeviceToken, error) {
	var resp types.DeviceToken
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/deviceTokens/%s/terminate", deviceTokenID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
