package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/mlbautomation/ProyectoEMLB/cmd/routes"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/invoice"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/login"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/paypal"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/product"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/purchaseorder"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/user"
	"github.com/mlbautomation/ProyectoEMLB/infrastructure/handlers/middle"
	"github.com/mlbautomation/ProyectoEMLB/infrastructure/handlers/response"
)

type Server struct {
	uHandler  user.Handler
	pHandler  product.Handler
	poHandler purchaseorder.Handler
	lHandler  login.Handler
	ppHandler paypal.Handler
	irHandler invoice.Handler
}

func NewServer(
	uHandler user.Handler,
	pHandler product.Handler,
	poHandler purchaseorder.Handler,
	lHandler login.Handler,
	ppHandler paypal.Handler,
	irHandler invoice.Handler) *Server {

	return &Server{
		uHandler:  uHandler,
		pHandler:  pHandler,
		poHandler: poHandler,
		lHandler:  lHandler,
		ppHandler: ppHandler,
		irHandler: irHandler,
	}
}

func (s *Server) Initialize() {

	e := NewHTTP(response.HTTPErrorHandler)

	health(e) //esto es para verificar que el servicio está funcionando

	authMiddleware := middle.New()

	routes.UserAdmin(e, s.uHandler, authMiddleware.IsValid, authMiddleware.IsAdmin)
	routes.UserPublic(e, s.uHandler)

	routes.ProductAdmin(e, s.pHandler, authMiddleware.IsValid, authMiddleware.IsAdmin)
	routes.ProductPublic(e, s.pHandler)

	routes.PurchaseOrderPrivate(e, s.poHandler, authMiddleware.IsValid)

	routes.LoginPublic(e, s.lHandler)

	routes.PaypalPublic(e, s.ppHandler)

	routes.InvoiceAdmin(e, s.irHandler, authMiddleware.IsValid, authMiddleware.IsAdmin)
	routes.InvoicePrivate(e, s.irHandler, authMiddleware.IsValid)

	err := e.Start(":" + os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal(err)
	}
}

func health(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(
			http.StatusOK,
			map[string]string{
				"time":         time.Now().String(),
				"message":      "Hola mundo desde MLB!",
				"service_name": "Proyecto Ecommerce",
			},
		)
	})
}
