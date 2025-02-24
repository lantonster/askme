// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repo/user.go

// Package repo is a generated GoMock package.
package repo

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/lantonster/askme/internal/model"
)

// MockUserRepo is a mock of UserRepo interface.
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo.
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance.
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserRepo) CreateUser(c context.Context, user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", c, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepoMockRecorder) CreateUser(c, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepo)(nil).CreateUser), c, user)
}

// GenerateUniqueUsername mocks base method.
func (m *MockUserRepo) GenerateUniqueUsername(c context.Context, username string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateUniqueUsername", c, username)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateUniqueUsername indicates an expected call of GenerateUniqueUsername.
func (mr *MockUserRepoMockRecorder) GenerateUniqueUsername(c, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateUniqueUsername", reflect.TypeOf((*MockUserRepo)(nil).GenerateUniqueUsername), c, username)
}

// GetUserByEmail mocks base method.
func (m *MockUserRepo) GetUserByEmail(c context.Context, email string) (*model.User, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", c, email)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserRepoMockRecorder) GetUserByEmail(c, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserRepo)(nil).GetUserByEmail), c, email)
}

// GetUserByUsername mocks base method.
func (m *MockUserRepo) GetUserByUsername(c context.Context, username string) (*model.User, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", c, username)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockUserRepoMockRecorder) GetUserByUsername(c, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockUserRepo)(nil).GetUserByUsername), c, username)
}

// IncrRank mocks base method.
func (m *MockUserRepo) IncrRank(c context.Context, userId, currentRank, deltaRank int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrRank", c, userId, currentRank, deltaRank)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrRank indicates an expected call of IncrRank.
func (mr *MockUserRepoMockRecorder) IncrRank(c, userId, currentRank, deltaRank interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrRank", reflect.TypeOf((*MockUserRepo)(nil).IncrRank), c, userId, currentRank, deltaRank)
}

// UpdateEmailStatus mocks base method.
func (m *MockUserRepo) UpdateEmailStatus(c context.Context, userId int64, emailStatus string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmailStatus", c, userId, emailStatus)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEmailStatus indicates an expected call of UpdateEmailStatus.
func (mr *MockUserRepoMockRecorder) UpdateEmailStatus(c, userId, emailStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmailStatus", reflect.TypeOf((*MockUserRepo)(nil).UpdateEmailStatus), c, userId, emailStatus)
}
