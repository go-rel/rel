package errors

var UnexpectedErrorCode = 0
var ChangesetErrorCode = 1
var NotFoundErrorCode = 2
var DuplicateErrorCode = 3

type Error struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
	Code    int    `json:"json,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) UnexpectedError() bool {
	return e.Code == UnexpectedErrorCode
}

func (e Error) ChangesetError() bool {
	return e.Code == ChangesetErrorCode
}

func (e Error) NotFoundError() bool {
	return e.Code == NotFoundErrorCode
}

func (e Error) DuplicateError() bool {
	return e.Code == DuplicateErrorCode
}

func New(message string, field string, code int) Error {
	return Error{message, field, code}
}

func UnexpectedError(message string) Error {
	return Error{
		Message: message,
		Code:    UnexpectedErrorCode,
	}
}

func NotFoundError(message string) Error {
	return Error{
		Message: message,
		Code:    NotFoundErrorCode,
	}
}

func ChangesetError(message string, field string) Error {
	return Error{
		Message: message,
		Field:   field,
		Code:    ChangesetErrorCode,
	}
}

func DuplicateError(message string, field string) Error {
	return Error{
		Message: message,
		Field:   field,
		Code:    DuplicateErrorCode,
	}
}

func Wrap(err error) error {
	if err == nil {
		return nil
	} else if _, ok := err.(Error); ok {
		return err
	} else {
		return UnexpectedError(err.Error())
	}
}
