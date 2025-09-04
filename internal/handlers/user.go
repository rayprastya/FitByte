package handlers

import (
	"context"
	"net/http"
	"strconv"

	"fitbyte/internal/entities"
	"fitbyte/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	limit := h.parseIntQueryParam(c, "limit", 5)
	offset := h.parseIntQueryParam(c, "offset", 0)

	if limit <= 0 {
		limit = 5
	}
	if offset < 0 {
		offset = 0
	}

	filter := services.UserServiceFilter{
		Limit:  limit,
		Offset: offset,
	}

	if isActiveStr := c.Query("isActive"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filter.IsActive = &isActive
		}
	}

	result, err := h.userService.GetUsersWithFilters(context.Background(), filter)
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

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, err := h.userService.GetUser(context.Background(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, entities.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, entities.APIResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req entities.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, err := h.userService.CreateUser(context.Background(), req)
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
		Message: "User created successfully",
		Data:    user,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req entities.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, err := h.userService.UpdateUser(context.Background(), uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, entities.APIResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = h.userService.DeleteUser(context.Background(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, entities.APIResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}

func (h *UserHandler) parseIntQueryParam(c *gin.Context, key string, defaultValue int) int {
	if valueStr := c.Query(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil && value >= 0 {
			return value
		}
	}
	return defaultValue
}