package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	OURO                 = "ouro"
	POINTS               = 6        // points after comma
	INITIAL              = 10000000 // initial emission
	PARAMINING_THRESHOLD = 2000000  // posmining threshold
)

func NewOuroCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(OURO, amount)
}

func NewCoin(amount sdk.Int) sdk.Coin {
	return sdk.NewCoin(OURO, amount)
}
func NewCoins(amount sdk.Int) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(OURO, amount))
}

// Since we cannot make constans as objects
func GetMaxLevel() sdk.Int {
	return sdk.NewInt(100)
}

func GetPosminingThreshold() sdk.Int {
	return sdk.NewIntWithDecimal(PARAMINING_THRESHOLD, POINTS)
}
