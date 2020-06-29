package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/ouroboros-crypto/node/x/coins/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the coins store
type Keeper struct {
	BankKeeper bank.Keeper
	storeKey   sdk.StoreKey
	cdc        *codec.Codec

	coinCreatedHooks []CoinCreatedHook
}

// NewKeeper creates a coins keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bankKeeper bank.Keeper) Keeper {
	keeper := Keeper{
		BankKeeper: bankKeeper,
		storeKey:   key,
		cdc:        cdc,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}


// Updates the required creation price based on the current coin price
func (k Keeper) UpdateCreationPrice (ctx sdk.Context, currentPrice sdk.Int) sdk.Int {
	newPrice := types.CreationPriceInUSD.Quo(currentPrice).MulRaw(1e6)

	k.SetCreationPrice(ctx, types.CreationPrice{
		Updated: ctx.BlockTime(),
		Price:   newPrice,
	})

	return newPrice
}