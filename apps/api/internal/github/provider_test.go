package github

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func generateTestPrivateKeyPEM(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return string(pem.EncodeToMemory(block))
}

func TestNewProvider(t *testing.T) {
	t.Run("builds a provider with a valid PEM private key", func(t *testing.T) {
		privateKey := generateTestPrivateKeyPEM(t)
		p, err := NewProvider(123, 456, privateKey)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if p.Get() == nil {
			t.Fatal("expected a *github.Client, got nil")
		}
	})

	t.Run("returns an error for an invalid PEM private key", func(t *testing.T) {
		p, err := NewProvider(123, 456, "not-a-valid-pem-key")
		if err == nil {
			t.Fatalf("expected an error, got nil provider=%+v", p)
		}
	})
}

func TestNewTokenProvider(t *testing.T) {
	t.Run("builds a provider from a static token", func(t *testing.T) {
		p := NewTokenProvider("test-token")
		if p.Get() == nil {
			t.Fatal("expected a *github.Client, got nil")
		}
	})
}
