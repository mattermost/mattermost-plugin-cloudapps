// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-apps/server/api (interfaces: Proxy)

// Package mock_api is a generated GoMock package.
package mock_api

import (
	gomock "github.com/golang/mock/gomock"
	apps "github.com/mattermost/mattermost-plugin-apps/apps"
	api "github.com/mattermost/mattermost-plugin-apps/server/api"
	io "io"
	reflect "reflect"
)

// MockProxy is a mock of Proxy interface
type MockProxy struct {
	ctrl     *gomock.Controller
	recorder *MockProxyMockRecorder
}

// MockProxyMockRecorder is the mock recorder for MockProxy
type MockProxyMockRecorder struct {
	mock *MockProxy
}

// NewMockProxy creates a new mock instance
func NewMockProxy(ctrl *gomock.Controller) *MockProxy {
	mock := &MockProxy{ctrl: ctrl}
	mock.recorder = &MockProxyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProxy) EXPECT() *MockProxyMockRecorder {
	return m.recorder
}

// Call mocks base method
func (m *MockProxy) Call(arg0 apps.SessionToken, arg1 *apps.CallRequest) *apps.CallResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", arg0, arg1)
	ret0, _ := ret[0].(*apps.CallResponse)
	return ret0
}

// Call indicates an expected call of Call
func (mr *MockProxyMockRecorder) Call(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockProxy)(nil).Call), arg0, arg1)
}

// GetAsset mocks base method
func (m *MockProxy) GetAsset(arg0 apps.AppID, arg1 string) (io.ReadCloser, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAsset", arg0, arg1)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAsset indicates an expected call of GetAsset
func (mr *MockProxyMockRecorder) GetAsset(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAsset", reflect.TypeOf((*MockProxy)(nil).GetAsset), arg0, arg1)
}

// GetBindings mocks base method
func (m *MockProxy) GetBindings(arg0 *apps.Context) ([]*apps.Binding, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBindings", arg0)
	ret0, _ := ret[0].([]*apps.Binding)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBindings indicates an expected call of GetBindings
func (mr *MockProxyMockRecorder) GetBindings(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBindings", reflect.TypeOf((*MockProxy)(nil).GetBindings), arg0)
}

// Notify mocks base method
func (m *MockProxy) Notify(arg0 *apps.Context, arg1 apps.Subject) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Notify", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Notify indicates an expected call of Notify
func (mr *MockProxyMockRecorder) Notify(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Notify", reflect.TypeOf((*MockProxy)(nil).Notify), arg0, arg1)
}

// ProvisionBuiltIn mocks base method
func (m *MockProxy) ProvisionBuiltIn(arg0 apps.AppID, arg1 api.Upstream) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ProvisionBuiltIn", arg0, arg1)
}

// ProvisionBuiltIn indicates an expected call of ProvisionBuiltIn
func (mr *MockProxyMockRecorder) ProvisionBuiltIn(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProvisionBuiltIn", reflect.TypeOf((*MockProxy)(nil).ProvisionBuiltIn), arg0, arg1)
}