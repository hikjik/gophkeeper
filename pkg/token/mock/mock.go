// Code generated by MockGen. DO NOT EDIT.
// Source: token.go

// Package mock_token is a generated GoMock package.
package mock_token

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	token "github.com/go-developer-ya-practicum/gophkeeper/pkg/token"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockManager) Create(userID int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockManagerMockRecorder) Create(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockManager)(nil).Create), userID)
}

// Validate mocks base method.
func (m *MockManager) Validate(accessToken string) (*token.Payload, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", accessToken)
	ret0, _ := ret[0].(*token.Payload)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Validate indicates an expected call of Validate.
func (mr *MockManagerMockRecorder) Validate(accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockManager)(nil).Validate), accessToken)
}