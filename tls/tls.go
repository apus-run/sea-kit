package tls

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// The defaults should be a safe configuration
const defaultMinTLSVersion = tls.VersionTLS12

// Uses the default MaxVersion from "crypto/tls"
const defaultMaxTLSVersion = 0

// Config TLS is the configuration for TLS files
type Config struct {
	// Enable TLS
	Enable bool

	// the CA file
	CA string
	// the cert file
	Cert string
	// the key file
	Key string

	// whether to skip the TLS verification
	Insecure bool

	// MinVersion sets the minimum TLS version that is acceptable.
	// If not set, TLS 1.2 will be used. (optional)
	MinVersion string
	// MaxVersion sets the maximum TLS version that is acceptable.
	// If not set, refer to crypto/tls for defaults. (optional)
	MaxVersion string
}

// Config return a tls.Config object
func (t *Config) Config() (*tls.Config, error) {
	if t.Enable == false {
		slog.Debug("[TLS] TLS is disabled")
		return nil, nil
	}
	if len(t.CA) <= 0 {
		// the insecure is true but no ca/cert/key, then return a tls config
		if t.Insecure == true {
			slog.Debug("[TLS] Insecure is true but the CA is empty, return a tls config")
			return &tls.Config{InsecureSkipVerify: true}, nil
		}
		return nil, nil
	}

	caCertPool, err := t.loadCert(t.CA)

	// only have CA file, go TLS
	if len(t.Cert) <= 0 || len(t.Key) <= 0 {
		slog.Debug("[TLS] Only have CA file, go TLS")
		return &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: t.Insecure,
		}, nil
	}

	// have both CA and cert/key, go mTLS way
	slog.Debug("[TLS] Have both CA and cert/key, go mTLS way")
	certificate, err := tls.LoadX509KeyPair(t.Cert, t.Key)
	if err != nil {
		return nil, fmt.Errorf("could not load TLS client key/certificate from %s:%s: %s", t.Key, t.Cert, err)
	}

	minVersion, err := convertVersion(t.MinVersion, defaultMinTLSVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid TLS min_version: %w", err)
	}
	maxVersion, err := convertVersion(t.MaxVersion, defaultMaxTLSVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid TLS max_version: %w", err)
	}
	return &tls.Config{
		RootCAs:            caCertPool,
		Certificates:       []tls.Certificate{certificate},
		InsecureSkipVerify: t.Insecure,
		MinVersion:         minVersion,
		MaxVersion:         maxVersion,
	}, nil
}

func (t *Config) loadCert(caPath string) (*x509.CertPool, error) {
	caPEM, err := os.ReadFile(filepath.Clean(caPath))
	if err != nil {
		return nil, fmt.Errorf("failed to load CA %s: %w", caPath, err)
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("failed to parse CA %s", caPath)
	}
	return certPool, nil
}

func convertVersion(version string, defaultVersion uint16) (uint16, error) {
	if version == "" {
		return defaultVersion, nil
	}
	val, ok := TlsVersion[version]
	if !ok {
		return 0, fmt.Errorf("unsupported TLS version: %q", version)
	}
	return val, nil
}

func Certificate(host ...string) (tls.Certificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 365)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range host {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	// create public key
	certOut := bytes.NewBuffer(nil)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// create private key
	keyOut := bytes.NewBuffer(nil)
	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return tls.Certificate{}, err
	}
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})

	return tls.X509KeyPair(certOut.Bytes(), keyOut.Bytes())
}
