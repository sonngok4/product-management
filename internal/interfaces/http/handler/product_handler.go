package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/product-management/internal/domain/entity"
	"github.com/product-management/internal/domain/repository"
	"github.com/product-management/internal/domain/service"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	productService service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided information
// @Tags products
// @Accept json
// @Produce json
// @Param product body service.ProductCreateRequest true "Product creation request"
// @Success 201 {object} entity.Product
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req service.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get a product by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} entity.Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid product ID",
		})
		return
	}

	product, err := h.productService.GetProductByID(c.Request.Context(), uint(id))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProducts godoc
// @Summary Get products with filtering and pagination
// @Description Get a list of products with optional filtering and pagination
// @Tags products
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Param category query string false "Filter by category"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param is_active query boolean false "Filter by active status"
// @Param search query string false "Search in name and description"
// @Success 200 {object} service.ProductListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Parse filter parameters
	filter := &repository.ProductFilter{}
	
	if category := c.Query("category"); category != "" {
		filter.Category = category
	}
	
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = &minPrice
		}
	}
	
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = &maxPrice
		}
	}
	
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filter.IsActive = &isActive
		}
	}
	
	if search := c.Query("search"); search != "" {
		filter.SearchTerm = search
	}

	products, err := h.productService.GetProducts(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product with the provided information
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body service.ProductUpdateRequest true "Product update request"
// @Success 200 {object} entity.Product
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid product ID",
		})
		return
	}

	var req service.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), uint(id), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid product ID",
		})
		return
	}

	err = h.productService.DeleteProduct(c.Request.Context(), uint(id))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product deleted successfully",
	})
}

// SearchProducts godoc
// @Summary Search products
// @Description Search for products by name or description
// @Tags products
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} service.ProductListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Search query is required",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	products, err := h.productService.SearchProducts(c.Request.Context(), query, page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

// UpdateProductStock godoc
// @Summary Update product stock
// @Description Update the stock quantity of a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body StockUpdateRequest true "Stock update request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/products/{id}/stock [put]
func (h *ProductHandler) UpdateProductStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid product ID",
		})
		return
	}

	var req StockUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	err = h.productService.UpdateProductStock(c.Request.Context(), uint(id), req.Stock)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product stock updated successfully",
	})
}

// handleError handles different types of errors and returns appropriate HTTP responses
func handleError(c *gin.Context, err error) {
	switch err {
	case entity.ErrProductNotFound:
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Not Found",
			Message: err.Error(),
		})
	case entity.ErrProductAlreadyExists:
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "Conflict",
			Message: err.Error(),
		})
	case entity.ErrProductNameRequired, entity.ErrProductNameTooShort, entity.ErrProductNameTooLong,
		 entity.ErrProductPriceInvalid, entity.ErrProductStockInvalid, entity.ErrInvalidInput:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "An unexpected error occurred",
			Details: err.Error(),
		})
	}
}