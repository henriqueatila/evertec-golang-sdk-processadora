package mtls

import (
	"crypto/tls"
	"os"
	"path/filepath"
	"testing"
)

func TestNewCertRotator(t *testing.T) {
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate cert: %v", err)
	}

	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "client.crt")
	keyFile := filepath.Join(tempDir, "client.key")

	_ = os.WriteFile(certFile, certPEM, 0644)
	_ = os.WriteFile(keyFile, keyPEM, 0600)

	rotator, err := NewCertRotator(certFile, keyFile)
	if err != nil {
		t.Fatalf("NewCertRotator failed: %v", err)
	}
	if rotator == nil {
		t.Fatal("rotator is nil")
	}
}

func TestNewCertRotator_InvalidFiles(t *testing.T) {
	_, err := NewCertRotator("/nonexistent/cert.pem", "/nonexistent/key.pem")
	if err == nil {
		t.Error("expected error for missing files")
	}
}

func TestCertRotator_GetClientCertificate(t *testing.T) {
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate cert: %v", err)
	}

	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "client.crt")
	keyFile := filepath.Join(tempDir, "client.key")

	_ = os.WriteFile(certFile, certPEM, 0644)
	_ = os.WriteFile(keyFile, keyPEM, 0600)

	rotator, err := NewCertRotator(certFile, keyFile)
	if err != nil {
		t.Fatalf("NewCertRotator failed: %v", err)
	}

	cert, err := rotator.GetClientCertificate(&tls.CertificateRequestInfo{})
	if err != nil {
		t.Fatalf("GetClientCertificate failed: %v", err)
	}
	if cert == nil {
		t.Fatal("cert is nil")
	}
}

func TestCertRotator_FallbackOnError(t *testing.T) {
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate cert: %v", err)
	}

	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "client.crt")
	keyFile := filepath.Join(tempDir, "client.key")

	_ = os.WriteFile(certFile, certPEM, 0644)
	_ = os.WriteFile(keyFile, keyPEM, 0600)

	rotator, err := NewCertRotator(certFile, keyFile)
	if err != nil {
		t.Fatalf("NewCertRotator failed: %v", err)
	}

	// Delete cert file to simulate rotation failure
	_ = os.Remove(certFile)

	// Should fall back to last good cert
	cert, err := rotator.GetClientCertificate(&tls.CertificateRequestInfo{})
	if err != nil {
		t.Fatalf("GetClientCertificate should fall back, got error: %v", err)
	}
	if cert == nil {
		t.Fatal("cert should not be nil (fallback)")
	}
}

func TestCertRotator_Refresh(t *testing.T) {
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate cert: %v", err)
	}

	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "client.crt")
	keyFile := filepath.Join(tempDir, "client.key")

	_ = os.WriteFile(certFile, certPEM, 0644)
	_ = os.WriteFile(keyFile, keyPEM, 0600)

	rotator, err := NewCertRotator(certFile, keyFile)
	if err != nil {
		t.Fatalf("NewCertRotator failed: %v", err)
	}

	// Refresh should succeed
	if err := rotator.Refresh(); err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}

	// Refresh after deleting should fail
	_ = os.Remove(certFile)
	if err := rotator.Refresh(); err == nil {
		t.Error("Refresh should fail after deleting cert")
	}
}
