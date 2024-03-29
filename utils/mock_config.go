// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/andrewelkin/trilib/utils (interfaces: IConfig)

// Package utils is a generated GoMock package.
package utils

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIConfig is a mock of IConfig interface.
type MockIConfig struct {
	ctrl     *gomock.Controller
	recorder *MockIConfigMockRecorder
}

// MockIConfigMockRecorder is the mock recorder for MockIConfig.
type MockIConfigMockRecorder struct {
	mock *MockIConfig
}

// NewMockIConfig creates a new mock instance.
func NewMockIConfig(ctrl *gomock.Controller) *MockIConfig {
	mock := &MockIConfig{ctrl: ctrl}
	mock.recorder = &MockIConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIConfig) EXPECT() *MockIConfigMockRecorder {
	return m.recorder
}

// FromKey mocks base method.
func (m *MockIConfig) FromKey(arg0 string) IConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FromKey", arg0)
	ret0, _ := ret[0].(IConfig)
	return ret0
}

// FromKey indicates an expected call of FromKey.
func (mr *MockIConfigMockRecorder) FromKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FromKey", reflect.TypeOf((*MockIConfig)(nil).FromKey), arg0)
}

// GetBoolDefault mocks base method.
func (m *MockIConfig) GetBoolDefault(arg0 string, arg1 bool) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoolDefault", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetBoolDefault indicates an expected call of GetBoolDefault.
func (mr *MockIConfigMockRecorder) GetBoolDefault(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoolDefault", reflect.TypeOf((*MockIConfig)(nil).GetBoolDefault), arg0, arg1)
}

// GetCfg mocks base method.
func (m *MockIConfig) GetCfg() map[string]interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCfg")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// GetCfg indicates an expected call of GetCfg.
func (mr *MockIConfigMockRecorder) GetCfg() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCfg", reflect.TypeOf((*MockIConfig)(nil).GetCfg))
}

// GetFileName mocks base method.
func (m *MockIConfig) GetFileName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetFileName indicates an expected call of GetFileName.
func (mr *MockIConfigMockRecorder) GetFileName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileName", reflect.TypeOf((*MockIConfig)(nil).GetFileName))
}

// GetFloatDefault mocks base method.
func (m *MockIConfig) GetFloatDefault(arg0 string, arg1 float64) float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFloatDefault", arg0, arg1)
	ret0, _ := ret[0].(float64)
	return ret0
}

// GetFloatDefault indicates an expected call of GetFloatDefault.
func (mr *MockIConfigMockRecorder) GetFloatDefault(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFloatDefault", reflect.TypeOf((*MockIConfig)(nil).GetFloatDefault), arg0, arg1)
}

// GetIntDefault mocks base method.
func (m *MockIConfig) GetIntDefault(arg0 string, arg1 int64) int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIntDefault", arg0, arg1)
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetIntDefault indicates an expected call of GetIntDefault.
func (mr *MockIConfigMockRecorder) GetIntDefault(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIntDefault", reflect.TypeOf((*MockIConfig)(nil).GetIntDefault), arg0, arg1)
}

// GetRO mocks base method.
func (m *MockIConfig) GetRO() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRO")
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetRO indicates an expected call of GetRO.
func (mr *MockIConfigMockRecorder) GetRO() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRO", reflect.TypeOf((*MockIConfig)(nil).GetRO))
}

// GetString mocks base method.
func (m *MockIConfig) GetString(arg0 string) *string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", arg0)
	ret0, _ := ret[0].(*string)
	return ret0
}

// GetString indicates an expected call of GetString.
func (mr *MockIConfigMockRecorder) GetString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockIConfig)(nil).GetString), arg0)
}

// GetStringDefault mocks base method.
func (m *MockIConfig) GetStringDefault(arg0, arg1 string) *string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStringDefault", arg0, arg1)
	ret0, _ := ret[0].(*string)
	return ret0
}

// GetStringDefault indicates an expected call of GetStringDefault.
func (mr *MockIConfigMockRecorder) GetStringDefault(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStringDefault", reflect.TypeOf((*MockIConfig)(nil).GetStringDefault), arg0, arg1)
}

// GetStringList mocks base method.
func (m *MockIConfig) GetStringList(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStringList", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetStringList indicates an expected call of GetStringList.
func (mr *MockIConfigMockRecorder) GetStringList(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStringList", reflect.TypeOf((*MockIConfig)(nil).GetStringList), arg0)
}

// GetValue mocks base method.
func (m *MockIConfig) GetValue(arg0 string) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValue", arg0)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// GetValue indicates an expected call of GetValue.
func (mr *MockIConfigMockRecorder) GetValue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValue", reflect.TypeOf((*MockIConfig)(nil).GetValue), arg0)
}

// ModifiedQ mocks base method.
func (m *MockIConfig) ModifiedQ() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModifiedQ")
	ret0, _ := ret[0].(bool)
	return ret0
}

// ModifiedQ indicates an expected call of ModifiedQ.
func (mr *MockIConfigMockRecorder) ModifiedQ() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModifiedQ", reflect.TypeOf((*MockIConfig)(nil).ModifiedQ))
}

// ReadConfig mocks base method.
func (m *MockIConfig) ReadConfig(arg0 string) IConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadConfig", arg0)
	ret0, _ := ret[0].(IConfig)
	return ret0
}

// ReadConfig indicates an expected call of ReadConfig.
func (mr *MockIConfigMockRecorder) ReadConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadConfig", reflect.TypeOf((*MockIConfig)(nil).ReadConfig), arg0)
}

// SetFileName mocks base method.
func (m *MockIConfig) SetFileName(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetFileName", arg0)
}

// SetFileName indicates an expected call of SetFileName.
func (mr *MockIConfigMockRecorder) SetFileName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFileName", reflect.TypeOf((*MockIConfig)(nil).SetFileName), arg0)
}

// SetRO mocks base method.
func (m *MockIConfig) SetRO(arg0 bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRO", arg0)
}

// SetRO indicates an expected call of SetRO.
func (mr *MockIConfigMockRecorder) SetRO(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRO", reflect.TypeOf((*MockIConfig)(nil).SetRO), arg0)
}

// WriteConfigX mocks base method.
func (m *MockIConfig) WriteConfigX() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteConfigX")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteConfigX indicates an expected call of WriteConfigX.
func (mr *MockIConfigMockRecorder) WriteConfigX() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteConfigX", reflect.TypeOf((*MockIConfig)(nil).WriteConfigX))
}
