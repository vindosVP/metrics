// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Counter is an autogenerated mock type for the Counter type
type Counter struct {
	mock.Mock
}

// Get provides a mock function with given fields: name
func (_m *Counter) Get(name string) (int64, error) {
	ret := _m.Called(name)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int64, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) int64); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *Counter) GetAll() (map[string]int64, error) {
	ret := _m.Called()

	var r0 map[string]int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (map[string]int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() map[string]int64); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]int64)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: name, v
func (_m *Counter) Set(name string, v int64) (int64, error) {
	ret := _m.Called(name, v)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, int64) (int64, error)); ok {
		return rf(name, v)
	}
	if rf, ok := ret.Get(0).(func(string, int64) int64); ok {
		r0 = rf(name, v)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string, int64) error); ok {
		r1 = rf(name, v)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: name, v
func (_m *Counter) Update(name string, v int64) (int64, error) {
	ret := _m.Called(name, v)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, int64) (int64, error)); ok {
		return rf(name, v)
	}
	if rf, ok := ret.Get(0).(func(string, int64) int64); ok {
		r0 = rf(name, v)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(string, int64) error); ok {
		r1 = rf(name, v)
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
