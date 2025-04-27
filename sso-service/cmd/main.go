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
	application := app.New(5000, time.Hour*1)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	logger.Log.Info("gracefully stopped")
}
