package service

import (
	"Product_Service/internal/erorrs"
	"context"

	"github.com/stretchr/testify/assert"

	"errors"
	"testing"

	"Product_Service/internal/domain"
	"Product_Service/internal/logger"
	"Product_Service/internal/repository/mocks"
	mocks2 "Product_Service/pkg/kafka/mocks"

	"github.com/stretchr/testify/mock"
)

func TestCreateCart(t *testing.T) {
	repoMock := mocks.NewRepositoryProduct(t)
	kafkaMock := mocks2.NewKafka(t)
	mockLogger := logger.NewNoOpLogger()

	service := NewProductService(repoMock, (*logger.Logger)(mockLogger), kafkaMock)

	ctx := context.Background()

	testCases := []struct {
		name        string
		dto         domain.CartDTO
		expectedErr error
		setupMocks  func()
	}{
		{
			name: "Success",
			dto: domain.CartDTO{
				UserId:    1,
				ProductId: 2,
				Quantity:  3,
				Name:      "Test Product",
				Category:  "Test Category",
			},
			expectedErr: nil,
			setupMocks: func() {
				repoMock.On("CreateCartItem", ctx, mock.AnythingOfType("*domain.Cart")).Return(nil)
				kafkaMock.On("SendMessage", mock.AnythingOfType("string"), mock.AnythingOfType("kafka.KafkaMessage")).Return(nil)
			},
		},
		{
			name: "Repository Error",
			dto: domain.CartDTO{
				UserId:    1,
				ProductId: 2,
				Quantity:  3,
				Name:      "Test Product",
				Category:  "Test Category",
			},
			expectedErr: errors.New("repository error"),
			setupMocks: func() {
				repoMock.On("CreateCartItem", ctx, mock.AnythingOfType("*domain.Cart")).Return(errors.New("repository error"))
				kafkaMock.AssertNotCalled(t, "SendMessage")
			},
		},
		{
			name: "Kafka Error",
			dto: domain.CartDTO{
				UserId:    1,
				ProductId: 2,
				Quantity:  3,
				Name:      "Test Product",
				Category:  "Test Category",
			},
			expectedErr: errors.New("kafka error"),
			setupMocks: func() {
				repoMock.On("CreateCartItem", ctx, mock.AnythingOfType("*domain.Cart")).Return(nil)
				kafkaMock.On("SendMessage", mock.AnythingOfType("string"), mock.AnythingOfType("kafka.KafkaMessage")).Return(errors.New("kafka error"))
			},
		},
		{
			name: "Empty input",
			dto: domain.CartDTO{
				UserId:    0,
				ProductId: 0,
				Quantity:  0,
				Name:      "Test Product",
				Category:  "Test Category",
			},
			expectedErr: erorrs.ErrEmptyInput,
			setupMocks: func() {
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			t.Cleanup(func() {
				repoMock.AssertExpectations(t)
				kafkaMock.AssertExpectations(t)
				repoMock.ExpectedCalls = nil
				kafkaMock.ExpectedCalls = nil
			})

			err := service.CreateCart(ctx, &tc.dto)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateCartItem(t *testing.T) {
	repoMock := mocks.NewRepositoryProduct(t)
	kafkaMock := mocks2.NewKafka(t)
	mockLogger := logger.NewNoOpLogger()

	service := NewProductService(repoMock, (*logger.Logger)(mockLogger), kafkaMock)

	ctx := context.Background()

	testCases := []struct {
		name        string
		dto         domain.CartDTO
		expectedErr error
		setupMocks  func()
	}{
		{
			name: "Success",
			dto: domain.CartDTO{
				UserId:    1,
				ProductId: 2,
				Quantity:  3,
				Name:      "Test Product",
				Category:  "Test Category",
			},
			expectedErr: nil,
			setupMocks: func() {
				kafkaMock.On("SendMessage", mock.AnythingOfType("string"), mock.AnythingOfType("kafka.KafkaMessage")).Return(nil)
				repoMock.On("UpdateCartItem", ctx, mock.AnythingOfType("*domain.Cart")).Return(nil)
			},
		},
		{
			name: "Kafka Error",
			dto: domain.CartDTO{
				UserId:    1,
				ProductId: 2,
				Quantity:  3,
				Name:      "Test Product",
				Category:  "Test Category",
			},
			expectedErr: errors.New("kafka error"),
			setupMocks: func() {
				kafkaMock.On("SendMessage", mock.AnythingOfType("string"), mock.AnythingOfType("kafka.KafkaMessage")).Return(errors.New("kafka error"))
				repoMock.AssertNotCalled(t, "UpdateCartItem")
			},
		},
		{
			name: "Repository Error",
			dto: domain.CartDTO{
				UserId:    1,
				ProductId: 2,
				Quantity:  3,
				Name:      "Test Product",
				Category:  "Test Category",
			},
			expectedErr: erorrs.ErrCartNotFound,
			setupMocks: func() {
				kafkaMock.On("SendMessage", mock.AnythingOfType("string"), mock.AnythingOfType("kafka.KafkaMessage")).Return(nil)
				repoMock.On("UpdateCartItem", ctx, mock.AnythingOfType("*domain.Cart")).Return(erorrs.ErrRowsNull)
			},
		},
		{
			name: "Empty input",
			dto: domain.CartDTO{
				UserId:    0,
				ProductId: 0,
				Quantity:  0,
				Name:      "Test Product",
				Category:  "Test Category",
			},
			expectedErr: erorrs.ErrEmptyInput,
			setupMocks: func() {
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			t.Cleanup(func() {
				repoMock.AssertExpectations(t)
				kafkaMock.AssertExpectations(t)
				repoMock.ExpectedCalls = nil
				kafkaMock.ExpectedCalls = nil
			})

			err := service.UpdateCartItem(ctx, &tc.dto)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteCartItem(t *testing.T) {
	repoMock := mocks.NewRepositoryProduct(t)
	kafkaMock := mocks2.NewKafka(t)
	mockLogger := logger.NewNoOpLogger()

	service := NewProductService(repoMock, (*logger.Logger)(mockLogger), kafkaMock)

	ctx := context.Background()

	testCases := []struct {
		name        string
		userID      int
		productID   int
		expectedErr error
		setupMocks  func()
	}{
		{
			name:        "Success",
			userID:      1,
			productID:   2,
			expectedErr: nil,
			setupMocks: func() {
				repoMock.On("DeleteCartItem", ctx, 1, 2).Return(nil)
			},
		},
		{
			name:        "Repository Error",
			userID:      1,
			productID:   2,
			expectedErr: erorrs.ErrCartNotFound,
			setupMocks: func() {
				repoMock.On("DeleteCartItem", ctx, 1, 2).Return(erorrs.ErrRowsNull)
			},
		},
		{
			name:        "Empty Input",
			userID:      0,
			productID:   0,
			expectedErr: erorrs.ErrEmptyInput,
			setupMocks: func() {
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			t.Cleanup(func() {
				repoMock.AssertExpectations(t)
				repoMock.ExpectedCalls = nil
			})

			err := service.DeleteCartItem(ctx, tc.userID, tc.productID)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAllProducts(t *testing.T) {
	repoMock := mocks.NewRepositoryProduct(t)
	kafkaMock := mocks2.NewKafka(t)
	mockLogger := logger.NewNoOpLogger()

	service := NewProductService(repoMock, (*logger.Logger)(mockLogger), kafkaMock)

	ctx := context.Background()

	testCases := []struct {
		name             string
		pageSize         int
		page             int
		expectedResponse domain.ProductsResponse
		expectedErr      error
		setupMocks       func()
	}{
		{
			name:     "Success",
			pageSize: 10,
			page:     1,
			expectedResponse: domain.ProductsResponse{
				Data: []domain.Product{
					{Id: 1, Name: "Product 1"},
					{Id: 2, Name: "Product 2"},
				},
				Total: 2,
			},
			expectedErr: nil,
			setupMocks: func() {
				repoMock.On("GetAllProducts", ctx, 10, 1).Return([]domain.Product{
					{Id: 1, Name: "Product 1"},
					{Id: 2, Name: "Product 2"},
				}, 2, nil)
			},
		},
		{
			name:             "Repository Error",
			pageSize:         10,
			page:             1,
			expectedResponse: domain.ProductsResponse{},
			expectedErr:      errors.New("repository error"),
			setupMocks: func() {
				repoMock.On("GetAllProducts", ctx, 10, 1).Return(nil, 0, errors.New("repository error"))
			},
		},
		{
			name:     "Empty Result",
			pageSize: 10,
			page:     1,
			expectedResponse: domain.ProductsResponse{
				Data:  []domain.Product{},
				Total: 0,
			},
			expectedErr: nil,
			setupMocks: func() {
				repoMock.On("GetAllProducts", ctx, 10, 1).Return([]domain.Product{}, 0, nil)
			},
		},
		{
			name:             "Negative Page Size",
			pageSize:         -10,
			page:             1,
			expectedResponse: domain.ProductsResponse{},
			expectedErr:      errors.New("invalid page or pageSize"),
			setupMocks:       func() {},
		},
		{
			name:             "Negative Page",
			pageSize:         10,
			page:             -1,
			expectedResponse: domain.ProductsResponse{},
			expectedErr:      errors.New("invalid page or pageSize"),
			setupMocks:       func() {},
		},
		{
			name:             "Negative Page and Page Size",
			pageSize:         -10,
			page:             -1,
			expectedResponse: domain.ProductsResponse{},
			expectedErr:      errors.New("invalid page or pageSize"),
			setupMocks:       func() {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()

			defer func() {
				if !t.Failed() {
					repoMock.AssertExpectations(t)
				}
				repoMock.ExpectedCalls = nil
			}()

			response, err := service.GetAllProducts(ctx, tc.pageSize, tc.page)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedResponse, response)
		})
	}
}
