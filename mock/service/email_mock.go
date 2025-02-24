// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/email.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/lantonster/askme/internal/model"
	schema "github.com/lantonster/askme/internal/schema"
)

// MockEmailService is a mock of EmailService interface.
type MockEmailService struct {
	ctrl     *gomock.Controller
	recorder *MockEmailServiceMockRecorder
}

// MockEmailServiceMockRecorder is the mock recorder for MockEmailService.
type MockEmailServiceMockRecorder struct {
	mock *MockEmailService
}

// NewMockEmailService creates a new mock instance.
func NewMockEmailService(ctrl *gomock.Controller) *MockEmailService {
	mock := &MockEmailService{ctrl: ctrl}
	mock.recorder = &MockEmailServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailService) EXPECT() *MockEmailServiceMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockEmailService) Send(c context.Context, email, subject, body string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", c, email, subject, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockEmailServiceMockRecorder) Send(c, email, subject, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockEmailService)(nil).Send), c, email, subject, body)
}

// SendRegisterVerificationEmail mocks base method.
func (m *MockEmailService) SendRegisterVerificationEmail(c context.Context, userId int64, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendRegisterVerificationEmail", c, userId, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendRegisterVerificationEmail indicates an expected call of SendRegisterVerificationEmail.
func (mr *MockEmailServiceMockRecorder) SendRegisterVerificationEmail(c, userId, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendRegisterVerificationEmail", reflect.TypeOf((*MockEmailService)(nil).SendRegisterVerificationEmail), c, userId, email)
}

// VerifyUrlExpired mocks base method.
func (m *MockEmailService) VerifyUrlExpired(c context.Context, code string) (*model.VerificationEmail, *schema.ForbiddenRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyUrlExpired", c, code)
	ret0, _ := ret[0].(*model.VerificationEmail)
	ret1, _ := ret[1].(*schema.ForbiddenRes)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// VerifyUrlExpired indicates an expected call of VerifyUrlExpired.
func (mr *MockEmailServiceMockRecorder) VerifyUrlExpired(c, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyUrlExpired", reflect.TypeOf((*MockEmailService)(nil).VerifyUrlExpired), c, code)
}
