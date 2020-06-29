package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ouroboros-crypto/node/x/coins"
	"github.com/ouroboros-crypto/node/x/structure/types"
)

// Returns "upper structure" that's just a pointer to the next structure above
func (k Keeper) GetUpperStructure(ctx sdk.Context, address sdk.AccAddress) types.UpperStructure {
	store := ctx.KVStore(k.fastAccessKey)

	if !store.Has(address.Bytes()) {
		return types.NewUpperStructure(address)
	}

	var upperStructure types.UpperStructure

	k.Cdc.MustUnmarshalBinaryBare(store.Get(address.Bytes()), &upperStructure)

	return upperStructure
}

// Saves pointer to the upper account in the sturcture
func (k Keeper) SetUpperStructure(ctx sdk.Context, address sdk.AccAddress, upperStructure types.UpperStructure) {
	store := ctx.KVStore(k.fastAccessKey)

	store.Set(address.Bytes(), k.Cdc.MustMarshalBinaryBare(upperStructure))
}

// Checks if the user in any structure yet
func (k Keeper) HasUpperStructure(ctx sdk.Context, address sdk.AccAddress) bool {
	return !k.GetUpperStructure(ctx, address).Owner.Empty()
}

// Iterator
func (k Keeper) GetUpperStructureIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.fastAccessKey)

	return sdk.KVStorePrefixIterator(store, nil)
}

// Returns the structure record by its owner
func (k Keeper) GetStructure(ctx sdk.Context, owner sdk.AccAddress, coin coins.Coin) types.Structure {
	store := ctx.KVStore(k.storeKey)

	structureKey := []byte("structure:" + owner.String())

	if !store.Has(structureKey) {
		return types.NewStructure(owner)
	}

	var structure types.Structure

	k.Cdc.MustUnmarshalBinaryBare(store.Get(structureKey), &structure)

	// Because of the custom coins, we should fetch the balance from another key
	realBalance := sdk.NewInt(0)

	balanceKey := []byte(coin.Symbol + owner.String())

	if store.Has(balanceKey) {
		k.Cdc.MustUnmarshalBinaryBare(store.Get([]byte(coin.Symbol + owner.String())), &realBalance)
	}

	structure.Balance = realBalance

	return structure
}

// Saves the sturcture record
func (k Keeper) SetStructure(ctx sdk.Context, structure types.Structure, coin coins.Coin) {
	store := ctx.KVStore(k.storeKey)

	structureKey := []byte("structure:" + structure.Owner.String())

	store.Set(structureKey, k.Cdc.MustMarshalBinaryBare(structure))

	// Because of the custom coins, we should fetch the balance from another key
	store.Set([]byte(coin.Symbol + structure.Owner.String()), k.Cdc.MustMarshalBinaryBare(structure.Balance))
}


// Iterator
func (k Keeper) GetStructureIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)

	return sdk.KVStorePrefixIterator(store, []byte("structure:"))
}
