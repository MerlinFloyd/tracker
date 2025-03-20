package models

import (
	"time"
)

// TokenBalanceRecord represents a stored ERC20 token balance
type TokenBalanceRecord struct {
	ID         int       `json:"id" db:"id"`
	Address    string    `json:"address" db:"address"`
	Balance    string    `json:"balance" db:"balance"`         // Raw balance as string
	BalanceETH string    `json:"balance_eth" db:"balance_eth"` // Formatted with decimals
	FetchedAt  time.Time `json:"fetched_at" db:"fetched_at"`
}
