package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// `from` and `to` can be nill
type CoinsTransferHook = func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins)


func (k *Keeper) AddBeforeHook(hook CoinsTransferHook) {
	k.beforeTransferHooks = append(k.beforeTransferHooks, hook)
}

func (k *Keeper) AddAfterHook(hook CoinsTransferHook) {
	k.afterTransferHooks = append(k.afterTransferHooks, hook)
}


func (k Keeper) beforeCoinsTransfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) {
	for _, hook := range k.beforeTransferHooks {
		hook(ctx, from, to, amt)
	}
}

func (k Keeper) afterCoinsTransfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amt sdk.Coins) {
	for _, hook := range k.afterTransferHooks {
		hook(ctx, from, to, amt)
	}
}