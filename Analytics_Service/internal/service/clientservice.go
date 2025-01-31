package service

import (
	"Analytics_Service/internal/domain"
	"Analytics_Service/internal/logger"
	"Analytics_Service/internal/repository"

	"go.uber.org/zap"
)

type ServiceClient struct {
	repo   repository.RepositoryClient
	logger logger.Logger
}

func NewClientService(repo repository.RepositoryClient, logger logger.Logger) *ServiceClient {
	return &ServiceClient{
		repo:   repo,
		logger: logger,
	}
}

func (s *ServiceClient) GetUserActions() ([]domain.UserAction, error) {
	actions, err := s.repo.GetUserActions()
	if err != nil {
		s.logger.Error("Error getting user actions",
			zap.String("operation", "service.Client.GetUserActions"),
			zap.Error(err))
		return nil, err
	}
	return actions, nil
}

func (s *ServiceClient) GetUserEvents() ([]domain.UserEvent, error) {
	events, err := s.repo.GetUserEvents()
	if err != nil {
		s.logger.Error("Error getting user events",
			zap.String("operation", "service.Client.GetUserEvents"),
			zap.Error(err))
		return nil, err
	}
	return events, nil
}

func (s *ServiceClient) GetProductStats() ([]domain.ProductStat, error) {
	stats, err := s.repo.GetProductStats()
	if err != nil {
		s.logger.Error("Error getting user events",
			zap.String("operation", "service.Client.GetProductStats"),
			zap.Error(err))
		return nil, err
	}
	return stats, nil
}
