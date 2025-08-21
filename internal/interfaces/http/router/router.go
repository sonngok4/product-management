package router

import (
	"github.com/gin-gonic/gin"
	"github.com/product-management/internal/config"
	"github.com/product-management/internal/domain/service"
	"github.com/product-management/internal/infrastructure/database"
	"github.com/product-management/internal/interfaces/http/handler"
	"github.com/product-management/internal/interfaces/http/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter sets up the HTTP router with all routes and middleware
func SetupRouter(
	cfg *config.Config,
	db *database.Database,
	productService service.ProductService,
	authService service.AuthService,
) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	r := gin.New()

	// Global middleware
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(middleware.CORSMiddleware(cfg))

	// Create handlers
	healthHandler := handler.NewHealthHandler(db)
	productHandler := handler.NewProductHandler(productService)
	authHandler := handler.NewAuthHandler(authService)

	// Health check routes (no authentication required)
	r.GET("/health", healthHandler.HealthCheck)
	r.GET("/ready", healthHandler.ReadinessCheck)
	r.GET("/live", healthHandler.LivenessCheck)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Authentication routes (no auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)

			// Authenticated routes
			authProtected := auth.Group("")
			authProtected.Use(middleware.AuthMiddleware(authService))
			{
				authProtected.GET("/profile", authHandler.GetProfile)
				authProtected.PUT("/profile", authHandler.UpdateProfile)
				authProtected.POST("/change-password", authHandler.ChangePassword)
				authProtected.POST("/logout", authHandler.Logout)

				// Admin only routes
				adminOnly := authProtected.Group("")
				adminOnly.Use(middleware.AdminMiddleware())
				{
					adminOnly.GET("/users/:id", authHandler.GetUser)
				}
			}
		}

		// Product routes
		products := v1.Group("/products")
		{
			// Public routes (no auth required)
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.GET("/search", productHandler.SearchProducts)

			// Protected routes (authentication required)
			protected := products.Group("")
			protected.Use(middleware.AuthMiddleware(authService))
			{
				protected.POST("", productHandler.CreateProduct)
				protected.PUT("/:id", productHandler.UpdateProduct)
				protected.DELETE("/:id", productHandler.DeleteProduct)
				protected.PUT("/:id/stock", productHandler.UpdateProductStock)
			}
		}
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Not Found",
			"message": "The requested resource was not found",
			"path":    c.Request.URL.Path,
		})
	})

	return r
}