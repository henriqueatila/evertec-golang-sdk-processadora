package webhook

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
//	server := webhook.NewServer(webhook.Config{
//	    Handler: handler,
//	    Hooks:   []webhook.Hook{webhook.NewLoggingHook(logger)},
//	})
func NewLoggingHook(logger *slog.Logger) *LoggingHook {
	return &LoggingHook{logger: logger}
}

func (h *LoggingHook) BeforeEvent(ctx context.Context, event *EventInfo) context.Context {
	return ctx
}

func (h *LoggingHook) AfterEvent(ctx context.Context, event *EventInfo, result *ProcessingInfo) {
	if h.logger == nil {
		return
	}

	// Determine log level based on result
	level := slog.LevelInfo
	if result.Error != nil || result.StatusCode >= 500 {
		level = slog.LevelError
	} else if !result.Processed {
		level = slog.LevelWarn
	}

	// Build attributes
	attrs := []slog.Attr{
		slog.String("event_type", string(event.EventType)),
		slog.Int("status", result.StatusCode),
		slog.Bool("processed", result.Processed),
		slog.Duration("duration", result.Duration),
		slog.Int("request_size", event.BodySize),
	}

	// Add event ID if present
	if event.EventID != "" {
		attrs = append(attrs, slog.String("event_id", event.EventID))
	}

	// Add duplicate flag if true
	if event.IsDuplicate {
		attrs = append(attrs, slog.Bool("duplicate", true))
	}

	// Add error details if present
	if result.Error != nil {
		attrs = append(attrs, slog.String("error", result.Error.Error()))
	}

	h.logger.LogAttrs(ctx, level, "evertec_webhook", attrs...)
}
