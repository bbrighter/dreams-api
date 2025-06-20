package entity

import (
	"fmt"
)

type CustomError struct {
	Message string
	Type    ErrorType
}

type ErrorType string

const (
	ErrorTypeBadParam = "bad parameter"
	ErrorTypeNotFound = "not found"
)

func (ce CustomError) Error() string {
	return ce.Message
}

func NewCustomError(msg string, t ErrorType) *CustomError {
	return &CustomError{Message: msg, Type: t}
}

func (ce CustomError) IsNotFound() bool {
	return ce.Type == ErrorTypeNotFound
}

func (ce CustomError) IsBadParam() bool {
	return ce.Type == ErrorTypeBadParam
}

var (
	ErrorNotFound *CustomError = &CustomError{Message: "not found", Type: ErrorTypeNotFound}
	ErrorBadParam *CustomError = &CustomError{Message: "bad parameter", Type: ErrorTypeBadParam}
)

func ErrorBadParamWithReasons(reason string) *CustomError {
	msg := fmt.Sprintf("bad parameter: %s", reason)
	return NewCustomError(msg, ErrorTypeBadParam)
}
