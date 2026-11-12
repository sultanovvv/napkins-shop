// Package customerr описывает прикладные ошибки независимо от транспорта.
package customerr

type Code string

const (
	Internal    Code = "INTERNAL_ERROR"
	Validation  Code = "VALIDATION_ERROR"
	NotFound    Code = "NOT_FOUND"
	Forbidden   Code = "FORBIDDEN"
	Conflict    Code = "CONFLICT"
	Unavailable Code = "SERVICE_UNAVAILABLE"
)

type Detail struct {
	Code    string
	Message string
	Target  string
}

// Message и Details предназначены для клиента. Исходная ошибка доступна
// через Unwrap для логирования и errors.Is/errors.As, но не для ответа API.
type Error struct {
	Code    Code
	Message string
	Details []Detail
	cause   error
}

func (e *Error) Error() string {
	message := string(e.Code) + ": " + e.Message
	if e.cause != nil {
		message += ": " + e.cause.Error()
	}
	return message
}

func (e *Error) Unwrap() error { return e.cause }

func NewInternal(message string, err error) *Error {
	return &Error{Code: Internal, Message: message, cause: err}
}

func NewValidation(message string, err error) *Error {
	return &Error{Code: Validation, Message: message, cause: err}
}

func NewNotFound(message string, err error) *Error {
	return &Error{Code: NotFound, Message: message, cause: err}
}

func NewForbidden(message string, err error) *Error {
	return &Error{Code: Forbidden, Message: message, cause: err}
}

func NewConflict(message string, err error) *Error {
	return &Error{Code: Conflict, Message: message, cause: err}
}

func NewUnavailable(message string, err error) *Error {
	return &Error{Code: Unavailable, Message: message, cause: err}
}
