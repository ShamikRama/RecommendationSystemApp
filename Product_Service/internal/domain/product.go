package domain

type Product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

type CartDTO struct {
	UserId    int    `json:"user_id" binding:"required"`
	ProductId int    `json:"product_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
	Category  string `json:"category" binding:"required"`
}

type Cart struct {
	Id        int
	UserId    int
	ProductId int
	Quantity  int
}

type CartDeleteDTO struct {
	UserId    int `json:"user_id" binding:"required"`
	ProductId int `json:"product_id" binding:"required"`
}

type ProductsResponse struct {
	Data  []Product `json:"data"`
	Total int       `json:"total"`
}
