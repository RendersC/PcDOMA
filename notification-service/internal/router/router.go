package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pcdoma/notification-service/internal/handler"
)

func Setup(h *handler.NotificationHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "notification-service"})
	})

	notif := r.Group("/notifications")
	{
		notif.GET("", h.List)
		notif.PATCH("/:id/read", h.MarkRead)
		notif.PATCH("/read-all", h.MarkAllRead)
		notif.DELETE("/:id", h.Delete)
	}

	return r
}
