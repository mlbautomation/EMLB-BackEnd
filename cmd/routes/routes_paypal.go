package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/paypal"
)

// publicRoutes handle the routes that not requires a validation of any kind to be use
func PaypalPublic(e *echo.Echo, h paypal.Handler) {
	//recomendación de no poner "paypal" muy obvio para ataques
	route := e.Group("/api/v1/public/paypal")

	route.POST("", h.Webhook)
}
