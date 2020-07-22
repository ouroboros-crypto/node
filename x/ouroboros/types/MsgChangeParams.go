package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const ChangeParamsConst = "change_params"

// Message for changing params
type MsgChangeParams struct {
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgChangeParams(owner sdk.AccAddress) MsgChangeParams {
	return MsgChangeParams{
		Owner: owner,
	}
}

func (msg MsgChangeParams) Route() string { return RouterKey }

func (msg MsgChangeParams) Type() string { return "change_params" }

func (msg MsgChangeParams) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner cannot be empty")
	}

	return nil
}

func (msg MsgChangeParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgChangeParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
