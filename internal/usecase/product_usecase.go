package usecase

import (
	"context"
	"math"

	"github.com/product-management/internal/domain/entity"
	"github.com/product-management/internal/domain/repository"
	"github.com/product-management/internal/domain/service"
)

// productUseCase implements the ProductService interface
type productUseCase struct {
	productRepo repository.ProductRepository
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(productRepo repository.ProductRepository) service.ProductService {
	return &productUseCase{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product
func (uc *productUseCase) CreateProduct(ctx context.Context, req *service.ProductCreateRequest) (*entity.Product, error) {
	// Check if product with same name already exists
	exists, err := uc.productRepo.ExistsByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entity.ErrProductAlreadyExists
	}

	// Create product entity
	product := &entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		IsActive:    true, // New products are active by default
	}

	// Validate product
	if err := product.Validate(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID retrieves a product by its ID
func (uc *productUseCase) GetProductByID(ctx context.Context, id uint) (*entity.Product, error) {
	return uc.productRepo.GetByID(ctx, id)
}

// GetProducts retrieves a paginated list of products with filtering
func (uc *productUseCase) GetProducts(ctx context.Context, filter *repository.ProductFilter, page, pageSize int) (*service.ProductListResponse, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get products
	products, err := uc.productRepo.GetAll(ctx, filter, offset, pageSize)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := uc.productRepo.GetTotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &service.ProductListResponse{
		Products:   products,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateProduct updates an existing product
func (uc *productUseCase) UpdateProduct(ctx context.Context, id uint, req *service.ProductUpdateRequest) (*entity.Product, error) {
	// Get existing product
	product, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		// Check if another product with this name exists
		if *req.Name != product.Name {
			exists, err := uc.productRepo.ExistsByName(ctx, *req.Name)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, entity.ErrProductAlreadyExists
			}
		}
		product.Name = *req.Name
	}
	
	if req.Description != nil {
		product.Description = *req.Description
	}
	
	if req.Price != nil {
		product.Price = *req.Price
	}
	
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	
	if req.Category != nil {
		product.Category = *req.Category
	}
	
	if req.ImageURL != nil {
		product.ImageURL = *req.ImageURL
	}
	
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	// Validate updated product
	if err := product.Validate(); err != nil {
		return nil, err
	}

	// Save updates
	if err := uc.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct deletes a product by its ID
func (uc *productUseCase) DeleteProduct(ctx context.Context, id uint) error {
	// Check if product exists
	_, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.productRepo.Delete(ctx, id)
}

// GetProductsByCategory retrieves products by category
func (uc *productUseCase) GetProductsByCategory(ctx context.Context, category string, page, pageSize int) (*service.ProductListResponse, error) {
	filter := &repository.ProductFilter{
		Category: category,
		IsActive: boolPtr(true),
	}
	
	return uc.GetProducts(ctx, filter, page, pageSize)
}

// SearchProducts searches for products by name or description
func (uc *productUseCase) SearchProducts(ctx context.Context, searchTerm string, page, pageSize int) (*service.ProductListResponse, error) {
	filter := &repository.ProductFilter{
		SearchTerm: searchTerm,
		IsActive:   boolPtr(true),
	}
	
	return uc.GetProducts(ctx, filter, page, pageSize)
}

// UpdateProductStock updates the stock quantity of a product
func (uc *productUseCase) UpdateProductStock(ctx context.Context, id uint, stock int) error {
	// Check if product exists
	_, err := uc.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Validate stock
	if stock < 0 {
		return entity.ErrProductStockInvalid
	}

	return uc.productRepo.UpdateStock(ctx, id, stock)
}

// BulkUpdateProductStatus updates the active status of multiple products
func (uc *productUseCase) BulkUpdateProductStatus(ctx context.Context, ids []uint, isActive bool) error {
	if len(ids) == 0 {
		return entity.ErrInvalidInput
	}

	return uc.productRepo.BulkUpdateStatus(ctx, ids, isActive)
}

// Helper function to create a pointer to a boolean
func boolPtr(b bool) *bool {
	return &b
}