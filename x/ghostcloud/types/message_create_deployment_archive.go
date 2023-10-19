package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateDeploymentArchive = "create_deployment_archive"

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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
