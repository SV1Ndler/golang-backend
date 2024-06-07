// Code generated by mockery v2.40.2. DO NOT EDIT.

package mocks

import (
	models "url-shortener/pkg/models"

	mock "github.com/stretchr/testify/mock"
)

// MappingGetterAll is an autogenerated mock type for the MappingGetterAll type
type MappingGetterAll struct {
	mock.Mock
}

// GetAllMappingsWithLink provides a mock function with given fields:
func (_m *MappingGetterAll) GetAllMappingsWithLink() ([]models.PostImageMappingWithLink, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAllMappingsWithLink")
	}

	var r0 []models.PostImageMappingWithLink
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]models.PostImageMappingWithLink, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []models.PostImageMappingWithLink); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.PostImageMappingWithLink)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMappingGetterAll creates a new instance of MappingGetterAll. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMappingGetterAll(t interface {
	mock.TestingT
	Cleanup(func())
}) *MappingGetterAll {
	mock := &MappingGetterAll{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
