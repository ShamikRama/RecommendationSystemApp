package kafka

import (
	"Recommendation_Service/internal/logger"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Handler interface {
	Handle(msg KafkaMessage) error
}

type kafkaLocal struct {
	reader  *kafka.Reader
	logger  logger.Logger
	handler Handler
}

func NewKafkaConsumer(logger logger.Logger, handler Handler, brokers []string, topic string, groupID string) *kafkaLocal {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 5,
		MaxBytes: 10e6,
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
			r.logger.Error("Error reading message", zap.Error(err))
			continue
		}

		var kafkaMsg KafkaMessage
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
