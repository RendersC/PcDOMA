// @title           PCDoma Payment Service
// @version         1.0
// @description     Mock payment processing and refunds for PCDoma PC rental service
// @termsOfService  http://swagger.io/terms/

// @contact.name   PCDoma Support
// @contact.email  support@pcdoma.kz

// @license.name  MIT

// @host      localhost:8084
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
	"github.com/pcdoma/payment-service/internal/config"
	"github.com/pcdoma/payment-service/internal/handler"
	postgresrepo "github.com/pcdoma/payment-service/internal/repository/postgres"
	"github.com/pcdoma/payment-service/internal/router"
	"github.com/pcdoma/payment-service/internal/service"
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

	repo := postgresrepo.NewPaymentRepository(db)
	svc := service.NewPaymentService(repo)
	h := handler.NewPaymentHandler(svc)
	r := router.Setup(h)

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}

	go func() {
		log.Printf("payment-service listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
