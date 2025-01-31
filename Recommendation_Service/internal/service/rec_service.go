package service

import (
	"Recommendation_Service/internal/domain"
	"Recommendation_Service/internal/logger"
	"Recommendation_Service/internal/repository"
	"Recommendation_Service/pkg/redis"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type ServiceRec interface {
	GetRecommendations(userId int) (map[string][]domain.Product, error)
	GetRec(category string, productId int) ([]domain.Product, error)
	GetProductCategory(userId int) ([]domain.Product, error)
}

type serviceRec struct {
	repo   repository.ProductRepository
	logger *logger.Logger
	redis  redis.RedisClient
}

func NewRecService(repo repository.ProductRepository, logger *logger.Logger, redis redis.RedisClient) *serviceRec {
	return &serviceRec{
		repo:   repo,
		logger: logger,
		redis:  redis,
	}
}

// извиняюсь, захардкодил, обещаю так больше не делать)
func (s *serviceRec) GetRecommendations(userId int) (map[string][]domain.Product, error) {
	cachedData, err := redis.GetRecommendationsFromCache(s.redis, userId)
	if err == nil {
		return cachedData, nil
	}

	userProducts, err := s.GetProductCategory(userId)
	if err != nil {
		s.logger.Error("Error get data",
			zap.String("operation", "service.RecService.GetRecommendation"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	if len(userProducts) == 0 {
		s.logger.Error("user products not found",
			zap.String("operation", "service.RecService.GetRecommendation"),
			zap.Error(err))
		return nil, fmt.Errorf("user products not found ")
	}

	recommendations := make(map[string][]domain.Product)
	for _, product := range userProducts {
		products, err := s.GetRec(product.Category, product.ProductId)

		if err != nil {
			s.logger.Error("failed to get recommendation",
				zap.String("operation", "service.RecService.GetRecommendation"),
				zap.Error(err))
			return nil, fmt.Errorf("failed to get recommendations: %w", err)
		}
		recommendations[product.Category] = products

		for _, rec := range products {
			err := s.repo.SaveRecommendation(userId, rec.ProductId, rec.Name, rec.Category)
			if err != nil {
				s.logger.Error("failed to save recommendation",
					zap.String("operation", "service.RecService.GetRecommendation"),
					zap.Error(err))
				return nil, fmt.Errorf("failed to save recommendation: %w", err)
			}
		}
	}

	if err := redis.SaveRecommendationsToCache(s.redis, userId, recommendations, time.Hour); err != nil {
		s.logger.Error("failed to save recommendation into cache",
			zap.String("operation", "service.RecService.GetRecommendation"),
			zap.Error(err))
		return nil, fmt.Errorf("failed to save recommendations to Redis: %w", err)
	}

	return recommendations, nil
}

func (s *serviceRec) GetRec(category string, productId int) ([]domain.Product, error) {
	return s.repo.Get(category, productId)
}

func (s *serviceRec) GetProductCategory(userId int) ([]domain.Product, error) {
	return s.repo.GetUserProducts(userId)
}
