package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/pcdoma/catalog-service/docs"
	"github.com/pcdoma/catalog-service/internal/handler"
)

func Setup(pcH *handler.PCHandler, peripheralH *handler.PeripheralHandler, setupH *handler.SetupHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "catalog-service"})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Internal (called by booking-service, no auth — internal Docker network)
	internal := r.Group("/internal")
	{
		internal.GET("/pcs/:id/availability", pcH.GetAvailability)
		internal.PATCH("/pcs/:id/status", pcH.UpdateStatus)
	}

	// Public catalog
	catalog := r.Group("/catalog")
	{
		catalog.GET("/pcs", pcH.List)
		catalog.GET("/pcs/:id", pcH.GetByID)
		catalog.GET("/pcs/:id/reviews", pcH.ListReviews)
		catalog.GET("/pcs/:id/available-peripherals", peripheralH.ListByLocation)
		catalog.GET("/locations", pcH.GetLocations)
		catalog.GET("/peripherals", peripheralH.List)
		catalog.GET("/peripherals/:id", peripheralH.GetByID)
		catalog.GET("/setups", setupH.List)
		catalog.GET("/setups/:id", setupH.GetByID)
	}

	// Protected (admin/user via X-User-Role header set by gateway)
	catalog.POST("/pcs/:id/reviews", pcH.CreateReview)
	catalog.DELETE("/reviews/:id", pcH.DeleteReview)

	catalog.POST("/pcs", pcH.Create)
	catalog.PUT("/pcs/:id", pcH.Update)
	catalog.DELETE("/pcs/:id", pcH.Delete)
	catalog.PATCH("/pcs/:id/status", pcH.UpdateStatus)

	catalog.POST("/peripherals", peripheralH.Create)
	catalog.PUT("/peripherals/:id", peripheralH.Update)
	catalog.DELETE("/peripherals/:id", peripheralH.Delete)

	catalog.POST("/setups", setupH.Create)
	catalog.PUT("/setups/:id", setupH.Update)
	catalog.DELETE("/setups/:id", setupH.Delete)

	return r
}
