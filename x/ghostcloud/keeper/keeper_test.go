package keeper_test

import (
	keepertest "ghostcloud/testutil/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSet_GetMeta(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, keepertest.NUM_DEPLOYMENT, keepertest.DATASET_SIZE)
	require.Len(t, metas, keepertest.NUM_DEPLOYMENT)

	for _, meta := range metas {
		creator, err := sdk.AccAddressFromBech32(meta.GetCreator())
		require.NoError(t, err)

		retrievedMeta, found := k.GetMeta(ctx, creator, meta.GetName())
		require.True(t, found)
		require.Equal(t, meta, &retrievedMeta)
	}
}

func TestSet_HasDeployment(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, keepertest.NUM_DEPLOYMENT, keepertest.DATASET_SIZE)
	require.Len(t, metas, keepertest.NUM_DEPLOYMENT)

	for _, meta := range metas {
		creator, err := sdk.AccAddressFromBech32(meta.GetCreator())
		require.NoError(t, err)

		has := k.HasDeployment(ctx, creator, meta.GetName())
		require.True(t, has)
	}
}

func TestSet_GetAllMeta(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, keepertest.NUM_DEPLOYMENT, keepertest.DATASET_SIZE)
	require.Len(t, metas, keepertest.NUM_DEPLOYMENT)

	all := k.GetAllMeta(ctx)
	require.Len(t, all, keepertest.NUM_DEPLOYMENT)
	require.ElementsMatch(t, metas, all)
}
