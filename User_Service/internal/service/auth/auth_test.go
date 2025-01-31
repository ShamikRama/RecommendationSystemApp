package auth

import (
	"User_Service/internal/erorrs"
	"User_Service/internal/logger"
	"User_Service/internal/repository/mocks"
	"User_Service/internal/service/model"
	"User_Service/internal/utils"
	mocks2 "User_Service/pkg/kafka/mocks"
	"context"
	"testing"

	"User_Service/internal/converter"

	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceAuth_CreateUser(t *testing.T) {
	mockAuth := mocks.NewAuth(t)
	mockKafka := mocks2.NewKafka(t)
	mockLogger := logger.NewNoOpLogger()

	authService := ServiceAuth{
		repo:   mockAuth,
		kafka:  mockKafka,
		logger: (*logger.Logger)(mockLogger),
	}

	ctx := context.Background()
	login := model.Login{
		Email:    "test@example.com",
		Password: "password123",
	}
	user := converter.FromLoginToUser(login)
	user.Password = utils.GeneratePasswordHash(user.Password)

	testCases := []struct {
		nametest      string
		login         model.Login
		expectedID    int
		expectedError error
		setupMocks    func()
	}{
		{
			nametest: "Success",
			login: model.Login{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedID:    1,
			expectedError: nil,
			setupMocks: func() {
				mockAuth.On("Create", mock.Anything, mock.Anything).Return(1, nil)
				mockKafka.On("SendMessage", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			nametest: "Empty Input",
			login: model.Login{
				Email:    "",
				Password: "",
			},
			expectedID:    0,
			expectedError: erorrs.ErrEmptyInput,
			setupMocks: func() {
				// не нужны моки потому что у нас не прохоидт валидацию данные и действие отменяется
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.nametest, func(t *testing.T) {
			tc.setupMocks()

			userID, err := authService.CreateUser(ctx, tc.login)

			assert.Equal(t, tc.expectedID, userID)
			if tc.expectedError != nil {
				assert.True(t, errors.Is(err, tc.expectedError))
			} else {
				assert.Equal(t, tc.expectedError, err)
			}

			mockAuth.AssertExpectations(t)
			mockKafka.AssertExpectations(t)
		})
	}
}
