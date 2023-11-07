package keeper

import (
	"context"

	"ghostcloud/x/ghostcloud/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func validateCreateDeploymentRequest(msg *types.MsgCreateDeploymentRequest) error {
	if err := validateMeta(msg.Meta); err != nil {
		return err
	}
	if err := validatePayload(msg.Payload); err != nil {
		return err
	}
	return nil
}

func (k msgServer) CreateDeployment(goCtx context.Context, msg *types.MsgCreateDeploymentRequest) (*types.MsgCreateDeploymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := validateCreateDeploymentRequest(msg); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	addr, err := sdk.AccAddressFromBech32(msg.Meta.Creator)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddr, err)
	}

	if k.HasDeployment(ctx, addr, msg.Meta.Name) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	dataset, err := handlePayload(msg.Payload)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	k.SetDeployment(
		ctx,
		msg.Meta,
		dataset,
	)
	return &types.MsgCreateDeploymentResponse{}, nil
}
