package main

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/loyalty-service/internal/process"
	"github.com/paranoiachains/loyalty-api/pkg/app"
	"github.com/paranoiachains/loyalty-api/pkg/database"
	"github.com/paranoiachains/loyalty-api/pkg/flags"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
)

func main() {
	var (
		once       sync.Once
		loyaltyApp app.App
	)

	once.Do(func() {
		logger.Log.Info("connecting to db...")
		db, err := database.Connect(flags.DatabaseDSN)
		if err != nil {
			panic(err)
		}
		logger.Log.Info("successfully connected to db!")

		loyaltyKafka := messaging.InitLoyaltyKafka()

		loyaltyApp = app.App{
			DB:    db,
			Kafka: loyaltyKafka,
		}
	})

	go process.Process(context.Background(), &loyaltyApp)

	r := gin.New()
	r.Run(flags.AccrualSystemAddress)
}
