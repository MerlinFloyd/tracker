package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
}

// BalanceRecord represents a stored Ethereum balance
type BalanceRecord struct {
	ID         int       `json:"id" db:"id"`
	Address    string    `json:"address" db:"address"`
	Balance    string    `json:"balance" db:"balance"`         // wei, stored as string due to large size
	BalanceETH string    `json:"balance_eth" db:"balance_eth"` // ETH value as string
	FetchedAt  time.Time `json:"fetched_at" db:"fetched_at"`
}
