package transport

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"testing"

	"nodex-desktop/engine/identity"
)

func testIdentity(t *testing.T) (identity.Public, ed25519.PrivateKey) {
	t.Helper()
	pub, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return identity.Public{NodeID: "123 456 789", UUID: "11111111-1111-4111-8111-111111111111", PublicKey: hex.EncodeToString(pub)}, key
}

func TestIdentityCertificateMatchesPinnedPeer(t *testing.T) {
	id, key := testIdentity(t)
	cert, err := newIdentityCertificate(id, key)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatal(err)
	}
	if err := verifyParsedIdentityCertificate(parsed, &id); err != nil {
		t.Fatalf("expected certificate to verify: %v", err)
	}
}

func TestIdentityCertificateRejectsWrongPinnedKey(t *testing.T) {
	id, key := testIdentity(t)
	cert, err := newIdentityCertificate(id, key)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatal(err)
	}
	other, _ := testIdentity(t)
	if err := verifyParsedIdentityCertificate(parsed, &other); err == nil {
		t.Fatal("expected mismatched identity to be rejected")
	}
}
