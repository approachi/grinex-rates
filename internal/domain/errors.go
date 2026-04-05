package domain

import "errors"

// calculator errors
var (
	ErrEmptyEntries = errors.New("entries is empty")
	ErrOutOfRange   = errors.New("position out of range")
	ErrInvalidRange = errors.New("position_from must be less than or equal to position_to")
)

// service errors
var (
	ErrUnknownMethod = errors.New("unknown calculation method")
)
