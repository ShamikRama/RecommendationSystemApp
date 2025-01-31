// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	model "User_Service/internal/repository/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// User is an autogenerated mock type for the User type
type User struct {
	mock.Mock
}

// UpdateUser provides a mock function with given fields: ctx, user
func (_m *User) UpdateUser(ctx context.Context, user model.User) (int, error) {
	ret := _m.Called(ctx, user)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.User) (int, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.User) int); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUser interface {
	mock.TestingT
	Cleanup(func())
}

// NewUser creates a new instance of User. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUser(t mockConstructorTestingTNewUser) *User {
	mock := &User{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
