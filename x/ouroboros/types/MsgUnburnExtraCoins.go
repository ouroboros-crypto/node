package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const UnburnExtraCoinsConst = "unburn_extra_coins"

// Message for changing params
type MsgUnburnExtraCoins struct {
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgUnburnExtraCoins(owner sdk.AccAddress) MsgUnburnExtraCoins {
	return MsgUnburnExtraCoins{
		Owner: owner,
	}
}

func (msg MsgUnburnExtraCoins) Route() string { return RouterKey }

func (msg MsgUnburnExtraCoins) Type() string { return "unburn_extra_coins" }

func (msg MsgUnburnExtraCoins) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner cannot be empty")
	}

	return nil
}

func (msg MsgUnburnExtraCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgUnburnExtraCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
