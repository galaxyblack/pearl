package mocks

import mock "github.com/stretchr/testify/mock"
import openssl "github.com/mmcloughlin/openssl"

// PrivateKey is an autogenerated mock type for the PrivateKey type
type PrivateKey struct {
	mock.Mock
}

// MarshalPKCS1PrivateKeyPEM provides a mock function with given fields:
func (_m *PrivateKey) MarshalPKCS1PrivateKeyPEM() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarshalPKCS1PublicKeyDER provides a mock function with given fields:
func (_m *PrivateKey) MarshalPKCS1PublicKeyDER() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarshalPKCS1PublicKeyPEM provides a mock function with given fields:
func (_m *PrivateKey) MarshalPKCS1PublicKeyPEM() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignPKCS1v15 provides a mock function with given fields: _a0, _a1
func (_m *PrivateKey) SignPKCS1v15(_a0 openssl.Method, _a1 []byte) ([]byte, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(openssl.Method, []byte) []byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(openssl.Method, []byte) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
