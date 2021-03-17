// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-apps/server/api (interfaces: AppStore)

// Package mock_api is a generated GoMock package.
package mock_api

import (
	gomock "github.com/golang/mock/gomock"
	apps "github.com/mattermost/mattermost-plugin-apps/apps"
	reflect "reflect"
)

// MockAppStore is a mock of AppStore interface
type MockAppStore struct {
	ctrl     *gomock.Controller
	recorder *MockAppStoreMockRecorder
}

// MockAppStoreMockRecorder is the mock recorder for MockAppStore
type MockAppStoreMockRecorder struct {
	mock *MockAppStore
}

// NewMockAppStore creates a new mock instance
func NewMockAppStore(ctrl *gomock.Controller) *MockAppStore {
	mock := &MockAppStore{ctrl: ctrl}
	mock.recorder = &MockAppStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppStore) EXPECT() *MockAppStoreMockRecorder {
	return m.recorder
}

// Delete mocks base method
func (m *MockAppStore) Delete(arg0 *apps.App) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockAppStoreMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAppStore)(nil).Delete), arg0)
}

// Get mocks base method
func (m *MockAppStore) Get(arg0 apps.AppID) (*apps.App, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*apps.App)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockAppStoreMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAppStore)(nil).Get), arg0)
}

// GetAll mocks base method
func (m *MockAppStore) GetAll() []*apps.App {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]*apps.App)
	return ret0
}

// GetAll indicates an expected call of GetAll
func (mr *MockAppStoreMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockAppStore)(nil).GetAll))
}

// Save mocks base method
func (m *MockAppStore) Save(arg0 *apps.App) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockAppStoreMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockAppStore)(nil).Save), arg0)
}