package httpchi

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

const (
	ServerErrorStatusText = "server error"
)

var (
	UserNotFound = errors.New("user_not_found")
)

type CreateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

type CreateUserResponse struct {
	ID string `json:"user_id"`
}

func (c *CreateUserRequest) Bind(r *http.Request) error { return nil }

type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
}

func (c *UpdateUserRequest) Bind(r *http.Request) error { return nil }

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrorResponse(err error, status int, statusText string) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: status,
		ErrorText:      err.Error(),
		StatusText:     statusText,
	}
}
