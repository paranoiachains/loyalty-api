package main

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/order-service/internal/auth"
	"github.com/paranoiachains/loyalty-api/order-service/internal/handlers"
	"github.com/paranoiachains/loyalty-api/pkg/app"
	"github.com/paranoiachains/loyalty-api/pkg/database"
	"github.com/paranoiachains/loyalty-api/pkg/flags"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
	"github.com/paranoiachains/loyalty-api/pkg/middleware"

	"go.uber.org/zap"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), middleware.Compression())

	var (
		once     sync.Once
		orderApp *app.App
	)

	once.Do(func() {
		logger.Log.Debug("DB connection", zap.String("DSN", flags.DatabaseDSN))
		db, err := database.Connect(flags.DatabaseDSN)
		if err != nil {
			panic(err)
		}

		orderKafka := messaging.InitOrderKafka()

		orderApp = &app.App{
			DB:    db,
			Kafka: orderKafka,
		}
	})

	r.POST("/api/user/register", handlers.Register(orderApp))
	r.POST("/api/user/login", handlers.Login(orderApp))

	authGroup := r.Group("/")
	authGroup.Use(auth.Auth())
	{
		authGroup.POST("/api/user/orders", handlers.LoadOrder(orderApp))
		authGroup.GET("/api/user/orders", handlers.GetOrder)
		authGroup.GET("/api/user/balance", handlers.GetBalance)
		authGroup.POST("/api/user/balance/withdraw", handlers.RequestWithdraw)
		authGroup.GET("/api/user/withdrawals", handlers.Withdrawals)
	}

	r.Run(flags.RunAddress)
}
