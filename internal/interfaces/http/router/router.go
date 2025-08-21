package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/product-management/internal/config"
	"github.com/product-management/internal/infrastructure/database"
	"github.com/product-management/internal/interfaces/http/handler"
	"github.com/product-management/internal/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter(
	cfg *config.Config,
	db *database.Database,
	productService *usecase.ProductUseCase,
	authService *usecase.AuthUseCase,
) *gin.Engine {
	// Set Gin mode
	if cfg.Server.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	r := gin.Default()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// Health check endpoint
	r.GET("/health", handler.HealthCheck)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handler.Login(authService))
			auth.POST("/register", handler.Register(authService))
		}

		// Product routes (protected)
		products := v1.Group("/products")
		products.Use(authMiddleware(authService))
		{
			products.GET("", handler.GetAllProducts(productService))
			products.GET("/:id", handler.GetProduct(productService))
			products.POST("", handler.CreateProduct(productService))
			products.PUT("/:id", handler.UpdateProduct(productService))
			products.DELETE("/:id", handler.DeleteProduct(productService))
			products.PATCH("/:id/stock", handler.UpdateProductStock(productService))
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(authMiddleware(authService))
		{
			users.GET("/profile", handler.GetUserProfile(authService))
			users.PUT("/profile", handler.UpdateUserProfile(authService))
		}
	}

	// Swagger documentation
	if cfg.Server.GinMode != "release" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// authMiddleware validates JWT tokens
func authMiddleware(authService *usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := authHeader[7:]

		// Validate token (this would need to be implemented in the auth service)
		// For now, we'll just check if the token exists
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Add user info to context if needed
		// c.Set("user", user)

		c.Next()
	}
}
