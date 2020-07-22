package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const RestoreUnbondingConst = "restore_unbonding"

// Message for changing params
type MsgRestoreUnbonding struct {
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgRestoreUnbonding(owner sdk.AccAddress) MsgRestoreUnbonding {
	return MsgRestoreUnbonding{
		Owner: owner,
	}
}

func (msg MsgRestoreUnbonding) Route() string { return RouterKey }

func (msg MsgRestoreUnbonding) Type() string { return "restore_unbonding" }

func (msg MsgRestoreUnbonding) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner cannot be empty")
	}

	return nil
}

func (msg MsgRestoreUnbonding) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRestoreUnbonding) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
