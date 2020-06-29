package emission

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ouroboros-crypto/node/x/coins"
	"github.com/ouroboros-crypto/node/x/emission/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	for _, record := range data.Records {
		k.SetEmission(ctx, record, coins.GetDefaultCoin())
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	var records []types.Emission

	iterator := k.GetIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		var emission types.Emission

		k.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &emission)

		records = append(records, emission)
	}

	return NewGenesisState(records)
}
