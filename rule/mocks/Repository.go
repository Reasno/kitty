// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	rule "glab.tagtic.cn/ad_gains/kitty/rule"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// GetCompiled provides a mock function with given fields: ruleName
func (_m *Repository) GetCompiled(ruleName string) []rule.Rule {
	ret := _m.Called(ruleName)

	var r0 []rule.Rule
	if rf, ok := ret.Get(0).(func(string) []rule.Rule); ok {
		r0 = rf(ruleName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]rule.Rule)
		}
	}

	return r0
}

// GetRaw provides a mock function with given fields: ctx, key
func (_m *Repository) GetRaw(ctx context.Context, key string) ([]byte, error) {
	ret := _m.Called(ctx, key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsNewest provides a mock function with given fields: ctx, key, value
func (_m *Repository) IsNewest(ctx context.Context, key string, value string) (bool, error) {
	ret := _m.Called(ctx, key, value)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, key, value)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, key, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetRaw provides a mock function with given fields: ctx, key, value
func (_m *Repository) SetRaw(ctx context.Context, key string, value string) error {
	ret := _m.Called(ctx, key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WatchConfigUpdate provides a mock function with given fields: ctx
func (_m *Repository) WatchConfigUpdate(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
