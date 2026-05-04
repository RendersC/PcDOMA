package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pcdoma/auth-service/internal/domain"
	"github.com/pcdoma/auth-service/internal/service"
)

type AuthHandler struct {
	authSvc service.AuthService
	userSvc service.UserService
}

func NewAuthHandler(authSvc service.AuthService, userSvc service.UserService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, userSvc: userSvc}
}

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body domain.RegisterRequest true "Registration data"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      409 {object} map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, tokens, err := h.authSvc.Register(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrEmailTaken {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

// Login godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body domain.LoginRequest true "Credentials"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, tokens, err := h.authSvc.Login(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

// Refresh godoc
// @Summary      Refresh access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body domain.RefreshRequest true "Refresh token"
// @Success      200 {object} domain.TokenPair
// @Failure      401 {object} map[string]string
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req domain.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authSvc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Logout godoc
// @Summary      Logout (invalidate refresh token)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body domain.RefreshRequest true "Refresh token"
// @Success      200 {object} map[string]string
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req domain.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.authSvc.Logout(c.Request.Context(), req.RefreshToken)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// Me godoc
// @Summary      Get current user profile
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} domain.User
// @Failure      401 {object} map[string]string
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := h.userSvc.GetByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ValidateToken godoc (internal endpoint, called by gateway)
// @Summary      Validate JWT token (internal)
// @Tags         internal
// @Accept       json
// @Produce      json
// @Success      200 {object} domain.ValidateTokenResponse
// @Router       /internal/validate-token [post]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusOK, domain.ValidateTokenResponse{Valid: false})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusOK, domain.ValidateTokenResponse{Valid: false})
		return
	}

	resp, err := h.authSvc.ValidateToken(c.Request.Context(), parts[1])
	if err != nil {
		c.JSON(http.StatusOK, domain.ValidateTokenResponse{Valid: false})
		return
	}

	c.JSON(http.StatusOK, resp)
}
