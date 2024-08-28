package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/product"
)

// Para autorizaciones:
// func productAdminRoutes(e *echo.Echo, h productPorts.ProductHandlers, middlewares ...echo.MiddlewareFunc) {
func ProductAdmin(e *echo.Echo, h product.Handler, middlewares ...echo.MiddlewareFunc) {
	// g := e.Group("/api/v1/admin/products, middlewares ...")
	g := e.Group("/api/v1/admin/products", middlewares...)

	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.GET("", h.GetAll)
	g.GET("/:id", h.GetByID)
}

func ProductPublic(e *echo.Echo, h product.Handler) {
	g := e.Group("/api/v1/public/products")

	g.GET("", h.GetAll)
	g.GET("/:id", h.GetByID)
}
