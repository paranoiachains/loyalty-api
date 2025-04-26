package app

import (
	"github.com/paranoiachains/loyalty-api/pkg/database"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
)

type App struct {
	Kafka *messaging.KafkaService
	DB    database.Storage
}
