package blockchain

import "errors"

var (
	// ErrClientNotInitialized is returned when the Ethereum client is not initialized
	ErrClientNotInitialized = errors.New("ethereum client not initialized")
	// ErrInvalidAddress is returned when an invalid Ethereum address is provided
	ErrInvalidAddress      = errors.New("invalid ethereum address format")
	ErrInvalidTokenAddress = errors.New("invalid token contract address")
	ErrTokenContract       = errors.New("error interacting with token contract")
)
