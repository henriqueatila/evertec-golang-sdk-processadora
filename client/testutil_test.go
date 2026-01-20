package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/henriqueatila/evertec-golang-sdk-processadora/types"
)

// ============================================================================
// LEGACY COMPATIBILITY FUNCTIONS (for existing tests)
// ============================================================================

// mockServer creates a test server that returns the given response (LEGACY)
func mockServer(t *testing.T, method, path string, statusCode int, response interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != method {
			t.Errorf("Method = %s, want %s", r.Method, method)
		}

		// Verify path (without base path)
		if !strings.HasSuffix(r.URL.Path, path) {
			t.Errorf("Path = %s, want suffix %s", r.URL.Path, path)
		}

		// Verify required headers
		if r.Header.Get("API-Key") == "" {
			t.Error("Missing API-Key header")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if response != nil {
			_ = json.NewEncoder(w).Encode(response)
		}
	}))
}

// testConfig creates a test client config (LEGACY)
func testConfig(baseURL string) Config {
	return Config{
		APIKey:    "test-key",
		UserAgent: "test/1.0",
		BaseURL:   baseURL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

// ============================================================================
// PROFESSIONAL TEST UTILITIES
// ============================================================================

// TestCase represents a generic test case structure for table-driven tests.
type TestCase[Req any, Resp any] struct {
	Name           string
	Request        Req
	MockStatus     int
	MockResponse   interface{}
	MockError      *types.Error
	WantErr        bool
	WantErrType    error
	WantErrContain string
	Validate       func(t *testing.T, resp Resp)
	ValidateReq    func(t *testing.T, r *http.Request, body []byte)
}

// MockServerConfig configures the mock server behavior.
type MockServerConfig struct {
	Method       string
	Path         string
	Status       int
	Response     interface{}
	Error        *types.Error
	ValidateReq  func(t *testing.T, r *http.Request, body []byte)
	Delay        time.Duration
	FailAfter    int // Fail after N requests (for retry testing)
	requestCount *atomic.Int32
}

// MockServer wraps httptest.Server with enhanced capabilities.
type MockServer struct {
	*httptest.Server
	t              *testing.T
	config         *MockServerConfig
	requestHistory []RequestRecord
	mu             sync.Mutex
}

// RequestRecord records details of each request for verification.
type RequestRecord struct {
	Method      string
	Path        string
	Headers     http.Header
	QueryParams map[string][]string
	Body        []byte
	Timestamp   time.Time
}

// NewMockServer creates a new enhanced mock server.
func NewMockServer(t *testing.T, config *MockServerConfig) *MockServer {
	t.Helper()

	if config.requestCount == nil {
		config.requestCount = &atomic.Int32{}
	}

	ms := &MockServer{
		t:              t,
		config:         config,
		requestHistory: make([]RequestRecord, 0),
	}

	ms.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ms.handleRequest(w, r)
	}))

	return ms
}

func (ms *MockServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	// Record the request
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	ms.mu.Lock()
	ms.requestHistory = append(ms.requestHistory, RequestRecord{
		Method:      r.Method,
		Path:        r.URL.Path,
		Headers:     r.Header.Clone(),
		QueryParams: r.URL.Query(),
		Body:        body,
		Timestamp:   time.Now(),
	})
	ms.mu.Unlock()

	// Verify method
	if ms.config.Method != "" && r.Method != ms.config.Method {
		ms.t.Errorf("MockServer: Method = %s, want %s", r.Method, ms.config.Method)
	}

	// Verify path (suffix match to account for base path)
	if ms.config.Path != "" && !strings.HasSuffix(r.URL.Path, ms.config.Path) {
		ms.t.Errorf("MockServer: Path = %s, want suffix %s", r.URL.Path, ms.config.Path)
	}

	// Verify required headers
	ms.verifyRequiredHeaders(r)

	// Custom request validation
	if ms.config.ValidateReq != nil {
		ms.config.ValidateReq(ms.t, r, body)
	}

	// Simulate delay if configured
	if ms.config.Delay > 0 {
		time.Sleep(ms.config.Delay)
	}

	// Check if should fail for retry testing
	count := ms.config.requestCount.Add(1)
	if ms.config.FailAfter > 0 && int(count) <= ms.config.FailAfter {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")

	if ms.config.Error != nil {
		w.WriteHeader(ms.config.Status)
		_ = json.NewEncoder(w).Encode(ms.config.Error)
		return
	}

	w.WriteHeader(ms.config.Status)
	if ms.config.Response != nil {
		_ = json.NewEncoder(w).Encode(ms.config.Response)
	}
}

func (ms *MockServer) verifyRequiredHeaders(r *http.Request) {
	// API-Key is required
	if r.Header.Get("API-Key") == "" {
		ms.t.Error("MockServer: Missing required API-Key header")
	}

	// User-Agent is required
	if r.Header.Get("User-Agent") == "" {
		ms.t.Error("MockServer: Missing required User-Agent header")
	}

	// Content-Type should be JSON for POST/PATCH/PUT
	if r.Method != http.MethodGet && r.Method != http.MethodDelete {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			ms.t.Errorf("MockServer: Content-Type = %s, want application/json", ct)
		}
	}

	// IdempotencyKey required for non-GET requests
	if r.Method != http.MethodGet {
		if r.Header.Get("IdempotencyKey") == "" {
			ms.t.Error("MockServer: Missing required IdempotencyKey header for non-GET request")
		}
	}
}

