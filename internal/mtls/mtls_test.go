package mtls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// generateTestCertificate generates a self-signed certificate for testing
func generateTestCertificate() (certPEM, keyPEM []byte, err error) {
	// Generate RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Organization"},
			CommonName:   "test.example.com",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	// Encode certificate to PEM
	certPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	// Encode private key to PEM
	keyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return certPEM, keyPEM, nil
}

// generateTestCA generates a CA certificate for testing
func generateTestCA() (caCertPEM, caKeyPEM []byte, err error) {
	// Generate RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Create CA certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test CA"},
			CommonName:   "Test CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	// Encode to PEM
	caCertPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	caKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return caCertPEM, caKeyPEM, nil
}

// TestLoadTLSConfig tests loading TLS config from files
func TestLoadTLSConfig(t *testing.T) {
	// Generate test certificates
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	// Create temp directory
	tempDir := t.TempDir()

	// Write cert and key files
	certFile := filepath.Join(tempDir, "client.crt")
	keyFile := filepath.Join(tempDir, "client.key")

	if err := os.WriteFile(certFile, certPEM, 0644); err != nil {
		t.Fatalf("Failed to write cert file: %v", err)
	}
	if err := os.WriteFile(keyFile, keyPEM, 0600); err != nil {
		t.Fatalf("Failed to write key file: %v", err)
	}

	// Test LoadTLSConfig
	tlsConfig, err := LoadTLSConfig(Config{
		CertFile: certFile,
		KeyFile:  keyFile,
	})

	if err != nil {
		t.Fatalf("LoadTLSConfig failed: %v", err)
	}

	if tlsConfig == nil {
		t.Fatal("LoadTLSConfig returned nil config")
	}

	if len(tlsConfig.Certificates) != 1 {
		t.Errorf("Certificates count = %d, want 1", len(tlsConfig.Certificates))
	}

	if tlsConfig.MinVersion != tls.VersionTLS12 {
		t.Errorf("MinVersion = %d, want %d", tlsConfig.MinVersion, tls.VersionTLS12)
	}
}

// TestLoadTLSConfig_WithCA tests loading TLS config with CA certificate
func TestLoadTLSConfig_WithCA(t *testing.T) {
	// Generate test certificates
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	caCertPEM, _, err := generateTestCA()
	if err != nil {
		t.Fatalf("Failed to generate CA certificate: %v", err)
	}

	// Create temp directory
	tempDir := t.TempDir()

	// Write files
	certFile := filepath.Join(tempDir, "client.crt")
	keyFile := filepath.Join(tempDir, "client.key")
	caFile := filepath.Join(tempDir, "ca.crt")

	_ = os.WriteFile(certFile, certPEM, 0644)
	_ = os.WriteFile(keyFile, keyPEM, 0600)
	_ = os.WriteFile(caFile, caCertPEM, 0644)

	// Test LoadTLSConfig with CA
	tlsConfig, err := LoadTLSConfig(Config{
		CertFile:   certFile,
		KeyFile:    keyFile,
		CACertFile: caFile,
	})

	if err != nil {
		t.Fatalf("LoadTLSConfig failed: %v", err)
	}

	if tlsConfig.RootCAs == nil {
		t.Error("RootCAs not set when CACertFile provided")
	}
}

// TestLoadTLSConfig_MissingCertFile tests error on missing cert file
func TestLoadTLSConfig_MissingCertFile(t *testing.T) {
	_, err := LoadTLSConfig(Config{
		CertFile: "/nonexistent/cert.pem",
		KeyFile:  "/nonexistent/key.pem",
	})

	if err == nil {
		t.Error("Expected error for missing cert file")
	}
}

// TestLoadTLSConfig_MissingKeyFile tests error on missing key file
func TestLoadTLSConfig_MissingKeyFile(t *testing.T) {
	// Generate cert only
	certPEM, _, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "client.crt")
	_ = os.WriteFile(certFile, certPEM, 0644)

	_, err = LoadTLSConfig(Config{
		CertFile: certFile,
		KeyFile:  "/nonexistent/key.pem",
	})

	if err == nil {
		t.Error("Expected error for missing key file")
	}
}

// TestLoadTLSConfig_MissingCAFile tests error on missing CA file
func TestLoadTLSConfig_MissingCAFile(t *testing.T) {
	// Generate test certificates
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "client.crt")
	keyFile := filepath.Join(tempDir, "client.key")

	_ = os.WriteFile(certFile, certPEM, 0644)
	_ = os.WriteFile(keyFile, keyPEM, 0600)

	_, err = LoadTLSConfig(Config{
		CertFile:   certFile,
		KeyFile:    keyFile,
		CACertFile: "/nonexistent/ca.pem",
	})

	if err == nil {
		t.Error("Expected error for missing CA file")
	}
}

// TestLoadTLSConfig_InvalidCertificate tests error on invalid certificate
func TestLoadTLSConfig_InvalidCertificate(t *testing.T) {
	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "invalid.crt")
	keyFile := filepath.Join(tempDir, "invalid.key")

	// Write invalid content
	_ = os.WriteFile(certFile, []byte("not a valid certificate"), 0644)
	_ = os.WriteFile(keyFile, []byte("not a valid key"), 0600)

	_, err := LoadTLSConfig(Config{
		CertFile: certFile,
		KeyFile:  keyFile,
	})

	if err == nil {
		t.Error("Expected error for invalid certificate")
	}
}

