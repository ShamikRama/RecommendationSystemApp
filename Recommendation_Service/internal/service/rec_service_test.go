package service

import (
	"Recommendation_Service/internal/domain"
	"Recommendation_Service/internal/logger"
	"Recommendation_Service/internal/repository/mocks"
	mocks2 "Recommendation_Service/pkg/redis/mocks"
	"errors"
	"fmt"
	"testing"
	"time"

	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceRec_GetRecommendations(t *testing.T) {
	mockRepo := mocks.NewProductRepository(t)
	mockRedis := mocks2.NewRedisClient(t)
	mockLogger := logger.NewLogger()

	service := NewRecService(mockRepo, mockLogger, mockRedis)

	userId := 1
	cachedRecommendations := map[string][]domain.Product{
		"electronics": {
			{ProductId: 1, Name: "Laptop", Category: "electronics"},
		},
	}

	userProducts := []domain.Product{
		{ProductId: 1, Name: "Laptop", Category: "electronics"},
	}

	recommendations := map[string][]domain.Product{
		"electronics": {
			{ProductId: 2, Name: "Smartphone", Category: "electronics"},
		},
	}

	testCases := []struct {
		name          string
		userId        int
		expectedData  map[string][]domain.Product
		expectedError error
		setupMocks    func()
	}{
		{
			name:          "Success with cache",
			userId:        userId,
			expectedData:  cachedRecommendations,
			expectedError: nil,
			setupMocks: func() {
				cachedDataJSON, _ := json.Marshal(cachedRecommendations)
				mockRedis.On("Get", mock.Anything, "recommendations:1").Return(string(cachedDataJSON), nil).Once()
			},
		},
		{
			name:          "Success without cache",
			userId:        userId,
			expectedData:  recommendations,
			expectedError: nil,
			setupMocks: func() {
				mockRedis.On("Get", mock.Anything, "recommendations:1").Return("", errors.New("cache miss")).Once()

				mockRepo.On("GetUserProducts", userId).Return(userProducts, nil).Once()

				mockRepo.On("Get", "electronics", 1).Return([]domain.Product{
					{ProductId: 2, Name: "Smartphone", Category: "electronics"},
				}, nil).Once()

				mockRepo.On("SaveRecommendation", userId, 2, "Smartphone", "electronics").Return(nil).Once()

				mockRedis.On("Set", mock.Anything, "recommendations:1", mock.Anything, time.Hour).Return(nil).Once()
			},
		},
		{
			name:          "Error getting user products",
			userId:        userId,
			expectedData:  nil,
			expectedError: fmt.Errorf("failed to get user data: %w", errors.New("database error")),
			setupMocks: func() {
				mockRedis.On("Get", mock.Anything, "recommendations:1").Return("", errors.New("cache miss")).Once()
				mockRepo.On("GetUserProducts", userId).Return(nil, errors.New("database error")).Once()
			},
		},
		{
			name:          "Error getting recommendations",
			userId:        userId,
			expectedData:  nil,
			expectedError: fmt.Errorf("failed to get recommendations: %w", errors.New("database error")),
			setupMocks: func() {
				mockRedis.On("Get", mock.Anything, "recommendations:1").Return("", errors.New("cache miss")).Once()
				mockRepo.On("GetUserProducts", userId).Return(userProducts, nil).Once()
				mockRepo.On("Get", "electronics", 1).Return(nil, errors.New("database error")).Once()
			},
		},
		{
			name:          "Error saving recommendations to cache",
			userId:        userId,
			expectedData:  nil,
			expectedError: fmt.Errorf("failed to save recommendations to Redis: %w", errors.New("cache save error")),
			setupMocks: func() {
				mockRedis.On("Get", mock.Anything, "recommendations:1").Return("", errors.New("cache miss")).Once()
				mockRepo.On("GetUserProducts", userId).Return(userProducts, nil).Once()
				mockRepo.On("Get", "electronics", 1).Return([]domain.Product{
					{ProductId: 2, Name: "Smartphone", Category: "electronics"},
				}, nil).Once()
				mockRepo.On("SaveRecommendation", userId, 2, "Smartphone", "electronics").Return(nil).Once()
				mockRedis.On("Set", mock.Anything, "recommendations:1", mock.Anything, time.Hour).Return(errors.New("cache save error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			result, err := service.GetRecommendations(tc.userId)

			assert.Equal(t, tc.expectedData, result)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockRedis.AssertExpectations(t)
		})
	}
}
