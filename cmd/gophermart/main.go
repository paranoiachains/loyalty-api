package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/internal/flags"
	"github.com/paranoiachains/loyalty-api/internal/handlers"
	"github.com/paranoiachains/loyalty-api/internal/middleware"
)

func main() {
	router := gin.New()
	router.Use(gin.Recovery(), middleware.Logger(), middleware.Compression())

	group := router.Group("/api/user/")
	group.POST("register", handlers.Register)
	group.POST("login", handlers.Login)
	group.POST("orders", handlers.LoadOrder)
	group.GET("orders", handlers.GetOrder)
	group.GET("balance", handlers.GetBalance)
	group.POST("balance/withdraw", handlers.RequestWithdraw)
	group.GET("withdrawals", handlers.Withdrawals)

	router.Run(flags.RunAddress)
}
