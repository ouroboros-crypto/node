package coins

import (
	"github.com/ouroboros-crypto/node/x/coins/keeper"
	"github.com/ouroboros-crypto/node/x/coins/types"
)

const (
	// TODO: define constants that you would like exposed from your module

	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	//QueryParams       = types.QueryParams
	QuerierRoute      = types.QuerierRoute
)

var (
	// functions aliases
	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	// TODO: Fill out function aliases

	// variable aliases
	ModuleCdc     = types.ModuleCdc
	GetDefaultCoin     = types.GetDefaultCoin
	GetCoinStub     = types.GetCoinStub
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params
	Coin       = types.Coin

	// TODO: Fill out module types
	MsgCreateCoin = types.MsgCreateCoin
)
