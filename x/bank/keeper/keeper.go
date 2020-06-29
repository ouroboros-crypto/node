package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkbank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/ouroboros-crypto/node/x/coins"
)

type Keeper struct {
	sdkbank.BaseKeeper

	ak         auth.AccountKeeper
	StakingKeeper    staking.Keeper
	paramSpace params.Subspace

	beforeTransferHooks []CoinsTransferHook
	afterTransferHooks []CoinsTransferHook
}

func NewKeeper(
	ak auth.AccountKeeper, paramSpace params.Subspace, blacklistedAddrs map[string]bool,
) Keeper {
	return Keeper{
		BaseKeeper: sdkbank.NewBaseKeeper(ak, paramSpace, blacklistedAddrs),
		ak:             ak,
		paramSpace:     paramSpace,
		beforeTransferHooks: make([]CoinsTransferHook, 0),
		afterTransferHooks: make([]CoinsTransferHook, 0),
	}
}


// Returns the balance that should be used during calculations of posmining
func (k Keeper) GetPosminableBalance(ctx sdk.Context, addr sdk.AccAddress, coin coins.Coin) sdk.Int {
	return k.GetCoins(ctx, addr).Add(k.GetStackedCoins(ctx, addr)).AmountOf(coin.Symbol)
}

// Returns both stacked and unbounding coins
func (k Keeper) GetStackedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coin {
	result := sdk.NewInt(0)

	// First let's get through the stakes
	stakes := k.StakingKeeper.GetAllDelegatorDelegations(ctx, addr)

	for _, value := range stakes {
		result = result.Add(value.GetShares().TruncateInt())
	}

	// Then let's get through the unbounding coins
	unbounding := k.StakingKeeper.GetAllUnbondingDelegations(ctx, addr)

	for _, value := range unbounding {
		for _, entry := range value.Entries {
			result = result.Add(entry.Balance)
		}
	}

	return sdk.NewCoin("ouro", result)
}


// SendCoins moves coins from one account to another
func (keeper Keeper) SendCoins(
	ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins,
) error {
	keeper.beforeCoinsTransfer(ctx, fromAddr, toAddr, amt)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdkbank.EventTypeTransfer,
			sdk.NewAttribute(sdkbank.AttributeKeyRecipient, toAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amt.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdkbank.AttributeKeySender, fromAddr.String()),
		),
	})

	_, err := keeper.SubtractCoins(ctx, fromAddr, amt)

	if err != nil {
		return err
	}

	_, err = keeper.AddCoins(ctx, toAddr, amt)

	if err != nil {
		return err
	}

	keeper.afterCoinsTransfer(ctx, fromAddr, toAddr, amt)

	return nil
}