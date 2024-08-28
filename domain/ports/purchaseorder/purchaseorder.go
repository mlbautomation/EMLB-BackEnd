package purchaseorder

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/mlbautomation/ProyectoEMLB/model"
)

type Repository interface {
	Create(m *model.PurchaseOrder) error
	GetByID(ID uuid.UUID) (model.PurchaseOrder, error)
}

type Service interface {
	Create(m *model.PurchaseOrder) error
	GetByID(ID uuid.UUID) (model.PurchaseOrder, error)
}

type ServiceInvoice interface {
	Validate(m *model.PurchaseOrder) error
}

type ServicePaypal interface {
	GetByID(ID uuid.UUID) (model.PurchaseOrder, error)
	TotalAmount(m model.PurchaseOrder) float64
}

type Handler interface {
	Create(c echo.Context) error
}
