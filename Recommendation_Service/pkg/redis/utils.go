package redis

import (
	"Recommendation_Service/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// InvalidateCache удаляет кэш рекомендаций для пользователя
func InvalidateCache(client RedisClient, userId int) error {
	ctx := context.Background()
	cacheKey := getRecommendationsCacheKey(userId)
	return client.Del(ctx, cacheKey)
}

// GetRecommendationsFromCache получает рекомендации из Redis
func GetRecommendationsFromCache(client RedisClient, userId int) (map[string][]domain.Product, error) {
	ctx := context.Background()
	cacheKey := getRecommendationsCacheKey(userId)

	cachedData, err := client.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	var recommendations map[string][]domain.Product
	if err := json.Unmarshal([]byte(cachedData), &recommendations); err != nil {
		return nil, err
	}

	return recommendations, nil
}

// SaveRecommendationsToCache сохраняет рекомендации в Redis
func SaveRecommendationsToCache(client RedisClient, userId int, recommendations map[string][]domain.Product, expiration time.Duration) error {
	ctx := context.Background()
	cacheKey := getRecommendationsCacheKey(userId)
	return client.Set(ctx, cacheKey, recommendations, expiration)
}

// getRecommendationsCacheKey возвращает ключ для кэша рекомендаций
func getRecommendationsCacheKey(userId int) string {
	return fmt.Sprintf("recommendations:%d", userId)
}
