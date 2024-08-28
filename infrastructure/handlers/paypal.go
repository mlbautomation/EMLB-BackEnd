package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/mlbautomation/ProyectoEMLB/domain/ports/paypal"
	"github.com/mlbautomation/ProyectoEMLB/infrastructure/handlers/response"
)

type Paypal struct {
	service   paypal.Service
	responser response.API
}

func NewPaypal(pps paypal.Service) *Paypal {
	return &Paypal{service: pps}
}

func (h *Paypal) Webhook(c echo.Context) error {
	//recibe el body que esta haciendo paypal a nosotros
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return h.responser.Error(c, "handlers-Paypal-Webhook-io.ReadAll(c.Request().Body)", err)
	}
	go func() {
		//1. Procesar la data para crear las estructuras que necesito
		//2. Validar esa data con paypal
		//3. Procesar el pago
		err = h.service.ProcessRequest(c.Request().Header, body)
		if err != nil {
			log.Printf("useCasePayPal.ProcessRequest(): %v", err)
		}
	}()

	/* Creo yo que si fuera una respuesta al cliente usaría el
	formato que estamos trabajando:
	return c.JSON(h.responser.OK("Webhook recibido!"))
	pero esta respuesta es a PAYPAL no al cliente, por lo tanto:*/

	return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
}
