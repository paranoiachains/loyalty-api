package main

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/order-service/internal/auth"
	"github.com/paranoiachains/loyalty-api/order-service/internal/database"
	"github.com/paranoiachains/loyalty-api/order-service/internal/handlers"
	"github.com/paranoiachains/loyalty-api/order-service/internal/kafka"
	"github.com/paranoiachains/loyalty-api/pkg/flags"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/middleware"

	"go.uber.org/zap"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), middleware.Compression())

	// connect to db only once
	var once sync.Once
	once.Do(func() {
		logger.Log.Debug("DB connection", zap.String("DSN", flags.DatabaseDSN))
		err := database.ConnectToPostgres(flags.DatabaseDSN)
		if err != nil {
			panic(err)
		}
	})

	kafka.StartKafkaServices(kafka.Messages, kafka.Output)

	r.POST("/api/user/register", handlers.Register)
	r.POST("/api/user/login", handlers.Login)

	authGroup := r.Group("/")
	authGroup.Use(auth.Auth())
	{
		authGroup.POST("/api/user/orders", handlers.LoadOrder)
		authGroup.GET("/api/user/orders", handlers.GetOrder)
		authGroup.GET("/api/user/balance", handlers.GetBalance)
		authGroup.POST("/api/user/balance/withdraw", handlers.RequestWithdraw)
		authGroup.GET("/api/user/withdrawals", handlers.Withdrawals)
	}

	r.Run(flags.RunAddress)
}
