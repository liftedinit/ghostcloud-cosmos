package keeper_test

import (
	"strconv"
	"testing"

	keepertest "ghostcloud/testutil/keeper"
	"ghostcloud/testutil/nullify"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestDeploymentGet(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	items := keepertest.CreateNDeployment(k, ctx, 10)
	for _, item := range items {
		addr, err := sdk.AccAddressFromBech32(item.Meta.Creator)
		require.NoError(t, err)
		rst, found := k.GetDeployment(ctx,
			addr,
			item.Meta.Name,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestDeploymentRemove(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	items := keepertest.CreateNDeployment(k, ctx, 10)
	for _, item := range items {
		addr, err := sdk.AccAddressFromBech32(item.Meta.Creator)
		require.NoError(t, err)
		k.RemoveDeployment(ctx,
			addr,
			item.Meta.Name,
		)
		_, found := k.GetDeployment(ctx,
			addr,
			item.Meta.Name,
		)
		require.False(t, found)
	}
}

func TestDeploymentGetAll(t *testing.T) {
	k, ctx := keepertest.GhostcloudKeeper(t)
	items := keepertest.CreateNDeployment(k, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(k.GetAllDeployment(ctx)),
	)
}
