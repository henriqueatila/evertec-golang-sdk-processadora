package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// ListAccountTransactions lists transactions for an account.
func (c *Client) ListAccountTransactions(ctx context.Context, accountID string, req *types.ListTransactionsRequest) (*types.ListTransactionsResponse, error) {
	path := fmt.Sprintf("/accounts/%s/transactions", pathParam(accountID))
	if req != nil {
		path = path + "?" + buildTransactionQuery(req)
	}

	var resp types.ListTransactionsResponse
	if err := c.request(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTransaction retrieves a transaction by ID.
func (c *Client) GetTransaction(ctx context.Context, transactionID string) (*types.Transaction, error) {
	var resp types.Transaction
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/transactions/%s", pathParam(transactionID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAllTransactions lists all transactions with filters (not account/card specific).
func (c *Client) ListAllTransactions(ctx context.Context, req *types.ListAllTransactionsRequest) (*types.ListTransactionsResponse, error) {
	path := "/transactions"
	if req != nil {
		path = path + "?" + buildAllTransactionQuery(req)
	}

	var resp types.ListTransactionsResponse
	if err := c.request(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func buildTransactionQuery(req *types.ListTransactionsRequest) string {
	q := url.Values{}
	addPaginationParams(q, req.Limit, req.StartingAfter, req.EndingBefore)
	if req.BeginningDate != "" {
		q.Set("beginning_date", req.BeginningDate)
	}
	if req.EndingDate != "" {
		q.Set("ending_date", req.EndingDate)
	}
	if req.TransactionType != "" {
		q.Set("transaction_type", req.TransactionType)
	}
	if req.TransactionStatus != "" {
		q.Set("transaction_status", req.TransactionStatus)
	}
	if req.TransactionApproved != nil {
		q.Set("transaction_approved", strconv.FormatBool(*req.TransactionApproved))
	}
	if req.TransactionDenialCode != "" {
		q.Set("transaction_denial_code", req.TransactionDenialCode)
	}
	if req.MinimumAmount > 0 {
		q.Set("minimum_amount", strconv.FormatInt(req.MinimumAmount, 10))
	}
	if req.MaxAmount > 0 {
		q.Set("max_amount", strconv.FormatInt(req.MaxAmount, 10))
	}
	if req.TransactionEntryMode != "" {
		q.Set("transaction_mode", req.TransactionEntryMode)
	}
	if req.Ordinated != "" {
		q.Set("ordinated", req.Ordinated)
	}
	return q.Encode()
}

func buildAllTransactionQuery(req *types.ListAllTransactionsRequest) string {
	q := url.Values{}
	addPaginationParams(q, req.Limit, req.StartingAfter, req.EndingBefore)
	if req.BeginningDate != "" {
		q.Set("beginning_date", req.BeginningDate)
	}
	if req.EndingDate != "" {
		q.Set("ending_date", req.EndingDate)
	}
	if req.TransactionType != "" {
		q.Set("transaction_type", req.TransactionType)
	}
	if req.TransactionStatus != "" {
		q.Set("transaction_status", req.TransactionStatus)
	}
	if req.TransactionApproved != nil {
		q.Set("transaction_approved", strconv.FormatBool(*req.TransactionApproved))
	}
	if req.TransactionDenialCode != "" {
		q.Set("transaction_denial_code", req.TransactionDenialCode)
	}
	if req.MinimumAmount > 0 {
		q.Set("minimum_amount", strconv.FormatInt(req.MinimumAmount, 10))
	}
	if req.MaxAmount > 0 {
		q.Set("max_amount", strconv.FormatInt(req.MaxAmount, 10))
	}
	if req.TransactionEntryMode != "" {
		q.Set("transaction_mode", req.TransactionEntryMode)
	}
	return q.Encode()
}
