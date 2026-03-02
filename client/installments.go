package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// SimulateInstallment simulates installment payment options for an account.
func (c *Client) SimulateInstallment(ctx context.Context, accountID string, req *types.InstallmentSimulationBody) (*types.InstallmentSimulationResult, error) {
	var resp types.InstallmentSimulationResult
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/installmentSimulation", pathParam(accountID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetInstallmentOptions retrieves available installment options for an account.
func (c *Client) GetInstallmentOptions(ctx context.Context, accountID string) (*types.InstallmentSimulationResult, error) {
	var resp types.InstallmentSimulationResult
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/accounts/%s/installmentSimulation", pathParam(accountID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateInstallmentRequest creates an installment request for an account.
func (c *Client) CreateInstallmentRequest(ctx context.Context, accountID string, req *types.InstallmentRequest) (*types.ResultData, error) {
	var resp types.ResultData
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/installmentRequest", pathParam(accountID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SimulateAdvancePayment simulates advance payment for installments on a transaction.
func (c *Client) SimulateAdvancePayment(ctx context.Context, accountID, transactionID string, req *types.AdvancePaymentRequest) (*types.InstallmentSimulationResult, error) {
	var resp types.InstallmentSimulationResult
	path := fmt.Sprintf("/accounts/%s/installmentAdvanceSimulation/%s", pathParam(accountID), pathParam(transactionID))
	if err := c.request(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RequestAdvancePayment requests advance payment for installments on a transaction.
func (c *Client) RequestAdvancePayment(ctx context.Context, accountID, transactionID string, req *types.AdvancePaymentRequest) (*types.ResultData, error) {
	var resp types.ResultData
	path := fmt.Sprintf("/accounts/%s/installmentAdvanceRequest/%s", pathParam(accountID), pathParam(transactionID))
	if err := c.request(ctx, http.MethodPost, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