// GetRequests returns all recorded requests.
func (ms *MockServer) GetRequests() []RequestRecord {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return append([]RequestRecord{}, ms.requestHistory...)
}

// GetLastRequest returns the most recent request.
func (ms *MockServer) GetLastRequest() *RequestRecord {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if len(ms.requestHistory) == 0 {
		return nil
	}
	return &ms.requestHistory[len(ms.requestHistory)-1]
}

// RequestCount returns the number of requests received.
func (ms *MockServer) RequestCount() int {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return len(ms.requestHistory)
}

// Reset clears the request history.
func (ms *MockServer) Reset() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.requestHistory = make([]RequestRecord, 0)
	ms.config.requestCount.Store(0)
}

// ============================================================================
// TEST CLIENT FACTORY
// ============================================================================

// TestClientConfig holds test client configuration.
type TestClientConfig struct {
	APIKey    string
	UserAgent string
	Timeout   time.Duration
}

// DefaultTestConfig returns default test configuration.
func DefaultTestConfig() TestClientConfig {
	return TestClientConfig{
		APIKey:    "test-api-key-12345",
		UserAgent: "TestEmissor/1.0.0",
		Timeout:   30 * time.Second,
	}
}

// NewTestClient creates a test client connected to a mock server.
func NewTestClient(t *testing.T, server *MockServer, cfg ...TestClientConfig) *Client {
	t.Helper()

	config := DefaultTestConfig()
	if len(cfg) > 0 {
		config = cfg[0]
	}

	clientCfg := Config{
		APIKey:    config.APIKey,
		UserAgent: config.UserAgent,
		BaseURL:   server.URL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Timeout:   config.Timeout,
	}

	client, err := New(clientCfg)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	return client
}

// ============================================================================
// ASSERTION HELPERS
// ============================================================================

// AssertNoError fails the test if err is not nil.
func AssertNoError(t *testing.T, err error, msg ...string) {
	t.Helper()
	if err != nil {
		prefix := ""
		if len(msg) > 0 {
			prefix = msg[0] + ": "
		}
		t.Fatalf("%sunexpected error: %v", prefix, err)
	}
}

// AssertError fails the test if err is nil.
func AssertError(t *testing.T, err error, msg ...string) {
	t.Helper()
	if err == nil {
		prefix := ""
		if len(msg) > 0 {
			prefix = msg[0] + ": "
		}
		t.Fatalf("%sexpected error but got nil", prefix)
	}
}

// AssertErrorContains fails if the error doesn't contain the substring.
func AssertErrorContains(t *testing.T, err error, substr string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error containing %q but got nil", substr)
	}
	if !strings.Contains(err.Error(), substr) {
		t.Errorf("error %q does not contain %q", err.Error(), substr)
	}
}

// AssertAPIError asserts the error is an APIError with expected status.
func AssertAPIError(t *testing.T, err error, expectedStatus int) *APIError {
	t.Helper()
	if err == nil {
		t.Fatalf("expected APIError but got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != expectedStatus {
		t.Errorf("APIError.StatusCode = %d, want %d", apiErr.StatusCode, expectedStatus)
	}
	return apiErr
}

// AssertEqual fails if expected != actual.
func AssertEqual(t *testing.T, expected, actual interface{}, field string) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s = %v (%T), want %v (%T)", field, actual, actual, expected, expected)
	}
}

// AssertNotNil fails if value is nil.
func AssertNotNil(t *testing.T, value interface{}, name string) {
	t.Helper()
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		t.Fatalf("%s is nil", name)
	}
}

// AssertNil fails if value is not nil.
func AssertNil(t *testing.T, value interface{}, name string) {
	t.Helper()
	if value != nil && !(reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		t.Fatalf("%s should be nil, got %v", name, value)
	}
}

