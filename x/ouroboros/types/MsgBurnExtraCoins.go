package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const BurnExtraCoinsConst = "burn_extra_coins"

// Message for changing params
type MsgBurnExtraCoins struct {
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgBurnExtraCoins(owner sdk.AccAddress) MsgBurnExtraCoins {
	return MsgBurnExtraCoins{
		Owner: owner,
	}
}

func (msg MsgBurnExtraCoins) Route() string { return RouterKey }

func (msg MsgBurnExtraCoins) Type() string { return "burn_extra_coins" }

func (msg MsgBurnExtraCoins) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner cannot be empty")
	}

	return nil
}

func (msg MsgBurnExtraCoins) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnExtraCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
