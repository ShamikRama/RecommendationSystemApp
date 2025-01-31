package kafka

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Kafka
type Kafka interface {
	Close()
	SendMessage(key string, value interface{}) error
}
