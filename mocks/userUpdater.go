// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/AlecSmith96/faceit-user-service/internal/usecases (interfaces: UserUpdater)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=../../mocks/userUpdater.go . UserUpdater
//
// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	context "context"
	reflect "reflect"

	entities "github.com/AlecSmith96/faceit-user-service/internal/entities"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockUserUpdater is a mock of UserUpdater interface.
type MockUserUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockUserUpdaterMockRecorder
}

// MockUserUpdaterMockRecorder is the mock recorder for MockUserUpdater.
type MockUserUpdaterMockRecorder struct {
	mock *MockUserUpdater
}

// NewMockUserUpdater creates a new mock instance.
func NewMockUserUpdater(ctrl *gomock.Controller) *MockUserUpdater {
	mock := &MockUserUpdater{ctrl: ctrl}
	mock.recorder = &MockUserUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUpdater) EXPECT() *MockUserUpdaterMockRecorder {
	return m.recorder
}

// UpdateUser mocks base method.
func (m *MockUserUpdater) UpdateUser(arg0 context.Context, arg1 uuid.UUID, arg2, arg3, arg4, arg5, arg6, arg7 string) (*entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	ret0, _ := ret[0].(*entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserUpdaterMockRecorder) UpdateUser(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserUpdater)(nil).UpdateUser), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
}