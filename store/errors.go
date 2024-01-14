package store

import "errors"

var (
	ErrorNotFound      = errors.New("Not found")
	ErrorAlreadyExists = errors.New("Already exists")
)
