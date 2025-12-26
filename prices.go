package birdeye

import (
	"context"
	"net/url"
	"strings"

	"github.com/shopspring/decimal"
)

// PriceData contains price information for a single token.
type PriceData struct {
	// Value is the current price in USD.
	Value decimal.Decimal `json:"value"`

	// UpdateUnixTime is when the price was last updated (Unix timestamp).
	UpdateUnixTime int64 `json:"updateUnixTime"`

	// UpdateHumanTime is a human-readable update timestamp.
	UpdateHumanTime string `json:"updateHumanTime"`

	// PriceChange24h is the 24-hour price change percentage.
	PriceChange24h decimal.Decimal `json:"priceChange24h"`
}

// GetPrice fetches the current price for a single token.
//
// Example:
//
//	price, err := client.GetPrice(ctx, "So11111111111111111111111111111111111111112")
//	if err != nil {
//	    return err
//	}
//	log.Printf("SOL price: $%s", price.Value)
func (c *Client) GetPrice(ctx context.Context, address string) (*PriceData, error) {
	if address == "" {
		return nil, &APIError{
			StatusCode: 400,
			Message:    "address is required",
			Path:       "/defi/price",
		}
	}

	params := url.Values{}
	params.Set("address", address)

	body, err := c.doGet(ctx, "/defi/price", params)
	if err != nil {
		return nil, err
	}

	price, err := parseResponse[PriceData](body)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("fetched token price",
		"address", address,
		"price", price.Value.String(),
		"change_24h", price.PriceChange24h.String(),
	)

	return price, nil
}

// GetMultiplePrices fetches prices for multiple tokens in a single request.
//
// Birdeye supports up to 100 addresses per request. This method handles
// batching automatically for larger lists.
//
// Returns a map of address -> price. Missing prices are omitted from the result.
//
// Example:
//
//	prices, err := client.GetMultiplePrices(ctx, []string{
//	    "So11111111111111111111111111111111111111112", // SOL
//	    "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
//	})
//	for addr, price := range prices {
//	    log.Printf("%s: $%s", addr, price)
//	}
func (c *Client) GetMultiplePrices(ctx context.Context, addresses []string) (map[string]decimal.Decimal, error) {
	if len(addresses) == 0 {
		return make(map[string]decimal.Decimal), nil
	}

	// Validate no empty addresses in the list.
	for _, addr := range addresses {
		if addr == "" {
			return nil, &APIError{
				StatusCode: 400,
				Message:    "address list contains empty string",
				Path:       "/defi/multi_price",
			}
		}
	}

	const batchSize = 100
	result := make(map[string]decimal.Decimal, len(addresses))

	// Process addresses in batches of 100.
	for i := 0; i < len(addresses); i += batchSize {
		end := i + batchSize
		if end > len(addresses) {
			end = len(addresses)
		}

		batch := addresses[i:end]
		listAddress := strings.Join(batch, ",")

		params := url.Values{}
		params.Set("list_address", listAddress)

		body, err := c.doGet(ctx, "/defi/multi_price", params)
		if err != nil {
			return nil, err
		}

		// Multi-price response is a map of address -> price directly.
		batchPrices, err := parseResponse[map[string]decimal.Decimal](body)
		if err != nil {
			return nil, err
		}

		for addr, price := range *batchPrices {
			result[addr] = price
		}
	}

	c.logger.Debug("fetched multiple token prices",
		"requested", len(addresses),
		"received", len(result),
	)

	return result, nil
}
