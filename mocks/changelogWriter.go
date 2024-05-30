// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/AlecSmith96/faceit-user-service/internal/usecases (interfaces: ChangelogWriter)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=../../mocks/changelogWriter.go . ChangelogWriter
//
// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	reflect "reflect"

	entities "github.com/AlecSmith96/faceit-user-service/internal/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockChangelogWriter is a mock of ChangelogWriter interface.
type MockChangelogWriter struct {
	ctrl     *gomock.Controller
	recorder *MockChangelogWriterMockRecorder
}

// MockChangelogWriterMockRecorder is the mock recorder for MockChangelogWriter.
type MockChangelogWriterMockRecorder struct {
	mock *MockChangelogWriter
}

// NewMockChangelogWriter creates a new mock instance.
func NewMockChangelogWriter(ctrl *gomock.Controller) *MockChangelogWriter {
	mock := &MockChangelogWriter{ctrl: ctrl}
	mock.recorder = &MockChangelogWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChangelogWriter) EXPECT() *MockChangelogWriterMockRecorder {
	return m.recorder
}

// PublishChangelogEntry mocks base method.
func (m *MockChangelogWriter) PublishChangelogEntry(arg0 entities.ChangelogEntry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishChangelogEntry", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishChangelogEntry indicates an expected call of PublishChangelogEntry.
func (mr *MockChangelogWriterMockRecorder) PublishChangelogEntry(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishChangelogEntry", reflect.TypeOf((*MockChangelogWriter)(nil).PublishChangelogEntry), arg0)
}
