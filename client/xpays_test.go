package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// ============================================================================
// XPAYS (DIGITAL WALLETS) - PROFESSIONAL TABLE-DRIVEN TESTS
// Updated to match official PaySmart API spec
// Reference: https://paysmart-api.gitlab.io/processadora/PT-br/docs/xpays-api
// ============================================================================

func TestClient_GetDeviceTokensFromCard_Pro(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cardID       string
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
		validate     func(t *testing.T, resp []types.DeviceToken)
	}{
		{
			name:   "success: get device tokens",
			cardID: "card-001",
			// API returns array directly per official spec
			mockResponse: []map[string]interface{}{
				{
					"deviceTokenId":    "token-001",
					"activationStatus": "ACTIVE",
					"wallet":           "apple",
					"dpanLast4Digits":  "1234",
				},
				{
					"deviceTokenId":    "token-002",
					"suspensionStatus": "SUSPENDED",
					"wallet":           "google",
					"cardReferenceId":  "ref-002",
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validate: func(t *testing.T, resp []types.DeviceToken) {
				if len(resp) != 2 {
					t.Errorf("expected 2 tokens, got %d", len(resp))
				}
			},
		},
		{
			name:         "success: no tokens",
			cardID:       "card-no-tokens",
			mockResponse: []map[string]interface{}{},
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
				Path:     "/cards/" + tt.cardID + "/deviceTokens",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetDeviceTokensFromCard(context.Background(), tt.cardID)

			if tt.wantErr {
				AssertError(t, err, "GetDeviceTokensFromCard")
			} else {
				AssertNoError(t, err, "GetDeviceTokensFromCard")
				if tt.validate != nil {
					tt.validate(t, resp)
				}
			}
		})
	}
}

