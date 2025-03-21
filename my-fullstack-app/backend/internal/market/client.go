package market

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"my-fullstack-app/backend/internal/logger"

	"github.com/adshao/go-binance/v2"
)

const (
	// Default values
	defaultInterval = "1d" // Daily klines/candlesticks
	dateFormat      = "2006-01-02"
)

// Client represents a market data client
type Client struct {
	binanceClient *binance.Client
}

// PriceData represents price information at a specific time
type PriceData struct {
	Symbol       string    `json:"symbol"`
	Price        float64   `json:"price"`
	Timestamp    time.Time `json:"timestamp"`
	USD          float64   `json:"usd,omitempty"` // Price in USD
	OpenTime     time.Time `json:"openTime,omitempty"`
	CloseTime    time.Time `json:"closeTime,omitempty"`
	High         float64   `json:"high,omitempty"`
	Low          float64   `json:"low,omitempty"`
	Volume       float64   `json:"volume,omitempty"`
	NumberTrades int64     `json:"numberTrades,omitempty"`
}

// NewClient creates a new client with the Binance API
func NewClient(apiKey, secretKey string) *Client {
	// Create Binance client
	binanceClient := binance.NewClient(apiKey, secretKey)

	logger.Info().Msg("Market price client initialized")

	return &Client{
		binanceClient: binanceClient,
	}
}

// GetCurrentPrice gets the latest price for a symbol
func (c *Client) GetCurrentPrice(ctx context.Context, symbol string) (*PriceData, error) {
	logger.Debug().
		Str("symbol", symbol).
		Msg("Getting current price")

	// Get ticker price from Binance
	prices, err := c.binanceClient.NewListPricesService().Symbol(symbol).Do(ctx)
	if err != nil {
		logger.Error().
			Err(err).
			Str("symbol", symbol).
			Msg("Failed to get current price")
		return nil, fmt.Errorf("failed to get ticker price: %w", err)
	}

	if len(prices) == 0 {
		return nil, fmt.Errorf("no price data found for symbol %s", symbol)
	}

	price, err := strconv.ParseFloat(prices[0].Price, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price: %w", err)
	}

	// Get more detailed ticker for additional info
	ticker24h, err := c.binanceClient.NewTradingDayTickerService().Symbol(symbol).Do(ctx)
	if err != nil {
		// Not critical, just log and continue
		logger.Warn().
			Err(err).
			Str("symbol", symbol).
			Msg("Failed to get 24h ticker data")
	}

	priceData := &PriceData{
		Symbol:    symbol,
		Price:     price,
		Timestamp: time.Now(),
		USD:       price, // For USDT pairs, this is already in USD equivalent
	}

	// Add additional data if available
	if ticker24h != nil {
		volume, _ := strconv.ParseFloat(ticker24h[0].Volume, 64)
		high, _ := strconv.ParseFloat(ticker24h[0].HighPrice, 64)
		low, _ := strconv.ParseFloat(ticker24h[0].LowPrice, 64)

		priceData.Volume = volume
		priceData.High = high
		priceData.Low = low
		priceData.NumberTrades = int64(ticker24h[0].Count)
	}

	logger.Info().
		Str("symbol", symbol).
		Float64("price", price).
		Msg("Successfully retrieved current price")

	return priceData, nil
}

