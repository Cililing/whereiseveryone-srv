package me

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"whereiseveryone/internal/users"
	"whereiseveryone/internal/webapi/binder"
	"whereiseveryone/internal/webapi/jsonerr"
	"whereiseveryone/pkg/timer"
)

type mux struct {
	userAdapter users.Adapter
	timer       timer.Timer
}

func NewMux(userAdapter users.Adapter, timer timer.Timer) *mux {
	return &mux{userAdapter: userAdapter, timer: timer}
}

func (m *mux) Route(g *echo.Group, _ echo.MiddlewareFunc) {
	g.PUT("/status", m.updateStatus)
	g.GET("/friends", m.getFriends)
	g.PUT("/location", m.updateLocation)
	g.POST("/observe", m.observe)
	g.DELETE("/observe", m.unobserve)
}

// updateStatus
//
// @summary update status
// @description updates logged user status (text status)
// @tags me
// @accept json
// @produce json
// @param status body updateStatusRequest true "update status object"
// @success 204
// @failure 400 {object} jsonerr.JSONError "invalid request"
// @failure 500 {object} jsonerr.JSONError "internal server error"
// @router /me/status [PUT]
func (m *mux) updateStatus(c echo.Context) error {
	request, bindErr := binder.BindRequest[updateStatusRequest](c, true)
	if bindErr != nil {
		return bindErr.Echo(c)
	}
	defer request.Cancel()

	requestData := request.Request
	status := requestData.Status

	err := m.userAdapter.UpdateStatus(request.Context(), request.UserID(), status)
	if err != nil {
		return jsonerr.EchoInternalError(err).Echo(c)
	}

	return c.NoContent(204)
}

// getFriends
//
// @summary get friends details
// @description returns all details about observed users
// @tags me
// @produce json
// @success 200 {object} getFriendsResponse
// @failure 500 {object} jsonerr.JSONError "internal server error"
// @router /me/friends [GET]
func (m *mux) getFriends(c echo.Context) error {
	request, bindErr := binder.BindRequest[binder.EmptyBody](c, true)
	if bindErr != nil {
		return bindErr.Echo(c)
	}
	defer request.Cancel()

	user, err := m.userAdapter.GetUser(request.Context(), request.UserID())
	if err != nil {
		return jsonerr.EchoInternalError(err).Echo(c)
	}

	observedUsersIDs := user.SubscribedUsers
	observedUsers, err := m.userAdapter.GetUsers(request.Context(), observedUsersIDs)
	if err != nil {
		return jsonerr.EchoInternalError(err).Echo(c)
	}

	var result getFriendsResponse
	for _, u := range observedUsers {
		if !u.SubscribeUser(request.UserID()) {
			continue
		}

		result = append(result, friendDetails{
			Username: u.Auth.Username,
			Status:   u.Status,
			Location: locationDetails{
				Longitude:  u.Location.Longitude,
				Latitude:   u.Location.Latitude,
				Altitude:   u.Location.Altitude,
				Bearing:    u.Location.Bearing,
				Accuracy:   u.Location.Accuracy,
				LastUpdate: u.Location.LastUpdate,
			},
		})
	}

	return c.JSON(http.StatusOK, result)
}

// updateLocation
//
// @summary update location
// @description update logged user location
// @tags me
// @accept json
// @param location body updateLocationRequest true "update location object"
// @success 204
// @failure 400 {object} jsonerr.JSONError "invalid request"
// @failure 500 {object} jsonerr.JSONError "internal server error"
// @router /me/updateLocation [PUT]
func (m *mux) updateLocation(c echo.Context) error {
	request, bindErr := binder.BindRequest[updateLocationRequest](c, true)
	if bindErr != nil {
		return bindErr.Echo(c)
	}
	defer request.Cancel()

	newLoc := request.Request
	err := m.userAdapter.UpdateLocation(request.Context(), request.UserID(), users.Location{
		Longitude:  newLoc.Longitude,
		Latitude:   newLoc.Latitude,
		Altitude:   newLoc.Altitude,
		Bearing:    newLoc.Bearing,
		Accuracy:   newLoc.Accuracy,
		LastUpdate: newLoc.LastUpdate,
	})
	if err != nil {
		return jsonerr.EchoInternalError(err).Echo(c)
	}

	return c.NoContent(204)
}

// observe
//
// @summary observe the user
// @description start observing the user, the second user must observe requester too to get his details
// @tags me
// @accept json
// @param user body observeRequest true "user to observe"
// @success 204
// @failure 400 {object} jsonerr.JSONError "invalid request"
// @failure 404 {object} jsonerr.JSONError "requested user not exists"
// @failure 500 {object} jsonerr.JSONError "internal server error"
// @router /me/observe [POST]
func (m *mux) observe(c echo.Context) error {
	request, bindErr := binder.BindRequest[observeRequest](c, true)
	if bindErr != nil {
		return bindErr.Echo(c)
	}
	defer request.Cancel()

	userToObserve, err := m.userAdapter.GetUserByUsername(request.Context(), request.Request.Username)
	if err != nil {
		if errors.Is(err, users.ErrUserNotExists) {
			return jsonerr.EchoNotFoundError(err).Echo(c)
		}
		return jsonerr.EchoInternalError(err).Echo(c)
	}

	err = m.userAdapter.ObserveUser(request.Context(), request.UserID(), userToObserve.ID)
	if err != nil {
		return jsonerr.EchoInternalError(err).Echo(c)
	}

	return c.NoContent(204)
}

// unobserve
//
// @summary unobserve the user
// @description stop observing the user, if user is not observed, nothing happen
// @tags me
// @accept json
// @param user body observeRequest true "user to unobserve"
// @success 204
// @failure 400 {object} jsonerr.JSONError "invalid request"
// @failure 404 {object} jsonerr.JSONError "requested user not exists"
// @failure 500 {object} jsonerr.JSONError "internal server error"
// @router /me/observe [DELETE]
func (m *mux) unobserve(c echo.Context) error {
	request, bindErr := binder.BindRequest[observeRequest](c, true)
	if bindErr != nil {
		return bindErr.Echo(c)
	}
	defer request.Cancel()

	userToUnobserve, err := m.userAdapter.GetUserByUsername(request.Context(), request.Request.Username)
	if err != nil {
		if errors.Is(err, users.ErrUserNotExists) {
			return jsonerr.EchoNotFoundError(err).Echo(c)
		}
		return jsonerr.EchoInternalError(err).Echo(c)
	}

	err = m.userAdapter.UnobserveUser(request.Context(), request.UserID(), userToUnobserve.ID)
	if err != nil {
		return jsonerr.EchoInternalError(err).Echo(c)
	}

	return c.NoContent(204)
}
