// Package errors wraps driver and changeset error as a grimoire's error.
package errors

// UnexpectedErrorCode defines default code for Unexpected Errors.
var UnexpectedErrorCode = 0

// ChangesetErrorCode defines default code for Changeset Errors.
var ChangesetErrorCode = 1

// NotFoundErrorCode defines default code for NotFound Errors.
var NotFoundErrorCode = 2

// DuplicateErrorCode defines default code for Duplicate Errors.
var DuplicateErrorCode = 3

// Error defines information about grimoire's error.
type Error struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// Error prints error message.
func (e Error) Error() string {
	return e.Message
}

// UnexpectedError returns true if error is an UnexpectedError.
func (e Error) UnexpectedError() bool {
	return e.Code == UnexpectedErrorCode
}

// ChangesetError returns true if error is an ChangesetError.
func (e Error) ChangesetError() bool {
	return e.Code == ChangesetErrorCode
}

// NotFoundError returns true if error is an NotFoundError.
func (e Error) NotFoundError() bool {
	return e.Code == NotFoundErrorCode
}

// DuplicateError returns true if error is an DuplicateError.
func (e Error) DuplicateError() bool {
	return e.Code == DuplicateErrorCode
}

// New creates an error with custom image, field and error code.
func New(message string, field string, code int) Error {
	return Error{message, field, code}
}

// UnexpectedError creates an unexpected error with custom message.
func UnexpectedError(message string) Error {
	return Error{
		Message: message,
		Code:    UnexpectedErrorCode,
	}
}

// NotFoundError creates a not found error with custom message.
func NotFoundError(message string) Error {
	return Error{
		Message: message,
		Code:    NotFoundErrorCode,
	}
}

// ChangesetError creates a changeset error with custom message.
func ChangesetError(message string, field string) Error {
	return Error{
		Message: message,
		Field:   field,
		Code:    ChangesetErrorCode,
	}
}

// DuplicateError creates a duplicate error with custom message.
func DuplicateError(message string, field string) Error {
	return Error{
		Message: message,
		Field:   field,
		Code:    DuplicateErrorCode,
	}
}

// Wrap errors as grimoire's error.
// If error is grimoire error, it'll remain as is.
// Otherwise it'll be wrapped as unexpected error.
func Wrap(err error) error {
	if err == nil {
		return nil
	} else if _, ok := err.(Error); ok {
		return err
	} else {
		return UnexpectedError(err.Error())
	}
}
