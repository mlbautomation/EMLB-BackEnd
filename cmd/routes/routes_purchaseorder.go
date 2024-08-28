package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/mlbautomation/ProyectoEMLB/domain/ports/purchaseorder"
)

func PurchaseOrderPrivate(e *echo.Echo, h purchaseorder.Handler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/private/purchaseorders", middlewares...)

	g.POST("", h.Create)
}
