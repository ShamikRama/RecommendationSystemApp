package kafka

type KafkaMessage struct {
	Id    int                    `json:"id"`
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
	Topic string
}
