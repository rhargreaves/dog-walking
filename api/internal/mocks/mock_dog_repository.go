// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	domain "github.com/rhargreaves/dog-walking/api/internal/dogs/domain"
	mock "github.com/stretchr/testify/mock"
)

// DogRepository is an autogenerated mock type for the DogRepository type
type DogRepository struct {
	mock.Mock
}

type DogRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *DogRepository) EXPECT() *DogRepository_Expecter {
	return &DogRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: dog
func (_m *DogRepository) Create(dog *domain.Dog) error {
	ret := _m.Called(dog)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Dog) error); ok {
		r0 = rf(dog)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DogRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type DogRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - dog *domain.Dog
func (_e *DogRepository_Expecter) Create(dog interface{}) *DogRepository_Create_Call {
	return &DogRepository_Create_Call{Call: _e.mock.On("Create", dog)}
}

func (_c *DogRepository_Create_Call) Run(run func(dog *domain.Dog)) *DogRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.Dog))
	})
	return _c
}

func (_c *DogRepository_Create_Call) Return(_a0 error) *DogRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DogRepository_Create_Call) RunAndReturn(run func(*domain.Dog) error) *DogRepository_Create_Call {
	_c.Call.Return(run)
	return _c
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

// DogRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type DogRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - id string
func (_e *DogRepository_Expecter) Delete(id interface{}) *DogRepository_Delete_Call {
	return &DogRepository_Delete_Call{Call: _e.mock.On("Delete", id)}
}

func (_c *DogRepository_Delete_Call) Run(run func(id string)) *DogRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *DogRepository_Delete_Call) Return(_a0 error) *DogRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DogRepository_Delete_Call) RunAndReturn(run func(string) error) *DogRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: id
func (_m *DogRepository) Get(id string) (*domain.Dog, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *domain.Dog
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Dog, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Dog); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Dog)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DogRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type DogRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - id string
func (_e *DogRepository_Expecter) Get(id interface{}) *DogRepository_Get_Call {
	return &DogRepository_Get_Call{Call: _e.mock.On("Get", id)}
}

func (_c *DogRepository_Get_Call) Run(run func(id string)) *DogRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *DogRepository_Get_Call) Return(_a0 *domain.Dog, _a1 error) *DogRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DogRepository_Get_Call) RunAndReturn(run func(string) (*domain.Dog, error)) *DogRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: limit, name, nextToken
func (_m *DogRepository) List(limit int, name string, nextToken string) (*domain.DogList, error) {
	ret := _m.Called(limit, name, nextToken)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 *domain.DogList
	var r1 error
	if rf, ok := ret.Get(0).(func(int, string, string) (*domain.DogList, error)); ok {
		return rf(limit, name, nextToken)
	}
	if rf, ok := ret.Get(0).(func(int, string, string) *domain.DogList); ok {
		r0 = rf(limit, name, nextToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.DogList)
		}
	}

	if rf, ok := ret.Get(1).(func(int, string, string) error); ok {
		r1 = rf(limit, name, nextToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DogRepository_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type DogRepository_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - limit int
//   - name string
//   - nextToken string
func (_e *DogRepository_Expecter) List(limit interface{}, name interface{}, nextToken interface{}) *DogRepository_List_Call {
	return &DogRepository_List_Call{Call: _e.mock.On("List", limit, name, nextToken)}
}

func (_c *DogRepository_List_Call) Run(run func(limit int, name string, nextToken string)) *DogRepository_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *DogRepository_List_Call) Return(_a0 *domain.DogList, _a1 error) *DogRepository_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DogRepository_List_Call) RunAndReturn(run func(int, string, string) (*domain.DogList, error)) *DogRepository_List_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: id, dog
func (_m *DogRepository) Update(id string, dog *domain.Dog) error {
	ret := _m.Called(id, dog)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *domain.Dog) error); ok {
		r0 = rf(id, dog)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DogRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type DogRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - id string
//   - dog *domain.Dog
func (_e *DogRepository_Expecter) Update(id interface{}, dog interface{}) *DogRepository_Update_Call {
	return &DogRepository_Update_Call{Call: _e.mock.On("Update", id, dog)}
}

func (_c *DogRepository_Update_Call) Run(run func(id string, dog *domain.Dog)) *DogRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*domain.Dog))
	})
	return _c
}

func (_c *DogRepository_Update_Call) Return(_a0 error) *DogRepository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DogRepository_Update_Call) RunAndReturn(run func(string, *domain.Dog) error) *DogRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePhotoHash provides a mock function with given fields: id, photoHash
func (_m *DogRepository) UpdatePhotoHash(id string, photoHash string) error {
	ret := _m.Called(id, photoHash)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePhotoHash")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(id, photoHash)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DogRepository_UpdatePhotoHash_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePhotoHash'
type DogRepository_UpdatePhotoHash_Call struct {
	*mock.Call
}

// UpdatePhotoHash is a helper method to define mock.On call
//   - id string
//   - photoHash string
func (_e *DogRepository_Expecter) UpdatePhotoHash(id interface{}, photoHash interface{}) *DogRepository_UpdatePhotoHash_Call {
	return &DogRepository_UpdatePhotoHash_Call{Call: _e.mock.On("UpdatePhotoHash", id, photoHash)}
}

func (_c *DogRepository_UpdatePhotoHash_Call) Run(run func(id string, photoHash string)) *DogRepository_UpdatePhotoHash_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *DogRepository_UpdatePhotoHash_Call) Return(_a0 error) *DogRepository_UpdatePhotoHash_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DogRepository_UpdatePhotoHash_Call) RunAndReturn(run func(string, string) error) *DogRepository_UpdatePhotoHash_Call {
	_c.Call.Return(run)
	return _c
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
