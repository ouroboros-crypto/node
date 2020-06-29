package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Creation price
type CreationPrice struct {
	Updated time.Time `json:"updated"`
	Price sdk.Int `json:"updated"`
}

type Coin struct {
	Creator            sdk.AccAddress           `json:"creator" yaml:"creator"`                         // address of the coin creator
	Name               string                   `json:"string" yaml:"string"`                           // name of the coin
	Symbol             string                   `json:"symbol" yaml:"symbol"`                           // identifier of the coin
	Emission           sdk.Int                  `json:"emission" yaml:"emission"`                       // initial emission of the coin
	Description        string                   `json:"description" yaml:"description"`                 // description of the coin
	PosminingEnabled   bool                     `json:"posmining_enabled" yaml:"posmining_enabled"`     // if posmining should be enabled
	PosminingBalance   []CoinBalancePosmining   `json:"posmining_balance" yaml:"posmining_balance"`     // all the daily percent conditions
	PosminingStructure []CoinStructurePosmining `json:"posmining_structure" yaml:"posmining_structure"` // all the structure coffs
	PosminingThreshold sdk.Int `json:"posmining_threshold" yaml:"posmining_threshold"` // Posmining threshold

	Default bool `json:"default" yaml:"default"` // if coin is default ouro
}

// Represents every posmining condition based on the balane
type CoinBalancePosmining struct {
	FromAmount   sdk.Int `json:"from_amount" yaml:"from_amount"`     // range
	ToAmount     sdk.Int `json:"to_amount" yaml:"to_amount"`         // range
	DailyPercent sdk.Int `json:"daily_percent" yaml:"daily_percent"` // Daily percent
}

// Represents every posmining condition based on the balane
type CoinStructurePosmining struct {
	FromAmount sdk.Int `json:"from_amount" yaml:"from_amount"` // range
	ToAmount   sdk.Int `json:"to_amount" yaml:"to_amount"`     // range
	Coff       sdk.Int `json:"coff" yaml:"coff"`               // Coff
}

// implement fmt.Stringer
func (c Coin) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Creator: %s
	Name: %s
	Symbol: %s
	Emission: %s
	Description: %s`,
		c.Creator,
		c.Name,
		c.Symbol,
		c.Emission,
		c.Description,
	))
}

// Returns daily percent based on the amount of the coins
func (c Coin) GetDailyPercent(amnt sdk.Int) sdk.Int {
	if c.Default {
		return GetDailyPercent(amnt)
	}

	for _, b := range c.PosminingBalance {
		if amnt.GTE(b.FromAmount) && amnt.LTE(b.ToAmount) {
			return b.DailyPercent
		}
	}

	return sdk.NewInt(0)
}

// Returns daily percent based on the volume of the structure
func (c Coin) GetStructureCoff(amnt sdk.Int) sdk.Int {
	if c.Default {
		return GetStructureCoff(amnt)
	}

	for _, b := range c.PosminingStructure {
		if amnt.GTE(b.FromAmount) && amnt.LTE(b.ToAmount) {
			return b.Coff
		}
	}

	return sdk.NewInt(0)
}

// implement fmt.Stringer
func (cbp CoinBalancePosmining) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Range: %s - %s
	Daily Percent: %s`,
		cbp.FromAmount,
		cbp.ToAmount,
		cbp.DailyPercent,
	))
}

func InBetween(i sdk.Int, minRaw, maxRaw int64) bool {
	min := sdk.NewInt(minRaw)
	max := sdk.NewInt(maxRaw)

	if i.GTE(min) && i.LTE(max) {
		return true
	} else {
		return false
	}
}

// Возвращает коэффициент структуры в зависимости от баланса
func GetStructureCoff(balance sdk.Int) sdk.Int {
	if balance.LT(sdk.NewInt(1000000000)) {
		return sdk.NewInt(0)
	}

	if InBetween(balance, 1000000000, 9999999999) {
		return sdk.NewInt(218)
	}

	if InBetween(balance, 10000000000, 99999999999) {
		return sdk.NewInt(236)
	}

	if InBetween(balance, 100000000000, 999999999999) {
		return sdk.NewInt(277)
	}

	if InBetween(balance, 1000000000000, 9999999999999) {
		return sdk.NewInt(305)
	}

	if InBetween(balance, 10000000000000, 99999999999999) {
		return sdk.NewInt(336)
	}

	if InBetween(balance, 100000000000000, 999999999999999) {
		return sdk.NewInt(388)
	}

	return sdk.NewInt(437)
}

// Возвращает дневной процент в зависимости от баланса
func GetDailyPercent(balance sdk.Int) sdk.Int {
	if balance.LT(sdk.NewInt(10000)) {
		return sdk.NewInt(0)
	}

	if InBetween(balance, 10000, 99999999) {
		return sdk.NewInt(6)
	}

	if InBetween(balance, 100000000, 999999999) {
		return sdk.NewInt(7)
	}

	if InBetween(balance, 1000000000, 9999999999) {
		return sdk.NewInt(9)
	}

	if InBetween(balance, 10000000000, 49999999999) {
		return sdk.NewInt(10)
	}

	if InBetween(balance, 50000000000, 99999999999) {
		return sdk.NewInt(12)
	}

	if InBetween(balance, 100000000000, 499999999999) {
		return sdk.NewInt(14)
	}

	return sdk.NewInt(16)
}