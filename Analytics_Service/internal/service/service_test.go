package service

import (
	"Analytics_Service/internal/logger"
	"Analytics_Service/internal/repository/mocks"
	"Analytics_Service/pkg/kafka"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Handle(t *testing.T) {
	mockRepo := mocks.NewAnaliticRepository(t)
	mockLogger := logger.NewNoOpLogger()
	handler := NewHandler((logger.Logger)(mockLogger), mockRepo)

	testCases := []struct {
		name        string
		msg         kafka.KafkaMessage
		expectedErr error
		setupMocks  func()
	}{
		{
			name: "Success - Cart Event",
			msg: kafka.KafkaMessage{
				Topic: "product_updates",
				Event: "cart_create",
				Data: map[string]interface{}{
					"UserId":    1.0,
					"ProductId": 2.0,
					"Name":      "Product A",
					"Category":  "Category A",
				},
			},
			expectedErr: nil,
			setupMocks: func() {
				mockRepo.On("SaveUserAction", uint32(1), uint32(2), "cart_create", "Product A", "Category A").Return(nil)
				mockRepo.On("UpdateProductStats", uint32(2), "cart_add_count").Return(nil)
			},
		},
		{
			name: "Success - User Event",
			msg: kafka.KafkaMessage{
				Topic: "user_updates",
				Event: "user_create",
				Id:    1,
				Data: map[string]interface{}{
					"email": "user@example.com",
				},
			},
			expectedErr: nil,
			setupMocks: func() {
				mockRepo.On("SaveUserEvent", uint32(1), "user@example.com", "user_create").Return(nil)
			},
		},
		{
			name: "Repository Error - Cart Event",
			msg: kafka.KafkaMessage{
				Topic: "product_updates",
				Event: "cart_create",
				Data: map[string]interface{}{
					"UserId":    1.0,
					"ProductId": 2.0,
					"Name":      "Product A",
					"Category":  "Category A",
				},
			},
			expectedErr: errors.New("repository error"),
			setupMocks: func() {
				mockRepo.On("SaveUserAction", uint32(1), uint32(2), "cart_create", "Product A", "Category A").Return(errors.New("repository error"))
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

			err := handler.Handle(tc.msg)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
