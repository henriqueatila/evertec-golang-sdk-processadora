// Package mtls provides mTLS configuration for Evertec API connections.
package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"
)

// Config holds mTLS configuration.
type Config struct {
	CertFile   string // Path to client certificate PEM file
	KeyFile    string // Path to client key PEM file
	CACertFile string // Path to CA certificate PEM file (optional)
}

// validateCertificateExpiration checks if the certificate has expired.
func validateCertificateExpiration(cert tls.Certificate) error {
	if len(cert.Certificate) == 0 {
		return fmt.Errorf("client certificate has no certificates")
	}

	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fmt.Errorf("failed to parse client certificate: %w", err)
	}

	if time.Now().After(x509Cert.NotAfter) {
		return fmt.Errorf("client certificate expired on %s", x509Cert.NotAfter.Format(time.RFC3339))
	}

	return nil
}

// LoadTLSConfig creates a TLS configuration from the provided mTLS config.
func LoadTLSConfig(cfg Config) (*tls.Config, error) {
	// Load client certificate
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %w", err)
	}

	// Validate certificate expiration
	if err := validateCertificateExpiration(cert); err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Load CA certificate if provided
	if cfg.CACertFile != "" {
		caCert, err := os.ReadFile(cfg.CACertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}

		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}

// LoadTLSConfigFromBytes creates a TLS configuration from certificate bytes.
func LoadTLSConfigFromBytes(certPEM, keyPEM, caCertPEM []byte) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse client certificate: %w", err)
	}

	// Validate certificate expiration
	if err := validateCertificateExpiration(cert); err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	if len(caCertPEM) > 0 {
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCertPEM) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}
