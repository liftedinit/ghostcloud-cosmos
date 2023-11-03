package types

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateDeploymentRequest = "create_deployment"
)

var _ sdk.Msg = &MsgCreateDeploymentRequest{}

func (msg *MsgCreateDeploymentRequest) Route() string {
	return RouterKey
}

func (msg *MsgCreateDeploymentRequest) Type() string {
	return TypeMsgCreateDeploymentRequest
}

func (msg *MsgCreateDeploymentRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Meta.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateDeploymentRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDeploymentRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Meta.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddress, err)
	}
	return nil
}
