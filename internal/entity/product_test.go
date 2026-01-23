package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct_NewProduct_Success(t *testing.T) {
	name := "Product A"
	price := 10.0

	product, err := NewProduct(name, price)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, name, product.Name)
	assert.Equal(t, price, product.Price)
}

func TestProduct_NewProduct_InvalidName(t *testing.T) {
	name := ""
	price := 10.0

	product, err := NewProduct(name, price)

	assert.Nil(t, product)
	assert.Equal(t, ErrInvalidProductName, err)
}

func TestProduct_NewProduct_InvalidPrice(t *testing.T) {
	name := "Product A"
	price := -5.0

	product, err := NewProduct(name, price)

	assert.Nil(t, product)
	assert.Equal(t, ErrInvalidProductPrice, err)
}
