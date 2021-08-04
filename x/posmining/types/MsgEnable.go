package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)


// Enables or disables posmining
type MsgEnable struct {
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgEnable(owner sdk.AccAddress) MsgEnable {
	return MsgEnable{
		Owner: owner,
	}
}

func (msg MsgEnable) Route() string { return RouterKey }

func (msg MsgEnable) Type() string { return "enable" }

func (msg MsgEnable) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner cannot be empty")
	}

	return nil
}

func (msg MsgEnable) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgEnable) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
