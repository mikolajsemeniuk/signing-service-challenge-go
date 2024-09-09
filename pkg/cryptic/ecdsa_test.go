package cryptic_test

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/cryptic"
	"github.com/stretchr/testify/assert"
)

func TestECCGenerator_Generate(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewECCGenerator()

	keyPair, err := generator.Generate()
	assert.NoError(t, err, "Failed to generate ECC key pair")

	assert.NotNil(t, keyPair.Public, "Generated public key should not be nil")
	assert.NotNil(t, keyPair.Private, "Generated private key should not be nil")
}

func TestECCMarshaler_Marshal(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewECCGenerator()
	keyPair, err := generator.Generate()
	assert.NoError(t, err, "Failed to generate ECC key pair")

	marshaler := cryptic.NewECCMarshaler()
	public, private, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err, "Failed to marshal ECC key pair")

	assert.NotEmpty(t, public, "Marshalled public key should not be empty")
	assert.NotEmpty(t, private, "Marshalled private key should not be empty")
}

func TestECCMarshaler_Unmarshal(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewECCGenerator()
	keyPair, err := generator.Generate()
	assert.NoError(t, err, "Failed to generate ECC key pair")

	marshaler := cryptic.NewECCMarshaler()
	_, private, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err, "Failed to marshal ECC key pair")

	unmarshaledKeyPair, err := marshaler.Unmarshal(private)
	assert.NoError(t, err, "Failed to unmarshal ECC private key")

	assert.NotNil(t, unmarshaledKeyPair.Private, "Unmarshaled private key should not be nil")
	assert.NotNil(t, unmarshaledKeyPair.Public, "Unmarshaled public key should not be nil")
}

func TestECDSASigner_Sign(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewECCGenerator()
	keyPair, err := generator.Generate()
	assert.NoError(t, err, "Failed to generate ECC key pair")

	signer := cryptic.NewECDSASigner(keyPair.Private)

	data := []byte("test data")
	signature, err := signer.Sign(data)
	assert.NoError(t, err, "Failed to sign data")
	assert.NotEmpty(t, signature, "Signature should not be empty")

	hashed := sha256.Sum256(data)

	isValid := ecdsa.Verify(keyPair.Public, hashed[:], new(big.Int).SetBytes(signature[:len(signature)/2]), new(big.Int).SetBytes(signature[len(signature)/2:]))
	assert.True(t, isValid, "Failed to verify the signature")
}

func TestGenerateECDSAWithMarshal(t *testing.T) {
	t.Parallel()

	public, private, err := cryptic.GenerateECDSAWithMarshal()
	assert.NoError(t, err, "Failed to generate and marshal ECC keys")

	assert.NotEmpty(t, public, "Generated public key should not be empty")
	assert.NotEmpty(t, private, "Generated private key should not be empty")
}

func TestUnmarshalECDSAWithSign(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewECCGenerator()
	keyPair, err := generator.Generate()
	assert.NoError(t, err, "Failed to generate ECC key pair")

	marshaler := cryptic.NewECCMarshaler()
	_, private, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err, "Failed to marshal ECC key pair")

	data := []byte("test data")
	signature, err := cryptic.UnmarshalECDSAWithSign(data, private)
	assert.NoError(t, err, "Failed to unmarshal and sign data")
	assert.NotEmpty(t, signature, "Signature should not be empty")

	hashed := sha256.Sum256(data)
	isValid := ecdsa.Verify(keyPair.Public, hashed[:], new(big.Int).SetBytes(signature[:len(signature)/2]), new(big.Int).SetBytes(signature[len(signature)/2:]))
	assert.True(t, isValid, "Failed to verify the signature")
}
