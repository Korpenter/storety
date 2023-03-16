// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// TokenAuth is an autogenerated mock type for the TokenAuth type
type TokenAuth struct {
	mock.Mock
}

type TokenAuth_Expecter struct {
	mock *mock.Mock
}

func (_m *TokenAuth) EXPECT() *TokenAuth_Expecter {
	return &TokenAuth_Expecter{mock: &_m.Mock}
}

// GenerateTokenPair provides a mock function with given fields: id, sessionID
func (_m *TokenAuth) GenerateTokenPair(id uuid.UUID, sessionID uuid.UUID) (string, string, error) {
	ret := _m.Called(id, sessionID)

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) (string, string, error)); ok {
		return rf(id, sessionID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID) string); ok {
		r0 = rf(id, sessionID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID, uuid.UUID) string); ok {
		r1 = rf(id, sessionID)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(uuid.UUID, uuid.UUID) error); ok {
		r2 = rf(id, sessionID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// TokenAuth_GenerateTokenPair_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateTokenPair'
type TokenAuth_GenerateTokenPair_Call struct {
	*mock.Call
}

// GenerateTokenPair is a helper method to define mock.On call
//   - id uuid.UUID
//   - sessionID uuid.UUID
func (_e *TokenAuth_Expecter) GenerateTokenPair(id interface{}, sessionID interface{}) *TokenAuth_GenerateTokenPair_Call {
	return &TokenAuth_GenerateTokenPair_Call{Call: _e.mock.On("GenerateTokenPair", id, sessionID)}
}

func (_c *TokenAuth_GenerateTokenPair_Call) Run(run func(id uuid.UUID, sessionID uuid.UUID)) *TokenAuth_GenerateTokenPair_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *TokenAuth_GenerateTokenPair_Call) Return(_a0 string, _a1 string, _a2 error) *TokenAuth_GenerateTokenPair_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *TokenAuth_GenerateTokenPair_Call) RunAndReturn(run func(uuid.UUID, uuid.UUID) (string, string, error)) *TokenAuth_GenerateTokenPair_Call {
	_c.Call.Return(run)
	return _c
}

// Verify provides a mock function with given fields: _a0
func (_m *TokenAuth) Verify(_a0 string) (uuid.UUID, error) {
	ret := _m.Called(_a0)

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (uuid.UUID, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) uuid.UUID); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenAuth_Verify_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Verify'
type TokenAuth_Verify_Call struct {
	*mock.Call
}

// Verify is a helper method to define mock.On call
//   - _a0 string
func (_e *TokenAuth_Expecter) Verify(_a0 interface{}) *TokenAuth_Verify_Call {
	return &TokenAuth_Verify_Call{Call: _e.mock.On("Verify", _a0)}
}

func (_c *TokenAuth_Verify_Call) Run(run func(_a0 string)) *TokenAuth_Verify_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *TokenAuth_Verify_Call) Return(_a0 uuid.UUID, _a1 error) *TokenAuth_Verify_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TokenAuth_Verify_Call) RunAndReturn(run func(string) (uuid.UUID, error)) *TokenAuth_Verify_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewTokenAuth interface {
	mock.TestingT
	Cleanup(func())
}

// NewTokenAuth creates a new instance of TokenAuth. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTokenAuth(t mockConstructorTestingTNewTokenAuth) *TokenAuth {
	mock := &TokenAuth{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}