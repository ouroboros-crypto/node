package posmining

import (
	"fmt"
	"github.com/ouroboros-crypto/node/x/posmining/types"
	"github.com/ouroboros-crypto/node/x/update"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the posmining type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgReinvest:
			return handleReinvest(ctx, k, msg)
		case types.MsgEnable:
			return handleEnable(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// Handles reinvest
func handleReinvest(ctx sdk.Context, k Keeper, msg types.MsgReinvest) (*sdk.Result, error) {
	realCoin, err := k.CoinsKeeper.GetCoin(ctx, msg.Coin.Symbol)

	if err != nil {
		return &sdk.Result{}, err
	}

	reinvested := k.ChargePosmining(ctx, msg.Owner, realCoin, true)

	// @todo Remove that once we're updated
	if ctx.BlockHeight() > update.UpdateAfterBlock {
		k.UpdateDailyPercent(ctx, msg.Owner, realCoin)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.ReinvestConst),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, reinvested.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// Handles enable
func handleEnable(ctx sdk.Context, k Keeper, msg types.MsgEnable) (*sdk.Result, error) {
	realCoin, err := k.CoinsKeeper.GetCoin(ctx, "ouro")

	if err != nil {
		return &sdk.Result{}, err
	}

	reinvested := k.ChargePosmining(ctx, msg.Owner, realCoin, true)

	k.UpdateDailyPercent(ctx, msg.Owner, realCoin)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.ReinvestConst),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, reinvested.String()),
		),
	)

	isEnabled := k.GetPosminingEnabled(ctx, msg.Owner)

	k.SetPosminingEnabled(ctx, msg.Owner, !isEnabled)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
