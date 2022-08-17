// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	core "github.com/retitle/go-sdk/core"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Impersonation is an autogenerated mock type for the Impersonation type
type Impersonation struct {
	mock.Mock
}

// GetAccessToken provides a mock function with given fields:
func (_m *Impersonation) GetAccessToken() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetAccessTokenExpires provides a mock function with given fields:
func (_m *Impersonation) GetAccessTokenExpires() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// GetScopes provides a mock function with given fields:
func (_m *Impersonation) GetScopes() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// GetSub provides a mock function with given fields:
func (_m *Impersonation) GetSub() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SetAccessToken provides a mock function with given fields: _a0
func (_m *Impersonation) SetAccessToken(_a0 string) core.Impersonation {
	ret := _m.Called(_a0)

	var r0 core.Impersonation
	if rf, ok := ret.Get(0).(func(string) core.Impersonation); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Impersonation)
		}
	}

	return r0
}

// SetAccessTokenExpires provides a mock function with given fields: _a0
func (_m *Impersonation) SetAccessTokenExpires(_a0 time.Time) core.Impersonation {
	ret := _m.Called(_a0)

	var r0 core.Impersonation
	if rf, ok := ret.Get(0).(func(time.Time) core.Impersonation); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Impersonation)
		}
	}

	return r0
}

// SetScopes provides a mock function with given fields: _a0
func (_m *Impersonation) SetScopes(_a0 []string) core.Impersonation {
	ret := _m.Called(_a0)

	var r0 core.Impersonation
	if rf, ok := ret.Get(0).(func([]string) core.Impersonation); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Impersonation)
		}
	}

	return r0
}

// SetSub provides a mock function with given fields: _a0
func (_m *Impersonation) SetSub(_a0 string) core.Impersonation {
	ret := _m.Called(_a0)

	var r0 core.Impersonation
	if rf, ok := ret.Get(0).(func(string) core.Impersonation); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Impersonation)
		}
	}

	return r0
}
