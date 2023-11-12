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
