package blockchain

import (
	"encoding/json"
	"net/http"

	"my-fullstack-app/backend/internal/api"
	"my-fullstack-app/backend/internal/database"

	"github.com/ethereum/go-ethereum/common"
)

// Handler handles blockchain-related HTTP requests
type Handler struct {
	client *Client
}

// NewHandler creates a new blockchain handler
func NewHandler() (*Handler, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	return &Handler{
		client: client,
	}, nil
}

// BlockNumberHandler returns the latest Ethereum block number
// @Summary      Get latest Ethereum block number
// @Description  Returns the latest block number from the Ethereum blockchain
// @Tags         ethereum
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.Response
// @Failure      500  {object}  api.Response
// @Router       /eth/block [get]
func (h *Handler) BlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the latest block number
	blockNumber, err := h.client.GetBlockNumber()
	if err != nil {
		http.Error(w, "Failed to get block number", http.StatusInternalServerError)
		return
	}

	response := api.Response{
		Message: "Current Ethereum block number",
		Data:    blockNumber,
	}
	json.NewEncoder(w).Encode(response)
}

// GetBalanceHandler returns the balance for a given Ethereum address
// @Summary      Get Ethereum address balance
// @Description  Returns the balance of an Ethereum address in wei and ETH
// @Tags         ethereum
// @Accept       json
// @Produce      json
// @Param        address  query     string  true  "Ethereum address (0x format)"
// @Failure      400      {object}  api.Response
// @Failure      500      {object}  api.Response
// @Router       /eth/balance [get]
func (h *Handler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	// Validate address
	if !common.IsHexAddress(address) {
		http.Error(w, "Invalid Ethereum address format", http.StatusBadRequest)
		return
	}

	// Get balance
	balance, ethBalance, err := h.client.GetBalanceInEth(address)
	if err != nil {
		http.Error(w, "Failed to get account balance", http.StatusInternalServerError)
		return
	}

	response := api.Response{
		Message: "Account balance retrieved",
		Data: map[string]interface{}{
			"wei": balance.String(),
			"eth": ethBalance.Text('f', 18),
		},
	}
	json.NewEncoder(w).Encode(response)
}

// StoreBalanceHandler retrieves the balance for an Ethereum address and stores it in the database
// @Summary      Store Ethereum address balance
// @Description  Retrieves and stores the balance of an Ethereum address
// @Tags         tokens
// @Accept       json
// @Produce      json
// @Param        address  query     string  true  "Ethereum address (0x format)"
// @Success      200      {object}  api.Response
// @Failure      400      {object}  api.Response
// @Failure      500      {object}  api.Response
// @Router       /eth/store-balance [get]
func (h *Handler) StoreBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get address from query parameters
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	// Create balance record
	balanceRecord, err := h.client.CreateBalanceRecord(address)
	if err == ErrInvalidAddress {
		http.Error(w, "Invalid Ethereum address format", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Failed to get account balance", http.StatusInternalServerError)
		return
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Store the balance in the database
	balanceID, err := database.StoreBalance(db, balanceRecord)
	if err != nil {
		http.Error(w, "Failed to store balance in database", http.StatusInternalServerError)
		return
	}

	response := api.Response{
		Message: "Account balance retrieved and stored",
		Data: map[string]interface{}{
			"id":        balanceID,
			"address":   balanceRecord.Address,
			"wei":       balanceRecord.Balance,
			"eth":       balanceRecord.BalanceETH,
			"timestamp": balanceRecord.FetchedAt,
		},
	}
	json.NewEncoder(w).Encode(response)
}

// GetTokenBalancesHandler returns all token balance records for a specific token address
// @Summary      Get token balance records
// @Description  Returns all balance records for a specific ERC20 token
// @Tags         tokens
// @Accept       json
// @Produce      json
// @Param        token_address  query  string  true  "ERC20 token address (0x format)"
// @Success      200  {object}  api.Response
// @Failure      400  {object}  api.Response
// @Failure      500  {object}  api.Response
// @Router       /eth/get-token-balances [get]
func (h *Handler) GetTokenBalancesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get token address from query parameters
	tokenAddress := r.URL.Query().Get("token_address")
	if tokenAddress == "" {
		http.Error(w, "token_address parameter is required", http.StatusBadRequest)
		return
	}

	// Validate token address
	if !common.IsHexAddress(tokenAddress) {
		http.Error(w, "Invalid token address format", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Retrieve token balances
	balances, err := database.GetTokenBalances(db, tokenAddress)
	if err != nil {
		http.Error(w, "Failed to retrieve token balances", http.StatusInternalServerError)
		return
	}

	// Check if any records were found
	if len(balances) == 0 {
		response := api.Response{
			Message: "No token balance records found for this token",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := api.Response{
		Message: "Token balance records retrieved",
		Data:    balances,
	}
	json.NewEncoder(w).Encode(response)
}
