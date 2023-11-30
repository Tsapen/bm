package bmhttp

import (
	"net/http"

	"github.com/Tsapen/bm/internal/bm"
)

func httpStatus(err error) int {
	switch err.(type) {
	case bm.ValidationError:
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError
	}
}
