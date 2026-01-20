// Package otel provides OpenTelemetry tracing integration for the Evertec Authorization Server.
//
// This package is a separate module to avoid adding OpenTelemetry dependencies
// to the core SDK. To use it, import this package and add the tracing hook:
//
//	import (
//	    "github.com/henriqueatila/evertec-golang-sdk-processadora/authorization"
//	    authorizationotel "github.com/henriqueatila/evertec-golang-sdk-processadora/authorization/otel"
//	)
//
//	server := authorization.NewServer(handler,
//	    authorization.WithHooks(authorizationotel.NewTracingHook()),
//	)
package otel

import (
	"context"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/authorization"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "github.com/henriqueatila/evertec-golang-sdk-processadora/authorization"
	sdkVersion = "1.0.0"
)

// TracingHook implements authorization.Hook with OpenTelemetry tracing.
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

// BeforeAuthorization starts a new span for the authorization request.
func (h *TracingHook) BeforeAuthorization(ctx context.Context, req *authorization.RequestInfo) context.Context {
	ctx, span := h.tracer.Start(ctx, "evertec.authorization."+req.OperationType,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attribute.String("evertec.operation", req.OperationType),
			attribute.String("http.route", req.Path),
			attribute.Int("http.request_content_length", req.BodySize),
			attribute.String("rpc.system", "evertec"),
			attribute.String("rpc.service", "authorization"),
			attribute.String("evertec.sdk.name", "evertec-golang-sdk-processadora"),
			attribute.String("evertec.sdk.version", sdkVersion),
		),
	)

	if req.TransactionID != "" {
		span.SetAttributes(attribute.String("evertec.transaction_id", req.TransactionID))
	}
	if req.AccountID != "" {
		span.SetAttributes(attribute.String("evertec.account_id", req.AccountID))
	}
	if req.CardID != "" {
		span.SetAttributes(attribute.String("evertec.card_id", req.CardID))
	}
	if req.Amount > 0 {
		span.SetAttributes(
			attribute.Int64("evertec.amount", req.Amount),
			attribute.String("evertec.currency", req.Currency),
		)
	}

	return context.WithValue(ctx, spanContextKey{}, span)
}

// AfterAuthorization ends the span and records response attributes.
func (h *TracingHook) AfterAuthorization(ctx context.Context, req *authorization.RequestInfo, resp *authorization.ResponseInfo) {
	span, ok := ctx.Value(spanContextKey{}).(trace.Span)
	if !ok || span == nil {
		return
	}
	defer span.End()

	span.SetAttributes(
		attribute.Int("http.status_code", resp.StatusCode),
		attribute.Int("evertec.response_code", resp.ResponseCode),
		attribute.Bool("evertec.approved", resp.Approved),
		attribute.Int64("evertec.duration_ms", resp.Duration.Milliseconds()),
	)

	if resp.Error != nil {
		span.RecordError(resp.Error)
		span.SetStatus(codes.Error, resp.Error.Error())
	} else if !resp.Approved {
		span.SetStatus(codes.Ok, "declined")
	} else {
		span.SetStatus(codes.Ok, "approved")
	}
}
