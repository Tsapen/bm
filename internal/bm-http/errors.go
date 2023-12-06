package bmhttp

import (
	"errors"
	"net/http"

	"github.com/Tsapen/bm/internal/bm"
)

func httpStatus(err error) int {
	switch {
	case errors.As(err, &bm.ValidationError{}):
		return http.StatusBadRequest

	case errors.As(err, &bm.NotFoundError{}):
		return http.StatusNotFound

	case errors.As(err, &bm.ConflictError{}):
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}
