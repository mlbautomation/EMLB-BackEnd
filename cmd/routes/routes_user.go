package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/user"
)

func UserAdmin(e *echo.Echo, h user.Handler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/admin/users", middlewares...)

	g.GET("", h.GetAll)
}

func UserPublic(e *echo.Echo, h user.Handler) {
	g := e.Group("/api/v1/public/users")

	g.POST("", h.Create)
}
