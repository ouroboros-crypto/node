package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ouroboros-crypto/node/x/coins"
	"github.com/ouroboros-crypto/node/x/emission/types"
)


type Keeper struct {
	stakingKeeper staking.Keeper

	storeKey      sdk.StoreKey

	Cdc *codec.Codec
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, stakingKeeper staking.Keeper) Keeper {
	return Keeper{
		storeKey:      storeKey,
		Cdc:           cdc,
		stakingKeeper:           stakingKeeper,
	}
}


// Returns the emission record
func (k Keeper) GetEmission(ctx sdk.Context, coin coins.Coin) types.Emission {
	store := ctx.KVStore(k.storeKey)

	key := []byte(types.StoreKey)

	if !coin.Default {
		key = []byte(coin.Symbol)
	}

	if !store.Has(key) {
		return types.NewEmission()
	}

	var emission types.Emission

	k.Cdc.MustUnmarshalBinaryBare(store.Get(key), &emission)

	return emission
}

// Saves the emission record
func (k Keeper) SetEmission(ctx sdk.Context, emission types.Emission, coin coins.Coin) {
	store := ctx.KVStore(k.storeKey)

	key := []byte(types.StoreKey)

	if !coin.Default {
		key = []byte(coin.Symbol)
	}

	store.Set(key, k.Cdc.MustMarshalBinaryBare(emission))
}

// Checks if the threshold has been reached - in that case, we won't do posmining
func (k Keeper) IsThresholdReached(ctx sdk.Context, coin coins.Coin) bool {
	if coin.Default {
		return k.GetEmission(ctx, coin).IsThresholdReached()
	}

	return false
}

// Adding new coins to emission
func (k Keeper) Add(ctx sdk.Context, amount sdk.Int, coin coins.Coin) {
	emission := k.GetEmission(ctx, coin)
	emission.Current = emission.Current.Add(amount)

	k.SetEmission(ctx, emission, coin)
}

// Adding new custom coins to emission
func (k Keeper) AddCoin(ctx sdk.Context, amount sdk.Int, coin coins.Coin) {
	emission := k.GetEmission(ctx, coin)
	emission.Current = emission.Current.Add(amount)

	k.SetEmission(ctx, emission, coin)
}

// Remove coins from emission
func (k Keeper) Sub(ctx sdk.Context, amount sdk.Int, coin coins.Coin) {
	emission := k.GetEmission(ctx, coin)
	emission.Current = emission.Current.Sub(amount)

	k.SetEmission(ctx, emission, coin)
}

// Iterator
func (k Keeper) GetIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)

	return sdk.KVStorePrefixIterator(store, nil)
}

// Remove coins from emission before slashing
func (k Keeper) UpdateBeforeSlashing(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	// Since we've disabled the fees, there is nothing to do
	return
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)

	if !found {
		panic("Validator does not exist but got slashed")
	}

	amount := validator.GetTokens()

	slashAmountDec := sdk.NewInt(amount.ToDec().Mul(fraction).Int64())


	if slashAmountDec.IsPositive() {
		k.Sub(ctx, slashAmountDec, coins.GetDefaultCoin())
	}
}
