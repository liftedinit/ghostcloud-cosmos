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

	err = validateName(msg.Meta.Name)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, InvalidName, err)
	}

	err = validateDomain(msg.Meta.Domain)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, InvalidDomain, err)
	}

	err = validateDescription(msg.Meta.Description)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, InvalidDescription, err)
	}
	return nil
}
