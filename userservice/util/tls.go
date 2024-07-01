package util

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"path/filepath"
	"user-service/dto"

	log "github.com/sirupsen/logrus"
)

// GetRefreshableTLSConfig returns a TLS configuration that can be refreshed when the files change
// GetRefreshableTLSConfig returns an error if the files cannot be read or the TLS configuration cannot be created
func GetRefreshableTLSConfig(ctx context.Context, tlsOptions dto.TLSConfig, onRefreshError func(error), onRefresh func()) (*tls.Config, error) {
	// Create TLS configuration
	tlsConfig := &tls.Config{}
	err := populateTLSConfig(tlsConfig, tlsOptions.CAFile, tlsOptions.CertFile, tlsOptions.KeyFile, tlsOptions.SkipVerify)
	if err != nil {
		return nil, err
	}
	fileDirs := make([]string, 3)
	// get directories each file is in
	fileDirs[0] = filepath.Dir(tlsOptions.CAFile)
	fileDirs[1] = filepath.Dir(tlsOptions.CertFile)
	fileDirs[2] = filepath.Dir(tlsOptions.KeyFile)
	err = WatchPath(ctx, func(file string) {
		log.Infof("TLS files changed %s, refreshing TLS configuration", file)
		err := populateTLSConfig(tlsConfig, tlsOptions.CAFile, tlsOptions.CertFile, tlsOptions.KeyFile, tlsOptions.SkipVerify)
		if err != nil {
			onRefreshError(err)
			return
		}
		onRefresh()
	}, fileDirs...)
	if err != nil {
		return nil, err
	}
	return tlsConfig, nil
}

func populateTLSConfig(tlsConfig *tls.Config, CAFile, CertFile, KeyFile string, skipVerify bool) error {
	// Load CA certificate
	caCert, err := os.ReadFile(CAFile)
	if err != nil {
		return errors.Join(errors.New("failed to read CA file"), err)
	}

	// Load server certificate and private key
	certs, err := tls.LoadX509KeyPair(CertFile, KeyFile)
	if err != nil {
		return errors.Join(errors.New("failed to load server certificate and private key"), err)
	}

	// Create a certificate pool and add CA certificate to it
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig.Certificates = []tls.Certificate{certs}
	tlsConfig.RootCAs = caCertPool
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	tlsConfig.ClientCAs = caCertPool
	tlsConfig.InsecureSkipVerify = skipVerify
	return nil
}
