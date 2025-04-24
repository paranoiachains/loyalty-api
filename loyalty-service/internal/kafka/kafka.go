package kafka

import (
	"context"

	"github.com/paranoiachains/loyalty-api/pkg/messaging"
)

var Input = make(chan []byte, 10)
var Processed = make(chan []byte, 10)

func StartKafkaServices(input chan []byte, processed chan []byte) {
	reader := messaging.CreateReader("kafka:9092", "order-placement")
	writer := messaging.CreateWriter("kafka:9092", "order-completed")

	go messaging.Consumer(context.Background(), reader, input)
	go messaging.Producer(context.Background(), writer, processed)
}
