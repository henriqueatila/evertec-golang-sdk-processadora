package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// CreateCampaign creates a new campaign.
func (c *Client) CreateCampaign(ctx context.Context, req *types.CampaignObject) (*types.CampaignResponseObject, error) {
	var resp types.CampaignResponseObject
	if err := c.request(ctx, http.MethodPost, "/campaigns", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCampaign retrieves a campaign by ID.
func (c *Client) GetCampaign(ctx context.Context, campaignID string) (*types.CampaignResponseObject, error) {
	var resp types.CampaignResponseObject
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/campaigns/%s", pathParam(campaignID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListCampaigns retrieves all campaigns.
func (c *Client) ListCampaigns(ctx context.Context) (*types.CampaignListResponseObject, error) {
	var resp types.CampaignListResponseObject
	if err := c.request(ctx, http.MethodGet, "/campaigns", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCampaign updates a campaign.
func (c *Client) UpdateCampaign(ctx context.Context, campaignID string, req *types.CampaignObject) (*types.CampaignResponseObject, error) {
	var resp types.CampaignResponseObject
	if err := c.request(ctx, http.MethodPut, fmt.Sprintf("/campaigns/%s", pathParam(campaignID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCampaignStatus updates the status of a campaign.
func (c *Client) UpdateCampaignStatus(ctx context.Context, campaignID string, status types.CampaignStatus) error {
	req := map[string]string{"status": string(status)}
	return c.request(ctx, http.MethodPut, fmt.Sprintf("/campaigns/%s/status", pathParam(campaignID)), req, nil)
}

// GetCampaignAccounts retrieves accounts associated with a campaign.
func (c *Client) GetCampaignAccounts(ctx context.Context, campaignID string) (*types.CampaignAccountsListResponseObject, error) {
	var resp types.CampaignAccountsListResponseObject
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/campaigns/%s/accounts", pathParam(campaignID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateAgent creates a new agent.
func (c *Client) CreateAgent(ctx context.Context, req *types.AgentObject) (*types.AgentResponseObject, error) {
	var resp types.AgentResponseObject
	if err := c.request(ctx, http.MethodPost, "/agents", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAgent retrieves an agent by ID.
func (c *Client) GetAgent(ctx context.Context, agentID string) (*types.AgentResponseObject, error) {
	var resp types.AgentResponseObject
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/agents/%s", pathParam(agentID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAgents retrieves all agents.
func (c *Client) ListAgents(ctx context.Context) ([]types.AgentResponseObject, error) {
	var resp []types.AgentResponseObject
	if err := c.request(ctx, http.MethodGet, "/agents", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateAgent updates an agent.
func (c *Client) UpdateAgent(ctx context.Context, agentID string, req *types.AgentObject) (*types.AgentResponseObject, error) {
	var resp types.AgentResponseObject
	if err := c.request(ctx, http.MethodPut, fmt.Sprintf("/agents/%s", pathParam(agentID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateAgentStatus updates the status of an agent.
func (c *Client) UpdateAgentStatus(ctx context.Context, agentID string, status string) error {
	req := map[string]string{"status": status}
	return c.request(ctx, http.MethodPut, fmt.Sprintf("/agents/%s/status", pathParam(agentID)), req, nil)
}
