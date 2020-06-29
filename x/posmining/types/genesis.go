package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// GenesisState - all posmining state that must be provided at genesis
type GenesisState struct {
	Correction Correction
	Records  []Posmining
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(correction Correction, records []Posmining) GenesisState {
	return GenesisState{
		Correction: correction,
		Records:  records,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	startDate := time.Date(2019, 9, 1, 0, 0, 0, 0, time.UTC)

	return GenesisState{
		Correction: Correction{StartDate: startDate, OpeningPrice: sdk.NewInt(100), CorrectionCoff: sdk.NewInt(100), PreviousCorrections: make([]PreviousCorrection, 0)},
		Records:  make([]Posmining, 0),
	}
}

// ValidateGenesis validates the posmining genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}
