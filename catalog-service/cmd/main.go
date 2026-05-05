// @title           PCDoma Catalog Service
// @version         1.0
// @description     PC, Peripheral and Setup catalog for PCDoma rental service
// @host            localhost:8082
// @BasePath        /

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
	"github.com/pcdoma/catalog-service/internal/config"
	"github.com/pcdoma/catalog-service/internal/handler"
	mongorepo "github.com/pcdoma/catalog-service/internal/repository/mongo"
	"github.com/pcdoma/catalog-service/internal/router"
	"github.com/pcdoma/catalog-service/internal/service"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	godotenv.Load()
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
	db := mongoClient.Database(cfg.MongoDB)

	var redisClient *redis.Client
	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err == nil {
		redisClient = redis.NewClient(redisOpt)
		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			log.Printf("warning: redis unavailable, caching disabled: %v", err)
			redisClient = nil
		}
	}

	pcRepo := mongorepo.NewPCRepository(db)
	reviewRepo := mongorepo.NewReviewRepository(db)
	peripheralRepo := mongorepo.NewPeripheralRepository(db)
	setupRepo := mongorepo.NewSetupRepository(db)

	pcSvc := service.NewPCService(pcRepo, reviewRepo, redisClient, cfg.CacheTTL)

	pcH := handler.NewPCHandler(pcSvc)
	peripheralH := handler.NewPeripheralHandler(peripheralRepo)
	setupH := handler.NewSetupHandler(setupRepo, pcRepo, peripheralRepo)

	r := router.Setup(pcH, peripheralH, setupH)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("catalog-service listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down catalog-service...")
	shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutCancel()
	srv.Shutdown(shutCtx)
	mongoClient.Disconnect(shutCtx)
}
