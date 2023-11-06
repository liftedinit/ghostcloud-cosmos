package keeper_test

import (
	"testing"

	testkeeper "ghostcloud/testutil/keeper"
	"ghostcloud/x/ghostcloud/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestContentQuery(t *testing.T) {
	keeper, ctx := testkeeper.GhostcloudKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	metas, dataset := testkeeper.CreateAndSetNDeployments(ctx, keeper, testkeeper.NUM_DEPLOYMENT, testkeeper.DATASET_SIZE)
	require.Len(t, metas, testkeeper.NUM_DEPLOYMENT)
	require.Len(t, dataset, testkeeper.NUM_DEPLOYMENT)

	for i, meta := range metas {
		for _, item := range dataset[i].Items {
			response, err := keeper.Content(wctx, &types.QueryContentRequest{
				Creator: meta.GetCreator(),
				Name:    meta.GetName(),
				Path:    item.GetMeta().GetPath(),
			})
			require.NoError(t, err)
			require.NotNil(t, response)
			require.Equal(t, item.GetContent().GetContent(), response.GetContent())
		}
	}
}
