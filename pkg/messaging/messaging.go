package messaging

import (
	"context"

	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaService struct {
	reader    *kafka.Reader
	writer    *kafka.Writer
	consumeCh chan []byte
	produceCh chan []byte
}

type MessageBroker interface {
	Send(msg []byte)
	Receive() <-chan []byte
}

func NewKafkaService(reader *kafka.Reader, writer *kafka.Writer) *KafkaService {
	return &KafkaService{
		reader:    reader,
		writer:    writer,
		consumeCh: make(chan []byte, 10),
		produceCh: make(chan []byte, 10),
	}
}

func (k *KafkaService) Start(ctx context.Context) {
	if k.reader != nil && k.consumeCh != nil {
		go k.consumer(ctx)
	}
	if k.writer != nil && k.produceCh != nil {
		go k.producer(ctx)
	}
}

func (k *KafkaService) Send(msg []byte) {
	k.produceCh <- msg
}

func (k *KafkaService) Receive() <-chan []byte {
	return k.consumeCh
}

func (k *KafkaService) producer(ctx context.Context) {
	logger.Log.Info("producer started!")
	for {
		logger.Log.Info("waiting for msg...")
		msg := <-k.produceCh
		logger.Log.Info("kafka", zap.ByteString("got message from messages channel, sending to kafka", msg))

		err := k.writer.WriteMessages(
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

func (k *KafkaService) consumer(ctx context.Context) {
	logger.Log.Info("consumer started!")
	for {
		m, err := k.reader.ReadMessage(ctx)
		if err != nil {
			logger.Log.Error("read message", zap.Error(err))
			break
		}
		logger.Log.Info("got message from kafka", zap.ByteString("message", m.Value))
		k.consumeCh <- m.Value
	}
	if err := k.reader.Close(); err != nil {
		logger.Log.Error("close reader", zap.Error(err))
	}
}

func CreateWriter(broker string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(broker),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
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

func InitOrderKafka() *KafkaService {
	return NewKafkaService(
		CreateReader("kafka:9092", "order-completed"),
		CreateWriter("kafka:9092", "order-created"),
	)
}

func InitLoyaltyKafka() *KafkaService {
	return NewKafkaService(
		CreateReader("kafka:9092", "order-created"),
		CreateWriter("kafka:9092", "order-completed"),
	)
}

func InitStatusOrder() *KafkaService {
	return &KafkaService{
		reader:    CreateReader("kafka:9092", "order-status"),
		consumeCh: make(chan []byte, 10),
	}
}

func InitStatusLoyalty() *KafkaService {
	return &KafkaService{
		writer:    CreateWriter("kafka:9092", "order-status"),
		produceCh: make(chan []byte, 10),
	}
}
