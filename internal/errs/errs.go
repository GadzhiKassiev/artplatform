package errs

import "fmt"

type Code string

const (
	CodeInternal        Code = "INTERNAL_ERROR"
	CodeUnauthenticated Code = "UNAUTHENTICATED"
	CodeInvalidInput    Code = "INVALID_INPUT"
	CodeNotFound        Code = "NOT_FOUND"
)

type Error struct {
	Code    Code
	Message string
	Details map[string]any
}

func (e *Error) Error() string {
	return e.Message
}

func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Newf(code Code, format string, args ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}

func (e *Error) WithDetail(key string, value any) *Error {
	details := make(map[string]any, len(e.Details)+1)
	for k, v := range e.Details {
		details[k] = v
	}
	details[key] = value
	return &Error{Code: e.Code, Message: e.Message, Details: details}
}

func GetCode(err error) Code {
	if e, ok := err.(*Error); ok {
		return e.Code
	}
	return CodeInternal
}

func Is(err error, code Code) bool {
	return GetCode(err) == code
}
