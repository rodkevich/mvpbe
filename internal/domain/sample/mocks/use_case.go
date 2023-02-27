// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/rodkevich/mvpbe/internal/domain/sample/model"
	mock "github.com/stretchr/testify/mock"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

// AllDatabases provides a mock function with given fields: ctx
func (_m *UseCase) AllDatabases(ctx context.Context) ([]string, error) {
	ret := _m.Called(ctx)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateItem provides a mock function with given fields: ctx, m
func (_m *UseCase) CreateItem(ctx context.Context, m *model.SampleItem) error {
	ret := _m.Called(ctx, m)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.SampleItem) error); ok {
		r0 = rf(ctx, m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetItem provides a mock function with given fields: ctx, id
func (_m *UseCase) GetItem(ctx context.Context, id string) (*model.SampleItem, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.SampleItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.SampleItem, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.SampleItem); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SampleItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListItems provides a mock function with given fields: ctx
func (_m *UseCase) ListItems(ctx context.Context) ([]*model.SampleItem, error) {
	ret := _m.Called(ctx)

	var r0 []*model.SampleItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.SampleItem, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.SampleItem); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.SampleItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Readiness provides a mock function with given fields:
func (_m *UseCase) Readiness() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateItem provides a mock function with given fields: ctx, m
func (_m *UseCase) UpdateItem(ctx context.Context, m *model.SampleItem) error {
	ret := _m.Called(ctx, m)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.SampleItem) error); ok {
		r0 = rf(ctx, m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUseCase interface {
	mock.TestingT
	Cleanup(func())
}

// NewUseCase creates a new instance of UseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUseCase(t mockConstructorTestingTNewUseCase) *UseCase {
	mock := &UseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
