package app

import (
	"encoding/json"
	"fmt"
	"github.com/ouroboros-crypto/node/x/coins"
	"github.com/ouroboros-crypto/node/x/emission"
	"github.com/ouroboros-crypto/node/x/ouroboros"
	"github.com/ouroboros-crypto/node/x/posmining"
	"github.com/ouroboros-crypto/node/x/structure"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"
	"io"
	"os"
	"runtime"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	xbank "github.com/ouroboros-crypto/node/x/bank"
)

const appName = "ouroboros"

var (
	// default home directories for the application CLI
	DefaultCLIHome = getCliPath()

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = getNodePath()

	// ModuleBasics The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.AppModuleBasic{},
		params.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		emission.AppModuleBasic{},
		coins.AppModuleBasic{},
		posmining.AppModuleBasic{},
		structure.AppModuleBasic{},
		ouroboros.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
	}
)

// MakeCodec creates the application codec. The codec is sealed before it is
// returned.
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc.Seal()
}

// NewApp extended ABCI application
type NewApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     xbank.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	distrKeeper    distr.Keeper
	govKeeper      gov.Keeper
	supplyKeeper   supply.Keeper
	paramsKeeper   params.Keeper

	// here we go
	emissionKeeper  emission.Keeper
	coinsKeeper     coins.Keeper
	posminingKeeper posmining.Keeper
	structureKeeper structure.Keeper
	ouroKeeper      ouroboros.Keeper

	// Module Manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

// verify app interface at compile time
var _ simapp.App = (*NewApp)(nil)

// NewnodeApp is a constructor function for nodeApp
func NewInitApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *NewApp {
	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	// TODO: Add the keys that module requires
	keys := sdk.NewKVStoreKeys(bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, gov.StoreKey, distr.StoreKey, slashing.StoreKey, params.StoreKey, emission.StoreKey,
		coins.StoreKey, posmining.StoreKey, structure.StoreKey, structure.FastAccessKey)

	tKeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	// Here you initialize your application with the store keys it requires
	var app = &NewApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		subspaces:      make(map[string]params.Subspace),
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	// Set specific supspaces
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		app.subspaces[auth.ModuleName],
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = xbank.NewKeeper(
		app.accountKeeper,
		app.subspaces[bank.ModuleName],
		app.ModuleAccountAddrs(),
	)

	// The SupplyKeeper collects transaction fees and renders them to the fee distribution module
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		maccPerms,
	)

	// The staking keeper
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		keys[staking.StoreKey],
		app.supplyKeeper,
		app.subspaces[staking.ModuleName],
	)

	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		keys[distr.StoreKey],
		app.subspaces[distr.ModuleName],
		&stakingKeeper,
		app.supplyKeeper,
		auth.FeeCollectorName,
		app.ModuleAccountAddrs(),
	)

	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		keys[slashing.StoreKey],
		&stakingKeeper,
		app.subspaces[slashing.ModuleName],
	)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper))

	app.govKeeper = gov.NewKeeper(app.cdc, keys[gov.StoreKey], app.subspaces[gov.ModuleName],
		app.supplyKeeper, &stakingKeeper, govRouter)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			app.distrKeeper.Hooks(),
			app.slashingKeeper.Hooks(),
			app.emissionKeeper.SlashingHooks(),
		),
	)

	// Since we cannot do that during the keepers creation
	app.bankKeeper.StakingKeeper = app.stakingKeeper

	// The emission keeper
	app.emissionKeeper = emission.NewKeeper(
		app.cdc,
		keys[emission.StoreKey],
		app.stakingKeeper,
	)

	// The coins keeper
	app.coinsKeeper = coins.NewKeeper(
		app.cdc,
		keys[coins.StoreKey],
		app.bankKeeper,
	)

	// Keeper that handles all the structures related stuff
	app.structureKeeper = structure.NewKeeper(
		app.cdc,
		keys[structure.StoreKey],
		keys[structure.FastAccessKey],
		app.coinsKeeper,
	)

	// The paramining keeper
	app.posminingKeeper = posmining.NewKeeper(
		app.cdc,
		keys[posmining.StoreKey],
		app.bankKeeper,
		app.stakingKeeper,
		app.emissionKeeper,
		app.coinsKeeper,
	)

	// Helping keeper that's mostly being used for the API calls
	app.ouroKeeper = ouroboros.NewKeeper(
		app.cdc,
		app.accountKeeper,
		app.bankKeeper,
		app.structureKeeper,
		app.posminingKeeper,
		app.emissionKeeper,
		app.supplyKeeper,
		app.slashingKeeper,
		app.coinsKeeper,
	)

	// Register hooks
	app.coinsKeeper.AddCoinCreatedHook(app.emissionKeeper.GenerateCoinCreatedHook())
	app.coinsKeeper.AddCoinCreatedHook(app.structureKeeper.GenerateCoinCreatedHook())

	app.structureKeeper.AddStructureChangedHook(app.posminingKeeper.GenerateStructureChangedHook())

	app.posminingKeeper.AddPosminingChargedHook(app.structureKeeper.GeneratePosminingChargedHook())
	app.posminingKeeper.AddPosminingChargedHook(app.emissionKeeper.GeneratePosminingChargedHook())

	app.bankKeeper.AddBeforeHook(app.posminingKeeper.GenerateBeforeTransferHook())
	app.bankKeeper.AddAfterHook(app.posminingKeeper.GenerateAfterTransferHook())
	app.bankKeeper.AddAfterHook(app.structureKeeper.GenerateAfterTransferHook())

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		emission.NewAppModule(app.emissionKeeper),
		coins.NewAppModule(app.coinsKeeper, app.bankKeeper),
		posmining.NewAppModule(app.posminingKeeper),
		structure.NewAppModule(app.structureKeeper),
		ouroboros.NewAppModule(app.ouroKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		// TODO: Add your module(s)
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),

	)
	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.

	app.mm.SetOrderBeginBlockers(distr.ModuleName, slashing.ModuleName)
	app.mm.SetOrderEndBlockers(staking.ModuleName, gov.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		distr.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		emission.ModuleName,
		coins.ModuleName,
		posmining.ModuleName,
		structure.ModuleName,
		ouroboros.ModuleName,
		// TODO: Add your module(s)
		supply.ModuleName,
		genutil.ModuleName,
	)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

