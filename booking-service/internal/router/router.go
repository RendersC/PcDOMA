package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/pcdoma/booking-service/docs"
	"github.com/pcdoma/booking-service/internal/handler"
)

func Setup(h *handler.BookingHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "booking-service"})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public: check slots
	r.GET("/bookings/slots", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"slots": []interface{}{}})
	})

	// User bookings
	bookings := r.Group("/bookings")
	{
		bookings.POST("", h.Create)
		bookings.GET("", h.ListByUser)
		bookings.GET("/:id", h.GetByID)
		bookings.DELETE("/:id", h.Cancel)
	}

	// Worker routes
	worker := r.Group("/worker/bookings")
	{
		worker.GET("", h.WorkerList)
		worker.GET("/active", h.WorkerActiveList)
		worker.GET("/history", h.WorkerHistory)
		worker.PATCH("/:id/accept", h.WorkerAccept)
		worker.PATCH("/:id/reject", h.WorkerReject)
		worker.PATCH("/:id/delivering", h.WorkerDeliver)
		worker.PATCH("/:id/complete", h.WorkerComplete)
	}

	// Admin routes
	admin := r.Group("/admin/bookings")
	{
		admin.GET("", h.AdminList)
	}

	return r
}
