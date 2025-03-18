package api

// BalanceResponse is a specific response for balance queries
// @Description Balance response format
type BalanceResponse struct {
	Address   string `json:"address" example:"0x742d35Cc6634C0532925a3b844Bc454e4438f44e"`
	Wei       string `json:"wei" example:"1000000000000000000"`
	Eth       string `json:"eth" example:"1.0"`
	Timestamp string `json:"timestamp,omitempty" example:"2025-03-18T12:00:00Z"`
}
