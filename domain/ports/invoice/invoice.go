package invoice

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mlbautomation/ProyectoEMLB/model"
)

type Repository interface {
	Create(m *model.Invoice, ms model.InvoiceDetails) error
}

// Dependecia de Repository: Encabezado
type RepositoryInvoiceReport interface {
	HeadByInvoiceID(ID uuid.UUID) (model.InvoiceReport, error)
	HeadsByUserID(userID uuid.UUID) (model.InvoicesReport, error)
	AllHead() (model.InvoicesReport, error)
	AllDetailsByInvoiceID(ID uuid.UUID) (model.InvoiceDetailsReports, error)
}

type Service interface {
	Create(m *model.PurchaseOrder) error
	GetByUserID(userID uuid.UUID) (model.InvoicesReport, error)
	GetAll() (model.InvoicesReport, error)
}

type ServicePaypal interface {
	Create(m *model.PurchaseOrder) error
}

type Handler interface {
	GetByUserID(c echo.Context) error
	GetAll(c echo.Context) error
}
