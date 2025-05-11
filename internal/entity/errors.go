package entity

import "errors"

var (
	ErrorNotFound error = errors.New("not found")
	ErrorBadParam error = errors.New("bad parameter")
)