// AssertJSONContains verifies the JSON body contains expected fields.
func AssertJSONContains(t *testing.T, body []byte, expected map[string]interface{}) {
	t.Helper()
	var actual map[string]interface{}
	if err := json.Unmarshal(body, &actual); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}
	for key, expectedVal := range expected {
		actualVal, ok := actual[key]
		if !ok {
			t.Errorf("missing expected field %q in body", key)
			continue
		}
		if !reflect.DeepEqual(expectedVal, actualVal) {
			t.Errorf("field %q = %v, want %v", key, actualVal, expectedVal)
		}
	}
}

// AssertQueryParam verifies a query parameter value.
func AssertQueryParam(t *testing.T, r *http.Request, key, expected string) {
	t.Helper()
	actual := r.URL.Query().Get(key)
	if actual != expected {
		t.Errorf("query param %q = %q, want %q", key, actual, expected)
	}
}

// AssertHeader verifies a header value.
func AssertHeader(t *testing.T, r *http.Request, key, expected string) {
	t.Helper()
	actual := r.Header.Get(key)
	if actual != expected {
		t.Errorf("header %q = %q, want %q", key, actual, expected)
	}
}

// ============================================================================
// TEST DATA BUILDERS
// ============================================================================

// AccountBuilder creates test Account data.
type AccountBuilder struct {
	account types.Account
}

// NewAccountBuilder creates a new AccountBuilder with defaults.
func NewAccountBuilder() *AccountBuilder {
	status := types.AccountStatusActive
	return &AccountBuilder{
		account: types.Account{
			AccountID:     "acc-test-001",
			Status:        &status,
			PsProductCode: "CREDIT_CARD",
			AccountOwner: map[string]interface{}{
				"fullName":               "Test User",
				"identityDocumentNumber": "12345678901",
			},
		},
	}
}

func (b *AccountBuilder) WithID(id string) *AccountBuilder {
	b.account.AccountID = id
	return b
}

func (b *AccountBuilder) WithStatus(status types.AccountStatus) *AccountBuilder {
	b.account.Status = &status
	return b
}

func (b *AccountBuilder) WithAccountOwner(owner interface{}) *AccountBuilder {
	b.account.AccountOwner = owner
	return b
}

func (b *AccountBuilder) Build() types.Account {
	return b.account
}

func (b *AccountBuilder) BuildPtr() *types.Account {
	acc := b.account
	return &acc
}

// CardBuilder creates test Card data.
type CardBuilder struct {
	card types.CardDetails
}

// NewCardBuilder creates a new CardBuilder with defaults.
func NewCardBuilder() *CardBuilder {
	return &CardBuilder{
		card: types.CardDetails{
			CardID:         "card-test-001",
			AccountID:      "acc-test-001",
			Status:         types.CardStatusActive,
			Physical:       true,
			Last4Digits:    "1234",
			ExpirationDate: "12/2028",
			Profile:        "credit",
		},
	}
}

func (b *CardBuilder) WithID(id string) *CardBuilder {
	b.card.CardID = id
	return b
}

func (b *CardBuilder) WithStatus(status types.CardStatus) *CardBuilder {
	b.card.Status = status
	return b
}

func (b *CardBuilder) WithPhysical(physical bool) *CardBuilder {
	b.card.Physical = physical
	return b
}

func (b *CardBuilder) WithAccountID(id string) *CardBuilder {
	b.card.AccountID = id
	return b
}

func (b *CardBuilder) Build() types.CardDetails {
	return b.card
}

func (b *CardBuilder) BuildPtr() *types.CardDetails {
	card := b.card
	return &card
}

// AmountBuilder creates test Amount data.
type AmountBuilder struct {
	amount types.Amount
}

// NewAmountBuilder creates a new AmountBuilder with defaults.
func NewAmountBuilder() *AmountBuilder {
	return &AmountBuilder{
		amount: types.Amount{
			Amount:       10000, // R$ 100,00
			CurrencyCode: 986,   // BRL
		},
	}
}

func (b *AmountBuilder) WithValue(value int64) *AmountBuilder {
	b.amount.Amount = value
	return b
}

func (b *AmountBuilder) WithCurrency(currency int) *AmountBuilder {
	b.amount.CurrencyCode = currency
	return b
}

func (b *AmountBuilder) Build() types.Amount {
	return b.amount
}

func (b *AmountBuilder) BuildPtr() *types.Amount {
	amt := b.amount
	return &amt
}

// ============================================================================
// CONCURRENCY TESTING UTILITIES
// ============================================================================

// ConcurrentTestRunner executes tests concurrently.
type ConcurrentTestRunner struct {
	t        *testing.T
	wg       sync.WaitGroup
	errors   []error
	errorsMu sync.Mutex
}

// NewConcurrentTestRunner creates a new concurrent test runner.
func NewConcurrentTestRunner(t *testing.T) *ConcurrentTestRunner {
	return &ConcurrentTestRunner{
		t:      t,
		errors: make([]error, 0),
	}
}

