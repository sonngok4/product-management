package service

import (
	"context"

	"github.com/product-management/internal/domain/entity"
	"github.com/product-management/internal/domain/repository"
)

// ProductCreateRequest represents a request to create a product
type ProductCreateRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Stock       int     `json:"stock" validate:"min=0"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}

// ProductUpdateRequest represents a request to update a product
type ProductUpdateRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,min=0"`
	Category    *string  `json:"category,omitempty"`
	ImageURL    *string  `json:"image_url,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

// ProductListResponse represents a paginated list of products
type ProductListResponse struct {
	Products   []*entity.Product `json:"products"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// ProductService defines the interface for product business logic operations
type ProductService interface {
	// CreateProduct creates a new product
	CreateProduct(ctx context.Context, req *ProductCreateRequest) (*entity.Product, error)
	
	// GetProductByID retrieves a product by its ID
	GetProductByID(ctx context.Context, id uint) (*entity.Product, error)
	
	// GetProducts retrieves a paginated list of products with filtering
	GetProducts(ctx context.Context, filter *repository.ProductFilter, page, pageSize int) (*ProductListResponse, error)
	
	// UpdateProduct updates an existing product
	UpdateProduct(ctx context.Context, id uint, req *ProductUpdateRequest) (*entity.Product, error)
	
	// DeleteProduct deletes a product by its ID
	DeleteProduct(ctx context.Context, id uint) error
	
	// GetProductsByCategory retrieves products by category
	GetProductsByCategory(ctx context.Context, category string, page, pageSize int) (*ProductListResponse, error)
	
	// SearchProducts searches for products by name or description
	SearchProducts(ctx context.Context, searchTerm string, page, pageSize int) (*ProductListResponse, error)
	
	// UpdateProductStock updates the stock quantity of a product
	UpdateProductStock(ctx context.Context, id uint, stock int) error
	
	// BulkUpdateProductStatus updates the active status of multiple products
	BulkUpdateProductStatus(ctx context.Context, ids []uint, isActive bool) error
}