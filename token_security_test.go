package birdeye

import (
	"context"
	"testing"
)

func TestTokenSecurity_HasMintAuthority(t *testing.T) {
	tests := []struct {
		name          string
		mintAuthority *string
		expected      bool
	}{
		{"nil authority", nil, false},
		{"empty string", strPtr(""), false},
		{"valid authority", strPtr("SomeAddress123"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TokenSecurity{MintAuthority: tt.mintAuthority}
			if ts.HasMintAuthority() != tt.expected {
				t.Errorf("HasMintAuthority() = %v, expected %v", ts.HasMintAuthority(), tt.expected)
			}
		})
	}
}

func TestTokenSecurity_HasFreezeAuthority(t *testing.T) {
	tests := []struct {
		name            string
		freezeAuthority *string
		expected        bool
	}{
		{"nil authority", nil, false},
		{"empty string", strPtr(""), false},
		{"valid authority", strPtr("SomeAddress123"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TokenSecurity{FreezeAuthority: tt.freezeAuthority}
			if ts.HasFreezeAuthority() != tt.expected {
				t.Errorf("HasFreezeAuthority() = %v, expected %v", ts.HasFreezeAuthority(), tt.expected)
			}
		})
	}
}

func TestGetTokenSecurity_Success(t *testing.T) {
	mintAuth := "MintAuth123"
	responses := map[string]interface{}{
		"/defi/token_security": wrapResponse(map[string]interface{}{
			"mintAuthority":      mintAuth,
			"freezeAuthority":    nil,
			"creatorAddress":     "CreatorAddr",
			"creatorBalance":     "1000000",
			"creatorPercentage":  "5.5",
			"ownerAddress":       "OwnerAddr",
			"ownerBalance":       "500000",
			"ownerPercentage":    "2.5",
			"top10HolderBalance": "5000000",
			"top10HolderPercent": "25.0",
			"top10UserBalance":   "4000000",
			"top10UserPercent":   "20.0",
			"totalSupply":        "20000000",
			"isToken2022":        false,
			"transferFeeEnable":  false,
			"nonTransferable":    false,
			"mutableMetadata":    true,
		}),
	}

	server := testServer(t, responses)
	defer server.Close()

	client := testClient(t, server.URL)
	security, err := client.GetTokenSecurity(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !security.HasMintAuthority() {
		t.Error("expected HasMintAuthority to be true")
	}

	if security.HasFreezeAuthority() {
		t.Error("expected HasFreezeAuthority to be false")
	}

	if security.CreatorAddress != "CreatorAddr" {
		t.Errorf("expected creatorAddress 'CreatorAddr', got '%s'", security.CreatorAddress)
	}

	if security.Top10HolderPercent != "25.0" {
		t.Errorf("expected top10HolderPercent '25.0', got '%s'", security.Top10HolderPercent)
	}

	if security.MutableMetadata != true {
		t.Error("expected mutableMetadata to be true")
	}
}

func TestGetTokenSecurity_Token2022(t *testing.T) {
	responses := map[string]interface{}{
		"/defi/token_security": wrapResponse(map[string]interface{}{
			"mintAuthority":     nil,
			"freezeAuthority":   nil,
			"creatorAddress":    "CreatorAddr",
			"creatorBalance":    "0",
			"creatorPercentage": "0",
			"ownerAddress":      "OwnerAddr",
			"ownerBalance":      "0",
			"ownerPercentage":   "0",
			"totalSupply":       "1000000000",
			"isToken2022":       true,
			"transferFeeEnable": true,
			"transferFeeData": map[string]interface{}{
				"transferFeeBps":    100,
				"maxFee":            "1000000",
				"feeAuthority":      "FeeAuth",
				"withdrawAuthority": "WithdrawAuth",
			},
			"nonTransferable": false,
			"mutableMetadata": false,
		}),
	}

	server := testServer(t, responses)
	defer server.Close()

	client := testClient(t, server.URL)
	security, err := client.GetTokenSecurity(context.Background(), "token2022-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !security.IsToken2022 {
		t.Error("expected IsToken2022 to be true")
	}

	if !security.TransferFeeEnable {
		t.Error("expected TransferFeeEnable to be true")
	}

	if security.TransferFeeData == nil {
		t.Fatal("expected TransferFeeData to be non-nil")
	}

	if security.TransferFeeData.TransferFeeBPS != 100 {
		t.Errorf("expected TransferFeeBPS 100, got %d", security.TransferFeeData.TransferFeeBPS)
	}

	if security.TransferFeeData.FeeAuthority != "FeeAuth" {
		t.Errorf("expected FeeAuthority 'FeeAuth', got '%s'", security.TransferFeeData.FeeAuthority)
	}
}

func TestGetTokenSecurity_EmptyAddress(t *testing.T) {
	client, _ := NewClient("test-key")
	_, err := client.GetTokenSecurity(context.Background(), "")
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

func TestGetTokenSecurity_NotFound(t *testing.T) {
	responses := map[string]interface{}{
		"/defi/token_security": 404,
	}

	server := testServer(t, responses)
	defer server.Close()

	client := testClient(t, server.URL)
	_, err := client.GetTokenSecurity(context.Background(), "unknown-token")
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

func TestGetTokenSecurity_SuccessFalse(t *testing.T) {
	responses := map[string]interface{}{
		"/defi/token_security": wrapFailure(),
	}

	server := testServer(t, responses)
	defer server.Close()

	client := testClient(t, server.URL)
	_, err := client.GetTokenSecurity(context.Background(), "test-token")
	if err == nil {
		t.Error("expected error for success=false response")
	}
}

// strPtr is a helper to create string pointers for tests.
func strPtr(s string) *string {
	return &s
}