// Run executes a function concurrently n times.
func (r *ConcurrentTestRunner) Run(n int, fn func() error) {
	r.wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer r.wg.Done()
			if err := fn(); err != nil {
				r.errorsMu.Lock()
				r.errors = append(r.errors, err)
				r.errorsMu.Unlock()
			}
		}()
	}
}

// Wait waits for all goroutines to complete and fails if any errors occurred.
func (r *ConcurrentTestRunner) Wait() {
	r.wg.Wait()
	if len(r.errors) > 0 {
		r.t.Errorf("concurrent test had %d errors:", len(r.errors))
		for i, err := range r.errors {
			r.t.Errorf("  error %d: %v", i+1, err)
		}
	}
}

// ============================================================================
// TIMEOUT AND CONTEXT UTILITIES
// ============================================================================

// ContextWithTestTimeout creates a context with test timeout.
func ContextWithTestTimeout(t *testing.T, timeout time.Duration) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), timeout)
}

// ShortContext returns a context with 100ms timeout for quick operations.
func ShortContext(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), 100*time.Millisecond)
}

// ============================================================================
// MOCK RESPONSE BUILDERS
// ============================================================================

// ResultDataResponse creates a standard result data response.
func ResultDataResponse(code int, description string) map[string]interface{} {
	return map[string]interface{}{
		"resultData": map[string]interface{}{
			"resultCode":        code,
			"resultDescription": description,
		},
	}
}

// SuccessResultData creates a success result data response.
func SuccessResultData() map[string]interface{} {
	return ResultDataResponse(0, "Success")
}

// ErrorResponse creates an error response.
func ErrorResponse(code, message string, details ...string) *types.Error {
	return &types.Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// PaginatedResponse creates a paginated list response.
func PaginatedResponse(data interface{}, hasMore bool, totalCount int) map[string]interface{} {
	return map[string]interface{}{
		"data":       data,
		"hasMore":    hasMore,
		"totalCount": totalCount,
	}
}

// ============================================================================
// EDGE CASE TEST HELPERS
// ============================================================================

// EdgeCaseStrings returns common edge case strings for testing.
func EdgeCaseStrings() []string {
	return []string{
		"",                              // empty
		" ",                             // whitespace
		"   ",                           // multiple spaces
		"\t",                            // tab
		"\n",                            // newline
		"a",                             // single char
		strings.Repeat("x", 1000),       // very long
		"<script>alert('xss')</script>", // XSS attempt
		"'; DROP TABLE accounts; --",    // SQL injection attempt
		"null",                          // literal null string
		"undefined",                     // literal undefined string
		"NaN",                           // literal NaN string
		"中文测试",                          // unicode
		"émoji 🚀",                       // emoji
		"path/../../../etc/passwd",      // path traversal attempt
	}
}

// EdgeCaseIDs returns common edge case IDs for testing.
func EdgeCaseIDs() []string {
	return []string{
		"",
		" ",
		"123",
		"abc",
		"acc-123-456",
		strings.Repeat("a", 256),
		"id with spaces",
		"id/with/slashes",
		"id?with=query",
		"id#with#hash",
	}
}

// EdgeCaseAmounts returns common edge case amounts for testing.
func EdgeCaseAmounts() []int64 {
	return []int64{
		0,                    // zero
		1,                    // minimum positive
		-1,                   // negative
		100,                  // R$ 1,00
		10000,                // R$ 100,00
		1000000,              // R$ 10.000,00
		9223372036854775807,  // max int64
		-9223372036854775808, // min int64
	}
}

// ============================================================================
// BENCHMARKING HELPERS
// ============================================================================

// NewBenchmarkMockServer creates a minimal mock server for benchmarks.
func NewBenchmarkMockServer(response interface{}) *httptest.Server {
	respBytes, _ := json.Marshal(response)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respBytes)
	}))
}

// NewBenchmarkClient creates a minimal client for benchmarks.
func NewBenchmarkClient(b *testing.B, serverURL string) *Client {
	client, err := New(Config{
		APIKey:    "bench-key",
		UserAgent: "bench/1.0",
		BaseURL:   serverURL,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		b.Fatalf("failed to create benchmark client: %v", err)
	}
	return client
}

// ============================================================================
// TEST DOCUMENTATION HELPERS
// ============================================================================

// FormatTestDoc returns formatted documentation for a test.
func FormatTestDoc(purpose, preconditions, steps, expected string) string {
	return fmt.Sprintf(`
Purpose: %s
Preconditions: %s
Steps: %s
Expected: %s
`, purpose, preconditions, steps, expected)
}
