package crypto_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/crypto"
)

func TestRSA(t *testing.T) {
	generator := crypto.NewRSAGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}

	marshaler := crypto.NewRSAMarshaler()
	encodedPublic, encodedPrivate, err := marshaler.Marshal(*keyPair)
	if err != nil {
		log.Fatalf("Failed to encode RSA key pair: %v", err)
	}

	fmt.Println("Encoded RSA Public Key:\n", string(encodedPublic))
	fmt.Println("Encoded RSA Private Key:\n", string(encodedPrivate))

	decodedRSAKeyPair, err := marshaler.Unmarshal(encodedPrivate)
	if err != nil {
		log.Fatalf("Failed to decode RSA private key: %v", err)
	}

	signer := crypto.NewRSASigner(decodedRSAKeyPair.Private)
	data := []byte("signing example")
	signature, err := signer.Sign(data)
	if err != nil {
		log.Fatalf("Failed to sign data: %v", err)
	}

	fmt.Printf("Signature: %x\n", signature)
}