func TestClient_MigrateDeviceTokens_Pro(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cardID       string
		request      *types.MigrateTokensRequest
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "success: migrate tokens to new card",
			cardID: "card-001",
			request: &types.MigrateTokensRequest{
				CardID: "card-002",
			},
			mockResponse: map[string]interface{}{
				"resultData": map[string]interface{}{
					"resultCode": 0,
				},
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]interface{}) {
				AssertEqual(t, "card-002", body["cardId"], "cardId")
			},
		},
		{
			name:         "error: target card not found",
			cardID:       "card-001",
			request:      &types.MigrateTokensRequest{CardID: "nonexistent"},
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
				Path:     "/cards/" + tt.cardID + "/deviceTokens/migrate",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]interface{}
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.MigrateDeviceTokens(context.Background(), tt.cardID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "MigrateDeviceTokens")
			} else {
				AssertNoError(t, err, "MigrateDeviceTokens")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_SuspendOrResumeDeviceTokensByCard_Pro(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cardID       string
		request      *types.SuspendResumeRequest
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "success: suspend all tokens",
			cardID: "card-001",
			request: &types.SuspendResumeRequest{
				Suspend:           true,
				ReasonDescription: "Card reported lost",
			},
			// API returns array directly per official spec
			mockResponse: []map[string]interface{}{},
			mockStatus:   http.StatusOK,
			wantErr:      false,
			validateReq: func(t *testing.T, body map[string]interface{}) {
				AssertEqual(t, true, body["suspend"], "suspend")
			},
		},
		{
			name:   "success: resume all tokens",
			cardID: "card-002",
			request: &types.SuspendResumeRequest{
				Suspend: false,
			},
			mockResponse: []map[string]interface{}{},
			mockStatus:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:   "success: suspend with reason description",
			cardID: "card-003",
			request: &types.SuspendResumeRequest{
				Suspend:           true,
				ReasonDescription: "Fraud investigation",
			},
			mockResponse: []map[string]interface{}{},
			mockStatus:   http.StatusOK,
			wantErr:      false,
			validateReq: func(t *testing.T, body map[string]interface{}) {
				AssertEqual(t, "Fraud investigation", body["reasonDescription"], "reasonDescription")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/cards/" + tt.cardID + "/deviceTokens/suspendOrResume",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]interface{}
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.SuspendOrResumeDeviceTokensByCard(context.Background(), tt.cardID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "SuspendOrResumeDeviceTokensByCard")
			} else {
				AssertNoError(t, err, "SuspendOrResumeDeviceTokensByCard")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_TerminateDeviceTokensByCard_Pro(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cardID       string
		request      *types.TerminateTokensRequest
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "success: terminate all tokens - card stolen",
			cardID: "card-001",
			request: &types.TerminateTokensRequest{
				Reason:            types.TerminateReasonRevoked,
				ReasonDescription: "Card reported stolen",
			},
			// API returns array directly per official spec
			mockResponse: []map[string]interface{}{},
			mockStatus:   http.StatusOK,
			wantErr:      false,
			validateReq: func(t *testing.T, body map[string]interface{}) {
				AssertEqual(t, "REVOKED", body["reason"], "reason")
			},
		},
		{
			name:   "success: terminate - card expired",
			cardID: "card-002",
			request: &types.TerminateTokensRequest{
				Reason: types.TerminateReasonExpired,
			},
			mockResponse: []map[string]interface{}{},
			mockStatus:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:   "success: terminate - card deleted",
			cardID: "card-003",
			request: &types.TerminateTokensRequest{
				Reason: types.TerminateReasonDeleted,
			},
			mockResponse: []map[string]interface{}{},
			mockStatus:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "error: card not found",
			cardID:       "card-nonexistent",
			request:      &types.TerminateTokensRequest{Reason: types.TerminateReasonDeleted},
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
				Path:     "/cards/" + tt.cardID + "/deviceTokens/terminate",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]interface{}
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.TerminateDeviceTokensByCard(context.Background(), tt.cardID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "TerminateDeviceTokensByCard")
			} else {
				AssertNoError(t, err, "TerminateDeviceTokensByCard")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_GetEncryptedCard_Pro(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cardID       string
		request      *types.EncryptedCardRequest
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
		validateReq  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:   "success: get encrypted card data",
			cardID: "card-001",
			request: &types.EncryptedCardRequest{
				// Base64-encoded strings per official API spec
				Nonce:          "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXowMTIzNDU2Nzg5",
				NonceSignature: "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXowMTIzNDU2Nzg5",
				Certificates:   []string{"Y2VydDE=", "Y2VydDI="},
			},
			mockResponse: map[string]interface{}{
				"ephemeralPublicKey": "ZXBoZW1lcmFsS2V5",
				"encryptedPassData":  "ZW5jcnlwdGVkRGF0YQ==",
				"activationData":     "YWN0aXZhdGlvbkRhdGE=",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]interface{}) {
				AssertNotNil(t, body["nonce"], "nonce")
			},
		},
		{
			name:   "success: minimal request",
			cardID: "card-002",
			request: &types.EncryptedCardRequest{
				Nonce: "c2ltcGxlLW5vbmNl",
			},
			mockResponse: map[string]interface{}{"encryptedPassData": "data"},
			mockStatus:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "error: card not found",
			cardID:       "card-nonexistent",
			request:      &types.EncryptedCardRequest{Nonce: "bm9uY2U="},
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
				Path:     "/cards/" + tt.cardID + "/encryptedCard",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]interface{}
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetEncryptedCard(context.Background(), tt.cardID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "GetEncryptedCard")
			} else {
				AssertNoError(t, err, "GetEncryptedCard")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_GetOpaqueCard_Pro(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cardID       string
		mockResponse interface{}
		mockStatus   int
		wantErr      bool
		validate     func(t *testing.T, resp *types.OpaqueCardResponse)
	}{
		{
			name:   "success: get opaque card data",
			cardID: "card-001",
			mockResponse: map[string]interface{}{
				"sender":         "issuer-001",
				"cardDescriptor": "Y2FyZERlc2NyaXB0b3I=", // base64
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validate: func(t *testing.T, resp *types.OpaqueCardResponse) {
				AssertNotNil(t, resp, "response")
				AssertEqual(t, "issuer-001", resp.Sender, "sender")
			},
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
				Path:     "/cards/" + tt.cardID + "/opaqueCard",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.GetOpaqueCard(context.Background(), tt.cardID)

			if tt.wantErr {
				AssertError(t, err, "GetOpaqueCard")
			} else {
				AssertNoError(t, err, "GetOpaqueCard")
				if tt.validate != nil {
					tt.validate(t, resp)
				}
			}
		})
	}
}

// ============================================================================
// DEVICE TOKEN OPERATIONS (Individual tokens)
// ============================================================================

func TestClient_ActivateDeviceToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		deviceTokenID string
		request       *types.ActivateTokenRequest
		mockResponse  interface{}
		mockStatus    int
		wantErr       bool
		validateReq   func(t *testing.T, body map[string]interface{})
	}{
		{
			name:          "success: activate token with App method",
			deviceTokenID: "token-001",
			request: &types.ActivateTokenRequest{
				IDVMethod: types.IDVMethodApp,
			},
			// API returns DeviceToken directly per official spec
			mockResponse: map[string]interface{}{
				"deviceTokenId":    "token-001",
				"activationStatus": "ACTIVE",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]interface{}) {
				AssertEqual(t, "App", body["idvMethod"], "idvMethod")
			},
		},
		{
			name:          "success: activate with OTP SMS",
			deviceTokenID: "token-002",
			request: &types.ActivateTokenRequest{
				IDVMethod:         types.IDVMethodOTPSMS,
				ReasonDescription: "Customer requested via support",
			},
			mockResponse: map[string]interface{}{"deviceTokenId": "token-002", "activationStatus": "ACTIVE"},
			mockStatus:   http.StatusOK,
			wantErr:      false,
			validateReq: func(t *testing.T, body map[string]interface{}) {
				AssertEqual(t, "OTPSMS", body["idvMethod"], "idvMethod")
			},
		},
		{
			name:          "success: activate with OTP Email",
			deviceTokenID: "token-003",
			request: &types.ActivateTokenRequest{
				IDVMethod: types.IDVMethodOTPEmail,
			},
			mockResponse: map[string]interface{}{"deviceTokenId": "token-003", "activationStatus": "ACTIVE"},
			mockStatus:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:          "error: token not found",
			deviceTokenID: "token-nonexistent",
			request:       &types.ActivateTokenRequest{IDVMethod: types.IDVMethodApp},
			mockResponse:  nil,
			mockStatus:    http.StatusNotFound,
			wantErr:       true,
		},
		{
			name:          "error: already active",
			deviceTokenID: "token-active",
			request:       &types.ActivateTokenRequest{IDVMethod: types.IDVMethodApp},
			mockResponse:  nil,
			mockStatus:    http.StatusConflict,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/deviceTokens/" + tt.deviceTokenID + "/activate",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]interface{}
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.ActivateDeviceToken(context.Background(), tt.deviceTokenID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "ActivateDeviceToken")
			} else {
				AssertNoError(t, err, "ActivateDeviceToken")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_SuspendOrResumeDeviceToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		deviceTokenID string
		request       *types.SuspendResumeRequest
		mockResponse  interface{}
		mockStatus    int
		wantErr       bool
		validateReq   func(t *testing.T, body map[string]interface{})
	}{
		{
			name:          "success: suspend token",
			deviceTokenID: "token-001",
			request: &types.SuspendResumeRequest{
				Suspend:           true,
				ReasonDescription: "Lost device",
			},
			// API returns DeviceToken directly per official spec
			mockResponse: map[string]interface{}{
				"deviceTokenId":    "token-001",
				"suspensionStatus": "SUSPENDED",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]interface{}) {
				AssertEqual(t, true, body["suspend"], "suspend")
				AssertEqual(t, "Lost device", body["reasonDescription"], "reasonDescription")
			},
		},
		{
			name:          "success: resume token",
			deviceTokenID: "token-002",
			request: &types.SuspendResumeRequest{
				Suspend: false,
			},
			mockResponse: map[string]interface{}{"deviceTokenId": "token-002", "suspensionStatus": "NOT_SUSPENDED"},
			mockStatus:   http.StatusOK,
			wantErr:      false,
			validateReq: func(t *testing.T, body map[string]interface{}) {
				AssertEqual(t, false, body["suspend"], "suspend")
			},
		},
		{
			name:          "error: token not found",
			deviceTokenID: "token-nonexistent",
			request:       &types.SuspendResumeRequest{Suspend: true},
			mockResponse:  nil,
			mockStatus:    http.StatusNotFound,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/deviceTokens/" + tt.deviceTokenID + "/suspendOrResume",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]interface{}
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.SuspendOrResumeDeviceToken(context.Background(), tt.deviceTokenID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "SuspendOrResumeDeviceToken")
			} else {
				AssertNoError(t, err, "SuspendOrResumeDeviceToken")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

func TestClient_TerminateDeviceToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		deviceTokenID string
		request       *types.TerminateTokenRequest
		mockResponse  interface{}
		mockStatus    int
		wantErr       bool
		validateReq   func(t *testing.T, body map[string]interface{})
	}{
		{
			name:          "success: terminate token - revoked",
			deviceTokenID: "token-001",
			request: &types.TerminateTokenRequest{
				Reason:            types.TerminateReasonRevoked,
				ReasonDescription: "Card reported stolen",
			},
			// API returns DeviceToken directly per official spec
			mockResponse: map[string]interface{}{
				"deviceTokenId":    "token-001",
				"activationStatus": "TERMINATED",
			},
			mockStatus: http.StatusOK,
			wantErr:    false,
			validateReq: func(t *testing.T, body map[string]interface{}) {
				AssertEqual(t, "REVOKED", body["reason"], "reason")
			},
		},
		{
			name:          "success: terminate - expired",
			deviceTokenID: "token-002",
			request: &types.TerminateTokenRequest{
				Reason: types.TerminateReasonExpired,
			},
			mockResponse: map[string]interface{}{"deviceTokenId": "token-002", "activationStatus": "TERMINATED"},
			mockStatus:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:          "success: terminate - deleted",
			deviceTokenID: "token-003",
			request: &types.TerminateTokenRequest{
				Reason:            types.TerminateReasonDeleted,
				ReasonDescription: "Card replacement",
			},
			mockResponse: map[string]interface{}{"deviceTokenId": "token-003", "activationStatus": "TERMINATED"},
			mockStatus:   http.StatusOK,
			wantErr:      false,
		},
		{
			name:          "error: token not found",
			deviceTokenID: "token-nonexistent",
			request:       &types.TerminateTokenRequest{Reason: types.TerminateReasonDeleted},
			mockResponse:  nil,
			mockStatus:    http.StatusNotFound,
			wantErr:       true,
		},
		{
			name:          "error: already terminated",
			deviceTokenID: "token-terminated",
			request:       &types.TerminateTokenRequest{Reason: types.TerminateReasonDeleted},
			mockResponse:  nil,
			mockStatus:    http.StatusConflict,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "/deviceTokens/" + tt.deviceTokenID + "/terminate",
				Status:   tt.mockStatus,
				Response: tt.mockResponse,
				ValidateReq: func(t *testing.T, r *http.Request, body []byte) {
					if tt.validateReq != nil {
						var reqBody map[string]interface{}
						_ = json.Unmarshal(body, &reqBody)
						tt.validateReq(t, reqBody)
					}
				},
			})
			defer server.Close()

			client := NewTestClient(t, server)
			resp, err := client.TerminateDeviceToken(context.Background(), tt.deviceTokenID, tt.request)

			if tt.wantErr {
				AssertError(t, err, "TerminateDeviceToken")
			} else {
				AssertNoError(t, err, "TerminateDeviceToken")
				AssertNotNil(t, resp, "response")
			}
		})
	}
}

