package cryptic_test

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/cryptic"
)

func TestRSAGenerator_Generate(t *testing.T) {
	t.Parallel()
	generator := cryptic.NewRSAGenerator()

	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	if keyPair.Public == nil || keyPair.Private == nil {
		t.Fatal("Generated RSA key pair should not be nil")
	}
}

func TestRSAMarshaler_Marshal(t *testing.T) {
	t.Parallel()
	generator := cryptic.NewRSAGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	marshaler := cryptic.NewRSAMarshaler()
	public, private, err := marshaler.Marshal(*keyPair)
	if err != nil {
		t.Fatalf("Failed to marshal RSA key pair: %v", err)
	}

	if len(public) == 0 || len(private) == 0 {
		t.Fatal("Marshalled public and private keys should not be empty")
	}
}

func TestRSAMarshaler_Unmarshal(t *testing.T) {
	t.Parallel()
	generator := cryptic.NewRSAGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	marshaler := cryptic.NewRSAMarshaler()
	_, private, err := marshaler.Marshal(*keyPair)
	if err != nil {
		t.Fatalf("Failed to marshal RSA key pair: %v", err)
	}

	unmarshaledKeyPair, err := marshaler.Unmarshal(private)
	if err != nil {
		t.Fatalf("Failed to unmarshal RSA private key: %v", err)
	}

	if unmarshaledKeyPair.Private == nil || unmarshaledKeyPair.Public == nil {
		t.Fatal("Unmarshaled RSA key pair should not be nil")
	}
}

func TestRSASigner_Sign(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewRSAGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	signer := cryptic.NewRSASigner(keyPair.Private)

	data := []byte("test data")
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if len(signature) == 0 {
		t.Fatal("Signature should not be empty")
	}
}

func TestGenerateRSAWithMarshal(t *testing.T) {
	t.Parallel()

	public, private, err := cryptic.GenerateRSAWithMarshal()
	if err != nil {
		t.Fatalf("Failed to generate and marshal RSA keys: %v", err)
	}

	if len(public) == 0 || len(private) == 0 {
		t.Fatal("Generated public and private keys should not be empty")
	}
}

func TestUnmarshalRSAWithSign(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewRSAGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	marshaler := cryptic.NewRSAMarshaler()
	_, private, err := marshaler.Marshal(*keyPair)
	if err != nil {
		t.Fatalf("Failed to marshal RSA key pair: %v", err)
	}

	data := []byte("test data")
	signature, err := cryptic.UnmarshalRSAWithSign(data, private)
	if err != nil {
		t.Fatalf("Failed to unmarshal and sign data: %v", err)
	}

	if len(signature) == 0 {
		t.Fatal("Signature should not be empty")
	}

	hash := sha256.Sum256(data)
	if err := rsa.VerifyPKCS1v15(keyPair.Public, crypto.SHA256, hash[:], signature); err != nil {
		t.Fatalf("Failed to verify signature: %v", err)
	}
}
