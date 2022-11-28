package error

import (
	"github.com/go-chi/render"
	"google.golang.org/grpc/codes"
	"net/http"
)

// -----------

func ErrInvalidRequest(message string) error {
	return toError(http.StatusBadRequest, codes.InvalidArgument, message)
}
func ErrInvalidRequestWithError(err error) error {
	return toError(http.StatusBadRequest, codes.InvalidArgument, err.Error())
}
func ErrRenderInvalidRequest(err error) render.Renderer {
	return ToErrorResponse(ErrInvalidRequestWithError(err))
}

// ----------

func ErrNoPermissionRequest(message string) error {
	return toError(http.StatusForbidden, codes.PermissionDenied, message)
}
func ErrRenderNoPermissionRequest(message string) render.Renderer {
	return ToErrorResponse(ErrNoPermissionRequest(message))
}

func ErrNotFound(message string) error {
	return toError(http.StatusNotFound, codes.NotFound, message)
}
func ErrNotFoundWithError(err error) error {
	return toError(http.StatusNotFound, codes.NotFound, err.Error())
}
func ErrRenderNotFound(err error) render.Renderer {
	return ToErrorResponse(ErrNotFoundWithError(err))
}

func ErrLocked(message string) error {
	return toError(http.StatusLocked, codes.NotFound, message)
}
func ErrLockedWithError(err error) error {
	return toError(http.StatusLocked, codes.NotFound, err.Error())
}
func ErrRenderLocked(err error) render.Renderer {
	return ToErrorResponse(ErrLockedWithError(err))
}

func ErrInternal(message string) error {
	return toError(http.StatusInternalServerError, codes.Internal, message)
}
func ErrInternalWithError(err error) error {
	return domainError{
		err:            err,
		httpStatusCode: http.StatusInternalServerError,
		code:           codes.Internal,
		message:        http.StatusText(http.StatusInternalServerError),
	}
}
func ErrRenderInternal(err error) render.Renderer {
	return ToErrorResponse(ErrInternalWithError(err))
}

func ErrRenderGeneral(err error) render.Renderer {
	return ToErrorResponse(err)
}
