package repository

import (
	"context"
	"fmt"

	"github.com/product-management/internal/domain/entity"
	"github.com/product-management/internal/domain/repository"
	"gorm.io/gorm"
)

// productRepositoryImpl implements the ProductRepository interface
type productRepositoryImpl struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

// Create creates a new product
func (r *productRepositoryImpl) Create(ctx context.Context, product *entity.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// GetByID retrieves a product by its ID
func (r *productRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}
	return &product, nil
}

// GetAll retrieves all products with optional filtering and pagination
func (r *productRepositoryImpl) GetAll(ctx context.Context, filter *repository.ProductFilter, offset, limit int) ([]*entity.Product, error) {
	var products []*entity.Product
	query := r.db.WithContext(ctx)

	if filter != nil {
		query = r.applyFilter(query, filter)
	}

	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

// GetTotalCount returns the total count of products with optional filtering
func (r *productRepositoryImpl) GetTotalCount(ctx context.Context, filter *repository.ProductFilter) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.Product{})

	if filter != nil {
		query = r.applyFilter(query, filter)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}

// Update updates an existing product
func (r *productRepositoryImpl) Update(ctx context.Context, product *entity.Product) error {
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

// Delete soft-deletes a product by its ID
func (r *productRepositoryImpl) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Product{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

// HardDelete permanently deletes a product by its ID
func (r *productRepositoryImpl) HardDelete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Unscoped().Delete(&entity.Product{}, id).Error; err != nil {
		return fmt.Errorf("failed to hard delete product: %w", err)
	}
	return nil
}

// GetByName retrieves a product by its name
func (r *productRepositoryImpl) GetByName(ctx context.Context, name string) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product by name: %w", err)
	}
	return &product, nil
}

// ExistsByName checks if a product with the given name exists
func (r *productRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Product{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check product existence by name: %w", err)
	}
	return count > 0, nil
}

// GetByCategory retrieves products by category
func (r *productRepositoryImpl) GetByCategory(ctx context.Context, category string, offset, limit int) ([]*entity.Product, error) {
	var products []*entity.Product
	if err := r.db.WithContext(ctx).
		Where("category = ? AND is_active = ?", category, true).
		Offset(offset).Limit(limit).
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to get products by category: %w", err)
	}
	return products, nil
}

// UpdateStock updates the stock quantity of a product
func (r *productRepositoryImpl) UpdateStock(ctx context.Context, id uint, stock int) error {
	if err := r.db.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", id).Update("stock", stock).Error; err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}
	return nil
}

// BulkUpdateStatus updates the active status of multiple products
func (r *productRepositoryImpl) BulkUpdateStatus(ctx context.Context, ids []uint, isActive bool) error {
	if err := r.db.WithContext(ctx).Model(&entity.Product{}).Where("id IN ?", ids).Update("is_active", isActive).Error; err != nil {
		return fmt.Errorf("failed to bulk update product status: %w", err)
	}
	return nil
}

// applyFilter applies filters to the query
func (r *productRepositoryImpl) applyFilter(query *gorm.DB, filter *repository.ProductFilter) *gorm.DB {
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}
	
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}
	
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	
	if filter.SearchTerm != "" {
		searchPattern := "%" + filter.SearchTerm + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}
	
	return query
}