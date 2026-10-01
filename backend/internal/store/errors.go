package store

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidQuantity = errors.New("invalid quantity")
)
