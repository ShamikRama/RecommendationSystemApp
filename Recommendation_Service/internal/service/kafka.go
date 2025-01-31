package service

import (
	"Recommendation_Service/internal/erorrs"
	"Recommendation_Service/internal/logger"
	"Recommendation_Service/internal/repository"
	"Recommendation_Service/pkg/kafka"
	"Recommendation_Service/pkg/redis"
	"fmt"

	"go.uber.org/zap"
)

type Handler struct {
	logger logger.Logger
	repo   repository.ProductRepository
	redis  redis.RedisClient
}

func NewHandler(logger logger.Logger, repo repository.ProductRepository, redis redis.RedisClient) Handler {
	return Handler{
		logger: logger,
		repo:   repo,
		redis:  redis,
	}
}

// извиняюсь, захардкодил, обещаю так больше не делать)
func (h *Handler) Handle(msg kafka.KafkaMessage) error {
	userIdFloat, ok := msg.Data["UserId"].(float64)
	if !ok {
		h.logger.Info("Not a float64",
			zap.String("operations", "service.Kafka.Handle"),
			zap.Bool("res", ok))
		return fmt.Errorf("invalid type for UserId")
	}
	userId := int(userIdFloat)

	productIdFloat, ok := msg.Data["ProductId"].(float64)
	if !ok {
		h.logger.Info("Not a float64",
			zap.String("operations", "service.Kafka.Handle"),
			zap.Bool("res", ok))
		return fmt.Errorf("invalid type for ProductId")
	}
	productId := int(productIdFloat)

	name, ok := msg.Data["Name"].(string)
	if !ok {
		h.logger.Info("Not a string",
			zap.String("operations", "service.Kafka.Handle"),
			zap.Bool("res", ok))
		return fmt.Errorf("invalid type for Name")
	}

	category, ok := msg.Data["Category"].(string)
	if !ok {
		h.logger.Info("Not a string",
			zap.String("operations", "service.Kafka.Handle"),
			zap.Bool("res", ok))
		return fmt.Errorf("invalid type for Category")
	}

	if userId == 0 || productId == 0 || name == "" || category == "" {
		h.logger.Error("Empty meassage")
		return erorrs.ErrEmptyMessage
	}

	err := h.repo.SaveDataProduct(productId, name, category)
	if err != nil {
		h.logger.Info("Failed to save product into product",
			zap.String("operations", "service.Kafka.Handle"),
			zap.Error(err))
		return err
	}

	err = h.repo.SaveUserActions(userId, productId, name, category)
	if err != nil {
		h.logger.Info("Failed to save user action",
			zap.String("operations", "service.Kafka.Handle"),
			zap.Error(err))
		return err
	}

	if err := redis.InvalidateCache(h.redis, userId); err != nil {
		h.logger.Info("Failed to invalidate cache for user",
			zap.Int("user", userId))
		zap.Error(err)
	}

	return nil
}
