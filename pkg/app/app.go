package app

import (
	"context"

	ssoauth "github.com/paranoiachains/loyalty-api/pkg/clients/sso/auth"
	ssowithdraw "github.com/paranoiachains/loyalty-api/pkg/clients/sso/withdraw"
	"github.com/paranoiachains/loyalty-api/pkg/database"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
)

type App struct {
	Kafka          *messaging.KafkaService
	DB             database.Storage
	Processor      MessageProcessor
	StatusKafka    *messaging.KafkaService
	AuthClient     *ssoauth.AuthClient
	WithdrawClient *ssowithdraw.WithdrawalsClient
}

type MessageProcessor interface {
	Process(ctx context.Context)
}
