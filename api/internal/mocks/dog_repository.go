// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	models "github.com/rhargreaves/dog-walking/api/internal/dogs/models"
	mock "github.com/stretchr/testify/mock"
)

// DogRepository is an autogenerated mock type for the DogRepository type
type DogRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: dog
func (_m *DogRepository) Create(dog *models.Dog) error {
	ret := _m.Called(dog)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Dog) error); ok {
		r0 = rf(dog)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: id
func (_m *DogRepository) Delete(id string) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: id
func (_m *DogRepository) Get(id string) (*models.Dog, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *models.Dog
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.Dog, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *models.Dog); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Dog)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with no fields
func (_m *DogRepository) List() ([]models.Dog, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []models.Dog
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]models.Dog, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []models.Dog); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Dog)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: id, dog
func (_m *DogRepository) Update(id string, dog *models.Dog) error {
	ret := _m.Called(id, dog)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *models.Dog) error); ok {
		r0 = rf(id, dog)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDogRepository creates a new instance of DogRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDogRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *DogRepository {
	mock := &DogRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
