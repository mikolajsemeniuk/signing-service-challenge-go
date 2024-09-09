package cryptic_test

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/cryptic"
	"github.com/stretchr/testify/assert"
)

func TestRSAGenerator_Generate(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewRSAGenerator()

	keyPair, err := generator.Generate()
	assert.NoError(t, err, "Failed to generate RSA key pair")

	assert.NotNil(t, keyPair.Public, "Generated public key should not be nil")
	assert.NotNil(t, keyPair.Private, "Generated private key should not be nil")
}

func TestRSAMarshaler_Marshal(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewRSAGenerator()
	keyPair, err := generator.Generate()
	assert.NoError(t, err, "Failed to generate RSA key pair")

	marshaler := cryptic.NewRSAMarshaler()
	public, private, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err, "Failed to marshal RSA key pair")

	assert.NotEmpty(t, public, "Marshalled public key should not be empty")
	assert.NotEmpty(t, private, "Marshalled private key should not be empty")
}

func TestRSAMarshaler_Unmarshal(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewRSAGenerator()
	keyPair, err := generator.Generate()
	assert.NoError(t, err, "Failed to generate RSA key pair")

	marshaler := cryptic.NewRSAMarshaler()
	_, private, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err, "Failed to marshal RSA key pair")

	unmarshaledKeyPair, err := marshaler.Unmarshal(private)
	assert.NoError(t, err, "Failed to unmarshal RSA private key")

	assert.NotNil(t, unmarshaledKeyPair.Private, "Unmarshaled private key should not be nil")
	assert.NotNil(t, unmarshaledKeyPair.Public, "Unmarshaled public key should not be nil")
}

func TestRSASigner_Sign(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewRSAGenerator()
	keyPair, err := generator.Generate()
	assert.NoError(t, err, "Failed to generate RSA key pair")

	signer := cryptic.NewRSASigner(keyPair.Private)

	data := []byte("test data")
	signature, err := signer.Sign(data)
	assert.NoError(t, err, "Failed to sign data")
	assert.NotEmpty(t, signature, "Signature should not be empty")
}

func TestGenerateRSAWithMarshal(t *testing.T) {
	t.Parallel()

	public, private, err := cryptic.GenerateRSAWithMarshal()
	assert.NoError(t, err, "Failed to generate and marshal RSA keys")

	assert.NotEmpty(t, public, "Generated public key should not be empty")
	assert.NotEmpty(t, private, "Generated private key should not be empty")
}

func TestUnmarshalRSAWithSign(t *testing.T) {
	t.Parallel()

	generator := cryptic.NewRSAGenerator()
	keyPair, err := generator.Generate()
	assert.NoError(t, err, "Failed to generate RSA key pair")

	marshaler := cryptic.NewRSAMarshaler()
	_, private, err := marshaler.Marshal(*keyPair)
	assert.NoError(t, err, "Failed to marshal RSA key pair")

	data := []byte("test data")
	signature, err := cryptic.UnmarshalRSAWithSign(data, private)
	assert.NoError(t, err, "Failed to unmarshal and sign data")
	assert.NotEmpty(t, signature, "Signature should not be empty")

	hash := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(keyPair.Public, crypto.SHA256, hash[:], signature)
	assert.NoError(t, err, "Failed to verify signature")
}
