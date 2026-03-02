package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// ListDisputes lists disputes with optional filters.
func (c *Client) ListDisputes(ctx context.Context, req *types.ListDisputesRequest) (*types.ListDisputesResponse, error) {
	path := "/disputes"
	if req != nil {
		path = path + "?" + buildDisputeQuery(req)
	}

	var resp types.ListDisputesResponse
	if err := c.request(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateDispute creates a new dispute.
func (c *Client) CreateDispute(ctx context.Context, req *types.CreateDisputeRequest) (*types.CreateDisputeResponse, error) {
	var resp types.CreateDisputeResponse
	if err := c.request(ctx, http.MethodPost, "/disputes", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDispute retrieves a dispute by ID.
func (c *Client) GetDispute(ctx context.Context, disputeID string) (*types.Dispute, error) {
	var resp types.Dispute
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/disputes/%s", pathParam(disputeID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelDispute cancels a dispute (undo).
func (c *Client) CancelDispute(ctx context.Context, disputeID string, req *types.CancelDisputeRequest) error {
	return c.request(ctx, http.MethodPost, fmt.Sprintf("/disputes/%s/undo", pathParam(disputeID)), req, nil)
}

// AttachDisputeDocument attaches a document to an existing dispute.
func (c *Client) AttachDisputeDocument(ctx context.Context, disputeID string, req *types.DisputeDocumentRequest) (*types.DisputeDocumentCreatedSuccessfully, error) {
	var resp types.DisputeDocumentCreatedSuccessfully
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/disputes/%s/documents", pathParam(disputeID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RespondToDispute responds to a dispute ruling (accept or continue to next phase).
func (c *Client) RespondToDispute(ctx context.Context, disputeID string, req *types.DisputeResponseRequest) (*types.DisputeResponseCreatedSuccessfully, error) {
	var resp types.DisputeResponseCreatedSuccessfully
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/disputes/%s/response", pathParam(disputeID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func buildDisputeQuery(req *types.ListDisputesRequest) string {
	q := url.Values{}
	addPaginationParams(q, req.Limit, req.StartingAfter, req.EndingBefore)
	if req.DisputeCode != "" {
		q.Set("dispute_reason", req.DisputeCode)
	}
	if req.DisputeStatus != "" {
		q.Set("dispute_status", req.DisputeStatus)
	}
	if req.BeginningDate != "" {
		q.Set("beginning_date", req.BeginningDate)
	}
	if req.EndingDate != "" {
		q.Set("ending_date", req.EndingDate)
	}
	return q.Encode()
}
