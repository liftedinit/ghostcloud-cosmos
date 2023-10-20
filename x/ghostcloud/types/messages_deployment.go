package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateDeployment        = "create_deployment"
	TypeMsgUpdateDeployment        = "update_deployment"
	TypeMsgUpdateDeploymentMeta    = "update_deployment_meta"
	TypeMsgDeleteDeployment        = "delete_deployment"
	TypeMsgCreateDeploymentArchive = "create_deployment_archive"

	InvalidCreatorAddress = "invalid creator address (%s)"
)

var _ sdk.Msg = &MsgCreateDeployment{}

func NewMsgCreateDeployment(
	creator string,
	meta *DeploymentMeta,
	files []*File,

) *MsgCreateDeployment {
	return &MsgCreateDeployment{
		Creator: creator,
		Meta:    meta,
		Files:   files,
	}
}

func (msg *MsgCreateDeployment) Route() string {
	return RouterKey
}

func (msg *MsgCreateDeployment) Type() string {
	return TypeMsgCreateDeployment
}

func (msg *MsgCreateDeployment) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateDeployment) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDeployment) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddress, err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateDeploymentArchive{}

func NewMsgCreateDeploymentArchive(creator string, meta *DeploymentMeta, websiteArchive []byte) *MsgCreateDeploymentArchive {
	return &MsgCreateDeploymentArchive{
		Creator:        creator,
		Meta:           meta,
		WebsiteArchive: websiteArchive,
	}
}

func (msg *MsgCreateDeploymentArchive) Route() string {
	return RouterKey
}

func (msg *MsgCreateDeploymentArchive) Type() string {
	return TypeMsgCreateDeploymentArchive
}

func (msg *MsgCreateDeploymentArchive) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateDeploymentArchive) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDeploymentArchive) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddress, err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateDeployment{}

func NewMsgUpdateDeployment(
	creator string,
	meta *DeploymentMeta,
	files []*File,

) *MsgUpdateDeployment {
	return &MsgUpdateDeployment{
		Creator: creator,
		Meta:    meta,
		Files:   files,
	}
}

func (msg *MsgUpdateDeployment) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDeployment) Type() string {
	return TypeMsgUpdateDeployment
}

func (msg *MsgUpdateDeployment) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateDeployment) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDeployment) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddress, err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateDeploymentMeta{}

func NewMsgUpdateDeploymentMeta(
	creator string,
	meta *DeploymentMeta,

) *MsgUpdateDeploymentMeta {
	return &MsgUpdateDeploymentMeta{
		Creator: creator,
		Meta:    meta,
	}
}

func (msg *MsgUpdateDeploymentMeta) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDeploymentMeta) Type() string {
	return TypeMsgUpdateDeploymentMeta
}

func (msg *MsgUpdateDeploymentMeta) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateDeploymentMeta) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDeploymentMeta) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddress, err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteDeployment{}

func NewMsgDeleteDeployment(
	creator string,
	name string,

) *MsgDeleteDeployment {
	return &MsgDeleteDeployment{
		Creator: creator,
		Name:    name,
	}
}
func (msg *MsgDeleteDeployment) Route() string {
	return RouterKey
}

func (msg *MsgDeleteDeployment) Type() string {
	return TypeMsgDeleteDeployment
}

func (msg *MsgDeleteDeployment) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteDeployment) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteDeployment) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddress, err)
	}
	return nil
}
