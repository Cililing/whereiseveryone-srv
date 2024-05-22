package me

import (
	"github.com/labstack/echo/v4"
	"whereiseveryone/internal/users"
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

}
