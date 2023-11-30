package keeper

import (
	"context"

	"ghostcloud/x/ghostcloud/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func validateUpdateDeploymentRequest(msg *types.MsgUpdateDeploymentRequest, params types.Params) error {
	if err := validateMeta(msg.Meta, params); err != nil {
		return err
	}
	if msg.GetPayload() != nil {
		if err := validatePayload(msg.Payload, params); err != nil {
			return err
		}
	}
	return nil
}

func (k msgServer) UpdateDeployment(goCtx context.Context, msg *types.MsgUpdateDeploymentRequest) (*types.MsgUpdateDeploymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	err := validateUpdateDeploymentRequest(msg, params)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	addr, err := sdk.AccAddressFromBech32(msg.Meta.Creator)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddr, err)
	}

	meta, found := k.GetMeta(ctx, addr, msg.GetMeta().GetName())
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "unable to update a non-existing deployment")
	}

	// The following should never happen since the store key uses the creator address
	if meta.GetCreator() != msg.Meta.GetCreator() {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "unauthorized")
	}

	meta.Description = msg.Meta.Description
	meta.Domain = msg.Meta.Domain

	k.SetMeta(ctx, addr, &meta)

	if msg.GetPayload() != nil {
		dataset, err := HandlePayload(msg.Payload)
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
		k.RemoveDataset(ctx, addr, msg.Meta.Name)
		k.SetDataset(ctx, addr, msg.Meta.Name, dataset)
	}

	return &types.MsgUpdateDeploymentResponse{}, nil
}
