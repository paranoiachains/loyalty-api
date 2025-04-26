package app

import (
	"context"

	"github.com/paranoiachains/loyalty-api/pkg/database"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
)

type App struct {
	Kafka       *messaging.KafkaService
	DB          database.Storage
	Processor   MessageProcessor
	StatusKafka *messaging.KafkaService
}

type MessageProcessor interface {
	Process(ctx context.Context)
}
