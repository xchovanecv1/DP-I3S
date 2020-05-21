package admin

import (
	"net/http"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
)

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText       string            `json:"status"`           // user-level status message
	AppCode          int64             `json:"code,omitempty"`   // application-specific error code
	ErrorText        string            `json:"error,omitempty"`  // application-level error message, for debugging
	ValidationErrors validation.Errors `json:"errors,omitempty"` // user level model validation errors
}

// Render sets the application-specific error code in AppCode.
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest returns status 422 Unprocessable Entity including error message.
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     http.StatusText(http.StatusUnprocessableEntity),
		ErrorText:      err.Error(),
	}
}

// ErrRender returns status 422 Unprocessable Entity rendering response error.
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

// ErrValidation returns status 422 Unprocessable Entity stating validation errors.
func ErrValidation(err error, valErr validation.Errors) render.Renderer {
	return &ErrResponse{
		Err:              err,
		HTTPStatusCode:   http.StatusUnprocessableEntity,
		StatusText:       http.StatusText(http.StatusUnprocessableEntity),
		ErrorText:        err.Error(),
		ValidationErrors: valErr,
	}
}

// ErrMediaType returns status 415 Unprocessable Entity stating validation errors.
func ErrMediaType(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnsupportedMediaType,
		StatusText:     http.StatusText(http.StatusUnsupportedMediaType),
		ErrorText:      err.Error(),
	}
}

var (
	// ErrBadRequest return status 400 Bad Request for malformed request body.
	ErrBadRequest = &ErrResponse{HTTPStatusCode: http.StatusBadRequest, StatusText: http.StatusText(http.StatusBadRequest)}

	// ErrUnauthorized returns 401 Unauthorized.
	ErrUnauthorized = &ErrResponse{HTTPStatusCode: http.StatusUnauthorized, StatusText: http.StatusText(http.StatusUnauthorized)}

	// ErrNotFound returns status 404 Not Found for invalid resource request.
	ErrNotFound = &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: http.StatusText(http.StatusNotFound)}

	// ErrUnsupportedMediaType returns status 415 Unsupported Media Type
	ErrUnsupportedMediaType = &ErrResponse{HTTPStatusCode: http.StatusUnsupportedMediaType, StatusText: http.StatusText(http.StatusUnsupportedMediaType)}

	// ErrUnprocessableEntity returns status 422 Unsupported Media Type
	ErrUnprocessableEntity = &ErrResponse{HTTPStatusCode: http.StatusUnprocessableEntity, StatusText: http.StatusText(http.StatusUnprocessableEntity)}

	// ErrInternalServerError returns status 500 Internal Server Error.
	ErrInternalServerError = &ErrResponse{HTTPStatusCode: http.StatusInternalServerError, StatusText: http.StatusText(http.StatusInternalServerError)}
)
