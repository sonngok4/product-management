package repository

import (
	"context"

	"github.com/product-management/internal/domain/entity"
)

// ProductFilter represents filtering criteria for products
type ProductFilter struct {
	Category    string
	MinPrice    *float64
	MaxPrice    *float64
	IsActive    *bool
	SearchTerm  string // for searching in name or description
}

// ProductRepository defines the interface for product repository operations
type ProductRepository interface {
	// Create creates a new product
	Create(ctx context.Context, product *entity.Product) error
	
	// GetByID retrieves a product by its ID
	GetByID(ctx context.Context, id uint) (*entity.Product, error)
	
	// GetAll retrieves all products with optional filtering and pagination
	GetAll(ctx context.Context, filter *ProductFilter, offset, limit int) ([]*entity.Product, error)
	
	// GetTotalCount returns the total count of products with optional filtering
	GetTotalCount(ctx context.Context, filter *ProductFilter) (int64, error)
	
	// Update updates an existing product
	Update(ctx context.Context, product *entity.Product) error
	
	// Delete soft-deletes a product by its ID
	Delete(ctx context.Context, id uint) error
	
	// HardDelete permanently deletes a product by its ID
	HardDelete(ctx context.Context, id uint) error
	
	// GetByName retrieves a product by its name
	GetByName(ctx context.Context, name string) (*entity.Product, error)
	
	// ExistsByName checks if a product with the given name exists
	ExistsByName(ctx context.Context, name string) (bool, error)
	
	// GetByCategory retrieves products by category
	GetByCategory(ctx context.Context, category string, offset, limit int) ([]*entity.Product, error)
	
	// UpdateStock updates the stock quantity of a product
	UpdateStock(ctx context.Context, id uint, stock int) error
	
	// BulkUpdateStatus updates the active status of multiple products
	BulkUpdateStatus(ctx context.Context, ids []uint, isActive bool) error
}
