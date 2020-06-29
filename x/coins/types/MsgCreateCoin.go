package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateCoin{}

type MsgCreateCoin struct {
	Creator     sdk.AccAddress `json:"creator" yaml:"creator"`         // address of the coin creator
	Name        string         `json:"string" yaml:"string"`           // name of the coin
	Symbol      string         `json:"symbol" yaml:"symbol"`           // identifier of the coin
	Description string         `json:"description" yaml:"description"` // description of the coin
	Emission    sdk.Int        `json:"emission" yaml:"emission"`       // initial emission of the coin

	PosminingEnabled bool                   `json:"posmining_enabled" yaml:"posmining_enabled"` // if posmining should be enabled
	PosminingBalance []CoinBalancePosmining `json:"posmining_balance" yaml:"posmining_balance"` // all the daily percent conditions
	StructurePosmining []CoinStructurePosmining `json:"posmining_balance" yaml:"posmining_balance"` // all the daily percent conditions
}

// NewMsgCreateCoin creates a new MsgCreateCoin instance
func NewMsgCreateCoin(creator sdk.AccAddress, name, symbol, description string, emission sdk.Int, posminingEnabled bool, posminingBalance []CoinBalancePosmining, posminingStructure []CoinStructurePosmining) MsgCreateCoin {
	return MsgCreateCoin{
		Creator:     creator,
		Name:        name,
		Symbol:      symbol,
		Description: description,
		Emission:    emission,
		PosminingEnabled: posminingEnabled,
		PosminingBalance: posminingBalance,
		StructurePosmining: posminingStructure,
	}
}

const CreateCoinConst = "CreateCoin"

// nolint
func (msg MsgCreateCoin) Route() string { return RouterKey }
func (msg MsgCreateCoin) Type() string  { return CreateCoinConst }

func (msg MsgCreateCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Creator)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgCreateCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)

	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgCreateCoin) ValidateBasic() error {
	if msg.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator can't be empty")
	}

	if msg.PosminingEnabled {
		i := 1

		for i < len(msg.PosminingBalance) - 1 {
			if !msg.PosminingBalance[i].FromAmount.Equal(msg.PosminingBalance[i-1].ToAmount) {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Posmining coffs should be the same")
			}

			i -= 1
		}

		i = 1

		for i < len(msg.StructurePosmining) - 1 {
			if !msg.StructurePosmining[i].FromAmount.Equal(msg.StructurePosmining[i-1].ToAmount) {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Structure coffs should be the same")
			}

			i -= 1
		}
	}

	return nil
}
