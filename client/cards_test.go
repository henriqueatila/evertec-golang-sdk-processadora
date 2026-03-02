package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// ============================================================================
// CARD OPERATIONS - PROFESSIONAL TABLE-DRIVEN TESTS
// ============================================================================

func TestClient_RequestNewCard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		accountID    string
		request      *types.NewCardRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
		validate     func(t *testing.T, resp *types.NewCardResponse)
	}{
		{
			name:      "success: request new card",
			accountID: "acc123",
			request: &types.NewCardRequest{
				IssuerRequestID: "req123",
				PsProductCode:   "VIRTUAL",
			},
			mockResponse: map[string]any{
				"cardId":      "card123",
				"cardStatus":  "ACTIVE",
				"last4Digits": "1234",
			},
			mockStatus: http.StatusCreated,
			wantErr:    false,
			validate: func(t *testing.T, resp *types.NewCardResponse) {
				AssertNotNil(t, resp, "response")
			},
		},
		{
			name:      "success: request physical card",
			accountID: "acc456",
			request: &types.NewCardRequest{
				IssuerRequestID:     "req456",
				PsProductCode:       "PHYSICAL",
				InhibitPhysicalCard: false,
			},
			mockResponse: map[string]any{"cardId": "card456"},
			mockStatus:   http.StatusCreated,
			wantErr:      false,
		},
		{
			name:         "error: account not found",
			accountID:    "acc-nonexistent",
			request:      &types.NewCardRequest{},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/accounts/" + tt.accountID + "/newCardRequest",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.RequestNewCard(context.Background(), tt.accountID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "RequestNewCard")
			} else {
				AssertNoError(t, err, "RequestNewCard")
				if tt.validate != nil {
					tt.validate(t, resp)
				}
			}
		})
	}
}

