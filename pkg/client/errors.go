package client

import (
	"errors"
	"fmt"
)

// ErrInvaliID returned when an invalid id is provided
type ErrInvalidID struct {
	field string
	err   error
}

// Error implements the error interface.
func (e *ErrInvalidID) Error() string {
	return fmt.Sprintf("%v, field: %s", e.err, e.field)
}

var (
	// ErrUnauthorized returned when the request is not authorized
	ErrUnauthorized = errors.New("client is unauthorized")

	// ErrPermissionDenied returned when the request is not authorized
	ErrPermissionDenied = errors.New("client does not have permission to perform this action")
)
