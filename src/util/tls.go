package util

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"user-service/src/dto"
)

// GetTLSConfig creates a TLS config from files and properties supplied by dto.TLSConfig
func GetTLSConfig(tlsOptions dto.TLSConfig) (*tls.Config, error) {
	// Load CA certificate
	caCert, err := os.ReadFile(tlsOptions.CAFile)
	if err != nil {
		return nil, err
	}

	// Load server certificate and private key
	certs, err := tls.LoadX509KeyPair(tlsOptions.CertFile, tlsOptions.KeyFile)
	if err != nil {
		return nil, err
	}

	// Create a certificate pool and add CA certificate to it
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create TLS configuration
	return &tls.Config{
		Certificates:       []tls.Certificate{certs},
		RootCAs:            caCertPool,
		ClientAuth:         tls.RequireAndVerifyClientCert, // Change as per your needs
		ClientCAs:          caCertPool,
		InsecureSkipVerify: tlsOptions.SkipVerify,
	}, nil
}
