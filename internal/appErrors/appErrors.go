package apperrors

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("not found")
	ErrBadParam = errors.New("bad param")
)

func HttpStatus(err error) int {
	switch {
	case
		errors.Is(err, ErrNotFound),
		errors.Is(err, gorm.ErrRecordNotFound),
		errors.Is(err, gorm.ErrForeignKeyViolated):
		return http.StatusNotFound

	case errors.Is(err, ErrBadParam):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
