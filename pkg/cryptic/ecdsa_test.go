package cryptic_test

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/cryptic"
)

func TestECCGenerator_Generate(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewECCGenerator()

	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate ECC key pair: %v", err)
	}

	if keyPair.Public == nil || keyPair.Private == nil {
		t.Fatal("Generated ECC key pair should not be nil")
	}
}

func TestECCMarshaler_Marshal(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewECCGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate ECC key pair: %v", err)
	}

	marshaler := cryptic.NewECCMarshaler()
	public, private, err := marshaler.Marshal(*keyPair)
	if err != nil {
		t.Fatalf("Failed to marshal ECC key pair: %v", err)
	}

	if len(public) == 0 || len(private) == 0 {
		t.Fatal("Marshalled public and private keys should not be empty")
	}
}

func TestECCMarshaler_Unmarshal(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewECCGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate ECC key pair: %v", err)
	}

	marshaler := cryptic.NewECCMarshaler()
	_, private, err := marshaler.Marshal(*keyPair)
	if err != nil {
		t.Fatalf("Failed to marshal ECC key pair: %v", err)
	}

	unmarshaledKeyPair, err := marshaler.Unmarshal(private)
	if err != nil {
		t.Fatalf("Failed to unmarshal ECC private key: %v", err)
	}

	if unmarshaledKeyPair.Private == nil || unmarshaledKeyPair.Public == nil {
		t.Fatal("Unmarshaled ECC key pair should not be nil")
	}
}

func TestECDSASigner_Sign(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewECCGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate ECC key pair: %v", err)
	}

	signer := cryptic.NewECDSASigner(keyPair.Private)

	data := []byte("test data")
	signature, err := signer.Sign(data)
	if err != nil {
		t.Fatalf("Failed to sign data: %v", err)
	}

	if len(signature) == 0 {
		t.Fatal("Signature should not be empty")
	}

	hashed := sha256.Sum256(data)
	isValid := ecdsa.Verify(keyPair.Public, hashed[:], new(big.Int).SetBytes(signature[:len(signature)/2]), new(big.Int).SetBytes(signature[len(signature)/2:]))
	if !isValid {
		t.Fatal("Failed to verify the signature")
	}
}

func TestGenerateECDSAWithMarshal(t *testing.T) {
	t.Parallel()

	public, private, err := cryptic.GenerateECDSAWithMarshal()
	if err != nil {
		t.Fatalf("Failed to generate and marshal ECC keys: %v", err)
	}

	if len(public) == 0 || len(private) == 0 {
		t.Fatal("Generated public and private keys should not be empty")
	}
}

func TestUnmarshalECDSAWithSign(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewECCGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate ECC key pair: %v", err)
	}

	marshaler := cryptic.NewECCMarshaler()
	_, private, err := marshaler.Marshal(*keyPair)
	if err != nil {
		t.Fatalf("Failed to marshal ECC key pair: %v", err)
	}

	data := []byte("test data")
	signature, err := cryptic.UnmarshalECDSAWithSign(data, private)
	if err != nil {
		t.Fatalf("Failed to unmarshal and sign data: %v", err)
	}

	if len(signature) == 0 {
		t.Fatal("Signature should not be empty")
	}

	hashed := sha256.Sum256(data)
	isValid := ecdsa.Verify(keyPair.Public, hashed[:], new(big.Int).SetBytes(signature[:len(signature)/2]), new(big.Int).SetBytes(signature[len(signature)/2:]))
	if !isValid {
		t.Fatal("Failed to verify the signature")
	}
}
