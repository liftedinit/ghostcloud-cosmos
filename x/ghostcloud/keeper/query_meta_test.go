package keeper_test

import (
	"context"
	"ghostcloud/x/ghostcloud/keeper"
	"testing"

	"ghostcloud/testutil/sample"

	testkeeper "ghostcloud/testutil/keeper"
	"ghostcloud/x/ghostcloud/types"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMetaQuery(t *testing.T) {
	keeper, ctx := testkeeper.GhostcloudKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	metas, _ := testkeeper.CreateAndSetNDeployments(ctx, keeper, testkeeper.NUM_DEPLOYMENT, testkeeper.DATASET_SIZE)

	response, err := keeper.Metas(wctx, &types.QueryMetasRequest{})
	require.NoError(t, err)
	require.ElementsMatch(t, metas, response.Meta)
}

func createQueryMetasRequest(filters []*types.Filter, nextKey []byte, limit uint64) *types.QueryMetasRequest {
	return &types.QueryMetasRequest{
		Filters: filters,
		Pagination: &query.PageRequest{
			Key:        nextKey,
			Limit:      limit,
			CountTotal: false,
			Reverse:    false,
		},
	}
}

func assertResponse(t *testing.T, response *types.QueryMetasResponse, expectedMetas []*types.Meta) {
	require.Equal(t, len(expectedMetas), len(response.Meta))
	for i, meta := range expectedMetas {
		require.Equal(t, meta, response.Meta[i])
	}
}

func createFilters(field types.Filter_Field, operator types.Filter_Operator, value string) []*types.Filter {
	return []*types.Filter{
		{
			Field:    field,
			Operator: operator,
			Value:    value,
		},
	}
}

func queryMetas(t *testing.T, keeper *keeper.Keeper, wctx context.Context, filters []*types.Filter, nextKey []byte, limit uint64) *types.QueryMetasResponse {
	pageReq := createQueryMetasRequest(filters, nextKey, limit)
	response, err := keeper.Metas(wctx, pageReq)
	require.NoError(t, err)
	return response
}

func TestFilteredPaginatedMetaQuery(t *testing.T) {
	gcKeeper, ctx := testkeeper.GhostcloudKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	addr := sample.AccAddress()
	metas, _ := testkeeper.CreateAndSetNDeploymentsWithAddr(ctx, gcKeeper, testkeeper.NUM_DEPLOYMENT, testkeeper.DATASET_SIZE, addr)

	// Test filtering by `addr`, equal, should return the first 2 results
	filters := createFilters(types.Filter_CREATOR, types.Filter_EQUAL, addr)
	response := queryMetas(t, gcKeeper, wctx, filters, nil, 2)
	assertResponse(t, response, metas[:2])

	// Test filtering by `addr`, equal, should return the next 2 results
	response = queryMetas(t, gcKeeper, wctx, filters, response.Pagination.NextKey, 2)
	assertResponse(t, response, metas[2:4])

	// Test filtering by a random address, should return no results
	filters = createFilters(types.Filter_CREATOR, types.Filter_EQUAL, sample.AccAddress())
	response = queryMetas(t, gcKeeper, wctx, filters, nil, 2)
	require.Equal(t, 0, len(response.Meta))

	// Test filtering by `addr`, not equal, should return no results
	filters = createFilters(types.Filter_CREATOR, types.Filter_NOT_EQUAL, addr)
	response = queryMetas(t, gcKeeper, wctx, filters, nil, 2)
	require.Equal(t, 0, len(response.Meta))

	// Test filtering by a random address, not equal, should return the first 2 results
	filters = createFilters(types.Filter_CREATOR, types.Filter_NOT_EQUAL, sample.AccAddress())
	response = queryMetas(t, gcKeeper, wctx, filters, nil, 2)
	assertResponse(t, response, metas[:2])

	// Test filtering by a random address, not equal, should return the next 2 results
	response = queryMetas(t, gcKeeper, wctx, filters, response.Pagination.NextKey, 2)
	assertResponse(t, response, metas[2:4])

	// Test filtering by `addr`, contains, should return the first 2 results
	filters = createFilters(types.Filter_CREATOR, types.Filter_CONTAINS, addr[:5])
	response = queryMetas(t, gcKeeper, wctx, filters, nil, 2)
	assertResponse(t, response, metas[:2])

	// Test filtering by `addr`, contains, should return the next 2 results
	response = queryMetas(t, gcKeeper, wctx, filters, response.Pagination.NextKey, 2)
	assertResponse(t, response, metas[2:4])

	// Test filtering by a `addr`, not contains, should return no results
	filters = createFilters(types.Filter_CREATOR, types.Filter_NOT_CONTAINS, addr[:5])
	response = queryMetas(t, gcKeeper, wctx, filters, nil, 2)
	require.Equal(t, 0, len(response.Meta))
}

func TestPaginatedMetaQuery(t *testing.T) {
	keeper, ctx := testkeeper.GhostcloudKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	metas, _ := testkeeper.CreateAndSetNDeploymentsWithAddr(ctx, keeper, testkeeper.NUM_DEPLOYMENT, testkeeper.DATASET_SIZE, sample.AccAddress())

	// Test first page
	pageReq := &types.QueryMetasRequest{
		Pagination: &query.PageRequest{
			Key:        nil,
			Limit:      2,
			CountTotal: false,
			Reverse:    false,
		},
	}
	response, err := keeper.Metas(wctx, pageReq)
	require.NoError(t, err)
	require.Equal(t, 2, len(response.Meta))
	require.Equal(t, metas[0], response.Meta[0])
	require.Equal(t, metas[1], response.Meta[1])

	// Test second page
	pageReq = &types.QueryMetasRequest{
		Pagination: &query.PageRequest{
			Key:        response.Pagination.NextKey,
			Limit:      2,
			CountTotal: false,
			Reverse:    false,
		},
	}
	response, err = keeper.Metas(wctx, pageReq)
	require.NoError(t, err)
	require.Equal(t, 2, len(response.Meta))
	require.Equal(t, metas[2], response.Meta[0])
	require.Equal(t, metas[3], response.Meta[1])
}
