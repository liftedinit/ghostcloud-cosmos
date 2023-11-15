package keeper

import (
	"context"
	"fmt"

	"ghostcloud/x/ghostcloud/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func validateUpdateData(meta *types.Meta, payload *types.Payload) error {
	const noUpdate = "nothing to update"
	if meta == nil && payload == nil {
		return fmt.Errorf(noUpdate)
	}

	if meta.GetDescription() == "" && meta.GetDomain() == "" && payload == nil {
		return fmt.Errorf(noUpdate)
	}

	return nil
}

func validateUpdateDeploymentRequest(msg *types.MsgUpdateDeploymentRequest, params types.Params) error {
	if err := validateUpdateData(msg.Meta, msg.Payload); err != nil {
		return err
	}
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
	addr, err := sdk.AccAddressFromBech32(msg.Meta.Creator)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddr, err)
	}
	params := k.GetParams(ctx)
	err = validateUpdateDeploymentRequest(msg, params)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	meta, found := k.GetMeta(ctx, addr, msg.GetMeta().GetName())
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "unable to update a non-existing deployment")
	}

	if meta.GetCreator() != msg.Meta.GetCreator() {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "unauthorized")
	}

	if msg.Meta.Description != "" {
		meta.Description = msg.Meta.Description
	}

	if msg.Meta.Domain != "" {
		meta.Domain = msg.Meta.Domain
	}

	k.SetMeta(ctx, addr, &meta)

	if msg.GetPayload() != nil {
		dataset, err := handlePayload(msg.Payload)
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
		k.RemoveDataset(ctx, addr, msg.Meta.Name)
		k.SetDataset(ctx, addr, msg.Meta.Name, dataset)
	}

	return &types.MsgUpdateDeploymentResponse{}, nil
}
