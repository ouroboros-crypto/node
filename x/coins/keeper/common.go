package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ouroboros-crypto/node/x/coins/types"
)

// Saves coin
func (k Keeper) SetCoin(ctx sdk.Context, coin types.Coin) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(coin)

	key := []byte(types.CoinPrefix + coin.Symbol)

	store.Set(key, bz)
}

// GetCoin returns the coin information
func (k Keeper) GetCoin(ctx sdk.Context, coinId string) (types.Coin, error) {
	if coinId == "ouro" {
		return types.GetDefaultCoin(), nil
	}

	store := ctx.KVStore(k.storeKey)
	var coin types.Coin
	byteKey := []byte(types.CoinPrefix + coinId)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &coin)

	if err != nil {
		return coin, err
	}

	return coin, nil
}


// Returns the creation price record
func (k Keeper) GetCreationPrice(ctx sdk.Context) (types.CreationPrice, error) {
	store := ctx.KVStore(k.storeKey)

	var result types.CreationPrice

	byteKey := []byte("creation-price")

	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &result)

	if err != nil {
		return result, err
	}

	return result, nil
}


// Sets the creation price record
func (k Keeper) SetCreationPrice(ctx sdk.Context, price types.CreationPrice) {
	store := ctx.KVStore(k.storeKey)

	byteKey := []byte("creation-price")

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(price)

	store.Set(byteKey, bz)
}


// GetScavengesIterator gets an iterator over all scavnges in which the keys are the solutionHashes and the values are the scavenges
func (k Keeper) GetCoinsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.CoinPrefix))
}
