package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ouroboros-crypto/node/x/coins/types"
)

type CoinCreatedHook = func(ctx sdk.Context, coin types.Coin, paymentForCreation sdk.Int)

// Adds new coin created hook
func (k *Keeper) AddCoinCreatedHook(hook CoinCreatedHook) {
	k.coinCreatedHooks = append(k.coinCreatedHooks, hook)
}

func (k Keeper) AfterCoinCreated(ctx sdk.Context, coin types.Coin, paymentForCreation sdk.Int) {
	for _, hook := range k.coinCreatedHooks {
		hook(ctx, coin, paymentForCreation)
	}
}
