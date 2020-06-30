package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ouroboros-crypto/node/x/bank"
	"github.com/ouroboros-crypto/node/x/coins"
	"github.com/ouroboros-crypto/node/x/emission"
	"github.com/ouroboros-crypto/node/x/posmining/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	BankKeeper     bank.Keeper
	CoinsKeeper    coins.Keeper
	stakingKeeper  staking.Keeper
	emissionKeeper emission.Keeper

	Cdc *codec.Codec // The wire codec for binary encoding/decoding.

	// Hooks
	posminingChargedHooks []PosminingChargedHook
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, coinKeeper bank.Keeper, stakingKeeper staking.Keeper, emissionKeeper emission.Keeper, coinsKeeper coins.Keeper) Keeper {
	return Keeper{
		storeKey:              storeKey,
		BankKeeper:            coinKeeper,
		stakingKeeper:         stakingKeeper,
		emissionKeeper:        emissionKeeper,
		CoinsKeeper:           coinsKeeper,
		Cdc:                   cdc,
		posminingChargedHooks: make([]PosminingChargedHook, 0),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Charges the posmining and resets the lastCharged field
func (k Keeper) ChargePosmining(ctx sdk.Context, addr sdk.AccAddress, coin coins.Coin, isReinvest bool) sdk.Int {
	balance := k.BankKeeper.GetPosminableBalance(ctx, addr, coin)

	posmining := k.GetPosmining(ctx, addr, coin)

	posMined := k.CalculatePosmined(ctx, posmining, coin, balance)

	posmining.Paramined = sdk.NewInt(0)

	// Since reinvest doesn't reset "last transaction" field
	if !isReinvest {
		posmining.LastTransaction = ctx.BlockHeader().Time
	}

	posmining.LastCharged = ctx.BlockHeader().Time

	k.SetPosmining(ctx, posmining, coin)

	// If we charged at least 0.000001
	if posMined.IsPositive() {
		_, err := k.BankKeeper.AddCoins(ctx, addr, sdk.NewCoins(sdk.NewCoin(coin.Symbol, posMined)))

		if err != nil {
			panic(err)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePosminingCharged,
				sdk.NewAttribute(sdk.AttributeKeySender, addr.String()),
				sdk.NewAttribute(types.AttributeParamined, posMined.String()),
			),
		)

		k.afterPosminingCharged(ctx, addr, posMined, coin)
	}

	return posMined
}

// Saves the posmined coins without charging it
func (k Keeper) SavePosmined(ctx sdk.Context, addr sdk.AccAddress, coin coins.Coin) sdk.Int {
	balance := k.BankKeeper.GetPosminableBalance(ctx, addr, coin)

	posmining := k.GetPosmining(ctx, addr, coin)

	posMined := k.CalculatePosmined(ctx, posmining, coin, balance)

	posmining.Paramined = posMined
	posmining.LastCharged = ctx.BlockHeader().Time

	k.SetPosmining(ctx, posmining, coin)

	return posMined
}

// Updates daily percent based on the posminable balance
func (k Keeper) UpdateDailyPercent(ctx sdk.Context, addr sdk.AccAddress, coin coins.Coin) {
	balance := k.BankKeeper.GetPosminableBalance(ctx, addr, coin)

	posmining := k.GetPosmining(ctx, addr, coin)

	newDailyPercent := coin.GetDailyPercent(balance)

	if !posmining.DailyPercent.Equal(newDailyPercent) {
		posmining.DailyPercent = newDailyPercent

		k.SetPosmining(ctx, posmining, coin)
	}
}

// Fetches the current price and updates posmining regulation
func (k Keeper) UpdateRegulation(ctx sdk.Context, currentPrice sdk.Int) {
	regulation := k.GetCorrection(ctx)

	coff := regulation.GetCoff(currentPrice)

	// If the coff should be changed and since the latest update passed at least types.CorrectionUpdatePeriod hours
	if !regulation.CorrectionCoff.Equal(coff) && ctx.BlockTime().Sub(regulation.StartDate).Hours() >= types.CorrectionUpdatePeriod {
		regulation.Update(ctx.BlockTime(), currentPrice, coff)

		k.SetCorrection(ctx, regulation)
	}
}

// We need to save delegators posmining before they get slashed
func (k Keeper) UpdateDelegatorsBeforeSlashing(ctx sdk.Context, valAddr sdk.ValAddress) {
	delegations := k.stakingKeeper.GetValidatorDelegations(ctx, valAddr)

	defaultCoin := coins.GetDefaultCoin()

	for _, delegation := range delegations {
		k.SavePosmined(ctx, delegation.DelegatorAddress, defaultCoin)
	}
}

// Resolves posmining, so we can get that data via API
func (k Keeper) GetPosminingResolve(ctx sdk.Context, owner sdk.AccAddress, coin coins.Coin) types.PosminingResolve {
	balance := k.BankKeeper.GetPosminableBalance(ctx, owner, coin)

	posmining := k.GetPosmining(ctx, owner, coin)

	posminingGroup := k.GetPosminingGroup(ctx, posmining, coin, balance)

	var currentPeriod types.PosminingPeriod

	if len(posminingGroup.Periods) > 0 {
		currentPeriod = posminingGroup.Periods[len(posminingGroup.Periods)-1]
	} else {
		currentPeriod = types.PosminingPeriod{SavingCoff: sdk.NewInt(0), CorrectionCoff: sdk.NewInt(0)}
	}

	if coin.PosminingThreshold.IsPositive() && balance.Add(posminingGroup.Paramined).GTE(coin.PosminingThreshold) {
		posminingGroup.Paramined = sdk.NewInt(0)

		if balance.GTE(coin.PosminingThreshold) {
			posmining.Paramined = sdk.NewInt(0)
		}
	}

	return types.PosminingResolve{
		Coin:           coin.Symbol,
		Posmining:      posmining,
		Paramining:     posmining,
		SavingsCoff:    currentPeriod.SavingCoff,
		CorrectionCoff: currentPeriod.CorrectionCoff,
		Posmined:       posmining.Paramined.Add(posminingGroup.Paramined),
		Paramined:      posmining.Paramined.Add(posminingGroup.Paramined),
		CoinsPerTime:   types.NewCoinsPerTime(balance, posmining.DailyPercent, posmining.StructureCoff, currentPeriod.SavingCoff, currentPeriod.CorrectionCoff),
	}
}
