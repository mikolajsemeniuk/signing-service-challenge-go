package signatures

import (
	"encoding/json"
	"time"

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
		return err
	}

	if v != "ECC" && v != "RSA" {
		return ErrInvalidAlgorithm
	}

	*a = Algorithm(v)
	return nil
}

type Device struct {
	Key          uuid.UUID
	PublicKey    []byte
	PrivateKey   []byte
	Algorithm    Algorithm
	Label        string
	Counter      int64
	Transactions []Transaction
}

type Transaction struct {
	Signature  string    `json:"signature"`
	SignedData string    `json:"signedData"`
	Created    time.Time `json:"created"`
}
