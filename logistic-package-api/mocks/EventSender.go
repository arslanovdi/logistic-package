// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/arslanovdi/logistic-package/logistic-package-api/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// EventSender is an autogenerated mock type for the EventSender type
type EventSender struct {
	mock.Mock
}

type EventSender_Expecter struct {
	mock *mock.Mock
}

func (_m *EventSender) EXPECT() *EventSender_Expecter {
	return &EventSender_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with no fields
func (_m *EventSender) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventSender_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type EventSender_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *EventSender_Expecter) Close() *EventSender_Close_Call {
	return &EventSender_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *EventSender_Close_Call) Run(run func()) *EventSender_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *EventSender_Close_Call) Return(_a0 error) *EventSender_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EventSender_Close_Call) RunAndReturn(run func() error) *EventSender_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Send provides a mock function with given fields: ctx, pkg, topic
func (_m *EventSender) Send(ctx context.Context, pkg *model.PackageEvent, topic string) error {
	ret := _m.Called(ctx, pkg, topic)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.PackageEvent, string) error); ok {
		r0 = rf(ctx, pkg, topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventSender_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type EventSender_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - ctx context.Context
//   - pkg *model.PackageEvent
//   - topic string
func (_e *EventSender_Expecter) Send(ctx interface{}, pkg interface{}, topic interface{}) *EventSender_Send_Call {
	return &EventSender_Send_Call{Call: _e.mock.On("Send", ctx, pkg, topic)}
}

func (_c *EventSender_Send_Call) Run(run func(ctx context.Context, pkg *model.PackageEvent, topic string)) *EventSender_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.PackageEvent), args[2].(string))
	})
	return _c
}

func (_c *EventSender_Send_Call) Return(_a0 error) *EventSender_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EventSender_Send_Call) RunAndReturn(run func(context.Context, *model.PackageEvent, string) error) *EventSender_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewEventSender creates a new instance of EventSender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventSender(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventSender {
	mock := &EventSender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}