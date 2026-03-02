package authorization

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// Non-standard HTTP status codes used by the Evertec authorization protocol.
const (
	StatusFraudSuspect        = 459
	StatusInvalidMCC          = 483
	StatusDeclinedByProcessor = 499
)

// Handler defines the interface for authorization handlers.
// Implement this interface to handle Evertec authorization callbacks.
// Reference: https://paysmart-api.gitlab.io/processadora/PT-br/docs/compra
type Handler interface {
	// HandlePurchase handles purchase authorization requests.
	// POST /purchases - Evertec calls this when a purchase needs authorization.
	HandlePurchase(ctx context.Context, req *PurchaseRequest) (*PurchaseResponse, error)

	// HandlePurchaseCancel handles purchase cancellation requests.
	// POST /purchases/cancel - Evertec calls this to cancel a purchase.
	HandlePurchaseCancel(ctx context.Context, req *PurchaseCancellationRequest) (*CancellationResponse, error)

	// HandleQuery handles balance query requests.
	// POST /queries - Evertec calls this to query account balance.
	HandleQuery(ctx context.Context, req *QueryRequest) (*QueryResponse, error)

	// HandleWithdrawal handles withdrawal authorization requests.
	// POST /withdrawals - Evertec calls this when a withdrawal needs authorization.
	HandleWithdrawal(ctx context.Context, req *WithdrawalRequest) (*WithdrawalResponse, error)

	// HandleWithdrawalQuery handles withdrawal query requests.
	// POST /withdrawalQueries - Evertec calls this to query withdrawal limits.
	HandleWithdrawalQuery(ctx context.Context, req *WithdrawalQueryRequest) (*WithdrawalQueryResponse, error)

	// HandleWithdrawalCancel handles withdrawal cancellation requests.
	// POST /withdrawals/cancel - Evertec calls this to cancel a withdrawal.
	HandleWithdrawalCancel(ctx context.Context, req *WithdrawalCancellationRequest) (*CancellationResponse, error)

	// HandleChargeback handles chargeback requests.
	// POST /chargebacks - Evertec calls this when a chargeback is initiated.
	HandleChargeback(ctx context.Context, req *ChargebackRequest) (*ChargebackResponse, error)

	// HandleChargebackCancel handles chargeback cancellation requests.
	// POST /chargebacks/cancel - Evertec calls this to cancel a chargeback.
	HandleChargebackCancel(ctx context.Context, req *ChargebackCancellationRequest) (*CancellationResponse, error)

	// HandleTransfer handles transfer authorization requests.
	// POST /transfers - Evertec calls this when a P2P transfer needs authorization.
	HandleTransfer(ctx context.Context, req *TransferRequest) (*TransferResponse, error)

	// HandleTransferCancel handles transfer cancellation requests.
	// POST /transfers/cancel - Evertec calls this to cancel a transfer.
	HandleTransferCancel(ctx context.Context, req *TransferCancellationRequest) (*CancellationResponse, error)

	// HandleGetOTPChannel handles OTP channel requests for 3DS.
	// POST /acs/getOTPChannel - Evertec calls this for 3DS authentication.
	HandleGetOTPChannel(ctx context.Context, req *OTPChannelRequest) (*OTPChannelResponse, error)

	// HandleVerifyTransaction handles 3DS verification requests.
	// POST /acs/verifyTransaction - Evertec calls this to verify 3DS OTP.
	HandleVerifyTransaction(ctx context.Context, req *VerifyTransactionRequest) (*VerifyTransactionResponse, error)

	// HandleXPaysOTP handles xPays OTP requests.
	// POST /xpays/otp - Evertec calls this for wallet provisioning OTP.
	HandleXPaysOTP(ctx context.Context, req *XPaysOTPRequest) (*XPaysOTPResponse, error)

	// HandleCustomProvisioningData handles xPays custom provisioning data requests.
	// POST /xpays/customProvisioningData - Evertec calls this for custom wallet data.
	HandleCustomProvisioningData(ctx context.Context, req *CustomProvisioningDataRequest) (*types.CustomProvisioningDataResponse, error)

	// HandleStatus handles health check requests.
	// GET /status - Evertec calls this to check if your server is healthy.
	HandleStatus(ctx context.Context) (*StatusResponse, error)
}

