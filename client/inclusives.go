package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// CreateInclusiveTransaction creates a new inclusive transaction (TE10/TE20).
func (c *Client) CreateInclusiveTransaction(ctx context.Context, req *types.InclusiveTransactionRequest) (*types.InclusiveTransactionCreationSuccess, error) {
	var resp types.InclusiveTransactionCreationSuccess
	if err := c.request(ctx, http.MethodPost, "/inclusives", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListInclusiveTransactionsParams represents parameters for listing inclusive transactions.
type ListInclusiveTransactionsParams struct {
	Limit         int
	StartingAfter string
	EndingBefore  string
	BeginningDate string
	EndingDate    string
}

// ListInclusiveTransactions retrieves all inclusive transactions for an issuer.
func (c *Client) ListInclusiveTransactions(ctx context.Context, params *ListInclusiveTransactionsParams) (*types.InclusiveTransactionsListResult, error) {
	query := url.Values{}
	if params != nil {
		addPaginationParams(query, params.Limit, params.StartingAfter, params.EndingBefore)
		if params.BeginningDate != "" {
			query.Set("beginning_date", params.BeginningDate)
		}
		if params.EndingDate != "" {
			query.Set("ending_date", params.EndingDate)
		}
	}

	path := "/inclusives"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	var resp types.InclusiveTransactionsListResult
	if err := c.request(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetInclusiveTransaction retrieves a specific inclusive transaction by ID.
func (c *Client) GetInclusiveTransaction(ctx context.Context, inclusiveTransactionID string) (*types.InclusiveTransaction, error) {
	var resp types.InclusiveTransaction
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/inclusives/%s", inclusiveTransactionID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UndoInclusiveTransaction undoes an inclusive transaction.
func (c *Client) UndoInclusiveTransaction(ctx context.Context, inclusiveTransactionID string, req *types.UndoInclusiveTransactionRequest) (*types.InclusiveTransactionUndoSuccess, error) {
	var resp types.InclusiveTransactionUndoSuccess
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/inclusives/%s/undo", inclusiveTransactionID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
