package product

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/mlbautomation/ProyectoEMLB/model"
)

type Repository interface {
	Create(m *model.Product) error
	Update(m *model.Product) error
	Delete(ID uuid.UUID) error

	GetByID(ID uuid.UUID) (model.Product, error)
	GetAll() (model.Products, error)
}

type RepositoryPurchaseOrder interface {
	GetByID(ID uuid.UUID) (model.Product, error)
}

type Service interface {
	Create(m *model.Product) error
	Update(m *model.Product) error
	Delete(ID uuid.UUID) error

	GetByID(ID uuid.UUID) (model.Product, error)
	GetAll() (model.Products, error)
}

type Handler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	GetByID(c echo.Context) error
	GetAll(c echo.Context) error
}
