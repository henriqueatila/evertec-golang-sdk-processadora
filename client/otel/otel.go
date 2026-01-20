// Package otel provides OpenTelemetry tracing integration for the Evertec SDK.
//
// This package is a separate module to avoid adding OpenTelemetry dependencies
// to the core SDK. To use it, import this package and add the tracing hook:
//
//	import (
//	    "github.com/henriqueatila/evertec-golang-sdk-processadora/client"
//	    "github.com/henriqueatila/evertec-golang-sdk-processadora/client/otel"
//	)
//
//	c, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
//	    client.WithHooks(otel.NewTracingHook()),
//	)
package otel

import (
	"context"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/client"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "github.com/henriqueatila/evertec-golang-sdk-processadora"
	sdkVersion = "1.0.0"
)

// TracingHook implements client.Hook with OpenTelemetry tracing.
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

// BeforeRequest starts a new span for the request.
func (h *TracingHook) BeforeRequest(ctx context.Context, req *client.RequestInfo) context.Context {
	ctx, span := h.tracer.Start(ctx, req.Method+" "+req.Path,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("http.method", req.Method),
			attribute.String("http.url", req.Path),
			attribute.Int("http.request_content_length", req.BodySize),
			attribute.String("rpc.system", "paysmart"),
			attribute.String("paysmart.sdk.name", "evertec-golang-sdk-processadora"),
			attribute.String("paysmart.sdk.version", sdkVersion),
		),
	)

	if req.IdempotencyKey != "" {
		span.SetAttributes(attribute.String("paysmart.idempotency_key", req.IdempotencyKey))
	}

	return context.WithValue(ctx, spanContextKey{}, span)
}

// AfterResponse ends the span and records response attributes.
func (h *TracingHook) AfterResponse(ctx context.Context, req *client.RequestInfo, resp *client.ResponseInfo) {
	span, ok := ctx.Value(spanContextKey{}).(trace.Span)
	if !ok || span == nil {
		return
	}
	defer span.End()

	span.SetAttributes(
		attribute.Int("http.status_code", resp.StatusCode),
		attribute.Int("http.response_content_length", resp.BodySize),
		attribute.Int64("paysmart.duration_ms", resp.Duration.Milliseconds()),
	)

	if resp.Error != nil {
		span.RecordError(resp.Error)
		span.SetStatus(codes.Error, resp.Error.Error())

		// Add API error details if available
		if apiErr, ok := resp.Error.(*client.APIError); ok {
			if apiErr.Code != "" {
				span.SetAttributes(attribute.String("paysmart.error_code", apiErr.Code))
			}
		}
	} else if resp.StatusCode >= 400 {
		span.SetStatus(codes.Error, "HTTP error")
	} else {
		span.SetStatus(codes.Ok, "")
	}
}
