package util

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"user-service/src/dto"
)

func GetTLSConfig(tlsConfigFiles dto.TLSConfig) *tls.Config {
	// Load CA certificate
	caCert, err := os.ReadFile(tlsConfigFiles.CAFile)
	if err != nil {
		panic(err)
	}

	// Load server certificate and private key
	certs, err := tls.LoadX509KeyPair(tlsConfigFiles.CertFile, tlsConfigFiles.KeyFile)
	if err != nil {
		panic(err)
	}

	// Create a certificate pool and add CA certificate to it
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create TLS configuration
	return &tls.Config{
		Certificates: []tls.Certificate{certs},
		RootCAs:      caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert, // Change as per your needs
		ClientCAs:    caCertPool,
	}
}
