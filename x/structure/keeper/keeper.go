package keeper

import (
	"fmt"
	"github.com/ouroboros-crypto/node/x/coins"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ouroboros-crypto/node/x/structure/types"
	ouroTypes "github.com/ouroboros-crypto/node/x/ouroboros/types"
)

// Keeper of the structure store
type Keeper struct {
	storeKey      sdk.StoreKey
	fastAccessKey sdk.StoreKey

	CoinsKeeper    coins.Keeper

	structureChangedHooks []StructureChangedHook

	Cdc        *codec.Codec
}

// NewKeeper creates a structure keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, fastAccessKey sdk.StoreKey, coinsKeeper coins.Keeper) Keeper {
	return Keeper{
		Cdc:           cdc,
		storeKey:      storeKey,
		fastAccessKey: fastAccessKey,
		CoinsKeeper: coinsKeeper,
		structureChangedHooks: make([]StructureChangedHook, 0),
	}
}


// Adds address to the owner's structure if he's already not a part of some structure
func (k Keeper) AddToStructure(ctx sdk.Context, owner sdk.AccAddress, address sdk.AccAddress, coinsAmount sdk.Int, coin coins.Coin) bool {
	// Already has an upper structure
	if k.HasUpperStructure(ctx, address) {
		return false
	}

	// First save pointer to the upper structure
	k.SetUpperStructure(ctx, address, types.UpperStructure{Owner: owner, Address: address})

	one := sdk.NewInt(1)

	ownerStructure := k.GetStructure(ctx, owner, coin)

	// Get upper structure and update all the info
	ownerPreviousBalance := ownerStructure.Balance

	ownerStructure.Balance = ownerStructure.Balance.Add(coinsAmount)
	ownerStructure.Followers = ownerStructure.Followers.Add(one)

	if ownerStructure.MaxLevel.LT(one) {
		ownerStructure.MaxLevel = one
	}

	k.SetStructure(ctx, ownerStructure, coin)

	// Calling hooks
	for _, hook := range k.structureChangedHooks {
		hook(ctx, owner, ownerStructure.Balance, ownerPreviousBalance, coin)
	}

	// Going up to the above structure
	nextOwner := k.GetUpperStructure(ctx, owner).Owner
	currentLevel := sdk.NewInt(2)
	maxLevel := ouroTypes.GetMaxLevel()


	for {
		// The end
		if nextOwner.Empty() || currentLevel.GT(maxLevel) {
			// Отнимаем баланс у 101 уровня, т.к. этот address уже не входит в его структуру
			if currentLevel.GT(maxLevel) && nextOwner.Empty() == false {
				topOwnerStructure := k.GetStructure(ctx, nextOwner, coin)
				previousBalance := topOwnerStructure.Balance
				topOwnerStructure.Balance = topOwnerStructure.Balance.Sub(coinsAmount)

				if topOwnerStructure.Balance.LT(sdk.NewInt(0)) {
					topOwnerStructure.Balance = sdk.NewInt(0)
				}

				k.SetStructure(ctx, topOwnerStructure, coin)

				// Calling hooks
				for _, hook := range k.structureChangedHooks {
					hook(ctx, nextOwner, topOwnerStructure.Balance, previousBalance, coin)
				}
			}

			break
		}

		// Updating the structure
		currentStructure := k.GetStructure(ctx, nextOwner, coin)

		currentStructure.Followers = currentStructure.Followers.Add(one)

		if currentStructure.MaxLevel.LT(currentLevel) {
			currentStructure.MaxLevel = currentLevel
		}

		k.SetStructure(ctx, currentStructure, coin)

		currentLevel = currentLevel.AddRaw(1)

		// Taking the next one
		nextOwner = k.GetUpperStructure(ctx, nextOwner).Owner
	}

	return true
}


// Increase structure balance by coinsAmount
func (k Keeper) IncreaseStructureBalance(ctx sdk.Context, address sdk.AccAddress, coinsAmount sdk.Int, coin coins.Coin) {
	if coinsAmount.IsZero() {
		return
	}

	nextOwner := k.GetUpperStructure(ctx, address).Owner
	currentLevel := sdk.NewInt(1)
	maxLevel := ouroTypes.GetMaxLevel()

	// Going through all the available level (until we reach 100 or genesis wallet)
	for {
		// If we reached the end
		if nextOwner.Empty() || currentLevel.GTE(maxLevel) {
			break
		}

		// Taking the current structure
		currentStructure := k.GetStructure(ctx, nextOwner, coin)

		previousBalance := currentStructure.Balance

		// Adding coins
		currentStructure.Balance = currentStructure.Balance.Add(coinsAmount)

		k.SetStructure(ctx, currentStructure, coin)

		// Calling the hooks
		for _, hook := range k.structureChangedHooks {
			hook(ctx, nextOwner, currentStructure.Balance, previousBalance, coin)
		}

		// Getting the next one
		nextOwner = k.GetUpperStructure(ctx, nextOwner).Owner

		currentLevel = currentLevel.AddRaw(1)
	}
}


// Decrease structure balance
func (k Keeper) DecreaseStructureBalance(ctx sdk.Context, address sdk.AccAddress, coinsAmount sdk.Int, coin coins.Coin) {
	if coinsAmount.IsZero() {
		return
	}

	nextOwner := k.GetUpperStructure(ctx, address).Owner
	currentLevel := sdk.NewInt(1)
	maxLevel := ouroTypes.GetMaxLevel()

	// Going through all the available level (until we reach 100 or genesis wallet)
	for {
		// If we reached the end
		if nextOwner.Empty() || currentLevel.GTE(maxLevel) {
			break
		}

		// Taking the current structure
		currentStructure := k.GetStructure(ctx, nextOwner, coin)

		previousBalance := currentStructure.Balance

		// Removing the coins
		currentStructure.Balance = currentStructure.Balance.Sub(coinsAmount)

		// This should never happen, but it's still a good practice to check those things
		if currentStructure.Balance.LT(sdk.NewInt(0)) {
			currentStructure.Balance = sdk.NewInt(0)
		}

		k.SetStructure(ctx, currentStructure, coin)

		// Calling the hooks
		for _, hook := range k.structureChangedHooks {
			hook(ctx, nextOwner, currentStructure.Balance, previousBalance, coin)
		}

		// Getting the next one
		nextOwner = k.GetUpperStructure(ctx, nextOwner).Owner

		currentLevel = currentLevel.AddRaw(1)
	}
}


// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
