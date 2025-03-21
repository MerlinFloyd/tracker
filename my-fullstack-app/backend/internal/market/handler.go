package market

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"my-fullstack-app/backend/internal/api"
	"my-fullstack-app/backend/internal/logger"
)

// Handler handles market data API requests
type Handler struct {
	client *Client
}

// NewHandler creates a new market data handler
func NewHandler() (*Handler, error) {
	// Get API keys from environment variables
	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")

	// Create client
	client := NewClient(apiKey, secretKey)

	logger.Info().Msg("Market data handler initialized")

	return &Handler{
		client: client,
	}, nil
}

// GetCurrentPriceHandler returns the current price of a symbol
// @Summary Get current price
// @Description Returns the current price of a cryptocurrency
// @Tags market
// @Accept json
// @Produce json
// @Param symbol query string true "Trading pair symbol (e.g., BTCUSDT)"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /market/price [get]
func (h *Handler) GetCurrentPriceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get symbol from query parameters
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		logger.Warn().Msg("Missing symbol parameter")
		api.RespondWithError(w, http.StatusBadRequest, "Symbol parameter is required")
		return
	}

	// Create context with timeout
	ctx := r.Context()

	logger.Info().
		Str("symbol", symbol).
		Str("remote_addr", r.RemoteAddr).
		Msg("Current price request received")

	// Get current price
	priceData, err := h.client.GetCurrentPrice(ctx, symbol)
	if err != nil {
		logger.Error().
			Err(err).
			Str("symbol", symbol).
			Msg("Failed to get current price")
		api.RespondWithError(w, http.StatusInternalServerError, "Failed to get current price")
		return
	}

	// Check if we need to convert to USD
	convertToUsd := r.URL.Query().Get("convert_usd") == "true"

	if convertToUsd && !symbolEndsWithUsdStablecoin(symbol) {
		usdPrice, err := h.client.ConvertToUSD(ctx, symbol, priceData.Price, nil)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("symbol", symbol).
				Msg("Failed to convert price to USD")
		} else {
			priceData.USD = usdPrice
		}
	}

	response := api.Response{
		Success: true,
		Message: "Current price retrieved successfully",
		Data:    priceData,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error().Err(err).Msg("Failed to encode response")
	}
}

// GetHistoricalPriceHandler returns the price of a symbol at a specific date
// @Summary Get historical price
// @Description Returns the price of a cryptocurrency at a specific date
// @Tags market
// @Accept json
// @Produce json
// @Param symbol query string true "Trading pair symbol (e.g., BTCUSDT)"
// @Param date query string true "Date in YYYY-MM-DD format"
// @Param convert_usd query boolean false "Convert price to USD equivalent"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /market/historical [get]
func (h *Handler) GetHistoricalPriceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get symbol from query parameters
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		logger.Warn().Msg("Missing symbol parameter")
		api.RespondWithError(w, http.StatusBadRequest, "Symbol parameter is required")
		return
	}

	// Get date from query parameters
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		logger.Warn().Msg("Missing date parameter")
		api.RespondWithError(w, http.StatusBadRequest, "Date parameter is required (format: YYYY-MM-DD)")
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		logger.Warn().
			Err(err).
			Str("date", dateStr).
			Msg("Invalid date format")
		api.RespondWithError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
		return
	}

	// Create context with timeout
	ctx := r.Context()

	logger.Info().
		Str("symbol", symbol).
		Str("date", dateStr).
		Str("remote_addr", r.RemoteAddr).
		Msg("Historical price request received")

	// Get historical price
	priceData, err := h.client.GetHistoricalPrice(ctx, symbol, date)
	if err != nil {
		logger.Error().
			Err(err).
			Str("symbol", symbol).
			Str("date", dateStr).
			Msg("Failed to get historical price")
		api.RespondWithError(w, http.StatusInternalServerError, "Failed to get historical price")
		return
	}

	// Check if we need to convert to USD
	convertToUsd := r.URL.Query().Get("convert_usd") == "true"

	if convertToUsd && !symbolEndsWithUsdStablecoin(symbol) {
		usdPrice, err := h.client.ConvertToUSD(ctx, symbol, priceData.Price, &date)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("symbol", symbol).
				Str("date", dateStr).
				Msg("Failed to convert price to USD")
		} else {
			priceData.USD = usdPrice
		}
	}

	response := api.Response{
		Success: true,
		Message: "Historical price retrieved successfully",
		Data:    priceData,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error().Err(err).Msg("Failed to encode response")
	}
}

// Helper function to check if a symbol ends with a USD stablecoin symbol
func symbolEndsWithUsdStablecoin(symbol string) bool {
	stablecoins := []string{"USDT", "USDC", "BUSD", "DAI", "TUSD", "USDP"}
	for _, coin := range stablecoins {
		if len(symbol) >= len(coin) && symbol[len(symbol)-len(coin):] == coin {
			return true
		}
	}
	return false
}
