// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Nats is an autogenerated mock type for the Nats type
type Nats struct {
	mock.Mock
}

// CloseConnect provides a mock function with given fields:
func (_m *Nats) CloseConnect() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetMsg provides a mock function with given fields: _a0
func (_m *Nats) GetMsg(_a0 chan []byte) {
	_m.Called(_a0)
}

type mockConstructorTestingTNewNats interface {
	mock.TestingT
	Cleanup(func())
}

// NewNats creates a new instance of Nats. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNats(t mockConstructorTestingTNewNats) *Nats {
	mock := &Nats{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
