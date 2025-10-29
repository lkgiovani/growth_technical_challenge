package growth_technical_challenge

import (
	"errors"
	"fmt"
)

const (
	EBADREQUEST     = "bad_request"
	ECONFLICT       = "conflict"
	EDUPLICATION    = "duplication"
	EEXPIRED        = "expired"
	EFORBIDDEN      = "forbidden"
	EINTERNAL       = "internal"
	EINVALID        = "invalid"
	ENOTFOUND       = "not_found"
	ENOTIMPLEMENTED = "not_implemented"
	EPRECONDITION   = "precondition_failed"
	ERATELIMIT      = "rate_limit"
	ETIMEOUT        = "timeout"
	EUNAUTHORIZED   = "unauthorized"
	EUNAVAILABLE    = "service_unavailable"
)

type Error struct {
	Code string

	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("wtf error: code=%s message=%s", e.Code, e.Message)
}

func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return EINTERNAL
}

func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "Internal error."
}

func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
