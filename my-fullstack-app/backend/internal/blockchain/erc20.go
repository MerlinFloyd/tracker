package blockchain

import (
	"context"
	"math/big"
	"strings"
	"time"

	"my-fullstack-app/backend/internal/models"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// Common ERC20 contract ABI for balanceOf method
const erc20ABIJson = `[
    {
        "constant": true,
        "inputs": [{"name": "_owner", "type": "address"}],
        "name": "balanceOf",
        "outputs": [{"name": "balance", "type": "uint256"}],
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "decimals",
        "outputs": [{"name": "", "type": "uint8"}],
        "type": "function"
    },
    {
        "constant": true,
        "inputs": [],
        "name": "symbol",
        "outputs": [{"name": "", "type": "string"}],
        "type": "function"
    }
]`

// Common errors
var ()

// TokenInfo represents information about an ERC20 token
type TokenInfo struct {
	Address   string
	Symbol    string
	Decimals  uint8
	TokenName string
}

// TokenBalance represents a balance of an ERC20 token
type TokenBalance struct {
	Token      TokenInfo
	WeiBalance *big.Int
	Balance    *big.Float
}

// Common tokens
var CommonTokens = map[string]TokenInfo{
	"USDT": {
		Address:  "0xdAC17F958D2ee523a2206206994597C13D831ec7",
		Symbol:   "USDT",
		Decimals: 6,
	},
	"USDC": {
		Address:  "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		Symbol:   "USDC",
		Decimals: 6,
	},
	"DAI": {
		Address:  "0x6B175474E89094C44Da98b954EedeAC495271d0F",
		Symbol:   "DAI",
		Decimals: 18,
	},
	"LINK": {
		Address:  "0x514910771AF9Ca656af840dff83E8264EcF986CA",
		Symbol:   "LINK",
		Decimals: 18,
	},
}

// ERC20 represents an ERC20 token contract
type ERC20 struct {
	client    *Client
	contract  *bind.BoundContract
	tokenInfo TokenInfo
	address   common.Address
}

// NewERC20 creates a new ERC20 token client
func (c *Client) NewERC20(tokenAddress string) (*ERC20, error) {
	// Validate token address
	if !common.IsHexAddress(tokenAddress) {
		return nil, ErrInvalidTokenAddress
	}

	address := common.HexToAddress(tokenAddress)

	// Parse ABI
	parsedAbi, err := abi.JSON(strings.NewReader(erc20ABIJson))
	if err != nil {
		return nil, err
	}

	// Create contract binding
	contract := bind.NewBoundContract(address, parsedAbi, c.ethClient, c.ethClient, c.ethClient)

	// Get token information
	var symbol string
	var decimals uint8

	// Try to get symbol
	var result []interface{}
	err = contract.Call(nil, &result, "symbol")
	if err == nil && len(result) > 0 {
		if str, ok := result[0].(string); ok {
			symbol = str
		}
	} else {
		symbol = "UNKNOWN"
	}
	if err != nil {
		symbol = "UNKNOWN"
	}

	// Try to get decimals
	result = nil // Clear the result slice before reuse
	err = contract.Call(nil, &result, "decimals")
	if err == nil && len(result) > 0 {
		if val, ok := result[0].(uint8); ok {
			decimals = val
		} else {
			decimals = 18 // Default to 18 if the type assertion fails
		}
	} else {
		decimals = 18 // Default to 18 if we can't get the value
	}

	tokenInfo := TokenInfo{
		Address:  tokenAddress,
		Symbol:   symbol,
		Decimals: decimals,
	}

	return &ERC20{
		client:    c,
		contract:  contract,
		tokenInfo: tokenInfo,
		address:   address,
	}, nil
}

// GetERC20Balance returns the balance of an ERC20 token for an address
func (e *ERC20) GetBalance(address string) (*big.Int, error) {
	// Validate address
	if !common.IsHexAddress(address) {
		return nil, ErrInvalidAddress
	}

	// Call balanceOf function
	ownerAddress := common.HexToAddress(address)
	var balance big.Int

	callOpts := &bind.CallOpts{Context: context.Background()}
	var result []interface{}
	err := e.contract.Call(callOpts, &result, "balanceOf", ownerAddress)
	if err != nil || len(result) == 0 {
		return nil, ErrTokenContract
	}

	if balanceValue, ok := result[0].(*big.Int); ok {
		balance.Set(balanceValue)
	} else {
		return nil, ErrTokenContract
	}

	return &balance, nil
}

// GetFormattedBalance returns the balance in token units (considering decimals)
func (e *ERC20) GetFormattedBalance(address string) (*big.Int, *big.Float, error) {
	// Get raw balance
	balance, err := e.GetBalance(address)
	if err != nil {
		return nil, nil, err
	}

	// Convert to token units based on decimals
	divisor := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10), big.NewInt(int64(e.tokenInfo.Decimals)), nil,
	))

	tokenBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), divisor)

	return balance, tokenBalance, nil
}

// CreateTokenBalanceRecord creates a token balance record
func (e *ERC20) CreateTokenBalanceRecord(address string) (models.TokenBalanceRecord, error) {
	// Validate address
	if !common.IsHexAddress(address) {
		return models.TokenBalanceRecord{}, ErrInvalidAddress
	}

	// Get balances
	rawBalance, formattedBalance, err := e.GetFormattedBalance(address)
	if err != nil {
		return models.TokenBalanceRecord{}, err
	}

	// Create record
	balanceRecord := models.TokenBalanceRecord{
		Address:    address,
		Balance:    rawBalance.String(),
		BalanceETH: formattedBalance.Text('f', int(e.tokenInfo.Decimals)),
		FetchedAt:  time.Now(),
	}

	return balanceRecord, nil
}

// GetMultipleTokenBalances fetches balances for multiple tokens
func (c *Client) GetMultipleTokenBalances(address string, tokenAddresses []string) ([]models.TokenBalanceRecord, error) {
	// Validate address
	if !common.IsHexAddress(address) {
		return nil, ErrInvalidAddress
	}

	var results []models.TokenBalanceRecord

	// Get balance for each token
	for _, tokenAddress := range tokenAddresses {
		// Create ERC20 client
		erc20, err := c.NewERC20(tokenAddress)
		if err != nil {
			continue // Skip tokens with errors
		}

		// Get token balance record
		balanceRecord, err := erc20.CreateTokenBalanceRecord(address)
		if err != nil {
			continue // Skip tokens with errors
		}

		results = append(results, balanceRecord)
	}

	return results, nil
}

// GetCommonTokenBalances fetches balances for common tokens
func (c *Client) GetCommonTokenBalances(address string) ([]models.TokenBalanceRecord, error) {
	// Extract token addresses
	var tokenAddresses []string
	for _, token := range CommonTokens {
		tokenAddresses = append(tokenAddresses, token.Address)
	}

	return c.GetMultipleTokenBalances(address, tokenAddresses)
}
