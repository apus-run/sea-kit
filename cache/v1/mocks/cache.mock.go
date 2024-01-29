// Code generated by MockGen. DO NOT EDIT.
// Source: ./cache.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockCache is a mock of Cache interface.
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance.
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// Calc mocks base method.
func (m *MockCache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Calc", ctx, key, step)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Calc indicates an expected call of Calc.
func (mr *MockCacheMockRecorder) Calc(ctx, key, step interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Calc", reflect.TypeOf((*MockCache)(nil).Calc), ctx, key, step)
}

// Contains mocks base method.
func (m *MockCache) Contains(key string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Contains", key)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Contains indicates an expected call of Contains.
func (mr *MockCacheMockRecorder) Contains(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Contains", reflect.TypeOf((*MockCache)(nil).Contains), key)
}

// Decrement mocks base method.
func (m *MockCache) Decrement(ctx context.Context, key string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decrement", ctx, key)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Decrement indicates an expected call of Decrement.
func (mr *MockCacheMockRecorder) Decrement(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decrement", reflect.TypeOf((*MockCache)(nil).Decrement), ctx, key)
}

// Del mocks base method.
func (m *MockCache) Del(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockCacheMockRecorder) Del(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockCache)(nil).Del), ctx, key)
}

// DelMany mocks base method.
func (m *MockCache) DelMany(ctx context.Context, keys []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelMany", ctx, keys)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelMany indicates an expected call of DelMany.
func (mr *MockCacheMockRecorder) DelMany(ctx, keys interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelMany", reflect.TypeOf((*MockCache)(nil).DelMany), ctx, keys)
}

// Get mocks base method.
func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCacheMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCache)(nil).Get), ctx, key)
}

// GetMany mocks base method.
func (m *MockCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMany", ctx, keys)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMany indicates an expected call of GetMany.
func (mr *MockCacheMockRecorder) GetMany(ctx, keys interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMany", reflect.TypeOf((*MockCache)(nil).GetMany), ctx, keys)
}

// GetObj mocks base method.
func (m *MockCache) GetObj(ctx context.Context, key string, model interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObj", ctx, key, model)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetObj indicates an expected call of GetObj.
func (mr *MockCacheMockRecorder) GetObj(ctx, key, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObj", reflect.TypeOf((*MockCache)(nil).GetObj), ctx, key, model)
}

// GetTTL mocks base method.
func (m *MockCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTTL", ctx, key)
	ret0, _ := ret[0].(time.Duration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTTL indicates an expected call of GetTTL.
func (mr *MockCacheMockRecorder) GetTTL(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTTL", reflect.TypeOf((*MockCache)(nil).GetTTL), ctx, key)
}

// Increment mocks base method.
func (m *MockCache) Increment(ctx context.Context, key string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Increment", ctx, key)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Increment indicates an expected call of Increment.
func (mr *MockCacheMockRecorder) Increment(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Increment", reflect.TypeOf((*MockCache)(nil).Increment), ctx, key)
}

// Set mocks base method.
func (m *MockCache) Set(ctx context.Context, key, val string, timeout time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, val, timeout)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockCacheMockRecorder) Set(ctx, key, val, timeout interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCache)(nil).Set), ctx, key, val, timeout)
}

// SetForever mocks base method.
func (m *MockCache) SetForever(ctx context.Context, key, val string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetForever", ctx, key, val)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetForever indicates an expected call of SetForever.
func (mr *MockCacheMockRecorder) SetForever(ctx, key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetForever", reflect.TypeOf((*MockCache)(nil).SetForever), ctx, key, val)
}

// SetForeverObj mocks base method.
func (m *MockCache) SetForeverObj(ctx context.Context, key string, val interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetForeverObj", ctx, key, val)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetForeverObj indicates an expected call of SetForeverObj.
func (mr *MockCacheMockRecorder) SetForeverObj(ctx, key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetForeverObj", reflect.TypeOf((*MockCache)(nil).SetForeverObj), ctx, key, val)
}

// SetMany mocks base method.
func (m *MockCache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMany", ctx, data, timeout)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetMany indicates an expected call of SetMany.
func (mr *MockCacheMockRecorder) SetMany(ctx, data, timeout interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMany", reflect.TypeOf((*MockCache)(nil).SetMany), ctx, data, timeout)
}

// SetObj mocks base method.
func (m *MockCache) SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetObj", ctx, key, val, timeout)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetObj indicates an expected call of SetObj.
func (mr *MockCacheMockRecorder) SetObj(ctx, key, val, timeout interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetObj", reflect.TypeOf((*MockCache)(nil).SetObj), ctx, key, val, timeout)
}

// SetTTL mocks base method.
func (m *MockCache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTTL", ctx, key, timeout)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTTL indicates an expected call of SetTTL.
func (mr *MockCacheMockRecorder) SetTTL(ctx, key, timeout interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTTL", reflect.TypeOf((*MockCache)(nil).SetTTL), ctx, key, timeout)
}