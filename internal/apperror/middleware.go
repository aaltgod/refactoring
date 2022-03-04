package apperror

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

type errHandler func(http.ResponseWriter, *http.Request) error

func Middleware(h errHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var errResponse *ErrResponse

		if err := h(w, r); err != nil {
			if errors.As(err, &errResponse) {
				err := err.(*ErrResponse)
				render.Render(w, r, err)
				return
			}
		}
	}
}
