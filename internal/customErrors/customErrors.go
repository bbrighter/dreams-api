package customerrors

import (
	"errors"

	"gorm.io/gorm"
)

func ErrorParameterMissing(param string) error {
	return errors.New("Param missing: " + param)
}

var (
	ErrorNotFound = gorm.ErrRecordNotFound
)
