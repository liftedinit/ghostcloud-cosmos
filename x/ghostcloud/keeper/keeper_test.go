package keeper_test

import (
	keepertest "ghostcloud/testutil/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSet_GetMeta(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, 10, 5)
	require.Len(t, metas, 10)

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
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, 10, 5)
	require.Len(t, metas, 10)

	for _, meta := range metas {
		creator, err := sdk.AccAddressFromBech32(meta.GetCreator())
		require.NoError(t, err)

		has := k.HasDeployment(ctx, creator, meta.GetName())
		require.True(t, has)
	}
}

func TestSet_GetAllMeta(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	metas, _ := keepertest.CreateAndSetNDeployments(ctx, k, 10, 5)
	require.Len(t, metas, 10)

	all := k.GetAllMeta(ctx)
	require.Len(t, all, 10)
	require.ElementsMatch(t, metas, all)
}