func TestClient_VerifyPin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cardID       string
		request      *types.VerifyPinRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, body map[string]any)
	}{
		{
			name:   "success: verify pin",
			cardID: "card123",
			request: &types.VerifyPinRequest{
				Pin: types.InputPin{
					IDTransportKey: "key-001",
					PinBlock:       "encrypted-pin-block",
				},
			},
			mockResponse: map[string]any{
				"resultData": map[string]any{"resultCode": 0},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]any) {
				pin, ok := body["pin"].(map[string]any)
				if ok {
					AssertEqual(t, "encrypted-pin-block", pin["pinBlock"], "pinBlock")
					AssertEqual(t, "key-001", pin["idTransportKey"], "idTransportKey")
				}
			},
		},
		{
			name:   "error: card not found",
			cardID: "card-nonexistent",
			request: &types.VerifyPinRequest{
				Pin: types.InputPin{
					IDTransportKey: "key",
					PinBlock:       "pin",
				},
			},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/cards/" + tt.cardID + "/validatePin",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]any
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.VerifyPin(context.Background(), tt.cardID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "VerifyPin")
			} else {
				AssertNoError(t, err, "VerifyPin")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_UpdateCard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cardID       string
		request      *types.UpdateCardRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, body map[string]any)
	}{
		{
			name:   "success: update card",
			cardID: "card123",
			request: &types.UpdateCardRequest{
				IssuerRequestID:  "req123",
				AllowContactless: true,
			},
			mockResponse: map[string]any{
				"cardId":           "card123",
				"allowContactless": true,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]any) {
				AssertEqual(t, true, body["allowContactless"], "allowContactless")
			},
		},
		{
			name:   "error: card not found",
			cardID: "card-nonexistent",
			request: &types.UpdateCardRequest{
				IssuerRequestID: "req456",
			},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPut,
				Path:     "/cards/" + tt.cardID,
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]any
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.UpdateCard(context.Background(), tt.cardID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "UpdateCard")
			} else {
				AssertNoError(t, err, "UpdateCard")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_CreateAnonymousCard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.AnonymousCardRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name: "success: create anonymous virtual card",
			request: &types.AnonymousCardRequest{
				IssuerRequestID:     "req123",
				PsProductCode:       "VIRTUAL",
				CardDeliveryAddress: map[string]any{"country": "BR"},
			},
			mockResponse: map[string]any{
				"cardId":    "anon123",
				"accountId": "acc123",
			},
			mockStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "error: bad request",
			request: &types.AnonymousCardRequest{
				PsProductCode: "INVALID",
			},
			mockResponse: nil,
			mockStatus:   http.StatusBadRequest,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/cards/newAnonymousCardRequest",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.CreateAnonymousCard(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "CreateAnonymousCard")
			} else {
				AssertNoError(t, err, "CreateAnonymousCard")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_BindAnonymousCard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.BindAnonymousCardRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, body map[string]any)
	}{
		{
			name: "success: bind anonymous card",
			request: &types.BindAnonymousCardRequest{
				PAN:        "4111111111111111",
				CVV:        "123",
				DateExp:    "12/28",
				Cardholder: map[string]any{"name": "Test User"},
			},
			mockResponse: map[string]any{
				"cardId":    "anon123",
				"accountId": "acc456",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]any) {
				AssertEqual(t, "4111111111111111", body["PAN"], "PAN")
			},
		},
		{
			name: "error: card not found",
			request: &types.BindAnonymousCardRequest{
				PAN:        "0000000000000000",
				CVV:        "000",
				DateExp:    "01/25",
				Cardholder: map[string]any{},
			},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/cards/bindAnonymousCardRequest",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]any
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.BindAnonymousCard(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "BindAnonymousCard")
			} else {
				AssertNoError(t, err, "BindAnonymousCard")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_FindCardByPAN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		request      *types.FindCardByPANRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, body map[string]any)
	}{
		{
			name: "success: find card by PAN",
			request: &types.FindCardByPANRequest{
				PAN: "4111111111111111",
			},
			mockResponse: map[string]any{
				"cardId": "card123",
				"status": "ACTIVE",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]any) {
				AssertEqual(t, "4111111111111111", body["PAN"], "PAN")
			},
		},
		{
			name: "error: card not found",
			request: &types.FindCardByPANRequest{
				PAN: "0000000000000000",
			},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/cards/findByPAN",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]any
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.FindCardByPAN(context.Background(), tt.request)

			if tt.wantErr {
				AssertError(t, err, "FindCardByPAN")
			} else {
				AssertNoError(t, err, "FindCardByPAN")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_ReissueCard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cardID       string
		request      *types.ReissueCardRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:   "success: reissue card",
			cardID: "card123",
			request: &types.ReissueCardRequest{
				IssuerCardReissueID: "reissue123",
				Reason:              "DAMAGED",
			},
			mockResponse: map[string]any{
				"cardId":    "card456",
				"oldCardId": "card123",
				"status":    "ACTIVE",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "error: card not found",
			cardID:       "card-nonexistent",
			request:      &types.ReissueCardRequest{Reason: "LOST"},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/cards/" + tt.cardID + "/reissueCardRequest",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.ReissueCard(context.Background(), tt.cardID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "ReissueCard")
			} else {
				AssertNoError(t, err, "ReissueCard")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_SetCardFunctions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cardID      string
		request     *types.CardFunctions
		mockStatus  int
		wantErr     bool
		validateReq func(t *testing.T, body map[string]any)
	}{
		{
			name:   "success: set card functions",
			cardID: "card123",
			request: &types.CardFunctions{
				ContactlessEnabled:   true,
				EcommerceEnabled:     true,
				ATMEnabled:           true,
				InternationalEnabled: false,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]any) {
				AssertEqual(t, true, body["contactlessEnabled"], "contactlessEnabled")
				AssertEqual(t, true, body["ecommerceEnabled"], "ecommerceEnabled")
			},
		},
		{
			name:   "error: card not found",
			cardID: "card-nonexistent",
			request: &types.CardFunctions{
				ContactlessEnabled: true,
			},
			mockStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/cards/" + tt.cardID + "/cardFunctions",
				Status:   tt.mockStatus,
				Response: nil,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]any
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			err := client.SetCardFunctions(context.Background(), tt.cardID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "SetCardFunctions")
			} else {
				AssertNoError(t, err, "SetCardFunctions")
			}
		})
	}
}

func TestClient_GetCardFunctions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cardID       string
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:   "success: get card functions",
			cardID: "card123",
			mockResponse: map[string]any{
				"allowedFunctions": []string{"PURCHASE", "WITHDRAW", "ONLINE"},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "error: card not found",
			cardID:       "card-nonexistent",
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/cards/" + tt.cardID + "/cardFunctions",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetCardFunctions(context.Background(), tt.cardID)

			if tt.wantErr {
				AssertError(t, err, "GetCardFunctions")
			} else {
				AssertNoError(t, err, "GetCardFunctions")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_ResetCardFunctions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		cardID     string
		mockStatus int
		wantErr    bool
	}{
		{
			name:       "success: reset card functions",
			cardID:     "card123",
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "error: card not found",
			cardID:     "card-nonexistent",
			mockStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPatch,
				Path:     "/cards/" + tt.cardID + "/resetCardFunctions",
				Status:   tt.mockStatus,
				Response: nil,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			err := client.ResetCardFunctions(context.Background(), tt.cardID)

			if tt.wantErr {
				AssertError(t, err, "ResetCardFunctions")
			} else {
				AssertNoError(t, err, "ResetCardFunctions")
			}
		})
	}
}

