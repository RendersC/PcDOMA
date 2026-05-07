package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pcdoma/gateway/internal/middleware"
	"github.com/pcdoma/gateway/internal/proxy"
)

type Services struct {
	Auth         string
	Catalog      string
	Booking      string
	Payment      string
	Notification string
}

func Setup(svc Services, auth *middleware.AuthMiddleware) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(corsMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "gateway"})
	})

	api := r.Group("/api/v1")

	// ── Auth (public) ──────────────────────────────────────────
	api.POST("/auth/register", proxy.To(svc.Auth+"/auth/register"))
	api.POST("/auth/login", proxy.To(svc.Auth+"/auth/login"))
	api.POST("/auth/refresh", proxy.To(svc.Auth+"/auth/refresh"))
	api.POST("/auth/logout", proxy.To(svc.Auth+"/auth/logout"))
	api.GET("/auth/me", auth.Required(), proxy.To(svc.Auth+"/auth/me"))

	// ── Users (protected) ─────────────────────────────────────
	users := api.Group("/users", auth.Required())
	{
		users.GET("", proxy.To(svc.Auth+"/users"))
		users.GET("/:id", proxy.To(svc.Auth+"/users/:id"))
		users.PUT("/:id", proxy.To(svc.Auth+"/users/:id"))
		users.DELETE("/:id", proxy.To(svc.Auth+"/users/:id"))
		users.PUT("/:id/password", proxy.To(svc.Auth+"/users/:id/password"))
	}

	// ── Catalog (public + protected) ─────────────────────────
	// Public catalog routes
	api.GET("/catalog/pcs", auth.Optional(), proxy.To(svc.Catalog+"/catalog/pcs"))
	api.GET("/catalog/pcs/:id", auth.Optional(), proxy.To(svc.Catalog+"/catalog/pcs/:id"))
	api.GET("/catalog/pcs/:id/reviews", proxy.To(svc.Catalog+"/catalog/pcs/:id/reviews"))
	api.GET("/catalog/pcs/:id/available-peripherals", proxy.To(svc.Catalog+"/catalog/pcs/:id/available-peripherals"))
	api.GET("/catalog/peripherals", proxy.To(svc.Catalog+"/catalog/peripherals"))
	api.GET("/catalog/peripherals/:id", proxy.To(svc.Catalog+"/catalog/peripherals/:id"))
	api.GET("/catalog/setups", proxy.To(svc.Catalog+"/catalog/setups"))
	api.GET("/catalog/setups/:id", proxy.To(svc.Catalog+"/catalog/setups/:id"))
	api.GET("/catalog/locations", proxy.To(svc.Catalog+"/catalog/locations"))

	// Protected catalog routes (authenticated)
	api.POST("/catalog/pcs/:id/reviews", auth.Required(), proxy.To(svc.Catalog+"/catalog/pcs/:id/reviews"))
	api.DELETE("/catalog/reviews/:id", auth.Required(), proxy.To(svc.Catalog+"/catalog/reviews/:id"))

	// Admin catalog routes
	catalog := api.Group("/catalog", auth.Required())
	{
		catalog.POST("/pcs", proxy.To(svc.Catalog+"/catalog/pcs"))
		catalog.PUT("/pcs/:id", proxy.To(svc.Catalog+"/catalog/pcs/:id"))
		catalog.DELETE("/pcs/:id", proxy.To(svc.Catalog+"/catalog/pcs/:id"))
		catalog.PATCH("/pcs/:id/status", proxy.To(svc.Catalog+"/catalog/pcs/:id/status"))
		catalog.POST("/peripherals", proxy.To(svc.Catalog+"/catalog/peripherals"))
		catalog.PUT("/peripherals/:id", proxy.To(svc.Catalog+"/catalog/peripherals/:id"))
		catalog.DELETE("/peripherals/:id", proxy.To(svc.Catalog+"/catalog/peripherals/:id"))
		catalog.POST("/setups", proxy.To(svc.Catalog+"/catalog/setups"))
		catalog.PUT("/setups/:id", proxy.To(svc.Catalog+"/catalog/setups/:id"))
		catalog.DELETE("/setups/:id", proxy.To(svc.Catalog+"/catalog/setups/:id"))
	}

	// ── Bookings ──────────────────────────────────────────────
	// Public slots check
	api.GET("/bookings/slots", proxy.To(svc.Booking+"/bookings/slots"))
	api.GET("/bookings/peripheral-slots", proxy.To(svc.Booking+"/bookings/peripheral-slots"))

	bookings := api.Group("/bookings", auth.Required())
	{
		bookings.POST("", proxy.To(svc.Booking+"/bookings"))
		bookings.GET("", proxy.To(svc.Booking+"/bookings"))
		bookings.GET("/:id", proxy.To(svc.Booking+"/bookings/:id"))
		bookings.DELETE("/:id", proxy.To(svc.Booking+"/bookings/:id"))
	}

	// Worker routes
	worker := api.Group("/worker", auth.Required())
	{
		worker.GET("/bookings", proxy.To(svc.Booking+"/worker/bookings"))
		worker.GET("/bookings/active", proxy.To(svc.Booking+"/worker/bookings/active"))
		worker.GET("/bookings/history", proxy.To(svc.Booking+"/worker/bookings/history"))
		worker.PATCH("/bookings/:id/accept", proxy.To(svc.Booking+"/worker/bookings/:id/accept"))
		worker.PATCH("/bookings/:id/reject", proxy.To(svc.Booking+"/worker/bookings/:id/reject"))
		worker.PATCH("/bookings/:id/delivering", proxy.To(svc.Booking+"/worker/bookings/:id/delivering"))
		worker.PATCH("/bookings/:id/complete", proxy.To(svc.Booking+"/worker/bookings/:id/complete"))
	}

	// Admin bookings
	adminBookings := api.Group("/admin/bookings", auth.Required())
	{
		adminBookings.GET("", proxy.To(svc.Booking+"/admin/bookings"))
		adminBookings.GET("/stats", proxy.To(svc.Booking+"/admin/bookings/stats"))
		adminBookings.PATCH("/:id", proxy.To(svc.Booking+"/admin/bookings/:id"))
	}

	// ── Payments ──────────────────────────────────────────────
	payments := api.Group("/payments", auth.Required())
	{
		payments.POST("", proxy.To(svc.Payment+"/payments"))
		payments.GET("", proxy.To(svc.Payment+"/payments"))
		payments.GET("/:id", proxy.To(svc.Payment+"/payments/:id"))
		payments.POST("/:id/process", proxy.To(svc.Payment+"/payments/:id/process"))
		payments.POST("/:id/refund", proxy.To(svc.Payment+"/payments/:id/refund"))
	}

	adminPayments := api.Group("/admin/payments", auth.Required())
	{
		adminPayments.GET("", proxy.To(svc.Payment+"/admin/payments"))
		adminPayments.GET("/stats", proxy.To(svc.Payment+"/admin/payments/stats"))
	}

	// ── Notifications ─────────────────────────────────────────
	notif := api.Group("/notifications", auth.Required())
	{
		notif.GET("", proxy.To(svc.Notification+"/notifications"))
		notif.PATCH("/:id/read", proxy.To(svc.Notification+"/notifications/:id/read"))
		notif.PATCH("/read-all", proxy.To(svc.Notification+"/notifications/read-all"))
		notif.DELETE("/:id", proxy.To(svc.Notification+"/notifications/:id"))
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-User-Role, X-User-Name")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
