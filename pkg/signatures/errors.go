package signatures

import "errors"

var (
	ErrInvalidAlgorithm = errors.New(`algorithm can be "RSA" or "ECC"`)
	ErrDeviceNotFound   = errors.New("device not found")
	ErrEmptyJSONBody    = errors.New("request body cannot be empty")
)
