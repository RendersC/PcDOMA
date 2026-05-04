// @title           PCDoma Auth Service
// @version         1.0
// @description     Authentication and user management for PCDoma PC rental service
// @termsOfService  http://swagger.io/terms/

// @contact.name   PCDoma Support
// @contact.email  support@pcdoma.kz

// @license.name  MIT

// @host      localhost:8081
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
	"github.com/pcdoma/auth-service/internal/config"
	"github.com/pcdoma/auth-service/internal/handler"
	postgresrepo "github.com/pcdoma/auth-service/internal/repository/postgres"
	"github.com/pcdoma/auth-service/internal/router"
	"github.com/pcdoma/auth-service/internal/service"
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

	userRepo := postgresrepo.NewUserRepository(db)
	tokenRepo := postgresrepo.NewRefreshTokenRepository(db)

	authSvc := service.NewAuthService(userRepo, tokenRepo, cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	userSvc := service.NewUserService(userRepo)

	authH := handler.NewAuthHandler(authSvc, userSvc)
	userH := handler.NewUserHandler(userSvc)

	r := router.Setup(authH, userH, cfg.JWTSecret)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("auth-service listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down auth-service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("auth-service stopped")
}
