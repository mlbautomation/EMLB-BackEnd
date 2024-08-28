package paypal

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Service interface {
	ProcessRequest(header http.Header, body []byte) error
}

type Handler interface {
	Webhook(c echo.Context) error
}
