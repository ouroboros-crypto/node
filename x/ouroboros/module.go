package ouroboros

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/ouroboros-crypto/node/x/emission"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/ouroboros-crypto/node/x/bank"
	"github.com/ouroboros-crypto/node/x/ouroboros/client/cli"
	"github.com/ouroboros-crypto/node/x/ouroboros/client/rest"
)

// Type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the ouroboros module.
type AppModuleBasic struct{}

// Name returns the ouroboros module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the ouroboros module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the ouroboros
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the ouroboros module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the ouroboros module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the ouroboros module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns no root query command for the ouroboros module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the ouroboros module.
type AppModule struct {
	AppModuleBasic

	keeper         Keeper
	paramsKeeper   params.Keeper
	bankKeeper     bank.Keeper
	emissionKeeper emission.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k Keeper, paramsKeeper params.Keeper, bankKeeper bank.Keeper, emissionKeeper emission.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		paramsKeeper:   paramsKeeper,
		bankKeeper:     bankKeeper,
		emissionKeeper: emissionKeeper,
	}
}

// Name returns the ouroboros module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the ouroboros module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the ouroboros module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the ouroboros module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper, am.paramsKeeper, am.bankKeeper, am.emissionKeeper)
}

// QuerierRoute returns the ouroboros module's querier route name.
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the ouroboros module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the ouroboros module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the ouroboros
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the ouroboros module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}

// EndBlock returns the end blocker for the ouroboros module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
