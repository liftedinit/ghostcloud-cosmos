package keeper

import (
	"context"
	"strings"

	"ghostcloud/x/ghostcloud/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func metaPassesAllFilters(item *types.Meta, filters []*types.Filter) bool {
	for _, filter := range filters {
		if !metaPassesFilter(item, filter) {
			return false
		}
	}
	return true
}

func metaPassesFilter(item *types.Meta, filter *types.Filter) bool {
	var itemValue string

	switch filter.Field {
	case types.Filter_CREATOR:
		itemValue = item.Creator
	default:
		return false
	}

	switch filter.Operator {
	case types.Filter_EQUAL:
		return itemValue == filter.Value
	case types.Filter_NOT_EQUAL:
		return itemValue != filter.Value
	case types.Filter_CONTAINS:
		return strings.Contains(itemValue, filter.Value)
	case types.Filter_NOT_CONTAINS:
		return !strings.Contains(itemValue, filter.Value)
	default:
		return false
	}
}

func (k Keeper) Metas(goCtx context.Context, req *types.QueryMetasRequest) (*types.QueryMetasResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	metas := k.GetAllMeta(ctx)

	var filteredMetas []*types.Meta
	// Filter the result if necessary
	if req.Filters != nil {
		for _, meta := range metas {
			if metaPassesAllFilters(meta, req.Filters) {
				filteredMetas = append(filteredMetas, meta)
			}
		}
	} else {
		filteredMetas = metas
	}

	return &types.QueryMetasResponse{Meta: filteredMetas}, nil
}
