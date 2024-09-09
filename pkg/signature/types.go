// Package signature provides implement an API service that allows to create `signature devices` with which they can sign arbitrary transaction data.
package signature

// types.go implements types shared across the package with it's validating rules.

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Algorithm represents the cryptographic algorithm used for the signature device (ECC or RSA).
type Algorithm string

const (
	ECC Algorithm = "ECC"
	RSA Algorithm = "RSA"
)

func (a *Algorithm) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return fmt.Errorf("error unmarshalling Algorithm: %w", err)
	}

	if name != "ECC" && name != "RSA" {
		return ErrInvalidAlgorithm
	}

	*a = Algorithm(name)

	return nil
}

// Label represents the label assigned to a signature device, with validation for maximum length (255).
type Label string

func (l *Label) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return fmt.Errorf("error unmarshalling Data: %w", err)
	}

	count := utf8.RuneCountInString(name)
	limit := 255

	if count > limit {
		return ErrLabelTooLong
	}

	*l = Label(name)

	return nil
}

// Data represents arbitrary transaction data that is signed, with validation for size limits (2-1024).
type Data string

func (d *Data) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("error unmarshalling Data: %w", err)
	}

	count := utf8.RuneCountInString(value)
	if count < 2 || count > 1024 {
		return ErrDataIncorrectSize
	}

	*d = Data(value)

	return nil
}

// Device represents a signature device, which includes cryptographic keys, algorithm, label, counter, and associated transactions.
type Device struct {
	Key          uuid.UUID     `json:"key"`
	PublicKey    []byte        `json:"publicKey"`
	PrivateKey   []byte        `json:"privateKey"`
	Algorithm    Algorithm     `json:"algorithm"`
	Label        Label         `json:"label"`
	Counter      int64         `json:"counter"`
	Transactions []Transaction `json:"transactions"`
}

// Transaction represents a transaction containing signed data and the corresponding signature.
type Transaction struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signedData"`
}
