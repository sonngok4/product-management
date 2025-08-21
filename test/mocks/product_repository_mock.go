package mocks

import (
	"context"

	"github.com/product-management/internal/domain/entity"
	"github.com/product-management/internal/domain/repository"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockProductRepository) Create(ctx context.Context, product *entity.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *MockProductRepository) GetByID(ctx context.Context, id uint) (*entity.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

// GetAll mocks the GetAll method
func (m *MockProductRepository) GetAll(ctx context.Context, filter *repository.ProductFilter, offset, limit int) ([]*entity.Product, error) {
	args := m.Called(ctx, filter, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Product), args.Error(1)
}

// GetTotalCount mocks the GetTotalCount method
func (m *MockProductRepository) GetTotalCount(ctx context.Context, filter *repository.ProductFilter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

// Update mocks the Update method
func (m *MockProductRepository) Update(ctx context.Context, product *entity.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockProductRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// HardDelete mocks the HardDelete method
func (m *MockProductRepository) HardDelete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetByName mocks the GetByName method
func (m *MockProductRepository) GetByName(ctx context.Context, name string) (*entity.Product, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

// ExistsByName mocks the ExistsByName method
func (m *MockProductRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(bool), args.Error(1)
}

// GetByCategory mocks the GetByCategory method
func (m *MockProductRepository) GetByCategory(ctx context.Context, category string, offset, limit int) ([]*entity.Product, error) {
	args := m.Called(ctx, category, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Product), args.Error(1)
}

// UpdateStock mocks the UpdateStock method
func (m *MockProductRepository) UpdateStock(ctx context.Context, id uint, stock int) error {
	args := m.Called(ctx, id, stock)
	return args.Error(0)
}

// BulkUpdateStatus mocks the BulkUpdateStatus method
func (m *MockProductRepository) BulkUpdateStatus(ctx context.Context, ids []uint, isActive bool) error {
	args := m.Called(ctx, ids, isActive)
	return args.Error(0)
}