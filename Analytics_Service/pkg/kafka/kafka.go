package kafka

import (
	"Analytics_Service/internal/logger"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Handler
type Handler interface {
	Handle(msg KafkaMessage) error
}

type kafkaLocal struct {
	reader  *kafka.Reader
	logger  logger.Logger
	handler Handler
}

func NewKafkaConsumer(logger logger.Logger, handler Handler, brokers []string, topics []string, groupID string) *kafkaLocal {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupTopics: topics,
		GroupID:     groupID,
		MinBytes:    5,
		MaxBytes:    10e6,
	})

	return &kafkaLocal{
		reader:  reader,
		logger:  logger,
		handler: handler,
	}
}

func (r *kafkaLocal) ReadMessages() {
	for {
		msg, err := r.reader.ReadMessage(context.Background())
		if err != nil {
			if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Temporary() {
				r.logger.Error("Temporary Kafka error, retrying...", zap.Error(err))
				time.Sleep(5 * time.Second)
				continue
			}
			r.logger.Error("Fatal Kafka error", zap.Error(err))
			continue
		}

		var kafkaMsg KafkaMessage
		kafkaMsg.Topic = msg.Topic
		if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
			r.logger.Error("Error unmarshalling message", zap.Error(err))
			continue
		}

		r.logger.Info("Received Kafka message",
			zap.Int("id", kafkaMsg.Id),
			zap.String("event", kafkaMsg.Event),
			zap.Any("data", kafkaMsg.Data),
		)

		if err := r.handler.Handle(kafkaMsg); err != nil {
			r.logger.Error("Error handling message", zap.Error(err))
			continue
		}
	}
}

func (kc *kafkaLocal) Close() {
	if err := kc.reader.Close(); err != nil {
		log.Printf("Ошибка закрытия Kafka Reader: %v\n", err)
	}
}
