package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ouroboros-crypto/node/x/coins"
)

// When the keeper charges posmining
type PosminingChargedHook = func(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int, coin coins.Coin)

// Adds new posmining charged hook
func (k *Keeper) AddPosminingChargedHook(hook PosminingChargedHook) {
	k.posminingChargedHooks = append(k.posminingChargedHooks, hook)
}

// Cal it after charging posmining
func (k Keeper) afterPosminingCharged(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int, coin coins.Coin) {
	for _, hook := range k.posminingChargedHooks {
		hook(ctx, addr, amt, coin)
	}
}

// Generates a hook that would be called before moving the coins from one address to another
func (k Keeper) GenerateBeforeTransferHook() func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amn sdk.Coins) {
	return func(ctx sdk.Context, sender sdk.AccAddress, receiver sdk.AccAddress, amn sdk.Coins) {
		for _, coin := range amn {
			coinRecord, err := k.CoinsKeeper.GetCoin(ctx, coin.Denom)

			if err != nil {
				panic(err)
			}

			// Charges posmining to the sender
			k.ChargePosmining(ctx, sender, coinRecord, false)

			// Saves posmined for receiver
			k.SavePosmined(ctx, receiver, coinRecord)
		}
	}
}

// Generates a hook that would be called after moving the coins from one address to another
func (k Keeper) GenerateAfterTransferHook() func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amn sdk.Coins) {
	return func(ctx sdk.Context, sender sdk.AccAddress, receiver sdk.AccAddress, amn sdk.Coins) {
		for _, coin := range amn {
			coinRecord, err := k.CoinsKeeper.GetCoin(ctx, coin.Denom)

			if err != nil {
				panic(err)
			}

			k.UpdateDailyPercent(ctx, sender, coinRecord)
			k.UpdateDailyPercent(ctx, receiver, coinRecord)
		}
	}
}

// Generates a hook that will be called when somebody changes the structure balance
func (k Keeper) GenerateStructureChangedHook() func(ctx sdk.Context, addr sdk.AccAddress, currentBalance sdk.Int, previousBalance sdk.Int, coin coins.Coin) {
	return func(ctx sdk.Context, addr sdk.AccAddress, currentBalance sdk.Int, previousBalance sdk.Int, coin coins.Coin) {
		// To avoid extra fetching of the posmining record
		currentStructureCoff := coin.GetStructureCoff(currentBalance)

		if !currentStructureCoff.Equal(coin.GetStructureCoff(previousBalance)) {
			k.SavePosmined(ctx, addr, coin)

			posmining := k.GetPosmining(ctx, addr, coin)
			posmining.StructureCoff = currentStructureCoff

			k.SetPosmining(ctx, posmining, coin)
		}
	}
}

//_________________________________________________________________________________________

// Slashing hooks
type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Create new distribution hooks
func (k Keeper) SlashingHooks() Hooks { return Hooks{k} }

// We should save posmining of every delegator before validator gets slashed
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	h.k.UpdateDelegatorsBeforeSlashing(ctx, valAddr)
}

// nolint - unused hooks
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                          {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)          {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)  {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress)                    {}
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {}
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