func TestClient_GetCardDataprepStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cardID       string
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:   "success: get dataprep status ready",
			cardID: "card123",
			mockResponse: map[string]any{
				"status":    "READY",
				"timestamp": "2024-01-15T10:30:00Z",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:         "success: get dataprep status pending",
			cardID:       "card456",
			mockResponse: map[string]any{"status": "PENDING"},
			mockStatus:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "error: card not found",
			cardID:       "card-nonexistent",
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/cards/" + tt.cardID + "/dataprepstatus",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetCardDataprepStatus(context.Background(), tt.cardID)

			if tt.wantErr {
				AssertError(t, err, "GetCardDataprepStatus")
			} else {
				AssertNoError(t, err, "GetCardDataprepStatus")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_RequestPhysicalCard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		accountID    string
		request      *types.PhysicalCardRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:      "success: request physical card",
			accountID: "acc123",
			request: &types.PhysicalCardRequest{
				IssuerRequestID: "req123",
				PsProductCode:   "CARD_STANDARD",
				Cardholder: map[string]any{
					"fullName": "JOAO DA SILVA",
				},
				CardDeliveryAddress: map[string]any{
					"street":  "Rua das Flores, 123",
					"city":    "Sao Paulo",
					"state":   "SP",
					"zipCode": "01234567",
					"country": "BR",
				},
			},
			mockResponse: map[string]any{
				"cardId":      "card123",
				"cardStatus":  "EMBOSSING",
				"last4Digits": "1234",
			},
			mockStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name:      "error: account not found",
			accountID: "acc-nonexistent",
			request: &types.PhysicalCardRequest{
				IssuerRequestID: "req456",
			},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/accounts/" + tt.accountID + "/newCardRequest",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.RequestPhysicalCard(context.Background(), tt.accountID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "RequestPhysicalCard")
			} else {
				AssertNoError(t, err, "RequestPhysicalCard")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_RequestVirtualCard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		accountID    string
		request      *types.CreateAccountVirtualCardRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:      "success: request virtual card",
			accountID: "acc123",
			request: &types.CreateAccountVirtualCardRequest{
				IssuerRequestID: "req123",
			},
			mockResponse: map[string]any{
				"virtualCardId": "vc123",
				"cardId":        "card123",
				"status":        "ACTIVE",
				"last4Digits":   "9876",
			},
			mockStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name:      "error: account not found",
			accountID: "acc-nonexistent",
			request: &types.CreateAccountVirtualCardRequest{
				IssuerRequestID: "req456",
			},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/accounts/" + tt.accountID + "/virtualcards",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.RequestVirtualCard(context.Background(), tt.accountID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "RequestVirtualCard")
			} else {
				AssertNoError(t, err, "RequestVirtualCard")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_BlockCard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cardID      string
		request     *types.BlockCardRequest
		mockStatus  int
		wantErr     bool
		validateReq func(t *testing.T, body map[string]any)
	}{
		{
			name:   "success: block card",
			cardID: "card123",
			request: &types.BlockCardRequest{
				BlockCode: 1,
				Reason:    "Lost card",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]any) {
				AssertEqual(t, "Lost card", body["reason"], "reason")
			},
		},
		{
			name:   "error: card not found",
			cardID: "card-nonexistent",
			request: &types.BlockCardRequest{
				BlockCode: 1,
			},
			mockStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/cards/" + tt.cardID + "/blockCardRequest",
				Status:   tt.mockStatus,
				Response: nil,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]any
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			err := client.BlockCard(context.Background(), tt.cardID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "BlockCard")
			} else {
				AssertNoError(t, err, "BlockCard")
			}
		})
	}
}

func TestClient_UnblockCard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		cardID     string
		request    *types.UnblockCardRequest
		mockStatus int
		wantErr    bool
	}{
		{
			name:   "success: unblock card",
			cardID: "card123",
			request: &types.UnblockCardRequest{
				Reason: "Customer request",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "success: unblock without reason",
			cardID:     "card456",
			request:    &types.UnblockCardRequest{},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:   "error: card not found",
			cardID: "card-nonexistent",
			request: &types.UnblockCardRequest{
				Reason: "Test",
			},
			mockStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/cards/" + tt.cardID + "/unblockCardRequest",
				Status:   tt.mockStatus,
				Response: nil,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			err := client.UnblockCard(context.Background(), tt.cardID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "UnblockCard")
			} else {
				AssertNoError(t, err, "UnblockCard")
			}
		})
	}
}

