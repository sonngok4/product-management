// Package main provides the entry point for the Product Management API
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/product-management/internal/config"
	"github.com/product-management/internal/infrastructure/database"
	"github.com/product-management/internal/infrastructure/repository"
	"github.com/product-management/internal/interfaces/http/router"
	"github.com/product-management/internal/usecase"
	"github.com/product-management/pkg/jwt"
	_ "github.com/product-management/docs"
)

// @title Product Management API
// @version 1.0
// @description A REST API for managing products with authentication and authorization
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	// Run migrations
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.GetDB())
	productRepo := repository.NewProductRepository(db.GetDB())

	// Initialize JWT token manager
	expiresIn, err := time.ParseDuration(cfg.JWT.ExpiresIn)
	if err != nil {
		log.Fatalf("Invalid JWT expires in duration: %v", err)
	}
	tokenManager := jwt.NewTokenManager(cfg.JWT.Secret, expiresIn)

	// Initialize use cases
	authService := usecase.NewAuthUseCase(userRepo, tokenManager)
	productService := usecase.NewProductUseCase(productRepo)

	// Setup router
	r := router.SetupRouter(cfg, db, productService, authService)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", cfg.Server.Port)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// createDefaultAdminUser creates a default admin user if it doesn't exist
func createDefaultAdminUser(authService interface{}) {
	// This is a placeholder for creating a default admin user
	// In a real application, you might want to implement this
	log.Println("Default admin user creation not implemented")
}