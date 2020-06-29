package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Returns default ouro coin
func GetDefaultCoin() Coin {
	// Default ouroboros record - threshold is 2m coins
	return Coin{Name: "Ouroboros", Symbol: "ouro", Default: true, PosminingThreshold: sdk.NewIntWithDecimal(2000000, 6)}
}

// Returns a stub record of the coin
func GetCoinStub(symbol string) Coin {
	if symbol == "ouro" {
		return GetDefaultCoin()
	}

	return Coin{ Symbol: symbol, Default: false, PosminingThreshold: sdk.NewInt(0)}
}

// Returns default sdk coins
func GetDefaultCoins(amnt sdk.Int) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin("ouro", amnt))
}

func GetGenesisWallet() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32("ouro1wjc2pg9e3p8sdll4kg9ssj44nhm7ce4prapsxf")

	return addr
}
