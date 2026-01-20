package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// CreateFraudNotification creates a fraud notification.
func (c *Client) CreateFraudNotification(ctx context.Context, req *types.FraudNotificationRequest) (*types.FraudNotification, error) {
	var resp types.FraudNotification
	if err := c.request(ctx, http.MethodPost, "/frauds/notification", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListFraudNotifications lists all fraud notifications.
func (c *Client) ListFraudNotifications(ctx context.Context) ([]types.FraudNotification, error) {
	var resp []types.FraudNotification
	if err := c.request(ctx, http.MethodGet, "/frauds/notification", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetFraudNotification retrieves a specific fraud notification.
func (c *Client) GetFraudNotification(ctx context.Context, accountID, transactionID string) (*types.FraudNotification, error) {
	var resp types.FraudNotification
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/frauds/notification/%s/%s", accountID, transactionID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UndoFraudNotification undoes a fraud notification.
func (c *Client) UndoFraudNotification(ctx context.Context, req *types.FraudNotificationUndoRequest) error {
	return c.request(ctx, http.MethodPost, "/frauds/notification/undo", req, nil)
}
