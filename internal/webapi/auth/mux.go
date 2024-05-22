package auth

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"time"
	"whereiseveryone/internal/users"
	"whereiseveryone/internal/webapi/jsonErr"
	"whereiseveryone/pkg/crypto"
	"whereiseveryone/pkg/id"
	"whereiseveryone/pkg/jwt"
	"whereiseveryone/pkg/timer"
)

type mux struct {
	userAdapter users.Adapter
	timer       timer.Timer
	jwt         *jwt.JWT
}

func NewMux(
	userAdapter users.Adapter,
	timer timer.Timer,
	jwt *jwt.JWT,
) *mux {
	return &mux{userAdapter, timer, jwt}
}

func (m *mux) Route(g *echo.Group, _ echo.MiddlewareFunc) {
	g.POST("/signup", m.signUp)
	g.POST("/login", m.logIn)
}

// signUp
//
// @summary sign up as a new user
// @description creates a new user
// @tags auth
// @accept json
// @produces json
// @param userDetails body signUpRequest true "sign up details"
// @success 200 {object} authResponse
// @failure 400 {object} jsonErr.JsonError "invalid request"
// @failure 409 {object} jsonErr.JsonError "conflict (user with such a name exists)
// @failure 500 {object} jsonErr.JsonError "internal server error"
// @router /auth/signup [POST]
func (m *mux) signUp(c echo.Context) error {
	reqCtx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(60)*time.Second)
	defer cancel()

	var request signUpRequest
	if err := c.Bind(&request); err != nil {
		return jsonErr.EchoInvalidRequestError(c, err)
	}
	if err := c.Validate(request); err != nil {
		return jsonErr.EchoInvalidRequestError(c, err)
	}

	encPass, err := crypto.HashPassword(request.Password)
	if err != nil {
		return jsonErr.EchoInvalidRequestError(c, err)
	}

	u := users.User{
		ID:           id.ID{}, // stub
		Username:     request.Name,
		Password:     encPass,
		Email:        request.Email,
		Token:        "",
		RefreshToken: "",
		CreatedAt:    m.timer.Now(),
		UpdatedAt:    m.timer.Now(),
		Location:     nil,
	}

	if u, err = m.userAdapter.NewUser(reqCtx, u); err != nil { // overwrite user for ID and generated data
		if errors.Is(err, users.ErrUserNameAlreadyExists) {
			return jsonErr.EchoConflictError(c, err)
		}

		return jsonErr.EchoInternalError(c, err)
	}

	token, refresh, err := m.jwt.GenerateTokens(u.Email, u.Username, u.ID)
	if err != nil {
		return jsonErr.EchoInternalError(c, err)
	}

	if err := m.userAdapter.UpdateTokens(reqCtx, u.ID, &token, &refresh); err != nil {
		return jsonErr.EchoInternalError(c, err)
	}

	return c.JSON(200, authResponse{
		ID:           u.ID.Hex(),
		Token:        token,
		RefreshToken: refresh,
	})
}

// logIn
//
// @summary log in
// @description logs in as an exiting users using login and passowrd
// @tags auth
// @accept json
// @produces json
// @param userDetails body logInRequest true "login details"
// @success 200 {object} authResponse
// @failure 400 {object} jsonErr.JsonError "invalid request"
// @failure 403 {object} jsonErr.JsonError "forbidden (invalid password)"
// @failure 404 {object} jsonErr.JsonError "user not exists"
// @failure 500 {object} jsonErr.JsonError "internal server error"
// @router /auth/login [POST]
func (m *mux) logIn(c echo.Context) error {
	reqCtx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(60)*time.Second)
	defer cancel()

	var request logInRequest
	if err := c.Bind(&request); err != nil {
		return jsonErr.EchoInvalidRequestError(c, err)
	}
	if err := c.Validate(request); err != nil {
		return jsonErr.EchoInvalidRequestError(c, err)
	}

	u, err := m.userAdapter.GetUserByName(reqCtx, request.Name)
	if err != nil {
		if errors.Is(err, users.ErrUserNotExists) {
			return jsonErr.EchoNotFoundError(c, err)
		}
		return jsonErr.EchoInternalError(c, err)
	}

	if err = crypto.VerifyPassword(u.Password, request.Password); err != nil {
		return jsonErr.EchoForbiddenError(c)
	}

	token, refresh, err := m.jwt.GenerateTokens(u.Email, u.Username, u.ID)
	if err != nil {
		return jsonErr.EchoInternalError(c, err)
	}

	if err := m.userAdapter.UpdateTokens(reqCtx, u.ID, &token, &refresh); err != nil {
		return jsonErr.EchoInternalError(c, err)
	}

	return c.JSON(200, authResponse{
		ID:           u.ID.Hex(),
		Token:        token,
		RefreshToken: refresh,
	})
}
