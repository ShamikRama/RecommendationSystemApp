// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	domain "Analytics_Service/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// RepositoryClient is an autogenerated mock type for the RepositoryClient type
type RepositoryClient struct {
	mock.Mock
}

// GetProductStats provides a mock function with given fields:
func (_m *RepositoryClient) GetProductStats() ([]domain.ProductStat, error) {
	ret := _m.Called()

	var r0 []domain.ProductStat
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]domain.ProductStat, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []domain.ProductStat); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.ProductStat)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserActions provides a mock function with given fields:
func (_m *RepositoryClient) GetUserActions() ([]domain.UserAction, error) {
	ret := _m.Called()

	var r0 []domain.UserAction
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]domain.UserAction, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []domain.UserAction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.UserAction)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserEvents provides a mock function with given fields:
func (_m *RepositoryClient) GetUserEvents() ([]domain.UserEvent, error) {
	ret := _m.Called()

	var r0 []domain.UserEvent
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]domain.UserEvent, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []domain.UserEvent); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.UserEvent)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepositoryClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepositoryClient creates a new instance of RepositoryClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepositoryClient(t mockConstructorTestingTNewRepositoryClient) *RepositoryClient {
	mock := &RepositoryClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
