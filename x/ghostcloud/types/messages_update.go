package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateDeploymentRequest = "update_deployment"
)

var _ sdk.Msg = &MsgUpdateDeploymentRequest{}

func (msg *MsgUpdateDeploymentRequest) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDeploymentRequest) Type() string {
	return TypeMsgUpdateDeploymentRequest
}

func (msg *MsgUpdateDeploymentRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Meta.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateDeploymentRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDeploymentRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Meta.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddress, err)
	}
	return nil
}
