package domain

type UserAction struct {
	UserID     uint32 `json:"user_id"`
	ProductID  uint32 `json:"product_id"`
	ActionType string `json:"action_type"`
	Name       string `json:"name"`
	Category   string `json:"category"`
}

type UserEvent struct {
	UserID    uint32 `json:"user_id"`
	EventType string `json:"event_type"`
	Email     string `json:"email"`
}

type ProductStat struct {
	ProductID    uint32 `json:"product_id"`
	CartAddCount uint32 `json:"cart_add_count"`
}
