// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "glab.tagtic.cn/ad_gains/kitty/app/entity"

	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// Exists provides a mock function with given fields: ctx, id
func (_m *UserRepository) Exists(ctx context.Context, id uint) bool {
	ret := _m.Called(ctx, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, uint) bool); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// UpdateCallback provides a mock function with given fields: ctx, id, f
func (_m *UserRepository) UpdateCallback(ctx context.Context, id uint, f func(*entity.User) error) error {
	ret := _m.Called(ctx, id, f)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, func(*entity.User) error) error); ok {
		r0 = rf(ctx, id, f)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
