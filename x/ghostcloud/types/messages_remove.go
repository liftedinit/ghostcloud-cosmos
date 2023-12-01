package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRemoveDeploymentRequest = "remove_deployment"
)

var _ sdk.Msg = &MsgRemoveDeploymentRequest{}

func (msg *MsgRemoveDeploymentRequest) Route() string {
	return RouterKey
}

func (msg *MsgRemoveDeploymentRequest) Type() string {
	return TypeMsgRemoveDeploymentRequest
}

func (msg *MsgRemoveDeploymentRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveDeploymentRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveDeploymentRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddress, err)
	}

	return nil
}
