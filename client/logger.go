package client

import (
	"context"
	"log/slog"
)

// logRequest logs request information using the configured slog logger.
// Logs at INFO level for successful requests, ERROR level for failures.
func (c *Client) logRequest(ctx context.Context, req *RequestInfo, resp *ResponseInfo) {
	if c.logger == nil {
		return
	}

	// Determine log level based on response
	level := slog.LevelInfo
	if resp.Error != nil || resp.StatusCode >= 400 {
		level = slog.LevelError
	}

	// Build attributes
	attrs := []slog.Attr{
		slog.String("method", req.Method),
		slog.String("path", req.Path),
		slog.Int("status", resp.StatusCode),
		slog.Duration("duration", resp.Duration),
		slog.Int("request_size", req.BodySize),
		slog.Int("response_size", resp.BodySize),
	}

	// Add idempotency key if present
	if req.IdempotencyKey != "" {
		attrs = append(attrs, slog.String("idempotency_key", req.IdempotencyKey))
	}

	// Add error details if present
	if resp.Error != nil {
		attrs = append(attrs, slog.String("error", resp.Error.Error()))

		// Add structured API error info if available
		if apiErr, ok := resp.Error.(*APIError); ok {
			if apiErr.Code != "" {
				attrs = append(attrs, slog.String("error_code", apiErr.Code))
			}
			if apiErr.Message != "" {
				attrs = append(attrs, slog.String("error_message", apiErr.Message))
			}
		}
	}

	c.logger.LogAttrs(ctx, level, "paysmart_request", attrs...)
}
