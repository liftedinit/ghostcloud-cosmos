package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"ghostcloud/x/ghostcloud/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DeploymentAll(goCtx context.Context, req *types.QueryAllDeploymentRequest) (*types.QueryAllDeploymentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var deployments []types.Deployment
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	deploymentStore := prefix.NewStore(store, types.DeploymentKeyPrefix)

	pageRes, err := query.Paginate(deploymentStore, req.Pagination, func(key []byte, value []byte) error {
		var deployment types.Deployment
		if err := k.cdc.Unmarshal(value, &deployment); err != nil {
			return err
		}

		deployments = append(deployments, deployment)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDeploymentResponse{Deployment: deployments, Pagination: pageRes}, nil
}

func (k Keeper) Deployment(goCtx context.Context, req *types.QueryGetDeploymentRequest) (*types.QueryGetDeploymentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	val, found := k.GetDeployment(
		ctx,
		addr,
		req.Name,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetDeploymentResponse{Deployment: val}, nil
}

func (k Keeper) DeploymentFileContent(goCtx context.Context, req *types.QueryGetDeploymentFileContentRequest) (*types.QueryGetDeploymentFileContentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	val, found := k.GetDeploymentFileContent(
		ctx,
		addr,
		req.SiteName,
		req.FileName,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetDeploymentFileContentResponse{Content: val}, nil
}
