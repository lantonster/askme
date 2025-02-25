// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/uploads.go

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockUploadsService is a mock of UploadsService interface.
type MockUploadsService struct {
	ctrl     *gomock.Controller
	recorder *MockUploadsServiceMockRecorder
}

// MockUploadsServiceMockRecorder is the mock recorder for MockUploadsService.
type MockUploadsServiceMockRecorder struct {
	mock *MockUploadsService
}

// NewMockUploadsService creates a new mock instance.
func NewMockUploadsService(ctrl *gomock.Controller) *MockUploadsService {
	mock := &MockUploadsService{ctrl: ctrl}
	mock.recorder = &MockUploadsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUploadsService) EXPECT() *MockUploadsServiceMockRecorder {
	return m.recorder
}

// AvatarThumbFile mocks base method.
func (m *MockUploadsService) AvatarThumbFile(c *gin.Context, fileName string, size int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AvatarThumbFile", c, fileName, size)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AvatarThumbFile indicates an expected call of AvatarThumbFile.
func (mr *MockUploadsServiceMockRecorder) AvatarThumbFile(c, fileName, size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AvatarThumbFile", reflect.TypeOf((*MockUploadsService)(nil).AvatarThumbFile), c, fileName, size)
}
