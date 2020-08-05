package ouroboros

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/ouroboros-crypto/node/x/bank"
	"github.com/ouroboros-crypto/node/x/coins"
	coinTypes "github.com/ouroboros-crypto/node/x/coins/types"
	"github.com/ouroboros-crypto/node/x/emission"
	"github.com/ouroboros-crypto/node/x/ouroboros/types"
	"github.com/ouroboros-crypto/node/x/posmining"
)


var StakingExtraAmount = sdk.NewInt(13569415899716) // 13 569 415.8997
var UnbondingExtraAmount = sdk.NewInt(455546963689) // 455 546.963689

// NewHandler creates an sdk.Handler for all the ouroboros type messages
func NewHandler(k Keeper, paramsKeeper params.Keeper, bankKeeper bank.Keeper, emissionKeeper emission.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgBurnExtraCoins:
			return handleBurnExtraCoins(ctx, k, msg, bankKeeper, emissionKeeper)
		case types.MsgUnburnExtraCoins:
			return handleUnburnExtraCoins(ctx, k, msg, bankKeeper, emissionKeeper)
		case types.MsgUpdateRegulation:
			return handleRegulation(ctx, k, msg, k.PosminingKeeper, k.CoinsKeeper)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleRegulation(ctx sdk.Context, k Keeper, msg types.MsgUpdateRegulation, posminingKeeper posmining.Keeper, coinsKeeper coins.Keeper) (*sdk.Result, error) {
	// Method updates the regulation coff based on the current coin price
	if !msg.Owner.Equals(types.GetRegulationWallet()) {
		return nil, sdkerrors.Wrapf(params.ErrSettingParameter, "only the regulation wallet can call this method")
	}

	posminingKeeper.UpdateRegulation(ctx, msg.CurrentPrice)
	coinsKeeper.UpdateCreationPrice(ctx, msg.CurrentPrice)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleBurnExtraCoins(ctx sdk.Context, k Keeper, msg types.MsgBurnExtraCoins, bankKeeper bank.Keeper, emissionKeeper emission.Keeper) (*sdk.Result, error) {
	// Method removes all the extra coins from both staking and unbonding addresses, also decreases the emission
	if !msg.Owner.Equals(types.GetGenesisWallet()) {
		return nil, sdkerrors.Wrapf(params.ErrSettingParameter, "only genesis can call this method")
	}


	stakingAddr, _ := sdk.AccAddressFromBech32("ouro1fl48vsnmsdzcv85q5d2q4z5ajdha8yu356ym48")

	unbondingAddr, _ := sdk.AccAddressFromBech32("ouro1tygms3xhhs3yv487phx3dw4a95jn7t7lq6c2rn")

	_, err := bankKeeper.SubtractCoins(ctx, stakingAddr, coinTypes.GetDefaultCoins(StakingExtraAmount))

	if err != nil {
		return nil, sdkerrors.Wrapf(params.ErrSettingParameter, err.Error())
	}

	_, err = bankKeeper.SubtractCoins(ctx, unbondingAddr, coinTypes.GetDefaultCoins(UnbondingExtraAmount))

	if err != nil {
		return nil, sdkerrors.Wrapf(params.ErrSettingParameter, err.Error())
	}

	// Decrease total emission
	emissionKeeper.Sub(ctx, StakingExtraAmount.Add(UnbondingExtraAmount), coins.GetDefaultCoin())

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleUnburnExtraCoins(ctx sdk.Context, k Keeper, msg types.MsgUnburnExtraCoins, bankKeeper bank.Keeper, emissionKeeper emission.Keeper) (*sdk.Result, error) {
	// Debug method for fixing blockchain if anything goes wrong with the method above
	if !msg.Owner.Equals(types.GetGenesisWallet()) {
		return nil, sdkerrors.Wrapf(params.ErrSettingParameter, "only genesis can call this method")
	}

	stakingAddr, _ := sdk.AccAddressFromBech32("ouro1fl48vsnmsdzcv85q5d2q4z5ajdha8yu356ym48")

	unbondingAddr, _ := sdk.AccAddressFromBech32("ouro1tygms3xhhs3yv487phx3dw4a95jn7t7lq6c2rn")

	_, err := bankKeeper.AddCoins(ctx, stakingAddr, coinTypes.GetDefaultCoins(StakingExtraAmount))

	if err != nil {
		return nil, sdkerrors.Wrapf(params.ErrSettingParameter, err.Error())
	}

	_, err = bankKeeper.AddCoins(ctx, unbondingAddr, coinTypes.GetDefaultCoins(UnbondingExtraAmount))

	if err != nil {
		return nil, sdkerrors.Wrapf(params.ErrSettingParameter, err.Error())
	}

	// Decrease total emission
	emissionKeeper.Add(ctx, StakingExtraAmount.Add(UnbondingExtraAmount), coins.GetDefaultCoin())

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}