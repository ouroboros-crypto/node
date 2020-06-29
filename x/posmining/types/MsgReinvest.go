package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ouroboros-crypto/node/x/coins"
)

const ReinvestConst = "reinvest"

// Реинвест пара
type MsgReinvest struct {
	Owner sdk.AccAddress `json:"owner"`
	Coin  coins.Coin     `json:"coin"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgReinvest(owner sdk.AccAddress, coin coins.Coin) MsgReinvest {
	return MsgReinvest{
		Owner: owner,
		Coin: coin,
	}
}

func (msg MsgReinvest) Route() string { return RouterKey }

func (msg MsgReinvest) Type() string { return "reinvest" }

func (msg MsgReinvest) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner cannot be empty")
	}

	return nil
}

func (msg MsgReinvest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgReinvest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
