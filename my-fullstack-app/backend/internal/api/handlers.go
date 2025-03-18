package api

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"os"
	"time"

	"my-fullstack-app/backend/internal/database"
	"my-fullstack-app/backend/internal/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Response is a struct for standard API responses
// @Description API response format
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Global ethclient
var ethClient *ethclient.Client

// InitEthClient initializes the Ethereum client connection
func InitEthClient() error {
	// Get Infura API key from environment variable
	infuraKey := os.Getenv("INFURA_API_KEY")
	if infuraKey == "" {
		infuraKey = "4b4eabeb1b8b4bfeaa4f29e754f2d282" // Replace with your actual Infura API key for testing
	}

	// Connect to Infura
	infuraURL := "https://mainnet.infura.io/v3/" + infuraKey
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		return err
	}

	ethClient = client
	return nil
}

// HealthCheckHandler handles the /health endpoint
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := "OK"
	healthDetails := map[string]interface{}{
		"api": "running",
	}

	// Check Ethereum client connection
	if ethClient == nil {
		if err := InitEthClient(); err != nil {
			healthDetails["ethereum"] = "disconnected"
			status = "degraded"
		} else {
			// Try getting a block to verify connection is working
			_, err := ethClient.BlockNumber(context.Background())
			if err != nil {
				healthDetails["ethereum"] = "error: " + err.Error()
				status = "degraded"
			} else {
				healthDetails["ethereum"] = "connected"
			}
		}
	} else {
		healthDetails["ethereum"] = "connected"
	}

	// Check database connection
	db, err := database.Connect()
	if err != nil {
		healthDetails["database"] = "disconnected"
		status = "degraded"
	} else {
		// Test the database with a simple ping
		err = db.Ping()
		if err != nil {
			healthDetails["database"] = "error: " + err.Error()
			status = "degraded"
		} else {
			healthDetails["database"] = "connected"
		}
		db.Close()
	}

	response := Response{
		Message: status,
		Data:    healthDetails,
	}

	// If any service is degraded, return a 503 status code
	if status != "OK" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(response)
}

// BlockNumberHandler returns the latest Ethereum block number
func BlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if ethClient == nil {
		if err := InitEthClient(); err != nil {
			http.Error(w, "Failed to initialize Ethereum client", http.StatusInternalServerError)
			return
		}
	}

	// Get the latest block number
	blockNumber, err := ethClient.BlockNumber(context.Background())
	if err != nil {
		http.Error(w, "Failed to get block number", http.StatusInternalServerError)
		return
	}

	response := Response{
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
// @Success      200      {object}  Response
// @Failure      400      {object}  Response
// @Failure      500      {object}  Response
// @Router       /eth/balance [get]
func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	if ethClient == nil {
		if err := InitEthClient(); err != nil {
			http.Error(w, "Failed to initialize Ethereum client", http.StatusInternalServerError)
			return
		}
	}

	// Get the balance
	account := common.HexToAddress(address)
	balance, err := ethClient.BalanceAt(context.Background(), account, nil)
	if err != nil {
		http.Error(w, "Failed to get account balance", http.StatusInternalServerError)
		return
	}

	// Convert wei to ETH
	// 1 ETH = 10^18 wei
	weiBalance := new(big.Float).SetInt(balance)
	ethBalance := new(big.Float).Quo(weiBalance, big.NewFloat(1e18))

	response := Response{
		Message: "Account balance retrieved",
		Data: map[string]interface{}{
			"wei": balance.String(),
			"eth": ethBalance.Text('f', 18), // Format with up to 18 decimal places
		},
	}
	json.NewEncoder(w).Encode(response)
}

// StoreBalanceHandler retrieves the balance for an Ethereum address and stores it in the database
// @Summary      Store Ethereum address balance
// @Description  Retrieves and stores the balance of an Ethereum address
// @Tags         ethereum
// @Accept       json
// @Produce      json
// @Param        address  query     string  true  "Ethereum address (0x format)"
// @Success      200      {object}  Response
// @Failure      400      {object}  Response
// @Failure      500      {object}  Response
// @Router       /eth/store-balance [get]
func StoreBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get address from query parameters
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	// Validate Ethereum address format
	if !common.IsHexAddress(address) {
		http.Error(w, "Invalid Ethereum address format", http.StatusBadRequest)
		return
	}

	if ethClient == nil {
		if err := InitEthClient(); err != nil {
			http.Error(w, "Failed to initialize Ethereum client", http.StatusInternalServerError)
			return
		}
	}

	// Get the balance
	account := common.HexToAddress(address)
	balance, err := ethClient.BalanceAt(context.Background(), account, nil)
	if err != nil {
		http.Error(w, "Failed to get account balance", http.StatusInternalServerError)
		return
	}

	// Convert wei to ETH
	weiBalance := new(big.Float).SetInt(balance)
	ethBalance := new(big.Float).Quo(weiBalance, big.NewFloat(1e18))
	ethBalanceStr := ethBalance.Text('f', 18)

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Create a balance record
	balanceRecord := models.BalanceRecord{
		Address:    address,
		Balance:    balance.String(),
		BalanceETH: ethBalanceStr,
		FetchedAt:  time.Now(),
	}

	// Store the balance in the database
	balanceID, err := database.StoreBalance(db, balanceRecord)
	if err != nil {
		http.Error(w, "Failed to store balance in database", http.StatusInternalServerError)
		return
	}

	response := Response{
		Message: "Account balance retrieved and stored",
		Data: map[string]interface{}{
			"id":        balanceID,
			"address":   address,
			"wei":       balance.String(),
			"eth":       ethBalanceStr,
			"timestamp": balanceRecord.FetchedAt,
		},
	}
	json.NewEncoder(w).Encode(response)
}
