package types

// GenesisState - all emission state that must be provided at genesis
type GenesisState struct {
	Records []Emission
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(records []Emission) GenesisState {
	return GenesisState{
		Records:records,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Records: make([]Emission, 0),
	}
}

// ValidateGenesis validates the emission genesis parameters
func ValidateGenesis(data GenesisState) error {

	return nil
}
