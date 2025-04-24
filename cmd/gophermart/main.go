package main

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/internal/database"
	"github.com/paranoiachains/loyalty-api/internal/flags"
	"github.com/paranoiachains/loyalty-api/internal/handlers"
	"github.com/paranoiachains/loyalty-api/internal/middleware"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), middleware.Compression())

	// connect to db only once
	var once sync.Once
	once.Do(func() {
		err := database.ConnectToPostgres(flags.DatabaseURI)
		if err != nil {
			panic(err)
		}
	})

	r.POST("/api/user/register", handlers.Register)
	r.POST("/api/user/login", handlers.Login)

	authGroup := r.Group("/")
	authGroup.Use(middleware.Auth())
	{
		authGroup.POST("/api/user/orders", handlers.LoadOrder)
		authGroup.GET("/api/user/orders", handlers.GetOrder)
		authGroup.GET("/api/user/balance", handlers.GetBalance)
		authGroup.POST("/api/user/balance/withdraw", handlers.RequestWithdraw)
		authGroup.GET("/api/user/withdrawals", handlers.Withdrawals)
	}

	r.Run(flags.RunAddress)
}
