package service

import (
	"Product_Service/internal/domain"
	"Product_Service/internal/erorrs"
	"Product_Service/internal/logger"
	"Product_Service/internal/repository"
	"Product_Service/pkg/kafka"
	"context"
	"errors"
	"strconv"

	"go.uber.org/zap"
)

type ServiceProduct interface {
	CreateCart(ctx context.Context, dto *domain.CartDTO) error
	UpdateCartItem(ctx context.Context, dto *domain.CartDTO) error
	DeleteCartItem(ctx context.Context, userID int, productId int) error
	GetAllProducts(ctx context.Context, pageSize int, page int) (domain.ProductsResponse, error)
}

type serviceProduct struct {
	repo   repository.RepositoryProduct
	logger *logger.Logger
	kafka  kafka.Kafka
}

func NewProductService(repo repository.RepositoryProduct, logger *logger.Logger, kafka kafka.Kafka) *serviceProduct {
	return &serviceProduct{
		repo:   repo,
		logger: logger,
		kafka:  kafka,
	}
}

func (s *serviceProduct) CreateCart(ctx context.Context, dto *domain.CartDTO) error {

	cartItem := &domain.Cart{
		ProductId: dto.ProductId,
		UserId:    dto.UserId,
		Quantity:  dto.Quantity,
	}

	if cartItem.ProductId == 0 || cartItem.UserId == 0 || cartItem.Quantity == 0 {
		s.logger.Error("Empty input",
			zap.String("operation", "service.ProductService.CreateCart"))
		return erorrs.ErrEmptyInput
	}

	err := s.repo.CreateCartItem(ctx, cartItem)
	if err != nil {
		s.logger.Error("Error creating cart",
			zap.String("operation", "service.ProductService.CreateCart"),
			zap.Error(err))
		return err
	}
	s.logger.Info("Cart created")

	msg := kafka.KafkaMessage{
		Id:    dto.UserId,
		Event: "cart_create",
		Data: map[string]interface{}{
			"UserId":    dto.UserId,
			"ProductId": dto.ProductId,
			"Name":      dto.Name,
			"Category":  dto.Category,
		},
	}

	key := strconv.Itoa(dto.UserId)

	err = s.kafka.SendMessage(key, msg)
	if err != nil {
		s.logger.Error("Error sending the message to kafka",
			zap.String("operation", "service.ProductService.CreateCart"),
			zap.Error(err))
		return err
	}

	s.logger.Info("Success sending the message to kafka")

	return nil
}

func (s *serviceProduct) UpdateCartItem(ctx context.Context, dto *domain.CartDTO) error {
	cartItem := &domain.Cart{
		UserId:    dto.UserId,
		ProductId: dto.ProductId,
		Quantity:  dto.Quantity,
	}

	if cartItem.ProductId == 0 || cartItem.UserId == 0 || cartItem.Quantity == 0 {
		s.logger.Error("Empty input",
			zap.String("operation", "service.ProductService.CreateCart"))
		return erorrs.ErrEmptyInput
	}

	msg := kafka.KafkaMessage{
		Id:    dto.UserId,
		Event: "cart_update",
		Data: map[string]interface{}{
			"UserId":    dto.UserId,
			"ProductId": dto.ProductId,
			"Name":      dto.Name,
			"Category":  dto.Category,
		},
	}

	key := strconv.Itoa(dto.UserId)

	err := s.kafka.SendMessage(key, msg)
	if err != nil {
		s.logger.Error("Error sending the message",
			zap.String("operation", "service.ProductService.UpdateCartItem"),
			zap.Error(err))
		return err
	}

	s.logger.Info("Success sending the message to kafka")

	err = s.repo.UpdateCartItem(ctx, cartItem)
	if errors.Is(err, erorrs.ErrRowsNull) {
		s.logger.Error("Cart item for user not found",
			zap.String("operation", "service.ProductService.UpdateCartItem"),
			zap.Error(err))
		return erorrs.ErrCartNotFound
	}

	return nil
}

func (s *serviceProduct) DeleteCartItem(ctx context.Context, userID int, productId int) error {
	if userID == 0 || productId == 0 {
		s.logger.Error("Empty input",
			zap.String("operation", "service.ProductService.DeleteCartItem"))
		return erorrs.ErrEmptyInput
	}

	err := s.repo.DeleteCartItem(ctx, userID, productId)
	if err != nil {
		s.logger.Error("Cart item for user not found",
			zap.String("operation", "service.ProductService.DeleteCartItem"),
			zap.Error(err))
		return erorrs.ErrCartNotFound
	}

	return nil
}

func (s *serviceProduct) GetAllProducts(ctx context.Context, pageSize int, page int) (domain.ProductsResponse, error) {
	if pageSize < 0 || page < 0 {
		s.logger.Error("Invalid page or pageSize",
			zap.String("operation", "service.ProductService.GetAllProducts"),
			zap.Int("pageSize", pageSize),
			zap.Int("page", page))
		return domain.ProductsResponse{}, errors.New("invalid page or pageSize")
	}

	products, total, err := s.repo.GetAllProducts(ctx, pageSize, page)
	if err != nil {
		s.logger.Error("Error getting products",
			zap.String("operation", "service.ProductService.GetAllProducts"),
			zap.Error(err))
		return domain.ProductsResponse{}, err
	}

	return domain.ProductsResponse{
		Data:  products,
		Total: total,
	}, nil
}
