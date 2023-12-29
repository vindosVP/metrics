// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Gauge is an autogenerated mock type for the Gauge type
type Gauge struct {
	mock.Mock
}

// Get provides a mock function with given fields: name
func (_m *Gauge) Get(name string) (float64, error) {
	ret := _m.Called(name)

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (float64, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) float64); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: name, v
func (_m *Gauge) Update(name string, v float64) (float64, error) {
	ret := _m.Called(name, v)

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(string, float64) (float64, error)); ok {
		return rf(name, v)
	}
	if rf, ok := ret.Get(0).(func(string, float64) float64); ok {
		r0 = rf(name, v)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(string, float64) error); ok {
		r1 = rf(name, v)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewGauge interface {
	mock.TestingT
	Cleanup(func())
}

// NewGauge creates a new instance of Gauge. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGauge(t mockConstructorTestingTNewGauge) *Gauge {
	mock := &Gauge{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
