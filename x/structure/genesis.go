package structure

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ouroboros-crypto/node/x/coins"
	"github.com/ouroboros-crypto/node/x/structure/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k Keeper /* TODO: Define what keepers the module needs */, data GenesisState) {
	defaultCoin := coins.GetDefaultCoin()

	for _, record := range data.UpperStructureRecords {
		k.SetUpperStructure(ctx, record.Address, record)
	}

	for _, record := range data.StructureRecords {
		k.SetStructure(ctx, record, defaultCoin)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	var upperStructureRecords []types.UpperStructure
	var structureRecords []types.Structure

	iterator := k.GetUpperStructureIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		var upperStructure types.UpperStructure

		k.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &upperStructure)

		upperStructureRecords = append(upperStructureRecords, upperStructure)
	}

	iterator = k.GetStructureIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		var structure types.Structure

		k.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &structure)

		structureRecords = append(structureRecords, structure)
	}

	return NewGenesisState(upperStructureRecords, structureRecords)
}
