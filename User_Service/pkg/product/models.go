package product

const ShortClient = "/products"

type Product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
}

type CartCreateRequest struct {
	UserId    int    `json:"user_id" binding:"required"`
	ProductId int    `json:"product_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
	Category  string `json:"category" binding:"required"`
}

type CartDeleteRequest struct {
	UserId    int `json:"user_id" binding:"required"`
	ProductId int `json:"product_id" binding:"required"`
}

type CartUpdateRequest struct {
	UserId    int `json:"user_id" binding:"required"`
	ProductId int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

type ProductsResponse struct {
	Data  []Product `json:"data"`
	Total int       `json:"total"`
}
