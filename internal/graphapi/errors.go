package graphapi

import (
	"errors"
	"fmt"
)

var (
	// ErrInternalServerError is returned when an internal error occurs.
	ErrInternalServerError = errors.New("internal server error")

	// ErrFieldEmpty is returned when a required field is empty.
	ErrFieldEmpty = errors.New("must not be empty")

	// ErrInvalidJSON is returned when invalid json is provided.
	ErrInvalidJSON = errors.New("invalid json data")
)

// ErrInvalidField is returned when an invalid ID is provided.
type ErrInvalidField struct {
	field string
	err   error
}

// Error implements the error interface.
func (e *ErrInvalidField) Error() string {
	return fmt.Sprintf("%v, field: %s", e.err, e.field)
}
