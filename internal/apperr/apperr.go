package apperr

import (
	"errors"
	"fmt"
	"net/http"
)

type Code string

const (
	CodeValidation Code = "VALIDATION_ERROR"
	CodeNotFound   Code = "NOT_FOUND"
	CodeConflict   Code = "CONFLICT"
	CodeInternal   Code = "INTERNAL"
)

type AppError struct {
	Code       Code
	Message    string
	HTTPStatus int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
}

func (e *AppError) Unwrap() error { return e.Err }

func Validation(msg string) *AppError {
	return &AppError{Code: CodeValidation, Message: msg, HTTPStatus: http.StatusBadRequest}
}

func NotFound(msg string) *AppError {
	return &AppError{Code: CodeNotFound, Message: msg, HTTPStatus: http.StatusNotFound}
}

func Conflict(msg string, err error) *AppError {
	return &AppError{Code: CodeConflict, Message: msg, HTTPStatus: http.StatusConflict, Err: err}
}

func Internal(msg string, err error) *AppError {
	return &AppError{Code: CodeInternal, Message: msg, HTTPStatus: http.StatusInternalServerError, Err: err}
}

func As(err error) (*AppError, bool) {
	var ae *AppError
	ok := errors.As(err, &ae)
	return ae, ok
}
