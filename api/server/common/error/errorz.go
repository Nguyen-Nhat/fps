package error

import (
	"google.golang.org/grpc/codes"
	"net/http"

	"github.com/go-chi/render"
)

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Error   codes.Code `json:"error"`
	Message string     `json:"message"`

	AppCode   int64  `json:"code,omitempty"`       // application-specific error code
	ErrorText string `json:"error_text,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,

		Error:     codes.InvalidArgument,
		Message:   "Invalid request.",
		ErrorText: err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,

		Error:     codes.Internal,
		Message:   "Error rendering response.",
		ErrorText: err.Error(),
	}
}

func ErrNoPermissionRequest(msg string) render.Renderer {
	return &ErrResponse{
		Err:            nil,
		HTTPStatusCode: 403,

		Error:     codes.PermissionDenied,
		Message:   "No permission to access",
		ErrorText: msg,
	}
}

func ErrInternal(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,

		Error:     codes.Internal,
		Message:   "Internal server error",
		ErrorText: err.Error(),
	}
}
