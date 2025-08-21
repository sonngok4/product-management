package usecase

import (
	"context"
	"errors"

	"github.com/product-management/internal/domain/entity"
	"github.com/product-management/internal/domain/repository"
)

// ProductUseCase handles product business logic
type ProductUseCase struct {
	productRepo repository.ProductRepository
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(productRepo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
	}
}

// CreateProductRequest represents create product request data
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock" binding:"gte=0"`
}

// UpdateProductRequest represents update product request data
type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Category    *string  `json:"category"`
	Stock       *int     `json:"stock"`
}

// CreateProduct creates a new product
func (uc *ProductUseCase) CreateProduct(req *CreateProductRequest) (*entity.Product, error) {
	product := &entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Stock:       req.Stock,
	}

	if err := uc.productRepo.Create(context.Background(), product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (uc *ProductUseCase) GetProduct(id uint) (*entity.Product, error) {
	product, err := uc.productRepo.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetAllProducts retrieves all products with optional filtering
func (uc *ProductUseCase) GetAllProducts(category string, limit, offset int) ([]*entity.Product, error) {
	filter := &repository.ProductFilter{
		Category: category,
	}
	
	products, err := uc.productRepo.GetAll(context.Background(), filter, offset, limit)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// UpdateProduct updates an existing product
func (uc *ProductUseCase) UpdateProduct(id uint, req *UpdateProductRequest) (*entity.Product, error) {
	product, err := uc.productRepo.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Category != nil {
		product.Category = *req.Category
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if err := uc.productRepo.Update(context.Background(), product); err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct deletes a product
func (uc *ProductUseCase) DeleteProduct(id uint) error {
	return uc.productRepo.Delete(context.Background(), id)
}

// UpdateStock updates product stock
func (uc *ProductUseCase) UpdateStock(id uint, quantity int) error {
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	return uc.productRepo.UpdateStock(context.Background(), id, quantity)
}