// ============================================================================
// EDGE CASE TESTS FOR XPAYS
// ============================================================================

func TestClient_XPays_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty card ID", func(t *testing.T) {
		t.Parallel()

		server := NewMockServer(t, &MockServerConfig{
			Method:   http.MethodGet,
			Path:     "",
			Status:   http.StatusOK,
			Response: map[string]interface{}{},
		})
		defer server.Close()

		client := NewTestClient(t, server)

		// Should not panic
		_, err := client.GetDeviceTokensFromCard(context.Background(), "")
		_ = err
	})

	t.Run("empty device token ID", func(t *testing.T) {
		t.Parallel()

		server := NewMockServer(t, &MockServerConfig{
			Method:   http.MethodPost,
			Path:     "",
			Status:   http.StatusOK,
			Response: map[string]interface{}{},
		})
		defer server.Close()

		client := NewTestClient(t, server)

		// Should not panic
		_, err := client.ActivateDeviceToken(context.Background(), "", &types.ActivateTokenRequest{
			IDVMethod: types.IDVMethodApp,
		})
		_ = err
	})

	t.Run("special characters in IDs", func(t *testing.T) {
		t.Parallel()

		specialIDs := []string{
			"token/with/slash",
			"token?query=param",
			"token#hash",
		}

		for _, id := range specialIDs {
			server := NewMockServer(t, &MockServerConfig{
				Method:   http.MethodPost,
				Path:     "",
				Status:   http.StatusOK,
				Response: map[string]interface{}{},
			})

			client := NewTestClient(t, server)

			// Should not panic
			_, err := client.TerminateDeviceToken(context.Background(), id, &types.TerminateTokenRequest{
				Reason: types.TerminateReasonDeleted,
			})
			_ = err

			server.Close()
		}
	})
}

// ============================================================================
// CONCURRENT XPAYS TESTS
// ============================================================================

func TestClient_XPays_Concurrent(t *testing.T) {
	t.Parallel()

	server := NewMockServer(t, &MockServerConfig{
		Method:   http.MethodPost,
		Path:     "",
		Status:   http.StatusOK,
		Response: map[string]interface{}{"status": "SUCCESS"},
	})
	defer server.Close()

	client := NewTestClient(t, server)

	// Run concurrent operations using ConcurrentTestRunner
	runner := NewConcurrentTestRunner(t)
	for i := 0; i < 10; i++ {
		idx := i
		runner.Run(1, func() error {
			ctx := context.Background()
			tokenID := "token-" + string(rune('A'+idx))

			_, _ = client.ActivateDeviceToken(ctx, tokenID, &types.ActivateTokenRequest{
				IDVMethod: types.IDVMethodApp,
			})

			_, _ = client.SuspendOrResumeDeviceToken(ctx, tokenID, &types.SuspendResumeRequest{
				Suspend: true,
			})

			_, _ = client.TerminateDeviceToken(ctx, tokenID, &types.TerminateTokenRequest{
				Reason: types.TerminateReasonDeleted,
			})
			return nil
		})
	}
	runner.Wait()
}
