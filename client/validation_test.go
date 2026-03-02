package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// validatableBody implements Validatable for testing
type validatableBody struct {
	Name string `json:"name"`
	fail bool
}

func (v *validatableBody) Validate() error {
	if v.fail {
		return fmt.Errorf("validation failed: name is required")
	}
	return nil
}

func TestValidatable_PassesValidation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "123"})
	}))
	defer server.Close()

	client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})

	body := &validatableBody{Name: "test"}
	var result map[string]any
	err := client.request(context.Background(), http.MethodPost, "/test", body, &result)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidatable_FailsValidation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("request should not have been sent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})

	body := &validatableBody{fail: true}
	var result map[string]any
	err := client.request(context.Background(), http.MethodPost, "/test", body, &result)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if err.Error() != "validation failed: name is required" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidatable_DisabledWithNoValidation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "123"})
	}))
	defer server.Close()

	// Create client with noValidation = true
	client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})
	client.noValidation = true

	// Even though validation would fail, it should be skipped
	body := &validatableBody{fail: true}
	var result map[string]any
	err := client.request(context.Background(), http.MethodPost, "/test", body, &result)
	if err != nil {
		t.Fatalf("expected no error with validation disabled, got: %v", err)
	}
}

func TestNonValidatable_SkipsValidation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "123"})
	}))
	defer server.Close()

	client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})

	// Regular struct without Validate() should work fine
	body := map[string]string{"name": "test"}
	var result map[string]any
	err := client.request(context.Background(), http.MethodPost, "/test", body, &result)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidatable_NilBody(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{})
	}))
	defer server.Close()

	client := NewTestClient(t, &MockServer{Server: server, t: t, config: &MockServerConfig{}})

	// nil body should not panic
	var result map[string]any
	err := client.request(context.Background(), http.MethodGet, "/test", nil, &result)
	if err != nil {
		t.Fatalf("expected no error for nil body, got: %v", err)
	}
}
