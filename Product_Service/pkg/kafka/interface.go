package kafka

//go:generate mockery --name=Kafka --dir=./pkg/kafka --output=./pkg/kafka/mocks --outpkg=mocks
type Kafka interface {
	SendMessage(key string, value interface{}) error
	Close()
}
