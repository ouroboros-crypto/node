package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/ouroboros-crypto/node/x/coins"
	"github.com/ouroboros-crypto/node/x/emission"
	"github.com/ouroboros-crypto/node/x/posmining"
	"github.com/ouroboros-crypto/node/x/structure"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ouroboros-crypto/node/x/ouroboros/types"
)

// Keeper of the ouroboros store
type Keeper struct {
	accountKeeper   auth.AccountKeeper
	coinKeeper      bank.Keeper
	CoinsKeeper     coins.Keeper
	structureKeeper structure.Keeper
	PosminingKeeper posmining.Keeper
	emissionKeeper  emission.Keeper
	supplyKeeper    supply.Keeper
	slashingKeeper  slashing.Keeper

	cdc *codec.Codec
}

// NewKeeper creates a ouroboros keeper
func NewKeeper(cdc *codec.Codec, accountKeeper auth.AccountKeeper, coinKeeper bank.Keeper, structureKeeper structure.Keeper, posminingKeeper posmining.Keeper, emissionKeeper emission.Keeper, supplyKeeper supply.Keeper, slashingKeeper slashing.Keeper, coinsKeeper coins.Keeper) Keeper {
	return Keeper{
		cdc:             cdc,
		accountKeeper:   accountKeeper,
		coinKeeper:      coinKeeper,
		structureKeeper: structureKeeper,
		PosminingKeeper: posminingKeeper,
		emissionKeeper:  emissionKeeper,
		supplyKeeper:    supplyKeeper,
		slashingKeeper:  slashingKeeper,
		CoinsKeeper:     coinsKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
