package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pcdoma/auth-service/internal/domain"
	"github.com/pcdoma/auth-service/internal/service"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// List godoc
// @Summary      List all users (admin only)
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        limit  query int false "Limit"  default(20)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} map[string]interface{}
// @Failure      403 {object} map[string]string
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, total, err := h.userSvc.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetByID godoc
// @Summary      Get user by ID
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} domain.User
// @Failure      404 {object} map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.userSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update godoc
// @Summary      Update user profile
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string true "User ID"
// @Param        body body domain.UpdateUserRequest true "Update data"
// @Success      200 {object} domain.User
// @Router       /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	callerID, _ := uuid.Parse(c.GetString("userID"))
	callerRole := c.GetString("userRole")

	user, err := h.userSvc.Update(c.Request.Context(), id, req, callerID, callerRole)
	if err != nil {
		switch err {
		case service.ErrForbidden:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case service.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete godoc
// @Summary      Delete user (admin only)
// @Tags         users
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      204
// @Failure      404 {object} map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.userSvc.Delete(c.Request.Context(), id); err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ChangePassword godoc
// @Summary      Change password
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string true "User ID"
// @Param        body body domain.ChangePasswordRequest true "Passwords"
// @Success      200 {object} map[string]string
// @Router       /users/{id}/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	callerID, _ := uuid.Parse(c.GetString("userID"))
	if callerID != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "can only change your own password"})
		return
	}

	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userSvc.ChangePassword(c.Request.Context(), id, req); err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect old password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}
