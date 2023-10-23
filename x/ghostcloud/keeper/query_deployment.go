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

const InvalidRequest = "invalid request"

func (k Keeper) DeploymentMetaAll(goCtx context.Context, req *types.QueryAllDeploymentMetaRequest) (*types.QueryAllDeploymentMetaResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, InvalidRequest)
	}

	var metas []types.DeploymentMeta
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	deploymentStore := prefix.NewStore(store, types.DeploymentMetaKeyPrefix)

	pageRes, err := query.Paginate(deploymentStore, req.Pagination, func(key []byte, value []byte) error {
		var meta types.DeploymentMeta
		if err := k.cdc.Unmarshal(value, &meta); err != nil {
			return err
		}

		metas = append(metas, meta)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDeploymentMetaResponse{DeploymentMeta: metas, Pagination: pageRes}, nil
}

func (k Keeper) Deployment(goCtx context.Context, req *types.QueryGetDeploymentRequest) (*types.QueryGetDeploymentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, InvalidRequest)
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	deploymentMeta, found := k.GetDeploymentMeta(
		ctx,
		addr,
		req.Name,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	filesMeta, found := k.GetDeploymentFileMeta(
		ctx,
		addr,
		req.Name,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetDeploymentResponse{DeploymentMeta: deploymentMeta, FileMeta: filesMeta}, nil
}

func (k Keeper) DeploymentFileNames(goCtx context.Context, req *types.QueryDeploymentFileNamesRequest) (*types.QueryDeploymentFileNamesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, InvalidRequest)
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(req.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	fileMeta, found := k.GetDeploymentFileMeta(
		ctx,
		addr,
		req.SiteName,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryDeploymentFileNamesResponse{Meta: fileMeta}, nil
}

func (k Keeper) DeploymentFileContent(goCtx context.Context, req *types.QueryGetDeploymentFileContentRequest) (*types.FileContent, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, InvalidRequest)
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

	return &val, nil
}
