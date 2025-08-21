package unit

import (
	"context"
	"testing"

	"github.com/product-management/internal/domain/entity"
	"github.com/product-management/internal/domain/repository"
	"github.com/product-management/internal/domain/service"
	"github.com/product-management/internal/usecase"
	"github.com/product-management/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ProductUseCaseTestSuite represents the test suite for product use case
type ProductUseCaseTestSuite struct {
	suite.Suite
	mockRepo       *mocks.MockProductRepository
	productService service.ProductService
	ctx            context.Context
}

// SetupTest sets up the test suite
func (suite *ProductUseCaseTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.MockProductRepository)
	suite.productService = usecase.NewProductUseCase(suite.mockRepo)
	suite.ctx = context.Background()
}

// TestCreateProduct_Success tests successful product creation
func (suite *ProductUseCaseTestSuite) TestCreateProduct_Success() {
	// Arrange
	req := &service.ProductCreateRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       29.99,
		Stock:       100,
		Category:    "Electronics",
		ImageURL:    "http://example.com/image.jpg",
	}

	suite.mockRepo.On("ExistsByName", suite.ctx, req.Name).Return(false, nil)
	suite.mockRepo.On("Create", suite.ctx, mock.AnythingOfType("*entity.Product")).Return(nil).Run(func(args mock.Arguments) {
		product := args.Get(1).(*entity.Product)
		product.ID = 1 // Simulate database assigning an ID
	})

	// Act
	result, err := suite.productService.CreateProduct(suite.ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), req.Name, result.Name)
	assert.Equal(suite.T(), req.Description, result.Description)
	assert.Equal(suite.T(), req.Price, result.Price)
	assert.Equal(suite.T(), req.Stock, result.Stock)
	assert.Equal(suite.T(), req.Category, result.Category)
	assert.Equal(suite.T(), req.ImageURL, result.ImageURL)
	assert.True(suite.T(), result.IsActive)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestCreateProduct_ProductAlreadyExists tests product creation when product name already exists
func (suite *ProductUseCaseTestSuite) TestCreateProduct_ProductAlreadyExists() {
	// Arrange
	req := &service.ProductCreateRequest{
		Name:        "Existing Product",
		Description: "Test Description",
		Price:       29.99,
		Stock:       100,
		Category:    "Electronics",
	}

	suite.mockRepo.On("ExistsByName", suite.ctx, req.Name).Return(true, nil)

	// Act
	result, err := suite.productService.CreateProduct(suite.ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), entity.ErrProductAlreadyExists, err)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestCreateProduct_InvalidData tests product creation with invalid data
func (suite *ProductUseCaseTestSuite) TestCreateProduct_InvalidData() {
	// Arrange
	req := &service.ProductCreateRequest{
		Name:        "", // Invalid: empty name
		Description: "Test Description",
		Price:       -10, // Invalid: negative price
		Stock:       100,
		Category:    "Electronics",
	}

	suite.mockRepo.On("ExistsByName", suite.ctx, req.Name).Return(false, nil)

	// Act
	result, err := suite.productService.CreateProduct(suite.ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), entity.ErrProductNameRequired, err)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetProductByID_Success tests successful product retrieval by ID
func (suite *ProductUseCaseTestSuite) TestGetProductByID_Success() {
	// Arrange
	expectedProduct := &entity.Product{
		ID:          1,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       29.99,
		Stock:       100,
		Category:    "Electronics",
		IsActive:    true,
	}

	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(expectedProduct, nil)

	// Act
	result, err := suite.productService.GetProductByID(suite.ctx, 1)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedProduct.ID, result.ID)
	assert.Equal(suite.T(), expectedProduct.Name, result.Name)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetProductByID_NotFound tests product retrieval when product doesn't exist
func (suite *ProductUseCaseTestSuite) TestGetProductByID_NotFound() {
	// Arrange
	suite.mockRepo.On("GetByID", suite.ctx, uint(999)).Return(nil, entity.ErrProductNotFound)

	// Act
	result, err := suite.productService.GetProductByID(suite.ctx, 999)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), entity.ErrProductNotFound, err)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestGetProducts_Success tests successful products listing with pagination
