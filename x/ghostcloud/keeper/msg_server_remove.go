package keeper

import (
	"context"

	"ghostcloud/x/ghostcloud/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RemoveDeployment(goCtx context.Context, msg *types.MsgRemoveDeploymentRequest) (*types.MsgRemoveDeploymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, InvalidCreatorAddr, err)
	}

	meta, found := k.GetMeta(ctx, addr, msg.Name)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "unable to remove this deployment")
	}

	if meta.GetCreator() != msg.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "unauthorized")
	}

	k.Remove(ctx, addr, msg.Name)

	return &types.MsgRemoveDeploymentResponse{}, nil
}
