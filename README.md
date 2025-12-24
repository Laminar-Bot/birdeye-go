# birdeye-go

[![Go Reference](https://pkg.go.dev/badge/github.com/Laminar-Bot/birdeye-go.svg)](https://pkg.go.dev/github.com/Laminar-Bot/birdeye-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Laminar-Bot/birdeye-go)](https://goreportcard.com/report/github.com/Laminar-Bot/birdeye-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go client for the [Birdeye](https://birdeye.so) API - comprehensive DeFi analytics and data for Solana.

## Features

- 💰 **Token Prices** - Real-time and historical prices
- 📊 **Token Analytics** - Volume, liquidity, holder stats
- 🔒 **Security Info** - LP status, authority checks, holder concentration
- 👛 **Wallet Tracking** - Portfolio and transaction history
- 🏆 **Top Traders** - Discover profitable wallets
- 📈 **OHLCV Data** - Candlestick data for charting

## Installation
```bash
go get github.com/Laminar-Bot/birdeye-go
```

## Quick Start
```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Laminar-Bot/birdeye-go"
)

func main() {
    client := birdeye.NewClient(birdeye.Config{
        APIKey: "your-api-key",
    })

    // Get token price
    price, err := client.GetPrice(context.Background(), "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BONK: $%f\n", price.Value)
}
```

## Token Security & Liquidity
```go
// Check if a token is safe
security, err := client.GetTokenLiquidity(ctx, tokenAddress)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Mint Authority Revoked: %v\n", security.IsMintAuthorityRevoked())
fmt.Printf("Freeze Authority Revoked: %v\n", security.IsFreezeAuthorityRevoked())
fmt.Printf("LP Burned: %.1f%%\n", security.LPBurnedPct)
fmt.Printf("Total Liquidity: $%.2f\n", security.TotalLiquidityUSD)
fmt.Printf("Top 10 Holders: %.1f%%\n", security.Top10HolderPercent)
```

## Wallet Analysis
```go
// Get wallet portfolio
portfolio, err := client.GetWalletPortfolio(ctx, walletAddress)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Value: $%.2f\n", portfolio.TotalUSD)
for _, token := range portfolio.Items {
    fmt.Printf("  %s: $%.2f\n", token.Symbol, token.ValueUSD)
}

// Get wallet transactions
txs, err := client.GetWalletTransactions(ctx, walletAddress, &birdeye.TxOptions{
    Limit:  50,
    TxType: "swap",
})
```

## Top Traders
```go
// Find successful traders of a token
traders, err := client.GetTopTraders(ctx, tokenAddress, &birdeye.TopTradersRequest{
    TimeFrame: "24h",
    SortBy:    "pnl",
    Limit:     20,
})

for _, trader := range traders {
    fmt.Printf("%s: PnL $%.2f (%.1f%% win rate)\n", 
        trader.Address[:8], trader.PnL, trader.WinRate)
}
```

## Price History (OHLCV)
```go
// Get candlestick data
ohlcv, err := client.GetOHLCV(ctx, tokenAddress, &birdeye.OHLCVRequest{
    TimeFrom: time.Now().Add(-24 * time.Hour).Unix(),
    TimeTo:   time.Now().Unix(),
    Type:     "15m", // 1m, 5m, 15m, 1H, 4H, 1D
})

for _, candle := range ohlcv.Items {
    fmt.Printf("%s O:%.6f H:%.6f L:%.6f C:%.6f V:%.0f\n",
        time.Unix(candle.UnixTime, 0).Format("15:04"),
        candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
}
```

## Batch Price Queries
```go
// Get multiple prices efficiently
prices, err := client.GetMultiplePrices(ctx, []string{
    "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263",
    "7GCihgDB8fe6KNjn2MYtkzZcRjQy3t9GHdC8uHYmW2hr",
    "EKpQGSJtjMFqKZ9KQanSqYXRcF8fBopzLHYxdM65zcjm",
})

for addr, price := range prices {
    fmt.Printf("%s: $%f\n", addr[:8], price.Value)
}
```

## Configuration
```go
client := birdeye.NewClient(birdeye.Config{
    APIKey:     "your-api-key",
    BaseURL:    "https://public-api.birdeye.so", // default
    TimeoutSec: 30,
    Chain:      "solana", // default, also supports "ethereum", "bsc", etc.
})
```

## Rate Limits

Birdeye has API rate limits based on your plan. This client includes automatic retry with backoff for rate limit errors.
```go
client := birdeye.NewClient(birdeye.Config{
    APIKey:        "your-api-key",
    MaxRetries:    3,
    RetryDelaySec: 1,
})
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) first.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Links

- [Birdeye Documentation](https://docs.birdeye.so)
- [Birdeye App](https://birdeye.so)
- [Go Package Documentation](https://pkg.go.dev/github.com/Laminar-Bot/birdeye-go)
