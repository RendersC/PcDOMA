package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type validateTokenResponse struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	Valid  bool   `json:"valid"`
}

type AuthMiddleware struct {
	authServiceURL string
	redisClient    *redis.Client
}

func NewAuthMiddleware(authServiceURL string, redisClient *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		authServiceURL: authServiceURL,
		redisClient:    redisClient,
	}
}

func (m *AuthMiddleware) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]
		resp, err := m.validate(c.Request.Context(), token)
		if err != nil || !resp.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Request.Header.Set("X-User-ID", resp.UserID)
		c.Request.Header.Set("X-User-Role", resp.Role)
		c.Request.Header.Set("X-User-Name", resp.Name)
		c.Next()
	}
}

func (m *AuthMiddleware) Optional() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			if resp, err := m.validate(c.Request.Context(), parts[1]); err == nil && resp.Valid {
				c.Request.Header.Set("X-User-ID", resp.UserID)
				c.Request.Header.Set("X-User-Role", resp.Role)
				c.Request.Header.Set("X-User-Name", resp.Name)
			}
		}
		c.Next()
	}
}

func (m *AuthMiddleware) validate(ctx context.Context, token string) (*validateTokenResponse, error) {
	cacheKey := "token:" + token
	if m.redisClient != nil {
		cached, err := m.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var resp validateTokenResponse
			if json.Unmarshal([]byte(cached), &resp) == nil {
				return &resp, nil
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		m.authServiceURL+"/internal/validate-token", bytes.NewReader(nil))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 3 * time.Second}
	httpResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	body, _ := io.ReadAll(httpResp.Body)
	var resp validateTokenResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("invalid response from auth service")
	}

	if m.redisClient != nil && resp.Valid {
		m.redisClient.Set(ctx, cacheKey, string(body), 60*time.Second)
	}

	return &resp, nil
}
