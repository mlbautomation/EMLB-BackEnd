package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/invoice"
)

// adminRoutes handle the routes that requires a token and permissions to certain users
func InvoiceAdmin(e *echo.Echo, h invoice.Handler, middlewares ...echo.MiddlewareFunc) {
	route := e.Group("/api/v1/admin/invoices", middlewares...)

	route.GET("", h.GetAll)
}

// privateRoutes handle the routes that requires a token
func InvoicePrivate(e *echo.Echo, h invoice.Handler, middlewares ...echo.MiddlewareFunc) {
	route := e.Group("/api/v1/private/invoices", middlewares...)

	route.GET("", h.GetByUserID)
}
