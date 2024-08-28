package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/login"
)

func LoginPublic(e *echo.Echo, h login.Handler) {
	g := e.Group("/api/v1/public/login")

	g.POST("", h.Login)
}
