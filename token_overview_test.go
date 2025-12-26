package birdeye

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
)

func TestGetTokenOverview_Success(t *testing.T) {
	responses := map[string]interface{}{
		"/defi/token_overview": wrapResponse(map[string]interface{}{
			"address":                     "TokenMint123",
			"symbol":                      "TEST",
			"name":                        "Test Token",
			"decimals":                    9,
			"logoURI":                     "https://example.com/logo.png",
			"liquidity":                   50000.50,
			"price":                       1.234567,
			"priceChange24hPercent":       5.25,
			"v24h":                        1000000.0,
			"v24hUSD":                     1234567.89,
			"v24hChangePercent":           12.5,
			"mc":                          10000000.0,
			"supply":                      1000000000.0,
			"circulatingSupply":           800000000.0,
			"holder":                      5000,
			"trade24h":                    1500,
			"trade24hChangePercent":       8.3,
			"buy24h":                      800,
			"sell24h":                     700,
			"uniqueWallet24h":             350,
			"uniqueWallet24hChangePercent": 15.2,
			"lastTradeUnixTime":           1703980800,
			"lastTradeHumanTime":          "2024-12-31T00:00:00Z",
		}),
	}

	server := testServer(t, responses)
	defer server.Close()

	client := testClient(t, server.URL)
	overview, err := client.GetTokenOverview(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if overview.Symbol != "TEST" {
		t.Errorf("expected symbol 'TEST', got '%s'", overview.Symbol)
	}

	if overview.Name != "Test Token" {
		t.Errorf("expected name 'Test Token', got '%s'", overview.Name)
	}

	if overview.Decimals != 9 {
		t.Errorf("expected decimals 9, got %d", overview.Decimals)
	}

	if overview.Holder != 5000 {
		t.Errorf("expected holder 5000, got %d", overview.Holder)
	}

	if overview.Trade24h != 1500 {
		t.Errorf("expected trade24h 1500, got %d", overview.Trade24h)
	}

	if overview.Buy24h != 800 {
		t.Errorf("expected buy24h 800, got %d", overview.Buy24h)
	}

	if overview.Sell24h != 700 {
		t.Errorf("expected sell24h 700, got %d", overview.Sell24h)
	}

	if overview.LastTradeUnixTime != 1703980800 {
		t.Errorf("expected lastTradeUnixTime 1703980800, got %d", overview.LastTradeUnixTime)
	}
}

func TestGetTokenOverview_DecimalPrecision(t *testing.T) {
	responses := map[string]interface{}{
		"/defi/token_overview": wrapResponse(map[string]interface{}{
			"address":    "TokenMint123",
			"symbol":     "PREC",
			"name":       "Precision Test",
			"decimals":   9,
			"liquidity":  123456.789012,
			"price":      0.00000123,
			"mc":         9999999.99,
			"supply":     1000000000000.123456,
			"v24hUSD":    0.01,
		}),
	}

	server := testServer(t, responses)
	defer server.Close()

	client := testClient(t, server.URL)
	overview, err := client.GetTokenOverview(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify decimal.Decimal correctly preserves precision
	expectedLiquidity := decimal.NewFromFloat(123456.789012)
	if !overview.Liquidity.Equal(expectedLiquidity) {
		t.Errorf("liquidity precision lost: expected %s, got %s", expectedLiquidity, overview.Liquidity)
	}

	// Very small price values
	expectedPrice := decimal.NewFromFloat(0.00000123)
	if !overview.Price.Equal(expectedPrice) {
		t.Errorf("price precision lost: expected %s, got %s", expectedPrice, overview.Price)
	}

	// Large market cap values
	expectedMC := decimal.NewFromFloat(9999999.99)
	if !overview.MarketCap.Equal(expectedMC) {
		t.Errorf("marketCap precision lost: expected %s, got %s", expectedMC, overview.MarketCap)
	}
}

func TestGetTokenOverview_WithExtensions(t *testing.T) {
	responses := map[string]interface{}{
		"/defi/token_overview": wrapResponse(map[string]interface{}{
			"address":  "TokenMint123",
			"symbol":   "EXT",
			"name":     "Extensions Test",
			"decimals": 9,
			"extensions": map[string]interface{}{
				"coingecko":   "test-token",
				"twitter":     "https://twitter.com/testtoken",
				"website":     "https://testtoken.io",
				"telegram":    "https://t.me/testtoken",
				"discord":     "https://discord.gg/testtoken",
				"description": "A test token with all extensions",
			},
		}),
	}

	server := testServer(t, responses)
	defer server.Close()

	client := testClient(t, server.URL)
	overview, err := client.GetTokenOverview(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if overview.Extensions == nil {
		t.Fatal("expected extensions to be non-nil")
	}

	if overview.Extensions.Coingecko != "test-token" {
		t.Errorf("expected coingecko 'test-token', got '%s'", overview.Extensions.Coingecko)
	}

	if overview.Extensions.Twitter != "https://twitter.com/testtoken" {
		t.Errorf("expected twitter 'https://twitter.com/testtoken', got '%s'", overview.Extensions.Twitter)
	}

	if overview.Extensions.Website != "https://testtoken.io" {
		t.Errorf("expected website 'https://testtoken.io', got '%s'", overview.Extensions.Website)
	}

	if overview.Extensions.Telegram != "https://t.me/testtoken" {
		t.Errorf("expected telegram 'https://t.me/testtoken', got '%s'", overview.Extensions.Telegram)
	}

	if overview.Extensions.Discord != "https://discord.gg/testtoken" {
		t.Errorf("expected discord 'https://discord.gg/testtoken', got '%s'", overview.Extensions.Discord)
	}

	if overview.Extensions.Description != "A test token with all extensions" {
		t.Errorf("expected description, got '%s'", overview.Extensions.Description)
	}
}

func TestGetTokenOverview_EmptyAddress(t *testing.T) {
	client, _ := NewClient("test-key")
	_, err := client.GetTokenOverview(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty address")
	}

	apiErr, ok := IsAPIError(err)
	if !ok {
		t.Error("expected APIError")
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
}

func TestGetTokenOverview_NotFound(t *testing.T) {
	responses := map[string]interface{}{
		"/defi/token_overview": 404,
	}

	server := testServer(t, responses)
	defer server.Close()

	client := testClient(t, server.URL)
	_, err := client.GetTokenOverview(context.Background(), "unknown-token")
	if err == nil {
		t.Error("expected error for unknown token")
	}

	apiErr, ok := IsAPIError(err)
	if !ok {
		t.Error("expected APIError")
	}
	if !apiErr.IsNotFound() {
		t.Error("expected IsNotFound to be true")
	}
}

func TestGetTokenOverview_SuccessFalse(t *testing.T) {
	responses := map[string]interface{}{
		"/defi/token_overview": wrapFailure(),
	}

	server := testServer(t, responses)
	defer server.Close()

	client := testClient(t, server.URL)
	_, err := client.GetTokenOverview(context.Background(), "test-token")
	if err == nil {
		t.Error("expected error for success=false response")
	}
}
