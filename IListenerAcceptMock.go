// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bhbosman/goCommsNetListener (interfaces: IListenerAccept)

// Package goCommsNetListener is a generated GoMock package.
package goCommsNetListener

import (
	context "context"
	net "net"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIListenerAccept is a mock of IListenerAccept interface.
type MockIListenerAccept struct {
	ctrl     *gomock.Controller
	recorder *MockIListenerAcceptMockRecorder
}

// MockIListenerAcceptMockRecorder is the mock recorder for MockIListenerAccept.
type MockIListenerAcceptMockRecorder struct {
	mock *MockIListenerAccept
}

// NewMockIListenerAccept creates a new mock instance.
func NewMockIListenerAccept(ctrl *gomock.Controller) *MockIListenerAccept {
	mock := &MockIListenerAccept{ctrl: ctrl}
	mock.recorder = &MockIListenerAcceptMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIListenerAccept) EXPECT() *MockIListenerAcceptMockRecorder {
	return m.recorder
}

// AcceptWithContext mocks base method.
func (m *MockIListenerAccept) AcceptWithContext() (net.Conn, context.CancelFunc, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptWithContext")
	ret0, _ := ret[0].(net.Conn)
	ret1, _ := ret[1].(context.CancelFunc)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AcceptWithContext indicates an expected call of AcceptWithContext.
func (mr *MockIListenerAcceptMockRecorder) AcceptWithContext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptWithContext", reflect.TypeOf((*MockIListenerAccept)(nil).AcceptWithContext))
}

// argNames: []
// defaultArgs: []
// defaultArgsAsString:
// argTypes: []
// argString:
// rets: [net.Conn context.CancelFunc error]
// retString: net.Conn, context.CancelFunc, error
// retString:  (net.Conn, context.CancelFunc, error)
// ia: map[]
// idRecv: mr
// 0
func (mr *MockIListenerAcceptMockRecorder) OnAcceptWithContextDoAndReturn(
	f func() (net.Conn, context.CancelFunc, error)) *gomock.Call {
	return mr.
		AcceptWithContext().
		DoAndReturn(f)
}

// 0
func (mr *MockIListenerAcceptMockRecorder) OnAcceptWithContextDo(
	f func()) *gomock.Call {
	return mr.
		AcceptWithContext().
		DoAndReturn(f)
}

// retNames: [ret0 ret1 ret2]
// retArgs: [ret0 net.Conn ret1 context.CancelFunc ret2 error]
// retArgs22: ret0 net.Conn,ret1 context.CancelFunc,ret2 error
// 1
func (mr *MockIListenerAcceptMockRecorder) OnAcceptWithContextReturn(ret0 net.Conn, ret1 context.CancelFunc, ret2 error) *gomock.Call {
	return mr.
		AcceptWithContext().
		Return(ret0, ret1, ret2)
}
