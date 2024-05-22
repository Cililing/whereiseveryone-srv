package location

import (
	"github.com/labstack/echo/v4"
	"whereiseveryone/internal/users"
	"whereiseveryone/pkg/logger"
)

type mux struct {
	locationAdapter users.LocationAdapter
	logger          logger.Logger
}

func NewMux(
	adapter users.LocationAdapter,
	logger logger.Logger,
) *mux {
	return &mux{
		locationAdapter: adapter,
		logger:          logger,
	}
}

func (m *mux) Route(g *echo.Group, _ echo.MiddlewareFunc) {
	g.POST("/update", m.updateLocation)
	g.POST("/fetch", m.fetchLocation)
}

// updateLocation
// @Summary updates user location
// @Tags Location
// @Router /update [POST]
func (m *mux) updateLocation(c echo.Context) error {
	panic("not implemented")
	return nil
}

func (m *mux) fetchLocation(c echo.Context) error {
	panic("not implemented")
	return nil
}
