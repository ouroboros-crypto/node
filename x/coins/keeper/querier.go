package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ouroboros-crypto/node/x/coins/types"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new querier for coins clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryListCoins:
			return listCoins(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown coins query endpoint")
		}
	}
}

// Returns list of the all available coins
func listCoins(ctx sdk.Context, k Keeper) ([]byte, error) {
	var coinsList types.QueryResCoins

	iterator := k.GetCoinsIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		var coin types.Coin

		_ = k.cdc.UnmarshalBinaryLengthPrefixed(iterator.Value(), &coin)

		coinsList = append(coinsList, coin)
	}

	res, err := codec.MarshalJSONIndent(k.cdc, coinsList)

	if err != nil {
		return res, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
