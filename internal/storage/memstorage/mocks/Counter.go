// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Counter is an autogenerated mock type for the Counter type
type Counter struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, name
func (_m *Counter) Get(ctx context.Context, name string) (int64, error) {
	ret := _m.Called(ctx, name)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int64, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int64); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx
func (_m *Counter) GetAll(ctx context.Context) (map[string]int64, error) {
	ret := _m.Called(ctx)

	var r0 map[string]int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (map[string]int64, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) map[string]int64); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]int64)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: ctx, name, v
func (_m *Counter) Set(ctx context.Context, name string, v int64) (int64, error) {
	ret := _m.Called(ctx, name, v)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) (int64, error)); ok {
		return rf(ctx, name, v)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) int64); ok {
		r0 = rf(ctx, name, v)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int64) error); ok {
		r1 = rf(ctx, name, v)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, name, v
func (_m *Counter) Update(ctx context.Context, name string, v int64) (int64, error) {
	ret := _m.Called(ctx, name, v)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) (int64, error)); ok {
		return rf(ctx, name, v)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) int64); ok {
		r0 = rf(ctx, name, v)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int64) error); ok {
		r1 = rf(ctx, name, v)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCounter interface {
	mock.TestingT
	Cleanup(func())
}

// NewCounter creates a new instance of Counter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCounter(t mockConstructorTestingTNewCounter) *Counter {
	mock := &Counter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
