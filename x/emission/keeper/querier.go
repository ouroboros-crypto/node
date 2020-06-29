package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ouroboros-crypto/node/x/coins"
	"github.com/ouroboros-crypto/node/x/emission/types"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new querier for emission clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGetEmission:
			return getEmission(ctx, path, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown emission query endpoint")
		}
	}
}

// Returns list of the all available coins
func getEmission(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	var coin coins.Coin

	if path[1] == "ouro" {
		coin = coins.GetDefaultCoin()
	} else {
		coin = coins.Coin{
			Symbol: path[1],
		}
	}

	emission := k.GetEmission(ctx, coin)

	res, err := codec.MarshalJSONIndent(k.Cdc, types.NewQueryResGetEmission(emission, coin))

	if err != nil {
		return res, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
