package mtls

import (
	"crypto/tls"
	"fmt"
	"sync"
)

// CertRotator provides hot-reload certificate rotation for mTLS connections.
// It re-reads the certificate from disk on each TLS handshake via GetClientCertificate.
// If loading a new cert fails, it falls back to the last successfully loaded cert.
type CertRotator struct {
	certFile string
	keyFile  string

	mu       sync.RWMutex
	lastCert *tls.Certificate
}

// NewCertRotator creates a new certificate rotator that reads certs from the given files.
// It performs an initial load and validation to ensure the cert files are valid.
func NewCertRotator(certFile, keyFile string) (*CertRotator, error) {
	r := &CertRotator{
		certFile: certFile,
		keyFile:  keyFile,
	}

	// Initial load — must succeed
	cert, err := r.loadAndValidate()
	if err != nil {
		return nil, fmt.Errorf("initial certificate load failed: %w", err)
	}

	r.lastCert = cert
	return r, nil
}

// GetClientCertificate implements the tls.Config.GetClientCertificate callback.
// It attempts to reload the certificate from disk on each TLS handshake.
// If the reload fails (file error or expired cert), it falls back to the last good cert.
func (r *CertRotator) GetClientCertificate(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
	cert, err := r.loadAndValidate()
	if err != nil {
		// Fall back to last good cert
		r.mu.RLock()
		last := r.lastCert
		r.mu.RUnlock()

		if last != nil {
			return last, nil
		}
		return nil, fmt.Errorf("no valid certificate available: %w", err)
	}

	r.mu.Lock()
	r.lastCert = cert
	r.mu.Unlock()

	return cert, nil
}

// Refresh forces a certificate reload from disk.
// Returns an error if the new certificate is invalid or expired.
// On error, the previously loaded certificate remains active.
func (r *CertRotator) Refresh() error {
	cert, err := r.loadAndValidate()
	if err != nil {
		return fmt.Errorf("certificate refresh failed: %w", err)
	}

	r.mu.Lock()
	r.lastCert = cert
	r.mu.Unlock()

	return nil
}

// loadAndValidate loads the certificate from disk and validates its expiration.
func (r *CertRotator) loadAndValidate() (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(r.certFile, r.keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	if err := validateCertificateExpiration(cert); err != nil {
		return nil, err
	}

	return &cert, nil
}
