package keeper

import (
	"context"

	"ghostcloud/x/ghostcloud/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateDeployment(goCtx context.Context, msg *types.MsgCreateDeployment) (*types.MsgCreateDeploymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Check if the value already exists
	_, isFound := k.GetDeployment(
		ctx,
		addr,
		msg.Meta.Name,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var deployment = types.Deployment{
		Creator: msg.Creator,
		Meta:    msg.Meta,
		Files:   msg.Files,
	}

	k.SetDeployment(
		ctx,
		deployment,
	)
	return &types.MsgCreateDeploymentResponse{}, nil
}

func (k msgServer) UpdateDeployment(goCtx context.Context, msg *types.MsgUpdateDeployment) (*types.MsgUpdateDeploymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Check if the value exists
	valFound, isFound := k.GetDeployment(
		ctx,
		addr,
		msg.Meta.Name,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var deployment = types.Deployment{
		Creator: msg.Creator,
		Meta:    msg.Meta,
		Files:   msg.Files,
	}

	k.SetDeployment(ctx, deployment)

	return &types.MsgUpdateDeploymentResponse{}, nil
}

func (k msgServer) UpdateDeploymentMeta(goCtx context.Context, msg *types.MsgUpdateDeploymentMeta) (*types.MsgUpdateDeploymentMetaResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	deployment, isFound := k.GetDeployment(
		ctx,
		addr,
		msg.Meta.Name,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	deployment.Meta = msg.Meta

	k.SetDeployment(ctx, deployment)

	return &types.MsgUpdateDeploymentMetaResponse{}, nil
}

func (k msgServer) DeleteDeployment(goCtx context.Context, msg *types.MsgDeleteDeployment) (*types.MsgDeleteDeploymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Check if the value exists
	valFound, isFound := k.GetDeployment(
		ctx,
		addr,
		msg.Name,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveDeployment(
		ctx,
		addr,
		msg.Name,
	)

	return &types.MsgDeleteDeploymentResponse{}, nil
}
