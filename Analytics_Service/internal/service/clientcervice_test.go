package service

import (
	"Analytics_Service/internal/domain"
	"Analytics_Service/internal/logger"
	"Analytics_Service/internal/repository/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceClient_GetUserActions(t *testing.T) {
	mockRepo := mocks.NewRepositoryClient(t)
	mockLogger := logger.NewNoOpLogger()
	service := NewClientService(mockRepo, (logger.Logger)(mockLogger))

	testCases := []struct {
		name        string
		expected    []domain.UserAction
		expectedErr error
		setupMocks  func()
	}{
		{
			name: "Success",
			expected: []domain.UserAction{
				{UserID: 1, ProductID: 1, ActionType: "click", Name: "Product A", Category: "Category A"},
			},
			expectedErr: nil,
			setupMocks: func() {
				mockRepo.On("GetUserActions").Return([]domain.UserAction{
					{UserID: 1, ProductID: 1, ActionType: "click", Name: "Product A", Category: "Category A"},
				}, nil)
			},
		},
		{
			name:        "Repository Error",
			expected:    nil,
			expectedErr: errors.New("repository error"),
			setupMocks: func() {
				mockRepo.On("GetUserActions").Return([]domain.UserAction(nil), errors.New("repository error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			t.Cleanup(func() {
				mockRepo.AssertExpectations(t)
				mockRepo.ExpectedCalls = nil
			})

			result, err := service.GetUserActions()

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestServiceClient_GetUserEvents(t *testing.T) {
	mockRepo := mocks.NewRepositoryClient(t)
	mockLogger := logger.NewNoOpLogger()
	service := NewClientService(mockRepo, (logger.Logger)(mockLogger))

	testCases := []struct {
		name        string
		expected    []domain.UserEvent
		expectedErr error
		setupMocks  func()
	}{
		{
			name: "Success",
			expected: []domain.UserEvent{
				{UserID: 1, EventType: "login", Email: "user@example.com"},
			},
			expectedErr: nil,
			setupMocks: func() {
				mockRepo.On("GetUserEvents").Return([]domain.UserEvent{
					{UserID: 1, EventType: "login", Email: "user@example.com"},
				}, nil)
			},
		},
		{
			name:        "Repository Error",
			expected:    nil,
			expectedErr: errors.New("repository error"),
			setupMocks: func() {
				mockRepo.On("GetUserEvents").Return([]domain.UserEvent(nil), errors.New("repository error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			t.Cleanup(func() {
				mockRepo.AssertExpectations(t)
				mockRepo.ExpectedCalls = nil
			})

			result, err := service.GetUserEvents()

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestServiceClient_GetProductStats(t *testing.T) {
	mockRepo := mocks.NewRepositoryClient(t)
	mockLogger := logger.NewNoOpLogger()
	service := NewClientService(mockRepo, (logger.Logger)(mockLogger))

	testCases := []struct {
		name        string
		expected    []domain.ProductStat
		expectedErr error
		setupMocks  func()
	}{
		{
			name: "Success",
			expected: []domain.ProductStat{
				{ProductID: 1, CartAddCount: 5},
			},
			expectedErr: nil,
			setupMocks: func() {
				mockRepo.On("GetProductStats").Return([]domain.ProductStat{
					{ProductID: 1, CartAddCount: 5},
				}, nil)
			},
		},
		{
			name:        "Repository Error",
			expected:    nil,
			expectedErr: errors.New("repository error"),
			setupMocks: func() {
				mockRepo.On("GetProductStats").Return([]domain.ProductStat(nil), errors.New("repository error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			t.Cleanup(func() {
				mockRepo.AssertExpectations(t)
				mockRepo.ExpectedCalls = nil
			})

			result, err := service.GetProductStats()

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
