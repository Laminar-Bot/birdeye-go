package birdeye

import (
	"context"
	"net/url"

	"github.com/shopspring/decimal"
)

// TokenOverview contains market and metadata information about a token.
//
// This data is used to check:
//   - Liquidity (minimum threshold for trading)
//   - Trading volume (activity indicator)
//   - Market cap and holder count
type TokenOverview struct {
	// Address is the token's mint address.
	Address string `json:"address"`

	// Symbol is the token's trading symbol (e.g., "SOL").
	Symbol string `json:"symbol"`

	// Name is the token's full name (e.g., "Solana").
	Name string `json:"name"`

	// Decimals is the number of decimal places for the token.
	Decimals int `json:"decimals"`

	// LogoURI is a URL to the token's logo image.
	LogoURI string `json:"logoURI"`

	// Liquidity is the total liquidity in USD across all pools.
	Liquidity decimal.Decimal `json:"liquidity"`

	// Price is the current price in USD.
	Price decimal.Decimal `json:"price"`

	// PriceChange24hPercent is the 24-hour price change percentage.
	PriceChange24hPercent decimal.Decimal `json:"priceChange24hPercent"`

	// Volume24h is the 24-hour trading volume in the token's native units.
	Volume24h decimal.Decimal `json:"v24h"`

	// Volume24hUSD is the 24-hour trading volume in USD.
	Volume24hUSD decimal.Decimal `json:"v24hUSD"`

	// Volume24hChangePercent is the change in volume vs previous 24h.
	Volume24hChangePercent decimal.Decimal `json:"v24hChangePercent"`

	// MarketCap is the market capitalization in USD.
	MarketCap decimal.Decimal `json:"mc"`

	// Supply is the total token supply.
	Supply decimal.Decimal `json:"supply"`

	// CirculatingSupply is the supply in circulation.
	CirculatingSupply decimal.Decimal `json:"circulatingSupply"`

	// Holder is the number of unique token holders.
	Holder int `json:"holder"`

	// Trade24h is the number of trades in the last 24 hours.
	Trade24h int `json:"trade24h"`

	// Trade24hChangePercent is the change in trade count vs previous 24h.
	Trade24hChangePercent decimal.Decimal `json:"trade24hChangePercent"`

	// Buy24h is the number of buy trades in the last 24 hours.
	Buy24h int `json:"buy24h"`

	// Sell24h is the number of sell trades in the last 24 hours.
	Sell24h int `json:"sell24h"`

	// UniqueWallet24h is the number of unique wallets trading in 24h.
	UniqueWallet24h int `json:"uniqueWallet24h"`

	// UniqueWallet24hChangePercent is the change vs previous 24h.
	UniqueWallet24hChangePercent decimal.Decimal `json:"uniqueWallet24hChangePercent"`

	// LastTradeUnixTime is the Unix timestamp of the last trade.
	LastTradeUnixTime int64 `json:"lastTradeUnixTime"`

	// LastTradeHumanTime is a human-readable timestamp of the last trade.
	LastTradeHumanTime string `json:"lastTradeHumanTime"`

	// Extensions contains optional metadata links.
	Extensions *TokenExtensions `json:"extensions,omitempty"`
}

// TokenExtensions contains optional metadata and social links.
type TokenExtensions struct {
	// Coingecko is the CoinGecko listing ID.
	Coingecko string `json:"coingecko,omitempty"`

	// Twitter is the project's Twitter handle or URL.
	Twitter string `json:"twitter,omitempty"`

	// Website is the project's website URL.
	Website string `json:"website,omitempty"`

	// Telegram is the project's Telegram group URL.
	Telegram string `json:"telegram,omitempty"`

	// Discord is the project's Discord server URL.
	Discord string `json:"discord,omitempty"`

	// Description is a brief description of the token.
	Description string `json:"description,omitempty"`
}

// GetTokenOverview fetches market overview data for a token.
//
// This endpoint provides data for token screening:
//   - Liquidity in USD (for minimum liquidity checks)
//   - Volume data (activity indicator)
//   - Holder count (distribution indicator)
//   - Price and market cap
//
// Example:
//
//	overview, err := client.GetTokenOverview(ctx, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
//	if err != nil {
//	    return err
//	}
//	if overview.Liquidity.LessThan(decimal.NewFromInt(50000)) {
//	    log.Warn("liquidity below threshold")
//	}
func (c *Client) GetTokenOverview(ctx context.Context, address string) (*TokenOverview, error) {
	if address == "" {
		return nil, &APIError{
			StatusCode: 400,
			Message:    "address is required",
			Path:       "/defi/token_overview",
		}
	}

	params := url.Values{}
	params.Set("address", address)

	body, err := c.doGet(ctx, "/defi/token_overview", params)
	if err != nil {
		return nil, err
	}

	overview, err := parseResponse[TokenOverview](body)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("fetched token overview",
		"address", address,
		"symbol", overview.Symbol,
		"name", overview.Name,
		"liquidity", overview.Liquidity.String(),
		"volume_24h", overview.Volume24hUSD.String(),
		"holders", overview.Holder,
	)

	return overview, nil
}
