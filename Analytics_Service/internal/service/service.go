package service

import (
	"Analytics_Service/internal/logger"
	"Analytics_Service/internal/repository"
	"Analytics_Service/pkg/kafka"

	"go.uber.org/zap"
)

type Handler struct {
	logger logger.Logger
	repo   repository.AnaliticRepository
}

func NewHandler(logger logger.Logger, repo repository.AnaliticRepository) Handler {
	return Handler{
		logger: logger,
		repo:   repo,
	}
}

func (h *Handler) Handle(msg kafka.KafkaMessage) error {
	event := msg.Event
	topic := msg.Topic

	switch topic {
	case "product_updates":
		if event == "cart_create" {
			err := h.handleCartEvent(msg)
			if err != nil {
				h.logger.Error("Error handling cart event",
					zap.String("operation", "service.Service.Handle"),
					zap.Error(err))
				return err
			}
		} else {
			h.logger.Info("Unknown event for product_updates",
				zap.String("event", event),
				zap.String("operation", "service.Service.Handle"))
		}
	case "user_updates":
		if event == "user_create" {
			err := h.handleUserEvent(msg)
			if err != nil {
				h.logger.Error("Error handling user event",
					zap.String("operation", "service.Service.Handle"),
					zap.Error(err))
				return err
			}
		} else {
			h.logger.Info("Unknown event for user_updates",
				zap.String("operation", "service.Service.Handle"),
				zap.String("event", event))
		}
	default:
		h.logger.Info("Unknown topic",
			zap.String("operation", "service.Service.Handle"),
			zap.String("topic", topic))
	}

	return nil
}

func (h *Handler) handleCartEvent(msg kafka.KafkaMessage) error {
	userId := uint32(msg.Data["UserId"].(float64))
	productId := uint32(msg.Data["ProductId"].(float64))
	productName := msg.Data["Name"].(string)
	category := msg.Data["Category"].(string)

	err := h.repo.SaveUserAction(userId, productId, "cart_create", productName, category)
	if err != nil {
		h.logger.Error("Error save user action",
			zap.String("operation", "service.Service.handleCart"),
			zap.Error(err))
		return err
	}

	err = h.repo.UpdateProductStats(productId, "cart_add_count")
	if err != nil {
		h.logger.Error("Error update count action",
			zap.String("operation", "service.Service.handleCart"),
			zap.Error(err))
		return err
	}

	return nil
}

func (h *Handler) handleUserEvent(msg kafka.KafkaMessage) error {
	userId := uint32(msg.Id)
	event := msg.Event
	email := msg.Data["email"].(string)

	err := h.repo.SaveUserEvent(userId, email, event)
	if err != nil {
		h.logger.Error("Error save user event",
			zap.String("operation", "service.Service.handleUser"),
			zap.Error(err))
		return err
	}

	return nil
}
