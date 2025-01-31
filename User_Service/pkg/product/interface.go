package product

import (
	"context"
)

type Client interface {
	GetProducts(ctx context.Context, pageSize int, page int) (ProductsResponse, error)
	CreateCart(ctx context.Context, request CartCreateRequest) error
	UpdateCart(ctx context.Context, request CartUpdateRequest) error
	DeleteCart(ctx context.Context, request CartDeleteRequest) error
}
