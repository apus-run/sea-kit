package tls

import (
	"testing"
)

func TestCertificate(t *testing.T) {
	certificate, err := Certificate(":8080")
	if err != nil {
		return
	}

	t.Logf("certificate: %#v", certificate.Certificate)
}

type Authentication struct {
	// TLS authentication
	TLS *Config
}

func TestTLSConfig(t *testing.T) {
	config := &Authentication{}
	if config.TLS != nil {
		tlsConfig, err := config.TLS.Config()
		if err != nil {
			t.Fatalf("error loading tls config: %v", err)
		}
		if tlsConfig != nil && tlsConfig.InsecureSkipVerify {
			t.Error("tls config should not be insecure")
		}
	}

	t.Logf("tls config: %#v", config.TLS)
}
