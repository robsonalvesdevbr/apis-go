package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/robsonalvesdevbr/apis-go/internal/entity"
	"gorm.io/gorm"
)

type Product struct {
	DB *gorm.DB
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{DB: db}
}

func (p *Product) Create(product *entity.Product) error {
	ctx := context.Background()
	return p.DB.WithContext(ctx).Create(product).Error
}

func (p *Product) FindAll(page, limit int, sort string) ([]entity.Product, error) {
	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	var products []entity.Product
	if page != 0 && limit != 0 {
		offset := (page - 1) * limit
		if err := p.DB.Order("created_at " + sort).Limit(limit).Offset(offset).Find(&products).Error; err != nil {
			return nil, err
		}
	} else {
		if err := p.DB.Order("created_at " + sort).Find(&products).Error; err != nil {
			return nil, err
		}
	}

	for i := range products {
		products[i].CreatedAt = products[i].CreatedAt.UTC()
		products[i].UpdatedAt = products[i].UpdatedAt.UTC()
	}

	return products, nil
}

func (p *Product) FindByID(id string) (*entity.Product, error) {
	idd, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var product entity.Product
	if err := p.DB.First(&product, "id = ?", idd).Error; err != nil {
		return nil, err
	}
	product.CreatedAt = product.CreatedAt.UTC()
	product.UpdatedAt = product.UpdatedAt.UTC()
	return &product, nil
}

func (p *Product) Update(product *entity.Product) error {
	ctx := context.Background()
	return p.DB.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", product.ID).Updates(product).Error
}

func (p *Product) Delete(id string) error {
	ctx := context.Background()
	return p.DB.WithContext(ctx).Delete(&entity.Product{}, "id = ?", id).Error
}
