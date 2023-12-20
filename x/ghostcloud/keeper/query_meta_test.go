package keeper_test

import (
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
