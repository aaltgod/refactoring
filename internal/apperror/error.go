package apperror

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrUserNotFound   = NewResponseErr(400, "Invalid request.", "user_not_found")
	ErrInternalServer = NewResponseErr(500, "", "internal_server_error")
)

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func NewResponseErr(statusCode int, statusText string, errText string) *ErrResponse {
	return &ErrResponse{
		Err:            fmt.Errorf(errText),
		HTTPStatusCode: statusCode,
		StatusText:     statusText,
		ErrorText:      errText,
	}
}

func (e *ErrResponse) Error() string {
	return e.ErrorText
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
