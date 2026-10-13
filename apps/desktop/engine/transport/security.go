package transport

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"nodex-desktop/engine/identity"
)

var identityExtensionOID = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 57264, 1, 1}

type certificateIdentity struct {
	NodeID string `json:"node_id"`
	UUID   string `json:"uuid"`
}

func newIdentityCertificate(id identity.Public, key ed25519.PrivateKey) (tls.Certificate, error) {
	if len(key) != ed25519.PrivateKeySize {
		return tls.Certificate{}, fmt.Errorf("transport: invalid Ed25519 signing key")
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("transport: generate certificate serial: %w", err)
	}
	ext, err := json.Marshal(certificateIdentity{NodeID: id.NodeID, UUID: id.UUID})
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("transport: encode certificate identity: %w", err)
	}
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: id.NodeID},
		DNSNames:              []string{id.UUID},
		NotBefore:             now.Add(-time.Minute),
		NotAfter:              now.Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		ExtraExtensions:       []pkix.Extension{{Id: identityExtensionOID, Value: ext}},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("transport: create identity certificate: %w", err)
	}
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}, nil
}

func serverTLSConfig(cert tls.Certificate) *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		ClientAuth:   tls.RequireAnyClientCert,
		VerifyPeerCertificate: func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
			return verifyIdentityCertificate(rawCerts, nil)
		},
	}
}

func clientTLSConfig(cert tls.Certificate, expected *identity.Public) *tls.Config {
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		MinVersion:         tls.VersionTLS13,
		InsecureSkipVerify: true, // verification is performed against the discovered Ed25519 key below
		VerifyConnection: func(state tls.ConnectionState) error {
			if len(state.PeerCertificates) == 0 {
				return fmt.Errorf("transport: peer did not present an identity certificate")
			}
			return verifyParsedIdentityCertificate(state.PeerCertificates[0], expected)
		},
	}
}

func verifyIdentityCertificate(rawCerts [][]byte, expected *identity.Public) error {
	if len(rawCerts) == 0 {
		return fmt.Errorf("transport: peer did not present an identity certificate")
	}
	cert, err := x509.ParseCertificate(rawCerts[0])
	if err != nil {
		return fmt.Errorf("transport: parse peer certificate: %w", err)
	}
	return verifyParsedIdentityCertificate(cert, expected)
}

func verifyParsedIdentityCertificate(cert *x509.Certificate, expected *identity.Public) error {
	pub, ok := cert.PublicKey.(ed25519.PublicKey)
	if !ok || len(pub) != ed25519.PublicKeySize {
		return fmt.Errorf("transport: peer certificate is not Ed25519")
	}
	pool := x509.NewCertPool()
	pool.AddCert(cert)
	if _, err := cert.Verify(x509.VerifyOptions{Roots: pool, Intermediates: x509.NewCertPool(), KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}}); err != nil {
		return fmt.Errorf("transport: peer certificate is not self signed: %w", err)
	}
	var embedded certificateIdentity
	for _, ext := range cert.Extensions {
		if ext.Id.Equal(identityExtensionOID) {
			if err := json.Unmarshal(ext.Value, &embedded); err != nil {
				return fmt.Errorf("transport: invalid identity certificate extension: %w", err)
			}
			break
		}
	}
	if embedded.NodeID == "" || embedded.UUID == "" || embedded.NodeID != cert.Subject.CommonName || len(cert.DNSNames) == 0 || embedded.UUID != cert.DNSNames[0] {
		return fmt.Errorf("transport: identity certificate metadata is invalid")
	}
	if expected != nil {
		expectedKey, err := hex.DecodeString(expected.PublicKey)
		if err != nil || len(expectedKey) != ed25519.PublicKeySize {
			return fmt.Errorf("transport: invalid expected peer public key")
		}
		if !equalBytes(pub, expectedKey) {
			return fmt.Errorf("transport: peer public key does not match discovered identity")
		}
		if cert.Subject.CommonName != expected.NodeID || len(cert.DNSNames) == 0 || cert.DNSNames[0] != expected.UUID {
			return fmt.Errorf("transport: peer certificate identity does not match discovery")
		}
	}
	return nil
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := range a {
		v |= a[i] ^ b[i]
	}
	return v == 0
}
