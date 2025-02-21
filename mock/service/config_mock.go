// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/config.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/config.go -destination=mock/service/config_mock.go -package=service
//

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	model "github.com/lantonster/askme/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockConfigService is a mock of ConfigService interface.
type MockConfigService struct {
	ctrl     *gomock.Controller
	recorder *MockConfigServiceMockRecorder
	isgomock struct{}
}

// MockConfigServiceMockRecorder is the mock recorder for MockConfigService.
type MockConfigServiceMockRecorder struct {
	mock *MockConfigService
}

// NewMockConfigService creates a new mock instance.
func NewMockConfigService(ctrl *gomock.Controller) *MockConfigService {
	mock := &MockConfigService{ctrl: ctrl}
	mock.recorder = &MockConfigServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigService) EXPECT() *MockConfigServiceMockRecorder {
	return m.recorder
}

// GetEmail mocks base method.
func (m *MockConfigService) GetEmail(c context.Context) (*model.ConfigEmail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmail", c)
	ret0, _ := ret[0].(*model.ConfigEmail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmail indicates an expected call of GetEmail.
func (mr *MockConfigServiceMockRecorder) GetEmail(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmail", reflect.TypeOf((*MockConfigService)(nil).GetEmail), c)
}
