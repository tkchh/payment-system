// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	models "infotecsTest/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// BalanceReceiver is an autogenerated mock type for the BalanceReceiver type
type BalanceReceiver struct {
	mock.Mock
}

// GetWalletBalance provides a mock function with given fields: address
func (_m *BalanceReceiver) GetWalletBalance(address string) (models.Wallet, error) {
	ret := _m.Called(address)

	if len(ret) == 0 {
		panic("no return value specified for GetWalletBalance")
	}

	var r0 models.Wallet
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (models.Wallet, error)); ok {
		return rf(address)
	}
	if rf, ok := ret.Get(0).(func(string) models.Wallet); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(models.Wallet)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBalanceReceiver creates a new instance of BalanceReceiver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBalanceReceiver(t interface {
	mock.TestingT
	Cleanup(func())
}) *BalanceReceiver {
	mock := &BalanceReceiver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
