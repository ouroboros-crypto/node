package types

// GenesisState - all structure state that must be provided at genesis
type GenesisState struct {
	UpperStructureRecords []UpperStructure
	StructureRecords      []Structure
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(upperStructureRecords []UpperStructure, structureRecords []Structure) GenesisState {
	return GenesisState{
		UpperStructureRecords: upperStructureRecords,
		StructureRecords: structureRecords,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		UpperStructureRecords: make([]UpperStructure, 0),
		StructureRecords:      make([]Structure, 0),
	}
}

// ValidateGenesis validates the structure genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}
