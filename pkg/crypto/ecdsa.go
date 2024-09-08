package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

// type ECDSA struct {
// 	keyPair    *ECCKeyPair
// 	marshaller ECCMarshaler
// 	signer     ECDSASigner
// }

// func NewECDSA() (*ECDSA, error) {
// 	generator := NewECCGenerator()
// 	marshaller := NewECCMarshaler()

// 	keys, err := generator.Generate()
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, private, err := marshaller.Marshal(*keys)
// 	if err != nil {
// 		return nil, err
// 	}

// 	decodedPrivate, err := marshaller.Unmarshal(private)
// 	if err != nil {
// 		return nil, err
// 	}

// 	signer := NewECDSASigner(decodedPrivate.Private)

// 	algorithm := ECDSA{
// 		keyPair:    keys,
// 		marshaller: marshaller,
// 		signer:     *signer,
// 	}

// 	return &algorithm, nil
// }

// func (e *ECDSA) Keys() ([]byte, []byte, error) {
// 	public, private, err := e.marshaller.Marshal(*e.keyPair)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return public, private, nil
// }

// func (e *ECDSA) Sign(data []byte) ([]byte, error) {
// 	signature, err := e.signer.Sign(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return signature, nil
// }

// ECCKeyPair is a DTO that holds ECC private and public keys.
type ECCKeyPair struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

// ECCGenerator generates an ECC key pair.
type ECCGenerator struct{}

// NewECCMarshaler creates a new ECCMarshaler.
func NewECCGenerator() ECCGenerator {
	return ECCGenerator{}
}

// Generate generates a new ECCKeyPair.
func (g *ECCGenerator) Generate() (*ECCKeyPair, error) {
	// Security has been ignored for the sake of simplicity.
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Public:  &key.PublicKey,
		Private: key,
	}, nil
}

// ECCMarshaler can encode and decode an ECC key pair.
type ECCMarshaler struct{}

// NewECCMarshaler creates a new ECCMarshaler.
func NewECCMarshaler() ECCMarshaler {
	return ECCMarshaler{}
}

// TODO: fix arg here
// Encode takes an ECCKeyPair and encodes it to be written on disk.
// It returns the public and the private key as a byte slice.
func (m ECCMarshaler) Marshal(keyPair ECCKeyPair) ([]byte, []byte, error) {
	privateKeyBytes, err := x509.MarshalECPrivateKey(keyPair.Private)
	if err != nil {
		return nil, nil, err
	}

	// Pass the public key directly, without the extra pointer
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(keyPair.Public)
	if err != nil {
		return nil, nil, err
	}

	encodedPrivate := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	encodedPublic := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return encodedPublic, encodedPrivate, nil
}

// Decode assembles an ECCKeyPair from an encoded private key.
func (m ECCMarshaler) Unmarshal(privateKeyBytes []byte) (*ECCKeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}

type ECDSASigner struct {
	privateKey *ecdsa.PrivateKey
}

// NewECDSASigner creates a new ECDSASigner with the provided ECDSA private key.
func NewECDSASigner(privateKey *ecdsa.PrivateKey) *ECDSASigner {
	return &ECDSASigner{privateKey: privateKey}
}

// Sign signs the given data using ECDSA and returns the signature.
func (s *ECDSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha256.Sum256(dataToBeSigned)

	r, sigS, err := ecdsa.Sign(rand.Reader, s.privateKey, hashed[:])
	if err != nil {
		return nil, err
	}

	signature := append(r.Bytes(), sigS.Bytes()...)

	return signature, nil
}

func GenerateECDSAWithMarshal() ([]byte, []byte, error) {
	generator := NewECCGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		return nil, nil, err
	}

	marshaler := NewECCMarshaler()
	public, private, err := marshaler.Marshal(*keyPair)
	if err != nil {
		return nil, nil, err
	}

	return public, private, nil
}
