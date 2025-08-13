package collector

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/mock"
)

// MockGodoClient is a mock implementation of the godo.Client
type MockGodoClient struct {
	mock.Mock
	MockBalance *MockBalanceService
	MockAccount *MockAccountService
}

// MockBalanceService mocks the godo BalanceService
type MockBalanceService struct {
	mock.Mock
}

func (m *MockBalanceService) Get(ctx context.Context) (*godo.Balance, *godo.Response, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*godo.Response), args.Error(2)
	}
	return args.Get(0).(*godo.Balance), args.Get(1).(*godo.Response), args.Error(2)
}

// MockAccountService mocks the godo AccountService  
type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) Get(ctx context.Context) (*godo.Account, *godo.Response, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*godo.Response), args.Error(2)
	}
	return args.Get(0).(*godo.Account), args.Get(1).(*godo.Response), args.Error(2)
}