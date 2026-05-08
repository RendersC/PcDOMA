// @title           PCDoma Booking Service
// @version         1.0
// @description     Booking orchestration (catalog + payment + worker workflow) for PCDoma PC rental service
// @termsOfService  http://swagger.io/terms/

// @contact.name   PCDoma Support
// @contact.email  support@pcdoma.kz

// @license.name  MIT

// @host      localhost:8083
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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
	"github.com/pcdoma/booking-service/internal/clients"
	"github.com/pcdoma/booking-service/internal/config"
	"github.com/pcdoma/booking-service/internal/events"
	"github.com/pcdoma/booking-service/internal/handler"
	postgresrepo "github.com/pcdoma/booking-service/internal/repository/postgres"
	"github.com/pcdoma/booking-service/internal/router"
	"github.com/pcdoma/booking-service/internal/service"
	"github.com/redis/go-redis/v9"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	db, err := gorm.Open(gormpostgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	var redisClient *redis.Client
	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err == nil {
		redisClient = redis.NewClient(redisOpt)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := redisClient.Ping(ctx).Err(); err != nil {
			log.Printf("warning: redis unavailable: %v", err)
			redisClient = nil
		}
	}

	repo := postgresrepo.NewBookingRepository(db)
	catalogClient := clients.NewCatalogClient(cfg.CatalogServiceURL)
	paymentClient := clients.NewPaymentClient(cfg.PaymentServiceURL)
	publisher := events.NewPublisher(redisClient)

	svc := service.NewBookingService(repo, catalogClient, paymentClient, publisher)
	h := handler.NewBookingHandler(svc)
	r := router.Setup(h)

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}

	go func() {
		log.Printf("booking-service listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutCancel()
	srv.Shutdown(shutCtx)
}
