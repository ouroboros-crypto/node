package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ouroboros-crypto/node/x/structure/types"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new querier for structure clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGet:
			return getStructure(ctx, path, k)
		case types.QueryUpper:
			return getUpperStructure(ctx, path, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown structure query endpoint")
		}
	}
}


// Returns the structure record
func getStructure(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(path[1])

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	coin, err := k.CoinsKeeper.GetCoin(ctx, path[2])

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	res, err := codec.MarshalJSONIndent(k.Cdc, k.GetStructure(ctx, addr, coin))

	if err != nil {
		return res, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}



// Returns the structure record
func getUpperStructure(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(path[1])

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	res, err := codec.MarshalJSONIndent(k.Cdc, k.GetUpperStructure(ctx, addr))

	if err != nil {
		return res, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}