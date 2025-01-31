package domain

type Product struct {
	ProductId int    `json:"product_id"`
	Name      string `json:"name"`
	Category  string `json:"category"`
}

type KafkaMessage struct {
	Id    int                    `json:"id"`
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}
