// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/site_info.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/site_info.go -destination=mock/service/site_info_mock.go -package=service
//

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	schema "github.com/lantonster/askme/internal/schema"
	gomock "go.uber.org/mock/gomock"
)

// MockSiteInfoService is a mock of SiteInfoService interface.
type MockSiteInfoService struct {
	ctrl     *gomock.Controller
	recorder *MockSiteInfoServiceMockRecorder
	isgomock struct{}
}

// MockSiteInfoServiceMockRecorder is the mock recorder for MockSiteInfoService.
type MockSiteInfoServiceMockRecorder struct {
	mock *MockSiteInfoService
}

// NewMockSiteInfoService creates a new mock instance.
func NewMockSiteInfoService(ctrl *gomock.Controller) *MockSiteInfoService {
	mock := &MockSiteInfoService{ctrl: ctrl}
	mock.recorder = &MockSiteInfoServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSiteInfoService) EXPECT() *MockSiteInfoServiceMockRecorder {
	return m.recorder
}

// GetSiteGeneral mocks base method.
func (m *MockSiteInfoService) GetSiteGeneral(c context.Context) (*schema.GetSiteGeneralRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSiteGeneral", c)
	ret0, _ := ret[0].(*schema.GetSiteGeneralRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSiteGeneral indicates an expected call of GetSiteGeneral.
func (mr *MockSiteInfoServiceMockRecorder) GetSiteGeneral(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSiteGeneral", reflect.TypeOf((*MockSiteInfoService)(nil).GetSiteGeneral), c)
}

// GetSiteLogin mocks base method.
func (m *MockSiteInfoService) GetSiteLogin(c context.Context) (*schema.GetSiteLoginRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSiteLogin", c)
	ret0, _ := ret[0].(*schema.GetSiteLoginRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSiteLogin indicates an expected call of GetSiteLogin.
func (mr *MockSiteInfoServiceMockRecorder) GetSiteLogin(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSiteLogin", reflect.TypeOf((*MockSiteInfoService)(nil).GetSiteLogin), c)
}

// GetSiteUsers mocks base method.
func (m *MockSiteInfoService) GetSiteUsers(c context.Context) (*schema.GetSiteUsersRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSiteUsers", c)
	ret0, _ := ret[0].(*schema.GetSiteUsersRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSiteUsers indicates an expected call of GetSiteUsers.
func (mr *MockSiteInfoServiceMockRecorder) GetSiteUsers(c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSiteUsers", reflect.TypeOf((*MockSiteInfoService)(nil).GetSiteUsers), c)
}
