package server

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, h *Handlers) {
	v1 := router.Group("/api/v1")
	{
		health := v1.Group("/health")
		{
			health.GET("/", h.HealthHandler.Health)
			health.GET("/ready", h.HealthHandler.Ready)
		}

		users := v1.Group("/users")
		{
			users.GET("/", h.UserHandler.GetUsers)
			users.GET("/:id", h.UserHandler.GetUser)
			users.POST("/", h.UserHandler.CreateUser)
			users.PUT("/:id", h.UserHandler.UpdateUser)
			users.DELETE("/:id", h.UserHandler.DeleteUser)
		}

		activities := v1.Group("/activity")
		{
			activities.GET("/", h.ActivityHandler.GetActivities)
			activities.POST("/", h.ActivityHandler.CreateActivity)
		}

		activityTypes := v1.Group("/activity-types")
		{
			activityTypes.GET("/", h.ActivityHandler.GetActivityTypes)
		}
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to FitByte API",
			"version": "1.0.0",
			"docs":    "/api/v1/health",
		})
	})
}