// TestLoadTLSConfig_InvalidCA tests error on invalid CA certificate
func TestLoadTLSConfig_InvalidCA(t *testing.T) {
	// Generate valid client cert
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "client.crt")
	keyFile := filepath.Join(tempDir, "client.key")
	caFile := filepath.Join(tempDir, "invalid-ca.crt")

	_ = os.WriteFile(certFile, certPEM, 0644)
	_ = os.WriteFile(keyFile, keyPEM, 0600)
	_ = os.WriteFile(caFile, []byte("not a valid CA certificate"), 0644)

	_, err = LoadTLSConfig(Config{
		CertFile:   certFile,
		KeyFile:    keyFile,
		CACertFile: caFile,
	})

	if err == nil {
		t.Error("Expected error for invalid CA certificate")
	}
}

// TestLoadTLSConfigFromBytes tests loading TLS config from byte slices
func TestLoadTLSConfigFromBytes(t *testing.T) {
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	tlsConfig, err := LoadTLSConfigFromBytes(certPEM, keyPEM, nil)
	if err != nil {
		t.Fatalf("LoadTLSConfigFromBytes failed: %v", err)
	}

	if tlsConfig == nil {
		t.Fatal("LoadTLSConfigFromBytes returned nil config")
	}

	if len(tlsConfig.Certificates) != 1 {
		t.Errorf("Certificates count = %d, want 1", len(tlsConfig.Certificates))
	}

	if tlsConfig.MinVersion != tls.VersionTLS12 {
		t.Errorf("MinVersion = %d, want %d", tlsConfig.MinVersion, tls.VersionTLS12)
	}
}

// TestLoadTLSConfigFromBytes_WithCA tests loading TLS config from bytes with CA
func TestLoadTLSConfigFromBytes_WithCA(t *testing.T) {
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	caCertPEM, _, err := generateTestCA()
	if err != nil {
		t.Fatalf("Failed to generate CA certificate: %v", err)
	}

	tlsConfig, err := LoadTLSConfigFromBytes(certPEM, keyPEM, caCertPEM)
	if err != nil {
		t.Fatalf("LoadTLSConfigFromBytes failed: %v", err)
	}

	if tlsConfig.RootCAs == nil {
		t.Error("RootCAs not set when caCertPEM provided")
	}
}

// TestLoadTLSConfigFromBytes_InvalidCert tests error on invalid certificate bytes
func TestLoadTLSConfigFromBytes_InvalidCert(t *testing.T) {
	_, err := LoadTLSConfigFromBytes([]byte("invalid"), []byte("invalid"), nil)

	if err == nil {
		t.Error("Expected error for invalid certificate bytes")
	}
}

// TestLoadTLSConfigFromBytes_InvalidCA tests error on invalid CA bytes
func TestLoadTLSConfigFromBytes_InvalidCA(t *testing.T) {
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	_, err = LoadTLSConfigFromBytes(certPEM, keyPEM, []byte("invalid CA"))

	if err == nil {
		t.Error("Expected error for invalid CA bytes")
	}
}

// TestLoadTLSConfigFromBytes_EmptyCA tests that empty CA is allowed
func TestLoadTLSConfigFromBytes_EmptyCA(t *testing.T) {
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	tlsConfig, err := LoadTLSConfigFromBytes(certPEM, keyPEM, []byte{})
	if err != nil {
		t.Fatalf("LoadTLSConfigFromBytes failed: %v", err)
	}

	if tlsConfig.RootCAs != nil {
		t.Error("RootCAs should be nil when empty CA provided")
	}
}

// TestConfig_Fields tests Config struct fields
func TestConfig_Fields(t *testing.T) {
	cfg := Config{
		CertFile:   "/path/to/cert.pem",
		KeyFile:    "/path/to/key.pem",
		CACertFile: "/path/to/ca.pem",
	}

	if cfg.CertFile != "/path/to/cert.pem" {
		t.Errorf("CertFile = %s, want /path/to/cert.pem", cfg.CertFile)
	}
	if cfg.KeyFile != "/path/to/key.pem" {
		t.Errorf("KeyFile = %s, want /path/to/key.pem", cfg.KeyFile)
	}
	if cfg.CACertFile != "/path/to/ca.pem" {
		t.Errorf("CACertFile = %s, want /path/to/ca.pem", cfg.CACertFile)
	}
}

// TestTLSConfig_MinVersion ensures TLS 1.2 minimum
func TestTLSConfig_MinVersion(t *testing.T) {
	certPEM, keyPEM, err := generateTestCertificate()
	if err != nil {
		t.Fatalf("Failed to generate test certificate: %v", err)
	}

	tlsConfig, err := LoadTLSConfigFromBytes(certPEM, keyPEM, nil)
	if err != nil {
		t.Fatalf("LoadTLSConfigFromBytes failed: %v", err)
	}

	// Verify minimum TLS version is 1.2 for security
	if tlsConfig.MinVersion < tls.VersionTLS12 {
		t.Errorf("MinVersion = %d, should be at least TLS 1.2 (%d)", tlsConfig.MinVersion, tls.VersionTLS12)
	}
}
