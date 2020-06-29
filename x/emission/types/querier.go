package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ouroboros-crypto/node/x/coins"
)

// Query endpoints supported by the emission querier
const (
	QueryGetEmission = "get"
)

// Returns emission of the coin
type QueryResGetEmission struct {
	Current sdk.Int `json:"current"`
	Threshold sdk.Int `json:"threshold"`
	Coin string `json:"coin"`
}

func NewQueryResGetEmission(emission Emission, coin coins.Coin) QueryResGetEmission {
	return QueryResGetEmission{
		Current: emission.Current,
		Coin:    coin.Symbol,
	}
}