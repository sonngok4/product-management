package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/product-management/internal/infrastructure/database"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *database.Database
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// HealthCheck godoc
// @Summary Health check
// @Description Check the health status of the API and its dependencies
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	services := make(map[string]string)
	overallStatus := "healthy"

	// Check database health
	if err := h.db.HealthCheck(); err != nil {
		services["database"] = "unhealthy"
		overallStatus = "unhealthy"
	} else {
		services["database"] = "healthy"
	}

	// Add more service checks here as needed
	services["api"] = "healthy"

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  services,
	}

	if overallStatus == "healthy" {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusServiceUnavailable, response)
	}
}

// HealthCheck is a standalone function for the health check endpoint
func HealthCheck(c *gin.Context) {
	services := make(map[string]string)
	overallStatus := "healthy"

	// For now, we'll just return a basic health check
	// In a real implementation, you might want to inject the database dependency
	services["api"] = "healthy"

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  services,
	}

	c.JSON(http.StatusOK, response)
}

// ReadinessCheck godoc
// @Summary Readiness check
// @Description Check if the API is ready to serve requests
// @Tags health
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 503 {object} ErrorResponse
// @Router /ready [get]
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// Check if database is ready
	if err := h.db.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error:   "Service Unavailable",
			Message: "Database is not ready",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Service is ready",
	})
}

// LivenessCheck godoc
// @Summary Liveness check
// @Description Check if the API is alive
// @Tags health
// @Produce json
// @Success 200 {object} SuccessResponse
// @Router /live [get]
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Service is alive",
	})
}
