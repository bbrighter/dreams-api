package entity

import "errors"

var (
	ErrorNotFound     error = errors.New("not found")
	ErrorBadParam     error = errors.New("bad parameter")
	ErrorMissingParam error = errors.New("missing parameter")
)
