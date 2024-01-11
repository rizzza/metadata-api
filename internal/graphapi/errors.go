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

	// ErrUniquenessConstraint is returned when an input field must be unique.
	ErrUniquenessConstraint = errors.New("must be unique")

	ErrNamespaceInUse = errors.New("namespace is in use and can't be deleted")
)

// ErrInvalidField is returned when an invalid input is provided.
type ErrInvalidField struct {
	field string
	err   error
}

// Error implements the error interface.
func (e *ErrInvalidField) Error() string {
	return fmt.Sprintf("%s: %v", e.field, e.err)
}

func NewInvalidFieldError(field string, err error) *ErrInvalidField {
	return &ErrInvalidField{field: field, err: err}
}
