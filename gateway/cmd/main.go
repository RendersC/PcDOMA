package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/pcdoma/gateway/internal/config"
	"github.com/pcdoma/gateway/internal/middleware"
	"github.com/pcdoma/gateway/internal/router"
	"github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("invalid redis url: %v", err)
	}
	redisClient := redis.NewClient(redisOpt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("warning: redis unavailable (%v), token caching disabled", err)
		redisClient = nil
	}

	authMiddleware := middleware.NewAuthMiddleware(cfg.AuthServiceURL, redisClient)

	r := router.Setup(router.Services{
		Auth:         cfg.AuthServiceURL,
		Catalog:      cfg.CatalogServiceURL,
		Booking:      cfg.BookingServiceURL,
		Payment:      cfg.PaymentServiceURL,
		Notification: cfg.NotificationServiceURL,
	}, authMiddleware)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("gateway listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down gateway...")
	shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutCancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("gateway stopped")
}
