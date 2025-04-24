package messaging

import (
	"context"
	"time"

	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func CreateWriter(broker string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(broker),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}
}

func Producer(ctx context.Context, writer *kafka.Writer, messages <-chan []byte) {
	for {
		// buffered channel, wait until message comes
		msg := <-messages
		logger.Log.Info("kafka", zap.ByteString("got message from messages channel, sending to kafka", msg))

		err := writer.WriteMessages(
			ctx,
			kafka.Message{
				Value: msg,
			},
		)
		if err != nil {
			logger.Log.Error("send message to kafka", zap.Error(err))
		}
		logger.Log.Info("sent message to kafka")
	}
}

func CreateReader(broker string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e6,
	})
}

func Consumer(ctx context.Context, reader *kafka.Reader, output chan<- []byte) {
	time.Sleep(15 * time.Second)
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			logger.Log.Error("read message", zap.Error(err))
			break
		}
		logger.Log.Info("got message from kafka", zap.ByteString("message", m.Value))
		output <- m.Value
	}
	if err := reader.Close(); err != nil {
		logger.Log.Error("close reader", zap.Error(err))
	}
}
