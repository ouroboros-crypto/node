package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func GetGenesisWallet() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(GenesisWallet)

	return addr
}
