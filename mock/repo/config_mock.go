// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repo/config.go

// Package repo is a generated GoMock package.
package repo

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/lantonster/askme/internal/model"
)

// MockConfigRepo is a mock of ConfigRepo interface.
type MockConfigRepo struct {
	ctrl     *gomock.Controller
	recorder *MockConfigRepoMockRecorder
}

// MockConfigRepoMockRecorder is the mock recorder for MockConfigRepo.
type MockConfigRepoMockRecorder struct {
	mock *MockConfigRepo
}

// NewMockConfigRepo creates a new mock instance.
func NewMockConfigRepo(ctrl *gomock.Controller) *MockConfigRepo {
	mock := &MockConfigRepo{ctrl: ctrl}
	mock.recorder = &MockConfigRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigRepo) EXPECT() *MockConfigRepoMockRecorder {
	return m.recorder
}

// FirstConfigByKey mocks base method.
func (m *MockConfigRepo) FirstConfigByKey(c context.Context, key string) (*model.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FirstConfigByKey", c, key)
	ret0, _ := ret[0].(*model.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FirstConfigByKey indicates an expected call of FirstConfigByKey.
func (mr *MockConfigRepoMockRecorder) FirstConfigByKey(c, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FirstConfigByKey", reflect.TypeOf((*MockConfigRepo)(nil).FirstConfigByKey), c, key)
}
