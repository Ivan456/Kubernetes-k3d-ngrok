package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const infuraURL = "https://mainnet.infura.io/v3/c543932173e54a3fbbe7ce8e4d0c1e78"

// EthereumClient wraps an ethclient.Client
type EthereumClient struct {
	client *ethclient.Client
}

func NewEthereumClient() (*EthereumClient, error) {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		return nil, err
	}
	return &EthereumClient{client: client}, nil
}

func (ec *EthereumClient) getLatestBlockNumber() (*big.Int, error) {
	header, err := ec.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

func (ec *EthereumClient) getBalance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := ec.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func main() {
	client, err := NewEthereumClient()
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum: %v", err)
	}

	http.HandleFunc("/latest-block", func(w http.ResponseWriter, r *http.Request) {
		blockNumber, err := client.getLatestBlockNumber()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"latest_block": blockNumber.String()})
	})

	http.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}
		balance, err := client.getBalance(address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"balance": balance.String()})
	})

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
