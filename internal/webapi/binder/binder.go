package binder

import (
	"context"
	"reflect"
	"time"
	"whereiseveryone/internal/webapi/jsonErr"

	"github.com/labstack/echo/v4"
	"whereiseveryone/internal/webapi"
	"whereiseveryone/pkg/id"
	"whereiseveryone/pkg/jwt"
)

// BaseContext is interface over Context without generic type
// Allows to use Context without generic type
type BaseContext interface {
	Context() context.Context
	Cancel() context.CancelFunc
	Echo() echo.Context
	UserID() id.ID
	TokenData() jwt.SignedToken
}

type Context[T any] struct {
	ctx    context.Context
	cancel context.CancelFunc
	echo   echo.Context

	userID    id.ID
	tokenData jwt.SignedToken

	Request T
}

func (c Context[T]) Context() context.Context {
	return c.ctx
}

func (c Context[T]) Cancel() context.CancelFunc {
	return c.cancel
}

func (c Context[T]) Echo() echo.Context { // nolint:ireturn // nolintlint
	return c.echo
}

func (c Context[T]) UserID() id.ID {
	return c.userID
}

func (c Context[T]) TokenData() jwt.SignedToken {
	return c.tokenData
}

type StructValidator interface {
	Struct(str any) error
}

// BindRequest bind requests returning Context, user data (if requireAuth) and an error.
// T must be a simple type to be validated (pointers are not validated).
func BindRequest[T any](
	c echo.Context,
	requireAuth bool,
) (*Context[T], error) {
	result := &Context[T]{
		echo: c,
	}
	var t T

	// Obtain context and cancel
	reqCtx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(60)*time.Second)
	result.ctx = reqCtx
	result.cancel = cancel

	if requireAuth {
		jwtToken, err := webapi.GetJWTToken(c)
		if err != nil {
			return result, jsonErr.EchoForbiddenError(c)
		}
		requesterID, err := id.FromString(jwtToken.ID)
		if err != nil {
			return result, jsonErr.EchoInvalidRequestError(c, err)
		}
		result.userID = requesterID
		result.tokenData = jwtToken
	}

	// Obtain request
	if err := c.Bind(&t); err != nil {
		return result, jsonErr.EchoInvalidRequestError(c, err)
	}

	if val := reflect.ValueOf(t); val.Kind() == reflect.Struct { // don't validate interface{} type
		if err := c.Validate(t); err != nil {
			return result, jsonErr.EchoInvalidRequestError(c, err)
		}
	}

	result.Request = t
	return result, nil
}
