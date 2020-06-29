package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ouroboros-crypto/node/x/coins"
)

func (k Keeper) GeneratePosminingChargedHook() func(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int, coin coins.Coin) {
	return func(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int, coin coins.Coin) {
		k.Add(ctx, amt, coin)
	}
}

// When somebody creates a coin, we should setup its emission
func (k Keeper) GenerateCoinCreatedHook() func(ctx sdk.Context, coin coins.Coin, paymentForCreation sdk.Int) {
	return func(ctx sdk.Context, coin coins.Coin, paymentForCreation sdk.Int) {
		k.Add(ctx, coin.Emission, coin)
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
	h.k.UpdateBeforeSlashing(ctx, valAddr, fraction)
}


// nolint - unused hooks
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                         {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)         {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)       {}
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {}
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) { }
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {}
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {}
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) { }

