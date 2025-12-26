# birdeye-go

[![CI](https://github.com/Laminar-Bot/birdeye-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Laminar-Bot/birdeye-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Laminar-Bot/birdeye-go.svg)](https://pkg.go.dev/github.com/Laminar-Bot/birdeye-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Laminar-Bot/birdeye-go)](https://goreportcard.com/report/github.com/Laminar-Bot/birdeye-go)
[![Coverage](https://img.shields.io/badge/coverage-90%25-brightgreen)](https://github.com/Laminar-Bot/birdeye-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go client for the [Birdeye](https://birdeye.so) API - comprehensive DeFi analytics and data for Solana.

## Features

- **Token Prices** - Real-time prices with `decimal.Decimal` precision
- **Token Security** - Authority checks, holder concentration, Token-2022 detection
- **Token Overview** - Market data, liquidity, volume, holder counts
- **Automatic Retries** - Exponential backoff for rate limits and server errors
- **Flexible Configuration** - Functional options pattern for clean API

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

    birdeye "github.com/Laminar-Bot/birdeye-go"
)

func main() {
    // Create client with API key
    client, err := birdeye.NewClient("your-api-key")
    if err != nil {
        log.Fatal(err)
    }

    // Get token price
    price, err := client.GetPrice(context.Background(), "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BONK: $%s\n", price.Value.String())
}
```

## Configuration Options

```go
import (
    "time"
    birdeye "github.com/Laminar-Bot/birdeye-go"
)

// Configure with options
client, err := birdeye.NewClient("your-api-key",
    birdeye.WithTimeout(30*time.Second),
    birdeye.WithMaxRetries(5),
    birdeye.WithBaseURL("https://custom-proxy.example.com"),
)
```

### Available Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithTimeout(d)` | HTTP request timeout | 10 seconds |
| `WithMaxRetries(n)` | Maximum retry attempts | 3 |
| `WithBaseURL(url)` | Custom API base URL | `https://public-api.birdeye.so` |
| `WithLogger(l)` | Custom logger implementation | No-op logger |
| `WithHTTPClient(c)` | Custom `*http.Client` | Default with timeout |

## Token Prices

```go
// Get single token price
price, err := client.GetPrice(ctx, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("USDC: $%s\n", price.Value.String())
fmt.Printf("24h Change: %s%%\n", price.PriceChange24h.String())

// Get multiple prices (automatically batched for >100 tokens)
prices, err := client.GetMultiplePrices(ctx, []string{
    "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263",
    "7GCihgDB8fe6KNjn2MYtkzZcRjQy3t9GHdC8uHYmW2hr",
})
for addr, price := range prices {
    fmt.Printf("%s: $%s\n", addr[:8], price.String())
}
```

## Token Security

Check for rug pull indicators:

```go
security, err := client.GetTokenSecurity(ctx, tokenAddress)
if err != nil {
    log.Fatal(err)
}

// Check authority status (active authority = higher risk)
if security.HasMintAuthority() {
    fmt.Println("WARNING: Token has active mint authority")
}
if security.HasFreezeAuthority() {
    fmt.Println("WARNING: Token has active freeze authority")
}

// Check holder concentration
fmt.Printf("Top 10 Holders: %s%%\n", security.Top10HolderPercent)

// Token-2022 specific checks
if security.IsToken2022 && security.TransferFeeEnable {
    fmt.Printf("Transfer Fee: %d bps\n", security.TransferFeeData.TransferFeeBPS)
}
```

## Token Overview

Get comprehensive market data:

```go
overview, err := client.GetTokenOverview(ctx, tokenAddress)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Token: %s (%s)\n", overview.Name, overview.Symbol)
fmt.Printf("Price: $%s\n", overview.Price.String())
fmt.Printf("Liquidity: $%s\n", overview.Liquidity.String())
fmt.Printf("24h Volume: $%s\n", overview.Volume24hUSD.String())
fmt.Printf("Market Cap: $%s\n", overview.MarketCap.String())
fmt.Printf("Holders: %d\n", overview.Holder)

// Social links (if available)
if overview.Extensions != nil {
    if overview.Extensions.Twitter != "" {
        fmt.Printf("Twitter: %s\n", overview.Extensions.Twitter)
    }
}
```

## Error Handling

All API errors are returned as `*APIError` with helpful methods:

```go
price, err := client.GetPrice(ctx, "invalid-token")
if err != nil {
    if apiErr, ok := birdeye.IsAPIError(err); ok {
        switch {
        case apiErr.IsNotFound():
            fmt.Println("Token not found")
        case apiErr.IsRateLimited():
            fmt.Println("Rate limited - slow down")
        case apiErr.IsServerError():
            fmt.Println("Birdeye server error")
        case apiErr.IsClientError():
            fmt.Printf("Bad request: %s\n", apiErr.Message)
        }
    }
    return
}
```

## Custom Logging

Implement the `Logger` interface for custom logging:

```go
type Logger interface {
    Debug(msg string, keysAndValues ...interface{})
    Info(msg string, keysAndValues ...interface{})
    Warn(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
}

// Example with zap
type zapAdapter struct{ *zap.SugaredLogger }

func (z *zapAdapter) Debug(msg string, kv ...interface{}) { z.Debugw(msg, kv...) }
func (z *zapAdapter) Info(msg string, kv ...interface{})  { z.Infow(msg, kv...) }
func (z *zapAdapter) Warn(msg string, kv ...interface{})  { z.Warnw(msg, kv...) }
func (z *zapAdapter) Error(msg string, kv ...interface{}) { z.Errorw(msg, kv...) }

client, _ := birdeye.NewClient("api-key",
    birdeye.WithLogger(&zapAdapter{sugar}),
)
```

## Rate Limits

Birdeye enforces API rate limits based on your plan. This client:

- Automatically retries on 429 (rate limit) responses
- Uses exponential backoff between retries
- Respects context cancellation

```go
// Configure retry behavior
client, _ := birdeye.NewClient("api-key",
    birdeye.WithMaxRetries(5),  // Up to 5 retry attempts
)
```

## Financial Precision

All price and amount values use `decimal.Decimal` from [shopspring/decimal](https://github.com/shopspring/decimal) to avoid floating-point precision issues:

```go
// Prices are decimal.Decimal, not float64
price := overview.Price                    // decimal.Decimal
liquidity := overview.Liquidity            // decimal.Decimal

// Safe arithmetic
total := price.Mul(decimal.NewFromInt(100))
formatted := price.StringFixed(8)  // "0.00001234"
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) first.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Links

- [Birdeye Documentation](https://docs.birdeye.so)
- [Birdeye App](https://birdeye.so)
- [Go Package Documentation](https://pkg.go.dev/github.com/Laminar-Bot/birdeye-go)
