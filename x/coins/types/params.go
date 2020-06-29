package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
)

// Parameter store keys
var (
	CreationPriceInUSD = sdk.NewInt(5000) // 50.00 USD
)

// ParamKeyTable for coins module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for coins at genesis
type Params struct {
}

// NewParams creates a new Params object
func NewParams() Params {
	return Params{
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`
	// TODO: Return all the params as a string
	`, )
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams( /* TODO: Pass in your default Params */ )
}
