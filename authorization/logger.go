package authorization

import (
	"context"
	"log/slog"
)

// LoggingHook implements Hook for structured logging using slog.
type LoggingHook struct {
	logger *slog.Logger
}

// NewLoggingHook creates a new logging hook with the given logger.
//
// Example:
//
//	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
//	server := authorization.NewServer(handler,
//	    authorization.WithHooks(authorization.NewLoggingHook(logger)),
//	)
func NewLoggingHook(logger *slog.Logger) *LoggingHook {
	return &LoggingHook{logger: logger}
}

func (h *LoggingHook) BeforeAuthorization(ctx context.Context, req *RequestInfo) context.Context {
	return ctx
}

func (h *LoggingHook) AfterAuthorization(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
	if h.logger == nil {
		return
	}

	// Determine log level based on response
	level := slog.LevelInfo
	if resp.Error != nil || resp.StatusCode >= 500 {
		level = slog.LevelError
	} else if !resp.Approved && resp.ResponseCode != 0 {
		level = slog.LevelWarn
	}

	// Build attributes
	attrs := []slog.Attr{
		slog.String("path", req.Path),
		slog.String("operation", req.OperationType),
		slog.Int("status", resp.StatusCode),
		slog.Int("response_code", resp.ResponseCode),
		slog.Bool("approved", resp.Approved),
		slog.Duration("duration", resp.Duration),
		slog.Int("request_size", req.BodySize),
	}

	// Add transaction details if present
	if req.TransactionID != "" {
		attrs = append(attrs, slog.String("transaction_id", req.TransactionID))
	}
	if req.AccountID != "" {
		attrs = append(attrs, slog.String("account_id", req.AccountID))
	}
	if req.CardID != "" {
		attrs = append(attrs, slog.String("card_id", req.CardID))
	}
	if req.Amount > 0 {
		attrs = append(attrs, slog.Int64("amount", req.Amount))
		if req.Currency != "" {
			attrs = append(attrs, slog.String("currency", req.Currency))
		}
	}

	// Add error details if present
	if resp.Error != nil {
		attrs = append(attrs, slog.String("error", resp.Error.Error()))
	}

	h.logger.LogAttrs(ctx, level, "evertec_authorization", attrs...)
}
