package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

// type RSA struct {
// 	keyPair    *RSAKeyPair
// 	marshaller RSAMarshaler
// 	signer     RSASigner
// }

// func NewRSA() (*RSA, error) {
// 	generator := NewRSAGenerator()
// 	keyPair, err := generator.Generate()
// 	if err != nil {
// 		return nil, err
// 	}

// 	marshaller := NewRSAMarshaler()
// 	_, private, err := marshaller.Marshal(*keyPair)
// 	if err != nil {
// 		return nil, err
// 	}

// 	decodedPrivate, err := marshaller.Unmarshal(private)
// 	if err != nil {
// 		return nil, err
// 	}

// 	signer := NewRSASigner(decodedPrivate.Private)

// 	algorithm := RSA{
// 		keyPair:    keyPair,
// 		marshaller: marshaller,
// 		signer:     *signer,
// 	}

// 	return &algorithm, nil
// }

// func (r *RSA) Keys() ([]byte, []byte, error) {
// 	public, private, err := r.marshaller.Marshal(*r.keyPair)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return public, private, nil
// }

// func (r *RSA) Sign(data []byte) ([]byte, error) {
// 	signature, err := r.signer.Sign(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return signature, nil
// }

// RSAKeyPair is a DTO that holds RSA private and public keys.
type RSAKeyPair struct {
	Public  *rsa.PublicKey
	Private *rsa.PrivateKey
}

// RSAGenerator generates a RSA key pair.
type RSAGenerator struct{}

func NewRSAGenerator() RSAGenerator {
	return RSAGenerator{}
}

// Generate generates a new RSAKeyPair.
func (g *RSAGenerator) Generate() (*RSAKeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return nil, err
	}

	return &RSAKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

// RSAMarshaler can encode and decode an RSA key pair.
type RSAMarshaler struct{}

// NewRSAMarshaler creates a new RSAMarshaler.
func NewRSAMarshaler() RSAMarshaler {
	return RSAMarshaler{}
}

// TODO: fix arg here
// Marshal takes an RSAKeyPair and encodes it to be written on disk.
// It returns the public and the private key as a byte slice.
func (m *RSAMarshaler) Marshal(keyPair RSAKeyPair) ([]byte, []byte, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(keyPair.Private)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(keyPair.Public)

	private := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA_PRIVATE_KEY",
		Bytes: privateKeyBytes,
	})

	public := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA_PUBLIC_KEY",
		Bytes: publicKeyBytes,
	})

	return public, private, nil
}

// Unmarshal takes an encoded RSA private key and transforms it into a rsa.PrivateKey.
func (m *RSAMarshaler) Unmarshal(privateKeyBytes []byte) (*RSAKeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pair := &RSAKeyPair{Private: privateKey, Public: &privateKey.PublicKey}
	return pair, nil
}

// RSASigner implements the Signer interface for RSA.
type RSASigner struct {
	privateKey *rsa.PrivateKey
}

// NewRSASigner creates a new RSASigner with the provided RSA private key.
func NewRSASigner(privateKey *rsa.PrivateKey) *RSASigner {
	return &RSASigner{privateKey: privateKey}
}

// Sign signs the given data using RSA and returns the signature.
func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha256.Sum256(dataToBeSigned)

	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func GenerateRSAWithMarshal() ([]byte, []byte, error) {
	generator := NewRSAGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		return nil, nil, err
	}

	marshaler := NewRSAMarshaler()
	public, private, err := marshaler.Marshal(*keyPair)
	if err != nil {
		return nil, nil, err
	}

	return public, private, nil
}

func UnmarshalRSAWithSign(data, private []byte) ([]byte, error) {
	marshaler := NewRSAMarshaler()
	keyPair, err := marshaler.Unmarshal(private)
	if err != nil {
		return nil, err
	}

	signer := NewRSASigner(keyPair.Private)
	signature, err := signer.Sign(data)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