// InitChainer application update at chain initialization
func (app *NewApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState

	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

// BeginBlocker application updates every begin block
func (app *NewApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *NewApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	// Check if we should change regulation based on the price every 100 blocks
	if ctx.BlockHeight() % 100 == 0 {
		fmt.Println(ctx.BlockHeight())
		/*client := http.Client{
			Timeout: 5 * time.Second, // 5 seconds timeout
		}

		resp, err := client.Get("https://api.ouroboros-crypto.com/correction/price")

		if err != nil {
			return app.mm.EndBlock(ctx, req)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		// Some problems with parsing the body
		if err != nil {
			return app.mm.EndBlock(ctx, req)
		}*/

		price, isOk := sdk.NewIntFromString("1")

		if !isOk || (ctx.BlockHeight() >= 272100 && ctx.BlockHeight() <= 272105) {
			return app.mm.EndBlock(ctx, req)
		}

		// Update posmining correction and also the creation price
		app.posminingKeeper.UpdateRegulation(ctx, price)
		app.coinsKeeper.UpdateCreationPrice(ctx, price)
	}

	return app.mm.EndBlock(ctx, req)
}

// LoadHeight loads a particular height
func (app *NewApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *NewApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// Codec returns the application's sealed codec.
func (app *NewApp) Codec() *codec.Codec {
	return app.cdc
}

// SimulationManager implements the SimulationApp interface
func (app *NewApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// GetMaccPerms returns a mapping of the application's module account permissions.
func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}

func getCliPath() string {
	if runtime.GOOS == "windows" {
		return os.ExpandEnv("$UserProfile/.ouroboroscli")
	}

	return os.ExpandEnv("$HOME/.ouroboroscli")
}

func getNodePath() string {
	if runtime.GOOS == "windows" {
		return os.ExpandEnv("$UserProfile/.ouroborosd")
	}

	return os.ExpandEnv("$HOME/.ouroborosd")
}