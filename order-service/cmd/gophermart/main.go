package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/order-service/internal/auth"
	"github.com/paranoiachains/loyalty-api/order-service/internal/database"
	"github.com/paranoiachains/loyalty-api/order-service/internal/handlers"
	"github.com/paranoiachains/loyalty-api/order-service/internal/process"
	"github.com/paranoiachains/loyalty-api/pkg/app"

	"github.com/paranoiachains/loyalty-api/pkg/flags"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
	"github.com/paranoiachains/loyalty-api/pkg/middleware"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), middleware.Compression())

	var orderApp *app.App

	db, err := database.Connect(flags.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	orderKafka := messaging.InitOrderKafka()
	orderStatusKafka := messaging.InitStatusOrder()

	orderApp = &app.App{
		DB:    db,
		Kafka: orderKafka,
		Processor: process.OrderProcessor{
			DB:           db,
			Broker:       orderKafka,
			StatusBroker: orderStatusKafka,
		},
		StatusKafka: orderStatusKafka,
	}

	orderApp.Kafka.Start(context.Background())
	orderApp.StatusKafka.Start(context.Background())

	go orderApp.Processor.Process(context.Background())

	r.POST("/api/user/register", handlers.Register(orderApp))
	r.POST("/api/user/login", handlers.Login(orderApp))

	authGroup := r.Group("/")
	authGroup.Use(auth.Auth())
	{
		authGroup.POST("/api/user/orders", handlers.LoadOrder(orderApp))
		authGroup.GET("/api/user/orders", handlers.GetOrders(orderApp))
		authGroup.GET("/api/user/balance", handlers.GetBalance)
		authGroup.POST("/api/user/balance/withdraw", handlers.RequestWithdraw)
		authGroup.GET("/api/user/withdrawals", handlers.Withdrawals)
	}

	r.Run(flags.RunAddress)
}
