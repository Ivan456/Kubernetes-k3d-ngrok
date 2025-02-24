package main

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEthereumClient is a mock implementation of EthereumClientInterface
type MockEthereumClient struct {
	mock.Mock
}

func (m *MockEthereumClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	args := m.Called(ctx, number)
	return args.Get(0).(*types.Header), args.Error(1)
}

func (m *MockEthereumClient) BalanceAt(ctx context.Context, account common.Address, number *big.Int) (*big.Int, error) {
	args := m.Called(ctx, account, number)
	return args.Get(0).(*big.Int), args.Error(1)
}

func TestGetLatestBlockNumber(t *testing.T) {
	mockClient := new(MockEthereumClient)
	expectedBlockNumber := big.NewInt(12345)
	mockClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(&types.Header{Number: expectedBlockNumber}, nil)

	ethClient := NewEthereumClient(mockClient)
	blockNumber, err := ethClient.getLatestBlockNumber()

	assert.NoError(t, err)
	assert.Equal(t, expectedBlockNumber, blockNumber)
	mockClient.AssertExpectations(t)
}

func TestGetBalance(t *testing.T) {
	mockClient := new(MockEthereumClient)
	expectedBalance := big.NewInt(1000)
	address := common.HexToAddress("0x0")
	mockClient.On("BalanceAt", mock.Anything, address, (*big.Int)(nil)).Return(expectedBalance, nil)

	ethClient := NewEthereumClient(mockClient)
	balance, err := ethClient.getBalance(address.Hex())

	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
	mockClient.AssertExpectations(t)
}
