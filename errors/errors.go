// Package errors wraps driver and changeset error as a grimoire's error.
package errors

type Kind int

const (
	Unexpected Kind = iota
	Changeset
	NotFound
	UniqueConstraint
	ForeignKeyConstraint
	CheckConstraint
)

// Error defines information about grimoire's error.
type Error struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Code    int    `json:"code,omitempty"`
	kind    Kind
}

// Error prints error message.
func (e Error) Error() string {
	return e.Message
}

// Kind of error.
func (e Error) Kind() Kind {
	return e.kind
}

// New creates an error.
func New(message string, field string, kind Kind) Error {
	return Error{message, field, 0, kind}
}

// NewWithCode creates an error with code.
func NewWithCode(message string, field string, code int, kind Kind) Error {
	return Error{message, field, code, kind}
}

// NewUnexpected creates an error.
func NewUnexpected(message string) Error {
	return Error{message, "", 0, Unexpected}
}
