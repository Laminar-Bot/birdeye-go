package birdeye

import (
	"context"
	"net/url"
)

// TokenSecurity contains security-related information about a token.
//
// This data is used to check for red flags like:
//   - Active mint authority (can create more tokens)
//   - Active freeze authority (can freeze user accounts)
//   - High holder concentration (rug pull risk)
type TokenSecurity struct {
	// MintAuthority is the address that can mint new tokens.
	// Nil or empty string means no mint authority (safer).
	MintAuthority *string `json:"mintAuthority"`

	// FreezeAuthority is the address that can freeze token accounts.
	// Nil or empty string means no freeze authority (safer).
	FreezeAuthority *string `json:"freezeAuthority"`

	// CreatorAddress is the token creator's wallet address.
	CreatorAddress string `json:"creatorAddress"`

	// CreatorBalance is the creator's current token balance.
	CreatorBalance string `json:"creatorBalance"`

	// CreatorPercentage is the percentage of supply held by creator.
	CreatorPercentage string `json:"creatorPercentage"`

	// OwnerAddress is the current owner of the token mint.
	OwnerAddress string `json:"ownerAddress"`

	// OwnerBalance is the owner's current token balance.
	OwnerBalance string `json:"ownerBalance"`

	// OwnerPercentage is the percentage of supply held by owner.
	OwnerPercentage string `json:"ownerPercentage"`

	// Top10HolderBalance is the combined balance of top 10 holders.
	Top10HolderBalance string `json:"top10HolderBalance"`

	// Top10HolderPercent is the percentage of supply held by top 10 holders.
	Top10HolderPercent string `json:"top10HolderPercent"`

	// Top10UserBalance is the combined balance of top 10 non-contract holders.
	Top10UserBalance string `json:"top10UserBalance"`

	// Top10UserPercent is the percentage of supply held by top 10 users.
	Top10UserPercent string `json:"top10UserPercent"`

	// TotalSupply is the total token supply.
	TotalSupply string `json:"totalSupply"`

	// IsToken2022 indicates if this is a Token-2022 (new token program) token.
	IsToken2022 bool `json:"isToken2022"`

	// TransferFeeEnable indicates if transfer fees are enabled (Token-2022).
	TransferFeeEnable bool `json:"transferFeeEnable"`

	// TransferFeeData contains fee configuration if enabled.
	TransferFeeData *TransferFeeData `json:"transferFeeData,omitempty"`

	// NonTransferable indicates if the token is non-transferable (soulbound).
	NonTransferable bool `json:"nonTransferable"`

	// MutableMetadata indicates if token metadata can be changed.
	MutableMetadata bool `json:"mutableMetadata"`
}

// TransferFeeData contains Token-2022 transfer fee configuration.
type TransferFeeData struct {
	// TransferFeeBPS is the fee in basis points (100 = 1%).
	TransferFeeBPS int `json:"transferFeeBps"`

	// MaxFee is the maximum fee amount.
	MaxFee string `json:"maxFee"`

	// FeeAuthority can update the fee configuration.
	FeeAuthority string `json:"feeAuthority"`

	// WithdrawAuthority can withdraw collected fees.
	WithdrawAuthority string `json:"withdrawAuthority"`
}

// HasMintAuthority returns true if the token has an active mint authority.
//
// Tokens with mint authority are higher risk because more tokens
// can be minted at any time, diluting existing holders.
func (ts *TokenSecurity) HasMintAuthority() bool {
	return ts.MintAuthority != nil && *ts.MintAuthority != ""
}

// HasFreezeAuthority returns true if the token has an active freeze authority.
//
// Tokens with freeze authority are higher risk because user accounts
// can be frozen, preventing transfers or sales.
func (ts *TokenSecurity) HasFreezeAuthority() bool {
	return ts.FreezeAuthority != nil && *ts.FreezeAuthority != ""
}

// GetTokenSecurity fetches security information for a token.
//
// This endpoint provides data for token screening:
//   - Mint and freeze authority status
//   - Holder concentration (top 10 holders percentage)
//   - Creator/owner holdings
//   - Token-2022 specific features (transfer fees, etc.)
//
// Example:
//
//	security, err := client.GetTokenSecurity(ctx, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
//	if err != nil {
//	    return err
//	}
//	if security.HasMintAuthority() {
//	    log.Warn("token has active mint authority")
//	}
func (c *Client) GetTokenSecurity(ctx context.Context, address string) (*TokenSecurity, error) {
	if address == "" {
		return nil, &APIError{
			StatusCode: 400,
			Message:    "address is required",
			Path:       "/defi/token_security",
		}
	}

	params := url.Values{}
	params.Set("address", address)

	body, err := c.doGet(ctx, "/defi/token_security", params)
	if err != nil {
		return nil, err
	}

	security, err := parseResponse[TokenSecurity](body)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("fetched token security",
		"address", address,
		"has_mint_auth", security.HasMintAuthority(),
		"has_freeze_auth", security.HasFreezeAuthority(),
		"creator_pct", security.CreatorPercentage,
		"top10_pct", security.Top10HolderPercent,
	)

	return security, nil
}
