package keeper

import (
	"context"
	"strings"

	"ghostcloud/x/ghostcloud/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/types/query"

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

	var metas []*types.Meta
	store := k.getDeploymentMetaStore(ctx)
	pageRes, err := query.FilteredPaginate(store, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var meta types.Meta
		if err := k.cdc.Unmarshal(value, &meta); err != nil {
			return false, err
		}
		if req.Filters != nil {
			if metaPassesAllFilters(&meta, req.Filters) {
				if accumulate {
					metas = append(metas, &meta)
				}
				return true, nil
			}
		} else {
			if accumulate {
				metas = append(metas, &meta)
			}
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "pagination error %v", err)
	}

	return &types.QueryMetasResponse{Meta: metas, Pagination: pageRes}, nil
}
