package service

import (
	"Recommendation_Service/internal/erorrs"
	"Recommendation_Service/internal/logger"
	"Recommendation_Service/internal/repository/mocks"
	"Recommendation_Service/pkg/kafka"
	mocks2 "Recommendation_Service/pkg/redis/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"errors"
	"testing"
)

func TestHandler_Handle(t *testing.T) {
	mockRepo := mocks.NewProductRepository(t)
	mockRedis := mocks2.NewRedisClient(t)
	mockLogger := logger.NewNoOpLogger()

	handler := &Handler{
		repo:   mockRepo,
		redis:  mockRedis,
		logger: (logger.Logger)(mockLogger),
	}

	testCases := []struct {
		name          string
		msg           kafka.KafkaMessage
		expectedError error
		setupMocks    func()
	}{
		{
			name: "Success",
			msg: kafka.KafkaMessage{
				Data: map[string]interface{}{
					"UserId":    float64(1),
					"ProductId": float64(2),
					"Name":      "Test Product",
					"Category":  "Test Category",
				},
			},
			expectedError: nil,
			setupMocks: func() {
				mockRepo.On("SaveDataProduct", 2, "Test Product", "Test Category").Return(nil).Once()
				mockRepo.On("SaveUserActions", 1, 2, "Test Product", "Test Category").Return(nil).Once()
				mockRedis.On("Del", mock.Anything, "recommendations:1").Return(nil).Once() // Мок для Del
			},
		},
		{
			name: "Empty Message",
			msg: kafka.KafkaMessage{
				Data: map[string]interface{}{
					"UserId":    float64(0),
					"ProductId": float64(0),
					"Name":      "",
					"Category":  "",
				},
			},
			expectedError: erorrs.ErrEmptyMessage,
			setupMocks: func() {
			},
		},
		{
			name: "SaveDataProduct Error",
			msg: kafka.KafkaMessage{
				Data: map[string]interface{}{
					"UserId":    float64(1),
					"ProductId": float64(2),
					"Name":      "Test Product",
					"Category":  "Test Category",
				},
			},
			expectedError: errors.New("save data product error"),
			setupMocks: func() {
				mockRepo.On("SaveDataProduct", 2, "Test Product", "Test Category").Return(errors.New("save data product error")).Once()
			},
		},
		{
			name: "SaveUserActions Error",
			msg: kafka.KafkaMessage{
				Data: map[string]interface{}{
					"UserId":    float64(1),
					"ProductId": float64(2),
					"Name":      "Test Product",
					"Category":  "Test Category",
				},
			},
			expectedError: errors.New("save user actions error"),
			setupMocks: func() {
				mockRepo.On("SaveDataProduct", 2, "Test Product", "Test Category").Return(nil).Once()
				mockRepo.On("SaveUserActions", 1, 2, "Test Product", "Test Category").Return(errors.New("save user actions error")).Once()
			},
		},
		{
			name: "InvalidateCache Error",
			msg: kafka.KafkaMessage{
				Data: map[string]interface{}{
					"UserId":    float64(1),
					"ProductId": float64(2),
					"Name":      "Test Product",
					"Category":  "Test Category",
				},
			},
			expectedError: nil,
			setupMocks: func() {
				mockRepo.On("SaveDataProduct", 2, "Test Product", "Test Category").Return(nil).Once()
				mockRepo.On("SaveUserActions", 1, 2, "Test Product", "Test Category").Return(nil).Once()
				mockRedis.On("Del", mock.Anything, "recommendations:1").Return(errors.New("invalidate cache error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			err := handler.Handle(tc.msg)

			assert.Equal(t, tc.expectedError, err)

			mockRepo.AssertExpectations(t)
			mockRedis.AssertExpectations(t)
		})
	}
}
