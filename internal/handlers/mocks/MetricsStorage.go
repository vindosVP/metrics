// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MetricsStorage is an autogenerated mock type for the MetricsStorage type
type MetricsStorage struct {
	mock.Mock
}

// GetAllCounter provides a mock function with given fields:
func (_m *MetricsStorage) GetAllCounter() (map[string]int64, error) {
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

// GetAllGauge provides a mock function with given fields:
func (_m *MetricsStorage) GetAllGauge() (map[string]float64, error) {
	ret := _m.Called()

	var r0 map[string]float64
	var r1 error
	if rf, ok := ret.Get(0).(func() (map[string]float64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() map[string]float64); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]float64)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCounter provides a mock function with given fields: name
func (_m *MetricsStorage) GetCounter(name string) (int64, error) {
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

// GetGauge provides a mock function with given fields: name
func (_m *MetricsStorage) GetGauge(name string) (float64, error) {
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

// SetCounter provides a mock function with given fields: name, v
func (_m *MetricsStorage) SetCounter(name string, v int64) (int64, error) {
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

// UpdateCounter provides a mock function with given fields: name, v
func (_m *MetricsStorage) UpdateCounter(name string, v int64) (int64, error) {
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

// UpdateGauge provides a mock function with given fields: name, v
func (_m *MetricsStorage) UpdateGauge(name string, v float64) (float64, error) {
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

type mockConstructorTestingTNewMetricsStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewMetricsStorage creates a new instance of MetricsStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMetricsStorage(t mockConstructorTestingTNewMetricsStorage) *MetricsStorage {
	mock := &MetricsStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
