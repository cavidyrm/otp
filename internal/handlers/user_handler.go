package handlers

import (
	"net/http"

	"otp/internal/models"
	"otp/internal/services"

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

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieve user details by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User ID is required"})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers godoc
// @Summary List users with pagination and search
// @Description Retrieve a paginated list of users with optional search
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Param search query string false "Search by phone number"
// @Success 200 {object} models.UserListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var query models.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid query parameters"})
		return
	}

	// Set default values
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	users, err := h.userService.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Description Delete a user by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "User ID is required"})
		return
	}

	err := h.userService.Delete(c.Request.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "User deleted successfully"})
}

type SuccessResponse struct {
	Message string `json:"message"`
}
