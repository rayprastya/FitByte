package handlers

import (
	"context"
	"net/http"
	"strconv"

	"fitbyte/internal/entities"
	"fitbyte/internal/services"

	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	activityService services.ActivityService
}

func NewActivityHandler(activityService services.ActivityService) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
	}
}

func (h *ActivityHandler) GetActivities(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, entities.ErrorResponse{
			Success: false,
			Error:   "User ID required",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	limit := h.parseIntQueryParam(c, "limit", 5)
	offset := h.parseIntQueryParam(c, "offset", 0)

	if limit <= 0 {
		limit = 5
	}
	if offset < 0 {
		offset = 0
	}

	filter := services.ActivityServiceFilter{
		Limit:  limit,
		Offset: offset,
	}

	if activityType := c.Query("activityType"); activityType != "" {
		filter.ActivityType = &activityType
	}

	if doneAtFrom := c.Query("doneAtFrom"); doneAtFrom != "" {
		filter.DoneAtFrom = &doneAtFrom
	}

	if doneAtTo := c.Query("doneAtTo"); doneAtTo != "" {
		filter.DoneAtTo = &doneAtTo
	}

	if caloriesBurnedMinStr := c.Query("caloriesBurnedMin"); caloriesBurnedMinStr != "" {
		if caloriesBurnedMin, err := strconv.ParseFloat(caloriesBurnedMinStr, 64); err == nil {
			filter.CaloriesBurnedMin = &caloriesBurnedMin
		}
	}

	if caloriesBurnedMaxStr := c.Query("caloriesBurnedMax"); caloriesBurnedMaxStr != "" {
		if caloriesBurnedMax, err := strconv.ParseFloat(caloriesBurnedMaxStr, 64); err == nil {
			filter.CaloriesBurnedMax = &caloriesBurnedMax
		}
	}

	result, err := h.activityService.GetActivitiesWithFilters(
		context.Background(),
		uint(userID),
		filter,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, entities.ErrorResponse{
			Success: false,
			Error:   "User ID required",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req entities.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	activity, err := h.activityService.CreateActivity(
		context.Background(),
		req,
		uint(userID),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusCreated, entities.APIResponse{
		Success: true,
		Message: "Activity created successfully",
		Data:    activity,
	})
}

func (h *ActivityHandler) GetActivityTypes(c *gin.Context) {
	activityTypes, err := h.activityService.GetActivityTypes(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, entities.APIResponse{
		Success: true,
		Message: "Activity types retrieved successfully",
		Data:    activityTypes,
	})
}

func (h *ActivityHandler) parseIntQueryParam(c *gin.Context, key string, defaultValue int) int {
	if valueStr := c.Query(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil && value >= 0 {
			return value
		}
	}
	return defaultValue
}
