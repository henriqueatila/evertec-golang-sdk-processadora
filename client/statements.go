package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// GetClosedStatements retrieves all closed statements for a specific date.
// Docs: https://paysmart-api.gitlab.io/processadora/PT-br/docs/statements-api
func (c *Client) GetClosedStatements(ctx context.Context, req *types.ClosedStatementsRequest) (*types.ClosedStatementsResponse, error) {
	path := "/accounts/statements/closed"
	if req != nil {
		q := url.Values{}
		if req.ClosingDateQuery != "" {
			q.Set("closing_date_query", req.ClosingDateQuery)
		}
		if req.StartingAfter != "" {
			q.Set("starting_after", req.StartingAfter)
		}
		if len(q) > 0 {
			path = path + "?" + q.Encode()
		}
	}

	var resp types.ClosedStatementsResponse
	if err := c.request(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
