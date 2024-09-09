package signature

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Algorithm string

const (
	ECC Algorithm = "ECC"
	RSA Algorithm = "RSA"
)

func (a *Algorithm) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("failed to unmarshal Algorithm: %w", err)
	}

	if v != "ECC" && v != "RSA" {
		return ErrInvalidAlgorithm
	}

	*a = Algorithm(v)

	return nil
}

type Device struct {
	Key          uuid.UUID     `json:"key"`
	PublicKey    []byte        `json:"publicKey"`
	PrivateKey   []byte        `json:"privateKey"`
	Algorithm    Algorithm     `json:"algorithm"`
	Label        string        `json:"label"`
	Counter      int64         `json:"counter"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signedData"`
}
