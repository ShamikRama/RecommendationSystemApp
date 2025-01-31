package kafka

import (
	"Product_Service/internal/logger"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type kafkaLocal struct {
	logger    logger.Logger
	kafkaProd *KafkaProducer
}

func NewKafkaLocal(logger logger.Logger) Kafka {
	return &kafkaLocal{
		logger:    logger,
		kafkaProd: NewKafkaProducer(),
	}
}

func (r *kafkaLocal) SendMessage(key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		r.logger.Error("Error marshaling the value", zap.Error(err))
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: jsonValue,
		Time:  time.Now(),
	}

	err = r.kafkaProd.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		r.logger.Error("Error writing the message", zap.Error(err))
		return err
	}

	return nil
}

func (r *kafkaLocal) Close() {
	if r.kafkaProd.writer != nil {
		r.kafkaProd.writer.Close()
	}
}

func CreateTopic() error {
	conn, err := kafka.Dial("tcp", "kafka:9092")
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %v", err)
	}
	defer conn.Close()

	// Проверка существования топика
	exists, err := topicExists(conn, "product_updates")
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %v", err)
	}

	if exists {
		fmt.Printf("Topic %s already exists\n", "product_updates")
		return nil
	}

	// Создание топика
	topicConfig := kafka.TopicConfig{
		Topic:             "product_updates",
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	err = conn.CreateTopics(topicConfig)
	if err != nil {
		return fmt.Errorf("failed to create topic %s: %v", "product_updates", err)
	}

	fmt.Printf("Topic %s created successfully\n", "product_updates")
	return nil
}

func topicExists(conn *kafka.Conn, topic string) (bool, error) {
	partitions, err := conn.ReadPartitions()
	if err != nil {
		return false, fmt.Errorf("failed to read partitions: %v", err)
	}

	for _, p := range partitions {
		if p.Topic == topic {
			return true, nil
		}
	}

	return false, nil
}
