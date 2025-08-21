package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/product-management/internal/config"
	"github.com/product-management/internal/domain/service"
	"github.com/product-management/internal/infrastructure/database"
	"github.com/product-management/internal/infrastructure/repository"
	"github.com/product-management/internal/interfaces/http/router"
	"github.com/product-management/internal/usecase"
	"github.com/product-management/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"time"
)

// IntegrationTestSuite represents the integration test suite
type IntegrationTestSuite struct {
	suite.Suite
	app            *gin.Engine
	db             *database.Database
	authService    service.AuthService
	productService service.ProductService
	testUser       *TestUser
}

// TestUser represents a test user
type TestUser struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// SetupSuite sets up the test suite
func (suite *IntegrationTestSuite) SetupSuite() {
	// Set test environment
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "product_management_test")
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("GIN_MODE", "test")

	// Load configuration
	cfg := config.LoadConfig()

	// Skip database tests if not available
	db, err := database.NewDatabase(cfg)
	if err != nil {
		suite.T().Skip("Database not available for integration tests")
		return
	}

	suite.db = db

	// Run migrations
	err = db.AutoMigrate()
	if err != nil {
		suite.T().Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.GetDB())
	productRepo := repository.NewProductRepository(db.GetDB())

	// Initialize services
	tokenManager := jwt.NewTokenManager(cfg.JWT.Secret, time.Hour)
	suite.authService = usecase.NewAuthUseCase(userRepo, tokenManager)
	suite.productService = usecase.NewProductUseCase(productRepo)

	// Setup router
	suite.app = router.SetupRouter(cfg, db, suite.productService, suite.authService)

	// Create a test user
	suite.createTestUser()
}

// TearDownSuite cleans up after tests
func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		// Clean up test data
		suite.db.GetDB().Exec("DELETE FROM products")
		suite.db.GetDB().Exec("DELETE FROM users")
		suite.db.Close()
	}
}

// createTestUser creates a test user for authentication
func (suite *IntegrationTestSuite) createTestUser() {
	registerReq := service.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	reqBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	if w.Code == http.StatusCreated {
		var response service.AuthResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		
		suite.testUser = &TestUser{
			ID:    response.User.ID,
			Email: response.User.Email,
			Token: response.Token.AccessToken,
		}
	}
}

// TestHealthCheck tests the health check endpoint
func (suite *IntegrationTestSuite) TestHealthCheck() {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	suite.app.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "healthy", response["status"])
}

// TestAuthRegisterAndLogin tests user registration and login
func (suite *IntegrationTestSuite) TestAuthRegisterAndLogin() {
	// Test Registration
	registerReq := service.RegisterRequest{
		Email:     "integration@example.com",
		Username:  "integration",
		Password:  "password123",
		FirstName: "Integration",
		LastName:  "Test",
	}

	reqBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var registerResponse service.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &registerResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), registerReq.Email, registerResponse.User.Email)
	assert.NotEmpty(suite.T(), registerResponse.Token.AccessToken)

	// Test Login
	loginReq := service.LoginRequest{
		Email:    registerReq.Email,
		Password: registerReq.Password,
	}

	reqBody, _ = json.Marshal(loginReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var loginResponse service.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), loginReq.Email, loginResponse.User.Email)
	assert.NotEmpty(suite.T(), loginResponse.Token.AccessToken)
}

// TestProductCRUD tests the complete CRUD operations for products
func (suite *IntegrationTestSuite) TestProductCRUD() {
	if suite.testUser == nil {
		suite.T().Skip("Test user not available")
		return
	}

	// Test Create Product
	createReq := service.ProductCreateRequest{
		Name:        "Integration Test Product",
		Description: "A product for integration testing",
		Price:       99.99,
		Stock:       50,
		Category:    "Test Category",
		ImageURL:    "http://example.com/image.jpg",
	}

	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.testUser.Token)

	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var createdProduct map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &createdProduct)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createReq.Name, createdProduct["name"])

	productID := uint(createdProduct["id"].(float64))

	// Test Get Product
	req = httptest.NewRequest(http.MethodGet, "/api/v1/products/"+string(rune(productID+'0')), nil)
	w = httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Test Update Product
	updateReq := service.ProductUpdateRequest{
		Name: stringPtr("Updated Integration Test Product"),
	}

	reqBody, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPut, "/api/v1/products/"+string(rune(productID+'0')), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.testUser.Token)

	w = httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Test Delete Product
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+string(rune(productID+'0')), nil)
	req.Header.Set("Authorization", "Bearer "+suite.testUser.Token)

	w = httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// TestUnauthorizedAccess tests unauthorized access to protected endpoints
func (suite *IntegrationTestSuite) TestUnauthorizedAccess() {
	createReq := service.ProductCreateRequest{
		Name:  "Unauthorized Test",
		Price: 10.0,
	}

	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}

// TestIntegrationTestSuite runs the integration test suite
func TestIntegrationTestSuite(t *testing.T) {
	// Skip integration tests if requested
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	suite.Run(t, new(IntegrationTestSuite))
}