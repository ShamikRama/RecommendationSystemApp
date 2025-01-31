package kafka

import "github.com/segmentio/kafka-go"

type KafkaProducer struct {
	writer *kafka.Writer
}

type KafkaMessage struct {
	Id    int                    `json:"id"`
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

func NewKafkaProducer() *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP("kafka:9092"),
			Topic:        "product_updates",
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
			Async:        false,
		},
	}
}
