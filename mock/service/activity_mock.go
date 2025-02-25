// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/activity.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/lantonster/askme/internal/model"
)

// MockActivityService is a mock of ActivityService interface.
type MockActivityService struct {
	ctrl     *gomock.Controller
	recorder *MockActivityServiceMockRecorder
}

// MockActivityServiceMockRecorder is the mock recorder for MockActivityService.
type MockActivityServiceMockRecorder struct {
	mock *MockActivityService
}

// NewMockActivityService creates a new mock instance.
func NewMockActivityService(ctrl *gomock.Controller) *MockActivityService {
	mock := &MockActivityService{ctrl: ctrl}
	mock.recorder = &MockActivityServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockActivityService) EXPECT() *MockActivityServiceMockRecorder {
	return m.recorder
}

// ActivateUser mocks base method.
func (m *MockActivityService) ActivateUser(c context.Context, user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ActivateUser", c, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// ActivateUser indicates an expected call of ActivateUser.
func (mr *MockActivityServiceMockRecorder) ActivateUser(c, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ActivateUser", reflect.TypeOf((*MockActivityService)(nil).ActivateUser), c, user)
}
