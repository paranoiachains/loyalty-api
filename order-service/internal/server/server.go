package server

import (
	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/order-service/internal/handlers"
	"github.com/paranoiachains/loyalty-api/order-service/internal/handlers/auth"
	"github.com/paranoiachains/loyalty-api/pkg/app"
	"github.com/paranoiachains/loyalty-api/pkg/middleware"
)

type Server struct {
	engine *gin.Engine
}

func New(a *app.App) *Server {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), middleware.Compression())

	r.POST("/api/user/register", auth.Register(a))
	r.POST("/api/user/login", auth.Login(a))

	authGroup := r.Group("/")
	authGroup.Use(middleware.Auth())
	{
		authGroup.POST("/api/user/orders", handlers.LoadOrder(a))
		authGroup.GET("/api/user/orders", handlers.GetOrders(a))
		authGroup.GET("/api/user/balance", handlers.Balance(a))
		authGroup.POST("/api/user/balance/withdraw", handlers.Withdraw(a))
		authGroup.GET("/api/user/withdrawals", handlers.Withdrawals(a))
	}

	return &Server{engine: r}
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}
