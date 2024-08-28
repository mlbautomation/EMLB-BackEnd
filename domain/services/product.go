package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mlbautomation/ProyectoEMLB/domain/ports/product"
	"github.com/mlbautomation/ProyectoEMLB/model"
)

type Product struct {
	Repository product.Repository
}

func NewProduct(pr product.Repository) Product {
	return Product{Repository: pr}
}

func (p Product) Create(m *model.Product) error {

	ID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("%s %w", "uuid.NewUUID()", err)
	}
	m.ID = ID

	if m.ProductName == "" {
		return fmt.Errorf("%s", "product name is empty!")
	}

	//Price could be 0 for free

	//Images could be empty?
	if len(m.Images) == 0 {
		m.Images = []byte(`[]`)
	}

	if len(m.Features) == 0 {
		m.Features = []byte(`[]`)
	}

	m.CreatedAt = time.Now().Unix()

	err = p.Repository.Create(m)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.Create(m)", err)
	}

	return nil
}

func (p Product) Update(m *model.Product) error {
	if !m.HasID() {
		return fmt.Errorf("product: %w", model.ErrInvalidID)
	}

	if len(m.Images) == 0 {
		m.Images = []byte(`[]`)
	}
	if len(m.Features) == 0 {
		m.Features = []byte(`[]`)
	}

	m.UpdatedAt = time.Now().Unix()

	err := p.Repository.Update(m)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.Update(m)", err)
	}
	return nil
}

func (p Product) Delete(ID uuid.UUID) error {
	err := p.Repository.Delete(ID)
	if err != nil {
		return fmt.Errorf("%s %w", "Repository.Delete(ID)", err)
	}
	return nil
}

func (p Product) GetByID(ID uuid.UUID) (model.Product, error) {
	product, err := p.Repository.GetByID(ID)
	if err != nil {
		return model.Product{}, fmt.Errorf("%s %w", "Repository.GetByID(ID)", err)
	}
	return product, nil
}

func (p Product) GetAll() (model.Products, error) {
	products, err := p.Repository.GetAll()
	if err != nil {
		return model.Products{}, fmt.Errorf("%s %w", "Repository.GetAll()", err)
	}
	return products, nil
}
