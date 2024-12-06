// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/arslanovdi/logistic-package/logistic-package-api/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// EventRepo is an autogenerated mock type for the EventRepo type
type EventRepo struct {
	mock.Mock
}

type EventRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *EventRepo) EXPECT() *EventRepo_Expecter {
	return &EventRepo_Expecter{mock: &_m.Mock}
}

// Lock provides a mock function with given fields: ctx, n
func (_m *EventRepo) Lock(ctx context.Context, n int) ([]model.PackageEvent, error) {
	ret := _m.Called(ctx, n)

	if len(ret) == 0 {
		panic("no return value specified for Lock")
	}

	var r0 []model.PackageEvent
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]model.PackageEvent, error)); ok {
		return rf(ctx, n)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []model.PackageEvent); ok {
		r0 = rf(ctx, n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.PackageEvent)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, n)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EventRepo_Lock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Lock'
type EventRepo_Lock_Call struct {
	*mock.Call
}

// Lock is a helper method to define mock.On call
//   - ctx context.Context
//   - n int
func (_e *EventRepo_Expecter) Lock(ctx interface{}, n interface{}) *EventRepo_Lock_Call {
	return &EventRepo_Lock_Call{Call: _e.mock.On("Lock", ctx, n)}
}

func (_c *EventRepo_Lock_Call) Run(run func(ctx context.Context, n int)) *EventRepo_Lock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *EventRepo_Lock_Call) Return(_a0 []model.PackageEvent, _a1 error) *EventRepo_Lock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *EventRepo_Lock_Call) RunAndReturn(run func(context.Context, int) ([]model.PackageEvent, error)) *EventRepo_Lock_Call {
	_c.Call.Return(run)
	return _c
}

// Remove provides a mock function with given fields: ctx, eventIDs
func (_m *EventRepo) Remove(ctx context.Context, eventIDs []int64) error {
	ret := _m.Called(ctx, eventIDs)

	if len(ret) == 0 {
		panic("no return value specified for Remove")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []int64) error); ok {
		r0 = rf(ctx, eventIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventRepo_Remove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remove'
type EventRepo_Remove_Call struct {
	*mock.Call
}

// Remove is a helper method to define mock.On call
//   - ctx context.Context
//   - eventIDs []int64
func (_e *EventRepo_Expecter) Remove(ctx interface{}, eventIDs interface{}) *EventRepo_Remove_Call {
	return &EventRepo_Remove_Call{Call: _e.mock.On("Remove", ctx, eventIDs)}
}

func (_c *EventRepo_Remove_Call) Run(run func(ctx context.Context, eventIDs []int64)) *EventRepo_Remove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]int64))
	})
	return _c
}

func (_c *EventRepo_Remove_Call) Return(_a0 error) *EventRepo_Remove_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EventRepo_Remove_Call) RunAndReturn(run func(context.Context, []int64) error) *EventRepo_Remove_Call {
	_c.Call.Return(run)
	return _c
}

// Unlock provides a mock function with given fields: ctx, eventID
func (_m *EventRepo) Unlock(ctx context.Context, eventID []int64) error {
	ret := _m.Called(ctx, eventID)

	if len(ret) == 0 {
		panic("no return value specified for Unlock")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []int64) error); ok {
		r0 = rf(ctx, eventID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventRepo_Unlock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unlock'
type EventRepo_Unlock_Call struct {
	*mock.Call
}

// Unlock is a helper method to define mock.On call
//   - ctx context.Context
//   - eventID []int64
func (_e *EventRepo_Expecter) Unlock(ctx interface{}, eventID interface{}) *EventRepo_Unlock_Call {
	return &EventRepo_Unlock_Call{Call: _e.mock.On("Unlock", ctx, eventID)}
}

func (_c *EventRepo_Unlock_Call) Run(run func(ctx context.Context, eventID []int64)) *EventRepo_Unlock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]int64))
	})
	return _c
}

func (_c *EventRepo_Unlock_Call) Return(_a0 error) *EventRepo_Unlock_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EventRepo_Unlock_Call) RunAndReturn(run func(context.Context, []int64) error) *EventRepo_Unlock_Call {
	_c.Call.Return(run)
	return _c
}

// UnlockAll provides a mock function with given fields: ctx
func (_m *EventRepo) UnlockAll(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for UnlockAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventRepo_UnlockAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnlockAll'
type EventRepo_UnlockAll_Call struct {
	*mock.Call
}

// UnlockAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *EventRepo_Expecter) UnlockAll(ctx interface{}) *EventRepo_UnlockAll_Call {
	return &EventRepo_UnlockAll_Call{Call: _e.mock.On("UnlockAll", ctx)}
}

func (_c *EventRepo_UnlockAll_Call) Run(run func(ctx context.Context)) *EventRepo_UnlockAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *EventRepo_UnlockAll_Call) Return(_a0 error) *EventRepo_UnlockAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EventRepo_UnlockAll_Call) RunAndReturn(run func(context.Context) error) *EventRepo_UnlockAll_Call {
	_c.Call.Return(run)
	return _c
}

// NewEventRepo creates a new instance of EventRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventRepo {
	mock := &EventRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
