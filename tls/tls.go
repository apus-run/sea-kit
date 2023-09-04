package tls

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/apus-run/sea-kit/log"
)

// TLS is the configuration for TLS files
type Config struct {
	// the CA file
	CA string
	// the cert file
	Cert string
	// the key file
	Key string

	// whether to skip the TLS verification
	Insecure bool
}

// Config return a tls.Config object
func (t *Config) Config() (*tls.Config, error) {
	if len(t.CA) <= 0 {
		// the insecure is true but no ca/cert/key, then return a tls config
		if t.Insecure == true {
			log.Debug("[TLS] Insecure is true but the CA is empty, return a tls config")
			return &tls.Config{InsecureSkipVerify: true}, nil
		}
		return nil, nil
	}

	cert, err := os.ReadFile(t.CA)

	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	// only have CA file, go TLS
	if len(t.Cert) <= 0 || len(t.Key) <= 0 {
		log.Debug("[TLS] Only have CA file, go TLS")
		return &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: t.Insecure,
		}, nil
	}

	// have both CA and cert/key, go mTLS way
	log.Debug("[TLS] Have both CA and cert/key, go mTLS way")
	certificate, err := tls.LoadX509KeyPair(t.Cert, t.Key)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		RootCAs:            caCertPool,
		Certificates:       []tls.Certificate{certificate},
		InsecureSkipVerify: t.Insecure,
	}, nil
}