// Server wraps a Handler and provides HTTP routing.
type Server struct {
	handler       Handler
	hooks         []Hook
	logger        *slog.Logger
	maxBodySize   int64
	panicRecovery bool
}

// ServerOption is a functional option for configuring the server.
type ServerOption func(*Server)

// WithHooks adds observability hooks to the server.
func WithHooks(hooks ...Hook) ServerOption {
	return func(s *Server) {
		s.hooks = append(s.hooks, hooks...)
	}
}

// WithLogger adds a structured logger to the server.
func WithLogger(logger *slog.Logger) ServerOption {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithMaxBodySize sets the maximum request body size in bytes.
// If 0 (default), no limit is applied.
// Example: WithMaxBodySize(10 * 1024 * 1024) for 10MB limit.
func WithMaxBodySize(size int64) ServerOption {
	return func(s *Server) {
		s.maxBodySize = size
	}
}

// WithPanicRecovery enables panic recovery middleware.
// When enabled, panics in handlers are caught and logged (if a logger is configured),
// returning a 500 Internal Server Error instead of crashing the server.
func WithPanicRecovery() ServerOption {
	return func(s *Server) {
		s.panicRecovery = true
	}
}

// NewServer creates a new authorization server.
func NewServer(handler Handler, opts ...ServerOption) *Server {
	s := &Server{handler: handler}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.handler == nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("handler not configured"))
		return
	}

	// Panic recovery if enabled
	if s.panicRecovery {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				if s.logger != nil {
					s.logger.Error("panic recovered in authorization handler",
						"error", fmt.Sprintf("%v", err),
						"stack", string(stack),
						"path", r.URL.Path,
					)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
			}
		}()
	}

	w.Header().Set("Content-Type", "application/json")

	// Apply request body size limit if configured
	if r.Body != nil && s.maxBodySize > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, s.maxBodySize)
	}

	// Read body for hooks
	var bodyBytes []byte
	if r.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			if s.logger != nil {
				s.logger.Error("failed to read request body", "error", err, "path", r.URL.Path)
			}
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to read request body"})
			return
		}
	}

	// Determine operation type from path
	operationType := pathToOperationType(r.URL.Path)

	// Create request info for hooks
	reqInfo := &RequestInfo{
		Path:          r.URL.Path,
		OperationType: operationType,
		BodySize:      len(bodyBytes),
	}

	// Execute before hooks
	start := time.Now()
	ctx := s.executeBeforeHooks(r.Context(), reqInfo)

	// Create response recorder to capture status code
	recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

	// Route the request
	switch r.URL.Path {
	case "/status":
		if r.Method == http.MethodGet {
			s.handleStatus(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/purchases":
		if r.Method == http.MethodPost {
			s.handlePurchase(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/purchases/cancel":
		if r.Method == http.MethodPost {
			s.handlePurchaseCancel(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/queries":
		if r.Method == http.MethodPost {
			s.handleQuery(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/withdrawalQueries":
		if r.Method == http.MethodPost {
			s.handleWithdrawalQuery(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/withdrawals":
		if r.Method == http.MethodPost {
			s.handleWithdrawal(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/withdrawals/cancel":
		if r.Method == http.MethodPost {
			s.handleWithdrawalCancel(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/chargebacks":
		if r.Method == http.MethodPost {
			s.handleChargeback(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/chargebacks/cancel":
		if r.Method == http.MethodPost {
			s.handleChargebackCancel(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/transfers":
		if r.Method == http.MethodPost {
			s.handleTransfer(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/transfers/cancel":
		if r.Method == http.MethodPost {
			s.handleTransferCancel(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/acs/getOTPChannel":
		if r.Method == http.MethodPost {
			s.handleGetOTPChannel(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/acs/verifyTransaction":
		if r.Method == http.MethodPost {
			s.handleVerifyTransaction(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/xpays/otp":
		if r.Method == http.MethodPost {
			s.handleXPaysOTP(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	case "/xpays/customProvisioningData":
		if r.Method == http.MethodPost {
			s.handleCustomProvisioningData(recorder, r, ctx, reqInfo, start, bodyBytes)
		} else {
			s.handleMethodNotAllowed(recorder, r, ctx, reqInfo, start)
		}

	default:
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusNotFound,
			Duration:   time.Since(start),
		})
		http.NotFound(w, r)
	}
}

// responseRecorder wraps http.ResponseWriter to capture the status code.
// It also implements http.Hijacker and http.Flusher interfaces for
// compatibility with WebSocket upgrades and streaming responses.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// Hijack implements http.Hijacker interface for WebSocket support.
func (r *responseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("response writer does not implement http.Hijacker")
	}
	return h.Hijack()
}

// Flush implements http.Flusher interface for streaming responses.
func (r *responseRecorder) Flush() {
	f, ok := r.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}
	f.Flush()
}

// pathToOperationType maps URL paths to operation types.
func pathToOperationType(path string) string {
	switch path {
	case "/status":
		return "status"
	case "/purchases":
		return "purchase"
	case "/purchases/cancel":
		return "purchase_cancel"
	case "/queries":
		return "query"
	case "/withdrawalQueries":
		return "withdrawal_query"
	case "/withdrawals":
		return "withdrawal"
	case "/withdrawals/cancel":
		return "withdrawal_cancel"
	case "/chargebacks":
		return "chargeback"
	case "/chargebacks/cancel":
		return "chargeback_cancel"
	case "/transfers":
		return "transfer"
	case "/transfers/cancel":
		return "transfer_cancel"
	case "/acs/getOTPChannel":
		return "otp_channel"
	case "/acs/verifyTransaction":
		return "verify_transaction"
	case "/xpays/otp":
		return "xpays_otp"
	case "/xpays/customProvisioningData":
		return "custom_provisioning"
	default:
		return "unknown"
	}
}

func (s *Server) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time) {
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode: http.StatusMethodNotAllowed,
		Duration:   time.Since(start),
	})
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, _ []byte) {
	resp, err := s.handler.HandleStatus(ctx)
	statusCode := http.StatusOK
	if err != nil {
		statusCode = http.StatusServiceUnavailable
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: statusCode,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, statusCode, err)
		return
	}
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode: statusCode,
		Duration:   time.Since(start),
		Approved:   true,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handlePurchase(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req PurchaseRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	// Enrich request info with transaction details
	reqInfo.TransactionID = req.PurchaseID
	reqInfo.AccountID = req.AccountID
	if req.Card != nil {
		reqInfo.CardID = req.Card.PaysmartID
	}
	if req.TotalAmount != nil {
		reqInfo.Amount = req.TotalAmount.Amount
	}

	resp, err := s.handler.HandlePurchase(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := mapResponseCodeToHTTP(resp.Code)
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode:   statusCode,
		ResponseCode: resp.Code,
		Duration:     time.Since(start),
		Approved:     resp.Code == 0,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handlePurchaseCancel(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req PurchaseCancellationRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.TransactionID = req.CancellationID
	reqInfo.AccountID = req.AccountID
	if req.Card != nil {
		reqInfo.CardID = req.Card.PaysmartID
	}

	resp, err := s.handler.HandlePurchaseCancel(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := mapResponseCodeToHTTP(resp.Code)
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode:   statusCode,
		ResponseCode: resp.Code,
		Duration:     time.Since(start),
		Approved:     resp.Code == 0,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handleQuery(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req QueryRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.TransactionID = req.QueryID
	reqInfo.AccountID = req.AccountID
	if req.Card != nil {
		reqInfo.CardID = req.Card.PaysmartID
	}

	resp, err := s.handler.HandleQuery(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := mapResponseCodeToHTTP(resp.Code)
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode:   statusCode,
		ResponseCode: resp.Code,
		Duration:     time.Since(start),
		Approved:     resp.Code == 0,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handleWithdrawal(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req WithdrawalRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.TransactionID = req.WithdrawalID
	reqInfo.AccountID = req.AccountID
	if req.Card != nil {
		reqInfo.CardID = req.Card.PaysmartID
	}
	if req.TotalAmount != nil {
		reqInfo.Amount = req.TotalAmount.Amount
	}

	resp, err := s.handler.HandleWithdrawal(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := mapResponseCodeToHTTP(resp.Code)
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode:   statusCode,
		ResponseCode: resp.Code,
		Duration:     time.Since(start),
		Approved:     resp.Code == 0,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handleWithdrawalQuery(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req WithdrawalQueryRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.TransactionID = req.QueryID
	reqInfo.AccountID = req.AccountID
	if req.Card != nil {
		reqInfo.CardID = req.Card.PaysmartID
	}

	resp, err := s.handler.HandleWithdrawalQuery(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := mapResponseCodeToHTTP(resp.Code)
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode:   statusCode,
		ResponseCode: resp.Code,
		Duration:     time.Since(start),
		Approved:     resp.Code == 0,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handleWithdrawalCancel(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req WithdrawalCancellationRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.TransactionID = req.CancellationID
	reqInfo.AccountID = req.AccountID
	if req.Card != nil {
		reqInfo.CardID = req.Card.PaysmartID
	}

	resp, err := s.handler.HandleWithdrawalCancel(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := mapResponseCodeToHTTP(resp.Code)
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode:   statusCode,
		ResponseCode: resp.Code,
		Duration:     time.Since(start),
		Approved:     resp.Code == 0,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handleChargeback(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req ChargebackRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.TransactionID = req.ChargebackID
	reqInfo.AccountID = req.AccountID
	if req.Card != nil {
		reqInfo.CardID = req.Card.PaysmartID
	}
	if req.ChargebackAmount != nil {
		reqInfo.Amount = req.ChargebackAmount.Amount
	}

	resp, err := s.handler.HandleChargeback(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := mapResponseCodeToHTTP(resp.Code)
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode:   statusCode,
		ResponseCode: resp.Code,
		Duration:     time.Since(start),
		Approved:     resp.Code == 0,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handleChargebackCancel(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req ChargebackCancellationRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.TransactionID = req.CancellationID
	reqInfo.AccountID = req.AccountID
	if req.Card != nil {
		reqInfo.CardID = req.Card.PaysmartID
	}

	resp, err := s.handler.HandleChargebackCancel(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := mapResponseCodeToHTTP(resp.Code)
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode:   statusCode,
		ResponseCode: resp.Code,
		Duration:     time.Since(start),
		Approved:     resp.Code == 0,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handleTransfer(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req TransferRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.TransactionID = req.TransferID
	reqInfo.AccountID = req.AccountID
	if req.SourceCard != nil {
		reqInfo.CardID = req.SourceCard.PaysmartID
	}
	if req.TotalAmount != nil {
		reqInfo.Amount = req.TotalAmount.Amount
	}

	resp, err := s.handler.HandleTransfer(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := mapResponseCodeToHTTP(resp.Code)
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode:   statusCode,
		ResponseCode: resp.Code,
		Duration:     time.Since(start),
		Approved:     resp.Code == 0,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handleTransferCancel(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req TransferCancellationRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.TransactionID = req.CancellationID
	reqInfo.AccountID = req.AccountID
	if req.SourceCard != nil {
		reqInfo.CardID = req.SourceCard.PaysmartID
	}

	resp, err := s.handler.HandleTransferCancel(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := mapResponseCodeToHTTP(resp.Code)
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode:   statusCode,
		ResponseCode: resp.Code,
		Duration:     time.Since(start),
		Approved:     resp.Code == 0,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handleGetOTPChannel(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req OTPChannelRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.TransactionID = req.TransactionID
	reqInfo.AccountID = req.AccountID
	reqInfo.CardID = req.CardID

	resp, err := s.handler.HandleGetOTPChannel(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode: http.StatusOK,
		Duration:   time.Since(start),
		Approved:   true,
	})
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleVerifyTransaction(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req VerifyTransactionRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.TransactionID = req.TransactionID
	reqInfo.AccountID = req.AccountID
	reqInfo.CardID = req.CardID

	resp, err := s.handler.HandleVerifyTransaction(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := http.StatusOK
	if resp.Code != 0 {
		statusCode = mapResponseCodeToHTTP(resp.Code)
	}
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode:   statusCode,
		ResponseCode: resp.Code,
		Duration:     time.Since(start),
		Approved:     resp.Code == 0,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handleXPaysOTP(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req XPaysOTPRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.AccountID = req.AccountID
	reqInfo.CardID = req.CardID

	resp, err := s.handler.HandleXPaysOTP(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := http.StatusOK
	if resp.Code != 0 {
		statusCode = mapResponseCodeToHTTP(resp.Code)
	}
	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode:   statusCode,
		ResponseCode: resp.Code,
		Duration:     time.Since(start),
		Approved:     resp.Code == 0,
	})
	writeJSON(w, statusCode, resp)
}

func (s *Server) handleCustomProvisioningData(w http.ResponseWriter, r *http.Request, ctx context.Context, reqInfo *RequestInfo, start time.Time, bodyBytes []byte) {
	var req CustomProvisioningDataRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusBadRequest,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusBadRequest, err)
		return
	}

	reqInfo.AccountID = req.AccountID
	reqInfo.CardID = req.CardID

	resp, err := s.handler.HandleCustomProvisioningData(ctx, &req)
	if err != nil {
		s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
			StatusCode: http.StatusInternalServerError,
			Duration:   time.Since(start),
			Error:      err,
		})
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	s.executeAfterHooks(ctx, reqInfo, &ResponseInfo{
		StatusCode: http.StatusOK,
		Duration:   time.Since(start),
		Approved:   true,
	})
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	msg := err.Error()
	if statusCode >= 500 {
		msg = "internal server error"
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// executeBeforeHooks runs all before hooks and returns the modified context.
func (s *Server) executeBeforeHooks(ctx context.Context, req *RequestInfo) context.Context {
	for _, hook := range s.hooks {
		func() {
			defer func() {
				if r := recover(); r != nil && s.logger != nil {
					s.logger.Error("hook panic recovered",
						slog.Any("panic", r),
						slog.String("hook_phase", "BeforeAuthorization"))
				}
			}()
			ctx = hook.BeforeAuthorization(ctx, req)
		}()
	}
	return ctx
}

// executeAfterHooks runs all after hooks.
func (s *Server) executeAfterHooks(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
	for _, hook := range s.hooks {
		func() {
			defer func() {
				if r := recover(); r != nil && s.logger != nil {
					s.logger.Error("hook panic recovered",
						slog.Any("panic", r),
						slog.String("hook_phase", "AfterAuthorization"))
				}
			}()
			hook.AfterAuthorization(ctx, req, resp)
		}()
	}
}

// mapResponseCodeToHTTP maps Evertec response codes to HTTP status codes.
// Reference: https://paysmart-api.gitlab.io/processadora/PT-br/docs/codigos-retorno
func mapResponseCodeToHTTP(code int) int {
	switch code {
	case 0: // Approved
		return http.StatusOK // 200
	case 51: // Insufficient funds
		return 400
	case 61, 64, 65: // Exceeds limit
		return 400
	case 14, 25: // Invalid card/account
		return http.StatusNotFound // 404
	case 54: // Expired card
		return http.StatusNotFound // 404
	case 59: // Fraud suspect
		return StatusFraudSuspect
	case 62, 78: // Restricted/blocked card
		return http.StatusPreconditionFailed // 412
	case 57, 58: // Invalid MCC
		return StatusInvalidMCC
	case 91, 96: // System error
		return http.StatusServiceUnavailable // 503
	default:
		// For other denial codes, return 499 (Acknowledge - declined by processor)
		if code != 0 {
			return StatusDeclinedByProcessor
		}
		return http.StatusOK
	}
}
