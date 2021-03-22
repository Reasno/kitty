// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	io "io"

	mock "github.com/stretchr/testify/mock"
)

// Uploader is an autogenerated mock type for the Uploader type
type Uploader struct {
	mock.Mock
}

// Upload provides a mock function with given fields: ctx, reader
func (_m *Uploader) Upload(ctx context.Context, reader io.Reader) (string, error) {
	ret := _m.Called(ctx, reader)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, io.Reader) string); ok {
		r0 = rf(ctx, reader)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, io.Reader) error); ok {
		r1 = rf(ctx, reader)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
