// nolint
package keeper // noalias

// DONTCOVER

import (
	"encoding/hex"
	"github.com/ouroboros-crypto/node/x/coins"
	"github.com/ouroboros-crypto/node/x/emission"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	xbank "github.com/ouroboros-crypto/node/x/bank"
)

// TODO move to common testing framework
func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd
}

func makeTestCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func createTestInput(t *testing.T, isCheckTx bool) (sdk.Context, auth.AccountKeeper, Keeper, staking.Keeper, types.SupplyKeeper) {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyPosmining := sdk.NewKVStoreKey(types.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyPosmining, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	require.Nil(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "ouroboros"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := makeTestCodec()

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		types.ModuleName:          nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}

	// create module accounts
	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	govAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Burner)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[govAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := xbank.NewKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), blacklistedAddrs)
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)

	sk := staking.NewKeeper(cdc, keyStaking, supplyKeeper, pk.Subspace(staking.DefaultParamspace))
	sk.SetParams(ctx, staking.DefaultParams())

	keeper := NewKeeper(
		cdc, keyPosmining, bankKeeper, sk, emission.Keeper{}, coins.Keeper{},
	)

	return ctx, accountKeeper, keeper, sk, supplyKeeper
}
