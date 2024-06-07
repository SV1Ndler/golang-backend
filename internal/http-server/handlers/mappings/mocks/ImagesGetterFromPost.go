// Code generated by mockery v2.40.2. DO NOT EDIT.

package mocks

import (
	models "url-shortener/pkg/models"

	mock "github.com/stretchr/testify/mock"
)

// ImagesGetterFromPost is an autogenerated mock type for the ImagesGetterFromPost type
type ImagesGetterFromPost struct {
	mock.Mock
}

// GetPostImages provides a mock function with given fields: postID
func (_m *ImagesGetterFromPost) GetPostImages(postID int) ([]models.Image, error) {
	ret := _m.Called(postID)

	if len(ret) == 0 {
		panic("no return value specified for GetPostImages")
	}

	var r0 []models.Image
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]models.Image, error)); ok {
		return rf(postID)
	}
	if rf, ok := ret.Get(0).(func(int) []models.Image); ok {
		r0 = rf(postID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Image)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(postID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewImagesGetterFromPost creates a new instance of ImagesGetterFromPost. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewImagesGetterFromPost(t interface {
	mock.TestingT
	Cleanup(func())
}) *ImagesGetterFromPost {
	mock := &ImagesGetterFromPost{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}