package keeper

import (
	"context"
	"fmt"

	"ghostcloud/x/ghostcloud/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) Content(goCtx context.Context, req *types.QueryContentRequest) (*types.QueryContentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(req.GetCreator())
	if err != nil {
		return nil, fmt.Errorf("invalid creator address: %v", err)
	}

	content, found := k.GetItemContent(ctx, creator, req.GetName(), req.GetPath())
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	response := &types.QueryContentResponse{
		Content: content.GetContent(),
	}

	return response, nil
}
