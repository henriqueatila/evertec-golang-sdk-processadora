# Evertec Processadora SDK for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/henriqueatila/evertec-golang-sdk-processadora.svg)](https://pkg.go.dev/github.com/henriqueatila/evertec-golang-sdk-processadora)
[![CI](https://github.com/henriqueatila/evertec-golang-sdk-processadora/actions/workflows/ci.yml/badge.svg)](https://github.com/henriqueatila/evertec-golang-sdk-processadora/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/henriqueatila/evertec-golang-sdk-processadora/branch/master/graph/badge.svg)](https://codecov.io/gh/henriqueatila/evertec-golang-sdk-processadora)
[![Go Report Card](https://goreportcard.com/badge/github.com/henriqueatila/evertec-golang-sdk-processadora)](https://goreportcard.com/report/github.com/henriqueatila/evertec-golang-sdk-processadora)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go SDK for the [Evertec Processadora](https://paysmart-api.gitlab.io/processadora/) card processing API.

## Features

- **113 API Methods** — Complete coverage of 19 API domains
- **mTLS Authentication** — Secure mutual TLS as required by the API
- **Authorization Server** — Real-time transaction authorization callbacks with `context.Context`
- **Webhook Handler** — Asynchronous event processing with `context.Context`
- **Retry with Backoff** — Exponential backoff with jitter and Retry-After support
- **Circuit Breaker** — Protects against cascading failures from upstream errors
- **Certificate Rotation** — Hot-reload mTLS certificates without restart
- **Auto Validation** — Request body validation before sending
- **Observability** — slog logging, metrics hooks, OpenTelemetry tracing
- **Type-Safe** — 263 struct types, 34 enum types with validation
- **High Coverage** — 91.3-100% test coverage with race detection

## Installation

```bash
go get github.com/henriqueatila/evertec-golang-sdk-processadora
```

Requires Go 1.21+

## Quick Start

Create the client with mTLS certificates:

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/henriqueatila/evertec-golang-sdk-processadora/client"
    "github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

func main() {
    c, err := client.NewWithCertFiles(
        "your-api-key",
        "YourBank/1.0.0",
        "certs/client.crt",
        "certs/client.key",
        client.WithHomolog(),
        client.WithTimeout(30*time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    account, err := c.CreateAccount(ctx, &types.CreateAccountRequest{
        ProductID: "prd-credit-gold",
        Holder: &types.Holder{
            Name:         "João Silva",
            Document:     "12345678900",
            DocumentType: "CPF",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Account created: %s", account.AccountID)
}
```

## API Reference

See [API domains table](./docs/codebase-summary.md#api-domains) for full method list.

### Core Domains
- **Accounts** — Account creation, management, balance queries
- **Cards** — Physical/virtual card management, PIN operations
- **Transactions** — List, query, manage transaction lifecycle
- **Disputes** — Create, track, resolve card disputes
- **Statements** — Post-paid account statements and invoices
- **xPays** — Digital wallet provisioning (Apple Pay, Google Pay)
- **QR Code** — QR code payment processing
- Plus 12 more domains (HCE, Travel, Campaigns, Cobranded, etc.)

## Authorization Server

Handle real-time purchase/withdrawal authorizations:

```go
import (
    "context"
    "log"
    "net/http"

    "github.com/henriqueatila/evertec-golang-sdk-processadora/authorization"
)

type AuthHandler struct{}

func (h *AuthHandler) HandlePurchase(ctx context.Context, req *authorization.PurchaseRequest) (*authorization.PurchaseResponse, error) {
    log.Printf("[AUTH] Purchase: account=%s amount=%d", req.AccountID, req.TotalAmount.Amount)
    // Validate card, check balance, apply fraud rules
    return &authorization.PurchaseResponse{
        Code:            0, // approved
        AuthorizationID: 123456,
        Balance: &authorization.AuthAmount{Amount: 500000, CurrencyCode: 986},
    }, nil
}

func main() {
    server := authorization.NewServer(&AuthHandler{})
    log.Fatal(http.ListenAndServeTLS(":8443", "server.crt", "server.key", server))
}
```

See [Authorization Server Guide](./docs/code-standards.md#webhook--authorization-server) for all 16 handler methods.

## Webhook Handler

Process asynchronous EventHub notifications:

```go
import (
    "context"
    "log"
    "net/http"

    "github.com/henriqueatila/evertec-golang-sdk-processadora/webhook"
)

type EventHandler struct {
    webhook.BaseHandler // no-op defaults
}

func (h *EventHandler) OnCardStatusChange(ctx context.Context, event *webhook.CardStatusChangeEvent) error {
    log.Printf("[WEBHOOK] Card %v: %s → %s",
        event.CardIDList,
        event.OldStatus.Status,
        event.NewStatus.Status)
    return nil
}

func main() {
    server := webhook.NewServer(webhook.Config{
        Handler:       &EventHandler{},
        WebhookSecret: "your-webhook-secret",
    })
    http.Handle("/webhooks/evertec", server)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Error Handling

Handle API errors with type checking or sentinel errors:

```go
var apiErr *client.APIError
if errors.As(err, &apiErr) {
    log.Printf("API Error [%d] %s", apiErr.StatusCode, apiErr.Code)
    // Handle specific errors: ACCOUNT_ALREADY_EXISTS, INVALID_DOCUMENT, etc.
}

// Or use sentinel errors for quick checks:
if errors.Is(err, client.ErrAccountNotFound) {
    // handle missing account
}
if errors.Is(err, client.ErrCardBlocked) {
    // handle blocked card
}
```

## Observability

### Logging

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
c, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
    client.WithLogger(logger),
)
```

### Metrics

```go
metrics := client.NewSimpleMetrics()
c, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
    client.WithHooks(client.NewMetricsHook(metrics)),
)
total, errors, avgMs := metrics.Stats()
```

### OpenTelemetry

```go
c, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
    client.WithHooks(otel.NewTracingHook()),
)
```

## Configuration

```go
// Environment
c, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
    client.WithHomolog(),  // or WithBaseURL("https://custom.api")
    client.WithTimeout(60*time.Second),
)

// Idempotency (safe retries)
ctx := client.WithIdempotencyKey(context.Background(), "uuid-v4-key")
resp, err := c.CreateAccount(ctx, req)
```

## Resilience

### Retry with Exponential Backoff

```go
c, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
    client.WithRetry(
        client.MaxRetries(5),
        client.InitialDelay(200*time.Millisecond),
        client.MaxDelay(30*time.Second),
    ),
)
```

Retries on 429, 500, 502, 503, 504 with jitter. Respects `Retry-After` header.

### Circuit Breaker

```go
c, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
    client.WithCircuitBreaker(
        client.FailureThreshold(5),
        client.ResetTimeout(60*time.Second),
    ),
)
```

Opens after consecutive 5xx failures. Returns `client.ErrCircuitOpen` when open.

### Rate Limiter

```go
c, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
    client.WithRateLimiter(
        client.RequestsPerSecond(50),
        client.BurstSize(100),
        client.WaitTimeout(time.Second),
    ),
)
```

Token bucket algorithm. Prevents flooding during retry storms. Returns `client.ErrRateLimited` when exhausted.

### Certificate Rotation

```go
c, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
    client.WithCertRotation(),
)

// Force manual reload
err := c.RefreshCertificates()
```

Hot-reloads certificates from disk on each TLS handshake. Falls back to last valid cert on error.

### Auto Request Validation

Request types implementing `Validatable` are validated before sending (enabled by default):

```go
// Disable if needed
c, _ := client.NewWithCertFiles(apiKey, userAgent, cert, key,
    client.WithNoValidation(),
)
```

### Health Check & Circuit Breaker Metrics

```go
// Health check (includes CB state + rate limiter status)
status := c.Health()
fmt.Printf("Healthy: %v, CB: %s\n", status.Healthy, status.CircuitBreaker.State)

// Circuit breaker metrics
m := c.CircuitBreakerMetrics()
fmt.Printf("Requests: %d, Failures: %d, Rejected: %d, Transitions: %d\n",
    m.TotalRequests, m.TotalFailures, m.TotalRejected, m.StateTransitions)
```

## Migrating from v1 to v2

### Breaking Changes

**Handler interfaces now require `context.Context` as the first parameter:**

```go
// v1
func (h *MyHandler) HandlePurchase(req *authorization.PurchaseRequest) (*authorization.PurchaseResponse, error)
func (h *MyHandler) OnCardStatusChange(event *webhook.CardStatusChangeEvent) error

// v2
func (h *MyHandler) HandlePurchase(ctx context.Context, req *authorization.PurchaseRequest) (*authorization.PurchaseResponse, error)
func (h *MyHandler) OnCardStatusChange(ctx context.Context, event *webhook.CardStatusChangeEvent) error
```

All 15 authorization handler methods and all 5 webhook handler methods require this change.

### New Features (opt-in)

No code changes required for new features. They are opt-in via functional options:
- `client.WithRetry()` — Retry with exponential backoff
- `client.WithCircuitBreaker()` — Circuit breaker pattern
- `client.WithRateLimiter()` — Client-side rate limiting
- `client.WithCertRotation()` — Certificate hot-reload
- `client.WithNoValidation()` — Disable auto-validation (auto-validation is on by default)
- `client.Health()` — Health check with CB state and rate limiter status
- `client.CircuitBreakerMetrics()` — CB counters (requests, failures, rejected, transitions)

## Testing

```bash
go test ./...              # All tests
go test -race ./...        # With race detection
go test -cover ./...       # With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Documentation

- [Project Overview & PDR](./docs/project-overview-pdr.md) — Requirements and specs
- [Code Standards](./docs/code-standards.md) — Coding patterns and best practices
- [Codebase Summary](./docs/codebase-summary.md) — Code structure and organization
- [System Architecture](./docs/system-architecture.md) — High-level design
- [Official API Docs](https://paysmart-api.gitlab.io/processadora/) — Evertec documentation

## Security

- mTLS enforced on all requests
- TLS 1.2+ minimum
- Webhook HMAC-SHA256 signature verification
- Sensitive data masking (PAN, CVV, PIN)
- 69 official error codes with Portuguese messages

## CI/CD Security

- Automated vulnerability scanning (gosec)
- Code analysis (CodeQL)
- Dependency review on PRs
- SBOM generation (SPDX)
- golangci-lint checks

## Project Statistics

| Metric | Value |
|--------|-------|
| Total Code | ~27K LOC |
| Test Coverage | 91.3-100% |
| API Methods | 113 |
| Struct Types | 263 |
| Dependencies | 1 (google/uuid) |
| Go Version | 1.21+ |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.
