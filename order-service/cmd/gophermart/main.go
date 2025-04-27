package main

import (
	"context"

	"github.com/paranoiachains/loyalty-api/order-service/internal/app"
	"github.com/paranoiachains/loyalty-api/order-service/internal/server"

	"github.com/paranoiachains/loyalty-api/pkg/flags"
)

func main() {
	ctx := context.Background()

	application, err := app.New(ctx)
	if err != nil {
		panic(err)
	}

	go application.Processor.Process(ctx)

	srv := server.New(application)
	srv.Run(flags.RunAddress)
}
