package grpcapp

import (
	"fmt"
	"net"

	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/sso-service/internal/grpc/auth"
	"github.com/paranoiachains/loyalty-api/sso-service/internal/grpc/withdraw"
	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
	port       int
}

func NewAuth(authService auth.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	auth.Register(gRPCServer, authService)

	return &App{
		gRPCServer: gRPCServer,
		port:       port,
	}

}

func NewWithdraw(withdrawService withdraw.Withdraw, port int) *App {
	gRPCServer := grpc.NewServer()

	withdraw.Register(gRPCServer, withdrawService)

	return &App{
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}

	logger.Log.Info("grpc server started!")

	if err := a.gRPCServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	logger.Log.Info("stopping grpc server")

	a.gRPCServer.GracefulStop()
}
