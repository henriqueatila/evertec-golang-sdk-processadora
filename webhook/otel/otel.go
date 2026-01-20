// Package otel provides OpenTelemetry tracing integration for the Evertec Webhook Server.
//
// This package is a separate module to avoid adding OpenTelemetry dependencies
// to the core SDK. To use it, import this package and add the tracing hook:
//
//	import (
//	    "github.com/henriqueatila/evertec-golang-sdk-processadora/webhook"
//	    webhookotel "github.com/henriqueatila/evertec-golang-sdk-processadora/webhook/otel"
//	)
//
//	server := webhook.NewServer(webhook.Config{
//	    Handler: handler,
//	    Hooks:   []webhook.Hook{webhookotel.NewTracingHook()},
//	})
package otel

import (
	"context"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/webhook"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "github.com/henriqueatila/evertec-golang-sdk-processadora/webhook"
	sdkVersion = "1.0.0"
)

// TracingHook implements webhook.Hook with OpenTelemetry tracing.
type TracingHook struct {
	tracer trace.Tracer
}

// NewTracingHook creates a new tracing hook using the global tracer provider.
// This is the recommended way to create a tracing hook.
//
// Example:
//
//	// Setup your tracer provider first
//	tp := sdktrace.NewTracerProvider(...)
//	otel.SetTracerProvider(tp)
//
//	// Then create the hook
//	hook := otel.NewTracingHook()
func NewTracingHook() *TracingHook {
	return &TracingHook{
		tracer: otel.Tracer(tracerName),
	}
}

// NewTracingHookWithProvider creates a tracing hook with a custom tracer provider.
// Use this when you don't want to use the global tracer provider.
//
// Example:
//
//	tp := sdktrace.NewTracerProvider(...)
//	hook := otel.NewTracingHookWithProvider(tp)
func NewTracingHookWithProvider(tp trace.TracerProvider) *TracingHook {
	return &TracingHook{
		tracer: tp.Tracer(tracerName),
	}
}

type spanContextKey struct{}

// BeforeEvent starts a new span for the webhook event.
func (h *TracingHook) BeforeEvent(ctx context.Context, event *webhook.EventInfo) context.Context {
	spanName := "evertec.webhook"
	if event.EventType != "" {
		spanName = "evertec.webhook." + string(event.EventType)
	}

	ctx, span := h.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attribute.String("evertec.event_type", string(event.EventType)),
			attribute.Int("http.request_content_length", event.BodySize),
			attribute.String("rpc.system", "evertec"),
			attribute.String("rpc.service", "webhook"),
			attribute.String("evertec.sdk.name", "evertec-golang-sdk-processadora"),
			attribute.String("evertec.sdk.version", sdkVersion),
		),
	)

	if event.EventID != "" {
		span.SetAttributes(attribute.String("evertec.event_id", event.EventID))
	}
	if event.IsDuplicate {
		span.SetAttributes(attribute.Bool("evertec.duplicate", true))
	}

	return context.WithValue(ctx, spanContextKey{}, span)
}

// AfterEvent ends the span and records processing result attributes.
func (h *TracingHook) AfterEvent(ctx context.Context, event *webhook.EventInfo, result *webhook.ProcessingInfo) {
	span, ok := ctx.Value(spanContextKey{}).(trace.Span)
	if !ok || span == nil {
		return
	}
	defer span.End()

	span.SetAttributes(
		attribute.Int("http.status_code", result.StatusCode),
		attribute.Bool("evertec.processed", result.Processed),
		attribute.Int64("evertec.duration_ms", result.Duration.Milliseconds()),
	)

	if result.Error != nil {
		span.RecordError(result.Error)
		span.SetStatus(codes.Error, result.Error.Error())
	} else if result.Processed {
		span.SetStatus(codes.Ok, "processed")
	} else {
		span.SetStatus(codes.Ok, "skipped")
	}
}
