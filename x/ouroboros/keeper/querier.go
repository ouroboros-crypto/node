package keeper

import (
	"github.com/ouroboros-crypto/node/x/coins"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ouroboros-crypto/node/x/ouroboros/types"
)

// NewQuerier creates a new querier for ouroboros clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryProfile:
			return queryProfile(ctx, path[1:], req, k, k.coinKeeper, k.CoinsKeeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown ouroboros query endpoint")
		}
	}
}

func queryProfile(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, bankKeeper bank.Keeper, coinsKeeper coins.Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(path[0])

	if err != nil {
		return []byte{}, err
	}

	coin, err := keeper.CoinsKeeper.GetCoin(ctx, path[1])

	if err != nil {
		return []byte{}, err
	}

	balance := bankKeeper.GetCoins(ctx, addr)

	posmining := keeper.posminingKeeper.GetPosminingResolve(ctx, addr, coin)

	res, codecErr := codec.MarshalJSONIndent(keeper.cdc, types.ProfileResolve{
		Owner: addr,
		Balance: balance.AmountOf(coin.Symbol),
		Posmining: posmining,
		Paramining: posmining,
		Structure: keeper.structureKeeper.GetStructure(ctx, addr, coin),
	})

	if codecErr != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