// GetHistoricalPrice gets the price at a specific date
func (c *Client) GetHistoricalPrice(ctx context.Context, symbol string, date time.Time) (*PriceData, error) {
	logger.Debug().
		Str("symbol", symbol).
		Time("date", date).
		Msg("Getting historical price")

	// Format date to start of day in UTC
	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	// End time is the end of the day
	endTime := startTime.Add(24 * time.Hour)

	// Fetch klines (candlestick) data
	klines, err := c.binanceClient.NewKlinesService().
		Symbol(symbol).
		Interval(defaultInterval).
		StartTime(startTime.UnixNano() / int64(time.Millisecond)).
		EndTime(endTime.UnixNano() / int64(time.Millisecond)).
		Limit(1).
		Do(ctx)

	if err != nil {
		logger.Error().
			Err(err).
			Str("symbol", symbol).
			Time("date", date).
			Msg("Failed to get historical klines")
		return nil, fmt.Errorf("failed to get klines: %w", err)
	}

	if len(klines) == 0 {
		logger.Warn().
			Str("symbol", symbol).
			Time("date", date).
			Msg("No historical data found")
		return nil, fmt.Errorf("no historical data found for %s on %s", symbol, date.Format(dateFormat))
	}

	// Parse the closing price
	closePrice, err := strconv.ParseFloat(klines[0].Close, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse close price: %w", err)
	}

	// Parse other kline data
	//openPrice, _ := strconv.ParseFloat(klines[0].Open, 64)
	highPrice, _ := strconv.ParseFloat(klines[0].High, 64)
	lowPrice, _ := strconv.ParseFloat(klines[0].Low, 64)
	volume, _ := strconv.ParseFloat(klines[0].Volume, 64)

	// Convert timestamps
	openTime := time.Unix(klines[0].OpenTime/1000, 0)
	closeTime := time.Unix(klines[0].CloseTime/1000, 0)

	priceData := &PriceData{
		Symbol:       symbol,
		Price:        closePrice,
		Timestamp:    closeTime,
		OpenTime:     openTime,
		CloseTime:    closeTime,
		High:         highPrice,
		Low:          lowPrice,
		Volume:       volume,
		NumberTrades: klines[0].TradeNum,
		USD:          closePrice, // For USDT pairs, this is already in USD
	}

	logger.Info().
		Str("symbol", symbol).
		Time("date", date).
		Float64("price", closePrice).
		Msg("Successfully retrieved historical price")

	return priceData, nil
}

// ConvertToUSD converts a crypto price to USD equivalent
func (c *Client) ConvertToUSD(ctx context.Context, symbol string, price float64, date *time.Time) (float64, error) {
	// If the symbol already ends with USDT, it's already in USD equivalent
	if len(symbol) > 4 && symbol[len(symbol)-4:] == "USDT" {
		return price, nil
	}

	// If the symbol already ends with BUSD or USDC, it's already in USD equivalent
	if len(symbol) > 4 && (symbol[len(symbol)-4:] == "BUSD" || symbol[len(symbol)-4:] == "USDC") {
		return price, nil
	}

	// Extract the quote asset from the symbol
	// This is a simplification assuming common formats like ETHBTC, BTCETH, etc.
	if len(symbol) < 6 {
		return 0, fmt.Errorf("symbol format not recognized for USD conversion: %s", symbol)
	}

	var baseAsset, quoteAsset string

	// Try to identify common quote assets
	commonQuotes := []string{"BTC", "ETH", "BNB"}
	for _, quote := range commonQuotes {
		if len(symbol) >= len(quote) && symbol[len(symbol)-len(quote):] == quote {
			baseAsset = symbol[:len(symbol)-len(quote)]
			quoteAsset = quote
			break
		}
	}

	if quoteAsset == "" {
		return 0, fmt.Errorf("couldn't determine quote asset for symbol %s", symbol)
	}

	logger.Debug().
		Str("symbol", symbol).
		Str("base", baseAsset).
		Str("quote", quoteAsset).
		Msg("Converting price to USD")

	// Get the quote asset's price in USDT
	var quotePrice *PriceData
	var err error

	quoteSymbol := quoteAsset + "USDT"

	if date == nil {
		// Get current price
		quotePrice, err = c.GetCurrentPrice(ctx, quoteSymbol)
	} else {
		// Get historical price
		quotePrice, err = c.GetHistoricalPrice(ctx, quoteSymbol, *date)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to get %s price: %w", quoteSymbol, err)
	}

	// Calculate USD equivalent
	usdPrice := price * quotePrice.Price

	logger.Info().
		Str("symbol", symbol).
		Float64("original_price", price).
		Str("quote_asset", quoteAsset).
		Float64("quote_price_usd", quotePrice.Price).
		Float64("converted_usd", usdPrice).
		Msg("Successfully converted price to USD")

	return usdPrice, nil
}
