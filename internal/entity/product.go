package entity

import (
	"errors"
	"time"

	"github.com/robsonalvesdevbr/apis-go/pkg/entity"
)

var (
	ErrInvalidProductName  = errors.New("product name must not be empty")
	ErrInvalidProductPrice = errors.New("product price must be greater than zero")
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProduct(name string, price float64) (*Product, error) {
	uuidNew, errGuid := entity.NewID()
	if errGuid != nil {
		return nil, errGuid
	}

	product := &Product{
		ID:        uuidNew,
		Name:      name,
		Price:     price,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}

	return product, nil
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrInvalidProductName
	}
	if p.Price <= 0 {
		return ErrInvalidProductPrice
	}
	return nil
}
