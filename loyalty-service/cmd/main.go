package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/loyalty-service/internal/database"
	"github.com/paranoiachains/loyalty-api/loyalty-service/internal/handlers"
	"github.com/paranoiachains/loyalty-api/loyalty-service/internal/process"
	"github.com/paranoiachains/loyalty-api/pkg/app"
	"github.com/paranoiachains/loyalty-api/pkg/flags"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
	"github.com/paranoiachains/loyalty-api/pkg/middleware"
	"go.uber.org/zap"
)

func main() {
	var loyaltyApp *app.App

	logger.Log.Debug("DSN", zap.String("postgres", flags.LoyaltyDatabaseDSN))
	db, err := database.Connect(flags.LoyaltyDatabaseDSN)
	if err != nil {
		panic(err)
	}

	loyaltyKafka := messaging.InitLoyaltyKafka()
	loyaltyStatus := messaging.InitStatusLoyalty()

	loyaltyApp = &app.App{
		DB:    db,
		Kafka: loyaltyKafka,
		Processor: &process.LoyaltyProcessor{
			DB:           db,
			Broker:       loyaltyKafka,
			StatusBroker: loyaltyStatus,
		},
		StatusKafka: loyaltyStatus,
	}

	loyaltyApp.Kafka.Start(context.Background())
	loyaltyApp.StatusKafka.Start(context.Background())

	go loyaltyApp.Processor.Process(context.Background())

	r := gin.New()
	r.Use(middleware.Logger(), middleware.Compression(), middleware.Auth(), middleware.RateLimitMiddleware())
	r.GET("/api/orders/:number", handlers.GetOrder(loyaltyApp))
	r.Run(flags.AccrualSystemAddress)
}
