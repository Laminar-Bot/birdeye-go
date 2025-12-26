// Package birdeye provides a Go client for the Birdeye DeFi analytics API.
//
// Birdeye (https://birdeye.so) provides comprehensive analytics for Solana tokens
// including price data, liquidity metrics, holder information, and security analysis.
//
// # Quick Start
//
//	client, err := birdeye.NewClient("your-api-key")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Get token price
//	price, err := client.GetPrice(ctx, "So11111111111111111111111111111111111111112")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("SOL: $%s\n", price.Value)
//
// # Configuration Options
//
// Use functional options to customize the client:
//
//	client, err := birdeye.NewClient("your-api-key",
//	    birdeye.WithTimeout(30 * time.Second),
//	    birdeye.WithMaxRetries(5),
//	    birdeye.WithBaseURL("https://custom-endpoint.example.com"),
//	)
//
// # Error Handling
//
// API errors are returned as *APIError which provides helper methods:
//
//	price, err := client.GetPrice(ctx, tokenAddress)
//	if err != nil {
//	    if apiErr, ok := birdeye.IsAPIError(err); ok {
//	        if apiErr.IsRateLimited() {
//	            // Handle rate limiting
//	        }
//	        if apiErr.IsNotFound() {
//	            // Token not found
//	        }
//	    }
//	    return err
//	}
//
// # Logging
//
// Logging is optional. To enable, provide a Logger implementation:
//
//	client, err := birdeye.NewClient("your-api-key",
//	    birdeye.WithLogger(myLogger),
//	)
//
// The Logger interface is minimal and can wrap any logging library.
//
// # Financial Precision
//
// All monetary values use github.com/shopspring/decimal for precise arithmetic.
// Never use float64 for financial calculations.
package birdeye
