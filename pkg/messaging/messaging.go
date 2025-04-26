package messaging

import (
	"context"
	"time"

	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaService struct {
	reader    *kafka.Reader
	writer    *kafka.Writer
	input     chan []byte
	processed chan []byte
}

// functional options parameter
type Option func(*KafkaService)

func WithReader(reader *kafka.Reader) Option {
	return func(k *KafkaService) {
		k.reader = reader
	}
}

func WithWriter(writer *kafka.Writer) Option {
	return func(k *KafkaService) {
		k.writer = writer
	}
}

func WithInputChannel(ch chan []byte) Option {
	return func(k *KafkaService) {
		k.input = ch
	}
}

func WithProcessedChannel(ch chan []byte) Option {
	return func(k *KafkaService) {
		k.processed = ch
	}
}

func NewKafkaService(opts ...Option) *KafkaService {
	k := &KafkaService{}
	for _, opt := range opts {
		opt(k)
	}
	return k
}

func (k *KafkaService) Start(ctx context.Context) {
	if k.reader != nil && k.processed != nil {
		go k.consumer(ctx)
		logger.Log.Info("consumer started!")
	}
	if k.writer != nil && k.input != nil {
		go k.producer(ctx)
		logger.Log.Info("producer started!")
	}
}

func (k *KafkaService) Send(msg []byte) {
	k.processed <- msg
}

func (k *KafkaService) Receive() <-chan []byte {
	return k.input
}

func (k *KafkaService) producer(ctx context.Context) {
	for {
		msg := <-k.input
		logger.Log.Info("kafka", zap.ByteString("got message from messages channel, sending to kafka", msg))

		err := k.writer.WriteMessages(
			ctx,
			kafka.Message{
				Value: msg,
			},
		)
		if err != nil {
			logger.Log.Error("send message to kafka", zap.Error(err))
		} else {
			logger.Log.Info("sent message to kafka")
		}
	}
}

func (k *KafkaService) consumer(ctx context.Context) {
	time.Sleep(15 * time.Second)
	for {
		m, err := k.reader.ReadMessage(ctx)
		if err != nil {
			logger.Log.Error("read message", zap.Error(err))
			break
		}
		logger.Log.Info("got message from kafka", zap.ByteString("message", m.Value))
		k.processed <- m.Value
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

var OrderKafka *KafkaService

func InitOrderKafka() *KafkaService {
	writer := CreateWriter("kafka:9092", "order-created")
	reader := CreateReader("kafka:9092", "order-completed")

	return NewKafkaService(
		WithReader(reader),
		WithWriter(writer),
		WithInputChannel(make(chan []byte, 1)),
		WithProcessedChannel(make(chan []byte, 1)),
	)
}

func InitLoyaltyKafka() *KafkaService {
	writer := CreateWriter("kafka:9092", "order-completed")
	reader := CreateReader("kafka:9092", "order-created")

	return NewKafkaService(
		WithReader(reader),
		WithWriter(writer),
		WithInputChannel(make(chan []byte, 1)),
		WithProcessedChannel(make(chan []byte, 1)),
	)
}
