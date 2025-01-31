package kafka

import (
	"User_Service/internal/logger"
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

func NewKafka(logger logger.Logger) Kafka {
	return &kafkaLocal{
		logger:    logger,
		kafkaProd: NewKafkaProducer(),
	}
}

func (kp *kafkaLocal) SendMessage(key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		kp.logger.Error("Failed to marshal value", zap.Error(err))
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: jsonValue,
		Time:  time.Now(),
	}

	err = kp.kafkaProd.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		kp.logger.Error("Failed to send msg", zap.Error(err))
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (kp *kafkaLocal) Close() {
	if kp.kafkaProd.writer != nil {
		kp.kafkaProd.writer.Close()
	}
}

func CreateTopic(addr string) error {
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %v", err)
	}
	defer conn.Close()

	exists, err := topicExists(conn, "user_updates")
	if err != nil {
		return fmt.Errorf("failed to check if topic exists: %v", err)
	}

	if exists {
		fmt.Printf("Topic %s already exists\n", "user_updates")
		return nil
	}

	topicConfig := kafka.TopicConfig{
		Topic:             "user_updates",
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	err = conn.CreateTopics(topicConfig)
	if err != nil {
		return fmt.Errorf("failed to create topic %s: %v", "user_updates", err)
	}

	fmt.Printf("Topic %s created successfully\n", "user_updates")
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
