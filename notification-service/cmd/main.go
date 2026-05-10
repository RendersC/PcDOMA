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
	"github.com/pcdoma/notification-service/internal/config"
	"github.com/pcdoma/notification-service/internal/events"
	"github.com/pcdoma/notification-service/internal/handler"
	mongorepo "github.com/pcdoma/notification-service/internal/repository/mongo"
	"github.com/pcdoma/notification-service/internal/router"
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
		log.Fatalf("failed to connect MongoDB: %v", err)
	}
	db := mongoClient.Database(cfg.MongoDB)
	repo := mongorepo.NewNotificationRepository(db)

	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("invalid redis URL: %v", err)
	}
	redisClient := redis.NewClient(redisOpt)

	mainCtx, mainCancel := context.WithCancel(context.Background())
	subscriber := events.NewSubscriber(redisClient, repo)
	subscriber.Start(mainCtx)

	h := handler.NewNotificationHandler(repo)
	r := router.Setup(h)

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}

	go func() {
		log.Printf("notification-service listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	mainCancel()
	shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutCancel()
	srv.Shutdown(shutCtx)
	mongoClient.Disconnect(shutCtx)
}