func (suite *ProductUseCaseTestSuite) TestGetProducts_Success() {
	// Arrange
	expectedProducts := []*entity.Product{
		{ID: 1, Name: "Product 1", Price: 10.00, IsActive: true},
		{ID: 2, Name: "Product 2", Price: 20.00, IsActive: true},
	}
	filter := &repository.ProductFilter{IsActive: boolPtr(true)}
	
	suite.mockRepo.On("GetAll", suite.ctx, filter, 0, 10).Return(expectedProducts, nil)
	suite.mockRepo.On("GetTotalCount", suite.ctx, filter).Return(int64(2), nil)

	// Act
	result, err := suite.productService.GetProducts(suite.ctx, filter, 1, 10)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), 2, len(result.Products))
	assert.Equal(suite.T(), int64(2), result.Total)
	assert.Equal(suite.T(), 1, result.Page)
	assert.Equal(suite.T(), 10, result.PageSize)
	assert.Equal(suite.T(), 1, result.TotalPages)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestUpdateProduct_Success tests successful product update
func (suite *ProductUseCaseTestSuite) TestUpdateProduct_Success() {
	// Arrange
	existingProduct := &entity.Product{
		ID:          1,
		Name:        "Old Product",
		Description: "Old Description",
		Price:       29.99,
		Stock:       100,
		Category:    "Electronics",
		IsActive:    true,
	}

	newName := "Updated Product"
	newPrice := 39.99
	req := &service.ProductUpdateRequest{
		Name:  &newName,
		Price: &newPrice,
	}

	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingProduct, nil)
	suite.mockRepo.On("ExistsByName", suite.ctx, newName).Return(false, nil)
	suite.mockRepo.On("Update", suite.ctx, mock.AnythingOfType("*entity.Product")).Return(nil)

	// Act
	result, err := suite.productService.UpdateProduct(suite.ctx, 1, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), newName, result.Name)
	assert.Equal(suite.T(), newPrice, result.Price)
	assert.Equal(suite.T(), existingProduct.Description, result.Description) // Unchanged

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestUpdateProduct_ProductNotFound tests product update when product doesn't exist
func (suite *ProductUseCaseTestSuite) TestUpdateProduct_ProductNotFound() {
	// Arrange
	newName := "Updated Product"
	req := &service.ProductUpdateRequest{
		Name: &newName,
	}

	suite.mockRepo.On("GetByID", suite.ctx, uint(999)).Return(nil, entity.ErrProductNotFound)

	// Act
	result, err := suite.productService.UpdateProduct(suite.ctx, 999, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), entity.ErrProductNotFound, err)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestDeleteProduct_Success tests successful product deletion
func (suite *ProductUseCaseTestSuite) TestDeleteProduct_Success() {
	// Arrange
	existingProduct := &entity.Product{
		ID:       1,
		Name:     "Product to Delete",
		IsActive: true,
	}

	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingProduct, nil)
	suite.mockRepo.On("Delete", suite.ctx, uint(1)).Return(nil)

	// Act
	err := suite.productService.DeleteProduct(suite.ctx, 1)

	// Assert
	assert.NoError(suite.T(), err)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestDeleteProduct_ProductNotFound tests product deletion when product doesn't exist
func (suite *ProductUseCaseTestSuite) TestDeleteProduct_ProductNotFound() {
	// Arrange
	suite.mockRepo.On("GetByID", suite.ctx, uint(999)).Return(nil, entity.ErrProductNotFound)

	// Act
	err := suite.productService.DeleteProduct(suite.ctx, 999)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), entity.ErrProductNotFound, err)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestUpdateProductStock_Success tests successful product stock update
func (suite *ProductUseCaseTestSuite) TestUpdateProductStock_Success() {
	// Arrange
	existingProduct := &entity.Product{
		ID:    1,
		Name:  "Product",
		Stock: 100,
	}

	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingProduct, nil)
	suite.mockRepo.On("UpdateStock", suite.ctx, uint(1), 200).Return(nil)

	// Act
	err := suite.productService.UpdateProductStock(suite.ctx, 1, 200)

	// Assert
	assert.NoError(suite.T(), err)

	suite.mockRepo.AssertExpectations(suite.T())
}

// TestUpdateProductStock_InvalidStock tests stock update with invalid stock value
func (suite *ProductUseCaseTestSuite) TestUpdateProductStock_InvalidStock() {
	// Arrange
	existingProduct := &entity.Product{
		ID:    1,
		Name:  "Product",
		Stock: 100,
	}

	suite.mockRepo.On("GetByID", suite.ctx, uint(1)).Return(existingProduct, nil)

	// Act
	err := suite.productService.UpdateProductStock(suite.ctx, 1, -10)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), entity.ErrProductStockInvalid, err)

	suite.mockRepo.AssertExpectations(suite.T())
}

// Helper function to create a pointer to a boolean
func boolPtr(b bool) *bool {
	return &b
}

// TestProductUseCaseTestSuite runs the test suite
func TestProductUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ProductUseCaseTestSuite))
}