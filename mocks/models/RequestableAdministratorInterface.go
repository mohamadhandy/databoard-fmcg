// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	models "klikdaily-databoard/models"

	mock "github.com/stretchr/testify/mock"
)

// RequestableAdministratorInterface is an autogenerated mock type for the RequestableAdministratorInterface type
type RequestableAdministratorInterface struct {
	mock.Mock
}

// ForCreation provides a mock function with given fields:
func (_m *RequestableAdministratorInterface) ForCreation() models.Admin {
	ret := _m.Called()

	var r0 models.Admin
	if rf, ok := ret.Get(0).(func() models.Admin); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(models.Admin)
	}

	return r0
}

type mockConstructorTestingTNewRequestableAdministratorInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewRequestableAdministratorInterface creates a new instance of RequestableAdministratorInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRequestableAdministratorInterface(t mockConstructorTestingTNewRequestableAdministratorInterface) *RequestableAdministratorInterface {
	mock := &RequestableAdministratorInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
