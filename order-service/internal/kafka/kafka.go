package kafka

import (
	"context"

	"github.com/paranoiachains/loyalty-api/pkg/messaging"
)

var Messages = make(chan []byte, 10)
var Output = make(chan []byte, 10)

func StartKafkaServices(messages chan []byte, output chan []byte) {
	reader := messaging.CreateReader("kafka:9092", "order-completed")
	writer := messaging.CreateWriter("kafka:9092", "order-placement")

	go messaging.Producer(context.Background(), writer, messages)
	go messaging.Consumer(context.Background(), reader, output)
}
