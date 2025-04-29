package app

import (
	"context"

	"github.com/paranoiachains/loyalty-api/order-service/internal/database"
	"github.com/paranoiachains/loyalty-api/order-service/internal/process"
	"github.com/paranoiachains/loyalty-api/pkg/app"
	auth "github.com/paranoiachains/loyalty-api/pkg/clients/sso/auth"
	withdraw "github.com/paranoiachains/loyalty-api/pkg/clients/sso/withdraw"
	"github.com/paranoiachains/loyalty-api/pkg/flags"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
	"go.uber.org/zap"
)

func New(ctx context.Context) (*app.App, error) {
	logger.Log.Debug("Connecting to database", zap.String("dsn", flags.DatabaseDSN))
	db, err := database.Connect(flags.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	authClient, err := auth.New("sso_service:5000")
	if err != nil {
		return nil, err
	}

	withdrawClient, err := withdraw.New("sso_service:5001")
	if err != nil {
		return nil, err
	}

	orderKafka := messaging.InitOrderKafka()
	statusKafka := messaging.InitStatusOrder()

	orderKafka.Start(ctx)
	statusKafka.Start(ctx)

	return &app.App{
		DB:          db,
		Kafka:       orderKafka,
		StatusKafka: statusKafka,
		Processor: process.OrderProcessor{
			DB:             db,
			Broker:         orderKafka,
			StatusBroker:   statusKafka,
			WithdrawClient: withdrawClient,
		},
		AuthClient:     authClient,
		WithdrawClient: withdrawClient,
	}, nil
}
