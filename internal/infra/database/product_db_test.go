package database

import (
	"testing"

	"github.com/robsonalvesdevbr/apis-go/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestProduct_DB_CreateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product 1", 10.0)
	productDB := NewProduct(db)

	err = productDB.Create(product)
	assert.NoError(t, err)

	var productFound entity.Product
	err = db.First(&productFound, "id = ?", product.ID).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID, productFound.ID)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestProduct_DB_FindAll(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&entity.Product{})
	productDB := NewProduct(db)

	products := []struct {
		name  string
		price float64
	}{
		{"Product 1", 10.0},
		{"Product 2", 20.0},
		{"Product 3", 30.0},
		{"Product 4", 40.0},
		{"Product 5", 50.0},
	}

	for _, p := range products {
		product, _ := entity.NewProduct(p.name, p.price)
		err := productDB.Create(product)
		assert.NoError(t, err)
	}

	page := 1
	limit := 2
	result, err := productDB.FindAll(page, limit, "asc")
	assert.NoError(t, err)
	assert.Len(t, result, limit)
	assert.Equal(t, "Product 1", result[0].Name)
	assert.Equal(t, "Product 2", result[1].Name)

	page = 2
	result, err = productDB.FindAll(page, limit, "asc")
	assert.NoError(t, err)
	assert.Len(t, result, limit)
	assert.Equal(t, "Product 3", result[0].Name)
	assert.Equal(t, "Product 4", result[1].Name)

	page = 3
	result, err = productDB.FindAll(page, limit, "asc")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Product 5", result[0].Name)
}

func TestProduct_DB_FindByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product 1", 10.0)
	productDB := NewProduct(db)

	err = productDB.Create(product)
	assert.NoError(t, err)

	productFound, err := productDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestProduct_DB_Update(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product 1", 10.0)
	productDB := NewProduct(db)

	err = productDB.Create(product)
	assert.NoError(t, err)

	product.Name = "Updated Product"
	product.Price = 20.0
	err = productDB.Update(product)
	assert.NoError(t, err)

	var productFound entity.Product
	err = db.First(&productFound, "id = ?", product.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Updated Product", productFound.Name)
	assert.Equal(t, 20.0, productFound.Price)
}

func TestProduct_DB_Delete(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product 1", 10.0)
	productDB := NewProduct(db)

	err = productDB.Create(product)
	assert.NoError(t, err)

	err = productDB.Delete(product.ID.String())
	assert.NoError(t, err)

	var productFound entity.Product
	err = db.First(&productFound, "id = ?", product.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
