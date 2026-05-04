package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/pcdoma/auth-service/docs"
	"github.com/pcdoma/auth-service/internal/handler"
	"github.com/pcdoma/auth-service/internal/middleware"
)

func Setup(
	authH *handler.AuthHandler,
	userH *handler.UserHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "auth-service"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Internal (called by gateway — no auth required, internal network only)
	internal := r.Group("/internal")
	{
		internal.POST("/validate-token", authH.ValidateToken)
	}

	// Public auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/refresh", authH.Refresh)
		auth.POST("/logout", authH.Logout)
		auth.GET("/me", middleware.Auth(jwtSecret), authH.Me)
	}

	// User management
	users := r.Group("/users", middleware.Auth(jwtSecret))
	{
		users.GET("", middleware.RequireRole("admin"), userH.List)
		users.GET("/:id", userH.GetByID)
		users.PUT("/:id", userH.Update)
		users.DELETE("/:id", middleware.RequireRole("admin"), userH.Delete)
		users.PUT("/:id/password", userH.ChangePassword)
	}

	return r
}
