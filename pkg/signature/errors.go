package signature

// errors.go defines common error messages used across the `signature` package.
// Errors are used for handling invalid algorithms, missing devices, existing devices, missing transactions, and empty request bodies.

import "errors"

var (
	ErrInvalidAlgorithm    = errors.New(`algorithm can be "RSA" or "ECC"`)
	ErrDeviceNotFound      = errors.New("device not found")
	ErrDeviceAlreadyExists = errors.New("device already exists")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrEmptyJSONBody       = errors.New("request body cannot be empty")
)
