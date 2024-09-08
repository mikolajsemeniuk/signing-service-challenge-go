package crypto_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/pkg/crypto"
)

func TestECDSAV2(t *testing.T) {
	// ================================
	// Initialize the ECDSA Algorithm
	// ================================
	ecdsa, err := crypto.NewRSA()
	if err != nil {
		log.Fatalf("Failed to initialize ECDSA algorithm: %v", err)
	}

	// Retrieve keys using the Keys method
	public, private, err := ecdsa.Keys()
	if err != nil {
		log.Fatalf("Failed to retrieve ECC keys: %v", err)
	}

	fmt.Println("Encoded ECDSA Public Key:\n", string(public))
	fmt.Println("Encoded ECDSA Private Key:\n", string(private))

	// Use the signer to sign data
	dataToSign := []byte("ECC signing example")
	signature, err := ecdsa.Sign(dataToSign)
	if err != nil {
		log.Fatalf("Failed to sign data with ECDSA: %v", err)
	}

	fmt.Printf("ECC Signature: %x\n", signature)
}

func TestECDSA(t *testing.T) {
	// ================================
	// Generate ECC Key Pair and Sign Data
	// ================================
	generator := crypto.NewECCGenerator()
	keyPair, err := generator.Generate()
	if err != nil {
		log.Fatalf("Failed to generate ECC key pair: %v", err)
	}

	// Create ECC marshaler and encode the keys
	eccMarshaler := crypto.NewECCMarshaler()
	public, private, err := eccMarshaler.Marshal(*keyPair)
	if err != nil {
		log.Fatalf("Failed to encode ECC key pair: %v", err)
	}

	fmt.Println("Encoded ECDSA Public Key:\n", string(public))
	fmt.Println("Encoded ECDSA Private Key:\n", string(private))

	// Decode ECC private key
	decodedECCKeyPair, err := eccMarshaler.Unmarshal(private)
	if err != nil {
		log.Fatalf("Failed to decode ECC private key: %v", err)
	}

	// Sign some data using ECC
	ecdsaSigner := crypto.NewECDSASigner(decodedECCKeyPair.Private)
	dataToSign := []byte("ECC signing example")
	eccSignature, err := ecdsaSigner.Sign(dataToSign)
	if err != nil {
		log.Fatalf("Failed to sign data with ECC: %v", err)
	}

	fmt.Printf("ECC Signature: %x\n", eccSignature)
}

func TestRSA(t *testing.T) {
	// ================================
	// Generate RSA Key Pair and Sign Data
	// ================================
	rsaGenerator := crypto.NewRSAGenerator()
	rsaKeyPair, err := rsaGenerator.Generate()
	if err != nil {
		log.Fatalf("Failed to generate RSA key pair: %v", err)
	}

	// Create RSA marshaler and encode the keys
	rsaMarshaler := crypto.NewRSAMarshaler()
	encodedPublic, encodedPrivate, err := rsaMarshaler.Marshal(*rsaKeyPair)
	if err != nil {
		log.Fatalf("Failed to encode RSA key pair: %v", err)
	}

	fmt.Println("Encoded RSA Public Key:\n", string(encodedPublic))
	fmt.Println("Encoded RSA Private Key:\n", string(encodedPrivate))

	// Decode RSA private key
	decodedRSAKeyPair, err := rsaMarshaler.Unmarshal(encodedPrivate)
	if err != nil {
		log.Fatalf("Failed to decode RSA private key: %v", err)
	}

	// Sign some data using RSA
	rsaSigner := crypto.NewRSASigner(decodedRSAKeyPair.Private)
	dataToSignRSA := []byte("RSA signing example")
	rsaSignature, err := rsaSigner.Sign(dataToSignRSA)
	if err != nil {
		log.Fatalf("Failed to sign data with RSA: %v", err)
	}

	fmt.Printf("RSA Signature: %x\n", rsaSignature)
}