func TestClient_ChangePin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cardID      string
		request     *types.ChangePinRequest
		mockStatus  int
		wantErr     bool
		validateReq func(t *testing.T, body map[string]any)
	}{
		{
			name:   "success: change pin",
			cardID: "card123",
			request: &types.ChangePinRequest{
				NewPin: types.Pin{
					IDTransportKey: "key-001",
					PinBlock:       "encrypted-pin-block",
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]any) {
				newPin, ok := body["newPin"].(map[string]any)
				if ok {
					AssertEqual(t, "encrypted-pin-block", newPin["pinBlock"], "pinBlock")
				}
			},
		},
		{
			name:   "error: card not found",
			cardID: "card-nonexistent",
			request: &types.ChangePinRequest{
				NewPin: types.Pin{
					IDTransportKey: "key",
					PinBlock:       "pin",
				},
			},
			mockStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/cards/" + tt.cardID + "/changePin",
				Status:   tt.mockStatus,
				Response: nil,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]any
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			err := client.ChangePin(context.Background(), tt.cardID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "ChangePin")
			} else {
				AssertNoError(t, err, "ChangePin")
			}
		})
	}
}

func TestClient_ListCards(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		params       *types.ListCardsParams
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name: "success: list cards with filters",
			params: &types.ListCardsParams{
				Limit:                  10,
				IssuerCardID:           "issuer-card-123",
				IdentityDocumentNumber: "12345678900",
			},
			mockResponse: map[string]any{
				"cards": []map[string]any{
					{"cardId": "card001", "cardStatus": "ACTIVE"},
					{"cardId": "card002", "cardStatus": "ACTIVE"},
				},
				"hasMore": false,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:   "success: list all cards",
			params: nil,
			mockResponse: map[string]any{
				"cards": []map[string]any{
					{"cardId": "card001"},
				},
				"hasMore": true,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: list with pagination",
			params: &types.ListCardsParams{
				Limit:          20,
				StartingAfter:  "card100",
				PANLast4Digits: "1234",
			},
			mockResponse: map[string]any{
				"cards":   []map[string]any{},
				"hasMore": false,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: list with all filters",
			params: &types.ListCardsParams{
				Limit:                  10,
				StartingAfter:          "card100",
				EndingBefore:           "card200",
				IssuerCardID:           "issuer-001",
				IdentityDocumentNumber: "12345678900",
				PANLast4Digits:         "5678",
				IssuedOnOrAfterDate:    "01/01/2024",
				PSProductCode:          1001,
				LinkID:                 "lnk-123",
				AlternativeBindingKey:  "ABK001",
			},
			mockResponse: map[string]any{
				"cards":   []map[string]any{},
				"hasMore": false,
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "/cards",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.ListCards(context.Background(), tt.params)

			if tt.wantErr {
				AssertError(t, err, "ListCards")
			} else {
				AssertNoError(t, err, "ListCards")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_BlockAndReissueCard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cardID       string
		request      *types.BlockAndReissueCardRequest
		mockResponse any
		mockStatus   int
		wantErr      bool
	}{
		{
			name:   "success: block and reissue card",
			cardID: "card123",
			request: &types.BlockAndReissueCardRequest{
				BlockCode: 3,
				Reason:    "Stolen card",
			},
			mockResponse: map[string]any{
				"newCardId":   "card456",
				"oldCardId":   "card123",
				"status":      "EMBOSSING",
				"last4Digits": "5678",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:   "error: card not found",
			cardID: "card-nonexistent",
			request: &types.BlockAndReissueCardRequest{
				BlockCode: 1,
				Reason:    "Lost",
			},
			mockResponse: nil,
			mockStatus:   http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/cards/" + tt.cardID + "/blockAndReissueCardRequest",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.BlockAndReissueCard(context.Background(), tt.cardID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "BlockAndReissueCard")
			} else {
				AssertNoError(t, err, "BlockAndReissueCard")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

// ============================================================================
// CARD EDGE CASES
// ============================================================================

func TestClient_Cards_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty card ID", func(t *testing.T) {
		t.Parallel()

		server := NewMockServer(t, &MockServerConfig{
			Method:   http.MethodGet,
			Path:     "",
			Status:   http.StatusOK,
			Response: map[string]any{},
		})
		defer server.Close()

		client := NewTestClient(t, server)

		// Should not panic
		_, err := client.GetCardFunctions(context.Background(), "")
		_ = err
	})

	t.Run("special characters in card ID", func(t *testing.T) {
		t.Parallel()

		specialIDs := []string{
			"card/with/slash",
			"card?query=param",
			"card#hash",
		}

		for _, id := range specialIDs {
			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodGet,
				Path:     "",
				Status:   http.StatusOK,
				Response: map[string]any{},
			})

			client := NewTestClient(t, server)
			_, err := client.GetCardDataprepStatus(context.Background(), id)
			_ = err

			server.Close()
		}
	})
}
