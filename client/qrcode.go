package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// ParseQRCode parses a QR code to extract payment information.
func (c *Client) ParseQRCode(ctx context.Context, req *types.ParseQrCodeParams) (*types.ParseQrCodeSuccess, error) {
	var resp types.ParseQrCodeSuccess
	if err := c.request(ctx, http.MethodPost, "/qrcode/parse", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SimplePayment processes a simple payment via QR code.
func (c *Client) SimplePayment(ctx context.Context, req *types.SimplePaymentCardParams) (*types.SimplePaymentResponse, error) {
	var resp types.SimplePaymentResponse
	if err := c.request(ctx, http.MethodPost, "/qrcode/simplePayment", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TransactionCallback sends a callback for a transaction.
func (c *Client) TransactionCallback(ctx context.Context, transactionID string, req *types.TransactionCallbackRequest) (*types.TransactionCallbackResponse, error) {
	var resp types.TransactionCallbackResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/callbacks/transactions/%s", transactionID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
