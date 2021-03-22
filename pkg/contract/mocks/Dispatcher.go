// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	contract "glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

// Dispatcher is an autogenerated mock type for the Dispatcher type
type Dispatcher struct {
	mock.Mock
}

// Dispatch provides a mock function with given fields: event
func (_m *Dispatcher) Dispatch(event contract.Event) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(contract.Event) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Subscribe provides a mock function with given fields: listener
func (_m *Dispatcher) Subscribe(listener contract.Listener) {
	_m.Called(listener)
}
