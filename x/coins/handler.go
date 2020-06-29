package coins

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ouroboros-crypto/node/x/coins/types"
)

// NewHandler creates an sdk.Handler for all the coins type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgCreateCoin:
			return handleMsgCreateCoin(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgCreateCoin handles coins creation
func handleMsgCreateCoin(ctx sdk.Context, k Keeper, msg MsgCreateCoin) (*sdk.Result, error) {
	var coin = types.Coin{
		Creator:          msg.Creator,
		Name:             msg.Name,
		Symbol:           msg.Symbol,
		Description:      msg.Description,
		Emission:         msg.Emission,
		PosminingEnabled: msg.PosminingEnabled,
		PosminingBalance: msg.PosminingBalance,
		PosminingStructure: msg.StructurePosmining,
	}

	_, err := k.GetCoin(ctx, coin.Symbol)

	if err == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Coin with that symbol already exists")
	}

	_, sdkError := k.BankKeeper.AddCoins(ctx, coin.Creator, sdk.NewCoins(sdk.NewCoin(coin.Symbol, coin.Emission)))

	if sdkError != nil {
		return nil, sdkError
	}

	creationPrice, err := k.GetCreationPrice(ctx)

	if err != nil {
		return nil, err
	}

	paymentForCreation := types.GetDefaultCoins(creationPrice.Price)

	_, err = k.BankKeeper.SubtractCoins(ctx, coin.Creator, paymentForCreation)

	if err != nil {
		return nil, err
	}

	_, err = k.BankKeeper.AddCoins(ctx, types.GetGenesisWallet(), paymentForCreation)

	if err != nil {
		return nil, err
	}

	k.SetCoin(ctx, coin)

	k.AfterCoinCreated(ctx, coin, paymentForCreation.AmountOf("ouro"))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeCreateCoin),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
			sdk.NewAttribute(types.AttributeName, msg.Name),
			sdk.NewAttribute(types.AttributeSymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeDescription, msg.Description),
			sdk.NewAttribute(types.AttributeEmission, msg.Emission.String()),

		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
