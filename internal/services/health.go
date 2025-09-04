// this will be change much
package handlers

import (
	models "fitbyte/internal/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health returns the health status of the API
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "API is healthy",
		Data: gin.H{
			"status":    "ok",
			"service":   "fitbyte-api",
			"version":   "1.0.0",
			"timestamp": gin.H{},
		},
	})
}

// Ready returns the readiness status of the API
func (h *HealthHandler) Ready(c *gin.Context) {
	// Add your readiness checks here (database connection, external services, etc.)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "API is ready",
		Data: gin.H{
			"status": "ready",
		},
	})
}
