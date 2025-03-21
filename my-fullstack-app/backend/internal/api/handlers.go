package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"my-fullstack-app/backend/internal/database"

	"github.com/ethereum/go-ethereum/ethclient"
)

// Response is a struct for standard API responses
// @Description API response format
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
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
// @Summary      Check API health status
// @Description  Returns health status of the API, including Ethereum client and database connectivity
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  Response  "API is healthy"
// @Failure      503  {object}  Response  "One or more components are in degraded state"
// @Router       /health [get]
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
