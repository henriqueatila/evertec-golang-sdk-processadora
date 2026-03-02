package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// RequestPhysicalCard requests a new physical card for an account.
// Docs: https://paysmart-api.gitlab.io/processadora/PT-br/docs/criacao-de-cartoes-fisicos
func (c *Client) RequestPhysicalCard(ctx context.Context, accountID string, req *types.PhysicalCardRequest) (*types.PhysicalCardResponse, error) {
	var resp types.PhysicalCardResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/newCardRequest", pathParam(accountID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RequestVirtualCard requests a new virtual card for an account.
// Virtual cards are immediately usable after creation without activation.
// Docs: https://paysmart-api.gitlab.io/processadora/PT-br/docs/criacao-de-cartoes-virtuais
func (c *Client) RequestVirtualCard(ctx context.Context, accountID string, req *types.CreateAccountVirtualCardRequest) (*types.CreateVirtualCardResponse, error) {
	var resp types.CreateVirtualCardResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/virtualcards", pathParam(accountID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RequestNewCard is deprecated, use RequestPhysicalCard or RequestVirtualCard instead.
// Deprecated: This method is kept for backward compatibility.
func (c *Client) RequestNewCard(ctx context.Context, accountID string, req *types.NewCardRequest) (*types.NewCardResponse, error) {
	var resp types.NewCardResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/accounts/%s/newCardRequest", pathParam(accountID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCard retrieves a card by ID.
func (c *Client) GetCard(ctx context.Context, cardID string) (*types.CardDetails, error) {
	var resp types.CardDetails
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/cards/%s", pathParam(cardID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BlockCard blocks a card.
func (c *Client) BlockCard(ctx context.Context, cardID string, req *types.BlockCardRequest) error {
	return c.request(ctx, http.MethodPost, fmt.Sprintf("/cards/%s/blockCardRequest", pathParam(cardID)), req, nil)
}

// UnblockCard unblocks a card.
func (c *Client) UnblockCard(ctx context.Context, cardID string, req *types.UnblockCardRequest) error {
	return c.request(ctx, http.MethodPost, fmt.Sprintf("/cards/%s/unblockCardRequest", pathParam(cardID)), req, nil)
}

// ChangePin changes the card PIN.
func (c *Client) ChangePin(ctx context.Context, cardID string, req *types.ChangePinRequest) error {
	return c.request(ctx, http.MethodPost, fmt.Sprintf("/cards/%s/changePin", pathParam(cardID)), req, nil)
}

// VerifyPin verifies the card PIN.
func (c *Client) VerifyPin(ctx context.Context, cardID string, req *types.VerifyPinRequest) (*types.VerifyPinResponse, error) {
	var resp types.VerifyPinResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/cards/%s/validatePin", pathParam(cardID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListCards retrieves cards based on filters.
func (c *Client) ListCards(ctx context.Context, params *types.ListCardsParams) (*types.ListCardsResponse, error) {
	path := "/cards"
	if params != nil {
		path = path + "?" + buildCardsQuery(params)
	}

	var resp types.ListCardsResponse
	if err := c.request(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func buildCardsQuery(params *types.ListCardsParams) string {
	q := url.Values{}
	addPaginationParams(q, params.Limit, params.StartingAfter, params.EndingBefore)
	if params.IssuerCardID != "" {
		q.Set("issuer_id", params.IssuerCardID)
	}
	if params.IdentityDocumentNumber != "" {
		q.Set("identity_document_number", params.IdentityDocumentNumber)
	}
	if params.PANLast4Digits != "" {
		q.Set("pan_last_4_digits", params.PANLast4Digits)
	}
	if params.IssuedOnOrAfterDate != "" {
		q.Set("issued_on_or_after_date", params.IssuedOnOrAfterDate)
	}
	if params.PSProductCode > 0 {
		q.Set("ps_product_code", strconv.Itoa(params.PSProductCode))
	}
	if params.LinkID != "" {
		q.Set("link_id", params.LinkID)
	}
	if params.AlternativeBindingKey != "" {
		q.Set("alternative_binding_key", params.AlternativeBindingKey)
	}
	return q.Encode()
}

// UpdateCard updates a card.
func (c *Client) UpdateCard(ctx context.Context, cardID string, req *types.UpdateCardRequest) (*types.CardDetails, error) {
	var resp types.CardDetails
	if err := c.request(ctx, http.MethodPut, fmt.Sprintf("/cards/%s", pathParam(cardID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateAnonymousCard creates an anonymous card (not linked to an account).
func (c *Client) CreateAnonymousCard(ctx context.Context, req *types.AnonymousCardRequest) (*types.AnonymousCardResponse, error) {
	var resp types.AnonymousCardResponse
	if err := c.request(ctx, http.MethodPost, "/cards/newAnonymousCardRequest", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BindAnonymousCard binds an anonymous card to an account.
func (c *Client) BindAnonymousCard(ctx context.Context, req *types.BindAnonymousCardRequest) (*types.CardDetails, error) {
	var resp types.CardDetails
	if err := c.request(ctx, http.MethodPost, "/cards/bindAnonymousCardRequest", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FindCardByPAN finds a card by its PAN (Primary Account Number).
func (c *Client) FindCardByPAN(ctx context.Context, req *types.FindCardByPANRequest) (*types.CardDetails, error) {
	var resp types.CardDetails
	if err := c.request(ctx, http.MethodPost, "/cards/findByPAN", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BlockAndReissueCard blocks a card and issues a replacement.
func (c *Client) BlockAndReissueCard(ctx context.Context, cardID string, req *types.BlockAndReissueCardRequest) (*types.NewCardResponse, error) {
	var resp types.NewCardResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/cards/%s/blockAndReissueCardRequest", pathParam(cardID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReissueCard reissues a card without blocking the current one.
func (c *Client) ReissueCard(ctx context.Context, cardID string, req *types.ReissueCardRequest) (*types.NewCardResponse, error) {
	var resp types.NewCardResponse
	if err := c.request(ctx, http.MethodPost, fmt.Sprintf("/cards/%s/reissueCardRequest", pathParam(cardID)), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetCardFunctions sets the enabled functions for a card.
func (c *Client) SetCardFunctions(ctx context.Context, cardID string, req *types.CardFunctions) error {
	return c.request(ctx, http.MethodPost, fmt.Sprintf("/cards/%s/cardFunctions", pathParam(cardID)), req, nil)
}

// GetCardFunctions retrieves the enabled functions for a card.
func (c *Client) GetCardFunctions(ctx context.Context, cardID string) (*types.CardFunctions, error) {
	var resp types.CardFunctions
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/cards/%s/cardFunctions", pathParam(cardID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ResetCardFunctions resets card functions to default settings.
func (c *Client) ResetCardFunctions(ctx context.Context, cardID string) error {
	return c.request(ctx, http.MethodPatch, fmt.Sprintf("/cards/%s/resetCardFunctions", pathParam(cardID)), nil, nil)
}

// GetCardDataprepStatus retrieves the dataprep status of a card.
func (c *Client) GetCardDataprepStatus(ctx context.Context, cardID string) (*types.CardDataprepStatus, error) {
	var resp types.CardDataprepStatus
	if err := c.request(ctx, http.MethodGet, fmt.Sprintf("/cards/%s/dataprepstatus", pathParam(cardID)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
