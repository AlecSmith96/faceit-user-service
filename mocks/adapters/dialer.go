// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/AlecSmith96/faceit-user-service/internal/adapters (interfaces: Dialer)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=../../mocks/adapters/dialer.go . Dialer
//
// Package mock_adapters is a generated GoMock package.
package mock_adapters

import (
	context "context"
	reflect "reflect"

	adapters "github.com/AlecSmith96/faceit-user-service/internal/adapters"
	gomock "go.uber.org/mock/gomock"
)

// MockDialer is a mock of Dialer interface.
type MockDialer struct {
	ctrl     *gomock.Controller
	recorder *MockDialerMockRecorder
}

// MockDialerMockRecorder is the mock recorder for MockDialer.
type MockDialerMockRecorder struct {
	mock *MockDialer
}

// NewMockDialer creates a new mock instance.
func NewMockDialer(ctrl *gomock.Controller) *MockDialer {
	mock := &MockDialer{ctrl: ctrl}
	mock.recorder = &MockDialerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDialer) EXPECT() *MockDialerMockRecorder {
	return m.recorder
}

// DialLeader mocks base method.
func (m *MockDialer) DialLeader(arg0 context.Context, arg1, arg2, arg3 string, arg4 int) (adapters.KafkaConnection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DialLeader", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(adapters.KafkaConnection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DialLeader indicates an expected call of DialLeader.
func (mr *MockDialerMockRecorder) DialLeader(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DialLeader", reflect.TypeOf((*MockDialer)(nil).DialLeader), arg0, arg1, arg2, arg3, arg4)
}
