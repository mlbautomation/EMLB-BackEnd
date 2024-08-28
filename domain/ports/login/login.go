package login

import (
	"github.com/labstack/echo/v4"
	"github.com/mlbautomation/ProyectoEMLB/model"
)

type Service interface {
	Login(email, password, jwtSecretKey string) (model.User, string, error)
}

type Handler interface {
	Login(c echo.Context) error
}
