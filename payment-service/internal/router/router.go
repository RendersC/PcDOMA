package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/pcdoma/payment-service/docs"
	"github.com/pcdoma/payment-service/internal/handler"
)

func Setup(h *handler.PaymentHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "payment-service"})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Internal (called by booking-service)
	r.POST("/internal/payments", h.InternalCreate)

	payments := r.Group("/payments")
	{
		payments.GET("", h.List)
		payments.GET("/:id", h.GetByID)
		payments.POST("/:id/process", h.Process)
		payments.POST("/:id/refund", h.Refund)
	}

	admin := r.Group("/admin/payments")
	{
		admin.GET("", h.AdminList)
	}

	return r
}
