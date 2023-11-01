package keeper

import (
	"context"

	"ghostcloud/x/ghostcloud/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) Metas(goCtx context.Context, req *types.QueryMetasRequest) (*types.QueryMetasResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	metas := k.GetAllMeta(ctx)

	return &types.QueryMetasResponse{Meta: metas}, nil
}
