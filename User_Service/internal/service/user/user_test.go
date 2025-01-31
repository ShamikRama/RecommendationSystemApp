package user

import (
	"User_Service/internal/converter"
	"User_Service/internal/erorrs"
	"User_Service/internal/logger"
	"User_Service/internal/repository/mocks"
	"User_Service/internal/service/model"
	"User_Service/internal/utils"
	mocks2 "User_Service/pkg/kafka/mocks"
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceUser_UpdateUser(t *testing.T) {
	mockRepo := mocks.NewUser(t)
	mockKafka := mocks2.NewKafka(t)
	mockLogger := logger.NewNoOpLogger()

	service := ServiceUser{
		repo:   mockRepo,
		kafka:  mockKafka,
		logger: (*logger.Logger)(mockLogger),
	}

	ctx := context.Background()

	testCases := []struct {
		name          string
		userID        int
		input         model.UpdateInfoUser
		expectedID    int
		expectedError error
		setupMocks    func(userID int, input model.UpdateInfoUser)
	}{
		{
			name:   "Success",
			userID: 1,
			input: model.UpdateInfoUser{
				Email:    "new@example.com",
				Password: "newpassword123",
			},
			expectedID:    1,
			expectedError: nil,
			setupMocks: func(userID int, input model.UpdateInfoUser) {
				user := converter.FromUpdateToUser(input)
				user.Id = userID
				user.Password = utils.GeneratePasswordHash(input.Password)

				mockRepo.On("UpdateUser", ctx, user).Return(userID, nil)
				mockKafka.On("SendMessage", strconv.Itoa(userID), mock.AnythingOfType("*kafka.KafkaMessage")).Return(nil)
			},
		},
		{
			name:   "ErrorUserNotFound",
			userID: 999,
			input: model.UpdateInfoUser{
				Email:    "notfound@example.com",
				Password: "invalidpass",
			},
			expectedID:    0,
			expectedError: erorrs.ErrUserNotFound,
			setupMocks: func(userID int, input model.UpdateInfoUser) {
				user := converter.FromUpdateToUser(input)
				user.Id = userID
				user.Password = utils.GeneratePasswordHash(input.Password)

				mockRepo.On("UpdateUser", ctx, user).Return(0, erorrs.ErrNoRows)
			},
		},
		{
			name:   "Empty Input",
			userID: 1,
			input: model.UpdateInfoUser{
				Email:    "",
				Password: "",
			},
			expectedID:    0,
			expectedError: erorrs.ErrEmptyInput,
			setupMocks: func(userID int, input model.UpdateInfoUser) {
				// здесь ничего не делаем так как дальше не идет
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks(tc.userID, tc.input)

			updatedID, err := service.UpdateUser(ctx, tc.userID, tc.input)

			assert.Equal(t, tc.expectedID, updatedID)
			if tc.expectedError != nil {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockKafka.AssertExpectations(t)
		})
	}
}
