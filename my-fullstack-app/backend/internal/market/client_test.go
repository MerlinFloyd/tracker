package market

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestGetCurrentPrice(t *testing.T) {
	// Skip if no API keys available
	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")
	if apiKey == "" || secretKey == "" {
		t.Skip("Skipping test: No Binance API keys available")
	}

	client := NewClient(apiKey, secretKey)
	ctx := context.Background()

	testCases := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:    "Valid BTC/USDT",
			symbol:  "BTCUSDT",
			wantErr: false,
		},
		{
			name:    "Valid ETH/USDT",
			symbol:  "ETHUSDT",
			wantErr: false,
		},
		{
			name:    "Valid BTC/ETH",
			symbol:  "BTCETH",
			wantErr: false,
		},
		{
			name:    "Invalid symbol",
			symbol:  "INVALIDCOIN",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			price, err := client.GetCurrentPrice(ctx, tc.symbol)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error for symbol %s, got nil", tc.symbol)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if price == nil {
				t.Errorf("Expected price data, got nil")
				return
			}

			if price.Symbol != tc.symbol {
				t.Errorf("Expected symbol %s, got %s", tc.symbol, price.Symbol)
			}

			if price.Price <= 0 {
				t.Errorf("Expected positive price, got %f", price.Price)
			}
		})
	}
}

func TestGetHistoricalPrice(t *testing.T) {
	// Skip if no API keys available
	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")
	if apiKey == "" || secretKey == "" {
		t.Skip("Skipping test: No Binance API keys available")
	}

	client := NewClient(apiKey, secretKey)
	ctx := context.Background()

	// Test with a known date (one year ago)
	date := time.Now().AddDate(-1, 0, 0)

	testCases := []struct {
		name    string
		symbol  string
		date    time.Time
		wantErr bool
	}{
		{
			name:    "BTC/USDT historical",
			symbol:  "BTCUSDT",
			date:    date,
			wantErr: false,
		},
		{
			name:    "ETH/USDT historical",
			symbol:  "ETHUSDT",
			date:    date,
			wantErr: false,
		},
		{
			name:    "Future date",
			symbol:  "BTCUSDT",
			date:    time.Now().AddDate(1, 0, 0), // 1 year in the future
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			price, err := client.GetHistoricalPrice(ctx, tc.symbol, tc.date)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error for symbol %s and date %v, got nil",
						tc.symbol, tc.date)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if price == nil {
				t.Errorf("Expected price data, got nil")
				return
			}

			if price.Symbol != tc.symbol {
				t.Errorf("Expected symbol %s, got %s", tc.symbol, price.Symbol)
			}

			if price.Price <= 0 {
				t.Errorf("Expected positive price, got %f", price.Price)
			}
		})
	}
}
