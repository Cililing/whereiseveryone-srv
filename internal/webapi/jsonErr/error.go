package jsonErr

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type JsonError struct {
	// Message is human friendly error message
	Message string `json:"message"`
	// Err is a golang error returned by the app
	// It is removed in production application (TBD)
	Err error `json:"error" swaggertype:"string"`
}

func (h JsonError) Error() string {
	if h.Err != nil {
		return fmt.Sprintf("%s: %s", h.Message, h.Err)
	}

	return fmt.Sprintf("%s", h.Message)
}

func EchoError(c echo.Context, code int, message string, err error) error {
	httpErr := JsonError{message, err}
	return c.JSON(code, httpErr)
}

func EchoInvalidRequestError(c echo.Context, err error) error {
	return EchoError(c, 400, "invalid request", err)
}

func EchoNotFoundError(c echo.Context, err error) error {
	return EchoError(c, 404, "not found", err)
}

func EchoInternalError(c echo.Context, err error) error {
	return EchoError(c, 500, "internal error", err)
}

func EchoForbiddenError(c echo.Context) error {
	return EchoError(c, 403, "forbidden", nil)
}

func EchoConflictError(c echo.Context, err error) error {
	return EchoError(c, 409, "conflict", err)
}
