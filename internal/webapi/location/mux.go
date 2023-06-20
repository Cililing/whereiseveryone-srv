package location

import (
	"github.com/labstack/echo/v4"
	"whereiseveryone/internal/users"
	"whereiseveryone/internal/webapi/binder"
	"whereiseveryone/internal/webapi/jsonErr"
	"whereiseveryone/pkg/id"
	"whereiseveryone/pkg/logger"
	"whereiseveryone/pkg/timer"
)

type mux struct {
	usersAdapter    users.Adapter
	locationAdapter users.LocationAdapter
	logger          logger.Logger
	timer           timer.Timer
}

func NewMux(
	usersAdapter users.Adapter,
	locationAdapter users.LocationAdapter,
	logger logger.Logger,
	timer timer.Timer,
) *mux {
	return &mux{
		usersAdapter:    usersAdapter,
		locationAdapter: locationAdapter,
		logger:          logger,
		timer:           timer,
	}
}

func (m *mux) Route(g *echo.Group, _ echo.MiddlewareFunc) {
	g.POST("/update", m.updateLocation)
	g.POST("/fetch", m.fetchLocation)
}

// updateLocation
//
// @summary update user's location
// @description updates user's location
// @tags location
// @accept json
// @produces json
// @security Bearer
// @param locationUpdate body updateLocationRequest true "location"
// @success 204
// @failure 400 {object} jsonErr.JsonError "invalid request"
// @failure 500 {object} jsonErr.JsonError "internal server error"
// @router /location/update [POST]
func (m *mux) updateLocation(c echo.Context) error {
	request, bindErr := binder.BindRequest[updateLocationRequest](c, true)
	if bindErr != nil {
		return bindErr.Echo(c)
	}
	defer request.Cancel()

	userID := request.UserID()
	requestData := request.Request

	err := m.locationAdapter.UpdateLocation(request.Context(), userID, users.Location{
		Longitude:  requestData.Longitude,
		Latitude:   requestData.Latitude,
		LastUpdate: m.timer.Now(),
	})

	if err != nil {
		return jsonErr.EchoInternalError(err).Echo(c)
	}

	return c.NoContent(204)
}

// fetchLocation
//
// @summary returns location of provided users
// @description fetches users location
// @tags location
// @accept json
// @produces json
// @security Bearer
// @param fetchLocation body fetchRequest true "arrays of ids or nicks"
// @success 200 {object} fetchUserResponse "list of user"
// @failure 500 {object} jsonErr.JsonError "internal server error"
// @router /location/fetch [POST]
func (m *mux) fetchLocation(c echo.Context) error {
	request, bindErr := binder.BindRequest[fetchRequest](c, true)
	if bindErr != nil {
		return bindErr.Echo(c)
	}
	defer request.Cancel()

	requestData := request.Request
	idsFromNames, err := m.usersAdapter.TranslateNamesToIDs(request.Context(), requestData.Nicks)
	if err != nil {
		return jsonErr.EchoInternalError(err).Echo(c)
	}

	var idsFromRequest []id.ID
	for _, rawID := range requestData.UserIDs {
		parsedID, _ := id.FromString(rawID)
		idsFromRequest = append(idsFromRequest, parsedID)
	}

	allIds := append(idsFromNames, idsFromRequest...)
	userDetails, err := m.usersAdapter.UsersDetails(request.Context(), allIds)
	if err != nil {
		return jsonErr.EchoInternalError(err).Echo(c)
	}

	var res fetchUserResponse
	for k, v := range userDetails {
		if v.Location == nil {
			v.Location = &users.Location{}
		}
		res = append(res, userLocation{
			UUID:       k.Hex(),
			Nick:       v.Username,
			Longitude:  v.Location.Longitude,
			Latitude:   v.Location.Latitude,
			LastUpdate: v.Location.LastUpdate,
		})
	}

	return c.JSON(200, res)
}
