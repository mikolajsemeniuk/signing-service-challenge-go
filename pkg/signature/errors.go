package signature

// errors.go defines common error messages used across the `signature` package.
// Errors are used for handling invalid algorithms, missing devices, existing devices, missing transactions, and empty request bodies.

import "errors"

var (
	ErrInvalidAlgorithm    = errors.New(`algorithm can be "RSA" or "ECC"`)
	ErrLabelTooLong        = errors.New("label cannot have more than 255 characters")
	ErrDataIncorrectSize   = errors.New("data characters has to be between 2 and 1024")
	ErrDeviceNotFound      = errors.New("device not found")
	ErrDeviceAlreadyExists = errors.New("device already exists")
	ErrTransactionNotFound = errors.New("transaction not found")
)
