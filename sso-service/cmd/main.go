package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/sso-service/internal/app"
)

func main() {
	auth := app.NewAuth(5000, time.Hour*1)
	withdraw := app.NewWithdraw(5001)

	go auth.GRPCServer.MustRun()
	go withdraw.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	auth.GRPCServer.Stop()
	withdraw.GRPCServer.Stop()
	logger.Log.Info("gracefully stopped")
}
