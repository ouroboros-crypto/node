package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const UpdateRegulation = "update_regulation"

// Message for changing params
type MsgUpdateRegulation struct {
	Owner sdk.AccAddress `json:"owner"`
	CurrentPrice sdk.Int `json:"current_price"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgUpdateRegulation(owner sdk.AccAddress, price sdk.Int) MsgUpdateRegulation {
	return MsgUpdateRegulation{
		Owner: owner,
		CurrentPrice: price,
	}
}

func (msg MsgUpdateRegulation) Route() string { return RouterKey }

func (msg MsgUpdateRegulation) Type() string { return "update_regulation" }

func (msg MsgUpdateRegulation) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner cannot be empty")
	}

	return nil
}

func (msg MsgUpdateRegulation) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgUpdateRegulation) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
