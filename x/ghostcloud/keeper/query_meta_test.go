package keeper_test

import (
	testkeeper "ghostcloud/testutil/keeper"
	"ghostcloud/x/ghostcloud/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetaQuery(t *testing.T) {
	keeper, ctx := testkeeper.GhostcloudKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	metas, _ := testkeeper.CreateAndSetNDeployments(ctx, keeper, testkeeper.NUM_DEPLOYMENT, testkeeper.DATASET_SIZE)

	response, err := keeper.Metas(wctx, &types.QueryMetasRequest{})
	require.NoError(t, err)
	require.ElementsMatch(t, metas, response.Meta)
}
