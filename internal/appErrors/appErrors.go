package apperrors

import (
	"errors"
	"net"
	"net/http"

	"gorm.io/gorm"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadParam   = errors.New("bad param")
	ErrWrongState = errors.New("wrong state")
	ErrWhisper    = errors.New("problem in whisper")
)

func HttpStatus(err error) int {
	var opErr *net.OpError
	switch {
	case
		errors.Is(err, ErrNotFound),
		errors.Is(err, gorm.ErrRecordNotFound),
		errors.Is(err, gorm.ErrForeignKeyViolated):
		return http.StatusNotFound

	case
		errors.Is(err, ErrBadParam),
		errors.Is(err, ErrWrongState),
		errors.Is(err, ErrWhisper),
		errors.As(err, &opErr):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
