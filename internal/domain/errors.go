package domain

import "errors"

// Common errors
var (
	ErrNotFound = errors.New("requested item not found")
	ErrInvalidInput = errors.New("invalid input provided")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrForbidden = errors.New("forbidden access")
	ErrInternal = errors.New("internal server error")
)