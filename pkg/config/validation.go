package config

import (
	"strconv"
)

// Validator is optionally implemented by configuration types that allow this
// package to automatically validate them before they are returned.
type Validator interface {
	Validate() error
}

// ValidationError is returned by a Validator to indicate what failed
// validation and why.
type ValidationError struct {
	// Field is the path to the field in the configuration object that could
	// not be validated.
	Field string

	// Message describes why the field is invalid.
	Message string
}

func (ve *ValidationError) Error() string {
	return "ValidationError:" + ve.Field + ": " + ve.Message
}

// Wrap returns the ValidationError but prepended with a parent path.
func (ve *ValidationError) Wrap(parent string) *ValidationError {
	return &ValidationError{
		Field:   parent + "." + ve.Field,
		Message: ve.Message,
	}
}

// WrapIdx returns the ValidationError but prepended by the parent path and
// array index of the item, where there are multiple.
func (ve *ValidationError) WrapIdx(parent string, idx int) *ValidationError {
	return &ValidationError{
		Field:   parent + "[" + strconv.Itoa(idx) + "]." + ve.Field,
		Message: ve.Message,
	}
}
