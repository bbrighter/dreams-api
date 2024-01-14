package controller

import "errors"

func ErrorParameterMissing(param string) error {
	return errors.New("Param missing: " + param)
}
