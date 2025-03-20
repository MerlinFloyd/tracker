package blockchain

import (
	"context"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"my-fullstack-app/backend/internal/logger"
	"my-fullstack-app/backend/internal/models"
)

// Error definitions
var ()

// Client represents an Ethereum blockchain client
type Client struct {
	ethClient *ethclient.Client
}

// NewClient creates a new blockchain client
func NewClient() (*Client, error) {
	// Get Infura API key from environment variable
	infuraKey := os.Getenv("INFURA_API_KEY")
	if infuraKey == "" {
		logger.Warn().Msg("INFURA_API_KEY not set, using default key for testing")
		infuraKey = "4b4eabeb1b8b4bfeaa4f29e754f2d282" // Replace with your actual Infura API key for testing
	}

	// Connect to Infura
	infuraURL := "https://mainnet.infura.io/v3/" + infuraKey
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to Infura")
		return nil, err
	}

	logger.Info().Msg("Successfully connected to Infura")
	return &Client{
		ethClient: client,
	}, nil
}

// CheckConnection tests if the Ethereum client is connected
func (c *Client) CheckConnection() error {
	_, err := c.ethClient.BlockNumber(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to check Ethereum client connection")
	}
	return err
}

// GetBlockNumber returns the latest block number
func (c *Client) GetBlockNumber() (uint64, error) {
	blockNumber, err := c.ethClient.BlockNumber(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get the latest block number")
		return 0, err
	}
	logger.Info().Uint64("blockNumber", blockNumber).Msg("Fetched latest block number")
	return blockNumber, nil
}

// GetBalance returns the balance of an Ethereum address in wei
func (c *Client) GetBalance(address string) (*big.Int, error) {
	if !common.IsHexAddress(address) {
		logger.Warn().Str("address", address).Msg("Invalid Ethereum address")
		return nil, ErrInvalidAddress
	}

	account := common.HexToAddress(address)
	balance, err := c.ethClient.BalanceAt(context.Background(), account, nil)
	if err != nil {
		logger.Error().Err(err).Str("address", address).Msg("Failed to fetch balance")
		return nil, err
	}

	logger.Info().Str("address", address).Msg("Fetched balance successfully")
	return balance, nil
}

// GetBalanceInEth returns the balance of an Ethereum address in ETH
func (c *Client) GetBalanceInEth(address string) (*big.Int, *big.Float, error) {
	// Get balance in wei
	balance, err := c.GetBalance(address)
	if err != nil {
		logger.Error().Err(err).Str("address", address).Msg("Failed to fetch balance in ETH")
		return nil, nil, err
	}

	// Convert wei to ETH
	weiBalance := new(big.Float).SetInt(balance)
	ethBalance := new(big.Float).Quo(weiBalance, big.NewFloat(1e18))

	logger.Info().Str("address", address).Msg("Fetched balance in ETH successfully")
	return balance, ethBalance, nil
}

// CreateBalanceRecord creates a balance record from address and balance
func (c *Client) CreateBalanceRecord(address string) (models.BalanceRecord, error) {
	// Validate address
	if !common.IsHexAddress(address) {
		return models.BalanceRecord{}, ErrInvalidAddress
	}

	// Get balances
	balance, ethBalance, err := c.GetBalanceInEth(address)
	if err != nil {
		logger.Error().Err(err).Str("address", address).Msg("Failed to create balance record")
		return models.BalanceRecord{}, err
	}

	// Create record
	balanceRecord := models.BalanceRecord{
		Address:    address,
		Balance:    balance.String(),
		BalanceETH: ethBalance.Text('f', 18),
		FetchedAt:  time.Now(),
	}

	logger.Info().Str("address", address).Msg("Created balance record successfully")
	return balanceRecord, nil
}
