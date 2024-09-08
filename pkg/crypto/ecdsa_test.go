package crypto_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/crypto"
)

func TestECDSA(t *testing.T) {
	generator := crypto.NewECCGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}

	marshaler := crypto.NewECCMarshaler()
	public, private, err := marshaler.Marshal(*keyPair)
	if err != nil {
		log.Fatalf("Failed to encode ECC key pair: %v", err)
	}

	fmt.Println("Encoded Public Key:\n", string(public))
	fmt.Println("Encoded Private Key:\n", string(private))

	decodedECCKeyPair, err := marshaler.Unmarshal(private)
	if err != nil {
		log.Fatalf("Failed to decode ECC private key: %v", err)
	}

	signer := crypto.NewECDSASigner(decodedECCKeyPair.Private)
	data := []byte("signing example")
	signature, err := signer.Sign(data)
	if err != nil {
		log.Fatalf("Failed to sign data: %v", err)
	}

	fmt.Printf("Signature: %x\n", signature)
}
