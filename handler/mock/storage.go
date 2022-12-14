// Code generated by MockGen. DO NOT EDIT.
// Source: ../storage/storage.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPayout is a mock of Payout interface
type MockPayout struct {
	ctrl     *gomock.Controller
	recorder *MockPayoutMockRecorder
}

// MockPayoutMockRecorder is the mock recorder for MockPayout
type MockPayoutMockRecorder struct {
	mock *MockPayout
}

// NewMockPayout creates a new mock instance
func NewMockPayout(ctrl *gomock.Controller) *MockPayout {
	mock := &MockPayout{ctrl: ctrl}
	mock.recorder = &MockPayoutMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPayout) EXPECT() *MockPayoutMockRecorder {
	return m.recorder
}

// Register mocks base method
func (m *MockPayout) Register(ctx context.Context, userID string, reqID int64, payout float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, userID, reqID, payout)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register
func (mr *MockPayoutMockRecorder) Register(ctx, userID, reqID, payout interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockPayout)(nil).Register), ctx, userID, reqID, payout)
}

// Count mocks base method
func (m *MockPayout) Count(ctx context.Context, userID string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, userID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count
func (mr *MockPayoutMockRecorder) Count(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockPayout)(nil).Count), ctx, userID)
}
