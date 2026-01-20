package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// CreateAccount creates a new account.
func (c *Client) CreateAccount(ctx context.Context, req *types.CreateAccountRequest) (*types.CreateAccountResponse, error) {
	var resp types.CreateAccountResponse
	if err := c.request(ctx, http.MethodPost, "/accounts", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAccount retrieves an account by ID.
func (c *Client) GetAccount(ctx context.Context, accountID string) (*types.Account, error) {
	var resp types.Account
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/accounts/%s", accountID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateAccount updates an account.
func (c *Client) UpdateAccount(ctx context.Context, accountID string, req *types.UpdateAccountRequest) error {
	return c.request(ctx, http.MethodPut, fmt.Sprintf("/accounts/%s", accountID), req, nil)
}

// BlockAccount blocks an account.
func (c *Client) BlockAccount(ctx context.Context, accountID string, req *types.BlockAccountRequest) error {
	return c.request(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/blockAccount", accountID), req, nil)
}

// UnblockAccount unblocks an account.
func (c *Client) UnblockAccount(ctx context.Context, accountID string, req *types.UnblockAccountRequest) error {
	return c.request(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/unblockAccount", accountID), req, nil)
}

// ListAccounts retrieves accounts based on filters.
// Query parameters are properly URL-encoded to handle special characters.
func (c *Client) ListAccounts(ctx context.Context, params *types.ListAccountsParams) (*types.ListAccountsResponse, error) {
	path := "/accounts"
	if params != nil {
		query := url.Values{}
		addPaginationParams(query, params.Limit, params.StartingAfter, params.EndingBefore)
		if params.IdentityDocumentNumber != "" {
			query.Set("identity_document_number", params.IdentityDocumentNumber)
		}
		if params.FullName != "" {
			query.Set("full_name", params.FullName)
		}
		if params.PSProductCode > 0 {
			query.Set("ps_product_code", strconv.Itoa(params.PSProductCode))
		}
		if params.AccountStatus != "" {
			query.Set("account_status", params.AccountStatus)
		}
		if params.IssuerAccountID != "" {
			query.Set("issuer_id", params.IssuerAccountID)
		}
		if params.IncludedSince != "" {
			query.Set("included_since", params.IncludedSince)
		}
		if params.Sort != "" {
			query.Set("sort", params.Sort)
		}
		if params.First {
			query.Set("first", "true")
		}
		if encoded := query.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}

	var resp types.ListAccountsResponse
	if err := c.request(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAccountBalance retrieves the balance and credit limits for an account.
// GET /accounts/{accountId}/balance
// Reference: https://paysmart-api.gitlab.io/processadora/PT-br/docs/limits-api
func (c *Client) GetAccountBalance(ctx context.Context, accountID string) (*types.AccountBalance, error) {
	var resp types.AccountBalance
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/accounts/%s/balance", accountID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelAccount cancels an account.
func (c *Client) CancelAccount(ctx context.Context, accountID string, req *types.CancelAccountRequest) error {
	return c.request(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/cancelAccount", accountID), req, nil)
}

// ListAccountCards retrieves all cards for an account.
func (c *Client) ListAccountCards(ctx context.Context, accountID string) ([]types.CardDetails, error) {
	var resp []types.CardDetails
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/accounts/%s/cards", accountID), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetCreditVerificationStatus retrieves credit verification status for an account.
func (c *Client) GetCreditVerificationStatus(ctx context.Context, accountID string) (*types.CreditVerificationStatus, error) {
	var resp types.CreditVerificationStatus
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/accounts/%s/creditVerificationStatus", accountID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListScheduledTransactions retrieves scheduled transactions for an account.
func (c *Client) ListScheduledTransactions(ctx context.Context, accountID string) ([]types.ScheduledTransaction, error) {
	var resp []types.ScheduledTransaction
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/accounts/%s/findScheduledTransactions", accountID), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateCredit creates a credit transaction for an account.
func (c *Client) CreateCredit(ctx context.Context, accountID string, req *types.CreateCreditRequest) (*types.CreditTransaction, error) {
	var resp types.CreditTransaction
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/transactions/credits", accountID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListCredits retrieves all credit transactions for an account.
func (c *Client) ListCredits(ctx context.Context, accountID string) (*types.ListCreditsResponse, error) {
	var resp types.ListCreditsResponse
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/accounts/%s/transactions/credits", accountID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateDebit creates a debit transaction for an account.
func (c *Client) CreateDebit(ctx context.Context, accountID string, req *types.CreateDebitRequest) (*types.CreditTransaction, error) {
	var resp types.CreditTransaction
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/transactions/debits", accountID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
