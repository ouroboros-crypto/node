package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func GetGenesisWallet() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(GenesisWallet)

	return addr
}

func GetRegulationWallet() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32("ouro1qxj94wtuet8ct3cdulw4f47lnmfjxweums0zg8")

	return addr
}
