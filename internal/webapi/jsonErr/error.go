package jsonErr

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type JsonError struct {
	// Message is human friendly error message
	Message string `json:"message"`
	// Code is desired http code for this error
	Code int `json:"code"`
	// Err is a golang error returned by the app
	// It is removed in production application (TBD)
	Err error `json:"error" swaggertype:"string"`
}

func (h *JsonError) Error() string {
	if h.Err != nil {
		return fmt.Sprintf("%s: %s", h.Message, h.Err)
	}

	return fmt.Sprintf("%s", h.Message)
}

func (h *JsonError) Echo(context echo.Context) error {
	return context.JSON(h.Code, h)
}

func EchoError(code int, message string, err error) *JsonError {
	httpErr := JsonError{message, code, err}
	return &httpErr
}

func EchoInvalidRequestError(err error) *JsonError {
	return EchoError(400, "invalid request", err)
}

func EchoNotFoundError(err error) *JsonError {
	return EchoError(404, "not found", err)
}

func EchoInternalError(err error) *JsonError {
	return EchoError(500, "internal error", err)
}

func EchoForbiddenError() *JsonError {
	return EchoError(403, "forbidden", nil)
}

func EchoConflictError(err error) *JsonError {
	return EchoError(409, "conflict", err)
}
