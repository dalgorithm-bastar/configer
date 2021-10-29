// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/configcenter/pkg/repository (interfaces: Storage)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AcidCommit mocks base method.
func (m *MockStorage) AcidCommit(arg0 map[string]string, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcidCommit", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AcidCommit indicates an expected call of AcidCommit.
func (mr *MockStorageMockRecorder) AcidCommit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcidCommit", reflect.TypeOf((*MockStorage)(nil).AcidCommit), arg0, arg1)
}

// Delete mocks base method.
func (m *MockStorage) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStorageMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStorage)(nil).Delete), arg0)
}

// DeletebyPrefix mocks base method.
func (m *MockStorage) DeletebyPrefix(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletebyPrefix", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletebyPrefix indicates an expected call of DeletebyPrefix.
func (mr *MockStorageMockRecorder) DeletebyPrefix(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletebyPrefix", reflect.TypeOf((*MockStorage)(nil).DeletebyPrefix), arg0)
}

// Get mocks base method.
func (m *MockStorage) Get(arg0 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStorageMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStorage)(nil).Get), arg0)
}

// GetSourceDataorOperator mocks base method.
func (m *MockStorage) GetSourceDataorOperator() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSourceDataorOperator")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// GetSourceDataorOperator indicates an expected call of GetSourceDataorOperator.
func (mr *MockStorageMockRecorder) GetSourceDataorOperator() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSourceDataorOperator", reflect.TypeOf((*MockStorage)(nil).GetSourceDataorOperator))
}

// GetbyPrefix mocks base method.
func (m *MockStorage) GetbyPrefix(arg0 string) (map[string][]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetbyPrefix", arg0)
	ret0, _ := ret[0].(map[string][]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetbyPrefix indicates an expected call of GetbyPrefix.
func (mr *MockStorageMockRecorder) GetbyPrefix(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetbyPrefix", reflect.TypeOf((*MockStorage)(nil).GetbyPrefix), arg0)
}

// GracefullyClose mocks base method.
func (m *MockStorage) GracefullyClose(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GracefullyClose", arg0)
}

// GracefullyClose indicates an expected call of GracefullyClose.
func (mr *MockStorageMockRecorder) GracefullyClose(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GracefullyClose", reflect.TypeOf((*MockStorage)(nil).GracefullyClose), arg0)
}

// Put mocks base method.
func (m *MockStorage) Put(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockStorageMockRecorder) Put(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockStorage)(nil).Put), arg0, arg1)
}
