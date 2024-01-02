package main

import "errors"

var (
	ErrorNotFound = errors.New("Not found")
)

func ErrorParameterMissing(param string) error {
	return errors.New("Param missing: " + param)
}
